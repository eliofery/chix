package controller

import (
	"github.com/eliofery/go-chix/internal/app/service"
	"github.com/eliofery/go-chix/pkg/log"
)

// ServiceController обработчик маршрутов
type ServiceController struct {
	authService service.AuthService
}

// NewServiceController конструктор
func NewServiceController(
	authService service.AuthService,
) ServiceController {
	log.Debug("Инициализация service controller")

	return ServiceController{
		authService: authService,
	}
}
