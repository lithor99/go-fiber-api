package models

import "gorm.io/gorm"

type Images struct {
	gorm.Model
	ProductId uint   `json:"product_id,omitempty" validate:"required"`
	Image     string `json:"image,omitempty" validate:"required"`
}
