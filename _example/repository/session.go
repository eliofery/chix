package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/pkg/log"
	"log/slog"
)

// SessionQuery запросы в базу данных для сессий
type SessionQuery interface {
	Create(userId int, token string) error
	GetIdByToken(token string) (id int, err error)
	DeleteByToken(token string) error
	DeleteByUserId(userId int) error
}

type sessionQuery struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

// Create создание токена
func (q *sessionQuery) Create(userId int, token string) error {
	qb := q.builder.
		Insert(model.SessionTableName).
		Columns("user_id", "token").
		Values(userId, token)

	if _, err := qb.Exec(); err != nil {
		log.Error("Не удалось создать сессионный токен", slog.String("err", err.Error()))
		return err
	}

	return nil
}

// GetIdByToken получение ID сессии по токену
func (q *sessionQuery) GetIdByToken(token string) (int, error) {
	query := "SELECT id FROM sessions WHERE token = $1"

	var id int
	err := q.db.QueryRow(query, token).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
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
