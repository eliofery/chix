package service

import (
	"github.com/eliofery/go-chix/internal/app/dto"
	"github.com/eliofery/go-chix/internal/app/repository"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/jwt"
	"github.com/eliofery/go-chix/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

// AuthService логика работы с пользователями
type AuthService interface {
	Register(ctx *chix.Ctx, user dto.UserSignUp) (token string, err error)
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

// Register регистрация пользователя
func (s *authService) Register(ctx *chix.Ctx, user dto.UserSignUp) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user.Password = string(passwordHash)
	userId, err := s.dao.NewUserQuery().Create(user)
	if err != nil {
		return "", err
	}

	token, err := s.tokenManager.GenerateToken(userId)
	if err != nil {
		return "", err
	}

	if err = s.dao.NewSessionQuery().Create(userId, token); err != nil {
		return "", err
	}

	s.tokenManager.SetCookieToken(ctx, token)

	return token, nil
}
