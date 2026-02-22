package handler

import (
	"context"
	"crist-blog/internal/model"
	config "crist-blog/internal/uploadConfig"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type COSHandler struct {
	uploadHandler *UploadHandler
	cosService    *config.COSService
	config        *config.Config
}

func NewCOSHandler(uploadHandler *UploadHandler, service *config.COSService, cfg *config.Config) *COSHandler {
	return &COSHandler{
		uploadHandler: uploadHandler,
		cosService:    service,
		config:        cfg,
	}
}

// UploadImage 单图上传
func (h *COSHandler) UploadImage(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "未找到上传的文件")
	}

	response, err := h.processFileUpload(file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

// UploadImageBatch 批量上传
func (h *COSHandler) UploadImageBatch(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "无效的表单数据")
	}

	files := form.File["images"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "未找到上传的文件")
	}

	if len(files) > h.config.MaxUploadCount {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("最多只能上传 %d 个文件", h.config.MaxUploadCount))
	}

	responses := make([]model.UploadResponse, 0, len(files))

	for _, file := range files {
		response, err := h.processFileUpload(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string][]model.UploadResponse{
		"images": responses,
	})
}

// processFileUpload 处理单个文件上传的通用逻辑（上传到COS）
func (h *COSHandler) processFileUpload(file *multipart.FileHeader) (model.UploadResponse, error) {
	// 验证文件大小
	if file.Size > h.config.MaxUploadSize {
		return model.UploadResponse{}, fmt.Errorf("文件大小不能超过 %dMB", h.config.MaxUploadSize/1024/1024)
	}

	// 验证文件类型
	if !h.isValidImageType(file.Header.Get("Content-Type")) {
		return model.UploadResponse{}, fmt.Errorf("不允许的文件类型: %s", file.Header.Get("Content-Type"))
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return model.UploadResponse{}, fmt.Errorf("无法打开上传的文件")
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	// 生成COS对象键（按日期分类存储）
	ext := filepath.Ext(file.Filename)
	objectKey := h.generateObjectKey(ext)

	// 上传
	ctx := context.Background()
	_, err = h.cosService.Client.Object.Put(ctx, objectKey, src, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	})
	if err != nil {
		return model.UploadResponse{}, fmt.Errorf("上传文件到COS失败: %v", err)
	}

	// 构建CDN访问URL
	cdnURL := fmt.Sprintf("%s/%s", h.cosService.Config.CDNDomain, objectKey)
	width, height := 0, 0
	return model.UploadResponse{
		URL:       cdnURL,
		ID:        uuid.New().String(),
		Filename:  filepath.Base(objectKey),
		Size:      file.Size,
		Width:     width,
		Height:    height,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// isValidImageType 验证图片类型
func (h *COSHandler) isValidImageType(contentType string) bool {
	for _, allowed := range h.config.AllowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

func (h *COSHandler) generateObjectKey(ext string) string {
	datePath := time.Now().Format("2006/01/02")
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return fmt.Sprintf("images/%s/%s", datePath, filename)
}
