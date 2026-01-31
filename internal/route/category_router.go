package route

import (
	"crist-blog/internal/handler"

	"github.com/labstack/echo/v4"
)

func SetupCategoryRouter(e *echo.Echo, categoryHandler *handler.CategoryHandler) {
	api := e.Group("/api")
	category := api.Group("/category")
	category.GET("/getAll", categoryHandler.ListAllCategories)

	// TODO 需要权限
	category.GET("/getFather/:son_id", categoryHandler.GetFatherCategory)
	category.POST("/create", categoryHandler.CreateCategory)
	category.POST("/createCategories", categoryHandler.CreateCategories)
	category.DELETE("/delete/:id", categoryHandler.DeleteCategory)
	category.POST("/update", categoryHandler.UpdateCategory)
	category.PUT("/addParent/:son_id/:father_id", categoryHandler.AddParentCategory)
	category.PUT("/removeParent/:son_id", categoryHandler.RemoveParentCategory)
}
