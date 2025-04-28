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

	//utils.GenerateFakeData()

	utils.InitValidators()
	utils.StartMatchmaking()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static("/static", "./static")
	router.LoadHTMLGlob("views/*.html")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "The page you are looking for does not exist.",
			"code":  404,
		})
	})

	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "An unexpected error occurred. Please try again later.",
				"code":  500,
			})
		}
	})

	// Routes
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)

	router.GET("/admin/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_login.html", nil)
	})
	router.POST("/admin/login", controllers.AdminLogin)

	authorized := router.Group("", middlewares.AuthMiddleware())
	{
		authorized.GET("/user/stats", controllers.GetStats)

		authorized.GET("/user/achievements", controllers.GetAchievements)

		// Matchmaking endpoints
		authorized.POST("/queue/join", controllers.JoinQueue)
		authorized.GET("/queue/status", controllers.QueueStatus)
		authorized.POST("/queue/leave", controllers.LeaveQueue)

	}

	server := router.Group("/server")
	{

		server.GET("/register", controllers.RegisterServer)
		server.GET("/startGame", controllers.StartGame)
		server.POST("/endGame", controllers.EndGame)
	}

	// Admin endpoints
	admin := router.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.GET("/dashboard", controllers.AdminDashboard)

		// Achievement management
		admin.POST("/achievements", controllers.CreateAchievement)
		admin.PATCH("/achievements/:id", controllers.UpdateAchievement)
		admin.DELETE("/achievements/:id", controllers.DeleteAchievement)

	}

	s := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
