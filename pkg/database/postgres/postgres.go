package postgres

import (
	"database/sql"
	"fmt"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/database"
	"github.com/eliofery/go-chix/pkg/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"strconv"
)

const (
	portDefault = 5432
)

// Postgres база данных Postgres
// Пример: database.Connect(postgres.New(config))
type Postgres interface {
	database.Database
}

type postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// New конструктор Postgres
func New(config config.Config) Postgres {
	log.Debug("Инициализация базы данных Postgres")

	port, err := strconv.Atoi(config.Get("postgres.port"))
	if err != nil {
		port = portDefault
	}

	return &postgres{
		Host:     config.Get("postgres.host"),
		Port:     port,
		User:     config.Get("postgres.user"),
		Password: config.Get("postgres.password"),
		Database: config.Get("postgres.database"),
		SSLMode:  config.Get("postgres.sslmode"),
	}
}

// Init инициализация БД
func (s *postgres) Init() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", s.Host, s.Port, s.User, s.Password, s.Database, s.SSLMode)

	db, err := sql.Open("pgx", dsn)
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
