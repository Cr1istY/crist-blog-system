package model

import (
	"time"
)

type Image struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	URL       string    `gorm:"type:varchar(512);not null" json:"url"`
	Filename  string    `gorm:"type:varchar(255);not null" json:"filename"`
	Size      int64     `gorm:"type:bigint;not null" json:"size"`
	Width     int       `gorm:"type:integer" json:"width,omitempty"`
	Height    int       `gorm:"type:integer" json:"height,omitempty"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}

func (Image) TableName() string {
	return "blog.images"
}

type UploadResponse struct {
	URL       string `json:"url"`
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	CreatedAt string `json:"created_at"`
}

func (r *UploadResponse) ToImage() *Image {
	return &Image{
		ID:       r.ID,
		URL:      r.URL,
		Filename: r.Filename,
		Size:     r.Size,
		Width:    r.Width,
		Height:   r.Height,
	}
}

type UploadRequest struct {
	Filename string `json:"filename"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
}

type ChunkUploadRequest struct {
	UploadID    string `form:"uploadId"`
	ChunkIndex  int    `form:"chunkIndex"`
	TotalChunks int    `form:"totalChunks"`
}

type MergeRequest struct {
	UploadID string `json:"uploadId"`
	Filename string `json:"filename"`
}

// CosUploadResponse cos上传响应
type CosUploadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		URL      string `json:"url"`       // CDN访问URL
		COSKey   string `json:"cos_key"`   // COS中的对象键
		FileName string `json:"file_name"` // 原始文件名
	} `json:"data"`
}
