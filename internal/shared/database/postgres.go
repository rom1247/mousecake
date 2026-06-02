// Package database 提供 PostgreSQL 数据库连接初始化和管理功能。
package database

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/mousecake-go/mousecake-go/config"
)

// slogWriter 桥接 slog 到 Gorm 日志接口。
type slogWriter struct {
	logger *slog.Logger
}

func (w *slogWriter) Printf(format string, args ...interface{}) {
	w.logger.Info(fmt.Sprintf(format, args...))
}

// NewPostgres 根据配置创建 PostgreSQL 数据库连接，配置连接池和日志适配。
func NewPostgres(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	gormCfg := &gorm.Config{
		Logger: newGormLogger(cfg),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

func newGormLogger(cfg config.DatabaseConfig) gormlogger.Interface {
	writer := &slogWriter{logger: slog.Default()}

	logLevel := gormlogger.Info
	switch cfg.LogLevel {
	case "warn":
		logLevel = gormlogger.Warn
	case "error":
		logLevel = gormlogger.Error
	case "silent":
		logLevel = gormlogger.Silent
	}

	return gormlogger.New(
		writer,
		gormlogger.Config{
			SlowThreshold:             cfg.SlowThreshold,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}
