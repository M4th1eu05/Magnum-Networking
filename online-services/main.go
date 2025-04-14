package main

import (
	"net/http"
	"online-services/controllers"
	"online-services/database"
	"online-services/middlewares"
	"online-services/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.ConnectDB()
	defer database.CloseDB()

	utils.InitValidators()
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)

	authorized := router.Group("/", middlewares.AuthMiddleware())
	{
		// User and Stats endpoints
		authorized.GET("/user/:username/stats", controllers.GetUserStats)
		
		// Achievement endpoints
		authorized.GET("/achievements", controllers.GetAchievements)
		authorized.GET("/user/:username/achievements", controllers.GetUserAchievements)
		
		// Matchmaking endpoints
		authorized.POST("/queue/join", controllers.JoinQueue)
		authorized.POST("/queue/leave", controllers.LeaveQueue)
	}
	// Store endpoints
	authorized.GET("/store/items", controllers.GetStoreItems)
	authorized.GET("/user/:username/items", controllers.GetUserItems)
	authorized.POST("/store/items/:id/purchase", controllers.PurchaseItem)
	authorized.POST("/store/items/:id/equip", controllers.EquipItem)

	// Admin endpoints
	admin := router.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		// Achievement management
		admin.POST("/achievements", controllers.CreateAchievement)
		admin.PUT("/achievements/:id", controllers.UpdateAchievement)
		admin.DELETE("/achievements/:id", controllers.DeleteAchievement)

		// Server management
		admin.POST("/servers", controllers.RegisterServer)
		admin.PUT("/servers/:id/status", controllers.UpdateServerStatus)

		// Store management
		admin.POST("/store/items", controllers.CreateStoreItem)
		admin.PUT("/store/items/:id", controllers.UpdateStoreItem)
		admin.DELETE("/store/items/:id", controllers.DeleteStoreItem)
	}

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
