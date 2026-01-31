package service

import (
	"crist-blog/internal/model"
	"crist-blog/internal/repository"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenExpire  = 15 * time.Minute
	RefreshTokenExpire = 7 * 24 * time.Hour
	refreshTokenLength = 64
)

type AuthService struct {
	userRepo           *repository.UserRepository
	refreshTokenRepo   *repository.RefreshTokenRepository
	jwtSecret          string
	refreshTokenLength int
}

func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		refreshTokenRepo:   refreshTokenRepo,
		jwtSecret:          jwtSecret,
		refreshTokenLength: refreshTokenLength}
}

// generateAccessToken 生成用于用户认证的访问令牌
// 参数:
//   - userID: 用户唯一标识符
//
// 返回值:
//   - string: 生成的JWT访问令牌
//   - error: 生成过程中可能出现的错误
func (s *AuthService) generateAccessToken(userID uuid.UUID) (string, error) {
	// 创建一个新的JWT令牌，使用HS256签名方法
	// 并设置自定义声明：用户ID和过期时间
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),                          // 将用户ID转换为字符串并存储在声明中
		"exp":     time.Now().Add(AccessTokenExpire).Unix(), // 设置令牌过期时间
	})
	// 使用JWT密钥对令牌进行签名，并返回签名字符串
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRandomToken 是一个方法，属于 AuthService 结构体
// 用于生成随机刷新令牌
// 返回一个字符串类型的令牌和可能的错误
func (s *AuthService) generateRandomToken() (string, error) {
	// 创建一个字节切片，长度为 refreshTokenLength 的一半
	bytes := make([]byte, s.refreshTokenLength/2)
	// 使用加密安全的随机数生成器填充字节切片
	// 如果读取失败，返回错误
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// 将字节切片编码为十六进制字符串并返回
	return hex.EncodeToString(bytes), nil
}

// GenerateTokens 为用户生成访问令牌和刷新令牌
// 参数:
//   - user: 用户模型指针，包含用户信息
//   - userAgent: 用户代理字符串，表示客户端类型
//   - ip: 用户IP地址
//
// 返回值:
//   - accessToken: 访问令牌字符串
//   - refreshToken: 刷新令牌字符串
//   - err: 错误信息，如果生成过程中出现错误
func (s *AuthService) GenerateTokens(user *model.User, userAgent, ip string) (accessToken, refreshToken string, err error) {
	// 生成访问令牌，使用用户ID作为参数，仅进行校验
	accessToken, err = s.generateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}
	// 生成随机刷新令牌，长期存储
	refreshToken, err = s.generateRandomToken()
	if err != nil {
		return "", "", err
	}

	// 使用bcrypt算法对刷新令牌进行哈希处理，增强安全性
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	// 撤销该用户之前的所有刷新令牌
	_ = s.refreshTokenRepo.RevokeAllByUserID(user.ID)

	// 创建新的刷新令牌记录
	rt := &model.RefreshToken{
		UserID:    user.ID,                            // 用户ID
		TokenHash: string(tokenHash),                  // 哈希后的令牌
		UserAgent: userAgent,                          // 用户代理
		IPAddress: ip,                                 // IP地址
		ExpiresAt: time.Now().Add(RefreshTokenExpire), // 过期时间
		Revoked:   false,                              // 未被撤销
	}

	// 将新的刷新令牌保存到数据库
	err = s.refreshTokenRepo.CreateRefreshToken(rt)
	if err != nil {
		return "", "", err
	}

	// 返回生成的访问令牌和刷新令牌
	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshAccessToken(refreshTokenStr string) (newAccessToken string, err error) {
	var rt *model.RefreshToken
	// rt, err = s.refreshTokenRepo.FindByTokenHash(refreshTokenStr)
	rt, err = s.refreshTokenRepo.ReturnAdminHash()
	if err != nil {
		return "", errors.New("!!!database error")
	}
	err = bcrypt.CompareHashAndPassword([]byte(rt.TokenHash), []byte(refreshTokenStr))
	if err != nil {
		return "", errors.New("invalid refresh token")
	}
	newAccessToken, err = s.generateAccessToken(rt.UserID)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}
	if err := s.refreshTokenRepo.Revoke(rt.ID); err != nil {
		log.Printf("warning: failed to revoke refresh token %s: %v", rt.ID, err)
	}

	return newAccessToken, nil
}

func (s *AuthService) GetTheRefreshTokenExpired() time.Duration {
	return RefreshTokenExpire
}

func (s *AuthService) JwtSecret() string {
	return s.jwtSecret
}
