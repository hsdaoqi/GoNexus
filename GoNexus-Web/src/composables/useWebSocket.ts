import { ref, onMounted, onUnmounted, readonly } from 'vue'
import { useUserStore } from '../store/user'
import { useFriendStore } from '../store/friend'
import { ElMessage } from 'element-plus'

interface MessageHandler {
  (message: any): boolean // 返回是否已处理消息
}

export function useWebSocket() {
  const userStore = useUserStore()
  const friendStore = useFriendStore()
  const messageHandlers = ref<MessageHandler[]>([])

  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const reconnectTimer = ref<number | null>(null)

  // 建立WebSocket连接
  const connect = () => {
    if (socket.value?.readyState === WebSocket.OPEN) {
      return // 已连接
    }

    const token = localStorage.getItem('token')
    if (!token) {
      console.warn('没有token，无法建立WebSocket连接')
      return
    }

    try {
      const wsUrl = `ws://localhost:8080/socket?token=${encodeURIComponent(token)}`
      socket.value = new WebSocket(wsUrl)

      socket.value.onopen = () => {
        console.log('WebSocket连接已建立')
        isConnected.value = true
        ElMessage.success('实时连接已建立')

        // 将连接暴露到全局，供其他组件使用
        ;(window as any).globalWebSocket = socket.value

        // 取消重连定时器
        if (reconnectTimer.value) {
          clearTimeout(reconnectTimer.value)
          reconnectTimer.value = null
        }
      }

      socket.value.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)

          // 先让自定义处理器处理
          let handled = false
          for (const handler of messageHandlers.value) {
            if (handler(message) === true) {
              handled = true
              break
            }
          }

          // 如果没有被自定义处理器处理，则使用默认处理
          if (!handled) {
            handleMessage(message)
          }
        } catch (error) {
          console.error('解析WebSocket消息失败:', error)
        }
      }

      socket.value.onclose = () => {
        console.log('WebSocket连接已断开')
        isConnected.value = false

        // 清理全局引用
        if ((window as any).globalWebSocket === socket.value) {
          ;(window as any).globalWebSocket = null
        }

        // 自动重连
        scheduleReconnect()
      }

      socket.value.onerror = (error) => {
        console.error('WebSocket连接错误:', error)
        isConnected.value = false
      }

    } catch (error) {
      console.error('建立WebSocket连接失败:', error)
      scheduleReconnect()
    }
  }

  // 处理接收到的消息（默认处理器）
  const handleMessage = (message: any): boolean => {
    console.log('收到WebSocket消息:', message)

    // 处理好友状态变化
    if (message.type === 'friend_status_change' ||
        message.type === 'user_online' ||
        message.type === 'user_offline') {
      friendStore.handleWebSocketMessage(message)
      return true
    }

    // 处理系统消息
    if (message.type === 'system') {
      ElMessage.info(message.content || '系统消息')
      return true
    }

    return false // 未处理
  }

  // 注册消息处理器
  const addMessageHandler = (handler: MessageHandler) => {
    messageHandlers.value.push(handler)
  }

  // 移除消息处理器
  const removeMessageHandler = (handler: MessageHandler) => {
    const index = messageHandlers.value.indexOf(handler)
    if (index > -1) {
      messageHandlers.value.splice(index, 1)
    }
  }

  // 断开连接
  const disconnect = () => {
    if (socket.value && socket.value.readyState === WebSocket.OPEN) {
      socket.value.close()
      socket.value = null
    }
    isConnected.value = false

    if (reconnectTimer.value) {
      clearTimeout(reconnectTimer.value)
      reconnectTimer.value = null
    }
  }

  // 计划重连
  const scheduleReconnect = () => {
    if (reconnectTimer.value) return

    reconnectTimer.value = setTimeout(() => {
      console.log('尝试重连WebSocket...')
      connect()
    }, 5000)
  }

  // 发送消息
  const sendMessage = (message: any) => {
    if (socket.value && socket.value.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify(message))
    } else {
      console.warn('WebSocket未连接，无法发送消息')
    }
  }

  // 生命周期管理
  onMounted(() => {
    // 只有在用户已登录的情况下才建立连接
    if (userStore.userInfo.id && localStorage.getItem('token')) {
      connect()
    }
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    socket: readonly(socket),
    isConnected: readonly(isConnected),
    connect,
    disconnect,
    sendMessage,
    addMessageHandler,
    removeMessageHandler
  }
}
