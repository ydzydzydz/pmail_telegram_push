package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/ydzydzydz/pmail_telegram_push/config"
	"github.com/ydzydzydz/pmail_telegram_push/dao"
	"github.com/ydzydzydz/pmail_telegram_push/model"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

// DataSource 数据库数据源
type DataSource struct {
	db *xorm.Engine
}

// NewDataSource 创建数据库数据源
func NewDataSource(cfg *config.Config) (*DataSource, error) {
	db, err := getDB(cfg)
	if err != nil {
		return nil, err
	}
	return &DataSource{db: db}, nil
}

// getDB 获取数据库连接
func getDB(cfg *config.Config) (*xorm.Engine, error) {
	var db *xorm.Engine
	var err error

	// 根据数据库类型初始化数据库连接
	dsn := cfg.MainConfig.DbDSN
	switch cfg.MainConfig.DbType {
	case "mysql":
		db, err = initMysql(dsn)
	case "sqlite":
		db, err = initSqlite(dsn)
	case "postgres":
		db, err = initPostgres(dsn)
	default:
		return nil, fmt.Errorf("database type %s is not supported", cfg.MainConfig.DbType)
	}

	if err != nil {
		return nil, err
	}

	// 设置连接最大生命周期为 30 分钟
	db.SetConnMaxLifetime(30 * time.Minute)
	// 关闭 SQL 日志记录
	db.ShowSQL(false)
	// 同步数据库模型
	db.Sync2(new(model.TelegramPushSetting))
	return db, nil
}

// initMysql 初始化 MySQL 数据库连接
func initMysql(dsn string) (*xorm.Engine, error) {
	db, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	return db, nil
}

// initPostgres 初始化 PostgreSQL 数据库连接
func initPostgres(dsn string) (*xorm.Engine, error) {
	db, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	return db, nil
}

// initSqlite 初始化 SQLite 数据库连接
func initSqlite(dsn string) (*xorm.Engine, error) {
	db, err := xorm.NewEngine("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, nil
}

// SettingDao 返回设置数据访问对象
func (d *DataSource) SettingDao() dao.ISettingDao {
	return dao.NewSettingDaoImpl(d.db)
}
