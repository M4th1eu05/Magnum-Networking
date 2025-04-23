package models
//
//import (
//	"gorm.io/gorm"
//)
//
//type GameServer struct {
//	gorm.Model
//	Address     string `json:"address" gorm:"unique"`
//	Port        int    `json:"port"`
//	MaxPlayers  int    `json:"max_players"`
//	CurrentGame *Game  `json:"current_game"`
//	Status      string `json:"status"` // "available", "in_game", "offline"
//}
//
//type Game struct {
//	gorm.Model
//	ServerID   uint        `json:"server_id"`
//	Server     GameServer  `json:"server" gorm:"foreignKey:ServerID"`
//	Players    []User      `json:"players" gorm:"many2many:game_players;"`
//	Status     string     `json:"status"` // "waiting", "in_progress", "finished"
//	StartTime  int64      `json:"start_time"`
//	EndTime    int64      `json:"end_time"`
//}
//
//type QueuedPlayer struct {
//	gorm.Model
//	UserID    uint      `json:"user_id"`
//	User      User      `json:"user" gorm:"foreignKey:UserID"`
//	JoinedAt  int64     `json:"joined_at"`
//	Status    string    `json:"status"` // "queued", "matched", "cancelled"
//}
