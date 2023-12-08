package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/lithor99/go-api-fiber-mysql/constants"
	"github.com/lithor99/go-api-fiber-mysql/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

// var DATABASE_URI string = "root:@tcp(127.0.0.1)/go_test?charset=utf8mb4&parseTime=True&loc=Local"

func ConnectDB() error {
	var err error
	godotenv.Load()

	DATABASE_URI := os.Getenv(constants.USER_NAME) + ":" + os.Getenv(constants.PASSWORD) + "@tcp" + "(" + os.Getenv(constants.SERVER_NAME) + ")/" + os.Getenv(constants.DATABASE) + "?charset=utf8mb4&parseTime=True&loc=Local"

	Database, err = gorm.Open(mysql.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic(err)
	}

	Database.AutoMigrate(&models.Users{}, &models.Products{}, models.Orders{}, models.OrderDetails{}, models.Images{})

	return nil
}
