package middleware

import (
	"github.com/akgupta-47/auth-gofib/helpers"
	"github.com/gofiber/fiber/v2"
)

func Authenticate(c *fiber.Ctx) error {
	// ReqHeaders := c.GetReqHeaders()
	cookieToken := c.Cookies("auth")
	// fmt.Println(cookieToken)
	// clientToken := strings.Split(ReqHeaders["Authorization"], " ")[1]
	if cookieToken == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "no authorization header provided!!",
		})
	}

	claims, err := helpers.ValidateToken(cookieToken)
	if err != "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error validating token:" + err,
		})
	}
	c.Locals("email", claims.Email)
	c.Locals("firat_name", claims.First_name)
	c.Locals("last_name", claims.Last_name)
	c.Locals("uid", claims.Uid)
	c.Locals("user_type", claims.User_type)
	// fmt.Println(c)
	return c.Next()
}
