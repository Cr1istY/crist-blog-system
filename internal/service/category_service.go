package service

import (
	"crist-blog/internal/model"
	"crist-blog/internal/repository"

	"github.com/google/uuid"
)

type CategoryService struct {
	CategoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		CategoryRepo: categoryRepo,
	}
}

func (s *CategoryService) GetNameByID(id uuid.UUID) (string, error) {
	return s.CategoryRepo.GetNameByID(id)
}

func (s *CategoryService) ListAllCategories() ([]model.Category, error) {
	return s.CategoryRepo.ListAllCategories()
}

func (s *CategoryService) CreateCategory(category *model.Category) error {
	return s.CategoryRepo.CreateCategory(category)
}

func (s *CategoryService) CreatCategories(categories []model.Category) error {
	return s.CategoryRepo.CreatCategories(categories)
}

func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	return s.CategoryRepo.GetCategoryByID(id)
}

func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	return s.CategoryRepo.DeleteCategory(id)
}

func (s *CategoryService) UpdateCategory(id uuid.UUID, category *model.Category) error {
	return s.CategoryRepo.UpdateCategory(id, category)
}

func (s *CategoryService) AddParentCategory(parentId, sonId uuid.UUID) error {
	return s.CategoryRepo.AddParentCategory(parentId, sonId)
}

func (s *CategoryService) GetFatherCategoryById(id uuid.UUID) (*model.Category, error) {
	return s.CategoryRepo.GetFatherCategoryById(id)
}

func (s *CategoryService) RemoveFatherCategory(id uuid.UUID) error {
	return s.CategoryRepo.RemoveFatherCategory(id)
}
