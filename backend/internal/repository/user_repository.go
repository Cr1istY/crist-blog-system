package repository

import (
	"context"
	"crist-blog/internal/model"
	"encoding/json"
	"errors"
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

func (r *UserRepository) ChangeUserInfo(id uuid.UUID, user *model.User) error {
	// 修改
	err := r.DB.Model(&model.User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		return err
	}
	// 写入redis
	ctx := context.Background()
	go func() {
		err = r.saveUserToRedis(ctx, user)
		if err != nil {
			log.Warn("save user to redis error: ", err)
		}
	}()
	return nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user *model.User
	ctx := context.Background()
	// 先从redis中查找
	user, err := r.getUserFromRedis(ctx, id)
	if err == nil {
		return user, err
	}
	// 没找到，从数据库中查找
	err = r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Errorf("数据库查询错误: %v", err)
		return nil, err
	}
	// 写入redis，附加随机时间，防止缓存雪崩
	go func() {
		err = r.saveUserToRedis(ctx, user)
		if err != nil {
			log.Errorf("redis写入错误: %v", err)
		}
	}()
	// 返回
	return user, err
}

func (r *UserRepository) saveUserToRedis(ctx context.Context, user *model.User) error {
	key := "user:" + user.ID.String()
	ttl := time.Hour*12 + time.Duration(rand.Intn(300))*time.Second
	data, err := json.Marshal(user)
	if err != nil {
		return errors.New("序列化错误")
	}
	return r.rdb.Set(ctx, key, data, ttl).Err()

}

func (r *UserRepository) getUserFromRedis(ctx context.Context, id uuid.UUID) (*model.User, error) {
	key := "user:" + id.String()
	var user model.User
	data, err := r.rdb.Get(ctx, key).Bytes()
	if err == nil {
		// 找到了，直接返回
		if err = json.Unmarshal(data, &user); err != nil {
			// 缓存数据损坏
			return nil, err
		}
		return &user, nil
	}
	return nil, err
}
