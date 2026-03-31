// internal/repository/friend_repo.go
package repository

import (
	"errors"
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/pkg/global"

	"gorm.io/gorm"
)

type FriendRepository struct{}

var FriendRepo = &FriendRepository{}

func (r *FriendRepository) IsFriend(userID, friendID uint) bool {
	var count int64
	global.DB.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ? AND deleted_at IS NULL", userID, friendID).
		Count(&count)
	return count > 0
}

func (r *FriendRepository) GetList(userID uint) ([]model.UserProfileResponse, error) {
	var friends []model.UserProfileResponse
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

func (r *FriendRepository) IncrementUnread(userID, friendID uint) error {
	return global.DB.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		UpdateColumn("unread_count", gorm.Expr("unread_count + ?", 1)).Error
}

func (r *FriendRepository) ClearUnread(userID, friendID uint) error {
	return global.DB.Model(&model.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("unread_count", 0).Error
}

func (r *FriendRepository) CreateRequest(req *model.FriendRequest) error {
	return global.DB.Create(req).Error
}

func (r *FriendRepository) GetPendingRequest(requestID uint) (*model.FriendRequest, error) {
	var req model.FriendRequest
	err := global.DB.Where("id = ? AND status = ?", requestID, model.RequestStatusPending).
		First(&req).Error
	return &req, err
}

func (r *FriendRepository) CheckRequestExist(requesterID, receiverID uint) bool {
	var count int64
	global.DB.Model(&model.FriendRequest{}).
		Where("requester_id = ? AND receiver_id = ? AND status = ?",
			requesterID, receiverID, model.RequestStatusPending).
		Count(&count)
	return count > 0
}

func (r *FriendRepository) DeleteRecord(userID, friendID uint) error {
	return global.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).
		Delete(&model.Friend{}).Error
}

func (r *FriendRepository) GetPendingRequests(receiverID uint) ([]dto.FriendRequestResponse, error) {
	var requests []dto.FriendRequestResponse
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

func (r *FriendRepository) GetFriendIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := global.DB.Model(&model.Friend{}).
		Where("user_id = ?", userID).
		Pluck("friend_id", &ids).Error
	return ids, err
}

func (r *FriendRepository) GetSecondDegreeFriends(myFriendIDs []uint, myUserID uint) (map[uint]int, error) {
	if len(myFriendIDs) == 0 {
		return nil, nil
	}
	type Result struct {
		FriendID uint
		Count    int
	}
	var results []Result
	err := global.DB.Model(&model.Friend{}).
		Select("friend_id, count(*) as count").
		Where("user_id IN ?", myFriendIDs).
		Where("friend_id != ?", myUserID).
		Group("friend_id").
		Scan(&results).Error

	candidateMap := make(map[uint]int)
	for _, r := range results {
		candidateMap[r.FriendID] = r.Count
	}
	return candidateMap, err
}

func (r *FriendRepository) GetUsersByIDs(ids []uint) ([]model.User, error) {
	var users []model.User
	if len(ids) == 0 {
		return users, nil
	}
	err := global.DB.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

// AcceptRequest 同意申请的事务：更新申请状态 + 双向建立好友关系
// 之前这个事务逻辑错误地放在 Service 层，现在归还给 Repository
func (r *FriendRepository) AcceptRequest(req *model.FriendRequest) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		req.Status = model.RequestStatusAccepted
		if err := tx.Save(req).Error; err != nil {
			return err
		}
		// 建立 我->他 的关系（支持恢复软删除）
		if err := r.upsertFriend(tx, req.ReceiverID, req.RequesterID); err != nil {
			return err
		}
		// 建立 他->我 的关系
		return r.upsertFriend(tx, req.RequesterID, req.ReceiverID)
	})
}

// upsertFriend 创建或恢复好友关系（内部方法，不对外暴露）
func (r *FriendRepository) upsertFriend(tx *gorm.DB, userID, friendID uint) error {
	var friend model.Friend
	result := tx.Unscoped().
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		First(&friend)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return tx.Create(&model.Friend{
				UserID:   userID,
				FriendID: friendID,
				Source:   model.SourceSearch,
			}).Error
		}
		return result.Error
	}
	// 记录存在（可能是软删除的），恢复它
	return tx.Unscoped().Model(&friend).Update("deleted_at", nil).Error
}
