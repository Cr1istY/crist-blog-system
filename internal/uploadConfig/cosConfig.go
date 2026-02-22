package config

import (
	"log"
	"net/http"
	"net/url"

	"github.com/joho/godotenv"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type CosConfig struct {
	SecretID  string
	SecretKey string
	BucketURL string
	CDNDomain string
}

type COSService struct {
	Client *cos.Client
	Config *CosConfig
}

func loadCosConfig() *CosConfig {
	_ = godotenv.Load()

	return &CosConfig{
		SecretID:  getEnv("COS_SECRET_ID", ""),
		SecretKey: getEnv("COS_SECRET_KEY", ""),
		BucketURL: getEnv("COS_BUCKET_URL", ""),
		CDNDomain: getEnv("COS_CDN_DOMAIN", ""),
	}
}

// NewCOSService 初始化COS客户端
func NewCOSService() *COSService {

	config := loadCosConfig()

	bucketURL, _ := url.Parse(config.BucketURL)
	httpClient := &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	}
	client := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, httpClient)
	log.Println("☁️ COS client init success")
	return &COSService{
		Client: client,
		Config: config}
}
