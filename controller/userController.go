package controller

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/akgupta-47/auth-gofib/db"
	"github.com/akgupta-47/auth-gofib/helpers"
	"github.com/akgupta-47/auth-gofib/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// var userCollection = db.Mgi.Db.Collection("user")
var validate = validator.New()

func TestRoute(c *fiber.Ctx) error {
	utype := c.Get("user_type")
	fmt.Println(utype)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password incorrect!!")
		check = false
	}
	return check, msg
}

func Signup(c *fiber.Ctx) error {
	var userCollection = db.GetUserCollection()
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorJson{Error: err.Error()})
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorJson{Error: validationErr.Error()})
	}

	// count, err := userCollection.CountDocuments(c.Context(), bson.M{"email": user.Email})

	count, err := userCollection.CountDocuments(c.Context(), bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: "Error while counting the documents!!"})
	}
	// fmt.Println(count)
	if count > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: "Email already exists!!"})
	}
	// return nil
	password := HashPassword(*user.Password)
	user.Password = &password

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertionNumber, insertionErr := userCollection.InsertOne(c.Context(), user)
	if insertionErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: "User Item was not created!!"})
	}
	return c.Status(fiber.StatusOK).JSON(resultInsertionNumber)
}

func Login(c *fiber.Ctx) error {
	var userCollection = db.GetUserCollection()
	user := new(models.User)
	foundUser := new(models.User)

	// while testing check what happens if empty email is sent
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorJson{Error: err.Error()})
	}

	err := userCollection.FindOne(c.Context(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: "email or password incorrect"})
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	if !passwordIsValid {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: msg})
	}

	if foundUser.Email == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: "user not found!!"})
	}
	token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
	helpers.UpdateAllTokens(c, token, refreshToken, foundUser.User_id)

	err = userCollection.FindOne(c.Context(), bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: err.Error()})
	}

	cookie := fiber.Cookie{
		Name:     "auth",
		Value:    token,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(foundUser)
}

func GetUsers(c *fiber.Ctx) error {
	var userCollection = db.GetUserCollection()
	if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorJson{Error: err.Error()})
	}

	recordsPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordsPerPage < 1 {
		recordsPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordsPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))

	matchStage := bson.D{{"$match", bson.D{{}}}}
	groupStage := bson.D{{"$group", bson.D{
		{"_id", bson.D{{"_id", "null"}}},
		{"total_count", bson.D{{"$sum", 1}}},
		{"data", bson.D{{"$push", "$$ROOT"}}}}}}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordsPerPage}}}}}}}

	result, err := userCollection.Aggregate(c.Context(), mongo.Pipeline{matchStage, groupStage, projectStage})

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: "error while listing user items"})
	}

	var allusers []bson.M
	if err = result.All(c.Context(), &allusers); err != nil {
		log.Fatal(err)
	}

	return c.Status(fiber.StatusOK).JSON(allusers[0])
}

func GetUser(c *fiber.Ctx) error {
	var userCollection = db.GetUserCollection()
	userId := c.Params("user_id")

	if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorJson{Error: err.Error()})
	}

	var user models.User
	query := bson.M{"user_id": userId}
	err := userCollection.FindOne(c.Context(), query).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorJson{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
