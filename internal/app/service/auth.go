package service

import (
	"errors"
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/internal/app/repository"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/jwt"
	"github.com/eliofery/go-chix/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

// AuthService логика связанная с авторизацией
type AuthService interface {
	Register(ctx *chix.Ctx, user model.User) (token string, err error)
	Auth(ctx *chix.Ctx, user model.User) (token string, err error)
	Logout(ctx *chix.Ctx, userId int64) error
}

type authService struct {
	dao          repository.DAO
	tokenManager jwt.TokenManager
}

func NewAuthService(dao repository.DAO, tokenManager jwt.TokenManager) AuthService {
	log.Debug("Инициализация auth service")

	return &authService{dao: dao, tokenManager: tokenManager}
}

// Register регистрация аккаунта
func (s *authService) Register(ctx *chix.Ctx, user model.User) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.PasswordHash = string(passwordHash)

	userId, err := s.dao.NewUserQuery().Create(user)
	if err != nil {
		return "", err
	}

	token, err := s.tokenManager.GenerateToken(int(*userId))
	if err != nil {
		return "", err
	}

	if err = s.dao.NewSessionQuery().Create(*userId, token); err != nil {
		return "", err
	}

	s.tokenManager.SetCookieToken(ctx, token)

	return token, nil
}

// Auth авторизация аккаунта
func (s *authService) Auth(ctx *chix.Ctx, user model.User) (string, error) {
	findUser, err := s.dao.NewUserQuery().GetUserByEmail(user.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(findUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		return "", errors.New("не верный логин или пароль")
	}

	token, err := s.dao.NewSessionQuery().GetTokenByUserId(findUser.ID)
	if err == nil {
		s.tokenManager.SetCookieToken(ctx, token)
		return token, nil
	}

	token, err = s.tokenManager.GenerateToken(int(findUser.ID))
	if err != nil {
		return "", err
	}

	if err = s.dao.NewSessionQuery().Create(findUser.ID, token); err != nil {
		return "", err
	}

	s.tokenManager.SetCookieToken(ctx, token)

	return token, nil
}

// Logout выход из аккаунта
func (s *authService) Logout(ctx *chix.Ctx, userId int64) error {
	if err := s.dao.NewSessionQuery().DeleteByUserId(userId); err != nil {
		return err
	}

	s.tokenManager.RemoveCookieToken(ctx)

	return nil
}
