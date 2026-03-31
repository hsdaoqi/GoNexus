package service

import (
	"errors"
	"go-nexus/internal/core"
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
	"go-nexus/pkg/global"
	"strconv"
	"time"
)

// SaveMessage 处理消息入库并转换回 DTO 以便发送
func SaveMessage(proto *dto.ProtocolMsg) error {
	// 1. DTO -> Model (准备入库)
	dbMsg := &model.Message{
		FromUserID: proto.FromUserID,
		ToUserID:   proto.ToUserID,
		ChatType:   proto.ChatType,
		MsgType:    proto.Type,
		Content:    proto.Content,
		Url:        proto.Url,
	}
	// 2. 落库
	if err := repository.MessageRepo.Save(dbMsg); err != nil {
		return err
	}

	// 2.5 增加未读数
	if proto.ChatType == dto.ChatTypePrivate {
		repository.FriendRepo.IncrementUnread(proto.ToUserID, proto.FromUserID)
	} else if proto.ChatType == dto.ChatTypeGroup {
		repository.GroupRepo.IncrementGroupUnread(proto.ToUserID, proto.FromUserID)
	}

	// 3. 补充 DTO 信息 (准备发送)
	// 填充发送时间
	proto.SendTime = dbMsg.CreatedAt.Format("2006-01-02 15:04:05")
	// 填充 DB 生成的 ID
	proto.MsgID = dbMsg.ID

	// 填充发送者信息 (查询用户信息)
	user, err := repository.UserRepo.GetByID(proto.FromUserID)
	if err == nil {
		proto.SenderNickname = user.Nickname
		proto.SenderAvatar = user.Avatar
	}

	// 4. 🔥 同步给 AI (只有文本消息)
	if proto.Type == dto.TypeText {
		msgID := strconv.Itoa(int(dbMsg.ID)) // 使用数据库ID
		sessionID := core.GetSessionID(proto.ChatType, proto.FromUserID, proto.ToUserID)
		core.AsyncSyncMessage(proto.FromUserID, proto.Content, msgID, proto.SenderNickname, sessionID)
	}

	return nil
}

// ReadMessage 标记消息为已读
func ReadMessage(userID uint, targetID uint, chatType int) error {
	if chatType == dto.ChatTypePrivate {
		return repository.FriendRepo.ClearUnread(userID, targetID)
	} else if chatType == dto.ChatTypeGroup {
		return repository.GroupRepo.ClearUnread(userID, targetID)
	}
	return nil
}

func RevokeMessage(userID uint, msgID uint) error {
	var msg model.Message
	// 1. 查找消息
	if err := global.DB.First(&msg, msgID).Error; err != nil {
		return errors.New("消息不存在")
	}

	// 2. 权限校验
	if msg.FromUserID != userID {
		return errors.New("只能撤回自己的消息")
	}

	// 3. 时间限制 (例如 2 分钟)
	if time.Since(msg.CreatedAt) > 2*time.Minute {
		return errors.New("发送时间超过2分钟，无法撤回")
	}

	// 4. 更新状态
	// 注意：这里只更新 IsRevoked，内容可以不删，或者替换为 "此消息已撤回"
	if err := global.DB.Model(&msg).Update("is_revoked", true).Error; err != nil {
		return err
	}

	// 5. 🔥 通知 AI 撤回
	core.AsyncRevokeMessage(userID, strconv.Itoa(int(msgID)))

	return nil
}
