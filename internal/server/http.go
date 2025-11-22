package server

import (
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/config"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/controller"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/controller/middleware"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

func NewHTTPServer(
	teamCtrl *controller.TeamController,
	userCtrl *controller.UserController,
	prCtrl *controller.PRController,
) *echo.Echo {
	cfg := config.C()
	e := echo.New()

	e.Server.ReadTimeout = cfg.HTTPReadTimeout
	e.Server.WriteTimeout = cfg.HTTPWriteTimeout
	e.Server.IdleTimeout = cfg.HTTPIdleTimeout

	e.Use(middleware.HTTPLogger())
	e.Use(mw.Recover())

	controller.RegisterTeamRoutes(e, teamCtrl)
	controller.RegisterUserRoutes(e, userCtrl)
	controller.RegisterPRRoutes(e, prCtrl)

	return e
}
