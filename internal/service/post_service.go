package service

import (
	"crist-blog/internal/model"
	"crist-blog/internal/repository"
	"time"
)

type PostService struct {
	PostRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{
		PostRepo: postRepo,
	}
}

func (s *PostService) CreatePost(post *model.Post) error {
	if post.Status == model.Published && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}
	return s.PostRepo.CreatePost(post)
}

func (s *PostService) GetByID(id uint) (*model.Post, error) {
	return s.PostRepo.GetByID(id)
}

func (s *PostService) Update(post *model.Post) error {
	existing, err := s.GetByID(post.ID)
	if err != nil {
		return err
	}
	existing.Title = post.Title
	existing.Slug = post.Slug
	existing.Content = post.Content
	existing.Excerpt = post.Excerpt
	existing.Status = post.Status
	existing.CategoryID = post.CategoryID
	existing.Tags = post.Tags
	existing.MetaTitle = post.MetaTitle
	existing.MetaDescription = post.MetaDescription
	existing.Thumbnail = post.Thumbnail

	if existing.Status == model.Published && existing.PublishedAt == nil {
		now := time.Now()
		existing.PublishedAt = &now
	}
	return s.PostRepo.Update(existing)
}

func (s *PostService) Delete(id uint) error {
	return s.PostRepo.Delete(id)
}

func (s *PostService) List() ([]*model.Post, error) {

	return s.PostRepo.List()
}

func (s *PostService) GetHotPosts() ([]*model.HotPost, error) {
	return s.PostRepo.GetHotPost()
}

func (s *PostService) GetLatestPosts() ([]*model.LatestPost, error) {
	return s.PostRepo.GetLatestPosts()
}

func (s *PostService) AddViews(id uint) error {
	return s.PostRepo.AddViews(id)
}

func (s *PostService) AddLikes(id uint) error {
	return s.PostRepo.AddLikes(id)
}
