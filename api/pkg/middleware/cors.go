package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/nvidia/carbide-rest/api/pkg/api/pagination"
)

// CORS configures echo to accept requests from scripts from domains different
// than the API will be served from
func CORS() echo.MiddlewareFunc {
	return echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		MaxAge:           86400,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", pagination.ResponseHeaderName},
		AllowCredentials: true,
	})
}
