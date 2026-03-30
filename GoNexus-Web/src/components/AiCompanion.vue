<template>
  <div class="ai-companion-container" :style="positionStyle" @mousedown="startDrag" @touchstart="startDrag">
    <!-- 对话气泡 -->
    <transition name="fade">
      <div v-if="showBubble" class="speech-bubble">
        <div class="bubble-content">
          <p>{{ currentMessage }}</p>
          <div class="actions" v-if="showActions">
            <el-button size="small" type="primary" link @click="handleAction('analyze')">分析当前心情</el-button>
            <el-button size="small" type="success" link @click="handleAction('chat')">闲聊一下</el-button>
          </div>
        </div>
      </div>
    </transition>

    <!-- 小精灵本体 -->
    <div class="companion-avatar" @click="toggleBubble" :class="['state-' + currentState, { 'is-active': showBubble }]">
      <!-- 这里可以使用 img 标签替换为您的银月图片 -->
      <img src="@/assets/silver-moon.png" alt="Silver Moon" />
      <!-- <div class="avatar-placeholder">
        <span class="emoji" v-if="currentState === 'idle'">🧚‍♀️</span>
        <span class="emoji" v-else-if="currentState === 'working'">🔮</span>
        <span class="emoji" v-else-if="currentState === 'alert'">🔔</span>
        <div class="glow-effect"></div>
      </div> -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useWebSocket } from '@/composables/useWebSocket'
import { useUserStore } from '@/store/user'

const emit = defineEmits(['action'])
const router = useRouter()
const { addMessageHandler, removeMessageHandler } = useWebSocket()
const userStore = useUserStore()

const showBubble = ref(false)
const showActions = ref(true)
const currentMessage = ref('主人，我在呢~')
const isAlert = ref(false) // 是否处于提醒状态

const currentState = ref<'idle' | 'working' | 'alert'>('idle')

// 监听 WebSocket 消息
const handleIncomingMessage = (msg: any) => {
  // 只关注聊天消息 (Type 1: 私聊, 2: 群聊)
  if (msg.type === 1 || msg.type === 2) {
    // 忽略自己发的消息
    if (msg.from_user_id === userStore.userInfo.id) return false

    isAlert.value = true
    currentState.value = 'alert' // 切换到提醒状态
    showBubble.value = true
    showActions.value = false // 消息提醒时不显示操作按钮

    // 判断是否有人 @我 (简单判断文本包含 @用户名)
    // 注意：这里假设 msg.content 是文本。如果是富文本可能需要解析
    const myName = userStore.userInfo.nickname || ''
    if (msg.type === 2 && myName && msg.content.includes('@' + myName)) {
      currentMessage.value = `主人，${msg.from_user_nickname || '有人'} 在群里 @你了！`
    } else {
      currentMessage.value = `主人，收到来自 ${msg.from_user_nickname || '好友'} 的新消息~`
    }

    // 5秒后自动收起气泡
    setTimeout(() => {
      if (isAlert.value) { // 如果还在提醒状态才收起
        showBubble.value = false
        isAlert.value = false
        currentState.value = 'idle' // 恢复待机
        showActions.value = true
      }
    }, 5000)
  }
  
  // 返回 false 以便让消息继续传递给 Store 处理
  return false
}

onMounted(() => {
  addMessageHandler(handleIncomingMessage)
})

onUnmounted(() => {
  removeMessageHandler(handleIncomingMessage)
})

// 拖拽逻辑
const isDragging = ref(false)
const position = ref({ x: window.innerWidth - 80, y: window.innerHeight - 80 })
const startPos = ref({ x: 0, y: 0 })

const positionStyle = computed(() => ({
  left: `${position.value.x}px`,
  top: `${position.value.y}px`,
  position: 'fixed' as const,
  zIndex: 9999
}))

