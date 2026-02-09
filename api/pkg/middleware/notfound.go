package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nvidia/carbide-rest/api/internal/config"
	ccu "github.com/nvidia/carbide-rest/common/pkg/util"
)

// NotFoundHandler returns a middleware that returns a 404 status code for unmatched routes
func NotFoundHandler(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip auth processing for unmatched path
			if c.Path() == fmt.Sprintf("/%s/*", cfg.GetAPIRouteVersion()) {
				return ccu.NewAPIErrorResponse(c, http.StatusNotFound, "The requested path could not be found", nil)
			}

			return next(c)
		}
	}
}
