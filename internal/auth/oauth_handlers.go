package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type OAuthHandlers struct {
	oauthService *OAuthService
	states       map[string]time.Time
}

func NewOAuthHandlers(oauthService *OAuthService) *OAuthHandlers {
	return &OAuthHandlers{
		oauthService: oauthService,
		states:       make(map[string]time.Time),
	}
}

func (h *OAuthHandlers) generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)
	h.states[state] = time.Now().Add(10 * time.Minute)
	return state, nil
}

func (h *OAuthHandlers) validateState(state string) bool {
	exp, exists := h.states[state]
	if !exists {
		return false
	}

	if time.Now().After(exp) {
		delete(h.states, state)
		return false
	}

	delete(h.states, state)
	return true
}

func (h *OAuthHandlers) cleanupExpiredStates() {
	for state, exp := range h.states {
		if time.Now().After(exp) {
			delete(h.states, state)
		}
	}
}

func (h *OAuthHandlers) LoginHandler(c echo.Context) error {
	h.cleanupExpiredStates()

	state, err := h.generateState()
	if err != nil {
		log.Printf("Error generating state: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Authentication error")
	}

	authURL := h.oauthService.GetAuthURL(state)
	fmt.Println("authURL", authURL)
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *OAuthHandlers) CallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	errorParam := c.QueryParam("error")

	if errorParam != "" {
		log.Printf("OAuth error: %s", errorParam)
		return echo.NewHTTPError(http.StatusBadRequest, "Authentication failed")
	}

	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Authorization code not provided")
	}

	if !h.validateState(state) {
		log.Printf("Invalid or expired state: %s", state)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid authentication state")
	}

	token, err := h.oauthService.ExchangeCode(c.Request().Context(), code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Token exchange failed")
	}

	userInfo, err := h.oauthService.GetUserInfo(c.Request().Context(), token)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user information")
	}

	if !h.oauthService.IsEmailAllowed(userInfo.Email) {
		log.Printf("Email not allowed: %s", userInfo.Email)
		return echo.NewHTTPError(http.StatusForbidden, "Access denied: email not authorized")
	}

	jwtToken, err := h.oauthService.GenerateJWT(userInfo)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Token generation failed")
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Request().Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60, // 24 hours
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Authentication successful",
		"user": map[string]string{
			"email": userInfo.Email,
			"name":  userInfo.Name,
		},
		"token": jwtToken,
	})
}

func (h *OAuthHandlers) LogoutHandler(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

func (h *OAuthHandlers) StatusHandler(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
	}

	tokenString := authHeader[7:]

	claims, err := h.oauthService.ValidateJWT(tokenString)
	if err != nil {
		log.Printf("JWT validation error: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication token")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user": map[string]string{
			"email": claims.Email,
			"name":  claims.Name,
		},
	})
}
