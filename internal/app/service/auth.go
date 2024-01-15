package service

import (
	"github.com/eliofery/go-chix/internal/app/repository"
	"github.com/eliofery/go-chix/pkg/jwt"
	"github.com/eliofery/go-chix/pkg/log"
)

// AuthService логика работы с пользователями
type AuthService interface {
	Signup() // Signup регистрация пользователя
}

type authService struct {
	dao          repository.DAO
	tokenManager jwt.TokenManager
}

// NewAuthService конструктор
func NewAuthService(dao repository.DAO, tokenManager jwt.TokenManager) AuthService {
	log.Debug("Инициализация AuthService")

	return &authService{dao: dao, tokenManager: tokenManager}
}

// Signup регистрация пользователя
func (s *authService) Signup() {
	s.dao.NewAuthQuery().Register()
}
