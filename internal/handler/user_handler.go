package handler

import (
	"crist-blog/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewUserHandler(authService *service.AuthService, userService *service.UserService) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
	}
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}
	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	userAgent := c.Request().UserAgent()
	ip := c.RealIP()

	accessToken, refreshToken, err := h.authService.GenerateTokensWithAgent(user, userAgent, ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   int(h.authService.GetTheRefreshTokenExpired().Seconds()),
	})

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func (h *UserHandler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Refresh token not found"})
	}

	accessToken, err := h.authService.RefreshAccessToken(cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}
