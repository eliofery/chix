package route

import (
	"github.com/eliofery/go-chix/internal/app/middleware"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
)

// UserRoute маршруты для пользователей
func (r *route) UserRoute(router *chix.Router) {
	log.Debug("Инициализация UserRoute")

	router.Group(func(rt *chix.Router) {
		rt.Use(middleware.IsAuth)

		rt.Get(apiV1+"/users", r.Handler.GetUsers)
	})
}
