package routers

import (
	"context"
	"net/http"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/config"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/dto"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/service"
	"github.com/labstack/echo/v4"
)

type PRController struct {
	prService *service.PRService
}

func NewPRController(prService *service.PRService) *PRController {
	return &PRController{
		prService: prService,
	}
}

func RegisterPRRoutes(e *echo.Echo, h *PRController) {
	e.POST("/pullRequest/create", h.Create)
	e.POST("/pullRequest/merge", h.Merge)
	e.POST("/pullRequest/reassign", h.Reassign)
}

func (h *PRController) Create(c echo.Context) error {
	var req dto.CreatePRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorObject{
				Code:    dto.ErrorCodeNotFound,
				Message: "invalid request body",
			},
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), config.C().PGTimeout)
	defer cancel()

	pr, err := h.prService.CreatePR(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		return writeDomainError(c, err)
	}

	resp := dto.CreatePRResponse{
		PR: dto.ToPullRequestDTO(pr),
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *PRController) Merge(c echo.Context) error {
	var req dto.MergePRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorObject{
				Code:    dto.ErrorCodeNotFound,
				Message: "invalid request body",
			},
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), config.C().PGTimeout)
	defer cancel()

	pr, err := h.prService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		return writeDomainError(c, err)
	}

	resp := dto.MergePRResponse{
		PR: dto.ToPullRequestDTO(pr),
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *PRController) Reassign(c echo.Context) error {
	var req dto.ReassignReviewerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorObject{
				Code:    dto.ErrorCodeNotFound,
				Message: "invalid request body",
			},
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), config.C().PGTimeout)
	defer cancel()

	pr, newID, err := h.prService.ReassignReviewer(ctx, req.PullRequestID, req.OldReviewerID)
	if err != nil {
		return writeDomainError(c, err)
	}

	resp := dto.ReassignReviewerResponse{
		PR:         dto.ToPullRequestDTO(pr),
		ReplacedBy: newID,
	}

	return c.JSON(http.StatusOK, resp)
}
