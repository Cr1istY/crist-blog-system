package handler

import (
	"context"
	"crist-blog/internal/assets"
	"crist-blog/internal/model"
	"crist-blog/internal/service"
	"crist-blog/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService     *service.PostService
	categoryService *service.CategoryService
	rdb             *redis.Client
	cacheKey        string
}

func NewPostHandler(postService *service.PostService, categoryService *service.CategoryService, redisService *redis.Client) *PostHandler {
	return &PostHandler{
		postService:     postService,
		categoryService: categoryService,
		rdb:             redisService,
		cacheKey:        "blog:posts:frontend:list:v1",
	}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var req model.CreatePostRequest
	userId := c.Get("user_id_str").(string)
	// userId := "00000000-0000-0000-0000-000000000001"
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
	if req.Title != "" {
		title = req.Title
	}
	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title is required"})
	}
	req.Title = title
	req.MetaTitle = title
	slug, err := h.generateSlug(title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate slug when create"})
	}
	req.Slug = slug
	excerpt := req.Excerpt
	if len([]rune(req.Excerpt)) > 20 {
		runes := []rune(excerpt)
		if len(runes) > 20 {
			excerpt = string(runes[:20]) + "..."
		}
	}
	req.MetaDescription = excerpt
	if req.Thumbnail == "" {
		req.Thumbnail = assets.GetThumbnail()
	}
	post := &model.Post{
		UserID:          userID,
		Title:           req.Title,
		Slug:            req.Slug,
		Content:         req.Content,
		Excerpt:         excerpt,
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

	go func() {
		_, err := h.postsSaveToRedis(h.cacheKey)
		if err != nil {
			println("Error while saving post list to redis", err.Error())
			return
		}
	}()

	println("Post list cache refreshed successfully after create",
		"post_id", post.ID, "cache_key", h.cacheKey)
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
	var category string
	if post.CategoryID != uuid.Nil && post.CategoryID.String() != model.RootCategoryID {
		var categoryNames []string
		rawCategory, err := h.categoryService.GetNameByID(post.CategoryID)
		if err != nil {
			category = ""
		}
		categoryNames = append(categoryNames, rawCategory)
		sonCategoryID := post.CategoryID
		for sonCategoryID != uuid.Nil && sonCategoryID.String() != model.RootCategoryID {
			fatherCategory, err := h.categoryService.GetFatherCategoryById(sonCategoryID)
			if err != nil {
				break
			}
			if fatherCategory.ID == uuid.Nil || fatherCategory.ID.String() == model.RootCategoryID {
				break
			}
			fatherCategoryName, err := h.categoryService.GetNameByID(fatherCategory.ID)
			if err != nil {
				break
			}
			categoryNames = append(categoryNames, fatherCategoryName)
			sonCategoryID = fatherCategory.ID
		}
		// 将 categoryNames 处理成 父/子/孙 的形式
		for i := len(categoryNames) - 1; i >= 0; i-- {
			if i == len(categoryNames)-1 {
				category += categoryNames[i]
			} else {
				category += "/" + categoryNames[i]
			}
		}
	}
	//categoryName, err := h.categoryService.GetNameByID(post.CategoryID)
	//if err != nil {
	//	categoryName = "未分类"
	//}
	if post.Thumbnail == "" {
		post.Thumbnail = assets.GetThumbnail()
	}
	var postToViewers = &model.PostDetail{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		Date:            dateStr,
		Tags:            post.Tags,
		CategoryID:      post.CategoryID.String(),
		Category:        category,
		Views:           post.Views,
		Likes:           post.Likes,
		Excerpt:         post.Excerpt,
		Status:          post.Status,
		MetaTitle:       post.MetaTitle,
		MetaDescription: post.MetaDescription,
		Thumbnail:       post.Thumbnail,
		IsPinned:        post.IsPinned,
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

func (h *PostHandler) GetBlogBySlugWithPinned(c echo.Context) error {
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
	req.Slug = rawPost.Slug
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

	go func() {
		_, goErr := h.postsSaveToRedis(h.cacheKey)
		if goErr != nil {
			println("Error saving posts to redis: ", goErr.Error())
			return
		}
		println("successful saving to redis when update")
	}()

	return c.JSON(http.StatusOK, map[string]string{"post_slug": post.Slug})
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
	go func() {
		_, err := h.postsSaveToRedis(h.cacheKey)
		if err != nil {
			println("Error saving posts to redis: ", err.Error())
		}
	}()
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

// ListToFrontend 已弃用，见ListToFrontendWithPinned
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

// 使用 redis 缓存博客文章JSON，当有文章创建时，重写写入 redis

func (h *PostHandler) ListToFrontendWithPinned(c echo.Context) error {
	// 1. 创建 redis 键
	ctx := c.Request().Context()
	cacheKey := h.cacheKey
	// 2. 检查是否已经存在
	cacheData, err := h.rdb.Get(ctx, cacheKey).Bytes()
	if err == nil && len(cacheData) > 0 {
		// 3. 如果存在，直接返回
		var blogPosts []*model.PostFrontendWithPinned
		if json.Unmarshal(cacheData, &blogPosts) == nil {
			return c.JSON(http.StatusOK, blogPosts)
		}
		println("Redis cache parse failed, rebuilding...", "key", cacheKey, "err", err)
	}
	// 4. 如果不存在，写入 redis，并返回
	posts, err := h.postService.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	blogPosts, err := h.blogPostsToPostWithPinned(posts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error: ": "blogPostsToPostWithPinned"})
	}
	// 异步写入缓存
	go func() {
		data, err := json.Marshal(blogPosts)
		if err != nil {
			println("Failed to marshal posts for cache", "err", err)
			return
		}
		if err := h.rdb.Set(context.Background(), cacheKey, data, 12*time.Hour).Err(); err != nil {
			println("Failed to set posts cache", "err", err)
		}
	}()
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

func (h *PostHandler) PinPost(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}
	// TODO 实装 pinnedOrder, PinedUntil
	err = h.postService.PinPost(uint(id), 1, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Post pinned successfully"})
}

func (h *PostHandler) UnpinPost(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}
	err = h.postService.UnpinPost(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Post unpinned successfully"})
}

func (h *PostHandler) blogPostsToPostWithPinned(posts []*model.Post) ([]*model.PostFrontendWithPinned, error) {
	var blogPosts []*model.PostFrontendWithPinned
	for _, post := range posts {
		if post.Status != model.Published {
			continue
		}
		if post.Thumbnail == "" {
			post.Thumbnail = assets.GetThumbnail()
		}
		category := ""
		if post.CategoryID != uuid.Nil && post.CategoryID.String() != model.RootCategoryID {
			category, _ = h.categoryService.GetNameByID(post.CategoryID)
		}
		blogPosts = append(blogPosts, &model.PostFrontendWithPinned{
			ID:          post.ID,
			Slug:        post.Slug,
			Title:       post.Title,
			Category:    category,
			Tags:        post.Tags,
			Date:        post.PublishedAt.Format("2006-01-02"),
			Excerpt:     post.Excerpt,
			Views:       post.Views,
			Likes:       post.Likes,
			Thumbnail:   post.Thumbnail,
			IsPinned:    post.IsPinned,
			PinnedOrder: post.PinnedOrder,
		})
	}
	return blogPosts, nil
}

func (h *PostHandler) postsSaveToRedis(cacheKey string) ([]*model.PostFrontendWithPinned, error) {
	ctx := context.Background()

	posts, goErr := h.postService.List()
	if goErr != nil {
		println("Failed to list posts for cache refresh", "err", goErr)
		return nil, goErr
	}

	frontendPosts, goErr := h.blogPostsToPostWithPinned(posts)
	if goErr != nil {
		println("Failed to convert posts for cache refresh", "err", goErr)
		return frontendPosts, goErr
	}
	// 写入 redis

	data, goErr := json.Marshal(frontendPosts)
	if goErr != nil {
		println("Failed to marshal posts for cache", "err", goErr)
		return frontendPosts, goErr
	}

	if goErr = h.rdb.Set(ctx, cacheKey, data, 12*time.Hour).Err(); goErr != nil {
		println("Failed to refresh post list cache", "key", cacheKey, "err", goErr)
		return frontendPosts, goErr
	}
	return frontendPosts, nil
}

func (h *PostHandler) generateSlug(title string) (string, error) {
	slug, err := utils.ToSlug(title)
	if err != nil {
		return "", err
	}

	for attemptGenerateSlug := 0; attemptGenerateSlug < 6; attemptGenerateSlug++ {
		_, err = h.postService.GetBySlug(slug)
		if err != nil && errors.Is(model.ErrSlugNotFound, err) {
			break
		}
		if err != nil {
			// 此时，是数据库的其他问题，应该直接退出程序
			return "", err
		}
		// slug重复，生成一个末尾带有随机数的新slug
		if attemptGenerateSlug == 5 {
			// 重复次数过多，直接退出
			return "", errors.New("try too much times when generate slug")
		}
		slug = utils.SlugToSlugWithRandom(slug)
	}

	return slug, nil

}
