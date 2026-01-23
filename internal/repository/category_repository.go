package repository

import (
	"crist-blog/internal/model"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetNameByID(id uuid.UUID) (string, error) {
	var name string
	err := r.DB.Model(&model.Category{}).
		Select("name").
		Where("id = ?", id).
		First(&name).Error
	return name, err
}

func (r *CategoryRepository) ListAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.DB.Model(&model.Category{}).
		Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) CreateCategory(category *model.Category) error {
	return r.DB.Create(category).Error
}

func (r *CategoryRepository) CreatCategories(categories []model.Category) error {
	if len(categories) == 0 {
		return errors.New("categories is empty")
	}
	return r.DB.Create(&categories).Error
}

func (r *CategoryRepository) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category
	err := r.DB.Model(&model.Category{}).
		Where("id = ?", id).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, err
}

func (r *CategoryRepository) DeleteCategory(id uuid.UUID) error {
	_, err := r.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	err = r.DB.Model(&model.Category{}).
		Where("id = ?", id).
		Update("deleted_flag", true).Error
	return err
}

// UpdateCategory 更新分类
func (r *CategoryRepository) UpdateCategory(id uuid.UUID, category *model.Category) error {
	_, err := r.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	err = r.DB.Model(&model.Category{}).
		Where("id = ?", id).
		Updates(category).Error
	return err
}

// AddParentCategory 给分类添加父分类
func (r *CategoryRepository) AddParentCategory(parentId, sonId uuid.UUID) error {
	_, err := r.GetCategoryByID(parentId)
	if err != nil {
		return errors.New("parent category not found")
	}
	_, err = r.GetCategoryByID(sonId)
	if err != nil {
		return errors.New("son category not found")
	}
	err = r.DB.Model(&model.Category{}).
		Where("id = ?", sonId).
		Update("parent_id", parentId).Error
	return err
}

// GetFatherCategoryById 输出分类的父分类
func (r *CategoryRepository) GetFatherCategoryById(id uuid.UUID) (*model.Category, error) {
	category, err := r.GetCategoryByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if category.ParentID == uuid.Nil {
		return nil, errors.New("have no father")
	}

	fatherCategory, err := r.GetCategoryByID(category.ParentID)
	if err != nil {
		return nil, errors.New("father category not found")
	}
	return fatherCategory, err

}

func (r *CategoryRepository) RemoveFatherCategory(id uuid.UUID) error {
	err := r.DB.Model(&model.Category{}).
		Where("id = ?", id).
		Update("parent_id", uuid.Nil).Error
	if err != nil {
		return err
	}
	return nil
}
