package socket

import "sync"

// ClientManager 连接管理器 (大管家)
type ClientManager struct {
	// Clients 存所有在线的连接
	// Key: UserID, Value: *Client
	// 为什么要用 RWMutex？因为会有很多协程同时读写这个 map，不加锁程序会崩溃 (复试必问)
	Clients     map[uint]*Client
	ClientsLock sync.RWMutex

	// Register 注册通道 (有人连上来了)
	Register chan *Client

	// Unregister 注销通道 (有人断开了)
	Unregister chan *Client

	// Broadcast 广播通道 (全员推送，比如系统公告)
	Broadcast chan []byte
}

// Manager 全局单例
// 这就是你在 client.go 里报错找不到的那个变量！
var Manager = ClientManager{
	Clients:    make(map[uint]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast:  make(chan []byte),
}

// Start 启动管理者 (需要在 main.go 里开启协程跑它)
func (manager *ClientManager) Start() {
	for {
		select {
		// 1. 处理注册 (上线)
		case client := <-manager.Register:
			manager.ClientsLock.Lock() // 加写锁
			manager.Clients[client.ID] = client
			manager.ClientsLock.Unlock()
			// 这里以后可以加日志：fmt.Printf("用户 %d 上线了\n", client.ID)

		// 2. 处理注销 (下线)
		case client := <-manager.Unregister:
			manager.ClientsLock.Lock()
			if _, ok := manager.Clients[client.ID]; ok {
				// 从 map 中删除
				delete(manager.Clients, client.ID)
				// 关闭通道，防止内存泄漏
				close(client.Send)
			}
			manager.ClientsLock.Unlock()
			// fmt.Printf("用户 %d 下线了\n", client.ID)

		// 3. 处理全员广播
		case message := <-manager.Broadcast:
			for _, conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn.ID)
				}
			}
		}
	}
}

// SendMessage 给指定用户发消息 (点对点)
func (manager *ClientManager) SendMessage(receiverID uint, message []byte) {
	manager.ClientsLock.RLock() // 加读锁 (不影响其他人并发读)
	defer manager.ClientsLock.RUnlock()

	if client, ok := manager.Clients[receiverID]; ok {
		// 只要把消息塞给那个人的 Send 通道，他的 WritePump 协程就会自动把消息发出去
		client.Send <- message
	} else {
		//TODO
		// 对方不在线
		// 思考题：如果不在线，这里应该调用 Service 把消息存入"离线消息表"或者推送到手机通知栏
	}
}

// IsUserOnline 实现OnlineStatusChecker接口，检查用户是否在线
func (manager *ClientManager) IsUserOnline(userID uint) bool {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	_, exists := manager.Clients[userID]
	return exists
}
