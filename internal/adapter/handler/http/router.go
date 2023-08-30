package handler

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

// NewRouter creates a new HTTP router
func NewRouter() *Router {
	// Disable debug mode and write logs to file in production
	env := os.Getenv("APP_ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)

		logFile, _ := os.Create("gin.log")
		gin.DefaultWriter = io.Writer(logFile)
	}

	router := gin.New()
	router.Use(gin.LoggerWithFormatter(customLogger))
	router.Use(gin.Recovery())

	// CORS
	config := cors.DefaultConfig()
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	originsList := strings.Split(allowedOrigins, ",")
	config.AllowOrigins = originsList
	router.Use(cors.New(config))

	return &Router{
		router,
	}
}

// InitRoutes configures the handler for each route
func (r *Router) InitRoutes(
	userHandler UserHandler,
	paymentHandler PaymentHandler,
) {
	v1 := r.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.GET("/", userHandler.ListUsers)
			user.GET("/:id", userHandler.GetUser)
			user.PUT("/:id", userHandler.UpdateUser)
			user.DELETE("/:id", userHandler.DeleteUser)
		}
		payment := v1.Group("/payments")
		{
			payment.POST("/", paymentHandler.CreatePayment)
			payment.GET("/", paymentHandler.ListPayments)
			payment.GET("/:id", paymentHandler.GetPayment)
			payment.PUT("/:id", paymentHandler.UpdatePayment)
			payment.DELETE("/:id", paymentHandler.DeletePayment)
		}
	}
}

// customLogger is a custom Gin logger
func customLogger(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] - %s \"%s %s %s %d %s [%s]\"\n",
		param.TimeStamp.Format(time.RFC1123),
		param.ClientIP,
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency.Round(time.Millisecond),
		param.Request.UserAgent(),
	)
}
