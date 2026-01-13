<template>
<div class="chat-panel">
    <!-- Header -->
    <div class="chat-header">
    <div class="header-info" v-if="chatStore.currentChat">
        <span class="target-name">{{ chatStore.currentChat.nickname || chatStore.currentChat.username }}</span>
        <span class="target-status">Level 99</span>
    </div>
    <div class="header-info" v-else>
        <span class="target-name">未选择目标</span>
    </div>
    <button class="ai-trigger-btn" @click="$emit('toggleAI')">
        <el-icon><Cpu /></el-icon> YUI-SYSTEM
    </button>
    </div>

    <!-- Message Area -->
    <div class="message-area" ref="msgAreaRef">
    <div v-if="!chatStore.currentChat" class="empty-state">
        <img src="../../assets/login1.png" style="width: 120px; opacity: 0.3; filter: grayscale(100%);" />
        <p>请在左侧选择联络对象</p>
    </div>
    <div v-else v-for="(msg, index) in chatStore.messages" 
        :key="index"
        :class="['message-row', { 'me': msg.from_user_id === userStore.userInfo.id }]"
    >
        <el-avatar class="msg-avatar" :size="36" :src="msg.sender_avatar" />
        <div class="msg-content">
            <div class="sender-name">{{ msg.sender_nickname }}</div>
            <div class="bubble">{{ msg.content }}</div>
        </div>
    </div>
    </div>

    <!-- Input Area -->
    <div class="input-area">
    <textarea v-model="content" class="sao-textarea" placeholder="发送指令..." @keydown.enter.prevent="send"></textarea>
    <div class="toolbar-bottom">
        <button class="send-btn" @click="send">SEND</button>
    </div>
    </div>
</div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import { Cpu } from '@element-plus/icons-vue'
import { useChatStore } from '../../store/chat'
import { useUserStore } from '../../store/user'

const emit = defineEmits(['sendMessage', 'toggleAI'])
const chatStore = useChatStore()
const userStore = useUserStore()
const content = ref('')
const msgAreaRef = ref()


const send = () => {
if(!content.value.trim()) return
// 触发父组件的事件发送 WebSocket
emit('sendMessage', content.value) 
content.value = ''
}



// 监听消息列表变化，自动滚动到底部
watch(() => chatStore.messages.length, () => {
nextTick(() => {
    if (msgAreaRef.value) msgAreaRef.value.scrollTop = msgAreaRef.value.scrollHeight
})
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700&display=swap');

.chat-panel {
width: 100%; height: 100%; display: flex; flex-direction: column; background: #f5f7fa; position: relative; z-index: 5;
}
.chat-header {
height: 60px; background: white; border-bottom: 1px solid #e0e0e0; display: flex; justify-content: space-between; align-items: center; padding: 0 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}
.target-name { font-size: 16px; font-weight: bold; color: #333; }
.target-status { margin-left: 10px; font-size: 12px; color: #2ecc71; }
.ai-trigger-btn { background: white; border: 1px solid #e67e22; color: #e67e22; padding: 6px 12px; border-radius: 20px; cursor: pointer; font-family: 'Orbitron'; display: flex; align-items: center; gap: 5px; font-size: 12px; }
.ai-trigger-btn:hover { background: #e67e22; color: white; }

.message-area { flex: 1; padding: 20px; overflow-y: auto; display: flex; flex-direction: column; gap: 15px; }
.message-row { display: flex; align-items: flex-start; gap: 10px; }
.message-row.me { flex-direction: row-reverse; }
.msg-content { max-width: 70%; }
.sender-name { font-size: 12px; color: #999; margin-bottom: 2px; text-align: left; }
.me .sender-name { text-align: right; }
.bubble { padding: 10px 14px; border-radius: 8px; font-size: 14px; line-height: 1.5; box-shadow: 0 1px 2px rgba(0,0,0,0.1); background: white; color: #333; }
.me .bubble { background: #95ec69; color: black; }

.empty-state { height: 100%; display: flex; flex-direction: column; justify-content: center; align-items: center; color: #999; }

.input-area { height: 160px; background: white; border-top: 1px solid #e0e0e0; display: flex; flex-direction: column; padding: 10px 20px; }
.sao-textarea { flex: 1; border: none; resize: none; outline: none; font-size: 14px; font-family: inherit; margin-bottom: 10px; }
.toolbar-bottom { display: flex; justify-content: flex-end; }
.send-btn { background: #f5f5f5; color: #606266; border: 1px solid #dcdfe6; padding: 6px 20px; border-radius: 4px; cursor: pointer; }
.send-btn:hover { background: #4facfe; color: white; border-color: #4facfe; }
</style>