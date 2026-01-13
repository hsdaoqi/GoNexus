<template>
<div class="app-container">
    <vue-particles id="tsparticles" :options="particlesOptions" />

    <!-- 1. 左侧组件 -->
    <div class="layout-side">
    <SidePanel ref="sidePanelRef" />
    </div>

    <!-- 2. 中间组件 -->
    <div class="layout-main">
    <ChatWindow @sendMessage="handleSendSocket" @toggleAI="showAI = !showAI"/>
    </div>

    <!-- 3. 右侧组件 (动画过渡) -->
    <transition name="slide-fade">
    <div v-if="showAI" class="layout-right">
        <AIPanel @close="showAI = false" />
    </div>
    </transition>
</div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import SidePanel from './SidePanel.vue'
import ChatWindow from './ChatWindow.vue'
import AIPanel from './AIPanel.vue'
import { useChatStore } from '../../store/chat'
import { useUserStore } from '@/store/user'
import { useFriendStore } from '@/store/friend'
import { useWebSocket } from '@/composables/useWebSocket'
import { ms } from 'element-plus/es/locales.mjs'

const userStore = useUserStore()
const router = useRouter()
const chatStore = useChatStore()
const friendStore = useFriendStore()
const { sendMessage, addMessageHandler, removeMessageHandler } = useWebSocket()
const showAI = ref(false)
const sidePanelRef = ref()

// Chat页面的消息处理器
const chatMessageHandler = (message: any) => {
  if (message.type === 0) {
    // 心跳包，忽略
    return true
  } else if (message.type === 4) {
    // 刷新左侧面板的待处理申请列表
    sidePanelRef.value?.refreshPendingList()
    return true
  } else if (message.type === 5) {
    sidePanelRef.value?.initData()
    return true
  } else {
    // 普通聊天消息
    chatStore.addMessage(message)
    return true
  }
}

// 注册消息处理器
onMounted(() => {
  addMessageHandler(chatMessageHandler)
})

onUnmounted(() => {
  removeMessageHandler(chatMessageHandler)
})

// Chat页面不再主动建立WebSocket连接，使用全局连接
// 消息处理逻辑已经移到useWebSocket组合函数中

// 处理发送
const handleSendSocket = (text: string) => {
    if (!chatStore.currentChat) return

    const msg = {
        type: 1,
        to_user_id: chatStore.currentChat.id,
        chat_type: 1,
        content: text,
        sender_nickname: userStore.userInfo.nickname
    }
    sendMessage(msg)

    // 乐观更新：立即显示自己发送的消息
    const currentTime = new Date().toISOString().slice(0, 19).replace('T', ' ')
    const optimisticMsg = {
        ...msg,
        from_user_id: parseInt(localStorage.getItem('user_id') || '0'),
        sender_avatar: userStore.userInfo.avatar, // 暂时为空，后续可从store获取
        sender_nickname: userStore.userInfo.nickname, // 或者从store获取当前用户昵称
        send_time: currentTime, // 添加发送时间用于去重
        created_at: currentTime
    }
    chatStore.addMessage(optimisticMsg)
}

// 粒子配置
const particlesOptions = {
background: { color: { value: "transparent" } },
fpsLimit: 60,
particles: { color: { value: "#333333" }, links: { color: "#333333", distance: 150, enable: true, opacity: 0.05, width: 1 }, move: { enable: true, speed: 0.5 }, number: { value: 30 }, opacity: { value: 0.1 } }
};
</script>

<style scoped>
.app-container {
width: 100vw; height: 100vh; display: flex; overflow: hidden; background-color: #f7f7f7; font-family: 'Segoe UI', sans-serif; position: relative;
}
#tsparticles { position: absolute; top: 0; left: 0; width: 100%; height: 100%; z-index: 0; pointer-events: none; }

/* 布局控制 */
.layout-side {
width: 260px;
height: 100%;
flex-shrink: 0;
z-index: 10;
}

.layout-main {
flex: 1;
height: 100%;
min-width: 0; /* 防止Flex子项溢出 */
z-index: 5;
}

.layout-right {
width: 320px;
height: 100%;
flex-shrink: 0;
z-index: 10;
}

/* 动画 */
.slide-fade-enter-active, .slide-fade-leave-active { transition: all 0.3s ease; }
.slide-fade-enter-from, .slide-fade-leave-to { transform: translateX(20px); opacity: 0; }
</style>