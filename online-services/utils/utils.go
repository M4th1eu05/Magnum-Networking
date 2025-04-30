package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"online-services/controllers"
	"online-services/database"
	"online-services/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notblank", NotBlank)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("convertFloat64", convertFloat64)
	}
}

func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return strings.TrimSpace(field) != ""
}

func convertFloat64(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	_, err := strconv.ParseFloat(field, 64)
	return err == nil
}

func StartMatchmaking() {
	controllers.Queue = controllers.NewMatchmakingQueue()
	matchSize := 4
	matchHandler := controllers.CreateMatch
	controllers.Queue.StartMatchmaking(matchSize, matchHandler)
}

func GenerateFakeData() {
	// This function is a placeholder for generating fake data.
	// You can implement your own logic to generate and insert fake data into the database.
	// For example, you can use a library like "github.com/bxcodec/faker/v3" to generate fake data.
	// Here, we will just print a message indicating that fake data generation is in progress.
	println("Generating fake data...")

	// add 25 fake users
	users := make([]models.User, 0)
	for i := 0; i < 25; i++ {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("password%d", i)), bcrypt.DefaultCost)
		user := models.User{
			Username: fmt.Sprintf("user%d", i),
			Password: string(hashedPassword),
			Stats: []models.Stat{
				{
					Name:  "kills",
					Value: rand.Float64() * 10,
				},
				{
					Name:  "deaths",
					Value: rand.Float64() * 10,
				},
			},
		}
		database.DB.Create(&user)
		users = append(users, user)
	}

	// add 2 achievements
	achievement := models.Achievement{
		Name:        "Achievement 1",
		Description: "Description for achievement 1",
		Condition:   "Condition for achievement 1",
		StatsName:   "kills",
		Threshold:   10,
	}
	database.DB.Create(&achievement)
	achievement = models.Achievement{
		Name:        "Achievement 2",
		Description: "Description for achievement 2",
		Condition:   "Condition for achievement 2",
		StatsName:   "deaths",
		Threshold:   10,
	}

	database.DB.Create(&achievement)

	// add 4 servers
	servers := []models.GameServer{
		{
			IP:   "127.0.0.1",
			Port: 8080,
		},
		{
			IP:   "127.0.0.1",
			Port: 8081,
		},
		{
			IP:   "127.0.0.1",
			Port: 8082,
			CurrentGame: &models.Game{
				Players: []models.User{users[4], users[5]},
			},
			Status: "in Game",
		},
		{
			IP:   "127.0.0.1",
			Port: 8083,
			CurrentGame: &models.Game{
				Players: []models.User{users[6], users[7]},
			},
			Status: "in Game",
		},
	}
	database.DB.Create(&servers)
}
