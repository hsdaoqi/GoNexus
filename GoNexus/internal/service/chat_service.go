package service

import (
	"go-nexus/internal/model"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
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
	if err := repository.SaveMessage(dbMsg); err != nil {
		return err
	}
	// 3. 补充 DTO 信息 (准备发送)
	// 填充发送时间
	proto.SendTime = dbMsg.CreatedAt.Format("2006-01-02 15:04:05")

	// 填充发送者信息 (查询用户信息)
	user, err := repository.GetUserByID(proto.FromUserID)
	if err == nil {
		proto.SenderNickname = user.Nickname
		proto.SenderAvatar = user.Avatar
	}
	return nil
}
