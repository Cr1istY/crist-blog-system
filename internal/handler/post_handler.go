package handler

import (
	"crist-blog/internal/assets"
	"crist-blog/internal/model"
	"crist-blog/internal/service"
	"crist-blog/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService     *service.PostService
	categoryService *service.CategoryService
}

func NewPostHandler(postService *service.PostService, categoryService *service.CategoryService) *PostHandler {
	return &PostHandler{
		postService:     postService,
		categoryService: categoryService,
	}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var req model.CreatePostRequest
	userId := "00000000-0000-0000-0000-000000000001"
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if req.UserID == "" {
		req.UserID = userId
	}
	if _, err := uuid.Parse(req.UserID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	if _, err := uuid.Parse(req.CategoryID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}
	defaultValue := 0
	userID, err := uuid.Parse(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	title := utils.ExtractPostTitle(req.Content)
	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title is required"})
	}
	req.Title = title
	req.MetaTitle = title
	slug := utils.ToSlug(title)
	req.Slug = slug
	req.MetaDescription = req.Excerpt
	if req.Thumbnail == "" {
		req.Thumbnail = assets.GetThumbnail()
	}
	post := &model.Post{
		UserID:          userID,
		Title:           req.Title,
		Slug:            req.Slug,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		Status:          model.PostStatus(req.Status),
		CategoryID:      uuid.MustParse(req.CategoryID),
		Tags:            req.Tags,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		Views:           defaultValue,
		Likes:           defaultValue,
		Thumbnail:       req.Thumbnail,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := h.postService.CreatePost(post); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Post created successfully", "id": post.ID})
}

func (h *PostHandler) GetBlogById(c echo.Context) error {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	post, err := h.postService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, post)
}

func (h *PostHandler) uniformizePostToViewers(post *model.Post) *model.PostDetail {
	dateStr := ""
	if post.PublishedAt != nil {
		// 格式：2025年12月15日
		dateStr = post.UpdatedAt.Format("2006-1-2")
	} else {
		dateStr = post.CreatedAt.Format("2006-1-2")
	}
	categoryName, err := h.categoryService.GetNameByID(post.CategoryID)
	if err != nil {
		categoryName = "未分类"
	}
	if post.Thumbnail == "" {
		post.Thumbnail = assets.GetThumbnail()
	}
	var postToViewers = &model.PostDetail{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		Date:            dateStr,
		Tags:            post.Tags,
		CategoryID:      post.CategoryID,
		Category:        categoryName,
		Views:           post.Views,
		Likes:           post.Likes,
		Excerpt:         post.Excerpt,
		Status:          post.Status,
		MetaTitle:       post.MetaTitle,
		MetaDescription: post.MetaDescription,
		Thumbnail:       post.Thumbnail,
	}
	return postToViewers
}

func (h *PostHandler) GetBlogBySlug(c echo.Context) error {
	slug := c.Param("slug")
	post, err := h.postService.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	postToViewers := h.uniformizePostToViewers(post)
	return c.JSON(http.StatusOK, postToViewers)
}

func (h *PostHandler) GetBlogByIdToViewers(c echo.Context) error {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	post, err := h.postService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	postToViewers := h.uniformizePostToViewers(post)
	return c.JSON(http.StatusOK, postToViewers)
}

func (h *PostHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	var req model.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	rawPost, err := h.postService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	title := utils.ExtractPostTitle(req.Content)
	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is empty"})
	}
	req.Title = title
	req.MetaTitle = title
	slug := utils.ToSlug(title)
	req.Slug = slug
	req.MetaDescription = req.Excerpt
	post := &model.Post{
		ID:              id,
		UserID:          rawPost.UserID,
		Title:           req.Title,
		Slug:            req.Slug,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		Status:          model.PostStatus(req.Status),
		CategoryID:      uuid.MustParse(req.CategoryID),
		Views:           rawPost.Views,
		Likes:           rawPost.Likes,
		Tags:            req.Tags,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		Thumbnail:       req.Thumbnail,
	}
	if rawPost.PublishedAt != nil {
		post.PublishedAt = rawPost.PublishedAt
	}

	if err := h.postService.Update(post); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	updated, _ := h.postService.GetByID(id)
	return c.JSON(http.StatusOK, updated)
}

func (h *PostHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := h.postService.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// List 已弃用 见ListToFrontend
func (h *PostHandler) List(c echo.Context) error {
	posts, err := h.postService.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	for _, post := range posts {
		post.Content = ""
	}
	return c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) ListToFrontend(c echo.Context) error {
	posts, err := h.postService.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	var blogPosts []*model.PostFrontend
	for _, post := range posts {
		if post.Status != model.Published {
			continue
		}
		if post.Thumbnail == "" {
			post.Thumbnail = assets.GetThumbnail()
		}
		blogPosts = append(blogPosts, &model.PostFrontend{
			ID:        post.ID,
			Slug:      post.Slug,
			Title:     post.Title,
			Tags:      post.Tags,
			Date:      post.PublishedAt.Format("2006-01-02"),
			Excerpt:   post.Excerpt,
			Views:     post.Views,
			Likes:     post.Likes,
			Thumbnail: post.Thumbnail,
		})
	}
	return c.JSON(http.StatusOK, blogPosts)
}

func (h *PostHandler) GetHotPosts(c echo.Context) error {
	posts, err := h.postService.GetHotPosts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	var blogPosts []*model.HotPostFrontend
	for _, post := range posts {
		category, err := h.categoryService.GetNameByID(post.CategoryID)
		if err != nil {
			category = "未分类"
		}
		blogPosts = append(blogPosts, &model.HotPostFrontend{
			Slug:     post.Slug,
			Title:    post.Title,
			Category: category,
			Date:     post.CreatedAt.Format("2006-01-02"),
			Excerpt:  post.Excerpt,
		})
	}
	return c.JSON(http.StatusOK, blogPosts)
}

func (h *PostHandler) GetLatestPosts(c echo.Context) error {
	posts, err := h.postService.GetLatestPosts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	var blogPosts []*model.LatestPostFrontend
	for _, post := range posts {
		category, err := h.categoryService.GetNameByID(post.CategoryID)
		if err != nil {
			category = "未分类"
		}
		blogPosts = append(blogPosts, &model.LatestPostFrontend{
			Slug:     post.Slug,
			Title:    post.Title,
			Category: category,
			Date:     post.CreatedAt.Format("2006-01-02"),
		})
	}
	return c.JSON(http.StatusOK, blogPosts)
}

func (h *PostHandler) AddViews(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}
	err = h.postService.AddViews(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Views added successfully"})
}

func (h *PostHandler) AddLikes(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}
	err = h.postService.AddLikes(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Likes added successfully"})
}
