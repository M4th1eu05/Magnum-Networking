package models

import (
	"gorm.io/gorm"
)

type Achievement struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Condition   string  `json:"condition"`
	StatsName   string  `json:"stats_name" gorm:"not null"`
	Threshold   float64 `json:"threshold" gorm:"not null"`
}
