package service

import (
	"context"
	"crist-blog/internal/model"
	"crist-blog/internal/repository"
)

type TweetService struct {
	TweetRepo *repository.TweetRepository
}

func NewTweetService(tweetRepo *repository.TweetRepository) *TweetService {
	return &TweetService{
		TweetRepo: tweetRepo,
	}
}

func (s *TweetService) CreateTweetWithImages(ctx context.Context, tweet *model.Tweet, imageIDs []string) error {
	return s.TweetRepo.CreateWithImages(ctx, tweet, imageIDs)
}

func (s *TweetService) GetAllWithImages(ctx context.Context, limit, offset int) ([]model.Tweet, error) {
	return s.TweetRepo.GetAllWithImages(ctx, limit, offset)
}

func (s *TweetService) GetTweetWithImagesByID(ctx context.Context, id string) (*model.Tweet, error) {
	return s.TweetRepo.GetByIDWithImages(ctx, id)
}
