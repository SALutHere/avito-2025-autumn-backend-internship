package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
)

type TeamPostgres struct {
	db *sql.DB
}

func NewTeamPostgres(db *sql.DB) repository.TeamRepository {
	return &TeamPostgres{db: db}
}

func (r *TeamPostgres) Create(ctx context.Context, team *domain.Team) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO teams (name) VALUES ($1)
    `, team.Name)
	return err
}

func (r *TeamPostgres) ExistsByName(ctx context.Context, name string) (bool, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT 1 FROM teams WHERE name = $1
    `, name)

	var dummy int
	err := row.Scan(&dummy)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *TeamPostgres) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT name FROM teams WHERE name = $1
    `, name)

	var t domain.Team
	if err := row.Scan(&t.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, err
	}

	return &t, nil
}
