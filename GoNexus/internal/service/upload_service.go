package service

import (
	"errors"
	"go-nexus/pkg/global"
	"mime/multipart"
)

func UploadFile(file *multipart.FileHeader) (string, error) {
	// 1. 大小限制
	if file.Size > 50*1024*1024 {
		return "", errors.New("文件大小不能超过 50MB")
	}

	// 2. 直接使用全局单例调用！
	// 不需要 if-else 判断，不需要 new 对象
	if global.OSS == nil {
		return "", errors.New("OSS 未初始化")
	}

	return global.OSS.UploadFile(file)
}
