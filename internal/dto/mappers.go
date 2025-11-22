package dto

import (
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
)

func ToUserDTO(u *domain.User) UserDTO {
	return UserDTO{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func ToPullRequestDTO(pr *domain.PullRequest) PullRequestDTO {
	return PullRequestDTO{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func ToPullRequestShortDTO(pr domain.PullRequest) PullRequestShortDTO {
	return PullRequestShortDTO{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
	}
}
