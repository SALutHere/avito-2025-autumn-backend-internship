package server

import (
	"time"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/internal/controller/middleware"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

func NewHTTPServer(
	readTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,

) *echo.Echo {
	e := echo.New()

	e.Server.ReadTimeout = readTimeout
	e.Server.WriteTimeout = writeTimeout
	e.Server.IdleTimeout = idleTimeout

	e.Use(middleware.HTTPLogger())
	e.Use(mw.Recover())

	// TODO: register routes

	return e
}
