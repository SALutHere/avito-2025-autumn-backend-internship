package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
)

type PRPostgres struct {
	db *sql.DB
}

func NewPRPostgres(db *sql.DB) repository.PRRepository {
	return &PRPostgres{db: db}
}

func (r *PRPostgres) Create(ctx context.Context, pr *domain.PullRequest) error {
	log := logger.L()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction", slog.Any("err", err))
		return err
	}

	q := `
        INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err = tx.ExecContext(ctx, q,
		pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt,
	)
	if err != nil {
		tx.Rollback()
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return err
	}

	q = `
        INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
        VALUES ($1, $2)
	`
	for _, reviewer := range pr.AssignedReviewers {
		_, err = tx.ExecContext(ctx, q, pr.ID, reviewer)
		if err != nil {
			tx.Rollback()
			log.Error("failed to execute SQL",
				slog.String("query", q),
				slog.Any("err", err),
			)
			return err
		}
	}

	return tx.Commit()
}

func (r *PRPostgres) Exists(ctx context.Context, id string) (bool, error) {
	log := logger.L()

	q := `
        SELECT 1 FROM pull_requests WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, q, id)

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

func (r *PRPostgres) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	log := logger.L()

	q := `
        SELECT id, name, author_id, status, created_at, merged_at
        FROM pull_requests
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, q, id)

	var pr domain.PullRequest
	if err := row.Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPRNotFound
		}
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}

	q = `
        SELECT reviewer_id FROM pull_request_reviewers
        WHERE pr_id = $1
    `
	rows, err := r.db.QueryContext(ctx, q, id)
	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, rid)
	}

	pr.AssignedReviewers = reviewers
	return &pr, nil
}

func (r *PRPostgres) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	log := logger.L()

	q := `
        SELECT pr.id, pr.name, pr.author_id, pr.status, pr.created_at, pr.merged_at
        FROM pull_requests pr
        JOIN pull_request_reviewers r ON pr.id = r.pr_id
        WHERE r.reviewer_id = $1
    `
	rows, err := r.db.QueryContext(ctx, q, reviewerID)
	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}
	defer rows.Close()

	var list []domain.PullRequest

	for rows.Next() {
		var pr domain.PullRequest
		if err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.MergedAt,
		); err != nil {
			return nil, err
		}

		revs, err := r.fetchReviewers(ctx, pr.ID)
		if err != nil {
			return nil, err
		}
		pr.AssignedReviewers = revs

		list = append(list, pr)
	}

	return list, nil
}

func (r *PRPostgres) UpdateReviewers(ctx context.Context, prID string, reviewers []string) error {
	log := logger.L()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := `
        DELETE FROM pull_request_reviewers WHERE pr_id = $1
    `
	_, err = tx.ExecContext(ctx, q, prID)
	if err != nil {
		tx.Rollback()
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return err
	}

	for _, rid := range reviewers {
		q = `
            INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
            VALUES ($1, $2)
        `
		_, err = tx.ExecContext(ctx, q, prID, rid)
		if err != nil {
			tx.Rollback()
			log.Error("failed to execute SQL",
				slog.String("query", q),
				slog.Any("err", err),
			)
			return err
		}
	}

	return tx.Commit()
}

func (r *PRPostgres) UpdateStatusAndMergedAt(
	ctx context.Context,
	id string,
	status domain.PRStatus,
	mergedAt *time.Time,
) error {
	log := logger.L()

	q := `
        UPDATE pull_requests
        SET status = $2, merged_at = $3
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, q,
		id, status, mergedAt,
	)
	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
	}
	return err
}

func (r *PRPostgres) fetchReviewers(ctx context.Context, prID string) ([]string, error) {
	log := logger.L()

	q := `
        SELECT reviewer_id FROM pull_request_reviewers
        WHERE pr_id = $1
    `
	rows, err := r.db.QueryContext(ctx, q, prID)
	if err != nil {
		log.Error("failed to execute SQL",
			slog.String("query", q),
			slog.Any("err", err),
		)
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, id)
	}

	return reviewers, nil
}
