package core

import (
	"context"
	"fmt"
	"go-nexus/internal/model/dto"
	"log"
	"time"

	// 👇 这里必须对应你 go.mod 里的模块名 + 路径
	pb "go-nexus/internal/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var AIClient pb.AIServiceClient

// InitAIClient 连接 Python 服务
func InitAIClient() {
	// 连接本地 Python 服务端口
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ 无法连接 AI 服务: %v", err)
	}

	AIClient = pb.NewAIServiceClient(conn)
	log.Println("✅ AI 服务连接成功 (gRPC)")
}

// GetSessionID 辅助工具，生成会话ID
func GetSessionID(chatType int, fromID, toID uint) string {
	if chatType == dto.ChatTypeGroup { //群聊
		return fmt.Sprintf("g_%d", toID) //toID就是groupID
	}
	if fromID < toID {
		return fmt.Sprintf("p_%d_%d", fromID, toID)
	}
	return fmt.Sprintf("p_%d_%d", toID, fromID)
}

// AsyncSyncMessage 异步把消息推给 AI (不阻塞主流程)
func AsyncSyncMessage(userID uint, content string, msgID string, nickname string, sessionID string) {
	go func() {
		// 设置 5 秒超时
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := AIClient.SyncMessage(ctx, &pb.SyncRequest{
			UserId:    uint32(userID),
			Content:   content,
			MsgId:     msgID,
			Nickname:  nickname,
			SessionId: sessionID,
		})
		if err != nil {
			log.Printf("⚠️ AI 同步失败: %v", err)
		}
	}()
}

// ChatSummary 同步调用 AI 生成摘要
func ChatSummary(chats []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 摘要生成可能较慢
	defer cancel()

	resp, err := AIClient.ChatSummary(ctx, &pb.SummaryRequest{
		Chats: chats,
	})
	if err != nil {
		return "", err
	}
	return resp.Summary, nil
}

// GetReplySuggestions 获取回复建议
func GetReplySuggestions(chats []string, myName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := AIClient.SuggestReply(ctx, &pb.SuggestRequest{
		RecentMessages: chats,
		MyName:         myName,
	})
	if err != nil {
		return nil, err
	}
	return resp.Suggestions, nil
}

// AsyncRevokeMessage 异步撤回消息
func AsyncRevokeMessage(userID uint, msgID string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := AIClient.RevokeMessage(ctx, &pb.RevokeRequest{
			UserId: uint32(userID),
			MsgId:  msgID,
		})
		if err != nil {
			log.Printf("⚠️ AI 撤回失败: %v", err)
		}
	}()
}
