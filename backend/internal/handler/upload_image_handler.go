package handler

import (
	"crist-blog/internal/model"
	config "crist-blog/internal/uploadConfig"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	config *config.Config
	// 分片上传存储
	chunks map[string]map[int][]byte
	mu     sync.Mutex
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	// 确保上传目录存在
	err := os.MkdirAll(cfg.UploadPath, 0755)
	if err != nil {
		return nil
	}

	return &UploadHandler{
		config: cfg,
		chunks: make(map[string]map[int][]byte),
	}
}

func (h *UploadHandler) saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	_, err = io.Copy(out, src)
	return err
}

// UploadImage 单图上传
func (h *UploadHandler) UploadImage(c echo.Context) error {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "未找到上传的文件")
	}

	response, err := h.processFileUpload(file)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UploadHandler) UploadImages(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "无效的表单数据")
	}

	files := form.File["images"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "未找到上传的文件")
	}

	if len(files) > 4 {
		return echo.NewHTTPError(http.StatusBadRequest, "最多只能上传 4 张图片")
	}

	responses := make([]model.UploadResponse, 0, len(files))

	for _, file := range files {
		response, err := h.processFileUpload(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string][]model.UploadResponse{
		"images": responses,
	})
}

// InitChunkedUpload 初始化分片上传
func (h *UploadHandler) InitChunkedUpload(c echo.Context) error {
	var req model.UploadRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "请求参数错误")
	}

	if !h.isValidImageType(req.FileType) {
		return echo.NewHTTPError(http.StatusBadRequest, "不支持的图片格式")
	}

	uploadID := uuid.New().String()

	h.mu.Lock()
	h.chunks[uploadID] = make(map[int][]byte)
	h.mu.Unlock()

	return c.JSON(http.StatusOK, map[string]string{
		"upload_id": uploadID,
	})
}

// UploadChunk 上传分片
func (h *UploadHandler) UploadChunk(c echo.Context) error {
	uploadID := c.FormValue("uploadId")
	chunkIndex := c.FormValue("chunkIndex")
	totalChunks := c.FormValue("totalChunks")

	if uploadID == "" || chunkIndex == "" || totalChunks == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "缺少必要参数")
	}

	file, err := c.FormFile("chunk")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "未找到分片文件")
	}

	// 读取分片数据
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "读取分片失败")
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	chunkData := make([]byte, file.Size)
	_, err = src.Read(chunkData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "读取分片数据失败")
	}

	// 存储分片
	h.mu.Lock()
	if _, exists := h.chunks[uploadID]; !exists {
		h.mu.Unlock()
		return echo.NewHTTPError(http.StatusBadRequest, "上传会话不存在")
	}

	idx := 0
	_, _ = fmt.Sscanf(chunkIndex, "%d", &idx)
	h.chunks[uploadID][idx] = chunkData
	h.mu.Unlock()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"chunk_index": idx,
		"uploaded":    true,
	})
}

// MergeChunks 合并分片
func (h *UploadHandler) MergeChunks(c echo.Context) error {
	var req model.MergeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "请求参数错误")
	}

	h.mu.Lock()
	chunks, exists := h.chunks[req.UploadID]
	if !exists {
		h.mu.Unlock()
		return echo.NewHTTPError(http.StatusBadRequest, "上传会话不存在")
	}

	// 获取所有分片并排序
	totalChunks := len(chunks)
	mergedData := make([]byte, 0)
	for i := 0; i < totalChunks; i++ {
		chunk, exists := chunks[i]
		if !exists {
			h.mu.Unlock()
			return echo.NewHTTPError(http.StatusBadRequest, "分片不完整")
		}
		mergedData = append(mergedData, chunk...)
	}

	// 清理分片数据
	delete(h.chunks, req.UploadID)
	h.mu.Unlock()

	// 生成文件名并保存
	ext := filepath.Ext(req.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	thefilepath := filepath.Join(h.config.UploadPath, filename)

	if err := os.WriteFile(thefilepath, mergedData, 0644); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "保存文件失败")
	}

	// 获取图片尺寸
	width, height := h.getImageDimensions(thefilepath)

	response := model.UploadResponse{
		URL:       fmt.Sprintf("/uploads/images/%s", filename),
		ID:        uuid.New().String(),
		Filename:  filename,
		Size:      int64(len(mergedData)),
		Width:     width,
		Height:    height,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, response)
}

// isValidImageType 辅助函数：验证图片类型
func (h *UploadHandler) isValidImageType(contentType string) bool {
	for _, allowed := range h.config.AllowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// getImageDimensions 辅助函数：获取图片尺寸
func (h *UploadHandler) getImageDimensions(filepath string) (int, int) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, 0
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	theConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0
	}

	return theConfig.Width, theConfig.Height
}

// processFileUpload 处理单个文件上传的通用逻辑
func (h *UploadHandler) processFileUpload(file *multipart.FileHeader) (model.UploadResponse, error) {
	// 验证文件大小
	if file.Size > h.config.MaxUploadSize {
		return model.UploadResponse{}, fmt.Errorf("文件大小不能超过 %dMB", h.config.MaxUploadSize/1024/1024)
	}

	// 验证文件类型
	if !h.isValidImageType(file.Header.Get("Content-Type")) {
		return model.UploadResponse{}, fmt.Errorf("不支持的图片格式")
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	theFilepath := filepath.Join(h.config.UploadPath, filename)

	// 保存文件
	if err := h.saveUploadedFile(file, theFilepath); err != nil {
		return model.UploadResponse{}, fmt.Errorf("保存文件失败")
	}

	// 获取图片尺寸
	width, height := h.getImageDimensions(theFilepath)

	// 构建响应
	return model.UploadResponse{
		URL:       fmt.Sprintf("/uploads/images/%s", filename),
		ID:        uuid.New().String(),
		Filename:  filename,
		Size:      file.Size,
		Width:     width,
		Height:    height,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}
