package models

import (
	"gorm.io/gorm"
)

type Stat struct {
	gorm.Model
	UserID *uint   `json:"user_id" gorm:"not null"`
	Name   string  `json:"name" gorm:"not null"`
	Value  float64 `json:"value" gorm:"not null"`
}
