package global

import (
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	mysqllogger "gorm.io/gorm/logger"
	"strings"
	"toes/internal/utils"
)

var (
	DB *gorm.DB
)

// InitStore 读取 db 配置，创建 gorm.DB 实例，并初始化 store 层.
func InitStore() error {
	var err error
	
	// 优先使用 SQLite（本地开发模式）
	if Cfg.Sqlite.Enabled {
		LogInfow("使用 SQLite 数据库", "path", Cfg.Sqlite.Path)
		DB, err = newSqlite(&Cfg.Sqlite)
		if err != nil {
			LogErrorw("sqlite连接失败", "Subject", "sqlite", "Result", err)
			cobra.CheckErr(err)
		}
		LogDebugw("init sqlite db success")
		return nil
	}

	// 使用 MySQL（生产环境）
	// DecryptString
	if strings.ToUpper(Cfg.Mysql.PasswordMode) == "AES" {
		Cfg.Mysql.Password = utils.DecryptInternalValue(Cfg.Seckey.JwtKey, Cfg.Mysql.Password, "mysql")
	}

	LogInfow("使用 MySQL 数据库", "host", Cfg.Mysql.Host, "database", Cfg.Mysql.Database)
	DB, err = newMySQL(&Cfg.Mysql)
	if err != nil {
		LogErrorw("mysql连接失败", "Subject", "mysql", "Result", err)
		cobra.CheckErr(err)
	}

	LogDebugw("init mysql db success")
	return nil
}

// DSN returns mysql dsn.
func (o *Mysql) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local",
	)
}

// newMySQL create a new gorm db instance with the given options.
func newMySQL(opts *Mysql) (*gorm.DB, error) {
	logLevel := mysqllogger.Silent
	if opts.LogLevel != 0 {
		logLevel = mysqllogger.LogLevel(opts.LogLevel)
	}
	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		Logger: mysqllogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}

// newSqlite create a new gorm db instance with sqlite.
func newSqlite(opts *Sqlite) (*gorm.DB, error) {
	logLevel := mysqllogger.Silent
	if opts.LogLevel != 0 {
		logLevel = mysqllogger.LogLevel(opts.LogLevel)
	}
	db, err := gorm.Open(sqlite.Open(opts.Path), &gorm.Config{
		Logger: mysqllogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
