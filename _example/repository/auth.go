package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

// AuthQuery запросы в базу данных для авторизации пользователей
type AuthQuery interface{}

type authQuery struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}
