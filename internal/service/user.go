package service

import (
	"context"
	"fmt"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *UserService) UpsertUser(
	ctx context.Context,
	userID string,
	username string,
	teamName string,
	isActive bool,
) (*domain.User, error) {
	if userID == "" || username == "" || teamName == "" {
		return nil, fmt.Errorf("invalid input: missing required fields")
	}

	exists, err := s.teamRepo.ExistsByName(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrTeamNotFound
	}

	user := &domain.User{
		ID:       userID,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	}

	if err := s.userRepo.Upsert(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("empty user id")
	}

	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("empty user id")
	}

	u, err := s.userRepo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) ListUsersByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	if teamName == "" {
		return nil, fmt.Errorf("empty team name")
	}

	users, err := s.userRepo.ListByTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	return users, nil
}
