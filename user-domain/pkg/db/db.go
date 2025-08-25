package database

import (
	"database/sql"
	"fmt"
	config "user-domain/configs"

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
	return gorm.Open(postgres.New(postgres.Config{
		Conn: db.conn,
	}), nil)
}

func NewDatabase() (Database, error) {
	cfg := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.GetHostName(), cfg.GetPort(), cfg.GetUser(), cfg.GetDBName(), cfg.GetPassword(), cfg.GetSSLMode())
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &databse{
		cfg:  cfg,
		conn: conn,
	}, nil
}
