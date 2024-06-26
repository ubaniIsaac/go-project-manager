package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {

	var secret = []byte(os.Getenv("jwtSecret"))
	type MyClaims struct {
		Role   string `json:"role"`
		UserID string `json:"userID"`
		jwt.RegisteredClaims
	}

	return func(c *gin.Context) {
		// Get the authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if the authorization header is missing
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Check if the authorization header has the correct format (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token from the authorization header
		tokenString := parts[1]

		// Parse and validate the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID from the token claims into the request context for further processing
		claims, ok := token.Claims.(*MyClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract token claims"})
			c.Abort()
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func CheckRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole != role {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
