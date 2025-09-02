package server

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"

	"medika-backend/internal/application/user"
	"medika-backend/internal/infrastructure/config"
	"medika-backend/internal/infrastructure/persistence/repositories"
	"medika-backend/internal/presentation/http/handlers"
	"medika-backend/internal/presentation/http/middleware"
	"medika-backend/pkg/logger"
)

type Server struct {
	app    *fiber.App
	config config.ServerConfig
	logger logger.Logger
}

func New(
	cfg config.ServerConfig,
	db *bun.DB,
	redis *redis.Client,
	logger logger.Logger,
) *Server {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Prefork:      cfg.Prefork,
		ErrorHandler: middleware.ErrorHandler,
	})

	// Initialize dependencies
	validator := validator.New()
	
	// Repositories
	userRepo := repositories.NewUserRepository(db)
	
	// Application services
	userService := user.NewService(userRepo, nil, logger) // eventBus would be injected
	
	// Handlers
	userHandler := handlers.NewUserHandler(userService, validator, logger)

	// Setup middleware
	setupMiddleware(app)
	
	// Setup routes
	setupRoutes(app, userHandler)

	return &Server{
		app:    app,
		config: cfg,
		logger: logger,
	}
}

func setupMiddleware(app *fiber.App) {
	// Security middleware
	app.Use(helmet.New())
	
	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Configure appropriately for production
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	
	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))
	
	// Recovery middleware
	app.Use(recover.New())
	
	// Simple status endpoint for development
	app.Get("/monitor", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "running",
			"title":  "Medika API Monitor",
		})
	})
}

func setupRoutes(app *fiber.App, userHandler *handlers.UserHandler) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		})
	})

	// API routes
	api := app.Group("/api/v1")
	
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", userHandler.Login)
	auth.Post("/register", userHandler.CreateUser)
	
	// User routes
	users := api.Group("/users")
	users.Get("/", userHandler.GetUsersByOrganization)
	users.Post("/", userHandler.CreateUser)
	users.Get("/me", middleware.AuthRequired(), userHandler.GetCurrentUser)
	users.Get("/:id", middleware.AuthRequired(), userHandler.GetUser)
	users.Put("/:id/profile", middleware.AuthRequired(), userHandler.UpdateUserProfile)
	users.Put("/:id/medical-info", middleware.AuthRequired(), userHandler.UpdateMedicalInfo)
	users.Put("/:id/avatar", middleware.AuthRequired(), userHandler.UpdateAvatar)
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	s.logger.Info(ctx, "🚀 Starting server", "address", addr)
	
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info(ctx, "Shutting down server...")
	return s.app.ShutdownWithContext(ctx)
}
