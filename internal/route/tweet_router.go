package route

import (
	"crist-blog/internal/handler"
	"crist-blog/internal/middleware"
	"crist-blog/internal/service"

	"github.com/labstack/echo/v4"
)

func SetupTweetRouter(e *echo.Echo, tweetHandler *handler.TweetHandler, authService *service.AuthService) {
	api := e.Group("/api")

	// tweetPublic := api.Group("/tweet")

	tweetAuth := api.Group("/tweet", middleware.AuthMiddleware(authService))
	tweetAuth.POST("/create", tweetHandler.CreateTweet)

}
