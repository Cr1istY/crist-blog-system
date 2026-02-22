package model

type UploadResponse struct {
	URL       string `json:"url"`
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	CreatedAt string `json:"created_at"`
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
