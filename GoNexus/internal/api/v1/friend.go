package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-nexus/internal/core/socket"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/service"
	"go-nexus/pkg/response"
	"net/http"
)

// SendFriendReq 请求参数
type SendFriendReq struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	VerifyMsg  string `json:"verify_msg"`
}

// HandleFriendReq 处理参数
type HandleFriendReq struct {
	RequestID uint `json:"request_id" binding:"required"`       // 申请记录的ID
	Action    int  `json:"action" binding:"required,oneof=1 2"` // 只能传1或2
}

// SendFriendRequest 接口：发送申请
func SendFriendRequest(c *gin.Context) {
	var req SendFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}

	userID := c.MustGet("userID").(uint)
	if err := service.SendFriendRequest(userID, req.ReceiverID, req.VerifyMsg); err != nil {
		response.FailWithMessage(c, response.ErrParamInvalid, err.Error())
		return
	}

	// 实时推送通知
	go func() {
		// 构造一条系统通知消息
		notifyMsg := &dto.ProtocolMsg{
			Type:     dto.TypeFriendReq, // 类型 4
			ToUserID: req.ReceiverID,    // 发给对方
			Content:  "您有一条新的好友申请",      // 内容随便写，前端主要看 Type
			SendTime: "now",
		}
		// 序列化并发送
		msgBytes, _ := json.Marshal(notifyMsg)
		socket.Manager.SendMessage(req.ReceiverID, msgBytes)
	}()

	response.Success(c, nil)
}

// HandleFriendRequest 接口：处理申请
func HandleFriendRequest(c *gin.Context) {
	var req HandleFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}

	userID := c.MustGet("userID").(uint)
	// 接收 requesterID
	requesterID, err := service.HandleFriendRequest(userID, req.RequestID, req.Action)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	// 如果是“同意”，给申请人(A)发通知
	if req.Action == 1 { // 1=同意
		go func() {
			notifyMsg := &dto.ProtocolMsg{
				Type:     dto.TypeFriendAns, // Type 5
				ToUserID: requesterID,       // 发给 A
				Content:  "对方同意了您的好友申请",
				SendTime: "now",
			}
			msgBytes, _ := json.Marshal(notifyMsg)
			socket.Manager.SendMessage(requesterID, msgBytes)
		}()
	}

	response.Success(c, "处理成功")
}

func GetFriendList(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		response.Fail(c, response.ErrAuthFailed)
		return
	}
	userID := userIDValue.(uint)

	// 2. 调用Service层，使用WebSocket管理器检查在线状态
	friends, err := service.GetFriendListWithOnlineStatus(userID, &socket.Manager)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.ErrSystemError, "获取好友列表失败")
		return
	}

	// 3. 返回成功响应
	response.Success(c, friends)
}

// DeleteFriendReq 删除好友请求参数
type DeleteFriendReq struct {
	FriendID uint `json:"friend_id" binding:"required"` // 要删除的好友ID
}

// DeleteFriendRecord 删除好友记录
func DeleteFriendRecord(c *gin.Context) {
	var req DeleteFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.DeleteFriendRecord(userID, req.FriendID); err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, "删除成功")
}

// GetPendingRequests 获取待处理的好友请求
func GetPendingRequests(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	// 调用Service层获取待处理请求
	requests, err := service.GetPendingRequests(userID)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}

	response.Success(c, requests)
}
