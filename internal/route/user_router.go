package route

import (
	"crist-blog/internal/handler"
	"crist-blog/internal/middleware"
	"crist-blog/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(e *echo.Echo,
	userHandler *handler.UserHandler,
	authService *service.AuthService) {

	admin := e.Group("/api")
	admin.GET("/check", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Admin!")
	})
	admin.POST("/login", userHandler.Login)
	admin.POST("/auth/refresh", userHandler.Refresh)

	auth := e.Group("/api")
	auth.Use(middleware.AuthMiddleware(authService))
	{
		auth.GET("/user", func(c echo.Context) error {
			return c.JSON(http.StatusOK, c.Get("user_id"))
		})
	}
}
