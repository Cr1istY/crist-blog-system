package route

import (
	"crist-blog/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupCategoryRouter(e *echo.Echo, categoryHandler *handler.CategoryHandler) {
	api := e.Group("/api")
	category := api.Group("/category")
	category.GET("/getAll", categoryHandler.ListAllCategories)
	category.GET("/getFather/:son_id", categoryHandler.GetFatherCategory)
	category.POST("/create", categoryHandler.CreateCategory)
}
