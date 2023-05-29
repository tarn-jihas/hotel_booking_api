package middleware

import (
	"hotel-reservation/api"
	"hotel-reservation/types"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return api.ErrUnauthorized()
	}

	if !user.IsAdmin {
		return api.ErrUnauthorized()
	}

	return c.Next()
}
