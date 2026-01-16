package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/momokii/go-cli-notes/internal/api/handler"
	"github.com/momokii/go-cli-notes/internal/api/middleware"
	"github.com/momokii/go-cli-notes/internal/api/router"
	"github.com/momokii/go-cli-notes/internal/config"
	"github.com/momokii/go-cli-notes/internal/repository"
	"github.com/momokii/go-cli-notes/internal/service"
	"github.com/momokii/go-cli-notes/internal/util"
)

const (
	API_VERSION = "1.1.0"
)

func main() {
	// Try to load .env file if it exists (silent fail if not found)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup logger
	setupLogger(cfg)

	slog.Info("Starting Knowledge Garden API", "version", API_VERSION)

	// Connect to database
	db, err := repository.NewDB(
		cfg.Database.DSN(),
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Log database connection info
	if cfg.Database.DatabaseURL != "" {
		slog.Info("Connected to database", "connection", "DATABASE_URL")
	} else {
		slog.Info("Connected to database", "host", cfg.Database.Host, "port", cfg.Database.Port)
	}

	// Initialize repositories
	repos := repository.NewRepository(db)

	// Initialize utilities
	hasher := util.NewPasswordHasher()
	jwtManager := util.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiration,
		cfg.JWT.RefreshExpiration,
	)
	linkParser := util.NewLinkParser()

	// Initialize services
	authService := service.NewAuthService(repos.User, repos.RefreshToken, hasher, jwtManager)
	noteService := service.NewNoteService(repos.Note, repos.Tag, repos.Link, repos.Activity, linkParser)
	tagService := service.NewTagService(repos.Tag, repos.Note, repos.Activity)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Knowledge Garden API " + API_VERSION,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.RequestID())

	// Setup handlers
	handlers := &handler.Handlers{
		Auth:     handler.NewAuthHandler(authService),
		Note:     handler.NewNoteHandler(noteService),
		Tag:      handler.NewTagHandler(tagService),
		Search:   handler.NewSearchHandler(noteService),
		Link:     handler.NewLinkHandler(noteService),
		Activity: handler.NewActivityHandler(repos.Activity, noteService),
	}

	// Setup routes
	router.Setup(app, handlers, jwtManager)

	// Start server in goroutine
	go func() {
		slog.Info("Server listening", "address", cfg.Server.Address())
		if err := app.Listen(cfg.Server.Address()); err != nil {
			slog.Error("Server failed", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Graceful shutdown
	if err := app.ShutdownWithContext(context.Background()); err != nil {
		slog.Error("Shutdown error", "error", err)
	}

	slog.Info("Server stopped")
}

// setupLogger configures the structured logger
func setupLogger(cfg *config.Config) {
	var handler slog.Handler
	logLevel := parseLogLevel(cfg.Log.Level)

	if cfg.Log.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// parseLogLevel parses log level string
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// customErrorHandler handles errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
