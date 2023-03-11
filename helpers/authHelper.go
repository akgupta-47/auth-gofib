package helpers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CheckUserType(c *fiber.Ctx, role interface{}) (err error) {
	userType := c.Locals("user_type")
	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *fiber.Ctx, userId string) (err error) {
	userType := c.Locals("user_type")
	uid := c.Locals("uid")
	err = nil

	fmt.Println(userType, uid)

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}
