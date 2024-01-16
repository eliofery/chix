package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

// SessionQuery запросы в базу данных для сессий
type SessionQuery interface {
	Create(userId int, token string) error
	DeleteByToken(token string) error
	DeleteByUserId(userId int) error
}

type sessionQuery struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

// Create создание токена
func (q *sessionQuery) Create(userId int, token string) error {
	query := "INSERT INTO sessions (token, user_id) VALUES ($1, $2)"
	_, err := q.db.Exec(query, token, userId)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByToken удаление сессии по токену
func (q *sessionQuery) DeleteByToken(token string) error {
	query := "DELETE FROM sessions WHERE token = $1"
	_, err := q.db.Exec(query, token)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByUserId удаление сессии по id пользователя
func (q *sessionQuery) DeleteByUserId(userId int) error {
	query := "DELETE FROM sessions WHERE user_id = $1"
	_, err := q.db.Exec(query, userId)
	if err != nil {
		return err
	}

	return nil
}
