package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"user-domain/infrastructure/config"
	applicationoutbound "user-domain/internal/application/outbound"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	NewGorm(l applicationoutbound.Logger) (*gorm.DB, error)
	GetConnect() *sql.DB
}

type database struct {
	conn *sql.DB
}

func (db *database) GetConnect() *sql.DB {
	return db.conn
}

func (db *database) GetConnection() *sql.DB {
	return db.conn
}
func (db *database) NewGorm(l applicationoutbound.Logger) (*gorm.DB, error) {

	psgConfig := postgres.Config{
		Conn: db.conn,
	}
	logger := NewGormLogger(logger.Info, l)
	dialector := postgres.New(psgConfig)
	return gorm.Open(dialector, &gorm.Config{Logger: logger})
}

func NewDatabase(cfg *config.Config) (Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresDatabase, cfg.PostgresPassword, cfg.PostgresSSLMode)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		println("ping err :", err.Error())
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Hour * 24)
	conn.SetConnMaxIdleTime(time.Hour * 12)
	return &database{
		conn: conn,
	}, nil
}

type gormLogger struct {
	LogLevel logger.LogLevel
	logger   applicationoutbound.Logger
}

func NewGormLogger(logLevel logger.LogLevel, logger applicationoutbound.Logger) logger.Interface {
	return &gormLogger{
		LogLevel: logLevel,
		logger:   logger,
	}
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.logger.WithContext(ctx).Info(msg, data...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.logger.WithContext(ctx).Warn(msg, data...)
	}
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.logger.WithContext(ctx).Error(msg, data...)
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	msg := fmt.Sprintf("elapsed=%v rows=%d sql=%s", elapsed, rows, sql)

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		l.logger.WithContext(ctx).Error(fmt.Sprintf("%v | %s", err, msg))

	case elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn:
		l.logger.WithContext(ctx).Warn("slow query: " + msg)

	case l.LogLevel >= logger.Info:
		l.logger.WithContext(ctx).Info(msg)
	}
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}
