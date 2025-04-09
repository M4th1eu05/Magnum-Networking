package tests

import (
	"net/http"
	"net/http/httptest"
	"online-services/middlewares"
	"strings"
	"testing"

	"online-services/controllers"
	"online-services/database"
	"online-services/models"
	"online-services/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	utils.InitValidators()
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.GET("/user/:username/stats", middlewares.AuthMiddleware(), controllers.GetUserStats)
	return r
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	// Run the tests
	m.Run()
}

func TestLoginSuccess(t *testing.T) {

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	router := setupRouter()

	// Perform the request
	w := httptest.NewRecorder()
	reqBody := `username=testuser&password=password123`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLoginInvalidCredentials(t *testing.T) {
	router := setupRouter()

	// Perform the request
	w := httptest.NewRecorder()
	reqBody := `username=invaliduser&password=wrongpassword`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

func TestRegisterSuccess(t *testing.T) {
	router := setupRouter()

	// Perform the request
	w := httptest.NewRecorder()
	reqBody := `username=newuser&password=newpassword123`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")

	// Verify user is created in the database
	var user models.User
	err := database.DB.Where("username = ?", "newuser").First(&user).Error
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Password)
}

func TestRegisterDuplicateUser(t *testing.T) {
	// Create a test user
	testUser := models.User{Username: "duplicateuser", Password: "password123"}
	database.DB.Create(&testUser)

	router := setupRouter()

	// Perform the request
	w := httptest.NewRecorder()
	reqBody := `username=duplicateuser&password=newpassword123`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegisterInvalidInput(t *testing.T) {
	router := setupRouter()

	// Perform the request
	w := httptest.NewRecorder()
	reqBody := `username=&password=`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserStatsSuccess(t *testing.T) {
	// Créer un utilisateur de test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	// Générer un token JWT valide
	token, _ := middlewares.GenerateToken(testUser)

	router := setupRouter()

	// Effectuer la requête
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/testuser/stats", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "score")
}

func TestGetUserStatsUnauthorized(t *testing.T) {
	router := setupRouter()

	// Effectuer la requête sans token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/testuser/stats", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserStatsUserNotFound(t *testing.T) {
	// Créer un utilisateur de test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	// Générer un token JWT valide
	token, _ := middlewares.GenerateToken(testUser)

	router := setupRouter()

	// Effectuer la requête pour un utilisateur inexistant
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/unknownuser/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
