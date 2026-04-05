package route

import (
	"crist-blog/internal/handler"
	"crist-blog/internal/middleware"
	"crist-blog/internal/service"

	"github.com/labstack/echo/v4"
)

func SetupCategoryRouter(e *echo.Echo, categoryHandler *handler.CategoryHandler, authService *service.AuthService) {
	api := e.Group("/api")
	category := api.Group("/category")
	category.GET("/getAll", categoryHandler.ListAllCategories)

	protected := api.Group("/category")
	protected.Use(middleware.AuthMiddleware(authService))
	protected.GET("/getFather/:son_id", categoryHandler.GetFatherCategory)
	protected.POST("/create", categoryHandler.CreateCategory)
	protected.POST("/createCategories", categoryHandler.CreateCategories)
	protected.DELETE("/delete/:id", categoryHandler.DeleteCategory)
	protected.POST("/update", categoryHandler.UpdateCategory)
	protected.PUT("/addParent/:son_id/:father_id", categoryHandler.AddParentCategory)
	protected.PUT("/removeParent/:son_id", categoryHandler.RemoveParentCategory)
}
