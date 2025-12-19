package main

import (
	"context"
	"fmt"
	"os"

	"searchav/internal/config"
	"searchav/internal/handler"
	"searchav/internal/service"
	"searchav/internal/source"

	_ "searchav/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// @title SearchAV API
// @version 1.0
// @description Video aggregation search service
// @host localhost:9898
// @BasePath /api
func main() {
	fx.New(
		// Configuration
		fx.Provide(config.New),

		// Logger
		fx.Provide(NewLogger),

		// Source client
		fx.Provide(source.NewClient),

		// Service layer
		fx.Provide(service.NewSearchService),
		fx.Provide(service.NewDetailService),

		// Handlers
		fx.Provide(handler.NewContextHandler),
		fx.Provide(handler.NewSearchHandler),
		fx.Provide(handler.NewDetailHandler),

		// Fiber App
		fx.Provide(NewFiberApp),

		// Start
		fx.Invoke(RegisterRoutes),
		fx.Invoke(StartServer),
	).Run()
}

// NewLogger creates a logger instance
func NewLogger(cfg *config.Config) *zerolog.Logger {
	var logger zerolog.Logger

	if cfg.Log.Format == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stdout).
			With().Timestamp().Logger()
	}

	level, _ := zerolog.ParseLevel(cfg.Log.Level)
	logger = logger.Level(level)

	return &logger
}

// NewFiberApp creates a Fiber application
func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())

	return app
}

// RegisterRoutes registers all routes
func RegisterRoutes(
	app *fiber.App,
	cfg *config.Config,
	ctxHandler *handler.ContextHandler,
	searchHandler *handler.SearchHandler,
	detailHandler *handler.DetailHandler,
) {
	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Swagger docs
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes with auth middleware
	api := app.Group("/api", handler.AuthMiddleware(cfg))
	api.Get("/search", ctxHandler.Wrap(searchHandler.Search))
	api.Get("/detail", ctxHandler.Wrap(detailHandler.GetDetail))
}

// StartServer starts the HTTP server
func StartServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, logger *zerolog.Logger) {
	// Log config on startup
	logger.Info().
		Bool("auth_enabled", cfg.Auth.Enabled).
		Int("passwords_count", len(cfg.Auth.Passwords)).
		Int("sources_count", len(cfg.Sources)).
		Msg("config loaded")

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
			logger.Info().Str("addr", addr).Msg("server starting")
			go func() {
				if err := app.Listen(addr); err != nil {
					logger.Fatal().Err(err).Msg("server start failed")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("server stopping")
			return app.Shutdown()
		},
	})
}
