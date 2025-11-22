package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
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
	log := logger.L()

	log.Info("upserting user",
		slog.String("userID", userID),
		slog.String("username", username),
		slog.String("teamName", teamName),
		slog.Bool("isActive", isActive),
	)

	if userID == "" || username == "" || teamName == "" {
		log.Warn("invalid input: missing required fields",
			slog.String("userID", userID),
			slog.String("username", username),
			slog.String("teamName", teamName),
		)
		return nil, fmt.Errorf("invalid input: missing required fields")
	}

	exists, err := s.teamRepo.ExistsByName(ctx, teamName)
	if err != nil {
		log.Error("failed to check if team exists",
			slog.String("teamName", teamName),
			slog.Any("err", err),
		)
		return nil, err
	}
	if !exists {
		log.Warn("team does not exist",
			slog.String("teamName", teamName),
		)
		return nil, domain.ErrTeamNotFound
	}

	user := &domain.User{
		ID:       userID,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	}

	if err := s.userRepo.Upsert(ctx, user); err != nil {
		log.Error("failed to upsert user",
			slog.String("userID", userID),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("user successfully upserted",
		slog.String("userID", userID),
	)

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	log := logger.L()

	log.Info("fetching user by id", slog.String("userID", userID))

	if userID == "" {
		log.Warn("empty user id provided")
		return nil, fmt.Errorf("empty user id")
	}

	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Warn("user not found", slog.String("userID", userID))
			return nil, err
		}
		log.Error("failed to fetch user",
			slog.String("userID", userID),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("user successfully fetched",
		slog.String("userID", userID),
	)

	return u, nil
}

func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	log := logger.L()

	log.Info("setting user active state",
		slog.String("userID", userID),
		slog.Bool("isActive", isActive),
	)

	if userID == "" {
		log.Warn("empty user id provided")
		return nil, fmt.Errorf("empty user id")
	}

	u, err := s.userRepo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Warn("user not found", slog.String("userID", userID))
			return nil, err
		}
		log.Error("failed to update user active status",
			slog.String("userID", userID),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("user active state updated",
		slog.String("userID", userID),
		slog.Bool("isActive", isActive),
	)

	return u, nil
}

func (s *UserService) ListUsersByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	log := logger.L()

	log.Info("listing users by team", slog.String("teamName", teamName))

	if teamName == "" {
		log.Warn("empty team name provided")
		return nil, fmt.Errorf("empty team name")
	}

	users, err := s.userRepo.ListByTeam(ctx, teamName)
	if err != nil {
		log.Error("failed to list users by team",
			slog.String("teamName", teamName),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("users successfully listed",
		slog.String("teamName", teamName),
		slog.Int("count", len(users)),
	)

	return users, nil
}
