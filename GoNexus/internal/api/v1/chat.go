package v1

import (
	"github.com/gin-gonic/gin"
	"go-nexus/internal/repository"
	"go-nexus/pkg/response"
)

// GetChatHistoryRequest 请求参数
type GetChatHistoryRequest struct {
	TargetID uint `json:"target_id" form:"target_id" binding:"required"` // 对方ID
	Page     int  `json:"page" form:"page"`                              // 第几页
	PageSize int  `json:"page_size" form:"page_size"`                    // 每页几条
}

// GetChatHistory 接口
func GetChatHistory(c *gin.Context) {
	var req GetChatHistoryRequest
	// GET 请求参数在 URL 里，用 ShouldBindQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}

	// 默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	offset := (req.Page - 1) * req.PageSize

	// 获取当前用户ID
	userID := c.MustGet("userID").(uint)

	// 查询
	messages, err := repository.GetChatHistoryWithUserInfo(userID, req.TargetID, offset, req.PageSize)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}

	response.Success(c, messages)
}
