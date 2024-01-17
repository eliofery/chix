package controller

import (
	"github.com/eliofery/go-chix/pkg/chix"
	"net/http"
)

// GetUsers получение всех пользователей
func (c ServiceController) GetUsers(ctx *chix.Ctx) error {
	userId := ctx.GetUserIdFromToken()

	users, err := c.userService.GetUsers()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return err
	}

	return ctx.Status(http.StatusOK).JSON(chix.Map{
		"success": true,
		"message": "пользователи получены",
		"users":   users,
		"user_id": userId,
	})
}
