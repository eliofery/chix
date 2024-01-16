package middleware

import (
	"errors"
	"github.com/eliofery/go-chix/internal/app/repository"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/jwt"
	"github.com/eliofery/go-chix/pkg/log"
	"log/slog"
	"strconv"
)

type key string

const IssuerKey key = "issuer"

var (
	ErrNotAllowed = errors.New("не допустимое действие")
)

// SetUserIdFromToken добавление ID авторизованного пользователя в контекст
func SetUserIdFromToken(dao repository.DAO, tokenManager jwt.TokenManager) chix.HandlerCtx {
	return func(ctx *chix.Ctx) error {
		cookieToken := ctx.Cookies(jwt.CookieTokenName)
		authToken := ctx.Get("Authorization")

		var tokenString string
		if cookieToken != "" {
			tokenString = cookieToken
		} else if authToken != "" {
			tokenString = authToken
		}

		if tokenString == "" {
			return ctx.Next()
		}

		issuer, err := tokenManager.VerifyToken(tokenString)
		if err != nil {
			tokenManager.RemoveCookieToken(ctx)
			if err = dao.NewSessionQuery().DeleteByToken(tokenString); err != nil {
				log.Error("Не удалось удалить сессионный токен", slog.String("err", err.Error()))
			}

			return ctx.Next()
		}

		userId, err := strconv.Atoi(issuer)
		if err != nil {
			log.Error("Не удалось получить идентификатор пользователя", slog.String("err", err.Error()))
			return ctx.Next()
		}

		ctx.Locals(IssuerKey, userId)

		return ctx.Next()
	}
}

// IsAuth доступ только для авторизованных пользователей
func IsAuth(ctx *chix.Ctx) error {
	_, ok := ctx.Locals(IssuerKey).(int)
	if !ok {
		return ErrNotAllowed
	}

	return ctx.Next()
}

// IsGuest доступ только для гостей
func IsGuest(ctx *chix.Ctx) error {
	_, ok := ctx.Locals(IssuerKey).(int)
	if !ok {
		return ctx.Next()
	}

	return ErrNotAllowed
}
