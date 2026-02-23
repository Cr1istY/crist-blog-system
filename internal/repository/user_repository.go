package repository

import (
	"crist-blog/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
func (r *UserRepository) GetByName(name string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", name).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}
