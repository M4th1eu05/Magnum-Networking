package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     uuid.UUID `gorm:"type:uuid"`
	Username string    `json:"name" gorm:"unique"`
	Password string    `json:"password" gorm:"not null"`
	IsAdmin  bool      `json:"is_admin" gorm:"default:false"`
	Stats    Stats     `json:"stats" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	// Set default stats
	u.Stats = Stats{
		Rank:              Unranked,
		NbrCubesSpawned:   0,
		NbrSpheresSpawned: 0,
		NbrGamesPlayed:    0,
		NbrGamesWon:       0,
	}
	return
}
