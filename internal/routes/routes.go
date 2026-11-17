package routes

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/handlers"
	"course-selection-platform/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Add this blank import so the generated docs.go init() function runs
	_ "course-selection-platform/docs"
	"net/http"
	"strings"
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
		api.DELETE("/user/:id",
			middleware.RequireAuth(app.JWTSecret),
			middleware.AdminOnly(),
			handlers.DeleteUserHandler(app),
		)

		api.GET("/subjects", handlers.GetAllSubjectsHandler(app.DB))
		api.GET("/faculties", handlers.GetFacultiesHandler(app))

		api.GET("/student/choosable-subjects", 
			middleware.RequireAuth(app.JWTSecret), 
			handlers.GetStudentChoosableSubjectsHandler(app),
		)		
		api.GET("/priority",
			middleware.RequireAuth(app.JWTSecret),
			handlers.GetStudentPrioritiesHandler(app.DB),
		)
		api.POST("/priority",
			middleware.RequireAuth(app.JWTSecret),
			handlers.AddStudentPriorityHandler(app.DB),
		)
		api.DELETE("/priority/:id",
			middleware.RequireAuth(app.JWTSecret),
			handlers.DeleteStudentPriorityHandler(app.DB),
		)

		api.GET("/manager/subject-student-overview", 
			middleware.RequireAuth(app.JWTSecret),
			middleware.ManagerOnly(),
			handlers.GetSubjectStudentOverviewHandler(app),
		)

		api.POST("/login", handlers.LoginHandler(app))
		api.POST("/logout", handlers.LogoutHandler(app))
	}

	// Access to Swagger only when in Dev mode.
	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}
