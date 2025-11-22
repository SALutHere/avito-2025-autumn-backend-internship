package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
)

type TeamPostgres struct {
	db *sql.DB
}

func NewTeamPostgres(db *sql.DB) repository.TeamRepository {
	return &TeamPostgres{db: db}
}

func (r *TeamPostgres) Create(ctx context.Context, team *domain.Team) error {
	log := logger.L()

	q := `
        INSERT INTO teams (name) VALUES ($1)
    `
	_, err := r.db.ExecContext(ctx, q, team.Name)

	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
	}

	return err
}

func (r *TeamPostgres) ExistsByName(ctx context.Context, name string) (bool, error) {
	log := logger.L()

	q := `
        SELECT 1 FROM teams WHERE name = $1
    `
	row := r.db.QueryRowContext(ctx, q, name)

	var dummy int
	err := row.Scan(&dummy)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return false, err
	}
	return true, nil
}

func (r *TeamPostgres) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	log := logger.L()

	q := `
        SELECT name FROM teams WHERE name = $1
    `
	row := r.db.QueryRowContext(ctx, q, name)

	var t domain.Team
	if err := row.Scan(&t.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}

	return &t, nil
}
