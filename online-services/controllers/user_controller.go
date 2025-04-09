package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"online-services/database"
	"online-services/middlewares"
	"online-services/models"
)

type Info struct {
	Username string `form:"username" json:"username" binding:"required,notblank"`
	Password string `form:"password" json:"password" binding:"required,notblank"`
}

func Login(c *gin.Context) {
	// Get the username and password from the request
	var loginInfo Info

	c.Bind(&loginInfo)

	// Check if the user exists in the database
	var user models.User
	if err := database.DB.Where("username = ?", loginInfo.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate a JWT token
	token, err := middlewares.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register(c *gin.Context) {
	// Get the username and password from the request
	var registerInfo Info
	c.Bind(&registerInfo)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Check if the username already exists
	var existingUser models.User
	if err := database.DB.Where("username = ?", registerInfo.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Create a new user
	user := models.User{Username: registerInfo.Username, Password: string(hashedPassword)}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate a JWT token
	token, err := middlewares.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GetUserStats(c *gin.Context) {
	username := c.Param("username")

	// Check if the user exists in the database
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
