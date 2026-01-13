package socket

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-nexus/internal/core"
	"go-nexus/internal/model/dto"
	"go-nexus/pkg/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

// --- 配置常量 (复试考点：系统调优) ---
const (
	// 写超时时间
	writeWait = 10 * time.Second
	// 读超时时间 (心跳间隔)
	// 如果 60秒 没收到客户端的 Pong，就认为它断了
	pongWait = 60 * time.Second
	// 发送 Ping 的间隔 (必须小于 pongWait)
	pingPeriod = (pongWait * 9) / 10
	// 最大消息大小 (防止恶意发大包把服务器内存撑爆)
	maxMessageSize = 5120 // 5KB
)

// Client 代表一个 WebSocket 连接用户
type Client struct {
	ID     uint            // 用户ID
	Socket *websocket.Conn // 具体的底层连接
	Send   chan []byte     // 待发送消息的缓冲通道 (Outbound)
}

var handlers = map[int]MsgHandler{
	dto.ChatTypePrivate: &PrivateHandler{},
	dto.ChatTypeGroup:   &GroupHandler{},
	dto.ChatTypeSystem:  &SystemHandler{},
}

// --------------------------------------------------------------------------------
// 1. ReadPump: 读泵 (只管从 Socket 读，然后扔给后端处理)
// --------------------------------------------------------------------------------
func (c *Client) ReadPump() {
	// 确保函数退出时关闭连接并注销
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	// 配置 Socket 参数
	c.Socket.SetReadLimit(maxMessageSize)
	c.Socket.SetReadDeadline(time.Now().Add(pongWait)) // 设置读取死线

	// 设置 Pong 处理器：收到客户端的 Pong 后，延长死线
	c.Socket.SetPongHandler(func(string) error {
		c.Socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// 1. 阻塞读取消息
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			// 如果读出错了（比如客户端强退），跳出循环，触发 defer 销毁连接
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 2. 解析 JSON 协议
		var proto dto.ProtocolMsg
		if err := json.Unmarshal(message, &proto); err != nil {
			log.Printf("Json解析错误: %v", err)
			continue // 格式不对，忽略这条，继续读下一条
		}

		// 3. 处理心跳 (Type=0)
		// 实际上有了 SetPongHandler，这里的业务层心跳可以简化，或者用于前端业务逻辑保活
		if proto.Type == dto.TypeHeartbeat {
			continue
		}

		// ==========================================
		// 🔥 [核心代码]：将消息投喂给 AI
		// ==========================================
		// 只有文本消息才存 RAG，图片/语音暂时不存
		if proto.Type == dto.TypeText {
			// 生成一个简单的唯一ID (实际项目可以用 UUID)
			msgID := strconv.FormatInt(time.Now().UnixNano(), 10)
			nickname := proto.SenderNickname
			//计算SessionID
			sessionID := core.GetSessionID(proto.ChatType, c.ID, proto.ToUserID)
			// 异步发送 (go func)，绝对不能阻塞聊天主线程！
			core.AsyncSyncMessage(c.ID, proto.Content, msgID, nickname, sessionID)
		}
		// 4. 【安全关键】强制绑定发送者 ID
		// 无论前端传什么 from_user_id，都覆盖为当前连接的 ID
		// 防止黑客拿 A 的 Token 连上来，却发包说自己是 B
		proto.FromUserID = c.ID

		//2.策略分发
		worker, ok := handlers[proto.ChatType]
		if !ok {
			log.Printf("未知的聊天类型：: %v", proto.ChatType)
			continue
		}
		if err := worker.Handle(&proto); err != nil {
			log.Printf("消息处理出错: %v", err)
		}
	}
}

// --------------------------------------------------------------------------------
// 2. WritePump: 写泵 (只管把 Send 通道里的数据写给 Socket)
// --------------------------------------------------------------------------------
func (c *Client) WritePump() {
	// 定时器：每隔 54秒 给前端发一个 Ping
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()

	for {
		select {
		// A. 业务消息：Send 通道有数据了
		case message, ok := <-c.Send:
			// 设置写入超时
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 通道被 Manager 关闭了 (比如踢人下线)
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 获取 Writer 对象
			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// 写入数据
			w.Write(message)

			// 优化：如果你连着发了 10 条消息，Send 通道里积压了 10 条
			// 这里会一次性把缓冲区的都拿出来，合并成一个 TCP 包发出去，减少网络 IO
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		// B. 心跳保活：定时器触发
		case <-ticker.C:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			// 发送 Ping 帧 (Control Frame)，前端浏览器会自动回复 Pong，不需要前端写代码处理
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// --------------------------------------------------------------------------------
// 3. HTTP 升级入口
// --------------------------------------------------------------------------------

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域 (必加！否则前端 Vue 连不上)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ConnectWebSocket 处理 WebSocket 连接请求
func ConnectWebSocket(c *gin.Context) {
	// 1. 获取 Token
	token := c.Query("token")
	if token == "" {
		// 为了调试方便，有时也允许 header 传
		token = c.GetHeader("sec-websocket-protocol")
	}

	// 2. 鉴权
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "鉴权失败"})
		return
	}

	// 3. 升级 HTTP -> WebSocket
	// responseHeader 传 nil
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("升级失败: %v", err)
		return
	}

	// 4. 初始化 Client 实例
	client := &Client{
		ID:     claims.UserID,
		Socket: conn,
		Send:   make(chan []byte, 256), // 带缓冲通道，防止发送方阻塞
	}

	// 5. 注册到管理器
	Manager.Register <- client

	// 6. 启动双工协程
	// 一个负责读，一个负责写，互不阻塞
	go client.ReadPump()
	go client.WritePump()
}
