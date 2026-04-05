package repository

import (
	"crist-blog/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageRepository struct {
	DB *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{DB: db}
}

func (r *ImageRepository) Create(img *model.Image) error {
	return r.DB.Create(img).Error
}

func (r *ImageRepository) FindByID(id uuid.UUID) (*model.Image, error) {
	var image model.Image
	err := r.DB.Where("id = ?", id).First(&image).Error
	return &image, err
}
