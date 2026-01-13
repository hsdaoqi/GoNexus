import { defineStore } from 'pinia'
import { getFriendList } from '../api/friend'

export interface Friend {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
  signature: string
  gender: string
  birthday: string
  location: string
  isOnline: boolean
  lastSeen: string
}

export const useFriendStore = defineStore('friend', {
  state: () => ({
    friends: [] as Friend[],
    onlineCount: 0,
    isLoading: false,
    lastUpdate: null as Date | null,
    pollingTimer: null as number | null,
    webSocket: null as WebSocket | null,
    reconnectTimer: null as number | null
  }),

  getters: {
    onlineFriends: (state) => state.friends.filter(friend => friend.isOnline === true),
    offlineFriends: (state) => state.friends.filter(friend => friend.isOnline === false),
    // 确保lastUpdate始终返回有效的Date对象
    lastUpdateDate: (state) => {
      if (!state.lastUpdate) return null
      return state.lastUpdate instanceof Date ? state.lastUpdate : new Date(state.lastUpdate)
    }
  },

  actions: {
    // 获取好友列表
    async fetchFriends() {
      this.isLoading = true
      try {
        const res: any = await getFriendList()
        console.log(res)
        // 假设后端返回的好友数据结构
        const friendsData = res || res.friends || res.data || []

        // 处理好友数据，确保有状态字段
        this.friends = friendsData.map((friend: any) => ({
          id: friend.id,
          username: friend.username,
          nickname: friend.nickname || friend.username,
          avatar: friend.avatar || '',
          email: friend.email || '',
          signature: friend.signature || '',
          gender: friend.gender || '',
          birthday: friend.birthday || '',
          location: friend.location || '',
          isOnline: friend.is_online || friend.isOnline || false,
          lastSeen: friend.last_seen || friend.lastSeen || new Date().toISOString()
        }))

        this.updateOnlineCount()
        this.lastUpdate = new Date()

        console.log('好友列表已更新:', this.friends.length, '个好友，在线:', this.onlineCount)
      } catch (error) {
        console.error('获取好友列表失败:', error)
        // 如果获取失败，设置为空数组
        this.friends = []
        this.onlineCount = 0
      } finally {
        this.isLoading = false
      }
    },

    // 更新在线人数统计
    updateOnlineCount() {
      this.onlineCount = this.friends.filter(friend => friend.isOnline === true).length
    },

    // 模拟更新好友状态（用于测试）
    updateFriendStatus(friendId: number, isOnline: boolean) {
      const friend = this.friends.find(f => f.id === friendId)
      if (friend) {
        friend.isOnline = isOnline
        friend.lastSeen = new Date().toISOString()
        this.updateOnlineCount()
      }
    },

    // 启动实时状态更新
    startStatusPolling(intervalMs: number = 30000) { // 默认30秒
      // 如果已经有定时器，先停止
      this.stopStatusPolling()

      // 先获取一次好友列表
      this.fetchFriends()

      // 检查是否已经有全局WebSocket连接（比如chat页面建立的）
      // 如果没有，则建立自己的连接
      if (!this.hasGlobalWebSocket()) {
        this.connectWebSocket()
      } else {
        console.log('检测到全局WebSocket连接，使用现有连接')
      }

      // 设置HTTP轮询作为备用方案（当WebSocket断开时）
      this.pollingTimer = setInterval(() => {
        // 只有在没有全局WebSocket或全局WebSocket断开时才使用HTTP轮询
        if (!this.hasGlobalWebSocket()) {
          this.fetchFriends()
        }
      }, intervalMs)

      console.log(`好友状态实时更新已启动，优先使用全局WebSocket`)
    },

    // 检查是否有全局WebSocket连接
    hasGlobalWebSocket(): boolean {
      // 检查window对象上是否有全局WebSocket实例
      // 这里可以根据实际情况调整检测逻辑
      return (window as any).globalWebSocket &&
             (window as any).globalWebSocket.readyState === WebSocket.OPEN
    },

    // 停止实时状态更新
    stopStatusPolling() {
      // 停止HTTP轮询
      if (this.pollingTimer) {
        clearInterval(this.pollingTimer)
        this.pollingTimer = null
      }

      // 只断开自己创建的WebSocket连接，不影响全局连接
      this.disconnectWebSocket()

      console.log('好友状态实时更新已停止')
    },

    // WebSocket连接管理
    connectWebSocket() {
      if (this.webSocket?.readyState === WebSocket.OPEN) {
        return // 已经连接
      }

      try {
        const token = localStorage.getItem('token')
        if (!token) {
          console.warn('没有token，无法建立WebSocket连接')
          return
        }

        // 建立WebSocket连接
        const wsUrl = `ws://localhost:8080/socket?token=${encodeURIComponent(token)}`
        this.webSocket = new WebSocket(wsUrl)

        this.webSocket.onopen = () => {
          console.log('好友状态WebSocket连接已建立')
          // 取消重连定时器
          if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer)
            this.reconnectTimer = null
          }
        }

        this.webSocket.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data)
            this.handleWebSocketMessage(message)
          } catch (error) {
            console.error('解析WebSocket消息失败:', error)
          }
        }

        this.webSocket.onclose = () => {
          console.log('好友状态WebSocket连接已断开')
          // 自动重连
          this.scheduleReconnect()
        }

        this.webSocket.onerror = (error) => {
          console.error('好友状态WebSocket连接错误:', error)
        }

      } catch (error) {
        console.error('建立WebSocket连接失败:', error)
        this.scheduleReconnect()
      }
    },

    // 处理WebSocket消息 (公开方法，供外部调用)
    handleWebSocketMessage(message: any) {
      console.log('收到WebSocket消息:', message)

      // 根据消息类型处理
      switch (message.type) {
        case 'friend_status_change':
          // 好友状态变化
          if (message.user_id && typeof message.is_online === 'boolean') {
            this.updateFriendStatus(message.user_id, message.is_online)
          }
          break
        case 'user_online':
          // 用户上线
          if (message.user_id) {
            this.updateFriendStatus(message.user_id, true)
          }
          break
        case 'user_offline':
          // 用户下线
          if (message.user_id) {
            this.updateFriendStatus(message.user_id, false)
          }
          break
        default:
          console.log('未知的WebSocket消息类型:', message.type)
      }
    },

    // 断开WebSocket连接
    disconnectWebSocket() {
      if (this.webSocket && this.webSocket.readyState === WebSocket.OPEN) {
        this.webSocket.close()
        this.webSocket = null
      }
      if (this.reconnectTimer) {
        clearTimeout(this.reconnectTimer)
        this.reconnectTimer = null
      }
    },

    // 计划重连
    scheduleReconnect() {
      if (this.reconnectTimer) return

      this.reconnectTimer = window.setTimeout(() => {
        console.log('尝试重连WebSocket...')
        this.connectWebSocket()
      }, 5000) // 5秒后重连
    },

    // 清空好友数据
    clearFriends() {
      this.friends = []
      this.onlineCount = 0
      this.lastUpdate = null
      this.disconnectWebSocket()
    }
  },

  // 持久化好友列表，但不持久化在线状态（因为状态会变化）
  persist: true
})
