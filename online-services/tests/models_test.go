package tests

import (
	"online-services/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "online-services/models"
)

func TestUserStatsRelation(t *testing.T) {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.Stats{})
	assert.NoError(t, err)

	user := models.User{
		Username: "testuser",
		Password: "password123",
		Stats: models.Stats{
			Rank:              models.Silver,
			NbrCubesSpawned:   10,
			NbrSpheresSpawned: 5,
			NbrGamesPlayed:    20,
			NbrGamesWon:       15,
		},
	}

	err = db.Create(&user).Error
	assert.NoError(t, err)

	var retrievedUser models.User
	err = db.Preload("Stats").First(&retrievedUser, "username = ?", "testuser").Error
	assert.NoError(t, err)

	assert.Equal(t, "testuser", retrievedUser.Username)
	assert.Equal(t, models.Silver, retrievedUser.Stats.Rank)
	assert.Equal(t, 10, retrievedUser.Stats.NbrCubesSpawned)
	assert.Equal(t, 5, retrievedUser.Stats.NbrSpheresSpawned)
	assert.Equal(t, 20, retrievedUser.Stats.NbrGamesPlayed)
	assert.Equal(t, 15, retrievedUser.Stats.NbrGamesWon)
}
