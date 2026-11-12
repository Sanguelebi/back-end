package routes

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/handlers"
	"course-selection-platform/internal/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// The main SetupRouter function.
func SetupRouter(app *config.AppData) *gin.Engine {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "To access the api, add /api and the endpoint you are looking for.",
			})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
	})

	{
		api := r.Group("/api")
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Pong!",
			})
		})

		api.GET("/users",
			middleware.RequireAuth(app.JWTSecret),
			middleware.AdminOnly(),
			handlers.GetAllUsersHandler(app),
		)
		api.POST("/user",
			middleware.RequireAuth(app.JWTSecret),
			middleware.AdminOnly(),
			handlers.AddUserHandler(app),
		)
		api.PATCH("/user/:id",
			middleware.RequireAuth(app.JWTSecret),
			middleware.AdminOnly(),
			handlers.EditUserHandler(app),
		)

		api.GET("/subjects", handlers.GetSubjects(app.DB))

		api.GET("/priority",
			middleware.RequireAuth(app.JWTSecret),
			handlers.GetPriorities(app.DB),
		)
		api.POST("/priority",
			middleware.RequireAuth(app.JWTSecret),
			handlers.AddPriority(app.DB),
		)
		api.DELETE("/priority/:id",
			middleware.RequireAuth(app.JWTSecret),
			handlers.DeletePriority(app.DB),
		)

		api.POST("/login", handlers.LoginHandler(app))
	}

	return r
}
