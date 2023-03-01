package helpers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func CheckUserType(c *fiber.Ctx, role string) (err error) {
	userType := c.Get("user_type")
	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *fiber.Ctx, userId string) (err error) {
	userType := c.Get("user_type")
	uid := c.Get("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}
