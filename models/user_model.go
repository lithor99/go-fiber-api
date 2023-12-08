package models

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Image    string `json:"image,omitempty"`
	Status   bool   `json:"status,omitempty" gorm:"default:true"`
}
