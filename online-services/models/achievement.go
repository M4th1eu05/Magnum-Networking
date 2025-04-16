package models

import (
	"gorm.io/gorm"
)

type Achievement struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	Condition   string `json:"condition"`
	Threshold   int    `json:"threshold"`
}
