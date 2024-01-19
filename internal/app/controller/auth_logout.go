package controller

import (
	"github.com/eliofery/go-chix/pkg/chix"
	"net/http"
)

// Logout выход пользователя
func (c ServiceController) Logout(ctx *chix.Ctx) error {
	userId := ctx.GetUserIdFromToken()

	if err := c.authService.Logout(ctx, *userId); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return err
	}

	return ctx.JSON(chix.Map{
		"success": true,
		"message": "Вы вышли из аккаунта",
	})
}
