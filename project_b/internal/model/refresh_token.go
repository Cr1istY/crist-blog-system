package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents the 'refresh_tokens' table.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TokenHash string    `gorm:"type:text;not null" json:"-"` // 绝不返回！
	UserAgent string    `gorm:"type:text" json:"user_agent,omitempty"`
	IPAddress string    `gorm:"type:inet" json:"ip_address,omitempty"` // GORM 不直接支持 INET，用 string 存 IP
	Province  string    `gorm:"type:text"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null" json:"expires_at"`
	Revoked   bool      `gorm:"type:boolean;not null;default:false" json:"revoked"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`

	// 关联 User（可选）
	User User `gorm:"foreignKey:UserID"`
}

func (RefreshToken) TableName() string {
	return "admin.refresh_tokens"
}
