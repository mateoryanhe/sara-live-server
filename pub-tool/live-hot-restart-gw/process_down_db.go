package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var processDownDB *sql.DB

func initProcessDownDB() {
	dsn := strings.TrimSpace(cfg.MysqlDSN)
	if dsn == "" {
		log.Println("process down db: skipped (mysqlDSN empty)")
		return
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("process down db: open failed: %v", err)
		return
	}
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.Ping(); err != nil {
		log.Printf("process down db: ping failed: %v", err)
		_ = db.Close()
		return
	}
	if err := ensureProcessDownTable(db); err != nil {
		log.Printf("process down db: ensure table failed: %v", err)
		_ = db.Close()
		return
	}
	processDownDB = db
	log.Println("process down db: ready")
}

func ensureProcessDownTable(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS process_down_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  host VARCHAR(128) NOT NULL DEFAULT '',
  app_name VARCHAR(64) NOT NULL DEFAULT '',
  missing_since DATETIME(3) NOT NULL,
  detected_at DATETIME(3) NOT NULL,
  cold_start_ok TINYINT(1) NOT NULL DEFAULT 0,
  new_pid INT NOT NULL DEFAULT 0,
  detail VARCHAR(1024) NOT NULL DEFAULT '',
  created_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_process_down_host (host),
  KEY idx_process_down_detected (detected_at),
  KEY idx_process_down_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	return err
}

// recordProcessDownEvent 进程连续缺失达阈值并完成冷启尝试后写库.
func recordProcessDownEvent(appName string, missingSince, detectedAt time.Time, coldOK bool, newPID int, detail string) {
	if processDownDB == nil {
		return
	}
	host, _ := os.Hostname()
	host = strings.TrimSpace(host)
	appName = strings.TrimSpace(appName)
	if appName == "" {
		appName = "xr-game-server"
	}
	detail = strings.TrimSpace(detail)
	if len(detail) > 1000 {
		detail = detail[:1000]
	}
	now := time.Now()
	_, err := processDownDB.Exec(`
INSERT INTO process_down_logs
  (host, app_name, missing_since, detected_at, cold_start_ok, new_pid, detail, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		host, appName, missingSince, detectedAt, coldOK, newPID, detail, now)
	if err != nil {
		log.Printf("process down db: insert failed: %v", err)
		return
	}
	log.Printf("process down db: recorded host=%s ok=%v pid=%d detail=%s", host, coldOK, newPID, detail)
}

func formatDownDetail(coldOK bool, result coldStartResult, coldErr error) string {
	if coldErr != nil {
		return fmt.Sprintf("cold start failed: %v", coldErr)
	}
	if result.Output == "already running" {
		return fmt.Sprintf("cold start skipped, already running pid=%d", result.PID)
	}
	if coldOK {
		return fmt.Sprintf("cold start ok pid=%d", result.PID)
	}
	return "cold start finished"
}
