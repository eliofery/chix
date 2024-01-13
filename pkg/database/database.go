package database

import (
	"database/sql"
	"github.com/eliofery/go-chix/pkg/log"
)

// Database интерфейс для инициализации базы данных
type Database interface {
	Init() (*sql.DB, error) // Init инициализация БД
}

// DB база данных
type DB struct {
	Conn *sql.DB
}

// Connect подключение к БД
// Пример: database.Connect(postgres.New(config))
func Connect(driver Database) (*DB, error) {
	log.Debug("Подключение к базе данных")

	var (
		db  DB
		err error
	)
	db.Conn, err = driver.Init()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// MustConnect подключение к БД с обработкой ошибок
func MustConnect(driver Database) *DB {
	db, err := Connect(driver)
	if err != nil {
		panic(err)
	}

	return db
}
