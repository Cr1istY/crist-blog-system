package middleware

import (
	"crist-blog/internal/service"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr, err := GetTokenStrInAuthHeader(c)
			if err != nil {
				if errors.Is(err, AuthHeaderEmptyErr) {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authHeader == \"\""})
				}
				if errors.Is(err, AuthUnauthorizedErr) {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unknown error"})
			}
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(authService.JwtSecret()), nil
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
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}
			c.Set("user_id", userID)
			c.Set("user_id_str", userIDStr)
			return next(c)
		}
	}
}

var (
	AuthHeaderEmptyErr  = errors.New("auth header is empty")
	AuthUnauthorizedErr = errors.New("unauthorized")
)

func GetTokenStrInAuthHeader(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", AuthHeaderEmptyErr
	}
	pars := strings.Split(authHeader, " ")
	if len(pars) != 2 || pars[0] != "Bearer" {
		return "", AuthUnauthorizedErr
	}
	tokenStr := pars[1]
	return tokenStr, nil
}
