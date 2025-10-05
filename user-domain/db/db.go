package database

import (
	"database/sql"
	"fmt"
	"time"
	"user-domain/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	NewGorm() (*gorm.DB, error)
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
func (db *database) NewGorm() (*gorm.DB, error) {

	psgConfig := postgres.Config{
		Conn: db.conn,
	}
	dialector := postgres.New(psgConfig)
	return gorm.Open(dialector, nil)
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
