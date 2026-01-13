package socket

import (
	"go-nexus/internal/model/dto"
)

type SystemHandler struct{}

func (h *SystemHandler) Handle(msg *dto.ProtocolMsg) error {
	// 系统通知通常不入库 message 表，或者入一张单独的 notification 表
	// 这里直接广播给所有在线用户
	msgBytes := msg.ToBytes()

	// Manager 需要加一个 Broadcast 方法 (向所有人推送)
	Manager.Broadcast <- msgBytes

	return nil
}
