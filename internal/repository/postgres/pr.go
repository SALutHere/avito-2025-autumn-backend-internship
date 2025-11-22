package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
)

type PRPostgres struct {
	db *sql.DB
}

func NewPRPostgres(db *sql.DB) repository.PRRepository {
	return &PRPostgres{db: db}
}

func (r *PRPostgres) Create(ctx context.Context, pr *domain.PullRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
        INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `,
		pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, reviewer := range pr.AssignedReviewers {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
            VALUES ($1, $2)
        `, pr.ID, reviewer)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *PRPostgres) Exists(ctx context.Context, id string) (bool, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT 1 FROM pull_requests WHERE id = $1
    `, id)

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

func (r *PRPostgres) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, author_id, status, created_at, merged_at
        FROM pull_requests
        WHERE id = $1
    `, id)

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
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
        SELECT reviewer_id FROM pull_request_reviewers
        WHERE pr_id = $1
    `, id)
	if err != nil {
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
	rows, err := r.db.QueryContext(ctx, `
        SELECT pr.id, pr.name, pr.author_id, pr.status, pr.created_at, pr.merged_at
        FROM pull_requests pr
        JOIN pull_request_reviewers r ON pr.id = r.pr_id
        WHERE r.reviewer_id = $1
    `, reviewerID)
	if err != nil {
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
        DELETE FROM pull_request_reviewers WHERE pr_id = $1
    `, prID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, rid := range reviewers {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO pull_request_reviewers (pr_id, reviewer_id)
            VALUES ($1, $2)
        `, prID, rid)
		if err != nil {
			tx.Rollback()
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
	_, err := r.db.ExecContext(ctx, `
        UPDATE pull_requests
        SET status = $2, merged_at = $3
        WHERE id = $1
    `,
		id, status, mergedAt,
	)
	return err
}

func (r *PRPostgres) fetchReviewers(ctx context.Context, prID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT reviewer_id FROM pull_request_reviewers
        WHERE pr_id = $1
    `, prID)
	if err != nil {
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
