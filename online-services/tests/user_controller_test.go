package tests

import (
	"net/http"
	"net/http/httptest"
	"online-services/controllers"
	"strings"
	"testing"

	"online-services/database"
	"online-services/models"
	"online-services/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginSuccess(t *testing.T) {

	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	router := gin.Default()
	router.POST("/login", controllers.Login)

	w := httptest.NewRecorder()
	reqBody := `username=testuser&password=password123`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLoginInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	router := gin.Default()
	utils.InitValidators()
	router.POST("/login", controllers.Login)

	w := httptest.NewRecorder()
	reqBody := `username=invaliduser&password=wrongpassword`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

func TestRegisterSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	router := gin.Default()
	utils.InitValidators()
	router.POST("/register", controllers.Register)

	w := httptest.NewRecorder()
	reqBody := `username=newuser&password=newpassword123`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")

	var user models.User
	err := database.DB.Where("username = ?", "newuser").First(&user).Error
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Password)
}

func TestRegisterDuplicateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	testUser := models.User{Username: "duplicateuser", Password: "password123"}
	database.DB.Create(&testUser)

	router := gin.Default()
	utils.InitValidators()
	router.POST("/register", controllers.Register)

	w := httptest.NewRecorder()
	reqBody := `username=duplicateuser&password=newpassword123`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegisterInvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	router := gin.Default()
	utils.InitValidators()
	router.POST("/register", controllers.Register)

	w := httptest.NewRecorder()
	reqBody := `username=&password=`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
