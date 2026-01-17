package route

import (
	"crist-blog/internal/handler"
	"crist-blog/internal/middleware"
	"crist-blog/internal/service"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

var allowedDomains = []string{
	"www.bing.com",
	"th.bing.com",
	"gd-hbimg.huaban.com",
	"image-assets.soutushenqi.com",
	"i0.hdslb.com",
}

func isAllowedHost(host string) bool {
	for _, allowed := range allowedDomains {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}

func proxyImage(c echo.Context) error {
	imageURL := c.QueryParam("url")
	if imageURL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing url parameter")
	}

	u, err := url.ParseRequestURI(imageURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL format")
	}

	if !isAllowedHost(u.Host) {
		return echo.NewHTTPError(http.StatusForbidden, "Domain not allowed")
	}

	// 创建带超时和 UA 的请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequestWithContext(c.Request().Context(), "GET", imageURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := client.Do(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to fetch image")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return echo.NewHTTPError(http.StatusBadGateway,
			fmt.Sprintf("Upstream returned %d", resp.StatusCode))
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Only images allowed")
	}

	// 设置响应头
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")

	// 流式返回，不加载到内存
	_, err = io.Copy(c.Response(), resp.Body)
	if err != nil {
		log.Printf("Warning: failed to stream image: %v", err)
		// 无法返回错误（headers 已发送），只能记录
	}
	return nil
}
