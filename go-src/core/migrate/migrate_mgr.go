package migrate

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"xr-game-server/core/cfg"
)

type OneModel struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type MoreModel struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	IsDeleted bool      `gorm:"default:0" json:"-"`
	DeletedAt time.Time `json:"-"`
}

var (
	migrateMu  sync.Mutex
	migrateDB  *gorm.DB
	migrateSQL *sql.DB
)

func mysqlDSN() string {
	link := cfg.DefaultDbCfg.Link
	if !strings.HasPrefix(link, "mysql:") {
		panic("database.default.link must use mysql driver")
	}
	dsn := strings.TrimPrefix(link, "mysql:")
	if extra := cfg.DefaultDbCfg.Extra; extra != "" {
		sep := "?"
		if strings.Contains(dsn, "?") {
			sep = "&"
		}
		dsn += sep + extra
	}
	return dsn
}

func openMigrateDB() *gorm.DB {
	cfg.EnsureErrorLogger()

	migrateMu.Lock()
	defer migrateMu.Unlock()
	if migrateDB != nil {
		return migrateDB
	}

	db, err := gorm.Open(mysql.Open(mysqlDSN()), &gorm.Config{})
	if err != nil {
		panic("表结构无法同步,连不上数据库")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("表结构无法同步,获取连接池失败")
	}
	applyMigratePool(sqlDB)

	migrateDB = db
	migrateSQL = sqlDB
	return migrateDB
}

func applyMigratePool(sqlDB *sql.DB) {
	dbCfg := cfg.DefaultDbCfg
	maxOpen := dbCfg.MaxOpen
	if maxOpen <= 0 {
		maxOpen = 2
	}
	maxIdle := dbCfg.MaxIdle
	if maxIdle <= 0 {
		maxIdle = 1
	}
	if maxIdle > maxOpen {
		maxIdle = maxOpen
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	if dbCfg.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(dbCfg.MaxLifetime)
	}
}

// AutoMigrate 代码启动,自动同步表结构,表结构保持驼峰
func AutoMigrate(dst ...interface{}) {
	db := openMigrateDB()
	for _, m := range dst {
		if err := db.AutoMigrate(m); err != nil {
			panic("表结构无法同步: " + err.Error())
		}
	}
}

// Close 关闭迁移专用连接池；entity.Init 完成后调用,避免长期占用数据库连接。
func Close() {
	migrateMu.Lock()
	defer migrateMu.Unlock()
	if migrateSQL != nil {
		_ = migrateSQL.Close()
	}
	migrateDB = nil
	migrateSQL = nil
}
