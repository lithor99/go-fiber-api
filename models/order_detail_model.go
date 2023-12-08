package models

import "gorm.io/gorm"

type OrderDetails struct {
	gorm.Model
	OrderId   uint `json:"order_id,omitempty" validate:"required"`
	ProductId uint `json:"product_id,omitempty" validate:"required"`
	Quantity  int  `json:"quantity,omitempty" validate:"required"`
}
