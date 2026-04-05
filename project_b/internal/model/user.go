package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the 'users' table in the database.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username     string    `gorm:"type:text;not null;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`
	Nickname     string    `gorm:"type:text" json:"nickname"`
	Email        string    `gorm:"type:text;not null;uniqueIndex" json:"email"`
	Avatar       string    `gorm:"type:text" json:"avatar"`
	Bio          string    `gorm:"type:text" json:"bio"`
	IsAdmin      bool      `gorm:"type:boolean;not null;default:false" json:"is_admin"`

	// GORM 自动管理时间戳（需嵌入 gorm.Model 或手动声明）
	CreatedAt time.Time      `gorm:"type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp with time zone" json:"-"` // 软删除（可选）
}

func (User) TableName() string {
	return "admin.users"
}
