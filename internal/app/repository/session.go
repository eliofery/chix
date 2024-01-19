package repository

import (
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/pkg/log"
	"log/slog"
)

// SessionQuery запросы в базу данных связанные с сессиями
type SessionQuery interface {
	Create(userId int64, token string) error
	CheckByToken(string) error
	GetTokenByUserId(userId int64) (token string, err error)
	DeleteByToken(string) error
	DeleteByUserId(userId int64) error
}

type sessionQuery struct {
	pgQb squirrel.StatementBuilderType
}

// Create создание сессии
func (q *sessionQuery) Create(userId int64, token string) error {
	qb := q.pgQb.Insert(model.SessionTableName).
		Columns("user_id", "token").
		Values(userId, token)

	if _, err := qb.Exec(); err != nil {
		log.Warn("Ошибка при создании сессии", slog.String("err", err.Error()))
		return errors.New("сессия не создана")
	}

	return nil
}

// CheckByToken проверка наличия сессии по токену
func (q *sessionQuery) CheckByToken(token string) error {
	qb := q.pgQb.Select("user_id").
		From(model.SessionTableName).
		Where(squirrel.Eq{"token": token})

	var userId int64
	if err := qb.Scan(&userId); err != nil {
		log.Warn("Ошибка при проверке сессии", slog.String("err", err.Error()))
		return errors.New("сессия не существует")
	}

	return nil
}

// GetTokenByUserId получение токена по id пользователя
func (q *sessionQuery) GetTokenByUserId(userId int64) (string, error) {
	qb := q.pgQb.Select("token").
		From(model.SessionTableName).
		Where(squirrel.Eq{"user_id": userId})

	var token string
	if err := qb.Scan(&token); err != nil {
		log.Warn("Ошибка при получении токена по id", slog.String("err", err.Error()))
		return "", errors.New("токен не существует")
	}

	return token, nil
}

// DeleteByToken удаление сессии по токену
func (q *sessionQuery) DeleteByToken(token string) error {
	qb := q.pgQb.Delete(model.SessionTableName).
		Where(squirrel.Eq{"token": token})

	if _, err := qb.Exec(); err != nil {
		log.Warn("Ошибка при удалении сессии", slog.String("err", err.Error()))
		return errors.New("сессия не удалена")
	}

	return nil
}

// DeleteByUserId удаление сессии по идентификатору пользователя
func (q *sessionQuery) DeleteByUserId(userId int64) error {
	qb := q.pgQb.Delete(model.SessionTableName).
		Where(squirrel.Eq{"user_id": userId})

	if _, err := qb.Exec(); err != nil {
		log.Warn("Ошибка при удалении сессии", slog.String("err", err.Error()))
		return errors.New("сессия не удалена")
	}

	return nil
}
