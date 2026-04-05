package route

import (
	"crist-blog/internal/handler"
	"crist-blog/internal/middleware"
	"crist-blog/internal/service"

	"github.com/labstack/echo/v4"
)

func SetupUploadRouter(e *echo.Echo, uploadHandler *handler.UploadHandler, uploadCOSHandler *handler.COSHandler, authService *service.AuthService) {
	api := e.Group("/api")

	upload := api.Group("/upload")
	upload.Use(middleware.AuthMiddleware(authService))
	upload.POST("/image", uploadCOSHandler.UploadImage)
	upload.POST("/images", uploadCOSHandler.UploadImageBatch)
	upload.POST("/init", uploadHandler.InitChunkedUpload)
	upload.POST("/chunk", uploadHandler.UploadChunk)
	upload.POST("/merge", uploadHandler.MergeChunks)
}
