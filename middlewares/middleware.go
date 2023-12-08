package middlewares

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/lithor99/go-api-fiber-mysql/configs"
	"github.com/lithor99/go-api-fiber-mysql/constants"
	"github.com/lithor99/go-api-fiber-mysql/models"
)

// generate token
func GenerateToken(user models.Users, c *fiber.Ctx) string {
	godotenv.Load()

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["username"] = user.Username
	claims["status"] = user.Status
	claims["created_at"] = user.CreatedAt
	claims["updated_at"] = user.UpdatedAt
	claims["exp"] = time.Now().Add(time.Hour * 30 * 24).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(os.Getenv(constants.SECRET_KEY)))

	if err != nil {
		// fmt.Sprintf("generating JWT Token failed: %v", err)
		return ""
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   60 * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return tokenString
}

// verify token
func VerifyToken(c *fiber.Ctx) error {
	godotenv.Load()
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}
		return []byte(os.Getenv(constants.SECRET_KEY)), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

	}

	var user models.Users
	configs.Database.First(&user, "id = ?", fmt.Sprint(claims["id"]))

	if strconv.Itoa(int(user.ID)) != fmt.Sprint(claims["id"]) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}
	return c.Next()
}

// get user id from token
func GetUserIdFromToken(c *fiber.Ctx) string {
	godotenv.Load()
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return ""
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(os.Getenv(constants.SECRET_KEY)), nil
	})
	if err != nil {
		return ""
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return ""

	}

	var user models.Users
	configs.Database.First(&user, "id = ?", fmt.Sprint(claims["id"]))

	if strconv.Itoa(int(user.ID)) != fmt.Sprint(claims["id"]) {
		return ""
	}
	return fmt.Sprint(claims["id"])
}
