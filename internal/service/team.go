package service

import (
	"context"
	"fmt"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
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
	if teamName == "" {
		return nil, fmt.Errorf("team name is required")
	}

	exists, err := s.teamRepo.ExistsByName(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTeamExists
	}

	team := &domain.Team{
		Name: teamName,
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	if name == "" {
		return nil, fmt.Errorf("empty team name")
	}

	team, err := s.teamRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return team, nil
}
