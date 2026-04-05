package repository

import (
	"crist-blog/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		DB: db,
	}
}

func (r *RefreshTokenRepository) CreateRefreshToken(token *model.RefreshToken) error {
	return r.DB.Create(token).Error
}

// FindByTokenHash 弃用，由于hash的特性，不能返回也不应该直接返回hash值
func (r *RefreshTokenRepository) FindByTokenHash(hash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	if err := r.DB.Where("token_hash = ?", hash).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) ReturnAdminHash() (*model.RefreshToken, error) {
	userID := "00000000-0000-0000-0000-000000000001"
	var token model.RefreshToken
	if err := r.DB.Where("user_id = ? AND revoked = false", userID).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) ReturnAdminHashWithIPAndAgent(userID, userAgent, ip string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	if err := r.DB.Where("user_id = ? AND user_agent = ? AND ip_address = ? AND revoked = false", userID, userAgent, ip).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) ReturnAdminHashWithProvinceAndAgent(userID, userAgent, province string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	if err := r.DB.Where("user_id = ? AND user_agent = ? AND province = ? AND revoked = false", userID, userAgent, province).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

// Revoke revoke refresh token
func (r *RefreshTokenRepository) Revoke(id uuid.UUID) error {
	return r.DB.Model(&model.RefreshToken{}).
		Where("id = ? AND revoked = false", id).
		Update("revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeAllByUserID(userID uuid.UUID) error {
	return r.DB.Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func (r *RefreshTokenRepository) RevokeAllByUserIDAndAgent(userID uuid.UUID, agent string) error {
	return r.DB.Model(&model.RefreshToken{}).
		Where("user_id = ? AND user_agent = ? AND revoked = false", userID, agent).
		Update("revoked", true).Error
}

func (r *RefreshTokenRepository) CleanExpiredTokens() error {
	return r.DB.Where("expires_at < ? OR revoked = true", time.Now()).Delete(&model.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) FindAllValid() ([]*model.RefreshToken, error) {
	var tokens []*model.RefreshToken
	err := r.DB.Where("revoked = false AND expires_at > ?", time.Now()).Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *RefreshTokenRepository) FindIpByUserIDAndUserAgent(userID uuid.UUID, agent string) (string, error) {
	var token model.RefreshToken
	err := r.DB.Where("user_id = ? AND user_agent = ? AND revoked = false", userID, agent).
		Select("ip_address").
		First(token).Error
	if err != nil {
		return "", err
	}
	return token.IPAddress, nil
}
