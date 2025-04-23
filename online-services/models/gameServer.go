package models

import "gorm.io/gorm"

type GameServer struct {
	gorm.Model
	IP     string `json:"ip" gorm:"not null"`
	Port   int    `json:"port" gorm:"not null"`
	Status string `json:"status" gorm:"default:'available'"`
}