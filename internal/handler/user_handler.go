package handler

import (
	"crist-blog/internal/service"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

	// 解析获取user_id, 由于每次程序重启后，JWT会重新生成密钥，所以，需要重新登录
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authHeader == \"\""})
	}
	patrs := strings.Split(authHeader, " ")
	if len(patrs) != 2 || patrs[0] != "Bearer" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	tokenStr := patrs[1]
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.authService.JwtSecret()), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "access token expired",
			})
		}
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid access token",
		})
	}

	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid access token",
		})
	}
	calims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	userIDStr, ok := calims["user_id"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	// 现在，程序的逻辑变成了，当刷新令牌没有过期时
	// 每一次刷新操作都会重置权限令牌
	// 但是，当权限令牌过期或者被修改
	// 由于此时不能正确读取user_id
	// 所以必须重新登录

	accessToken, err := h.authService.RefreshAccessTokenWithIpAndAgent(userIDStr, cookie.Value, userAgent, ip)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}
