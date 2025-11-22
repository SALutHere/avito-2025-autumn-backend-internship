package service_test

import (
	"context"
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

func TestPRService_CreatePR_Success(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	userRepo := mocks.NewUserRepository(t)
	teamRepo := mocks.NewTeamRepository(t)

	svc := service.NewPRService(prRepo, userRepo, teamRepo)

	author := &domain.User{ID: "u1", TeamName: "backend"}

	prRepo.
		On("Exists", mock.Anything, "pr1").
		Return(false, nil).
		Once()

	userRepo.
		On("GetByID", mock.Anything, "u1").
		Return(author, nil).
		Once()

	teamRepo.
		On("GetByName", mock.Anything, "backend").
		Return(&domain.Team{Name: "backend"}, nil).
		Once()

	userRepo.
		On("ListActiveByTeam", mock.Anything, "backend").
		Return([]domain.User{{ID: "u2"}, {ID: "u3"}}, nil).
		Once()

	prRepo.
		On("Create", mock.Anything, mock.AnythingOfType("*domain.PullRequest")).
		Return(nil).
		Once()

	pr, err := svc.CreatePR(context.Background(), "pr1", "Fix bug", "u1")

	require.NoError(t, err)
	require.Equal(t, "pr1", pr.ID)
	require.Equal(t, "u1", pr.AuthorID)

	prRepo.AssertExpectations(t)
}

func TestPRService_CreatePR_InvalidInput(t *testing.T) {
	svc := service.NewPRService(nil, nil, nil)

	pr, err := svc.CreatePR(context.Background(), "", "name", "u1")

	require.Error(t, err)
	require.Nil(t, pr)
}

func TestPRService_CreatePR_AuthorNotFound(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	userRepo := mocks.NewUserRepository(t)

	svc := service.NewPRService(prRepo, userRepo, nil)

	prRepo.
		On("Exists", mock.Anything, "pr1").
		Return(false, nil).
		Once()

	userRepo.
		On("GetByID", mock.Anything, "u1").
		Return(nil, domain.ErrUserNotFound).
		Once()

	pr, err := svc.CreatePR(context.Background(), "pr1", "Test", "u1")

	require.Error(t, err)
	require.Equal(t, domain.ErrUserNotFound, err)
	require.Nil(t, pr)

	userRepo.AssertExpectations(t)
}

func TestPRService_CreatePR_TeamNotFound(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	userRepo := mocks.NewUserRepository(t)
	teamRepo := mocks.NewTeamRepository(t)

	svc := service.NewPRService(prRepo, userRepo, teamRepo)

	userRepo.
		On("GetByID", mock.Anything, "u1").
		Return(&domain.User{ID: "u1", TeamName: "mobile"}, nil).
		Once()

	prRepo.
		On("Exists", mock.Anything, "pr1").
		Return(false, nil).
		Once()

	teamRepo.
		On("GetByName", mock.Anything, "mobile").
		Return(nil, domain.ErrTeamNotFound).
		Once()

	pr, err := svc.CreatePR(context.Background(), "pr1", "Test", "u1")

	require.Error(t, err)
	require.Equal(t, domain.ErrTeamNotFound, err)
	require.Nil(t, pr)

	teamRepo.AssertExpectations(t)
}

func TestPRService_MergePR_Success(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	svc := service.NewPRService(prRepo, nil, nil)

	existing := &domain.PullRequest{ID: "pr1", Status: domain.PRStatusOpen}

	prRepo.
		On("GetByID", mock.Anything, "pr1").
		Return(existing, nil).
		Once()

	prRepo.
		On("UpdateStatusAndMergedAt", mock.Anything, "pr1", domain.PRStatusMerged, mock.AnythingOfType("*time.Time")).
		Return(nil).
		Once()

	pr, err := svc.MergePR(context.Background(), "pr1")

	require.NoError(t, err)
	require.Equal(t, domain.PRStatusMerged, pr.Status)

	prRepo.AssertExpectations(t)
}

func TestPRService_MergePR_NotFound(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	svc := service.NewPRService(prRepo, nil, nil)

	prRepo.
		On("GetByID", mock.Anything, "pr1").
		Return(nil, domain.ErrPRNotFound).
		Once()

	pr, err := svc.MergePR(context.Background(), "pr1")

	require.Error(t, err)
	require.Nil(t, pr)
	require.Equal(t, domain.ErrPRNotFound, err)

	prRepo.AssertExpectations(t)
}

func TestPRService_MergePR_AlreadyMerged(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	svc := service.NewPRService(prRepo, nil, nil)

	existing := &domain.PullRequest{ID: "pr1", Status: domain.PRStatusMerged}

	prRepo.
		On("GetByID", mock.Anything, "pr1").
		Return(existing, nil).
		Once()

	pr, err := svc.MergePR(context.Background(), "pr1")

	require.NoError(t, err)
	require.Equal(t, domain.PRStatusMerged, pr.Status)

	prRepo.AssertExpectations(t)
}

func TestPRService_ReassignReviewer_NotAssigned(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	svc := service.NewPRService(prRepo, nil, nil)

	existing := &domain.PullRequest{
		ID:                "pr1",
		AssignedReviewers: []string{"u2", "u3"},
	}

	prRepo.
		On("GetByID", mock.Anything, "pr1").
		Return(existing, nil).
		Once()

	pr, newID, err := svc.ReassignReviewer(context.Background(), "pr1", "u10")

	require.Error(t, err)
	require.Equal(t, domain.ErrNotAssigned, err)
	require.Nil(t, pr)
	require.Empty(t, newID)

	prRepo.AssertExpectations(t)
}

func TestPRService_ReassignReviewer_NoCandidate(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	userRepo := mocks.NewUserRepository(t)

	svc := service.NewPRService(prRepo, userRepo, nil)

	existing := &domain.PullRequest{
		ID:                "pr1",
		AuthorID:          "u1",
		AssignedReviewers: []string{"u2"},
	}

	prRepo.
		On("GetByID", mock.Anything, "pr1").
		Return(existing, nil).
		Once()

	userRepo.
		On("GetByID", mock.Anything, "u2").
		Return(&domain.User{ID: "u2", TeamName: "backend"}, nil).
		Once()

	// no candidates
	userRepo.
		On("ListActiveByTeam", mock.Anything, "backend").
		Return([]domain.User{}, nil).
		Once()

	pr, newID, err := svc.ReassignReviewer(context.Background(), "pr1", "u2")

	require.Error(t, err)
	require.Equal(t, domain.ErrNoCandidate, err)
	require.Nil(t, pr)
	require.Empty(t, newID)

	userRepo.AssertExpectations(t)
}

func TestPRService_GetPRsByReviewer_Success(t *testing.T) {
	prRepo := mocks.NewPRRepository(t)
	svc := service.NewPRService(prRepo, nil, nil)

	prRepo.
		On("ListByReviewer", mock.Anything, "u1").
		Return([]domain.PullRequest{{ID: "pr1"}}, nil).
		Once()

	prs, err := svc.GetPRsByReviewer(context.Background(), "u1")

	require.NoError(t, err)
	require.Len(t, prs, 1)

	prRepo.AssertExpectations(t)
}

func TestPRService_GetPRsByReviewer_InvalidInput(t *testing.T) {
	svc := service.NewPRService(nil, nil, nil)

	prs, err := svc.GetPRsByReviewer(context.Background(), "")

	require.Error(t, err)
	require.Nil(t, prs)
}
