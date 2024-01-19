package controller

import (
	"github.com/eliofery/go-chix/internal/app/dto"
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/pkg/chix"
	"net/http"
)

// SignIn авторизация пользователя
func (c ServiceController) SignIn(ctx *chix.Ctx) error {
	var userDto dto.UserSignIn
	if err := ctx.Decode(&userDto); err != nil {
		return err
	}

	user := model.User{
		Email:        userDto.Email,
		PasswordHash: userDto.Password,
	}

	token, err := c.authService.Auth(ctx, user)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return err
	}

	return ctx.JSON(chix.Map{
		"success": true,
		"message": "Пользователь авторизован",
		"token":   token,
	})
}
