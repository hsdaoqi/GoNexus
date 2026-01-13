package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "go-nexus/internal/api/v1"
	"go-nexus/internal/middleware"
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
	}
	return r
}
