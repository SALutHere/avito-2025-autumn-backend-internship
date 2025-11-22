package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
)

type StatsPostgres struct {
	db *sql.DB
}

func NewStatsPostgres(db *sql.DB) repository.StatsRepository {
	return &StatsPostgres{db: db}
}

func (r *StatsPostgres) CountAssignmentsByUser(ctx context.Context) ([]domain.UserAssignmentStat, error) {
	log := logger.L()

	q := `
        SELECT reviewer_id, COUNT(*) AS assignments
        FROM pull_request_reviewers
        GROUP BY reviewer_id
        ORDER BY assignments DESC
    `
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		log.Error("failed stats query", slog.String("query", q), slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	var stats []domain.UserAssignmentStat
	for rows.Next() {
		var s domain.UserAssignmentStat
		if err := rows.Scan(&s.UserID, &s.Assignments); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

func (r *StatsPostgres) CountReviewersByPR(ctx context.Context) ([]domain.PRReviewerStat, error) {
	log := logger.L()

	q := `
        SELECT pr_id, COUNT(*) AS reviewers
        FROM pull_request_reviewers
        GROUP BY pr_id
        ORDER BY reviewers DESC
    `
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		log.Error("failed stats query", slog.String("query", q), slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	var stats []domain.PRReviewerStat
	for rows.Next() {
		var s domain.PRReviewerStat
		if err := rows.Scan(&s.PRID, &s.Reviewers); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}
