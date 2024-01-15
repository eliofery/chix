package service

import (
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/internal/app/repository"
	"github.com/eliofery/go-chix/pkg/log"
)

// UserService логика работы с пользователями
type UserService interface {
	GetUsers() (*[]model.User, error) // GetUsers получение всех пользователей
}

type userService struct {
	dao repository.DAO
}

// NewUserService конструктор
func NewUserService(dao repository.DAO) UserService {
	log.Debug("Инициализация UserService")

	return &userService{dao: dao}
}

// GetUsers получение всех пользователей
func (s *userService) GetUsers() (*[]model.User, error) {
	users, err := s.dao.NewUserQuery().GetUsers()
	if err != nil {
		return users, err
	}

	return users, err
}
