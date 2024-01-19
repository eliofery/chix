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

// SetUserIdFromToken добавление ID авторизованного пользователя в контекст
func SetUserIdFromToken(dao repository.DAO, tokenManager jwt.TokenManager) chix.Handler {
	return func(ctx *chix.Ctx) error {
		token := tokenManager.GetToken(ctx)
		if token == "" {
			return ctx.Next()
		}

		if err := dao.NewSessionQuery().CheckByToken(token); err != nil {
			tokenManager.RemoveCookieToken(ctx)

			return ctx.Next()
		}

		issuer, err := tokenManager.VerifyToken(token)
		if err != nil {
			tokenManager.RemoveCookieToken(ctx)
			_ = dao.NewSessionQuery().DeleteByToken(token)

			return ctx.Next()
		}

		userId, err := strconv.Atoi(issuer)
		if err != nil {
			log.Error("Не удалось получить идентификатор пользователя", slog.String("err", err.Error()))

			return ctx.Next()
		}

		ctx.Locals(chix.IssuerKey, int64(userId))

		return ctx.Next()
	}
}

// IsAuth доступ только для авторизованных пользователей
func IsAuth(ctx *chix.Ctx) error {
	if ctx.GetUserIdFromToken() == nil {
		return errors.New("доступ только для авторизованных пользователей")
	}

	return ctx.Next()
}

// IsGuest доступ только для гостей
func IsGuest(ctx *chix.Ctx) error {
	if ctx.GetUserIdFromToken() == nil {
		return ctx.Next()
	}

	return errors.New("доступ только для не авторизованных пользователей")
}
