package route

import (
	"bytes"
	"crist-blog/internal/handler"
	"crist-blog/internal/middleware"
	"crist-blog/internal/service"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func SetupBlogRouter(e *echo.Echo, postHandler *handler.PostHandler, authService *service.AuthService) {
	api := e.Group("/api")
	api.GET("/proxy/image", proxyImage)

	// 公开路由
	posts := api.Group("/posts")
	posts.GET("/getAllPosts", postHandler.ListToFrontend)
	posts.GET("/get/:id", postHandler.GetBlogToViewers)
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
}

// proxyImage 处理图片代理请求
func proxyImage(c echo.Context) error {
	// 获取目标图片URL
	imageURL := c.QueryParam("url")
	if imageURL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing url parameter")
	}

	// 验证URL格式
	_, err := url.ParseRequestURI(imageURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL format")
	}

	// 创建HTTP请求
	resp, err := http.Get(imageURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Failed to fetch image: %v", err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(resp.Body)

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Image server returned status: %d", resp.StatusCode))
	}

	// 读取图片数据
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Failed to read image data: %v", err))
	}

	// 设置响应头
	c.Response().Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	c.Response().Header().Set("Cache-Control", "public, max-age=3600") // 缓存1小时

	// 返回图片数据
	return c.Stream(http.StatusOK, resp.Header.Get("Content-Type"),
		bytes.NewReader(buf.Bytes()))
}
