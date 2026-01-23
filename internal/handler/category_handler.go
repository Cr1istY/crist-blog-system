package handler

import (
	"crist-blog/internal/model"
	"crist-blog/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) ListAllCategories(c echo.Context) error {
	var categories []model.CreatePostCategory
	rawCategories, err := h.categoryService.ListAllCategories()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "network error, so there cannot list all categories:)"})
	}
	for _, rawCategory := range rawCategories {
		categories = append(categories, model.CreatePostCategory{
			ID:   rawCategory.ID,
			Name: rawCategory.Name,
		})
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var createCategory model.CreateCategory
	err := c.Bind(&createCategory)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong request body of category in create:)"})
	}
	var category = model.Category{
		Name:        createCategory.Name,
		Description: createCategory.Description,
		ParentID:    createCategory.ParentID,
	}
	err = h.categoryService.CreateCategory(&category)
	return c.JSON(http.StatusOK, map[string]string{"message": "create category successfully", "category_id": category.ID.String()})
}
