package socket

import (
	"errors"
	"go-nexus/internal/model/dto"
	"go-nexus/internal/repository"
	"go-nexus/internal/service"
)

type GroupHandler struct{}

func (h *GroupHandler) Handle(msg *dto.ProtocolMsg) error {
	// 0. 检查是否是群成员
	if !repository.CheckGroupMember(msg.ToUserID, msg.FromUserID) {
		return errors.New("you are not a member of this group")
	}

	// 1. 检查是否被禁言
	isMuted, err := repository.IsMemberMuted(msg.ToUserID, msg.FromUserID)
	if err == nil && isMuted {
		// 发送禁言提示给发送者
		errMsg := dto.ProtocolMsg{
			Type:     dto.TypeSystem,
			ChatType: dto.ChatTypeGroup,
			Content:  "您已被禁言，无法发言",
			ToUserID: msg.ToUserID, // 必须是群ID，否则前端无法识别是哪个群的消息
		}
		Manager.SendMessage(msg.FromUserID, errMsg.ToBytes())
		return errors.New("member is muted")
	}

	// 2. 入库
	if err := service.SaveMessage(msg); err != nil {
		return err
	}

	// 2. 查群成员 (ToUserID 就是 GroupID)
	memberIDs, err := repository.GetGroupMemberIDs(msg.ToUserID)
	if err != nil {
		return err
	}

	// 3. 广播 (写扩散)
	msgBytes := msg.ToBytes()
	for _, memberID := range memberIDs {
		// 发给所有成员，包括自己；客户端做去重
		Manager.SendMessage(memberID, msgBytes)
	}
	return nil
}
