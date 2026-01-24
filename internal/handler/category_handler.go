package handler

import (
	"crist-blog/internal/model"
	"crist-blog/internal/service"
	"crist-blog/internal/utils"
	"net/http"

	"github.com/google/uuid"
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
		if !rawCategory.DeletedFlag {
			continue
		}
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
	slug := utils.ToSlug(createCategory.Name)
	var category = model.Category{
		Name:        createCategory.Name,
		Description: createCategory.Description,
		ParentID:    createCategory.ParentID,
		Slug:        slug,
	}
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	err = h.categoryService.CreateCategory(&category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "please check your categories"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "create category successfully", "category_id": category.ID.String()})
}

func (h *CategoryHandler) CreateCategories(c echo.Context) error {
	var createCategories []model.CreateCategory
	err := c.Bind(&createCategories)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong request body of categories in bind:)"})
	}
	var categories []model.Category
	for _, rawCategory := range createCategories {
		category := model.Category{
			Name:        rawCategory.Name,
			Description: rawCategory.Description,
			ParentID:    rawCategory.ParentID,
		}
		categories = append(categories, category)
	}
	err = h.categoryService.CreateCategories(categories)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong request body of categories in create:)"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "create categories successfully"})
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong category id in param:)"})
	}
	err = h.categoryService.DeleteCategory(categoryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong category id in delete:)"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "delete category successfully"})

}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	var updateCategory model.UpdateCategory
	err := c.Bind(&updateCategory)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong category id in param:)"})
	}
	var category = model.Category{
		ID:          updateCategory.ID,
		Name:        updateCategory.Name,
		Description: updateCategory.Description,
		ParentID:    updateCategory.ParentID,
	}
	err = h.categoryService.UpdateCategory(&category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong category id in update:)"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "update category successfully"})
}

func (h *CategoryHandler) AddParentCategory(c echo.Context) error {
	sonID, err := uuid.Parse(c.Param("son_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong son category id in param:)"})
	}
	fatherID, err := uuid.Parse(c.Param("father_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong father category id in param:)"})
	}
	if fatherID == sonID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "father and son category id can't be the same:)"})
	}
	fatherCat, err := h.categoryService.GetCategoryByID(fatherID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong father category id in database"})
	}
	var fatherCategories []model.CategoryToFrontend
	var times = 0 // 限制循环次数，防止死循环
	var flag = true
	if fatherCat.ParentID == uuid.Nil {
		flag = false
	}
	for flag {
		rowFatherCategory, err := h.categoryService.GetFatherCategoryById(fatherID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong in get father category"})
		}
		var fatherCategory = model.CategoryToFrontend{
			ID:          rowFatherCategory.ID,
			Name:        rowFatherCategory.Name,
			Slug:        rowFatherCategory.Slug,
			Description: rowFatherCategory.Description,
		}
		fatherCategories = append(fatherCategories, fatherCategory)
		if rowFatherCategory.ParentID == uuid.Nil {
			flag = false
		}
		times++
		if times > 10 {
			flag = false
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "it has too many fathers"})
		}
	}
	for _, father := range fatherCategories {
		if father.ID == sonID {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "father`s father cannot be father`s son"})
		}
	}
	err = h.categoryService.AddParentCategory(fatherID, sonID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong in add parent:)"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "add parent category successfully:)"})
}

func (h *CategoryHandler) RemoveParentCategory(c echo.Context) error {
	sonID, err := uuid.Parse(c.Param("son_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong son category id in param"})
	}
	err = h.categoryService.RemoveFatherCategory(sonID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong in remove parent"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "remove parent category successfully"})
}

// 循环输出父亲路径

func (h *CategoryHandler) GetFatherCategory(c echo.Context) error {
	finalSonID, err := uuid.Parse(c.Param("son_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong son category id in param"})
	}
	var flag = true
	finalSonCategory, err := h.categoryService.GetCategoryByID(finalSonID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong in get son category"})
	}
	if finalSonCategory.ParentID == uuid.Nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "have no father"})
	}
	var fatherCategories []model.CategoryToFrontend
	var times = 0
	var tempID = finalSonCategory.ParentID
	for flag {
		rowFatherCategory, err := h.categoryService.GetCategoryByID(tempID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong in get father category"})
		}
		var fatherCategory = model.CategoryToFrontend{
			ID:          rowFatherCategory.ID,
			Name:        rowFatherCategory.Name,
			Slug:        rowFatherCategory.Slug,
			Description: rowFatherCategory.Description,
		}
		fatherCategories = append(fatherCategories, fatherCategory)
		if rowFatherCategory.ParentID == uuid.Nil {
			flag = false
		} else {
			tempID = rowFatherCategory.ParentID
		}
		times++
		if times > 10 {
			flag = false
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "too many times"})
		}
	}
	return c.JSON(http.StatusOK, fatherCategories)
}
