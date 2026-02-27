package handlers

import (
	"net/http"
	"todo_api/internal/models"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(context *gin.Context) {
		var input CreateTodoInput
		var err error = context.ShouldBindJSON(&input)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var todo *models.Todo
		todo, err = repository.CreateTodo(pool, input.Title, input.Completed)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		context.JSON(http.StatusCreated, todo)
	}
}
