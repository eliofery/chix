package route

import (
	"github.com/eliofery/go-chix/internal/app/middleware"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
)

// UserRoute маршруты для пользователей
func (r *route) UserRoute(router *chix.Router) {
	log.Debug("Инициализация UserRoute")

	api := apiV1 + "/users"

	router.Route(api, func(rt *chix.Router) {
		rt.Group(func(rt *chix.Router) {
			rt.Use(middleware.IsAuth)
			rt.Get("/", r.Handler.GetUsers)
		})
	})
}
