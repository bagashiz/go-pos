package main

import (
	"context"
	"log/slog"
	"os"

	_ "github.com/bagashiz/go-pos/docs"
	"github.com/bagashiz/go-pos/internal/adapter/auth/paseto"
	"github.com/bagashiz/go-pos/internal/adapter/handler/http"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres/repository"
	"github.com/bagashiz/go-pos/internal/adapter/storage/redis"
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

//	@title						Go POS (Point of Sale) API
//	@version					1.0
//	@description				This is a simple RESTful Point of Sale (POS) Service API written in Go using Gin web framework, PostgreSQL database, and Redis cache.
//
//	@contact.name				Bagas Hizbullah
//	@contact.url				https://github.com/bagashiz/go-pos
//	@contact.email				bagash.office@simplelogin.com
//
//	@license.name				MIT
//	@license.url				https://github.com/bagashiz/go-pos/blob/main/LICENSE
//
//	@host						gopos.bagashiz.me
//	@BasePath					/v1
//	@schemes					http https
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and the access token.
func main() {
	appName := os.Getenv("APP_NAME")
	env := os.Getenv("APP_ENV")
	dbConn := os.Getenv("DB_CONNECTION")
	tokenSymmetricKey := os.Getenv("TOKEN_SYMMETRIC_KEY")
	httpUrl := os.Getenv("HTTP_URL")
	httpPort := os.Getenv("HTTP_PORT")
	listenAddr := httpUrl + ":" + httpPort

	slog.Info("Starting the application", "app", appName, "env", env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Successfully connected to the database", "db", dbConn)

	// Init cache service
	cache, err := redis.New(ctx)
	if err != nil {
		slog.Error("Error initializing cache connection", "error", err)
		os.Exit(1)
	}
	defer cache.Close()

	slog.Info("Successfully connected to the cache server")

	// Init token service
	token, err := paseto.New(tokenSymmetricKey)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	// Dependency injection
	// User
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cache)
	userHandler := http.NewUserHandler(userService)

	// Auth
	authService := service.NewAuthService(userRepo, token)
	authHandler := http.NewAuthHandler(authService)

	// Payment
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, cache)
	paymentHandler := http.NewPaymentHandler(paymentService)

	// Category
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo, cache)
	categoryHandler := http.NewCategoryHandler(categoryService)

	// Product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo, categoryRepo, cache)
	productHandler := http.NewProductHandler(productService)

	// Order
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, productRepo, categoryRepo, userRepo, paymentRepo, cache)
	orderHandler := http.NewOrderHandler(orderService)

	// Init router
	router, err := http.NewRouter(
		token,
		*userHandler,
		*authHandler,
		*paymentHandler,
		*categoryHandler,
		*productHandler,
		*orderHandler,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
