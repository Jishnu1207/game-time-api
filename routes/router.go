package routes

import (
	"game-time-api/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
		}

		// User routes
		users := api.Group("/users")
		{
			users.POST("/register", handlers.Register)
		}
	}

	return router
}
