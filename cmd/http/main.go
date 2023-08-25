package main

import (
	"log/slog"
	"os"

	"github.com/bagashiz/go-pos/internal/adapter/handler"
	"github.com/bagashiz/go-pos/internal/adapter/postgres"
	"github.com/bagashiz/go-pos/internal/adapter/postgres/repository"
	"github.com/bagashiz/go-pos/internal/core/service"
	"github.com/joho/godotenv"
)

func init() {
	// Init logger
	var logHandler *slog.JSONHandler

	env := os.Getenv("APP_ENV")
	if env == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})

		// Load .env file
		err := godotenv.Load()
		if err != nil {
			slog.Error("Error loading .env file", "error", err)
			os.Exit(1)
		}
	}

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

func main() {
	appName := os.Getenv("APP_NAME")
	env := os.Getenv("APP_ENV")
	url := os.Getenv("APP_URL")
	port := os.Getenv("APP_PORT")
	listenAddr := url + ":" + port

	slog.Info("Starting the application", "app", appName, "env", env)

	// Init database
	db, err := postgres.NewDB()
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Successfully connected to the database")

	// Dependency injection
	// User
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Init router
	router := handler.NewRouter()
	router.InitRoutes(*userHandler)

	// Start server
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	router.Start(listenAddr)
}
