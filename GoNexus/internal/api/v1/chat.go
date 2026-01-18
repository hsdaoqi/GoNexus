package v1

import (
	"go-nexus/internal/core/socket"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
	"go-nexus/internal/service"
	"go-nexus/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetChatHistoryRequest 请求参数
type GetChatHistoryRequest struct {
	TargetID uint `json:"target_id" form:"target_id" binding:"required"` // 对方ID或群ID
	ChatType int  `json:"chat_type" form:"chat_type"`                    // 1 单聊 2 群聊
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

	// 默认聊天类型：单聊
	if req.ChatType == 0 {
		req.ChatType = dto.ChatTypePrivate
	}

	// 获取当前用户ID
	userID := c.MustGet("userID").(uint)

	var (
		messages interface{}
		err      error
	)

	// 根据聊天类型选择不同的查询逻辑
	if req.ChatType == dto.ChatTypeGroup {
		messages, err = repository.GetGroupChatHistoryWithUserInfo(userID, req.TargetID, offset, req.PageSize)
	} else {
		messages, err = repository.GetChatHistoryWithUserInfo(userID, req.TargetID, offset, req.PageSize)
	}
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}

	response.Success(c, messages)
}

// 定义请求参数
type RevokeMessageRequest struct {
	MsgID    uint `json:"msg_id" binding:"required"`
	ChatType int  `json:"chat_type" binding:"required"` // 1单聊 2群聊
	TargetID uint `json:"target_id" binding:"required"` // 对方ID或群ID
}

// RevokeMessage 撤回消息接口
func RevokeMessage(c *gin.Context) {
	var req RevokeMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}

	userID := c.MustGet("userID").(uint)

	// 1. 调用 Service 更新数据库状态
	if err := service.RevokeMessage(userID, req.MsgID); err != nil {
		response.FailWithMessage(c, response.ErrBusiness, err.Error())
		return
	}

	// 2. 构造 WebSocket 通知消息 (Type=7)
	// Content 存放被撤回的消息ID，方便前端定位
	notifyMsg := dto.ProtocolMsg{
		Type:       dto.TypeRevoke, // 假设 7 为撤回信令
		Content:    strconv.Itoa(int(req.MsgID)),
		ChatType:   req.ChatType,
		FromUserID: userID,
		ToUserID:   req.TargetID,
	}

	// 3. 广播通知
	if req.ChatType == 2 {
		// 群聊：查出所有群成员并推送
		memberIDs, _ := repository.GetGroupMemberIDs(req.TargetID)
		for _, mid := range memberIDs {
			socket.Manager.SendMessage(mid, notifyMsg.ToBytes())
		}
	} else {
		// 单聊：推给对方和自己(多端同步)
		socket.Manager.SendMessage(req.TargetID, notifyMsg.ToBytes())
		socket.Manager.SendMessage(userID, notifyMsg.ToBytes())
	}

	response.Success(c, nil)
}

type ReadMessageRequest struct {
	TargetID uint `json:"target_id" binding:"required"`
	ChatType int  `json:"chat_type" binding:"required"`
}

func ReadMessage(c *gin.Context) {
	var req ReadMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}

	userID := c.MustGet("userID").(uint)
	if err := service.ReadMessage(userID, req.TargetID, req.ChatType); err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}

	response.Success(c, nil)
}
