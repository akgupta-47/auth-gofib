package middleware

import (
	"github.com/akgupta-47/auth-gofib/helpers"
	"github.com/gofiber/fiber/v2"
)

func Authenticate(c *fiber.Ctx) error {
	clientToken := c.GetReqHeaders()
	if clientToken["token"] == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "no authorization header provided!!",
		})
	}

	claims, err := helpers.ValidateToken(clientToken["token"])
	if err != "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	c.Set("email", claims.Email)
	c.Set("firat_name", claims.First_name)
	c.Set("last_name", claims.Last_name)
	c.Set("uid", claims.Uid)
	c.Set("user_type", claims.User_type)

	return c.Next()
}
