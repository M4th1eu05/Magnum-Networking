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
	UserID            uint64 `json:"user_id" gorm:"not null"`
	Rank              int    `json:"rank" gorm:"default:-1"`
	NbrCubesSpawned   int    `json:"nbr_cubes_spawned" gorm:"default:0"`
	NbrSpheresSpawned int    `json:"nbr_spheres_spawned" gorm:"default:0"`
	NbrGamesPlayed    int    `json:"nbr_games_played" gorm:"default:0"`
	NbrGamesWon       int    `json:"nbr_games_won" gorm:"default:0"`
}
