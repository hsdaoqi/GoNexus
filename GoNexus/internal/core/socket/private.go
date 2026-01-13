package socket

import (
	"go-nexus/internal/model/dto"
	"go-nexus/internal/service"
)

type PrivateHandler struct{}

func (h *PrivateHandler) Handle(msg *dto.ProtocolMsg) error {
	// 1. 调用 Service 入库 (复用之前的逻辑)
	// 注意：service.SaveAndTransform 需要适配一下，只负责存，不负责转 protocol
	// 这里假设 service.SaveMsg(msg) 已经写好了
	if err := service.SaveMessage(msg); err != nil {
		return err
	}

	// 2. 发送给目标用户
	Manager.SendMessage(msg.ToUserID, msg.ToBytes())
	return nil
}
