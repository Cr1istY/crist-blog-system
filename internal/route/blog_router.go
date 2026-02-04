package route

import (
	"crist-blog/internal/handler"
	"crist-blog/internal/middleware"
	"crist-blog/internal/service"

	"github.com/labstack/echo/v4"
)

func SetupBlogRouter(e *echo.Echo, postHandler *handler.PostHandler, imageHandler *handler.ImageHandler, authService *service.AuthService) {
	api := e.Group("/api")
	api.GET("/proxy/image", imageHandler.ProxyImage)

	// 公开路由
	posts := api.Group("/posts")
	posts.GET("/getAllPosts", postHandler.ListToFrontendWithPinned)
	posts.GET("/get/:id", postHandler.GetBlogByIdToViewers)
	posts.GET("/getBySlug/:slug", postHandler.GetBlogBySlug)
	posts.GET("/hot", postHandler.GetHotPosts)
	posts.GET("/latest", postHandler.GetLatestPosts)
	posts.GET("/addViews/:id", postHandler.AddViews)
	posts.GET("/addLikes/:id", postHandler.AddLikes)

	// 需要认证的路由
	protected := api.Group("/posts")
	protected.Use(middleware.AuthMiddleware(authService))
	protected.POST("/create", postHandler.CreatePost)
	protected.PUT("/update/:id", postHandler.Update)
	protected.DELETE("/delete/:id", postHandler.Delete)
	protected.PUT("/pin/:id", postHandler.PinPost)
	protected.PUT("/unpin/:id", postHandler.UnpinPost)
}
