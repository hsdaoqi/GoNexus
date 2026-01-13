# GoNexus - 基于RAG的智能社交聊天平台

## 项目简介

GoNexus是一个集成了AI能力的智能社交聊天平台，基于RAG（检索增强生成）架构实现智能聊天总结、语义搜索和智能回复建议等功能。

## 核心功能

- **实时聊天**: WebSocket实现的实时通信
- **好友系统**: 用户好友管理和群聊功能
- **AI聊天总结**: 自动生成群聊摘要
- **语义搜索**: 基于向量的智能搜索
- **智能回复建议**: AI生成回复候选

## 技术栈

### 后端 (Go)
- Gin: Web框架
- GORM: ORM框架
- gorilla/websocket: WebSocket通信
- gRPC: 微服务通信

### AI服务 (Python)
- FastAPI: Web框架
- LangChain: LLM应用框架
- Milvus/Faiss: 向量数据库
- OpenAI/Claude: LLM API

### 前端
- React/Vue.js + TypeScript
- Socket.io-client: 实时通信

### 数据库
- MySQL: 关系型数据存储
- Milvus: 向量数据库

## 项目结构

```
GoNexus/
├── backend/go-nexus/     # Go后端服务
├── ai-service/          # Python AI服务
├── chat-frontend/       # 前端应用
├── docker/              # Docker配置
├── docs/                # 项目文档
└── scripts/             # 部署脚本
```

## 快速开始

### 环境要求
- Go 1.21+
- Python 3.9+
- MySQL 8.0+
- Node.js 18+

### 安装步骤
1. 克隆项目
2. 配置环境变量
3. 启动数据库
4. 启动后端服务
5. 启动AI服务
6. 启动前端

## 开发计划

- [x] 项目初始化
- [ ] Go后端基础架构
- [ ] 数据库设计
- [ ] WebSocket聊天功能
- [ ] gRPC服务通信
- [ ] Python AI服务
- [ ] 向量数据库集成
- [ ] AI聊天总结功能
- [ ] 语义搜索功能
- [ ] 智能回复建议
- [ ] 前端界面开发
- [ ] Docker部署

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License


