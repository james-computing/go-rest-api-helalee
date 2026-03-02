package main

import (
	"log"
	"net/http"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"
	"todo_api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.ConnectionString)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	// Public routes
	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleWare(cfg))

	// Protected routes
	protected.POST("", handlers.CreateTodoHandler(pool))
	protected.GET("", handlers.GetAllTodosHandler(pool))
	protected.GET("/:id", handlers.GetTodoByIdHandler(pool))
	protected.PUT("/:id", handlers.UpdateTodoHandler(pool))
	protected.DELETE("/:id", handlers.DeleteTodoHandler(pool))

	// Middleware test route
	router.GET("/protected-test", middleware.AuthMiddleWare(cfg), handlers.TestProtectedHandler())

	router.Run(":" + cfg.Port)
}
