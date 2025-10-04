package database

import (
	"database/sql"
	"fmt"
	"time"
	config "user-domain/configs"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	NewGorm() (*gorm.DB, error)
	GetConnect() *sql.DB
}

type databse struct {
	cfg  config.DBConfig
	conn *sql.DB
}

func (db *databse) GetConnect() *sql.DB {
	return db.conn
}

func (db *databse) GetConnection() *sql.DB {
	return db.conn
}
func (db *databse) NewGorm() (*gorm.DB, error) {

	psgConfig := postgres.Config{
		Conn: db.conn,
	}
	dialector := postgres.New(psgConfig)
	return gorm.Open(dialector, nil)
}

func NewDatabase() (Database, error) {
	cfg := config.NewConfig()
	cfg.Load()
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.GetHostName(), cfg.GetPort(), cfg.GetUser(), cfg.GetDBName(), cfg.GetPassword(), cfg.GetSSLMode())
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Hour * 24)
	conn.SetConnMaxIdleTime(time.Hour * 12)
	return &databse{
		cfg:  cfg,
		conn: conn,
	}, nil
}
