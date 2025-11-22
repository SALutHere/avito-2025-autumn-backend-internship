package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTeamNotFound = errors.New("team not found")
	ErrPRNotFound   = errors.New("pull request not found")

	ErrTeamExists = errors.New("team already exists")
	ErrPRExists   = errors.New("pull request already exists")

	ErrPRAlreadyMerged = errors.New("pull request is already merged")

	ErrNotAssigned = errors.New("user is not assigned as reviewer")

	ErrNoCandidate = errors.New("no candidate reviewer available")
)
