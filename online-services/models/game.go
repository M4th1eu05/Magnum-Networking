package models

import "gorm.io/gorm"

type Game struct {
	gorm.Model
	GameServerID uint       `json:"game_server_id" gorm:"not null"`
	GameServer   GameServer `json:"game_server" gorm:"not null"`
	Players      []User
}
