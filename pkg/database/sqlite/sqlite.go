package sqlite

import (
	"database/sql"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/database"
	"github.com/eliofery/go-chix/pkg/log"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

// Sqlite база данных sqlite
// Пример: database.Connect(sqlite.New(config))
type Sqlite interface {
	database.Database
}

type sqlite struct {
	Path string
}

// New конструктор Sqlite
func New(config config.Config) Sqlite {
	log.Debug("Инициализация базы данных Sqlite")

	return &sqlite{
		Path: config.Get("sqlite.path"),
	}
}

// Init инициализация БД
func (s *sqlite) Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		log.Error("Не удалось подключиться к базе данных", slog.String("err", err.Error()))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error("Не удалось подключиться к базе данных", slog.String("err", err.Error()))
		return nil, err
	}

	return db, nil
}
