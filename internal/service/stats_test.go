package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository/mocks"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/service"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
)

func init() {
	logger.Setup("local")
}

func TestStatsService_GetStats_Success(t *testing.T) {
	ctx := context.Background()

	statsRepo := mocks.NewStatsRepository(t)
	svc := service.NewStatsService(statsRepo)

	expectedByUser := []domain.UserAssignmentStat{
		{UserID: "u1", Assignments: 5},
		{UserID: "u2", Assignments: 3},
	}
	expectedByPR := []domain.PRReviewerStat{
		{PRID: "pr-101", Reviewers: 2},
		{PRID: "pr-102", Reviewers: 1},
	}

	statsRepo.
		On("CountAssignmentsByUser", ctx).
		Return(expectedByUser, nil).
		Once()

	statsRepo.
		On("CountReviewersByPR", ctx).
		Return(expectedByPR, nil).
		Once()

	byUser, byPR, err := svc.GetStats(ctx)

	require.NoError(t, err)
	require.Equal(t, expectedByUser, byUser)
	require.Equal(t, expectedByPR, byPR)
}

func TestStatsService_GetStats_CountAssignmentsError(t *testing.T) {
	ctx := context.Background()

	statsRepo := mocks.NewStatsRepository(t)
	svc := service.NewStatsService(statsRepo)

	expectedErr := assert.AnError

	statsRepo.
		On("CountAssignmentsByUser", ctx).
		Return(nil, expectedErr).
		Once()

	byUser, byPR, err := svc.GetStats(ctx)

	require.Error(t, err)
	require.Same(t, expectedErr, err)
	require.Nil(t, byUser)
	require.Nil(t, byPR)
}

func TestStatsService_GetStats_CountReviewersError(t *testing.T) {
	ctx := context.Background()

	statsRepo := mocks.NewStatsRepository(t)
	svc := service.NewStatsService(statsRepo)

	statsRepo.
		On("CountAssignmentsByUser", ctx).
		Return([]domain.UserAssignmentStat{
			{UserID: "u1", Assignments: 5},
		}, nil).
		Once()

	expectedErr := assert.AnError

	statsRepo.
		On("CountReviewersByPR", ctx).
		Return(nil, expectedErr).
		Once()

	byUser, byPR, err := svc.GetStats(ctx)

	require.Error(t, err)
	require.Nil(t, byUser)
	require.Nil(t, byPR)
	require.Same(t, expectedErr, err)
}
