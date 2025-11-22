package service

import (
	"context"
	"log/slog"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
)

type StatsService struct {
	statsRepo repository.StatsRepository
}

func NewStatsService(statsRepo repository.StatsRepository) *StatsService {
	return &StatsService{statsRepo: statsRepo}
}

func (s *StatsService) GetStats(ctx context.Context) ([]domain.UserAssignmentStat, []domain.PRReviewerStat, error) {
	log := logger.L()
	log.Info("collecting statistics")

	byUser, err := s.statsRepo.CountAssignmentsByUser(ctx)
	if err != nil {
		log.Error("failed get count assignments by user",
			slog.Any("err", err),
		)
		return nil, nil, err
	}

	byPR, err := s.statsRepo.CountReviewersByPR(ctx)
	if err != nil {
		log.Error("failed get count reviewers by pull request",
			slog.Any("err", err),
		)
		return nil, nil, err
	}

	log.Info("statistics successfully collected")

	return byUser, byPR, nil
}
