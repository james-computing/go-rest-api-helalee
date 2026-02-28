package handlers

import (
	"net/http"
	"strings"
	"time"
	"todo_api/internal/config"
	"todo_api/internal/models"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerRequest RegisterRequest
		var err error
		err = c.BindJSON(&registerRequest)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(registerRequest.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must have at least 6 characters"})
			return
		}

		var bytes []byte = []byte(registerRequest.Password)
		var hashedPasswordBytes []byte
		hashedPasswordBytes, err = bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
		if err != nil {
			// I think the client shouldn't receive this information...
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		var user models.User
		user.Email = registerRequest.Email
		user.Password = string(hashedPasswordBytes)

		var createdUser *models.User
		createdUser, err = repository.CreateUser(pool, &user)

		if err != nil {
			// Not good to compare to a string...
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	}
}

func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		var err error
		err = c.BindJSON(&loginRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user *models.User
		user, err = repository.GetUserByEmail(pool, loginRequest.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		claims := jwt.MapClaims{
			"user_id": user.Id,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		var tokenString string
		tokenString, err = token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
			return
		}

		var loginResponse LoginResponse = LoginResponse{
			AccessToken: tokenString,
		}
		c.JSON(http.StatusOK, loginResponse)
	}
}
