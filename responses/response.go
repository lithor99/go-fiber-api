package responses

import "github.com/gofiber/fiber/v2"

type MultiData struct {
	Status     string     `json:"status"`
	TotalData  int        `json:"total_data"`
	TotalPage  int        `json:"total_page"`
	CurentPage int        `json:"curent_page"`
	Value      *fiber.Map `json:"value"`
}

type SingleData struct {
	Status string     `json:"status"`
	Value  *fiber.Map `json:"value"`
}

type ExcelData struct {
	Status string     `json:"status"`
	Total  int        `json:"total"`
	Value  *fiber.Map `json:"value"`
}
