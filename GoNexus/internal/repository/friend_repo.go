package repository

import (
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/pkg/global"

	"gorm.io/gorm"
)

//--- 好友关系操作 ---

// IsFriend 检查两人是否是好友
func IsFriend(userID, friendID uint) bool {
	var count int64
	global.DB.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ? AND deleted_at IS NULL", userID, friendID).
		Count(&count)
	return count > 0
}

// GetFriendList 连表查询我的好友
// 返回值：一个包含 User 信息的切片
func GetFriendList(userID uint) ([]model.UserProfileResponse, error) {
	var friends []model.UserProfileResponse

	//逻辑：查 users 表，条件是 users.id 必须在 friends 表里出现，且对应的 user_id 是我
	//添加条件过滤软删除的好友记录
	// 增加子查询：获取最后一条消息
	lastMsgQuery := `(SELECT CASE 
		WHEN msg_type = 1 THEN content 
		WHEN msg_type = 2 THEN '[图片]' 
		WHEN msg_type = 3 THEN '[语音]' 
		ELSE '[文件]' END 
		FROM messages 
		WHERE ((from_user_id = users.id AND to_user_id = ? AND chat_type = 1) 
		OR (from_user_id = ? AND to_user_id = users.id AND chat_type = 1))
		AND deleted_at IS NULL 
		ORDER BY created_at DESC LIMIT 1) as last_msg`

	err := global.DB.Table("users").
		Select("users.*, friends.unread_count, "+lastMsgQuery, userID, userID).
		Joins("JOIN friends on friends.friend_id = users.id").
		Where("friends.user_id = ? AND friends.deleted_at IS NULL", userID).
		Find(&friends).Error

	return friends, err
}

// IncrementFriendUnread 增加未读数
func IncrementFriendUnread(userID, friendID uint) error {
	return global.DB.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		UpdateColumn("unread_count", gorm.Expr("unread_count + ?", 1)).Error
}

// ClearFriendUnread 清除未读数
func ClearFriendUnread(userID, friendID uint) error {
	return global.DB.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("unread_count", 0).Error
}

// --- 申请记录操作 ---
// CreateFriendRequest 创建申请
func CreateFriendRequest(req *model.FriendRequest) error {
	return global.DB.Create(req).Error
}

// GetPendingRequest 查找一条具体的待处理申请
func GetPendingRequest(requestID uint) (*model.FriendRequest, error) {
	var req model.FriendRequest
	err := global.DB.Where("id = ? AND status = ?", requestID, model.RequestStatusPending).First(&req).Error
	return &req, err
}

// CheckRequestExist 检查是否重复申请 (防止那个人一直点申请按钮骚扰)
func CheckRequestExist(requesterID, receiverID uint) bool {
	var count int64
	global.DB.Model(&model.FriendRequest{}).
		Where("requester_id = ? AND receiver_id = ? AND status = ?", requesterID, receiverID, model.RequestStatusPending).
		Count(&count)
	return count > 0
}

// DeleteFriendRecord 删除好友记录
func DeleteFriendRecord(userID, FriendID uint) error {
	return global.DB.Where("user_id = ? AND friend_id = ?", userID, FriendID).Delete(&model.Friend{}).Error
}

func GetPendingRequests(receiverID uint) ([]dto.FriendRequestResponse, error) {
	var requests []dto.FriendRequestResponse

	// 联表查询：friend_requests JOIN users (申请者信息)
	err := global.DB.Table("friend_requests").
		Select(`friend_requests.id, friend_requests.requester_id, friend_requests.verify_msg,
		        friend_requests.created_at, users.nickname as requester_name, users.avatar as requester_avatar`).
		Joins("LEFT JOIN users ON friend_requests.requester_id = users.id").
		Where("friend_requests.receiver_id = ? AND friend_requests.status = ?",
			receiverID, model.RequestStatusPending).
		Order("friend_requests.created_at DESC").
		Find(&requests).Error

	return requests, err
}
