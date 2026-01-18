package repository

import (
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/pkg/global"
)

// SaveMessage 保存信息
func SaveMessage(msg *model.Message) error {
	return global.DB.Create(msg).Error
}

// GetChatHistory 获取单聊历史记录
// 参数：uid(我), targetID(对方), offset(偏移量), limit(条数)
func GetChatHistory(uid, targetID uint, offset, limit int) ([]model.Message, error) {
	var message []model.Message

	// SQL: SELECT * FROM messages
	// WHERE (from_user_id = A AND to_user_id = B)
	//    OR (from_user_id = B AND to_user_id = A)
	// ORDER BY created_at DESC (按时间倒序，最新的在前面)
	err := global.DB.Where("chat_type = ? AND ((from_user_id = ? AND to_user_id = ?) or (from_user_id = ? AND to_user_id = ?))",
		model.ChatTypePrivate, uid, targetID, targetID, uid).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&message).Error
	return message, err
}
func GetChatHistoryWithUserInfo(uid, targetID uint, offset, limit int) ([]dto.ChatHistoryResponse, error) {
	var messages []dto.ChatHistoryResponse

	// 联表查询：messages 表 JOIN users 表
	err := global.DB.Table("messages").
		Select(`messages.*, users.avatar as sender_avatar, users.nickname as sender_nickname`).
		Joins("LEFT JOIN users ON messages.from_user_id = users.id").
		Where("messages.chat_type = ? AND ((messages.from_user_id = ? AND messages.to_user_id = ?) OR (messages.from_user_id = ? AND messages.to_user_id = ?))",
			model.ChatTypePrivate, uid, targetID, targetID, uid).
		Order("messages.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

// GetGroupChatHistoryWithUserInfo 获取群聊历史记录（带用户信息）
func GetGroupChatHistoryWithUserInfo(uid, groupID uint, offset, limit int) ([]dto.ChatHistoryResponse, error) {
	var messages []dto.ChatHistoryResponse

	err := global.DB.Table("messages").
		Select(`messages.*, users.avatar as sender_avatar, users.nickname as sender_nickname`).
		Joins("LEFT JOIN users ON messages.from_user_id = users.id").
		Where("messages.chat_type = ? AND messages.to_user_id = ?", model.ChatTypeGroup, groupID).
		Order("messages.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error

	return messages, err
}
