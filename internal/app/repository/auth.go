package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

// AuthQuery интерфейс для запросов связанных с пользователями
type AuthQuery interface {
	Register() // Register регистрация пользователя
}

type authQuery struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

// Register регистрация пользователя
func (q *authQuery) Register() {
	// TODO
}
