package routers

import (
	"errors"
	"net/http"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/domain"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/dto"
	"github.com/labstack/echo/v4"
)

func writeDomainError(c echo.Context, err error) error {
	var (
		status int
		code   dto.ErrorCode
	)

	switch {
	case errors.Is(err, domain.ErrTeamExists):
		status = http.StatusBadRequest
		code = dto.ErrorCodeTeamExists

	case errors.Is(err, domain.ErrPRExists):
		status = http.StatusConflict
		code = dto.ErrorCodePRExists

	case errors.Is(err, domain.ErrPRAlreadyMerged):
		status = http.StatusConflict
		code = dto.ErrorCodePRMerged

	case errors.Is(err, domain.ErrNotAssigned):
		status = http.StatusConflict
		code = dto.ErrorCodeNotAssigned

	case errors.Is(err, domain.ErrNoCandidate):
		status = http.StatusConflict
		code = dto.ErrorCodeNoCandidate

	case errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrTeamNotFound),
		errors.Is(err, domain.ErrPRNotFound):
		status = http.StatusNotFound
		code = dto.ErrorCodeNotFound

	default:
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorObject{
				Code:    dto.ErrorCodeNotFound,
				Message: err.Error(),
			},
		})
	}

	return c.JSON(status, dto.ErrorResponse{
		Error: dto.ErrorObject{
			Code:    code,
			Message: err.Error(),
		},
	})
}
