//package middlewares
//
//import (
//	"net/http"
//	"online-services/database"
//	"online-services/models"
//
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt"
//)
//
//func AdminMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		tokenString := c.GetHeader("Authorization")
//		if tokenString == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
//			c.Abort()
//			return
//		}
//
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, jwt.ErrSignatureInvalid
//			}
//			return jwtSecret, nil
//		})
//
//		if err != nil || !token.Valid {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
//			c.Abort()
//			return
//		}
//
//		claims := token.Claims.(jwt.MapClaims)
//		username := claims["username"].(string)
//
//		// Check if user is an admin
//		var user models.User
//		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
//			c.Abort()
//			return
//		}
//
//		// Add admin field to user model
//		if !user.IsAdmin {
//			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
//			c.Abort()
//			return
//		}
//
//		c.Next()
//	}
//}
