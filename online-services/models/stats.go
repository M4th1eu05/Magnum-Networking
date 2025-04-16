package models

import (
	"gorm.io/gorm"
)

const (
	Unranked = -1
	Bronze   = iota
	Silver
	Gold
	Diamond
	Master
)

type Stats struct {
	gorm.Model
	UserID uint    `json:"user_id" gorm:"not null"`
	Name   string  `json:"name" gorm:"unique"`
	Value  float64 `json:"value" gorm:"not null"`
}
