package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lithor99/go-api-fiber-mysql/configs"
	"github.com/lithor99/go-api-fiber-mysql/middlewares"
	"github.com/lithor99/go-api-fiber-mysql/models"
	"github.com/lithor99/go-api-fiber-mysql/responses"
)

type ProductImage struct {
	Id    uint   `json:"id"`
	Image string `json:"url"`
}

type Product struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"created_by"`
}

type Data struct {
	Id            uint           `json:"id"`
	Name          string         `json:"name"`
	Price         float64        `json:"price"`
	Quantity      int            `json:"quantity"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Username      string         `json:"created_by"`
	ProductImages []ProductImage `json:"product_images"`
}

func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Products)
	user_id, _ := strconv.Atoi(middlewares.GetUserIdFromToken(c))

	if user_id != 0 {
		//validate the request body
		if bodyErr := c.BodyParser(product); bodyErr != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": bodyErr.Error()}})
		}

		//use the validator library to validate required fields
		product.CreatedBy = uint(user_id)
		if validationErr := validate.Struct(product); validationErr != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": validationErr.Error()}})
		}
		newProduct := models.Products{
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  product.Quantity,
			CreatedBy: product.CreatedBy,
		}

		configs.Database.Create(&newProduct)
		return c.Status(http.StatusCreated).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": newProduct}})
	}
	return c.Status(http.StatusUnauthorized).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": "You are not logged in"}})
}

func UploadProductImage(c *fiber.Ctx) error {
	id := c.Params("id")
	product_id, _ := strconv.Atoi(id)
	files, fileErr := c.MultipartForm()
	images := []models.Images{}
	if fileErr == nil {
		for _, fileHeaders := range files.File {
			for _, fileHeader := range fileHeaders {
				gen := strconv.Itoa(int(time.Now().UnixMilli()))
				replacer := strings.NewReplacer("/", "", "\\", "", " ", "")
				newFileName := replacer.Replace(fileHeader.Filename)
				c.SaveFile(fileHeader, fmt.Sprintf("./uploads/products/%s.%s", gen, newFileName))
				newImage := models.Images{
					ProductId: uint(product_id),
					Image:     "upload/products/" + gen + "." + newFileName,
				}
				images = append(images, newImage)
				configs.Database.Create(&newImage)
			}
		}
		return c.Status(http.StatusCreated).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": images}})
	}
	return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": fileErr}})
}

func GetProducts(c *fiber.Ctx) error {
	var datas []Data
	var products []Product
	var images []ProductImage
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

	configs.Database.Model(models.Products{}).Count(&count).Where("deleted_at IS NULL")
	if count%int64(limit) == 0 {
		totalPage = count / int64(limit)
	} else {
		totalPage = count/int64(limit) + 1
	}

	fields := "products.id, products.name, products.price, products.quantity, products.created_at, products.updated_at, products.deleted_at, users.username"
	joins := "left join users on users.id = products.created_by"
	res := configs.Database.Table("products").Select(fields).Offset(limit*page - limit).Limit(limit).Order("products.created_at desc").Joins(joins).Where("products.deleted_at IS NULL").Scan(&products)
	if res.RowsAffected == 0 {
		return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": "No data available"}})
	}
	for i := 0; i < len(products); i++ {
		configs.Database.Table("images").Select("images.id, images.image").Where("images.product_id = ?", products[i].Id).Scan(&images)
		datas = append(datas, Data{products[i].Id, products[i].Name, products[i].Price, products[i].Quantity, products[i].CreatedAt, products[i].UpdatedAt, products[i].Username, images})
	}
	return c.Status(http.StatusOK).JSON(responses.MultiData{Status: "success", TotalData: int(count), TotalPage: int(totalPage), CurentPage: page, Value: &fiber.Map{"data": datas}})
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var data Data
	var product Product
	var image []ProductImage

	fields := "products.id, products.name, products.price, products.quantity, products.created_at, products.updated_at, users.username"
	joins := "left join users on users.id = products.created_by"
	res := configs.Database.Table("products").Select(fields).Joins(joins).Where("products.id =?", id).Scan(&product)
	if res.RowsAffected == 0 {
		return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": "No data available"}})
	}
	configs.Database.Table("images").Select("images.id, images.image").Where("images.product_id = ?", product.Id).Scan(&image)
	data = Data{product.Id, product.Name, product.Price, product.Quantity, product.CreatedAt, product.UpdatedAt, product.Username, image}
	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": data}})
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product := new(models.Products)
	var data Data

	if err := c.BodyParser(product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": err.Error()}})
	}
	configs.Database.Where("id = ?", id).Updates(&product)

	fields := "products.id, products.name, products.price, products.quantity, products.created_at, products.updated_at, users.username"
	joins := "left join users on users.id = products.created_by"
	configs.Database.Table("products").Select(fields).Joins(joins).Where("products.id =?", id).Scan(&data)

	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": data}})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.Products
	var data Data

	configs.Database.Delete(&user, id)

	fields := "products.id, products.name, products.price, products.quantity, products.created_at, products.updated_at, users.username"
	joins := "left join users on users.id = products.created_by"
	configs.Database.Table("products").Select(fields).Joins(joins).Where("products.id =?", id).Scan(&data)
	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": user}})
}
