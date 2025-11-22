package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/config"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/controller/http/v1"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/controller/http/v1/routers"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/migrate"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository/postgres"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/service"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"

	_ "github.com/lib/pq"
)

func Run(configPath string) {
	// Loading configuration
	cfg := config.Load(configPath)

	// Setting up logger
	log := logger.Setup(cfg.Env)
	log.Info("starting application...")

	// Initializing database connection
	db, err := sql.Open("postgres", cfg.PostgresURL())
	if err != nil {
		log.Error("can not open postgres db", slog.Any("err", err))
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		log.Error("can not connect to the postgres db", slog.Any("err", err))
		os.Exit(1)
	}

	// Running migrations
	if err = migrate.Run(db); err != nil {
		log.Error("migration failed", slog.Any("err", err))
		os.Exit(1)
	}

	log.Info("Connected to PostgreSQL")

	// Initializing repositories
	log.Info("Initializing repositories...")
	teamRepo := postgres.NewTeamPostgres(db)
	userRepo := postgres.NewUserPostgres(db)
	prRepo := postgres.NewPRPostgres(db)
	statsRepo := postgres.NewStatsPostgres(db)
	log.Info("Repositories are ready")

	// Initializing services
	log.Info("Initializing services...")
	teamSvc := service.NewTeamService(teamRepo)
	userSvc := service.NewUserService(userRepo, teamRepo)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo)
	statsSvc := service.NewStatsService(statsRepo)
	log.Info("Services are ready")

	// Initializing controllers
	log.Info("Initializing controllers...")
	teamCtrl := routers.NewTeamController(teamSvc, userSvc)
	userCtrl := routers.NewUserController(userSvc, prSvc)
	prCtrl := routers.NewPRController(prSvc)
	statsCtrl := routers.NewStatsController(statsSvc)
	log.Info("Controllers are ready")

	// Initializing router
	log.Info("Initializing router...")
	e := v1.NewHTTPServer(teamCtrl, userCtrl, prCtrl, statsCtrl)
	log.Info("Router is ready")

	// Running server
	log.Info("Server is running", slog.Int("port", cfg.HTTPPort))

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)

	if err = e.Start(addr); err != nil {
		log.Error("server stopped", slog.Any("err", err))
		os.Exit(1)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = e.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown error", slog.Any("err", err))
		os.Exit(1)
	}

	if err = db.Close(); err != nil {
		log.Error("DB close error", slog.Any("err", err))
		os.Exit(1)
	}

	log.Info("Server stopped successfully")
}
