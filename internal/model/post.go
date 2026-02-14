package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PostStatus string

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
	Private   PostStatus = "private"
)

var ErrSlugNotFound = errors.New("post not found by slug")

type Post struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Title           string         `gorm:"type:text;not null" json:"title"`
	Slug            string         `gorm:"type:text;not null;uniqueIndex" json:"slug"`
	Content         string         `gorm:"type:text" json:"content"`
	Excerpt         string         `gorm:"type:text" json:"excerpt"`
	Status          PostStatus     `gorm:"type:post_status_enum;not null" json:"status"`
	CategoryID      uuid.UUID      `gorm:"type:uuid;not null" json:"category_id"`
	Tags            pq.StringArray `gorm:"type:text[]" json:"tags"`
	Views           int            `gorm:"default:0" json:"views"`
	Likes           int            `gorm:"default:0" json:"likes"`
	Thumbnail       string         `gorm:"type:text" json:"thumbnail"`
	PublishedAt     *time.Time     `json:"published_at"`
	MetaTitle       string         `gorm:"type:text" json:"meta_title"`
	MetaDescription string         `gorm:"type:text" json:"meta_description"`
	SearchVector    interface{}    `gorm:"type:tsvector" json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	IsPinned        bool           `gorm:"default:false;not null" json:"is_pinned"`        // 置顶开关
	PinnedOrder     int            `gorm:"default:0;not null" json:"pinned_order"`         // 置顶内排序
	PinnedUntil     *time.Time     `gorm:"type:timestamptz" json:"pinned_until,omitempty"` // 过期时间（可空）
}

// CreatePostRequest 创建文章请求结构体
type CreatePostRequest struct {
	UserID          string     `json:"user_id" validate:"required,uuid4"`
	Title           string     `json:"title" validate:"required"`
	Slug            string     `json:"slug" validate:"required"`
	Content         string     `json:"content"`
	Excerpt         string     `json:"excerpt"`
	Status          string     `json:"status" validate:"oneof=draft published private"`
	CategoryID      string     `json:"category_id" validate:"required,uuid4"`
	Tags            []string   `json:"tags"`
	MetaTitle       string     `json:"meta_title"`
	PublishedAt     *time.Time `json:"published_at"`
	MetaDescription string     `json:"meta_description"`
	Thumbnail       string     `json:"thumbnail"`
}

type PostFrontend struct {
	ID        uint     `json:"id"`
	Slug      string   `json:"slug" validate:"required"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Date      string   `json:"date"`
	Excerpt   string   `json:"excerpt"`
	Views     int      `gorm:"default:0" json:"views"`
	Likes     int      `gorm:"default:0" json:"likes"`
	Thumbnail string   `json:"thumbnail,omitempty"`
}

type PostFrontendWithPinned struct {
	ID          uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug        string   `json:"slug" validate:"required"`
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
	Date        string   `json:"date"`
	Excerpt     string   `json:"excerpt"`
	Views       int      `gorm:"default:0" json:"views"`
	Likes       int      `gorm:"default:0" json:"likes"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
	IsPinned    bool     `gorm:"default:false;not null" json:"is_pinned"` // 置顶开关
	PinnedOrder int      `gorm:"default:0;not null" json:"pinned_order"`  // 置顶内排序
}

type HotPost struct {
	Slug       string    `gorm:"type:text;not null;uniqueIndex" json:"slug"`
	Title      string    `gorm:"type:text;not null" json:"title"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	Excerpt    string    `gorm:"type:text" json:"excerpt"`
}

type HotPostFrontend struct {
	Slug     string `gorm:"type:text;not null;uniqueIndex" json:"slug"`
	Title    string `gorm:"type:text;not null" json:"title"`
	Category string `json:"category"`
	Date     string `json:"date"`
	Excerpt  string `gorm:"type:text" json:"excerpt"`
}

type LatestPost struct {
	Slug       string    `gorm:"type:text;not null;uniqueIndex" json:"slug"`
	Title      string    `gorm:"type:text;not null" json:"title"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type LatestPostFrontend struct {
	Slug     string `gorm:"type:text;not null;uniqueIndex" json:"slug"`
	Title    string `gorm:"type:text;not null" json:"title"`
	Date     string `json:"date"`
	Category string `json:"category"`
}

func (Post) TableName() string {
	return "blog.posts"
}

// PostDetail 是博客详情页返回给前端的数据结构
type PostDetail struct {
	ID              uint       `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"` // Markdown 原文
	Date            string     `json:"date"`    // 格式化后的发布日期，如 "2025年12月15日"
	Tags            []string   `json:"tags"`
	CategoryID      uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	Category        string     `json:"category"` // 分类名称，非 ID
	Views           int        `json:"views"`
	Likes           int        `json:"likes"`
	Excerpt         string     `json:"excerpt,omitempty"`
	MetaTitle       string     `json:"meta_title,omitempty"`
	MetaDescription string     `json:"meta_description,omitempty"`
	Status          PostStatus `gorm:"type:post_status_enum;not null" json:"status"`
	Thumbnail       string     `json:"thumbnail"`
	IsPinned        bool       `gorm:"default:false;not null" json:"is_pinned"` // 置顶开关
}
