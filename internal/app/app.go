package app

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/migrate"
	_ "github.com/lib/pq"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/config"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
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
	if err := migrate.Run(db); err != nil {
		log.Error("migration failed", slog.Any("err", err))
	}

	log.Info("Connected to PostgreSQL")

	// TODO: init repositories

	// TODO: init services

	// TODO: init controllers

	// TODO: init main router

	// TODO: graceful shutdown
}
