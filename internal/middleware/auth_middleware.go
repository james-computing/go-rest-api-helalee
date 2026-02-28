package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"todo_api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader string = c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		var tokenString string = strings.TrimPrefix(authHeader, "Bearer ")
		// First case is when the header had "Bearer " as prefix, but nothing more.
		// Second case is when it didn't have "Bearer " as prefix, so nothing was trimmed.
		if tokenString == "" || tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodES256.Alg() {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		var claims jwt.MapClaims
		var ok bool
		claims, ok = token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		var userId string
		userId, ok = claims["user_id"].(string)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		var exp float64
		exp, ok = claims["exp"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		expirationTime := time.Unix(int64(exp), 0)

		if time.Now().After(expirationTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		c.Set("user_id", userId)
		c.Next()
	}
}
