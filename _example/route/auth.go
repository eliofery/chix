package route

import (
	"github.com/eliofery/go-chix/internal/app/middleware"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/log"
)

// AuthRoute маршруты для авторизации
func (r *route) AuthRoute(router *chix.Router) {
	log.Debug("Инициализация AuthRoute")

	router.Group(func(rt *chix.Router) {
		rt.Use(middleware.IsGuest)

		rt.Get(apiV1+"/signup", r.Handler.SignUp)
	})
}
