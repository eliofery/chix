package controller

import (
	"github.com/eliofery/go-chix/internal/app/service"
	"github.com/eliofery/go-chix/pkg/log"
)

// ServiceController обработчик маршрутов
type ServiceController struct {
	authService service.AuthService
	userService service.UserService
}

// NewServiceController конструктор
func NewServiceController(
	authService service.AuthService,
	userService service.UserService,
) ServiceController {
	log.Debug("Инициализация ServiceController")

	return ServiceController{
		authService: authService,
		userService: userService,
	}
}
