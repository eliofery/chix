package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/eliofery/go-chix/pkg/log"
)

var defaultFormat = squirrel.Dollar

// DAO интерфейс для обращения к БД
type DAO interface {
	NewAuthQuery() AuthQuery       // NewAuthQuery конструктор для запросов связанных с авторизацией
	NewUserQuery() UserQuery       // NewAuthQuery конструктор для запросов связанных с авторизацией
	NewSessionQuery() SessionQuery // NewSessionQuery конструктор для запросов связанных с сессиями
}

type dao struct {
	db *sql.DB
}

// NewDAO конструктор dao
func NewDAO(db *sql.DB) DAO {
	log.Debug("Инициализация DAO")

	return &dao{db: db}
}

// queryBuilder создание запросов в postgres базу данных
func (d *dao) queryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(defaultFormat).RunWith(d.db)
}

// NewAuthQuery конструктор для запросов связанных с авторизацией
func (d *dao) NewAuthQuery() AuthQuery {
	log.Debug("Инициализация AuthQuery")

	return &authQuery{db: d.db, builder: d.queryBuilder()}
}

// NewUserQuery конструктор для запросов связанных с пользователями
func (d *dao) NewUserQuery() UserQuery {
	log.Debug("Инициализация UserQuery")

	return &userQuery{db: d.db, builder: d.queryBuilder()}
}

// NewSessionQuery запросы в базу данных для сессий
func (d *dao) NewSessionQuery() SessionQuery {
	return &sessionQuery{db: d.db}
}
