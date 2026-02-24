package repository

import (
	"context"
	"crist-blog/internal/model"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB  *gorm.DB
	rdb *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{
		DB:  db,
		rdb: redis,
	}
}
func (r *UserRepository) GetByName(name string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", name).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	ctx := context.Background()
	// 先从redis中查找
	key := "user:" + id.String()
	data, err := r.rdb.Get(ctx, key).Bytes()
	if err == nil {
		// 找到了，直接返回
		if err = json.Unmarshal(data, &user); err != nil {
			// 缓存数据损坏
			log.Warnf("缓存数据损坏: %v", err)
		} else {
			return &user, nil
		}
	}
	// 没找到，从数据库中查找
	err = r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Errorf("数据库查询错误: %v", err)
		return nil, err
	}
	// 写入redis，附加随机时间，防止缓存雪崩
	ttl := time.Hour*12 + time.Duration(rand.Intn(300))*time.Second
	data, err = json.Marshal(user)
	if err != nil {
		log.Errorf("序列化错误: %v", err)
	} else {
		r.rdb.Set(ctx, key, data, ttl)
	}

	// 返回
	return &user, err
}
