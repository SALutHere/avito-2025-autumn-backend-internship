package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
)

type PRService struct {
	prRepo   repository.PRRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPRService(
	prRepo repository.PRRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *PRService) CreatePR(ctx context.Context, prID string, prReqName string, authorID string) (*domain.PullRequest, error) {
	if prID == "" || prReqName == "" || authorID == "" {
		return nil, fmt.Errorf("invalid input: empty fields")
	}

	exists, err := s.prRepo.Exists(ctx, prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrPRExists
	}

	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		return nil, err
	}

	team, err := s.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			return nil, err
		}
		return nil, err
	}
	if team == nil {
		return nil, domain.ErrTeamNotFound
	}

	candidates, err := s.userRepo.ListActiveByTeam(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	filtered := make([]domain.User, 0, len(candidates))
	for _, u := range candidates {
		if u.ID == author.ID {
			continue
		}
		filtered = append(filtered, u)
	}

	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	var reviewers []string
	for i := 0; i < len(filtered) && i < 2; i++ {
		reviewers = append(reviewers, filtered[i].ID)
	}

	now := time.Now().UTC()

	pr := &domain.PullRequest{
		ID:                prID,
		Name:              prReqName,
		AuthorID:          authorID,
		Status:            domain.PRStatusOpen,
		AssignedReviewers: reviewers,
		CreatedAt:         &now,
		MergedAt:          nil,
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	if prID == "" {
		return nil, fmt.Errorf("empty prID")
	}

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			return nil, err
		}
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	now := time.Now().UTC()
	if err := s.prRepo.UpdateStatusAndMergedAt(ctx, pr.ID, domain.PRStatusMerged, &now); err != nil {
		return nil, err
	}

	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	return pr, nil
}

func (s *PRService) ReassignReviewer(
	ctx context.Context,
	prID string,
	oldReviewerID string,
) (*domain.PullRequest, string, error) {
	if prID == "" || oldReviewerID == "" {
		return nil, "", fmt.Errorf("invalid input: empty fields")
	}

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			return nil, "", err
		}
		return nil, "", err
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, "", domain.ErrPRAlreadyMerged
	}

	index := -1
	for i, id := range pr.AssignedReviewers {
		if id == oldReviewerID {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, "", domain.ErrNotAssigned
	}

	oldReviewer, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, "", err
		}
		return nil, "", err
	}

	teamName := oldReviewer.TeamName

	candidates, err := s.userRepo.ListActiveByTeam(ctx, teamName)
	if err != nil {
		return nil, "", err
	}

	reviewersSet := make(map[string]struct{}, len(pr.AssignedReviewers))
	for _, id := range pr.AssignedReviewers {
		reviewersSet[id] = struct{}{}
	}

	filtered := make([]domain.User, 0, len(candidates))
	for _, u := range candidates {
		if u.ID == oldReviewerID {
			continue
		}
		if u.ID == pr.AuthorID {
			continue
		}
		if _, exists := reviewersSet[u.ID]; exists {
			continue
		}
		filtered = append(filtered, u)
	}

	if len(filtered) == 0 {
		return nil, "", domain.ErrNoCandidate
	}

	newIdx := rand.Intn(len(filtered))
	newReviewer := filtered[newIdx]

	newReviewers := make([]string, len(pr.AssignedReviewers))
	copy(newReviewers, pr.AssignedReviewers)
	newReviewers[index] = newReviewer.ID

	if err := s.prRepo.UpdateReviewers(ctx, pr.ID, newReviewers); err != nil {
		return nil, "", err
	}

	pr.AssignedReviewers = newReviewers

	return pr, newReviewer.ID, nil
}

func (s *PRService) GetPRsByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	if reviewerID == "" {
		return nil, fmt.Errorf("empty reviewerID")
	}

	prs, err := s.prRepo.ListByReviewer(ctx, reviewerID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
