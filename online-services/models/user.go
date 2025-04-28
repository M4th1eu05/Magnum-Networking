package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID         uuid.UUID     `json:"UUID" gorm:"type:uuid"`
	Username     string        `json:"name" gorm:"unique"`
	Password     string        `json:"password" gorm:"not null"`
	Role         string        `json:"role" gorm:"default:'user'"`
	Stats        []Stat        `json:"stats" gorm:"constraint:OnDelete:CASCADE"`
	Achievements []Achievement `json:"achievements" gorm:"many2many:user_achievements;constraint:OnDelete:CASCADE"`
	GameID       *uint         `json:"game_id" gorm:"constraint:OnDelete:SET NULL;default:null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	u.Role = "user"
	return
}
