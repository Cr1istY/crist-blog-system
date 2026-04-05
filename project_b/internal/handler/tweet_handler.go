package handler

import (
	"crist-blog/internal/model"
	"crist-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TweetHandler struct {
	tweetService *service.TweetService
	userService  *service.UserService
}

func NewTweetHandler(tweetService *service.TweetService, userService *service.UserService) *TweetHandler {
	return &TweetHandler{
		tweetService: tweetService,
		userService:  userService,
	}
}

func (h *TweetHandler) CreateTweet(c echo.Context) error {
	ctx := c.Request().Context()
	var req model.CreateTweetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "请求参数错误"})
	}

	userID, ok := c.Get("user_id_str").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取用户ID失败"})
	}

	tweet := &model.Tweet{
		ID:      uuid.New().String(),
		UserID:  userID,
		Content: req.Content,
	}

	if err := h.tweetService.CreateTweetWithImages(ctx, tweet, req.ImageIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "创建推文失败"})
	}

	tweetWithImages, err := h.tweetService.GetTweetWithImagesByID(ctx, tweet.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取推文失败"})
	}

	return c.JSON(http.StatusCreated, h.toTweetResponse(tweetWithImages))

}

func (h *TweetHandler) GetAllTweets(c echo.Context) error {
	ctx := c.Request().Context()

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 {
		limit = 20
	}

	if limit > 50 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	tweets, err := h.tweetService.GetAllWithImages(ctx, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取推文失败"})
	}

	tweetResponses := make([]model.TweetResponse, 0, len(tweets))
	for _, tweet := range tweets {
		if tweet.DeletedFlag {
			continue
		}
		tweetResponses = append(tweetResponses, h.toTweetResponse(&tweet))

	}

	return c.JSON(http.StatusOK, tweetResponses)
}

func (h *TweetHandler) GetCurrentUserInTweet(c echo.Context) error {
	userIDStr, ok := c.Get("user_id_str").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取用户失败"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取用户uuid失败"})
	}
	user, err := h.userService.GetCurrentTweetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "在数据库中获取用户失败"})
	}
	return c.JSON(http.StatusOK, user)

}

func (h *TweetHandler) DeleteTweetByID(c echo.Context) error {
	tweetID := c.Param("id")
	if tweetID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "推文ID不能为空"})
	}
	// 检查是否为作者
	// 获取当前用户ID
	userIDStr, ok := c.Get("user_id_str").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "获取用户失败"})
	}
	ctx := c.Request().Context()
	err := h.tweetService.DeleteTweet(ctx, tweetID, userIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "删除推文失败"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "推文删除成功"})
}

func (h *TweetHandler) toTweetResponse(tweet *model.Tweet) model.TweetResponse {
	imageURLs := make([]string, 0, len(tweet.Images))
	for _, img := range tweet.Images {
		imageURLs = append(imageURLs, img.URL)
	}

	userID, err := uuid.Parse(tweet.UserID)
	if err != nil {
		return model.TweetResponse{}
	}

	user, _ := h.userService.GetCurrentTweetUserByID(userID)

	return model.TweetResponse{
		ID:        tweet.ID,
		User:      *user,
		Content:   tweet.Content,
		Timestamp: tweet.CreatedAt,
		Likes:     tweet.Likes,
		Images:    imageURLs,
	}
}
