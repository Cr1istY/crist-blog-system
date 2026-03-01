package main

import (
	"crist-blog/internal/blogConfig"
	"crist-blog/internal/handler"
	"crist-blog/internal/repository"
	"crist-blog/internal/route"
	"crist-blog/internal/service"
	config "crist-blog/internal/uploadConfig"
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 使用示例
func main() {

	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Failed to generate JWT secret")
	}
	jwtSecret := base64.URLEncoding.EncodeToString(bytes)

	db := blogConfig.ConnectDB()
	redis := blogConfig.ConnectRedis()
	searcher := blogConfig.ConnectIp2Region()

	defer func() {
		if redis != nil {
			_ = redis.Close()
		}
	}()

	// 初始化配置, 导入.env文件
	uploadCfg := config.Load()

	userRepo := repository.NewUserRepository(db, redis)
	authRepo := repository.NewRefreshTokenRepository(db)
	postRepo := repository.NewPostRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	imageRepo := repository.NewImageRepository(db)
	tweetRepo := repository.NewTweetRepository(db, redis)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, authRepo, searcher, jwtSecret)
	postService := service.NewPostService(postRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	imageService := service.NewImageService(imageRepo)
	tweetService := service.NewTweetService(tweetRepo)
	cos := config.NewCOSService()

	uploadHandler := handler.NewUploadHandler(uploadCfg)
	uploadCOSHandler := handler.NewCOSHandler(uploadHandler, imageService, cos, uploadCfg)
	postHandler := handler.NewPostHandler(postService, categoryService, redis)
	userHandler := handler.NewUserHandler(authService, userService)
	imageHandler := handler.NewImageHandler(redis)
	tweetHandler := handler.NewTweetHandler(tweetService, userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	e := echo.New()
	e.Use(middleware.BodyLimit("10M"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
	e.Use(middleware.Recover())

	route.SetupUserRoutes(e, userHandler, authService)
	route.SetupBlogRouter(e, postHandler, imageHandler, authService)
	route.SetupCategoryRouter(e, categoryHandler, authService)
	route.SetupUploadRouter(e, uploadHandler, uploadCOSHandler, authService)
	route.SetupTweetRouter(e, tweetHandler, authService)
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server is running on port", port)
	e.Logger.Fatal(e.Start("127.0.0.1:" + port))

}
