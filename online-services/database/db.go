package database

import (
	"log"
	"online-services/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
    if gin.Mode() == gin.TestMode {
        DB, err = gorm.Open(sqlite.Open("../database/test.db"), &gorm.Config{})
    } else {
        DB, err = gorm.Open(sqlite.Open("./database/data.db"), &gorm.Config{})
        
    }
	if err != nil {
		log.Fatal(err)
	}
    initializeDB()
}

func CloseDB() {
    if sqlDB, err := DB.DB(); err == nil {
        if gin.Mode() == gin.TestMode {
            // wipe the test database
            sqlDB.Exec("DELETE FROM users")
        }
        sqlDB.Close()
    }
}

func initializeDB() {
    // AutoMigrate all models    if err := DB.AutoMigrate(
        &models.User{},
        &models.Stats{},
        &models.Achievement{},
        &models.UserAchievement{},
        &models.GameServer{},
        &models.Game{},
        &models.QueuedPlayer{},
        &models.StoreItem{},
        &models.UserItem{},
    ); err != nil {
        log.Fatal(err)
    }
}
