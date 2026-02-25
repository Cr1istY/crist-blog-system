package repository

import (
	"context"
	"crist-blog/internal/model"
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type TweetRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewTweetRepository(db *gorm.DB, redis *redis.Client) *TweetRepository {
	return &TweetRepository{
		db:    db,
		redis: redis,
	}
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
	key := "tweets:" + strconv.Itoa(limit) + ":" + strconv.Itoa(offset)
	// 在redis中查询
	data, err := r.redis.Get(ctx, key).Bytes()
	if err == nil {
		// 找到了，直接返回
		if err = json.Unmarshal(data, &tweets); err != nil {
			// 缓存数据损坏
			return nil, err
		}
		return tweets, nil
	}

	// 查询推文，存入redis
	err = r.db.WithContext(ctx).
		Preload("Images").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tweets).Error
	if err != nil {
		return nil, err
	}
	if len(tweets) <= 0 {
		return nil, errors.New("no tweets found")
	}
	ttl := time.Hour*12 + time.Duration(rand.Intn(300))*time.Second
	data, err = json.Marshal(tweets)
	if err != nil {
		return nil, errors.New("marshal tweets error")
	}
	err = r.redis.Set(ctx, key, data, ttl).Err()
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
	// 刷新redis缓存
	if err != nil {
		return err
	}

	return r.ClearTweetCache(ctx)
}

func (r *TweetRepository) ClearTweetCache(ctx context.Context) error {
	// 1. 获取所有匹配 tweets:* 的键
	keys, err := r.redis.Keys(ctx, "tweets:*").Result()
	if err != nil {
		return err
	}

	// 2. 如果没有匹配的键，直接返回
	if len(keys) == 0 {
		return nil
	}

	// 3. 删除所有匹配的键
	return r.redis.Del(ctx, keys...).Err()
}
