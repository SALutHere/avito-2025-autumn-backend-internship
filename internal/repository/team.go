package repository

import (
	"context"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
)

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error

	ExistsByName(ctx context.Context, name string) (bool, error)

	GetByName(ctx context.Context, name string) (*domain.Team, error)
}
