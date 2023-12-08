package models

import (
	"gorm.io/gorm"
)

type Products struct {
	gorm.Model
	Name      string  `json:"name,omitempty" validate:"required"`
	Price     float64 `json:"price,omitempty" validate:"required"`
	Quantity  int     `json:"quantity,omitempty" validate:"required"`
	CreatedBy uint    `json:"created_by,omitempty" validate:"required"`
}
