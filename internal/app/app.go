package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"backend/internal/controller/http/middleware"
	repocat "backend/internal/storage/postgres/cat"
	repomission "backend/internal/storage/postgres/mission"
	repotarget "backend/internal/storage/postgres/target"

	svccat "backend/internal/service/cat"
	svcmission "backend/internal/service/mission"

	handlercat "backend/internal/controller/http/v1/cat"
	handlermission "backend/internal/controller/http/v1/mission"

	"backend/config"
	"backend/internal/entity/cat"
	"backend/pkg/httpserver"
	"backend/pkg/postgres"
	"backend/pkg/validator/breed"
	structvalidator "backend/pkg/validator/struct"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg config.Config) {
	logLevel := slog.LevelInfo
	if cfg.Server.IsDev {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := postgres.New(ctx,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Server.IsDev,
	)
	if err != nil {
		logger.Error("unable to connect to PostgreSQL", "err", err)
		return
	}

	logger.Info("connected to PostgreSQL", "db", cfg.Postgres.DBName)

	defer func() {
		if err := client.Close(); err != nil {
			logger.Error("error disconnecting PostgreSQL", "err", err)
		}
	}()

	if err = runMigrations(client); err != nil {
		logger.Error("db migrations failed", "err", err)
		return
	}

	breedValidator := breed.NewValidator()
	validator := structvalidator.NewValidator()

	catRepo := repocat.NewRepo(client)
	missionRepo := repomission.NewRepo(client)
	targetRepo := repotarget.NewRepo(client)

	catSvc := svccat.NewService(catRepo, breedValidator, logger)
	missionSvc := svcmission.NewService(
		missionRepo,
		targetRepo,
		catRepo,
		logger,
	)

	// HTTP server

	mw := middleware.NewMiddleware(logger)
	if !cfg.Server.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()

	// Could use either gin's logger, or customer logger middleware
	g.Use(
		// gin.Logger(), gin.Recovery(),
		mw.Logger(), mw.Recovery(),
	)

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	handlercat.InitHandler(
		g, logger,
		catSvc,
		validator,
	)

	handlermission.InitHandler(
		g, logger,
		missionSvc,
		validator,
	)

	server := httpserver.New(
		g,
		httpserver.Port(cfg.Server.Port),
	)

	logger.Info("starting http server", "port", cfg.Server.Port)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Info("shutdown signal received", "signal", sig.String())
	case err := <-server.Notify():
		logger.Error("HTTP server error", "err", err)
	}

	if err := server.Shutdown(); err != nil {
		logger.Error("server shutdown error", "err", err)
	} else {
		logger.Info("server stopped gracefully")
	}
}

func runMigrations(client *postgres.Postgres) error {
	slog.Info("running migrations...")
	return client.Instance().AutoMigrate(
		&cat.Cat{},
		&cat.Mission{},
		&cat.Target{},
	)
}
