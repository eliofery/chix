package repository

import (
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

// UserQuery запросы в базу данных связанные с пользователями
type UserQuery interface {
	Create(user model.User) (*int64, error)
	GetUserByEmail(email string) (*model.User, error)
}

type userQuery struct {
	pgQb squirrel.StatementBuilderType
}

// Create создание пользователя
func (q *userQuery) Create(user model.User) (*int64, error) {
	qb := q.pgQb.Insert(model.UserTableName).
		Columns("first_name", "last_name", "email", "password_hash").
		Values(user.FirstName, user.LastName, user.Email, user.PasswordHash).
		Suffix("RETURNING id")

	if err := qb.Scan(&user.ID); err != nil {
		log.Warn("Ошибка при создании пользователя", slog.String("err", err.Error()))

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, errors.New("пользователь уже существует")
			}
		}

		return nil, errors.New("пользователь не создан")
	}

	return &user.ID, nil
}

func (q *userQuery) GetUserByEmail(email string) (*model.User, error) {
	qb := q.pgQb.Select("id, first_name, last_name, email, password_hash").
		From(model.UserTableName).
		Where(squirrel.Eq{"email": email})

	var user model.User
	err := qb.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash)
	if err != nil {
		log.Warn("Ошибка при получении пользователя", slog.String("err", err.Error()))
		return nil, errors.New("не верный логин или пароль")
	}

	return &user, nil
}
