package repository

import (
	"context"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)

	ListActiveByTeam(ctx context.Context, teamName string) ([]domain.User, error)

	ListByTeam(ctx context.Context, teamName string) ([]domain.User, error)

	SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error)

	Upsert(ctx context.Context, user *domain.User) error
}
