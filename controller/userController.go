package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/akgupta-47/auth-gofib/db"
	helpers "github.com/akgupta-47/auth-gofib/helpers"
	"github.com/akgupta-47/auth-gofib/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = db.GetUserCollection()
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		return c.Status(http.StatusBadRequest).SendString(validationErr.Error())
	}

	count, err := userCollection.CountDocuments(c.Context(), bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		return c.Status(http.StatusInternalServerError).SendString("Error while counting documents!!")
	}

	if count > 0 {
		return c.Status(http.StatusInternalServerError).SendString("Email already exists!!")
	}

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertionNumber, insertionErr := userCollection.InsertOne(c.Context(), user)
	if insertionErr != nil {
		return c.Status(http.StatusInternalServerError).SendString("User Item was not created!!")
	}
	return c.Status(http.StatusOK).JSON(resultInsertionNumber)
}

func Login()

func GetUsers()

func GetUser(c *fiber.Ctx) error {
	userId := c.Params("user_id")

	if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	var user models.User
	query := bson.M{"user_id": userId}
	err := userCollection.FindOne(c.Context(), query).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(user)
}
