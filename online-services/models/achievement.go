package models

import (
	"gorm.io/gorm"
)

type Achievement struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	Condition   string `json:"condition"`
	StatsName   string `json:"stats_name"`
	Threshold   float64    `json:"threshold"`
}
