package main

import (
	"fmt"
	"go-nexus/internal/api"
	"go-nexus/internal/core"
	"go-nexus/internal/core/socket"
	"go-nexus/internal/model"
	"go-nexus/internal/service"
	"go-nexus/pkg/global"
	"go-nexus/pkg/initialize"
)

func main() {
	//初始化
	initialize.InitConfig()
	initialize.InitMySQL()
	initialize.InitOSS()
	core.InitAIClient()
	service.FriendSvc = service.NewFriendService(&socket.Manager)
	//global.DB.AutoMigrate(&model.Message{})
	//global.DB.AutoMigrate(&model.User{})
	//global.DB.AutoMigrate(&model.Friend{})
	//global.DB.AutoMigrate(&model.FriendRequest{})
	//global.DB.AutoMigrate(&model.Group{})
	//global.DB.AutoMigrate(&model.GroupMember{})
	global.DB.AutoMigrate(&model.Post{}, &model.Comment{}, &model.Like{})
	//初始化路由
	r := api.InitRouter()

	// ---------------------------------------------------------
	// 🔥 新增：WebSocket 路由
	// ---------------------------------------------------------
	// 注意 1: 这是一个 GET 请求
	// 注意 2: 它不使用 middleware.Auth()，因为它在内部自己处理了 Query Token 鉴权
	r.GET("/socket", socket.ConnectWebSocket)

	go socket.Manager.Start()

	//启动服务
	port := fmt.Sprintf(":%d", global.Config.Server.Port)
	fmt.Printf("服务运行在 %s\n", port)
	r.Run(port)
}
