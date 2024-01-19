package controller

import (
	"github.com/eliofery/go-chix/internal/app/dto"
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/pkg/chix"
	"net/http"
)

// SignUp регистрация пользователя
func (c ServiceController) SignUp(ctx *chix.Ctx) error {
	var userDto dto.UserSignUp
	if err := ctx.Decode(&userDto); err != nil {
		return err
	}

	user := model.User{
		FirstName:    userDto.FirstName,
		LastName:     userDto.LastName,
		Email:        userDto.Email,
		PasswordHash: userDto.Password,
	}

	token, err := c.authService.Register(ctx, user)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(chix.Map{
		"success": true,
		"message": "Пользователь зарегистрирован",
		"token":   token,
	})
}
