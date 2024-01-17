package controller

import (
	"github.com/eliofery/go-chix/internal/app/dto"
	"github.com/eliofery/go-chix/pkg/chix"
	"net/http"
)

// SignUp авторизация нового пользователя
func (c ServiceController) SignUp(ctx *chix.Ctx) error {
	var user dto.UserSignUp
	if err := ctx.Decode(&user); err != nil {
		ctx.Status(http.StatusBadRequest)
		return err
	}

	token, err := c.authService.Register(ctx, user)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return err
	}

	return ctx.Status(http.StatusOK).JSON(chix.Map{
		"success": true,
		"message": "пользователь создан и авторизован",
		"token":   token,
	})
}
