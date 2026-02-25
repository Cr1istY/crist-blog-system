package service

import (
	"crist-blog/internal/model"
	"crist-blog/internal/repository"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Login(username, password string) (*model.User, error) {
	user, err := s.userRepo.GetByName(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

func (s *UserService) GetCurrentTweetUserByID(userID uuid.UUID) (*model.TweetListUser, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return &model.TweetListUser{
		ID:          user.ID.String(),
		UserName:    user.Username,
		DisplayName: user.Nickname,
		Avatar:      user.Avatar,
		Verified:    user.IsAdmin,
		Email:       user.Email,
		Bio:         user.Bio,
	}, nil

}

func (s *UserService) ChangeUserInfo(id uuid.UUID, user *model.User) error {
	existingUser, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	// 确保用户重要信息不更改
	user.ID = existingUser.ID
	user.Username = existingUser.Username
	user.Email = existingUser.Email
	if err = s.userRepo.ChangeUserInfo(id, user); err != nil {
		return err
	}
	return nil
}
