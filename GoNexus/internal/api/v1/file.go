package v1

import (
	"github.com/gin-gonic/gin"
	"go-nexus/internal/service"
	"go-nexus/pkg/response"
)

// Upload 上传接口
func Upload(c *gin.Context) {
	// 1. 获取文件 (前端字段名叫 "file")
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, "请上传文件")
		return
	}

	// 2. 调用服务
	url, err := service.UploadFile(file)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "上传失败: "+err.Error())
		return
	}

	// 3. 返回 URL
	response.Success(c, gin.H{"url": url})
}
