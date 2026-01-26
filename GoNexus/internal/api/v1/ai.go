package v1

import (
	"context"
	"fmt"
	pb "go-nexus/internal/api/proto" // 引入生成的 proto
	"go-nexus/internal/core"         // 引入 AIClient
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
	"go-nexus/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// AISearchReq 搜索请求参数
type AISearchReq struct {
	Query    string `json:"query" form:"query" binding:"required"`         // 用户问什么
	TargetID uint   `json:"target_id" form:"target_id" binding:"required"` //正在跟谁聊
	ChatType int    `json:"chat_type" form:"chat_type" binding:"required"` //群聊 or 单聊
}

// SemanticSearch 语义搜索/问答接口
func SemanticSearch(c *gin.Context) {
	var req AISearchReq
	// 支持 GET Query 或 POST JSON
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	//计算当前窗口的SessionID
	sessionID := core.GetSessionID(req.ChatType, userID, req.TargetID)
	// 调用 gRPC (设置 10秒超时，因为 AI 思考比较慢)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 核心：Go 调 Python
	resp, err := core.AIClient.SemanticSearch(ctx, &pb.SearchRequest{
		Query:     req.Query,
		Limit:     3,         // 搜前3条相关记录
		SessionId: sessionID, // 🔥 告诉 AI：只在这个房间里搜！
	})

	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "AI 服务暂时不可用: "+err.Error())
		return
	}

	// 直接返回 AI 的回答
	response.Success(c, gin.H{
		"answer": resp.Answer,
	})
}

// SummaryReq 总结请求参数
type SummaryReq struct {
	TargetID uint `json:"target_id" form:"target_id" binding:"required"`
	ChatType int  `json:"chat_type" form:"chat_type" binding:"required"`
}

// ChatSummary 聊天总结接口
func ChatSummary(c *gin.Context) {
	var req SummaryReq
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)

	// 1. 获取最近 50 条消息
	limit := 50
	var messages []dto.ChatHistoryResponse
	var err error

	if req.ChatType == model.ChatTypeGroup {
		messages, err = repository.GetGroupChatHistoryWithUserInfo(userID, req.TargetID, 0, limit)
	} else {
		messages, err = repository.GetChatHistoryWithUserInfo(userID, req.TargetID, 0, limit)
	}

	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "获取聊天记录失败")
		return
	}

	if len(messages) == 0 {
		response.Success(c, gin.H{"summary": "暂无聊天记录"})
		return
	}

	// 2. 格式化消息 (注意：数据库返回的是倒序，我们需要正序给 AI)
	var chats []string
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		// 简单的格式：昵称: 内容
		// 过滤非文本消息 (MsgType=1 是文本)
		if msg.MsgType == model.MsgTypeText {
			chats = append(chats, fmt.Sprintf("%s: %s", msg.SenderNickname, msg.Content))
		}
	}

	if len(chats) == 0 {
		response.Success(c, gin.H{"summary": "没有足够的文本消息进行总结"})
		return
	}

	// 3. 调用 AI
	summary, err := core.ChatSummary(chats)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "AI 生成总结失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"summary": summary,
	})
}

// SuggestReq 建议请求参数
type SuggestReq struct {
	TargetID uint `json:"target_id" form:"target_id" binding:"required"`
	ChatType int  `json:"chat_type" form:"chat_type" binding:"required"`
}

// ReplySuggestion 回复建议接口
func ReplySuggestion(c *gin.Context) {
	var req SuggestReq
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)

	// 获取我的信息 (为了告诉 AI 我是谁)
	me, err := repository.GetUserByID(userID)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "获取用户信息失败")
		return
	}

	// 1. 获取最近 10 条消息 (建议只需要少量上下文)
	limit := 10
	var messages []dto.ChatHistoryResponse

	if req.ChatType == model.ChatTypeGroup {
		messages, err = repository.GetGroupChatHistoryWithUserInfo(userID, req.TargetID, 0, limit)
	} else {
		messages, err = repository.GetChatHistoryWithUserInfo(userID, req.TargetID, 0, limit)
	}

	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "获取聊天记录失败")
		return
	}

	// 如果没有聊天记录，无法生成建议
	if len(messages) == 0 {
		response.Success(c, gin.H{"suggestions": []string{}})
		return
	}

	// 2. 格式化消息
	var chats []string
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.MsgType == model.MsgTypeText {
			chats = append(chats, fmt.Sprintf("%s: %s", msg.SenderNickname, msg.Content))
		}
	}

	if len(chats) == 0 {
		response.Success(c, gin.H{"suggestions": []string{}})
		return
	}

	// 3. 调用 AI
	suggestions, err := core.GetReplySuggestions(chats, me.Nickname)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, "AI 生成建议失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"suggestions": suggestions,
	})
}
