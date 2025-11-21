package logger

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func HTTPLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			err := next(c)

			log.Info("request completed",
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.Int("status", res.Status),
				slog.Int64("bytes_out", res.Size),
				slog.String("remote_addr", c.RealIP()),
				slog.String("user_agent", req.UserAgent()),
				slog.Duration("duration", time.Since(start)),
			)

			return err
		}
	}
}
