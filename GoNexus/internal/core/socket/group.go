package socket

import (
	"go-nexus/internal/model/dto"
)

type GroupHandler struct{}

func (h *GroupHandler) Handle(msg *dto.ProtocolMsg) error {
	//// 1. 入库
	//if err := service.SaveMessage(msg); err != nil {
	//	return err
	//}
	//
	//// 2. 查群成员 (ToUserID 就是 GroupID)
	//memberIDs, err := repository.GetGroupMemberIDs(msg.ToUserID)
	//if err != nil {
	//	return err
	//}
	//
	//// 3. 广播 (写扩散)
	//msgBytes := msg.ToBytes()
	//for _, memberID := range memberIDs {
	//	// 不发给自己
	//	if memberID == msg.FromUserID {
	//		continue
	//	}
	//	socket.Manager.SendMessage(memberID, msgBytes)
	//}
	return nil
}
