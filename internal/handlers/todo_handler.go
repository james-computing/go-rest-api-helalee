package handlers

import (
	"net/http"
	"strconv"
	"todo_api/internal/models"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	//Id        int    `json: "id"`
	Title     string `json:"title"`
	Completed *bool  `json:"completed"` // *bool works as a nullable bool
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

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(context *gin.Context) {
		todos, err := repository.GetAllTodos(pool)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, todos)
	}
}

func GetTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(context *gin.Context) {
		// Get the id
		var idStr string = context.Param("id")
		var id int
		var err error
		id, err = strconv.Atoi(idStr)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo Id"})
			return
		}

		// Get the todo
		var todo *models.Todo
		todo, err = repository.GetTodoById(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, todo)
	}
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var idStr string = c.Param("id")
		var id int
		var err error
		id, err = strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Todo id"})
			return
		}

		var input UpdateTodoInput
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == "" && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Must provide at least title or completed"})
			return
		}

		var completed bool = false
		if input.Completed != nil {
			completed = *input.Completed
		}

		var todo *models.Todo
		todo, err = repository.UpdateTodo(pool, id, input.Title, completed)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusFound, gin.H{"error": "Todo not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}
