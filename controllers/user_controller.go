package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/lithor99/go-api-fiber-mysql/configs"
	"github.com/lithor99/go-api-fiber-mysql/middlewares"
	"github.com/lithor99/go-api-fiber-mysql/models"
	"github.com/lithor99/go-api-fiber-mysql/responses"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes)
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.Users)

	//validate the request body
	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": err.Error()}})
	}

	file, err := c.FormFile("image")

	// Check for errors:
	if err == nil {
		//upload image file
		gen := strconv.Itoa(int(time.Now().UnixMilli()))
		replacer := strings.NewReplacer("/", "", "\\", "", " ", "")
		newFileName := replacer.Replace(file.Filename)
		c.SaveFile(file, fmt.Sprintf("./uploads/%s.%s", gen, newFileName))
		user.Image = fmt.Sprintf("/uploads/%s.%s", gen, newFileName)

		//use the validator library to validate required fields
		if validationErr := validate.Struct(user); validationErr != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": validationErr.Error()}})
		}

		newUser := models.Users{
			Username: user.Username,
			Password: HashPassword(user.Password),
			Status:   user.Status,
			Image:    user.Image,
		}

		configs.Database.Create(&newUser)

		return c.Status(http.StatusCreated).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": newUser}})
	}
	return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error2", Value: &fiber.Map{"data": err}})
}

func GetUsers(c *fiber.Ctx) error {
	var users []models.Users
	var count int64
	var totalPage int64
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	configs.Database.Model(models.Users{}).Count(&count).Where("deleted_at =?", nil)
	if count%int64(limit) == 0 {
		totalPage = count / int64(limit)
	} else {
		totalPage = count/int64(limit) + 1
	}

	configs.Database.Offset(limit*page-limit).Limit(limit).Order("created_at desc").Find(&users).Where("deleted_at =?", nil)
	return c.Status(http.StatusOK).JSON(responses.MultiData{Status: "success", TotalData: int(count), TotalPage: int(totalPage), CurentPage: page, Value: &fiber.Map{"data": users}})
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.Users

	result := configs.Database.Find(&user, id)

	if result.RowsAffected == 0 {
		return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": "no data available"}})
	}
	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": user}})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": err.Error()}})
	}

	configs.Database.Where("id = ?", id).Updates(&user)
	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": user}})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.Users

	configs.Database.Delete(&user, id)
	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": user}})
}

func Login(c *fiber.Ctx) error {
	var user models.Users
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	fmt.Println(payload.Username, payload.Password)

	result := configs.Database.Where(&models.Users{Username: payload.Username}).Find(&user)

	match := ComparePassword(payload.Password, user.Password)

	fmt.Println(user.Username, user.Password)

	if result.RowsAffected == 0 {
		return c.SendStatus(400)
	}

	if !match {
		return c.SendStatus(400)
	}

	tokenString := middlewares.GenerateToken(user, c)
	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": tokenString}})
}
