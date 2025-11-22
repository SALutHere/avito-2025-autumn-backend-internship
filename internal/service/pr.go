package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/repository"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
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

func (s *PRService) CreatePR(
	ctx context.Context,
	prID string,
	prName string,
	authorID string,
) (*domain.PullRequest, error) {
	log := logger.L()

	log.Info("creating pull request",
		slog.String("prID", prID),
		slog.String("name", prName),
		slog.String("authorID", authorID),
	)

	if prID == "" || prName == "" || authorID == "" {
		log.Warn("invalid input: empty required fields",
			slog.String("prID", prID),
			slog.String("name", prName),
			slog.String("authorID", authorID),
		)
		return nil, fmt.Errorf("invalid input: empty fields")
	}

	exists, err := s.prRepo.Exists(ctx, prID)
	if err != nil {
		log.Error("failed to check PR existence",
			slog.String("prID", prID),
			slog.Any("err", err),
		)
		return nil, err
	}
	if exists {
		log.Warn("pull request already exists",
			slog.String("prID", prID),
		)
		return nil, domain.ErrPRExists
	}

	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Warn("author not found",
				slog.String("authorID", authorID),
			)
			return nil, err
		}
		log.Error("failed to fetch author",
			slog.String("authorID", authorID),
			slog.Any("err", err),
		)
		return nil, err
	}

	team, err := s.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) {
			log.Warn("team not found for author",
				slog.String("teamName", author.TeamName),
				slog.String("authorID", authorID),
			)
			return nil, err
		}
		log.Error("failed to fetch team",
			slog.String("teamName", author.TeamName),
			slog.Any("err", err),
		)
		return nil, err
	}
	if team == nil {
		log.Warn("team lookup returned nil",
			slog.String("teamName", author.TeamName),
		)
		return nil, domain.ErrTeamNotFound
	}

	candidates, err := s.userRepo.ListActiveByTeam(ctx, author.TeamName)
	if err != nil {
		log.Error("failed to list active reviewers",
			slog.String("teamName", author.TeamName),
			slog.Any("err", err),
		)
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
		Name:              prName,
		AuthorID:          authorID,
		Status:            domain.PRStatusOpen,
		AssignedReviewers: reviewers,
		CreatedAt:         &now,
		MergedAt:          nil,
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		log.Error("failed to create pull request",
			slog.String("prID", prID),
			slog.String("name", prName),
			slog.String("authorID", authorID),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("pull request successfully created",
		slog.String("prID", prID),
		slog.String("authorID", authorID),
		slog.Int("reviewersCount", len(reviewers)),
	)

	return pr, nil
}

func (s *PRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	log := logger.L()

	log.Info("merging pull request",
		slog.String("prID", prID),
	)

	if prID == "" {
		log.Warn("empty prID provided")
		return nil, fmt.Errorf("empty prID")
	}

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			log.Warn("pull request not found", slog.String("prID", prID))
			return nil, err
		}
		log.Error("failed to get pull request",
			slog.String("prID", prID),
			slog.Any("err", err),
		)
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		log.Info("pull request already merged", slog.String("prID", prID))
		return pr, nil
	}

	now := time.Now().UTC()
	if err := s.prRepo.UpdateStatusAndMergedAt(ctx, pr.ID, domain.PRStatusMerged, &now); err != nil {
		log.Error("failed to update pull request status",
			slog.String("prID", pr.ID),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("pull request successfully merged", slog.String("prID", prID))

	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	return pr, nil
}

func (s *PRService) ReassignReviewer(
	ctx context.Context,
	prID string,
	oldReviewerID string,
) (*domain.PullRequest, string, error) {
	log := logger.L()

	log.Info("reassigning reviewer",
		slog.String("prID", prID),
		slog.String("oldReviewerID", oldReviewerID),
	)

	if prID == "" || oldReviewerID == "" {
		log.Warn("invalid input: empty fields",
			slog.String("prID", prID),
			slog.String("oldReviewerID", oldReviewerID),
		)
		return nil, "", fmt.Errorf("invalid input: empty fields")
	}

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, domain.ErrPRNotFound) {
			log.Warn("pull request not found", slog.String("prID", prID))
			return nil, "", err
		}
		log.Error("failed to get pull request",
			slog.String("prID", prID),
			slog.Any("err", err),
		)
		return nil, "", err
	}

	if pr.Status == domain.PRStatusMerged {
		log.Warn("attempt to reassign reviewer for merged pull request",
			slog.String("prID", prID),
		)
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
		log.Warn("old reviewer is not assigned to the pull request",
			slog.String("prID", prID),
			slog.String("oldReviewerID", oldReviewerID),
		)
		return nil, "", domain.ErrNotAssigned
	}

	oldReviewer, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Warn("old reviewer not found", slog.String("oldReviewerID", oldReviewerID))
			return nil, "", err
		}
		log.Error("failed to get old reviewer",
			slog.String("oldReviewerID", oldReviewerID),
			slog.Any("err", err),
		)
		return nil, "", err
	}

	candidates, err := s.userRepo.ListActiveByTeam(ctx, oldReviewer.TeamName)
	if err != nil {
		log.Error("failed to list active candidates",
			slog.String("teamName", oldReviewer.TeamName),
			slog.Any("err", err),
		)
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
		log.Error("failed to update reviewers",
			slog.String("prID", prID),
			slog.Any("err", err),
		)
		return nil, "", err
	}

	log.Info("reviewer successfully reassigned",
		slog.String("prID", prID),
		slog.String("oldReviewerID", oldReviewerID),
		slog.String("newReviewerID", newReviewer.ID),
	)

	pr.AssignedReviewers = newReviewers

	return pr, newReviewer.ID, nil
}

func (s *PRService) GetPRsByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	log := logger.L()

	log.Info("listing pull requests by reviewer",
		slog.String("reviewerID", reviewerID),
	)

	if reviewerID == "" {
		log.Warn("empty reviewerID provided")
		return nil, fmt.Errorf("empty reviewerID")
	}

	prs, err := s.prRepo.ListByReviewer(ctx, reviewerID)
	if err != nil {
		log.Error("failed to list pull requests by reviewer",
			slog.String("reviewerID", reviewerID),
			slog.Any("err", err),
		)
		return nil, err
	}

	log.Info("successfully listed pull requests",
		slog.String("reviewerID", reviewerID),
		slog.Int("count", len(prs)),
	)

	return prs, nil
}
