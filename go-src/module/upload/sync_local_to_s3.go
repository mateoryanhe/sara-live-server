package upload

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dto/uploaddto"
	"xr-game-server/errercode"
)

// 跳过临时导出目录,不刷到云桶
var syncSkipTopDirs = map[string]struct{}{
	StoreCatExport: {},
}

type localToS3SyncState struct {
	mu         sync.Mutex
	running    atomic.Bool
	root       string
	keyPrefix  string
	startedAt  time.Time
	finishedAt time.Time
	total      int
	done       int
	success    int
	failed     int
	skipped    int
	lastKey    string
	lastError  string
}

var localToS3Sync = &localToS3SyncState{}

// StartSyncLocalStorageToS3 扫描 GetStoragePath 下文件并 PutObject;异步执行
func StartSyncLocalStorageToS3() (*uploaddto.SyncLocalStorageToS3Res, error) {
	if !IsS3Enabled() {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	snap := getResourceCfgCache()
	if !s3Ready(snap) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	root := strings.TrimSpace(GetStoragePath())
	info, err := os.Stat(root)
	if err != nil || info == nil || !info.IsDir() {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	st := localToS3Sync
	if !st.running.CompareAndSwap(false, true) {
		return &uploaddto.SyncLocalStorageToS3Res{
			Started:             false,
			AlreadyRunning:      true,
			Message:             "already running",
			SyncLocalToS3Status: st.snapshot(),
		}, nil
	}

	st.mu.Lock()
	st.root = root
	st.keyPrefix = strings.Trim(GetS3KeyPrefix(), "/")
	st.startedAt = time.Now()
	st.finishedAt = time.Time{}
	st.total = 0
	st.done = 0
	st.success = 0
	st.failed = 0
	st.skipped = 0
	st.lastKey = ""
	st.lastError = ""
	st.mu.Unlock()

	go runSyncLocalStorageToS3(root)

	return &uploaddto.SyncLocalStorageToS3Res{
		Started:             true,
		AlreadyRunning:      false,
		Message:             "started",
		SyncLocalToS3Status: st.snapshot(),
	}, nil
}

// GetSyncLocalStorageToS3Status 查询刷桶进度
func GetSyncLocalStorageToS3Status() *uploaddto.GetSyncLocalStorageToS3StatusRes {
	return &uploaddto.GetSyncLocalStorageToS3StatusRes{
		SyncLocalToS3Status: localToS3Sync.snapshot(),
	}
}

func (s *localToS3SyncState) snapshot() *uploaddto.SyncLocalToS3Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := &uploaddto.SyncLocalToS3Status{
		Running:   s.running.Load(),
		Root:      s.root,
		KeyPrefix: s.keyPrefix,
		Total:     s.total,
		Done:      s.done,
		Success:   s.success,
		Failed:    s.failed,
		Skipped:   s.skipped,
		LastKey:   s.lastKey,
		LastError: s.lastError,
	}
	if !s.startedAt.IsZero() {
		out.StartedAt = s.startedAt.Format("2006-01-02 15:04:05")
	}
	if !s.finishedAt.IsZero() {
		out.FinishedAt = s.finishedAt.Format("2006-01-02 15:04:05")
	}
	return out
}

func runSyncLocalStorageToS3(root string) {
	ctx := gctx.New()
	defer func() {
		if r := recover(); r != nil {
			g.Log().Errorf(ctx, "syncLocalToS3 panic: %v", r)
			localToS3Sync.mu.Lock()
			localToS3Sync.lastError = "panic"
			localToS3Sync.mu.Unlock()
		}
		localToS3Sync.mu.Lock()
		localToS3Sync.finishedAt = time.Now()
		localToS3Sync.mu.Unlock()
		localToS3Sync.running.Store(false)
	}()

	files, err := listLocalFilesForS3Sync(root)
	if err != nil {
		g.Log().Warningf(ctx, "syncLocalToS3 list failed root=%s err=%v", root, err)
		localToS3Sync.mu.Lock()
		localToS3Sync.lastError = err.Error()
		localToS3Sync.mu.Unlock()
		return
	}

	localToS3Sync.mu.Lock()
	localToS3Sync.total = len(files)
	localToS3Sync.mu.Unlock()

	g.Log().Infof(ctx, "syncLocalToS3 start root=%s files=%d prefix=%s", root, len(files), GetS3KeyPrefix())

	const workers = 4
	type job struct {
		absPath    string
		storedName string
	}
	ch := make(chan job, workers*2)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range ch {
				putErr := putLocalFileToS3(j.storedName, j.absPath)
				key := objectKeyFromStored(j.storedName)
				localToS3Sync.mu.Lock()
				localToS3Sync.done++
				localToS3Sync.lastKey = key
				if putErr != nil {
					localToS3Sync.failed++
					localToS3Sync.lastError = putErr.Error()
					g.Log().Warningf(ctx, "syncLocalToS3 put failed key=%s err=%v", key, putErr)
				} else {
					localToS3Sync.success++
				}
				localToS3Sync.mu.Unlock()
			}
		}()
	}

	for _, f := range files {
		ch <- job{absPath: f.absPath, storedName: f.storedName}
	}
	close(ch)
	wg.Wait()

	snap := localToS3Sync.snapshot()
	g.Log().Infof(ctx, "syncLocalToS3 done root=%s total=%d success=%d failed=%d",
		root, snap.Total, snap.Success, snap.Failed)
}

type syncFileItem struct {
	absPath    string
	storedName string
}

func listLocalFilesForS3Sync(root string) ([]syncFileItem, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	var out []syncFileItem
	err = filepath.WalkDir(rootAbs, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d == nil {
			return nil
		}
		name := d.Name()
		if name == "." || name == ".." {
			return nil
		}
		if strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(rootAbs, path)
		if err != nil {
			return err
		}
		relSlash := filepath.ToSlash(rel)
		if relSlash == "." {
			return nil
		}
		top := strings.SplitN(relSlash, "/", 2)[0]
		if d.IsDir() {
			if _, skip := syncSkipTopDirs[top]; skip && relSlash == top {
				localToS3Sync.mu.Lock()
				localToS3Sync.skipped++
				localToS3Sync.mu.Unlock()
				return filepath.SkipDir
			}
			return nil
		}
		info, err := d.Info()
		if err != nil || info == nil || !info.Mode().IsRegular() {
			return nil
		}
		stored := sanitizeStoredRelativePath(relSlash)
		if stored == "" {
			localToS3Sync.mu.Lock()
			localToS3Sync.skipped++
			localToS3Sync.mu.Unlock()
			return nil
		}
		out = append(out, syncFileItem{absPath: path, storedName: stored})
		return nil
	})
	return out, err
}