const startDrag = (e: MouseEvent | TouchEvent) => {
  isDragging.value = false
  const clientX = e instanceof MouseEvent ? e.clientX : e.touches[0]?.clientX ?? 0
  const clientY = e instanceof MouseEvent ? e.clientY : e.touches[0]?.clientY ?? 0
  
  startPos.value = {
    x: clientX - position.value.x,
    y: clientY - position.value.y
  }

  const moveHandler = (moveEvent: MouseEvent | TouchEvent) => {
    isDragging.value = true
    const moveClientX = moveEvent instanceof MouseEvent ? moveEvent.clientX : moveEvent.touches[0]?.clientX ?? 0
    const moveClientY = moveEvent instanceof MouseEvent ? moveEvent.clientY : moveEvent.touches[0]?.clientY ?? 0
    
    let newX = moveClientX - startPos.value.x
    let newY = moveClientY - startPos.value.y
    
    // 边界限制
    const maxX = window.innerWidth - 60
    const maxY = window.innerHeight - 60
    
    if (newX < 0) newX = 0
    if (newX > maxX) newX = maxX
    if (newY < 0) newY = 0
    if (newY > maxY) newY = maxY
    
    position.value = { x: newX, y: newY }
  }

  const upHandler = () => {
    document.removeEventListener('mousemove', moveHandler)
    document.removeEventListener('mouseup', upHandler)
    document.removeEventListener('touchmove', moveHandler)
    document.removeEventListener('touchend', upHandler)
  }

  document.addEventListener('mousemove', moveHandler)
  document.addEventListener('mouseup', upHandler)
  document.addEventListener('touchmove', moveHandler)
  document.addEventListener('touchend', upHandler)
}

const toggleBubble = () => {
  if (isDragging.value) return // 如果是拖拽结束，不触发点击
  
  showBubble.value = !showBubble.value
  if (showBubble.value) {
    // 随机问候语
    const greetings = [
      '有什么我可以帮您的吗？',
      '今天天气真不错呢~',
      '我在观察大家的心情哦...',
      '银月一直在这里陪着主人。'
    ]
    const index = Math.floor(Math.random() * greetings.length)
    currentMessage.value = greetings[index] ?? '有什么我可以帮您的吗？'
  }
}

const handleAction = (type: string) => {
  emit('action', type)
  if (type === 'analyze') {
    currentState.value = 'working' // 切换到工作状态
    currentMessage.value = '正在分析广场上的气氛...'
    setTimeout(() => {
      showBubble.value = false
      currentState.value = 'idle' // 恢复待机
    }, 2000)
  }
}
</script>

<style scoped>
.ai-companion-container {
  cursor: pointer;
  user-select: none;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

/* 呼吸浮动动画 */
@keyframes float {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0px); }
}

.companion-avatar {
  width: 60px;
  height: 60px;
  position: relative;
  animation: float 3s ease-in-out infinite;
  transition: transform 0.3s;
}

/* 状态动画：工作 */
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 状态动画：提醒 */
@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-15px); }
}

.companion-avatar.state-working .avatar-placeholder {
  animation: spin 2s linear infinite;
  background: radial-gradient(circle at 30% 30%, #e1bee7, #ba68c8);
  box-shadow: 0 0 20px rgba(186, 104, 200, 0.6);
}

.companion-avatar.state-alert .avatar-placeholder {
  animation: bounce 0.5s ease-in-out infinite;
  background: radial-gradient(circle at 30% 30%, #ffcdd2, #e57373);
  box-shadow: 0 0 20px rgba(229, 115, 115, 0.6);
}

.companion-avatar.state-alert .glow-effect {
  box-shadow: 0 0 30px #ef5350;
  animation: pulse 0.5s infinite;
}

.companion-avatar:hover {
  transform: scale(1.1);
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: radial-gradient(circle at 30% 30%, #e0f7fa, #4dd0e1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  box-shadow: 0 0 20px rgba(77, 208, 225, 0.6);
  position: relative;
  transition: all 0.5s;
}

.emoji {
  position: relative;
  z-index: 2;
}

.glow-effect {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 50%;
  box-shadow: 0 0 30px #4dd0e1;
  opacity: 0.5;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 0.4; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
  100% { opacity: 0.4; transform: scale(1); }
}

.speech-bubble {
  background: white;
  padding: 12px 16px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  margin-bottom: 12px;
  max-width: 200px;
  position: relative;
  margin-right: 10px;
}

.speech-bubble::after {
  content: '';
  position: absolute;
  bottom: -8px;
  right: 20px;
  border-width: 8px 8px 0;
  border-style: solid;
  border-color: white transparent transparent transparent;
}

.bubble-content p {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #333;
  line-height: 1.4;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-start;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
