package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
	UUID     uuid.UUID `gorm:"type:uuid"`
	Username  string `json:"name" gorm:"unique"`
	Password  string `json:"password" gorm:"not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.UUID = uuid.New()
    return
}