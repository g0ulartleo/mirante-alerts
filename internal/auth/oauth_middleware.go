package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/labstack/echo/v4"
)

func OAuthMiddleware(authConfig *config.AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !authConfig.OAuth.Enabled {
				return APIKeyAuthMiddleware(authConfig.APIKey)(next)(c)
			}

			oauthService, err := NewOAuthService(authConfig)
			if err != nil {
				log.Printf("Error creating OAuth service: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Authentication service error")
			}

			var tokenString string

			authHeader := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
			fmt.Println("tokenString", tokenString)

			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			claims, err := oauthService.ValidateJWT(tokenString)
			if err != nil {
				log.Printf("JWT validation error: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication token")
			}

			c.Set("user", claims)
			c.Set("user_email", claims.Email)

			return next(c)
		}
	}
}

func APIKeyAuthMiddleware(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if apiKey == "" {
				return echo.NewHTTPError(http.StatusInternalServerError, "No authentication method configured")
			}

			requestAPIKey := c.Request().Header.Get("X-API-Key")
			if requestAPIKey != apiKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key")
			}

			return next(c)
		}
	}
}

func AuthenticationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authConfig, err := config.LoadAuthConfig()
			if err != nil {
				log.Printf("Error loading auth config: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Authentication configuration error")
			}

			return OAuthMiddleware(authConfig)(next)(c)
		}
	}
}
