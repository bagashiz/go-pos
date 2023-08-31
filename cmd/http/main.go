package main

import (
	"log/slog"
	"os"

	handler "github.com/bagashiz/go-pos/internal/adapter/handler/http"
	repo "github.com/bagashiz/go-pos/internal/adapter/repository/postgres"
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

	slog.Info("Starting the application...", "app", appName, "env", env)

	// Init database
	slog.Info("Connecting to the database...")
	db, err := repo.NewDB()
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		slog.Error("Error connecting database", "error", err)
		os.Exit(1)
	} else {
		slog.Info("Successfully connected to the database")
	}

	// Dependency injection
	// User
	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Payment
	paymentRepo := repo.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Category
	categoryRepo := repo.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Product
	productRepo := repo.NewProductRepository(db)
	productService := service.NewProductService(productRepo, categoryRepo)
	productHandler := handler.NewProductHandler(productService)

	// Init router
	router := handler.NewRouter()
	router.InitRoutes(
		*userHandler,
		*paymentHandler,
		*categoryHandler,
		*productHandler,
	)

	// Start server
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Run(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
