package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	OrderedBy uint `json:"ordered_by,omitempty" validate:"required"`
}
