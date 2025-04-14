package models

import (
	"gorm.io/gorm"
)

type StoreItem struct {
	gorm.Model
	Name        string  `json:"name" gorm:"unique"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Type        string  `json:"type"` // e.g., "skin", "emote", "effect"
	IconURL     string  `json:"icon_url"`
}

type UserItem struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	ItemID    uint      `json:"item_id"`
	User      User      `gorm:"foreignKey:UserID"`
	StoreItem StoreItem `gorm:"foreignKey:ItemID"`
	Equipped  bool      `json:"equipped" gorm:"default:false"`
}
