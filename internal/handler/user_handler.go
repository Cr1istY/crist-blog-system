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
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	userAgent := c.Request().UserAgent()
	ip := c.RealIP()

	accessToken, refreshToken, err := h.authService.GenerateTokensWithAgent(user, userAgent, ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	isProduction := c.Request().TLS != nil || c.Request().Header.Get("X-Forwarded-Proto") == "https"

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   isProduction,         // 仅生产环境启用
		SameSite: http.SameSiteLaxMode, // 开发环境用 Lax
		Path:     "/",
		MaxAge:   int(h.authService.GetTheRefreshTokenExpired().Seconds()),
	}

	// 生产环境才设置 Domain
	if isProduction {
		cookie.Domain = "foreveryang.cn"
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func (h *UserHandler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	userAgent := c.Request().UserAgent()
	ip := c.RealIP()
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Refresh token not found"})
	}

	accessToken, err := h.authService.RefreshAccessTokenWithIpAndAgent(cookie.Value, userAgent, ip)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}
