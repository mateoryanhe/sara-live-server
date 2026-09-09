package dbbackup

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/core/cfg"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/dbbackupdto"
	sysentity "xr-game-server/entity/sys"
	"xr-game-server/module/upload"
)

var (
	jobMu      sync.Mutex
	jobRunning bool
)

var mysqlLinkRe = regexp.MustCompile(`(?i)^mysql:([^:]+):([^@]*)@tcp\(([^)]+)\)/([^?\s]+)`)

type mysqlConnParts struct {
	User     string
	Password string
	HostPort string
	Database string
}

func parseMysqlLink(link string) (*mysqlConnParts, error) {
	link = strings.TrimSpace(link)
	m := mysqlLinkRe.FindStringSubmatch(link)
	if len(m) != 5 {
		return nil, fmt.Errorf("invalid database link")
	}
	return &mysqlConnParts{
		User:     m[1],
		Password: m[2],
		HostPort: m[3],
		Database: strings.TrimSpace(m[4]),
	}, nil
}

func currentMysqlParts() (*mysqlConnParts, error) {
	return parseMysqlLink(cfg.DefaultDbCfg.Link)
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return gtime.New(t.UTC()).Format("Y-m-d H:i:s")
}

func ensureStoragePrefix(row *sysentity.DbBackupCfg) {
	if row == nil {
		return
	}
	if strings.TrimSpace(row.StoragePrefix) == "" {
		row.StoragePrefix = strings.Trim(guid.S(), "/") + "/db_backup"
	}
	if row.RetainDays <= 0 {
		row.RetainDays = 1
	}
}

func toCfgItem(row *sysentity.DbBackupCfg) *dbbackupdto.DbBackupCfgItem {
	srcDB := ""
	if parts, err := currentMysqlParts(); err == nil && parts != nil {
		srcDB = parts.Database
	}
	item := &dbbackupdto.DbBackupCfgItem{
		ID:             "0",
		Enabled:        false,
		RetainDays:     1,
		S3Enabled:      upload.IsS3Enabled(),
		SourceDatabase: srcDB,
	}
	if row == nil || row.ID == 0 {
		return item
	}
	item.ID = strconv.FormatUint(row.ID, 10)
	item.Enabled = row.Enabled
	item.StoragePrefix = row.StoragePrefix
	item.RetainDays = row.RetainDays
	if item.RetainDays <= 0 {
		item.RetainDays = 1
	}
	item.LastSuccessAt = formatCfgTime(row.LastSuccessAt)
	item.LastError = row.LastError
	item.LastObjectKey = row.LastObjectKey
	item.LastFileSize = row.LastFileSize
	item.LastTargetHint = row.LastTargetHint
	item.CreatedAt = formatCfgTime(row.CreatedAt)
	item.UpdatedAt = formatCfgTime(row.UpdatedAt)
	return item
}

func persistCfg(row *sysentity.DbBackupCfg) error {
	if err := cfgdao.SaveDbBackupCfg(row); err != nil {
		return err
	}
	cfgdao.ReloadDbBackupCfgCache()
	return nil
}

func tryBeginJob() bool {
	jobMu.Lock()
	defer jobMu.Unlock()
	if jobRunning {
		return false
	}
	jobRunning = true
	return true
}

func endJob() {
	jobMu.Lock()
	jobRunning = false
	jobMu.Unlock()
}

func markFailure(row *sysentity.DbBackupCfg, err error) {
	if row == nil || err == nil {
		return
	}
	msg := err.Error()
	if len(msg) > 500 {
		msg = msg[:500]
	}
	row.LastError = msg
	_ = persistCfg(row)
}

func markSuccess(row *sysentity.DbBackupCfg, objectKey string, size int64, dbName string) {
	if row == nil {
		return
	}
	row.LastSuccessAt = time.Now().UTC()
	row.LastError = ""
	row.LastObjectKey = objectKey
	row.LastFileSize = size
	row.LastTargetHint = dbName
	_ = persistCfg(row)
}
