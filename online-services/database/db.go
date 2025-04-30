package database

import (
	"golang.org/x/crypto/bcrypt"
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
			var tables, err = DB.Migrator().GetTables()
			if err != nil {
				log.Fatal(err)
			}

			for _, table := range tables {
				if table == "sqlite_sequence" {
					continue
				}
				if err := DB.Migrator().DropTable(table); err != nil {
					log.Fatal(err)
				}
			}
		}

		err := sqlDB.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func initializeDB() {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Stat{},
		&models.Achievement{},
		&models.GameServer{},
		&models.Game{},
	); err != nil {
		log.Fatal(err)
	}

	var admin models.User
	if err := DB.Where("username = ?", "admin").First(&admin).Error; err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		admin = models.User{
			Username: "admin",
			Password: string(hashedPassword),
		}
		DB.Create(&admin)
	}
	admin.Role = "admin"
	DB.Save(&admin)
}
