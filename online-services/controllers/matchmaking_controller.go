package controllers

import (
	"fmt"
	"log"
	"net/http"
	"online-services/database"
	"online-services/models"
	"sort"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

type MatchmakingQueue struct {
	sync.Mutex
	players []models.User
	cond    *sync.Cond
}

var Queue *MatchmakingQueue

func NewMatchmakingQueue() *MatchmakingQueue {
	queue := &MatchmakingQueue{
		players: []models.User{},
	}
	queue.cond = sync.NewCond(&queue.Mutex)
	return queue
}

func (q *MatchmakingQueue) Enqueue(player models.User) {
	q.Lock()
	defer q.Unlock()
	q.players = append(q.players, player)
	q.cond.Signal()
}

func (q *MatchmakingQueue) Dequeue() models.User {
	q.Lock()
	defer q.Unlock()
	for len(q.players) == 0 {
		q.cond.Wait()
	}
	player := q.players[0]
	q.players = q.players[1:]
	return player
}

func (q *MatchmakingQueue) PlayerInQueue(player models.User) bool {
	q.Lock()
	defer q.Unlock()
	for _, p := range q.players {
		if p.ID == player.ID {
			return true
		}
	}
	return false
}

func (q *MatchmakingQueue) StartMatchmaking(matchSize int, matchHandler func([]models.User)) {
	go func() {
		for {
			q.Lock()
			for len(q.players) < matchSize {
				fmt.Println("Waiting for more players to join the queue...")
				q.cond.Wait()
			}

			sort.Slice(q.players, func(i, j int) bool {
				return q.players[i].SkillLevel < q.players[j].SkillLevel
			})

			var matchPlayers []models.User
			matchPlayers = q.players[:matchSize]
			q.players = q.players[matchSize:]

			q.Unlock()

			// Handle the match
			matchHandler(matchPlayers)
		}
	}()
}

type GameServerInfo struct {
	IP   string `form:"ip" json:"ip" binding:"required"`
	Port int    `form:"port" json:"port" binding:"required"`
}

func RegisterServer(c *gin.Context) {
	var serverInfo GameServerInfo
	if err := c.ShouldBindJSON(&serverInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	server := models.GameServer{
		IP:   serverInfo.IP,
		Port: serverInfo.Port,
	}
	if err := database.DB.Where("ip = ? AND port = ?", serverInfo.IP, serverInfo.Port).First(&server).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Server already registered"})
		return
	}

	if err := database.DB.Create(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register server"})
		return
	}

	c.JSON(http.StatusCreated, server)
}

// JoinQueue adds a player to the matchmaking queue
func JoinQueue(c *gin.Context) {
	UUID, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found"})
		return
	}

	var user models.User
	if err := database.DB.Where("uuid = ?", UUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	Queue.Enqueue(user)

	c.JSON(http.StatusOK, gin.H{"message": "Joined matchmaking queue"})
}

// Pulling request from players
func QueueStatus(c *gin.Context) {
	UUID, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found"})
		return
	}

	var user models.User
	if err := database.DB.Where("uuid = ?", UUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	inQueue := Queue.PlayerInQueue(user)

	if inQueue {
		c.JSON(http.StatusTooEarly, gin.H{"status": "in queue"})
	} else {
		// Check if the user has a game
		if user.GameID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "not in queue"})
			return
		}

		var game models.Game
		if err := database.DB.Model(&game).Where("id = ?", user.GameID).Preload(clause.Associations).First(&game).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			return
		}

		var server models.GameServer
		if err := database.DB.Where("id = ?", game.GameServerID).First(&server).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Game server not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"game": game, "server": server})
	}
}

// LeaveQueue removes a player from the matchmaking queue
func LeaveQueue(c *gin.Context) {
	UUID, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found"})
		return
	}

	var user models.User
	if err := database.DB.Where("uuid = ?", UUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	Queue.Dequeue()

	c.JSON(http.StatusOK, gin.H{"message": "Left matchmaking queue"})
}

func CreateMatch(players []models.User) {
	// find a game server
	var server models.GameServer
	if err := database.DB.Where("status = ?", "available").First(&server).Error; err != nil {
		log.Println("No available game server")
		return
	}

	game := models.Game{
		Players:      players,
		GameServerID: server.ID,
	}
	if err := database.DB.Create(&game).Error; err != nil {
		log.Fatal(err)
		return
	}

	if err := database.DB.Model(&players).Update("game_id", game.ID).Error; err != nil {
		log.Fatal(err)
		return
	}

	server.CurrentGame = &game
	server.Status = "in_game"
	if err := database.DB.Save(&server).Error; err != nil {
		log.Fatal(err)
		return
	}

	return
}

// StartGame Called when a match starts by the game server
func StartGame(c *gin.Context) {
	// get server info
	var serverInfo GameServerInfo
	if err := c.ShouldBindJSON(&serverInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find the game server
	var server models.GameServer
	if err := database.DB.Preload("CurrentGame").Preload("CurrentGame.Players").Where("ip = ? AND port = ?", serverInfo.IP, serverInfo.Port).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match started", "server": server})
}

type StatsInfo struct {
	UserUUID uuid.UUID `json:"user_uuid" binding:"required"`
	StatName string    `json:"stat_name" binding:"required"`
	Value    float64   `json:"value" binding:"required"`
}
type GameInfo struct {
	GameID uint        `json:"game_id" binding:"required"`
	Stats  []StatsInfo `json:"stats" binding:"required"`
}

// EndGame Called when a match ends by the game server
func EndGame(c *gin.Context) {
	var gameInfo GameInfo
	if err := c.ShouldBindJSON(&gameInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var players []models.User
	if err := database.DB.Where("game_id = ?", gameInfo.GameID).Find(&players).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Players not found"})
		return
	}

	if err := database.DB.Model(&players).Update("game_id", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear game association"})
		return
	}

	for _, player := range players {
		playerStats := make([]StatsInfo, 0)
		for _, stat := range gameInfo.Stats {
			if stat.UserUUID == player.UUID {
				playerStats = append(playerStats, stat)
			}
		}

		if len(playerStats) > 0 {
			if err := UpdateStats(player.UUID, playerStats); err != nil {
				log.Printf("Failed to update stats for player %s: %v", player.UUID, err)
			}
		}
	}

	if err := database.DB.Delete(&models.Game{}, gameInfo.GameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match ended"})
}
