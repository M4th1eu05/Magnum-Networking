package tests

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"online-services/database"
	"online-services/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAndStatsRelation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Créer un utilisateur
	user := models.User{
		Username: "testUser",
		Password: "testPass",
	}
	database.DB.Create(&user)

	// Ajouter des statistiques à l'utilisateur
	stats := []models.Stat{
		{Name: "kills", Value: 10, UserID: &user.ID},
		{Name: "deaths", Value: 5, UserID: &user.ID},
	}
	for _, stat := range stats {
		database.DB.Create(&stat)
	}

	// Vérifier la relation User -> Stats
	var retrievedUser models.User
	database.DB.Preload("Stats").First(&retrievedUser, user.ID)
	assert.Equal(t, 2, len(retrievedUser.Stats))
	assert.Equal(t, "kills", retrievedUser.Stats[0].Name)
}

func TestGameServerAndGameRelation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Créer un serveur de jeu
	server := models.GameServer{
		IP:   "127.0.0.1",
		Port: 8080,
	}
	database.DB.Create(&server)

	// Associer un jeu au serveur
	game := models.Game{
		GameServerID: server.ID,
	}
	database.DB.Create(&game)

	// Vérifier la relation GameServer -> Game
	var retrievedServer models.GameServer
	database.DB.Preload("CurrentGame").First(&retrievedServer, server.ID)
	assert.NotNil(t, retrievedServer.CurrentGame)
	assert.Equal(t, game.ID, retrievedServer.CurrentGame.ID)
}

func TestGameAndPlayersRelation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Créer un jeu
	game := models.Game{}
	database.DB.Create(&game)

	// Ajouter des joueurs au jeu
	players := []models.User{
		{Username: "player1", Password: "pass1", GameID: &game.ID},
		{Username: "player2", Password: "pass2", GameID: &game.ID},
	}
	for _, player := range players {
		database.DB.Create(&player)
	}

	// Vérifier la relation Game -> Players
	var retrievedGame models.Game
	database.DB.Preload("Players").First(&retrievedGame, game.ID)
	assert.Equal(t, 2, len(retrievedGame.Players))
	assert.Equal(t, "player1", retrievedGame.Players[0].Username)
}

func TestUserAndAchievementsRelation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Créer un utilisateur
	user := models.User{
		Username: "achievementTester",
		Password: "testPass",
	}
	database.DB.Create(&user)

	// Créer des achievements
	achievement := models.Achievement{
		Name:        "First Blood",
		Description: "Get your first kill",
		Condition:   "kills >= 1",
		StatsName:   "kills",
		Threshold:   1,
	}
	database.DB.Create(&achievement)

	// Associer l'achievement à l'utilisateur
	database.DB.Model(&user).Association("Achievements").Append(&achievement)

	// Vérifier la relation User -> Achievements
	var retrievedUser models.User
	database.DB.Preload("Achievements").First(&retrievedUser, user.ID)
	assert.Equal(t, 1, len(retrievedUser.Achievements))
	assert.Equal(t, "First Blood", retrievedUser.Achievements[0].Name)
}

func TestUserUniqueConstraint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Créer un utilisateur
	user := models.User{
		Username: "uniqueUser",
		Password: "testPass",
	}
	result := database.DB.Create(&user)
	assert.Nil(t, result.Error)

	// Tenter de créer un utilisateur avec le même nom d'utilisateur
	duplicateUser := models.User{
		Username: "uniqueUser",
		Password: "anotherPass",
	}
	result = database.DB.Create(&duplicateUser)
	assert.NotNil(t, result.Error)
	assert.Contains(t, result.Error.Error(), "UNIQUE constraint failed")
}

func TestStatNotNullConstraint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Tenter de créer une statistique sans UserID
	stat := models.Stat{
		Name:  "kills",
		Value: 10,
	}
	result := database.DB.Create(&stat)
	assert.NotNil(t, result.Error)
	assert.Contains(t, result.Error.Error(), "NOT NULL constraint failed")
}

func TestCascadeDeleteUserStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Créer un utilisateur avec des statistiques
	user := models.User{
		Username: "cascadeUser",
		Password: "testPass",
	}
	database.DB.Create(&user)

	stats := []models.Stat{
		{Name: "kills", Value: 10, UserID: &user.ID},
		{Name: "deaths", Value: 5, UserID: &user.ID},
	}
	for _, stat := range stats {
		database.DB.Create(&stat)
	}

	// Supprimer l'utilisateur
	database.DB.Select(clause.Associations).Delete(&user)

	// Vérifier que les statistiques associées ont été supprimées
	var count int64
	database.DB.Model(&models.Stat{}).Where("user_id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestCascadeDeleteGamePlayers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Créer un jeu avec des joueurs
	game := models.Game{}
	database.DB.Create(&game)

	players := []models.User{
		{Username: "player1", Password: "pass1", GameID: &game.ID},
		{Username: "player2", Password: "pass2", GameID: &game.ID},
	}
	for _, player := range players {
		database.DB.Create(&player)
	}

	// Supprimer le jeu
	database.DB.Select(clause.Associations).Delete(&game)

	// Vérifier que les joueurs associés ont leur GameID mis à NULL
	var retrievedPlayers []models.User
	database.DB.Where("game_id IS NOT NULL").Find(&retrievedPlayers)
	assert.Equal(t, 0, len(retrievedPlayers))
}

func TestAchievementUniqueConstraint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Créer un achievement
	achievement := models.Achievement{
		Name:        "UniqueAchievement",
		Description: "Test unique constraint",
		Condition:   "kills >= 1",
		StatsName:   "kills",
		Threshold:   1,
	}
	result := database.DB.Create(&achievement)
	assert.Nil(t, result.Error)

	// Tenter de créer un achievement avec le même nom
	duplicateAchievement := models.Achievement{
		Name:        "UniqueAchievement",
		Description: "Duplicate test",
		Condition:   "kills >= 5",
		StatsName:   "kills",
		Threshold:   5,
	}
	result = database.DB.Create(&duplicateAchievement)
	assert.NotNil(t, result.Error)
	assert.Contains(t, result.Error.Error(), "UNIQUE constraint failed")
}
