package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
	"net/http"
	"online-services/database"
	"online-services/middlewares"
	"online-services/models"
)

type UserInfo struct {
	Username string `form:"username" json:"username" binding:"required,notblank"`
	Password string `form:"password" json:"password" binding:"required,notblank"`
}

func Login(c *gin.Context) {
	// Get the username and password from the request
	var loginInfo UserInfo

	err := c.Bind(&loginInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

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
	var registerInfo UserInfo
	err := c.Bind(&registerInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

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

func AdminLogin(c *gin.Context) {
	var loginInfo UserInfo

	if err := c.Bind(&loginInfo); err != nil {
		c.HTML(http.StatusBadRequest, "admin_login.html", gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ? AND role = ?", loginInfo.Username, "admin").First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := middlewares.GenerateToken(user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_login.html", gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", token, 3600, "/", "", true, true)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func AdminDashboard(c *gin.Context) {
	var servers []models.GameServer
	if err := database.DB.Preload(clause.Associations).Preload("CurrentGame.Players").Find(&servers).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_dashboard.html", gin.H{"error": "Failed to load games"})
		return
	}

	var achievements []models.Achievement
	if err := database.DB.Find(&achievements).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_dashboard.html", gin.H{"error": "Failed to load achievements"})
		return
	}

	var users []models.User
	if err := database.DB.Preload(clause.Associations).Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_dashboard.html", gin.H{"error": "Failed to load users"})
		return
	}

	serverJSON, _ := json.Marshal(servers)
	achievementsJSON, _ := json.Marshal(achievements)
	usersJSON, _ := json.Marshal(users)
	token, _ := c.Cookie("token")

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"servers":     string(serverJSON),
		"achivements": string(achievementsJSON),
		"users":       string(usersJSON),
		"token":       token,
	})
}
