package handler

import (
	"crist-blog/internal/model"
	"crist-blog/internal/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TweetHandler struct {
	tweetService *service.TweetService
}

func NewTweetHandler(tweetService *service.TweetService) *TweetHandler {
	return &TweetHandler{
		tweetService: tweetService,
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

func (h *TweetHandler) toTweetResponse(tweet *model.Tweet) model.TweetResponse {
	imageURLs := make([]string, 0, len(tweet.Images))
	for _, img := range tweet.Images {
		imageURLs = append(imageURLs, img.URL)
	}

	return model.TweetResponse{
		ID:        tweet.ID,
		Content:   tweet.Content,
		Timestamp: tweet.CreatedAt,
		Likes:     tweet.Likes,
		Images:    imageURLs,
	}
}
