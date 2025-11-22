package repository

import (
	"context"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
)

type StatsRepository interface {
	CountAssignmentsByUser(ctx context.Context) ([]domain.UserAssignmentStat, error)
	CountReviewersByPR(ctx context.Context) ([]domain.PRReviewerStat, error)
}
