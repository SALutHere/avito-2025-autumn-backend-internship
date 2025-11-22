package dto

type UserStatsDTO struct {
	UserID      string `json:"user_id"`
	Assignments int    `json:"assignments"`
}

type PRStatsDTO struct {
	PullRequestID string `json:"pull_request_id"`
	Reviewers     int    `json:"reviewers"`
}

type StatsResponse struct {
	ByUser []UserStatsDTO `json:"by_user"`
	ByPR   []PRStatsDTO   `json:"by_pr"`
}
