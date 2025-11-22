package v1

import (
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/config"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/controller/http/v1/middleware"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/controller/http/v1/routers"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

func NewHTTPServer(
	teamCtrl *routers.TeamController,
	userCtrl *routers.UserController,
	prCtrl *routers.PRController,
	statsCtrl *routers.StatsController,
) *echo.Echo {
	cfg := config.C()
	e := echo.New()

	e.Server.ReadTimeout = cfg.HTTPReadTimeout
	e.Server.WriteTimeout = cfg.HTTPWriteTimeout
	e.Server.IdleTimeout = cfg.HTTPIdleTimeout

	e.Use(middleware.HTTPLogger())
	e.Use(mw.Recover())

	routers.RegisterTeamRoutes(e, teamCtrl)
	routers.RegisterUserRoutes(e, userCtrl)
	routers.RegisterPRRoutes(e, prCtrl)
	routers.RegisterStatsRoutes(e, statsCtrl)

	return e
}
