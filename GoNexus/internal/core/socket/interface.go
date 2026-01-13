package socket

import (
	"go-nexus/internal/model/dto"
)

// MsgHandler 消息处理器接口
// 所有类型的消息（单聊、群聊、通知）都必须实现这个接口
type MsgHandler interface {
	Handle(msg *dto.ProtocolMsg) error
}
