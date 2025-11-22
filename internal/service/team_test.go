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

func TestTeamService_CreateTeam_Success(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	teamRepo.
		On("ExistsByName", mock.Anything, "backend").
		Return(false, nil).
		Once()

	teamRepo.
		On("Create", mock.Anything, mock.AnythingOfType("*domain.Team")).
		Return(nil).
		Once()

	team, err := svc.CreateTeam(context.Background(), "backend")

	require.NoError(t, err)
	require.NotNil(t, team)
	require.Equal(t, "backend", team.Name)

	teamRepo.AssertExpectations(t)
}

func TestTeamService_CreateTeam_EmptyName(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	team, err := svc.CreateTeam(context.Background(), "")

	require.Error(t, err)
	require.Nil(t, team)
}

func TestTeamService_CreateTeam_ExistsErr(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	expectedErr := errors.New("db failure")

	teamRepo.
		On("ExistsByName", mock.Anything, "backend").
		Return(false, expectedErr).
		Once()

	team, err := svc.CreateTeam(context.Background(), "backend")

	require.Error(t, err)
	require.Nil(t, team)
	require.Equal(t, expectedErr, err)

	teamRepo.AssertExpectations(t)
}

func TestTeamService_CreateTeam_AlreadyExists(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	teamRepo.
		On("ExistsByName", mock.Anything, "mobile").
		Return(true, nil).
		Once()

	team, err := svc.CreateTeam(context.Background(), "mobile")

	require.Error(t, err)
	require.Nil(t, team)
	require.Equal(t, domain.ErrTeamExists, err)

	teamRepo.AssertExpectations(t)
}

func TestTeamService_CreateTeam_CreateErr(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	expectedErr := errors.New("insert failed")

	teamRepo.
		On("ExistsByName", mock.Anything, "qa").
		Return(false, nil).
		Once()

	teamRepo.
		On("Create", mock.Anything, mock.AnythingOfType("*domain.Team")).
		Return(expectedErr).
		Once()

	team, err := svc.CreateTeam(context.Background(), "qa")

	require.Error(t, err)
	require.Nil(t, team)
	require.Equal(t, expectedErr, err)

	teamRepo.AssertExpectations(t)
}

func TestTeamService_GetTeam_Success(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	expected := &domain.Team{Name: "backend"}

	teamRepo.
		On("GetByName", mock.Anything, "backend").
		Return(expected, nil).
		Once()

	team, err := svc.GetTeam(context.Background(), "backend")

	require.NoError(t, err)
	require.Equal(t, expected, team)

	teamRepo.AssertExpectations(t)
}

func TestTeamService_GetTeam_EmptyName(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	team, err := svc.GetTeam(context.Background(), "")

	require.Error(t, err)
	require.Nil(t, team)
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	teamRepo.
		On("GetByName", mock.Anything, "mobile").
		Return(nil, domain.ErrTeamNotFound).
		Once()

	team, err := svc.GetTeam(context.Background(), "mobile")

	require.Error(t, err)
	require.Nil(t, team)
	require.Equal(t, domain.ErrTeamNotFound, err)

	teamRepo.AssertExpectations(t)
}

func TestTeamService_GetTeam_RepoErr(t *testing.T) {
	teamRepo := mocks.NewTeamRepository(t)
	svc := service.NewTeamService(teamRepo)

	expectedErr := errors.New("db error")

	teamRepo.
		On("GetByName", mock.Anything, "data").
		Return(nil, expectedErr).
		Once()

	team, err := svc.GetTeam(context.Background(), "data")

	require.Error(t, err)
	require.Nil(t, team)
	require.Equal(t, expectedErr, err)

	teamRepo.AssertExpectations(t)
}
