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
	db, err := db.Setup(os.Getenv("DB_CONNECTION_STRING"), &migrations.SQLMigration)
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
	router := gin.Default()

	// Health check for k8s
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.POST("/api/login", controllers.Login(env))

	api := router.Group("/api")
	api.Use(middlewares.JWTAuthMiddleware(env))

	// Server object
	s := &http.Server{
		Addr:         "127.0.0.1:8080",
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
