package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository/mocks"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/service"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.Setup("test")
}

func TestUserService_UpsertUser_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := mocks.NewUserRepository(t)
	teamRepo := mocks.NewTeamRepository(t)

	svc := service.NewUserService(userRepo, teamRepo)

	teamRepo.
		On("ExistsByName", ctx, "backend").
		Return(true, nil)

	userRepo.
		On("Upsert", ctx, mock.AnythingOfType("*domain.User")).
		Return(nil)

	user, err := svc.UpsertUser(ctx, "u1", "Alice", "backend", true)

	require.NoError(t, err)
	require.Equal(t, "u1", user.ID)
	require.Equal(t, "Alice", user.Username)
	require.Equal(t, "backend", user.TeamName)
	require.True(t, user.IsActive)

	teamRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestUserService_UpsertUser_TeamNotFound(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	teamRepo := mocks.NewTeamRepository(t)

	svc := service.NewUserService(userRepo, teamRepo)

	teamRepo.
		On("ExistsByName", mock.Anything, "mobile").
		Return(false, nil).
		Once()

	user, err := svc.UpsertUser(context.Background(), "u2", "Bob", "mobile", true)

	require.Error(t, err)
	require.Nil(t, user)
	require.Equal(t, domain.ErrTeamNotFound, err)

	teamRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestUserService_UpsertUser_ExistsErr(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	teamRepo := mocks.NewTeamRepository(t)

	svc := service.NewUserService(userRepo, teamRepo)

	expectedErr := errors.New("db failure")

	teamRepo.
		On("ExistsByName", mock.Anything, "backend").
		Return(false, expectedErr).
		Once()

	user, err := svc.UpsertUser(context.Background(), "u1", "Alice", "backend", true)

	require.Error(t, err)
	require.Nil(t, user)
	require.Equal(t, expectedErr, err)

	teamRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}
