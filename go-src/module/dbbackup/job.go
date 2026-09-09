package dbbackup

import (
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/dao/cfgdao"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/gameevent"
	"xr-game-server/module/upload"
)

var dbNameRe = regexp.MustCompile(`^[A-Za-z0-9_]{1,64}$`)

func stagingDir() string {
	root := strings.TrimSpace(upload.GetStoragePath())
	if root == "" {
		return filepath.Join(os.TempDir(), "sara-db-backup")
	}
	return filepath.Join(filepath.Dir(root), "staging", "db-backup")
}

func findBinary(names ...string) (string, error) {
	for _, name := range names {
		path, err := exec.LookPath(name)
		if err == nil && path != "" {
			return path, nil
		}
	}
	return "", fmt.Errorf("未找到命令: %s", strings.Join(names, "/"))
}

func runBackupOnce(row *sysentity.DbBackupCfg, force bool) (objectKey string, size int64, err error) {
	if row == nil {
		return "", 0, fmt.Errorf("cfg nil")
	}
	if !force && !row.Enabled {
		return "", 0, nil
	}
	if !upload.IsS3Enabled() {
		err = fmt.Errorf("云桶未开启,无法上传备份")
		markFailure(row, err)
		return "", 0, err
	}
	if !tryBeginJob() {
		err = fmt.Errorf("已有备份/还原任务进行中")
		return "", 0, err
	}
	defer endJob()

	ensureStoragePrefix(row)
	parts, err := currentMysqlParts()
	if err != nil {
		markFailure(row, err)
		return "", 0, err
	}
	dumpBin, err := findBinary("mysqldump", "mariadb-dump")
	if err != nil {
		markFailure(row, err)
		return "", 0, err
	}

	dir := stagingDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		markFailure(row, err)
		return "", 0, err
	}
	fileName := time.Now().UTC().Format("20060102_150405") + ".sql.gz"
	localPath := filepath.Join(dir, fileName)
	defer os.Remove(localPath)

	host, port := splitHostPort(parts.HostPort)
	args := []string{
		"--single-transaction",
		"--routines",
		"--triggers",
		"--events",
		"--hex-blob",
		"-h", host,
		"-P", port,
		"-u", parts.User,
		parts.Database,
	}
	cmd := exec.Command(dumpBin, args...)
	cmd.Env = append(os.Environ(), "MYSQL_PWD="+parts.Password)

	outFile, err := os.Create(localPath)
	if err != nil {
		markFailure(row, err)
		return "", 0, err
	}
	gzWriter := gzip.NewWriter(outFile)
	cmd.Stdout = gzWriter
	var stderr strings.Builder
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	_ = gzWriter.Close()
	_ = outFile.Close()
	if runErr != nil {
		err = fmt.Errorf("mysqldump failed: %v, stderr=%s", runErr, strings.TrimSpace(stderr.String()))
		markFailure(row, err)
		return "", 0, err
	}
	fi, err := os.Stat(localPath)
	if err != nil {
		markFailure(row, err)
		return "", 0, err
	}
	size = fi.Size()
	if size == 0 {
		err = fmt.Errorf("备份文件为空")
		markFailure(row, err)
		return "", 0, err
	}

	objectKey = strings.Trim(row.StoragePrefix, "/") + "/" + fileName
	if err := upload.PutStoredFileToS3(objectKey, localPath); err != nil {
		markFailure(row, err)
		return "", 0, err
	}
	markSuccess(row, objectKey, size, parts.Database)
	_ = cleanupExpired(row)
	g.Log().Infof(gctx.New(), "db backup ok key=%s size=%d db=%s", objectKey, size, parts.Database)
	return objectKey, size, nil
}

func splitHostPort(hostPort string) (host, port string) {
	hostPort = strings.TrimSpace(hostPort)
	if hostPort == "" {
		return "127.0.0.1", "3306"
	}
	if i := strings.LastIndex(hostPort, ":"); i >= 0 {
		return hostPort[:i], hostPort[i+1:]
	}
	return hostPort, "3306"
}

