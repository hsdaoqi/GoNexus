package v1

import (
	"go-nexus/internal/service"
	"go-nexus/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateGroupReq struct {
	Name   string `json:"name" binding:"required"`
	Avatar string `json:"avatar"`
	Notice string `json:"notice"`
}

// UpdateGroupReq 群组更新请求
type UpdateGroupReq struct {
	ID     uint   `json:"id" binding:"required"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Notice string `json:"notice"`
}

type InviteGroupReq struct {
	GroupID  uint `json:"group_id" binding:"required"`
	FriendID uint `json:"friend_id" binding:"required"`
}

type KickMemberReq struct {
	GroupID  uint `json:"group_id" binding:"required"`
	MemberID uint `json:"member_id" binding:"required"`
}

type MuteMemberReq struct {
	GroupID  uint `json:"group_id" binding:"required"`
	MemberID uint `json:"member_id" binding:"required"`
	Mute     int  `json:"mute" binding:"oneof=0 1"` // 0: unmute, 1: mute
}

// CreateGroup 接口
func CreateGroup(c *gin.Context) {
	var req CreateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}

	userID := c.MustGet("userID").(uint)
	group, err := service.CreateGroup(userID, req.Name, req.Avatar, req.Notice)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, group)
}

// GetMyGroups 接口
func GetMyGroups(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	list, err := service.GetJoinedGroups(userID)
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, list)
}

// UpdateGroup 更新群组信息（仅群主可修改）
func UpdateGroup(c *gin.Context) {
	var req UpdateGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.UpdateGroup(userID, req.ID, req.Name, req.Avatar, req.Notice); err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, gin.H{"id": req.ID})
}

// GetGroupMembers 获取群成员
func GetGroupMembers(c *gin.Context) {
	groupIDStr := c.Query("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	members, err := service.GetGroupMembers(uint(groupID))
	if err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, members)
}

// InviteMember 邀请好友
func InviteMember(c *gin.Context) {
	var req InviteGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.InviteFriendToGroup(req.GroupID, userID, req.FriendID); err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, nil)
}

// KickMember 踢人接口
func KickMember(c *gin.Context) {
	var req KickMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.KickMember(userID, req.GroupID, req.MemberID); err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, "移除成功")
}

// MuteMember 禁言成员
func MuteMember(c *gin.Context) {
	var req MuteMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.MuteMember(userID, req.GroupID, req.MemberID, req.Mute); err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, nil)
}

// SetAdmin 设置管理员
func SetAdmin(c *gin.Context) {
	var req SetAdminReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.SetGroupAdmin(userID, req.GroupID, req.MemberID, req.IsAdmin); err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, nil)
}

// TransferGroup 转让群主
func TransferGroup(c *gin.Context) {
	var req TransferGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParamInvalid)
		return
	}
	userID := c.MustGet("userID").(uint)
	if err := service.TransferGroupOwner(userID, req.GroupID, req.MemberID); err != nil {
		response.FailWithMessage(c, response.ErrSystemError, err.Error())
		return
	}
	response.Success(c, nil)
}
