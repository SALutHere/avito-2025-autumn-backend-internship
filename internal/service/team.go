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

type TeamService struct {
	teamRepo repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	log := logger.L()

	log.Info("creating team", slog.String("teamName", teamName))

	if teamName == "" {
		log.Warn("empty team name provided")
		return nil, fmt.Errorf("team name is required")
	}

	exists, err := s.teamRepo.ExistsByName(ctx, teamName)
	if err != nil {
		log.Error("failed to check if team exists",
			slog.String("teamName", teamName),
			slog.Any("err", err),
		)
		return nil, err
	}
	if exists {
		log.Warn("team already exists", slog.String("teamName", teamName))
		return nil, domain.ErrTeamExists
	}

	team := &domain.Team{
		Name: teamName,
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		log.Error("failed to create team",
			slog.String("teamName", teamName),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("team successfully created", slog.String("teamName", teamName))

	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	log := logger.L()

	log.Info("fetching team", slog.String("teamName", name))

	if name == "" {
		log.Warn("empty team name provided")
		return nil, fmt.Errorf("empty team name")
	}

	team, err := s.teamRepo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			log.Warn("team not found", slog.String("teamName", name))
			return nil, err
		}
		log.Error("failed to fetch team",
			slog.String("teamName", name),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("team successfully fetched", slog.String("teamName", name))

	return team, nil
}
