package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// Secure middleware configures echo with secure headers
func Secure() echo.MiddlewareFunc {
	return echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "",
	})
}
