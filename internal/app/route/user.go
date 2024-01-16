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
		rt.Use(middleware.Example2())

		rt.Route(apiV1, func(rt *chix.Router) {
			rt.Get("/users", r.Handler.GetUsers)
		})
	})
}
