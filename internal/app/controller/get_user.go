package controller

import (
	"github.com/eliofery/go-chix/internal/app/dto"
	"github.com/eliofery/go-chix/pkg/chix"
	"net/http"
)

// GetUsers получение всех пользователей
func (c ServiceController) GetUsers(ctx *chix.Ctx) error {
	var user dto.User
	if err := ctx.Decode(&user, "fr"); err != nil {
		ctx.Status(http.StatusBadRequest)
		return err
	}

	users, err := c.userService.GetUsers()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return err
	}

	return ctx.Status(http.StatusOK).JSON(chix.Map{
		"name":  user.Name,
		"age":   user.Age,
		"users": users,
	})
}
