package route

import (
	"github.com/eliofery/go-chix/internal/app/middleware"
	"github.com/eliofery/go-chix/pkg/chix"
)

// AuthRoute маршруты для авторизации
func (r *route) AuthRoute(router *chix.Router) {
	router.Group(func(rt *chix.Router) {
		rt.Use(middleware.IsGuest)

		rt.Post(apiV1+"/signup", r.Handler.SignUp)
		rt.Post(apiV1+"/signin", r.Handler.SignIn)
	})

	router.Group(func(rt *chix.Router) {
		rt.Use(middleware.IsAuth)

		rt.Post(apiV1+"/logout", r.Handler.Logout)
	})
}
