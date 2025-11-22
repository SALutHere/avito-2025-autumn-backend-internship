package routers

import (
	"context"
	"net/http"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/config"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/dto"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/service"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService *service.UserService
	prService   *service.PRService
}

func NewUserController(
	userService *service.UserService,
	prService *service.PRService,
) *UserController {
	return &UserController{
		userService: userService,
		prService:   prService,
	}
}

func RegisterUserRoutes(e *echo.Echo, h *UserController) {
	e.POST("/users/setIsActive", h.SetIsActive)
	e.GET("/users/getReview", h.GetReview)
}

func (h *UserController) SetIsActive(c echo.Context) error {
	var req dto.SetIsActiveUserRequest
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

	user, err := h.userService.SetUserActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		return writeDomainError(c, err)
	}

	resp := dto.SetIsActiveUserResponse{
		User: dto.ToUserDTO(user),
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *UserController) GetReview(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorObject{
				Code:    dto.ErrorCodeNotFound,
				Message: "user_id is required",
			},
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), config.C().PGTimeout)
	defer cancel()

	prs, err := h.prService.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return writeDomainError(c, err)
	}

	shorts := make([]dto.PullRequestShortDTO, 0, len(prs))
	for _, pr := range prs {
		shorts = append(shorts, dto.ToPullRequestShortDTO(pr))
	}

	resp := dto.GetReviewUserResponse{
		UserID:       userID,
		PullRequests: shorts,
	}

	return c.JSON(http.StatusOK, resp)
}
