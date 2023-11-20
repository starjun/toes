package global

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mysqllogger "gorm.io/gorm/logger"
	"toes/internal/utils"
)

var (
	DB *gorm.DB
)

// InitStore 读取 db 配置，创建 gorm.DB 实例，并初始化 store 层.
func InitStore() error {
	// DecryptString
	if Cfg.Mysql.PasswordMode == "AES" {
		Cfg.Mysql.Password = utils.DecryptInternalValue(Cfg.Mysql.Password, Cfg.Seckey.Basekey, "mysql")
	}

	// Get password mod
	//if config.Cfg.Mysql.PasswordMode == db.PasswordModeMist {
	//	var _err error
	//	config.Cfg.Mysql.Password, _err = GetSecretByMist()
	//	if _err != nil {
	//		return _err
	//	}
	//}

	var err error
	DB, err = newMySQL(&Cfg.Mysql)
	if err != nil {
		LogErrorw("mysql连接失败", "Subject", "mysql", "Result", err)
		return err
	}

	LogDebugw("init db success")
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
		Logger:                 mysqllogger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
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
