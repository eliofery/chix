package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/eliofery/go-chix/internal/app/dto"
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

// UserQuery запросы в базу данных для пользователей
type UserQuery interface {
	Create(user dto.UserSignUp) (userId int, err error)
	GetUsers() ([]model.User, error)
}

type userQuery struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

// Create создание пользователя
func (q *userQuery) Create(user dto.UserSignUp) (int, error) {
	var userId int

	query := "INSERT INTO users (first_name, last_name, age, email, password_hash) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := q.db.QueryRow(query, user.FirstName, user.LastName, user.Age, user.Email, user.Password).Scan(&userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return 0, errors.New("пользователь уже существует")
			}
		}
		return 0, err
	}

	return userId, nil
}

// GetUsers Получение всех пользователей
func (q *userQuery) GetUsers() ([]model.User, error) {
	qb := q.builder.Select("id", "first_name", "last_name", "age", "email").
		From(model.UserTableName)

	var users []model.User
	rows, err := qb.Query()
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age, &user.Email); err != nil {
			log.Error("Не удалось получить пользователей", slog.String("err", err.Error()))
			return nil, fmt.Errorf("не удалось получить пользователей")
		}
		users = append(users, user)
	}

	return users, nil
}
