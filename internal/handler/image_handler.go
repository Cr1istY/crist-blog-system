package handler

import (
	"bytes"
	"context"
	"crist-blog/internal/assets"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type ImageHandler struct {
	rdb *redis.Client
}

func NewImageHandler(rdb *redis.Client) *ImageHandler {
	return &ImageHandler{rdb: rdb}
}

func generateCacheKey(imageURL string) string {
	hash := sha256.Sum256([]byte(imageURL))
	return "image_proxy:" + hex.EncodeToString(hash[:])
}

// isLikelyImage 通过文件头魔数判断是否为常见图片格式
func isLikelyImage(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	// JPEG: FF D8 FF
	if bytes.HasPrefix(data, []byte("\xFF\xD8\xFF")) {
		return true
	}
	// PNG: 89 50 4E 47
	if bytes.HasPrefix(data, []byte("\x89PNG")) {
		return true
	}
	// GIF: 47 49 46 38 (GIF8)
	if bytes.HasPrefix(data, []byte("GIF8")) {
		return true
	}
	// BMP: 42 4D (BM)
	if len(data) >= 2 && bytes.Equal(data[:2], []byte{0x42, 0x4D}) {
		return true
	}
	// WebP: RIFF .... WEBP
	if len(data) >= 12 && bytes.HasPrefix(data, []byte("RIFF")) && bytes.HasSuffix(data[:12], []byte("WEBP")) {
		return true
	}
	return false
}

func (h *ImageHandler) ProxyImage(c echo.Context) error {
	const (
		defaultCacheTTL   = 24 * time.Hour
		browserCacheShort = 3600  // 1 hour
		browserCacheLong  = 86400 // 1 day
		requestTimeout    = 10 * time.Second
	)

	imageURL := c.QueryParam("url")
	if imageURL == "" {
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusBadRequest, "Missing url parameter")
	}

	// 解析并校验 URL
	parsedURL, err := url.ParseRequestURI(imageURL)
	if err != nil || parsedURL.Scheme == "" {
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL format")
	}

	// 仅允许 http/https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusBadRequest, "Only HTTP/HTTPS URLs allowed")
	}

	// 校验 Host 是否在白名单中
	if !assets.IsAllowedHost(parsedURL.Host) {
		log.Printf("Blocked request to non-allowed host: %s", parsedURL.Host)
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusForbidden, "Domain not allowed")
	}

	// 生成缓存键
	cacheKey := generateCacheKey(imageURL)
	ctx := c.Request().Context()

	// 1. 尝试从 Redis 获取缓存
	if cacheData, err := h.rdb.Get(ctx, cacheKey).Bytes(); err == nil {
		// 即使从 Redis 读取，也做一次轻量校验（防御历史坏数据）
		if !isLikelyImage(cacheData) {
			log.Printf("Cached data for %s is not a valid image, purging", imageURL)
			h.rdb.Del(ctx, cacheKey) // 清除坏缓存
			// 继续走远程获取流程
		} else {
			contentType := http.DetectContentType(cacheData)
			c.Response().Header().Set("Content-Type", contentType)
			c.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", browserCacheLong))
			c.Response().Header().Set("X-Cache", "HIT")
			return c.Blob(http.StatusOK, contentType, cacheData)
		}
	}

	// 2. 缓存未命中或缓存无效：从远程获取
	client := &http.Client{
		Timeout: requestTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch image from %s: %v", imageURL, err)
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusBadGateway, "Failed to fetch image")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Upstream returned %d for %s", resp.StatusCode, imageURL)
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("Upstream returned %d", resp.StatusCode))
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		log.Printf("Non-image content type %s from %s", contentType, imageURL)
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Only images allowed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body from %s: %v", imageURL, err)
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read image data")
	}

	// 验证 body 是否真的是图片
	if !isLikelyImage(body) {
		log.Printf("Body does not appear to be a valid image despite content-type '%s' from %s", contentType, imageURL)
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Invalid or corrupted image data")
	}

	// 仅有效图片才异步写入 Redis
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := h.rdb.Set(cacheCtx, cacheKey, body, defaultCacheTTL).Err(); err != nil {
			log.Printf("Failed to cache image (key=%s): %v", cacheKey, err)
		}
	}()

	// 3. 返回图片
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", browserCacheShort))
	c.Response().Header().Set("X-Cache", "MISS")
	return c.Blob(http.StatusOK, contentType, body)
}
