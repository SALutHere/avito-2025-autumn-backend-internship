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

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) repository.UserRepository {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetByID(ctx context.Context, id string) (*domain.User, error) {
	log := logger.L()

	q := `
        SELECT id, username, team_name, is_active
        FROM users
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, q, id)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}

	return &u, nil
}

func (r *UserPostgres) ListActiveByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	log := logger.L()

	q := `
        SELECT id, username, team_name, is_active
        FROM users
        WHERE team_name = $1 AND is_active = TRUE
    `
	rows, err := r.db.QueryContext(ctx, q, teamName)
	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}
	defer rows.Close()

	var list []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		list = append(list, u)
	}

	return list, nil
}

func (r *UserPostgres) ListByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	log := logger.L()

	q := `
        SELECT id, username, team_name, is_active
        FROM users
        WHERE team_name = $1
    `
	rows, err := r.db.QueryContext(ctx, q, teamName)
	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}
	defer rows.Close()

	var list []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		list = append(list, u)
	}

	return list, nil
}

func (r *UserPostgres) SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error) {
	log := logger.L()

	q := `
        UPDATE users SET is_active = $2 WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, q, id, isActive)
	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *UserPostgres) Upsert(ctx context.Context, user *domain.User) error {
	log := logger.L()

	q := `
        INSERT INTO users (id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE SET
            username = EXCLUDED.username,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active
    `
	_, err := r.db.ExecContext(ctx, q,
		user.ID, user.Username, user.TeamName, user.IsActive,
	)

	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
	}

	return err
}
