package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lithor99/go-api-fiber-mysql/configs"
	"github.com/lithor99/go-api-fiber-mysql/middlewares"
	"github.com/lithor99/go-api-fiber-mysql/models"
	"github.com/lithor99/go-api-fiber-mysql/responses"
)

type OrderDetailData struct {
	Id        uint    `json:"order_detail_id"`
	ProductId uint    `json:"product_id"`
	Name      string  `json:"product_name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type OrderData struct {
	Id          uint              `json:"order_id"`
	Username    string            `json:"ordered_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	OrderDetail []OrderDetailData `json:"order_detail"`
}

func CreateOrder(c *fiber.Ctx) error {
	type Item struct {
		Items []models.OrderDetails `json:"items,omitempty" validate:"required"`
	}
	item := new(Item)
	order := new(models.Orders)
	var order_data OrderData
	order_detail := []OrderDetailData{}
	var user models.Users
	user_id, _ := strconv.Atoi(middlewares.GetUserIdFromToken(c))
	if user_id != 0 {

		//validate the request body
		if err := c.BodyParser(item); err != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": err.Error()}})
		}

		//use the validator library to validate required fields
		order.OrderedBy = uint(user_id)
		if validationErr := validate.Struct(order); validationErr != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": validationErr.Error()}})
		}

		newOrder := models.Orders{
			OrderedBy: order.OrderedBy,
		}

		result := configs.Database.Create(&newOrder)
		if result.RowsAffected == 0 {
			return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": result.Error}})
		}

		for i := 0; i < len(item.Items); i++ {
			newOrderDetail := models.OrderDetails{
				OrderId:   newOrder.ID,
				ProductId: item.Items[i].ProductId,
				Quantity:  item.Items[i].Quantity,
			}
			// order_detail=append(order_detail, order_detail{newOrderDetail.ID, newOrderDetail.ProductId, newOrderDetail.Quantity, })
			configs.Database.Create(&newOrderDetail)
		}
		configs.Database.Find(&order, newOrder.ID)
		fields := "order_details.id, order_details.product_id, products.name, products.price, order_details.quantity"
		joins := "left join products on products.id = order_details.product_id"
		od_res := configs.Database.Table("order_details").Select(fields).Joins(joins).Where("order_details.order_id = ?", newOrder.ID).Scan(&order_detail)
		us_res := configs.Database.Find(&user, order.OrderedBy)
		if od_res.RowsAffected == 0 || us_res.RowsAffected == 0 {
			return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": "No data available"}})
		}
		order_data = OrderData{order.ID, user.Username, order.CreatedAt, order.UpdatedAt, order_detail}

		return c.Status(http.StatusCreated).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": order_data}})
	}
	return c.Status(http.StatusUnauthorized).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": "You are not logged in"}})
}

func GetOrders(c *fiber.Ctx) error {
	var orders []models.Orders
	// var order_details models.OrderDetails
	var user models.Users
	var orderDatas []OrderData
	var order_detail []OrderDetailData
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

	if count == 0 {
		configs.Database.Model(models.Orders{}).Count(&count).Where("deleted_at IS NULL")
		if count%int64(limit) == 0 {
			totalPage = count / int64(limit)
		} else {
			totalPage = (count / int64(limit)) + 1
		}
	}

	configs.Database.Offset(limit*page - limit).Limit(limit).Order("created_at desc").Find(&orders).Where("deleted_at IS NULL")
	for i := 0; i < len(orders); i++ {
		fields := "order_details.id, order_details.product_id, products.name, products.price, order_details.quantity"
		joins := "left join products on products.id = order_details.product_id"
		od_res := configs.Database.Table("order_details").Select(fields).Joins(joins).Where("order_details.deleted_at IS NULL AND order_details.order_id = ?", orders[i].ID).Scan(&order_detail)
		us_res := configs.Database.Find(&user, orders[i].OrderedBy)
		if od_res.RowsAffected == 0 || us_res.RowsAffected == 0 {
			return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": "no data available"}})
		}
		orderDatas = append(orderDatas, OrderData{orders[i].ID, user.Username, orders[i].CreatedAt, orders[i].UpdatedAt, order_detail})
	}

	// order_fields := "orders.id, orders.created_at, orders.updated_at, orders.deleted_at, users.username"
	// order_joins := "left join users on users.id = orders.ordered_by left join order_details on orders.id = order_details.order_id left join products on products.id = order_details.product_id"
	// configs.Database.Table("orders").Select(order_fields).Offset(limit*page - limit).Limit(limit).Order("orders.created_at desc").Joins(order_joins).Where("orders.deleted_at IS NULL").Scan(&orders)
	// for i := 0; i < len(orders); i++ {
	// 	detail_fields := "orders.id, orders.created_at, orders.updated_at, orders.deleted_at, users.username"
	// 	detail_joins := "left join users on users.id = orders.ordered_by left join order_details on orders.id = order_details.order_id left join products on products.id = order_details.product_id"
	// 	configs.Database.Table("orders").Select(fields).Offset(limit*page - limit).Limit(limit).Order("orders.created_at desc").Joins(joins).Where("orders.deleted_at IS NULL").Scan(&orders).Scan(&order_details)

	// }
	return c.Status(http.StatusOK).JSON(responses.MultiData{Status: "success", TotalData: int(count), TotalPage: int(totalPage), CurentPage: page, Value: &fiber.Map{"data": orderDatas}})
}

func GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	order := models.Orders{}
	order_detail := []OrderDetailData{}
	var order_data OrderData
	var user models.Users

	configs.Database.Find(&order).Where("deleted_at IS NULL")
	fields := "order_details.id, order_details.product_id, products.name, products.price, order_details.quantity"
	joins := "left join products on products.id = order_details.product_id"
	od_res := configs.Database.Table("order_details").Select(fields).Joins(joins).Where("order_details.deleted_at IS NULL AND order_details.order_id = ?", id).Scan(&order_detail)
	us_res := configs.Database.Find(&user, order.OrderedBy)
	if od_res.RowsAffected == 0 || us_res.RowsAffected == 0 {
		return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": "No data available"}})
	}
	order_id, _ := strconv.Atoi(id)
	order_data = OrderData{uint(order_id), user.Username, order.CreatedAt, order.UpdatedAt, order_detail}

	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": order_data}})
}

func UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	order := new(models.Orders)

	if err := c.BodyParser(order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SingleData{Status: "error", Value: &fiber.Map{"data": err.Error()}})
	}
	configs.Database.Where("id = ?", id).Updates(&order)
	configs.Database.Find(&order).Where("id = ?", id)
	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": order}})
}

func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	order := new(models.Orders)
	order_detail := new(models.OrderDetails)

	configs.Database.Delete(&order, id)
	configs.Database.Delete(&order_detail).Where("order_id = ?", id)
	configs.Database.Find(&order).Where("id = ?", id)
	return c.Status(http.StatusOK).JSON(responses.SingleData{Status: "success", Value: &fiber.Map{"data": order}})
}
