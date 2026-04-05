package model

import (
	"time"
)

type Tweet struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Images      []Image   `gorm:"many2many:blog.tweet_images;" json:"images,omitempty"`
	Likes       int       `gorm:"type:integer;default:0" json:"likes"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now();index" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedFlag bool      `gorm:"type:bool; not null; default:false"`
}

type CreateTweetRequest struct {
	Content  string   `json:"content" validate:"required"`
	ImageIDs []string `json:"image_ids,omitempty"`
}

func (Tweet) TableName() string {
	return "blog.tweets"
}

type TweetResponse struct {
	ID        string        `json:"id"`
	User      TweetListUser `json:"user"`
	Content   string        `json:"content"`
	Timestamp time.Time     `json:"timestamp"`
	Likes     int           `json:"likes"`
	Images    []string      `json:"images,omitempty"`
}

type TweetListUser struct {
	ID          string `json:"id"`
	UserName    string `json:"username"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
	Verified    bool   `json:"verified"`
	Bio         string `json:"bio"`
	Email       string `json:"email"`
}
