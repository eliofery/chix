package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/eliofery/go-chix/pkg/log"
)

var defaultFormat squirrel.PlaceholderFormat = squirrel.Dollar

// DAO интерфейс для обращения к БД
type DAO interface {
	NewUserQuery() UserQuery
	NewSessionQuery() SessionQuery
}

type dao struct {
	db *sql.DB
}

// NewDAO конструктор dao
func NewDAO(db *sql.DB) DAO {
	log.Debug("Инициализация dao")

	return &dao{db: db}
}

// queryBuilder создание запросов в postgres базу данных
func (d *dao) queryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(defaultFormat).RunWith(d.db)
}

// NewUserQuery запросы в базу данных связанные с пользователями
func (d *dao) NewUserQuery() UserQuery {
	return &userQuery{pgQb: d.queryBuilder()}
}

// NewSessionQuery запросы в базу данных связанные с сессиями
func (d *dao) NewSessionQuery() SessionQuery {
	return &sessionQuery{pgQb: d.queryBuilder()}
}
