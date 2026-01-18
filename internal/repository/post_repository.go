package repository

import (
	"crist-blog/internal/model"

	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) CreatePost(post *model.Post) error {
	return r.DB.Create(post).Error
}

func (r *PostRepository) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.DB.Where("id = ?", id).First(&post).Error
	return &post, err
}

func (r *PostRepository) GetBySlug(slug string) (*model.Post, error) {
	var post model.Post
	err := r.DB.Where("slug = ?", slug).First(&post).Error
	return &post, err
}

func (r *PostRepository) Update(post *model.Post) error {
	return r.DB.Save(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.DB.Where("id = ?", id).Delete(&model.Post{}).Error
}

func (r *PostRepository) List() ([]*model.Post, error) {
	var posts []*model.Post
	return posts, r.DB.Find(&posts).Error
}

func (r *PostRepository) GetHotPost() ([]*model.HotPost, error) {
	var hotPosts []*model.HotPost
	err := r.DB.Model(&model.Post{}).
		Select("id, title, category_id, created_at, excerpt").
		Where("status = ?", model.Published).
		Order("likes desc").
		Limit(2).
		Find(&hotPosts).Error
	return hotPosts, err
}

func (r *PostRepository) GetLatestPosts() ([]*model.LatestPost, error) {
	var latestPosts []*model.LatestPost
	err := r.DB.Model(&model.Post{}).
		Select("id, title, category_id, created_at").
		Where("status = ?", model.Published).
		Order("created_at desc").
		Limit(2).
		Find(&latestPosts).Error
	return latestPosts, err
}

func (r *PostRepository) AddViews(id uint) error {
	return r.DB.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

func (r *PostRepository) AddLikes(id uint) error {
	return r.DB.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
}
