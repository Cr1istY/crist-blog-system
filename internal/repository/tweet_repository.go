package repository

import (
	"context"
	"crist-blog/internal/model"

	"gorm.io/gorm"
)

type TweetRepository struct {
	db *gorm.DB
}

func NewTweetRepository(db *gorm.DB) *TweetRepository {
	return &TweetRepository{db: db}
}

func (r *TweetRepository) CreateWithImages(ctx context.Context, tweet *model.Tweet, imageIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建推文
		if err := tx.Create(tweet).Error; err != nil {
			return err
		}

		// 关联图片
		if len(imageIDs) > 0 {
			tweetImages := make([]model.TweetImage, 0, len(imageIDs))
			for i, imgID := range imageIDs {
				tweetImages = append(tweetImages, model.TweetImage{
					TweetID:      tweet.ID,
					ImageID:      imgID,
					DisplayOrder: i,
				})
			}
			if err := tx.Create(&tweetImages).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TweetRepository) GetAllWithImages(ctx context.Context, limit, offset int) ([]model.Tweet, error) {
	var tweets []model.Tweet
	err := r.db.WithContext(ctx).
		Preload("Images").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tweets).Error
	return tweets, err
}

// GetByIDWithImages 根据 ID 获取推文（带图片）
func (r *TweetRepository) GetByIDWithImages(ctx context.Context, id string) (*model.Tweet, error) {
	var tweet model.Tweet
	err := r.db.WithContext(ctx).
		Preload("Images").
		First(&tweet, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tweet, nil
}

func (r *TweetRepository) DeleteByID(ctx context.Context, id, userID string) error {
	err := r.db.WithContext(ctx).Model(&model.Tweet{}).Where("id = ? AND user_id = ?", id, userID).Update("deleted_flag", true).Error
	return err
}