func cleanupExpired(row *sysentity.DbBackupCfg) error {
	if row == nil || !upload.IsS3Enabled() {
		return nil
	}
	retainDays := row.RetainDays
	if retainDays <= 0 {
		retainDays = 1
	}
	cutoff := time.Now().UTC().Add(-time.Duration(retainDays) * 24 * time.Hour)
	objs, err := upload.ListStoredObjectsByPrefix(row.StoragePrefix)
	if err != nil {
		return err
	}
	for _, obj := range objs {
		if obj.StoredName == "" {
			continue
		}
		mod := obj.LastModified.UTC()
		if mod.IsZero() || !mod.Before(cutoff) {
			continue
		}
		upload.DeleteStoredFileFromS3(obj.StoredName)
		g.Log().Infof(gctx.New(), "db backup expired deleted key=%s modified=%v", obj.StoredName, mod)
	}
	return nil
}

func restoreFromObject(row *sysentity.DbBackupCfg, objectKey, targetDB string) error {
	objectKey = strings.Trim(strings.ReplaceAll(objectKey, "\\", "/"), "/")
	targetDB = strings.TrimSpace(targetDB)
	if objectKey == "" || targetDB == "" {
		return fmt.Errorf("参数无效")
	}
	if !dbNameRe.MatchString(targetDB) {
		return fmt.Errorf("数据库名仅允许字母数字下划线,最长64")
	}
	prefix := strings.Trim(row.StoragePrefix, "/") + "/"
	if !strings.HasPrefix(objectKey, prefix) {
		return fmt.Errorf("备份文件不在当前备份目录下")
	}
	if !strings.HasSuffix(strings.ToLower(objectKey), ".sql.gz") {
		return fmt.Errorf("仅支持 .sql.gz 备份文件")
	}
	if !upload.IsS3Enabled() {
		return fmt.Errorf("云桶未开启")
	}
	if !tryBeginJob() {
		return fmt.Errorf("已有备份/还原任务进行中")
	}
	defer endJob()

	parts, err := currentMysqlParts()
	if err != nil {
		return err
	}
	mysqlBin, err := findBinary("mysql", "mariadb")
	if err != nil {
		return err
	}

	dir := stagingDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	localPath := filepath.Join(dir, "restore_"+filepath.Base(objectKey))
	defer os.Remove(localPath)
	if err := upload.DownloadStoredFileFromS3(objectKey, localPath); err != nil {
		return fmt.Errorf("下载备份失败: %w", err)
	}

	host, port := splitHostPort(parts.HostPort)
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", targetDB)
	createCmd := exec.Command(mysqlBin, "-h", host, "-P", port, "-u", parts.User, "-e", createSQL)
	createCmd.Env = append(os.Environ(), "MYSQL_PWD="+parts.Password)
	if out, err := createCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("创建数据库失败: %v, out=%s", err, strings.TrimSpace(string(out)))
	}

	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("解压备份失败: %w", err)
	}
	defer gz.Close()

	importCmd := exec.Command(mysqlBin, "-h", host, "-P", port, "-u", parts.User, targetDB)
	importCmd.Env = append(os.Environ(), "MYSQL_PWD="+parts.Password)
	importCmd.Stdin = gz
	var stderr strings.Builder
	importCmd.Stderr = &stderr
	if err := importCmd.Run(); err != nil {
		return fmt.Errorf("还原失败: %v, stderr=%s", err, strings.TrimSpace(stderr.String()))
	}
	g.Log().Infof(gctx.New(), "db restore ok key=%s target=%s", objectKey, targetDB)
	return nil
}

func onDayBackup(_ any) {
	row := cfgdao.GetDbBackupCfgCached()
	if row == nil || !row.Enabled {
		g.Log().Info(gctx.New(), "db backup skipped: disabled or no cfg")
		return
	}
	if _, _, err := runBackupOnce(row, false); err != nil {
		g.Log().Warningf(gctx.New(), "db backup day job failed: %v", err)
	}
}

func subscribeDayEvent() {
	event.Sub(gameevent.DayEvent, onDayBackup)
}
