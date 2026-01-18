package api

import (
	v1 "go-nexus/internal/api/v1"
	"go-nexus/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175"} // 只允许你的 Vue 前端
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))
	//1.公开路由
	publicGroup := r.Group("/api/v1")
	{
		// 用户相关路由
		userGroup := publicGroup.Group("/user")
		{
			userGroup.POST("/register", v1.Register)
			userGroup.POST("/login", v1.Login)
		}
	}
	//2.私有路由
	privateGroup := r.Group("/api/v1")
	privateGroup.Use(middleware.Auth())
	{

		userGroup := privateGroup.Group("/user")
		{
			userGroup.GET("/info", v1.GetUserInfo)
			userGroup.POST("/avatar", v1.UpdateAvatar)
			userGroup.POST("/info", v1.UpdateUserInfo)
		}
		friendGroup := privateGroup.Group("/friend")
		{
			friendGroup.POST("/request", v1.SendFriendRequest)
			friendGroup.POST("/process", v1.HandleFriendRequest)
			friendGroup.GET("/list", v1.GetFriendList)
			friendGroup.GET("/pending", v1.GetPendingRequests)
			friendGroup.POST("/delete", v1.DeleteFriendRecord)
		}
		// 聊天相关路由
		chatGroup := privateGroup.Group("/chat")
		{
			chatGroup.GET("/history", v1.GetChatHistory)
			chatGroup.POST("/revoke", v1.RevokeMessage)
			chatGroup.POST("/read", v1.ReadMessage)
		}

		aiGroup := privateGroup.Group("/ai")
		{
			// GET /api/v1/ai/search?query=刚才谁说了去吃火锅
			aiGroup.GET("/search", v1.SemanticSearch)
		}
		//文件上传
		fileGroup := privateGroup.Group("/file")
		{
			fileGroup.POST("/upload", v1.Upload)
		}

		groupGroup := privateGroup.Group("/group")
		{
			// ... 之前的 public, join ...
			groupGroup.POST("/create", v1.CreateGroup)
			groupGroup.GET("/mine", v1.GetMyGroups)        // 获取我的群组
			groupGroup.POST("/update", v1.UpdateGroup)     // 更新群资料
			groupGroup.GET("/members", v1.GetGroupMembers) // 获取群成员
			groupGroup.POST("/invite", v1.InviteMember)    // 邀请好友
			groupGroup.POST("/kick", v1.KickMember)        // 踢人
			groupGroup.POST("/mute", v1.MuteMember)        // 禁言
			groupGroup.POST("/admin", v1.SetAdmin)         // 设置管理员
			groupGroup.POST("/transfer", v1.TransferGroup) // 转让群主
		}
	}
	return r
}
