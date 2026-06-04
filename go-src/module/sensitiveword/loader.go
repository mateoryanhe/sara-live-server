package sensitiveword

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"

	"xr-game-server/constants/lang"
	"xr-game-server/core/cfg"
)

const downloadTimeout = 30 * time.Second

var supportedLangs = []lang.Lang{lang.LangEN}

func langFileName(l lang.Lang) string {
	return string(l) + ".txt"
}

func langFilePath(localDir string, l lang.Lang) string {
	return filepath.Join(localDir, langFileName(l))
}

func loadAll(ctx context.Context) error {
	swCfg := cfg.GetSensitiveWordCfg()
	if swCfg == nil {
		return nil
	}
	localDir := swCfg.LocalDir
	if err := gfile.Mkdir(localDir); err != nil {
		return fmt.Errorf("mkdir sensitive word dir: %w", err)
	}

	for _, l := range supportedLangs {
		key := string(l)
		url := strings.TrimSpace(swCfg.DownloadUrls[key])
		path := langFilePath(localDir, l)

		words, err := readWordsFromFile(path)
		if shouldDownloadWordFile(l, url, words, err) {
			if dlErr := downloadToFile(ctx, url, path); dlErr != nil {
				g.Log().Warningf(ctx, "download sensitive words lang=%s err=%v, use local file", l, dlErr)
			} else {
				g.Log().Infof(ctx, "downloaded sensitive words lang=%s path=%s", l, path)
				words, err = readWordsFromFile(path)
			}
		}
		if err != nil {
			if os.IsNotExist(err) {
				g.Log().Warningf(ctx, "sensitive word file missing lang=%s path=%s", l, path)
				setWords(l, nil)
				continue
			}
			return fmt.Errorf("read sensitive words lang=%s: %w", l, err)
		}
		prepared := prepareWords(l, words)
		setWords(l, prepared)
		g.Log().Infof(ctx, "loaded sensitive words lang=%s count=%d path=%s", l, len(prepared), path)
	}
	return nil
}

// shouldDownloadWordFile 本地无文件或无有效词条时下载一次;已有词库则跳过
func shouldDownloadWordFile(l lang.Lang, url string, words []string, readErr error) bool {
	if url == "" {
		return false
	}
	if readErr != nil {
		return os.IsNotExist(readErr)
	}
	return len(prepareWords(l, words)) == 0
}

func downloadToFile(ctx context.Context, url, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: downloadTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status %d", resp.StatusCode)
	}
	tmpPath := destPath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, resp.Body)
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(tmpPath)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmpPath)
		return closeErr
	}
	return os.Rename(tmpPath, destPath)
}

func readWordsFromFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	words := make([]string, 0, 256)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

// prepareWords 加载后预处理(英文转小写、去重)
func prepareWords(l lang.Lang, words []string) []string {
	if len(words) == 0 {
		return words
	}
	seen := make(map[string]struct{}, len(words))
	out := make([]string, 0, len(words))
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		if l == lang.LangEN {
			w = strings.ToLower(w)
		}
		if _, ok := seen[w]; ok {
			continue
		}
		seen[w] = struct{}{}
		out = append(out, w)
	}
	return out
}

// Refresh 从本地词库文件重新加载(不重复下载)
func Refresh() error {
	return loadAll(gctx.New())
}

// LocalDir 返回当前配置的本地词库目录
func LocalDir() string {
	swCfg := cfg.GetSensitiveWordCfg()
	if swCfg == nil || swCfg.LocalDir == "" {
		return "data/sensitive_words"
	}
	return swCfg.LocalDir
}
