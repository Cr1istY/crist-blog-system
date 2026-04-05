package service

import (
	"crist-blog/internal/model"
	"crist-blog/internal/repository"

	"github.com/google/uuid"
)

type ImageService struct {
	ImageRepo *repository.ImageRepository
}

func NewImageService(imageService *repository.ImageRepository) *ImageService {
	return &ImageService{
		ImageRepo: imageService,
	}
}

func (s *ImageService) CreateImage(image *model.Image) error {
	return s.ImageRepo.Create(image)
}

func (s *ImageService) GetImageByID(id uuid.UUID) (*model.Image, error) {
	return s.ImageRepo.FindByID(id)
}
