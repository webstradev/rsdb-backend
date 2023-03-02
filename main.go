package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/controllers"
	"github.com/webstradev/rsdb-backend/db"
	"github.com/webstradev/rsdb-backend/middlewares"
	"github.com/webstradev/rsdb-backend/migrations"
	"github.com/webstradev/rsdb-backend/utils"
)

func main() {
	// Load environment variables
	loadEnvironmentVariables()

	// Set up database instance
	db, err := db.Setup(os.Getenv("DB_CONNECTION_STRING"), migrations.LoadMigrations())
	if err != nil {
		log.Fatal(err)
	}

	// Test database connection
	db.Ping()

	// Migrate Database
	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	jwtService, err := auth.CreateJWTService(os.Getenv("JWT_SIGNING_SECRET"), os.Getenv("JWT_ISSUER"), 24*time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Environment (for dependency injection)
	env := &utils.Environment{
		DB:  db,
		JWT: jwtService,
	}

	// Initialise router
	router := gin.New()

	// Disable logging for health check endpoint
	logger := gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/api/health"}})

	// Use the logger and recovery middleware
	router.Use(logger, gin.Recovery())

	// Health check for k8s
	router.GET("/api/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.POST("/api/v1/login", controllers.Login(env))

	// All the calls to the api group will require authentication
	api := router.Group("/api/v1")
	api.Use(middlewares.JWTAuthMiddleware(env))

	// General
	api.GET("/counts", controllers.GetCounts(env))

	// Platforms
	api.GET("/platforms", middlewares.PaginationMiddleware(), controllers.GetPlatforms(env))
	api.POST("/platforms", controllers.CreatePlatform(env))
	api.GET("/platforms/:id", controllers.GetPlatform(env))
	api.PUT("/platforms/:id", controllers.EditPlatform(env))
	api.DELETE("/platforms/:id", controllers.DeletePlatform(env))

	// Contacts
	api.GET("/platforms/:id/contacts", controllers.GetContacts(env))
	api.PUT("/contacts/:id", controllers.EditContact(env))

	// Server object
	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := s.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Println("Failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln("Server forced to shutdown")
	}

	log.Println("Server exiting.")
}

func loadEnvironmentVariables() {
	// If a database connection string is not yet set in environment variables (or by kube secrets) then load the .env file
	if os.Getenv("DB_CONNECTION_STRING") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
}
