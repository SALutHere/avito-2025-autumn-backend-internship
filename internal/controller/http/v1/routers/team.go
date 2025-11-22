package routers

import (
	"context"
	"net/http"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/config"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/dto"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/service"
	"github.com/labstack/echo/v4"
)

type TeamController struct {
	teamService *service.TeamService
	userService *service.UserService
}

func NewTeamController(
	teamService *service.TeamService,
	userService *service.UserService,
) *TeamController {
	return &TeamController{
		teamService: teamService,
		userService: userService,
	}
}

func RegisterTeamRoutes(e *echo.Echo, h *TeamController) {
	e.POST("/team/add", h.AddTeam)
	e.GET("/team/get", h.GetTeam)
}

func (h *TeamController) AddTeam(c echo.Context) error {
	var req dto.TeamDTO
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

	team, err := h.teamService.CreateTeam(ctx, req.TeamName)
	if err != nil {
		return writeDomainError(c, err)
	}

	for _, m := range req.Members {
		_, err := h.userService.UpsertUser(ctx, m.UserID, m.Username, team.Name, m.IsActive)
		if err != nil {
			return writeDomainError(c, err)
		}
	}

	resp := dto.AddTeamResponse{
		Team: dto.TeamDTO{
			TeamName: team.Name,
			Members:  req.Members,
		},
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *TeamController) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorObject{
				Code:    dto.ErrorCodeNotFound,
				Message: "team_name is required",
			},
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), config.C().PGTimeout)
	defer cancel()

	team, err := h.teamService.GetTeam(ctx, teamName)
	if err != nil {
		return writeDomainError(c, err)
	}

	users, err := h.userService.ListUsersByTeam(ctx, team.Name)
	if err != nil {
		return writeDomainError(c, err)
	}

	members := make([]dto.TeamMemberDTO, 0, len(users))
	for _, u := range users {
		members = append(members, dto.TeamMemberDTO{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	resp := dto.TeamDTO{
		TeamName: team.Name,
		Members:  members,
	}

	return c.JSON(http.StatusOK, resp)
}
