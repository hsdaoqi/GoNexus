## GoNexus项目接下来要完成的功能

### 1. 当前项目状态分析
- ✅ 用户注册、登录、JWT鉴权功能已完成
- ✅ 好友系统（申请、处理、列表、删除）已完成
- ✅ 基本WebSocket通信框架已搭建
- ✅ 消息存储和历史记录查询功能已实现
- ⚠️ 聊天历史记录API接口未实现
- ⚠️ WebSocket管理器未启动
- ⚠️ WebSocket消息处理逻辑不完整
- ⚠️ 离线消息功能未实现
- ⚠️ AI相关功能未实现

### 2. 接下来要完成的功能

#### 2.1 聊天历史记录API实现
- **文件**：`internal/api/v1/chat.go`
- **功能**：实现获取单聊历史记录的HTTP接口
- **路由**：`GET /api/v1/chat/history`
- **参数**：targetID（对方ID）、offset（偏移量）、limit（条数）
- **返回**：聊天历史记录列表

#### 2.2 WebSocket管理器启动
- **文件**：`cmd/server/main.go`
- **功能**：在main函数中启动WebSocket管理器协程
- **代码**：`go socket.Manager.Start()`

#### 2.3 WebSocket路由添加
- **文件**：`internal/api/router.go`
- **功能**：添加WebSocket连接的路由
- **路由**：`GET /ws`
- **处理函数**：实现WebSocket握手和连接管理

#### 2.4 WebSocket消息处理完善
- **文件**：`internal/core/socket/client.go`
- **功能**：完善消息接收、处理和转发逻辑
- **流程**：接收消息 → 解析协议 → 调用SaveAndTransform存储消息 → 转发给目标用户

#### 2.5 离线消息功能实现
- **文件**：`internal/repository/message_repo.go`、`internal/service/chat_service.go`
- **功能**：当接收用户不在线时，存储离线消息，待用户上线后推送

#### 2.6 AI聊天总结功能
- **文件**：`internal/service/ai_service.go`
- **功能**：调用AI服务生成聊天摘要
- **触发条件**：定时或手动触发群聊总结

#### 2.7 群聊功能实现
- **文件**：`internal/model/message.go`、`internal/repository/message_repo.go`
- **功能**：实现群聊创建、消息发送、历史记录查询等

### 3. 实施优先级
1. **聊天历史记录API**（最基础的HTTP接口，前端急需）
2. **WebSocket管理器启动**（核心通信功能）
3. **WebSocket路由和消息处理**（实时聊天功能）
4. **离线消息功能**（提升用户体验）
5. **群聊功能**（扩展聊天场景）
6. **AI相关功能**（增强智能特性）

### 4. 预期效果
- 用户可以通过HTTP接口获取聊天历史记录
- 用户可以通过WebSocket进行实时聊天
- 离线消息能够被正确存储和推送
- 群聊功能正常工作
- AI能够生成聊天摘要和智能回复建议

### 5. 技术要点
- 严格遵守分层架构（API → Service → Repository）
- 使用GORM进行数据库操作
- 使用WebSocket实现实时通信
- 考虑并发安全和性能优化
- 实现良好的错误处理和日志记录