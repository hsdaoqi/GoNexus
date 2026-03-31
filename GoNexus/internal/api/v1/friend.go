// internal/api/v1/friend.go
package v1

import (
	"encoding/json"
	"errors"
	"go-nexus/internal/core/socket"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/service"
	"go-nexus/pkg/response"

	"github.com/gin-gonic/gin"
)

type SendFriendReq struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	VerifyMsg  string `json:"verify_msg"`
}

func SendFriendRequest(c *gin.Context) {
	var req SendFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.FriendSvc.SendFriendRequest(userID, req.ReceiverID, req.VerifyMsg); err != nil {
		switch {
		case errors.Is(err, service.ErrCannotAddSelf),
			errors.Is(err, service.ErrAlreadyFriend),
			errors.Is(err, service.ErrRequestPending):
			response.FailWithMessage(c, response.ErrBusiness, err.Error())
		case errors.Is(err, service.ErrUserNotFound):
			response.Fail(c, response.ErrUserNotExist)
		default:
			response.Fail(c, response.ErrSystemError)
		}
		return
	}
	// 实时推送通知
	go func() {
		msg := &dto.ProtocolMsg{
			Type:     dto.TypeFriendReq,
			ToUserID: req.ReceiverID,
			Content:  "您有一条新的好友申请",
		}
		socket.Manager.SendMessage(req.ReceiverID, msg.ToBytes())
	}()
	response.Success(c, nil)
}

type HandleFriendReq struct {
	RequestID uint `json:"request_id" binding:"required"`
	Action    int  `json:"action" binding:"required,oneof=1 2"`
}

func HandleFriendRequest(c *gin.Context) {
	var req HandleFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	requesterID, err := service.FriendSvc.HandleFriendRequest(userID, req.RequestID, req.Action)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRequestNotFound):
			response.FailWithMessage(c, response.ErrBusiness, err.Error())
		case errors.Is(err, service.ErrNoPermission):
			response.FailWithMessage(c, response.ErrBusiness, err.Error())
		case errors.Is(err, service.ErrRequestAlreadyDone):
			response.FailWithMessage(c, response.ErrBusiness, err.Error())
		default:
			response.Fail(c, response.ErrSystemError)
		}
		return
	}
	if req.Action == 1 {
		go func() {
			msg := &dto.ProtocolMsg{
				Type:     dto.TypeFriendAns,
				ToUserID: requesterID,
				Content:  "对方同意了您的好友申请",
			}
			msgBytes, _ := json.Marshal(msg)
			socket.Manager.SendMessage(requesterID, msgBytes)
		}()
	}
	response.Success(c, "处理成功")
}

func GetFriendList(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	friends, err := service.FriendSvc.GetFriendList(userID)
	if err != nil {
		response.Fail(c, response.ErrSystemError)
		return
	}
	response.Success(c, friends)
}

type DeleteFriendReq struct {
	FriendID uint `json:"friend_id" binding:"required"`
}

func DeleteFriendRecord(c *gin.Context) {
	var req DeleteFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.FriendSvc.DeleteFriend(userID, req.FriendID); err != nil {
		response.Fail(c, response.ErrSystemError)
		return
	}
	response.Success(c, nil)
}

func GetPendingRequests(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	requests, err := service.FriendSvc.GetPendingRequests(userID)
	if err != nil {
		response.Fail(c, response.ErrSystemError)
		return
	}
	response.Success(c, requests)
}

func RecommendFriends(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	list, err := service.FriendSvc.RecommendFriends(userID)
	if err != nil {
		response.Fail(c, response.ErrSystemError)
		return
	}
	response.Success(c, list)
}
