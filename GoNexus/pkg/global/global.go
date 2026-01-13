package global

import (
	"go-nexus/configs"
	"gorm.io/gorm"
	"mime/multipart"
)

// OSS 接口定义
type OSSClient interface {
	UploadFile(file *multipart.FileHeader) (string, error)
	DeleteFile(key string) error
}

var (
	DB     *gorm.DB
	Config *configs.Config
	OSS    OSSClient
)
