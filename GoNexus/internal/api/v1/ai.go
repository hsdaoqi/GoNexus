package v1

import (
	"context"
	pb "go-nexus/internal/api/proto" // 引入生成的 proto
	"go-nexus/internal/core"         // 引入 AIClient
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
