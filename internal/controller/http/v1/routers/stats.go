package routers

import (
	"context"
	"net/http"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/config"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/dto"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/service"
	"github.com/labstack/echo/v4"
)

type StatsController struct {
	statsService *service.StatsService
}

func NewStatsController(statsService *service.StatsService) *StatsController {
	return &StatsController{statsService: statsService}
}

func RegisterStatsRoutes(e *echo.Echo, h *StatsController) {
	e.GET("/stats", h.Get)
}

func (h *StatsController) Get(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), config.C().PGTimeout)
	defer cancel()

	byUser, byPR, err := h.statsService.GetStats(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorObject{
				Code:    dto.ErrorCodeNotFound,
				Message: err.Error(),
			},
		})
	}

	resp := dto.StatsResponse{
		ByUser: make([]dto.UserStatsDTO, 0, len(byUser)),
		ByPR:   make([]dto.PRStatsDTO, 0, len(byPR)),
	}

	for _, s := range byUser {
		resp.ByUser = append(resp.ByUser, dto.UserStatsDTO{
			UserID:      s.UserID,
			Assignments: s.Assignments,
		})
	}

	for _, s := range byPR {
		resp.ByPR = append(resp.ByPR, dto.PRStatsDTO{
			PullRequestID: s.PRID,
			Reviewers:     s.Reviewers,
		})
	}

	return c.JSON(http.StatusOK, resp)
}
