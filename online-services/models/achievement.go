package models

import (
	"gorm.io/gorm"
)

type Achievement struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	Condition   string `json:"condition"`
	Threshold   int    `json:"threshold"`
	Type        string `json:"type"` // e.g., "games_won", "cubes_cleared"
	IconURL     string `json:"icon_url"`
}

type UserAchievement struct {
	gorm.Model
	UserID        uint         `json:"user_id"`
	AchievementID uint         `json:"achievement_id"`
	UnlockedAt    int64       `json:"unlocked_at"`
	User          User         `gorm:"foreignKey:UserID"`
	Achievement   Achievement  `gorm:"foreignKey:AchievementID"`
}
