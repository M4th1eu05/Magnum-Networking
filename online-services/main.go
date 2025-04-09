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
		authorized.GET("/user/:username/stats", controllers.GetUserStats)
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
