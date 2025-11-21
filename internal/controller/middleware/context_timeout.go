package middleware

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

func WithContextTimeout(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			return next(c)
		}
	}
}
