package handlers

import (
	"net/http"
	"todo_api/internal/models"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
			if err.Error() != "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	}
}
