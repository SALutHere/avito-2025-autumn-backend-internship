package repository

import (
	"context"
	"time"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
)

type PRRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) error

	Exists(ctx context.Context, id string) (bool, error)

	GetByID(ctx context.Context, id string) (*domain.PullRequest, error)

	ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error)

	UpdateReviewers(ctx context.Context, id string, reviewers []string) error

	UpdateStatusAndMergedAt(ctx context.Context, id string, status domain.PRStatus, mergedAt *time.Time) error
}
