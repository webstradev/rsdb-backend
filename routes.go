package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/webstradev/rsdb-backend/controllers"
	"github.com/webstradev/rsdb-backend/controllers/articles"
	"github.com/webstradev/rsdb-backend/controllers/platforms"
	"github.com/webstradev/rsdb-backend/controllers/projects"
	"github.com/webstradev/rsdb-backend/controllers/users"
	"github.com/webstradev/rsdb-backend/middlewares"
	"github.com/webstradev/rsdb-backend/utils"
)

func registerRoutes(env *utils.Environment) *gin.Engine {
	// Initialise router
	router := gin.New()

	// Disable logging for health check endpoint
	logger := gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/api/health"}})

	// Use the logger and recovery middleware
	router.Use(logger, gin.Recovery())

	// CORS Setup
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"localhost", "https://dev.rsdb.webstra.dev", "https://rsdb.webstra.dev"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
		MaxAge:        12 * time.Hour,
	}))

	// Health check for k8s
	router.GET("/api/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Users (unauthenticated)
	router.POST("/api/v1/login", users.Login(env))
	router.POST("api/v1/users/register", users.Register(env))

	// All the calls to the api group will require authentication
	api := router.Group("/api/v1")
	api.Use(middlewares.JWTAuthMiddleware(env))

	// General
	api.GET("/counts", controllers.GetCounts(env))

	// Platforms
	api.GET("/platforms", middlewares.PaginationMiddleware(), platforms.GetPlatforms(env))
	api.POST("/platforms", platforms.CreatePlatform(env))
	api.GET("/platforms/:platformId", platforms.GetPlatform(env))
	api.PUT("/platforms/:platformId", platforms.EditPlatform(env))
	api.DELETE("/platforms/:platformId", platforms.DeletePlatform(env))

	// Contacts
	api.GET("/platforms/:platformId/contacts", platforms.GetContacts(env))
	api.POST("/platforms/:platformId/contacts", platforms.CreateContact(env))
	api.PUT("/platforms/:platformId/contacts/:id", platforms.EditContact(env))
	api.DELETE("/platforms/:platformId/contacts/:id", platforms.DeleteContact(env))

	// Articles
	api.GET("/articles", middlewares.PaginationMiddleware(), articles.GetArticles(env))
	api.POST("/articles", articles.CreateArticle(env))
	api.GET("/articles/:articleId", articles.GetArticle(env))
	api.PUT("/articles/:articleId", articles.EditArticle(env))
	api.DELETE("/articles/:articleId", articles.DeleteArticle(env))

	// Projects
	api.GET("/projects", middlewares.PaginationMiddleware(), projects.GetProjects(env))
	api.POST("/projects", projects.CreateProject(env))
	api.GET("/projects/:projectId", projects.GetProject(env))
	api.PUT("/projects/:projectId", projects.EditProject(env))
	api.DELETE("/projects/:projectId", projects.DeleteProject(env))

	// Users (authenticated)
	api.PUT("/users/password", users.EditPassword(env))

	// Admin Routes
	admin := api.Group("/admin")
	admin.Use(middlewares.AdminAuthMiddleware())

	// Users (admin)
	admin.GET("/users/token", users.GetRegistrationToken(env))
	admin.GET("/users/:userId/resettoken", users.GetPasswordResetToken(env))

	return router
}
