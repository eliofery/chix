package chix

import (
	"github.com/eliofery/go-chix/pkg/log"
	"log/slog"
	"net/http"
	"time"
)

// NewCookie создание куки
func (ctx *Ctx) NewCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}

// Cookie cохранение куки
func (ctx *Ctx) Cookie(cookie *http.Cookie) {
	http.SetCookie(ctx.ResponseWriter, cookie)
}

// Cookies получение куки
func (ctx *Ctx) Cookies(name string, defaultValue ...string) string {
	cookie, err := ctx.Request.Cookie(name)
	if err != nil {
		log.Debug("Cookie не найдено", slog.Any("cookie", name))

		if len(defaultValue) == 0 {
			return ""
		}

		return defaultValue[0]
	}

	return cookie.Value
}

// ClearCookie удаление куки
func (ctx *Ctx) ClearCookie(name string) {
	cookie := ctx.NewCookie(name, "")
	cookie.MaxAge = -1
	cookie.Expires = time.Now().Add(-time.Hour)

	http.SetCookie(ctx.ResponseWriter, cookie)
}
