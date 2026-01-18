<template>
<div class="chat-panel">
    <!-- Header -->
    <div class="chat-header">
    <div class="header-info" v-if="chatStore.currentChat">
        <span class="target-name">{{ chatStore.currentChat.nickname || chatStore.currentChat.username }}</span>
        <span class="target-status" v-if="!chatStore.currentChat.isGroup">
            <span :class="['status-dot-small', { online: currentOnlineStatus }]"></span>
            {{ currentOnlineStatus ? 'Online' : 'Offline' }}
        </span>
        <span class="target-status" v-else>Level 99</span>
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
        :class="['message-row', { 'me': msg.from_user_id === userStore.userInfo.id, 'system-row': msg.type === 6 || msg.msg_type === 6 || msg.is_revoked }]"
    >
        <!-- 系统消息 -->
        <template v-if="msg.type === 6 || msg.msg_type === 6 || msg.is_revoked">
             <div class="system-msg-content">{{ msg.is_revoked ? (msg.from_user_id === userStore.userInfo.id ? '你撤回了一条消息' : '对方撤回了一条消息') : msg.content }}</div>
        </template>
        <!-- 普通消息 -->
        <template v-else>
            <el-avatar class="msg-avatar" :size="36" :src="msg.sender_avatar" />
            <div class="msg-content">
                <div class="sender-name">{{ msg.sender_nickname }}</div>
                <!-- 文本消息 -->
                <div v-if="msg.type === 1 || msg.msg_type === 1 || (!msg.type && !msg.msg_type)" 
                     class="bubble"
                     @contextmenu="showContextMenu($event, msg)"
                >{{ msg.content }}</div>
                <!-- 图片消息 -->
                <div v-else-if="msg.type === 2 || msg.msg_type === 2" 
                     class="bubble image-bubble"
                     @contextmenu="showContextMenu($event, msg)"
                >
                <el-image 
                    :src="msg.url" 
                    :preview-src-list="[msg.url]"
                    fit="cover"
                    style="max-width: 200px; border-radius: 4px;"
                />
                </div>
                <!-- 语音消息 -->
                <div v-else-if="msg.type === 3 || msg.msg_type === 3" 
                     class="bubble audio-bubble" 
                     style="padding: 5px 10px;"
                     @contextmenu="showContextMenu($event, msg)"
                >
                <audio controls :src="msg.url" style="height: 32px; width: 220px; vertical-align: middle;"></audio>
                </div>
            </div>
        </template>
    </div>
    </div>

    <!-- 右键菜单 -->
    <div v-if="contextMenu.visible" 
         class="msg-context-menu" 
         :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
         @click.stop
    >
      <div class="menu-item" @click="doRevoke">撤回</div>
    </div>

    <!-- Input Area -->
    <div class="input-area">
      <div class="toolbar-top">
         <el-icon :class="['tool-icon', { 'disabled': isMuted }]" @click="!isMuted && triggerUpload()" title="发送图片"><Paperclip /></el-icon>
         <input type="file" ref="fileInput" style="display: none" @change="handleFileChange" accept="image/*" :disabled="isMuted" />
         
         <el-icon 
            :class="['tool-icon', { 'recording-active': isRecording, 'disabled': isMuted }]" 
            @mousedown="!isMuted && startRecording()" 
            @mouseup="!isMuted && stopRecording()"
            @mouseleave="!isMuted && stopRecording()"
            title="按住说话"
            style="margin-left: 10px;"
         >
            <Microphone />
         </el-icon>
         <span v-if="isRecording" class="recording-tip">正在录音...松开发送</span>
      </div>
      <textarea 
        v-model="content" 
        class="sao-textarea" 
        :placeholder="isMuted ? '您已被禁言' : '发送指令...'" 
        @keydown.enter.prevent="send"
        :disabled="isMuted"
      ></textarea>
      <div class="toolbar-bottom">
        <button class="send-btn" @click="send" :disabled="isMuted" :class="{ 'btn-disabled': isMuted }">SEND</button>
      </div>
    </div>
</div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch, onMounted, onUnmounted, computed } from 'vue'
import { Cpu, Paperclip, Microphone, Document, Download } from '@element-plus/icons-vue'
import { useChatStore } from '../../store/chat'
import { useUserStore } from '../../store/user'
import { useFriendStore } from '../../store/friend'
import { uploadFile } from '../../api/file'
import { getGroupMembers } from '../../api/group'
import { revokeMessage } from '../../api/chat'
import { ElMessage } from 'element-plus'

const emit = defineEmits(['sendMessage', 'toggleAI'])
const chatStore = useChatStore()
const userStore = useUserStore()
const friendStore = useFriendStore()
const content = ref('')
const msgAreaRef = ref()
const fileInput = ref()
const isMuted = ref(false)

// 计算当前聊天对象的在线状态
const currentOnlineStatus = computed(() => {
  if (!chatStore.currentChat || chatStore.currentChat.isGroup) return false
  // 优先从 friendStore 获取实时状态
  const friend = friendStore.friends.find(f => f.id === chatStore.currentChat.id)
  return friend ? friend.isOnline : (chatStore.currentChat.isOnline || false)
})

// 右键菜单状态
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  msg: null as any
})

// 显示右键菜单
const showContextMenu = (e: MouseEvent, msg: any) => {
  // 只能撤回自己的消息
  if (msg.from_user_id !== userStore.userInfo.id) return
  // 系统消息或已撤回消息不能撤回
  if (msg.type === 6 || msg.msg_type === 6) return
  
  e.preventDefault()
  contextMenu.value = {
    visible: true,
    x: e.clientX,
    y: e.clientY,
    msg
  }
}

// 隐藏右键菜单
const hideContextMenu = () => {
  contextMenu.value.visible = false
}

// 执行撤回
const doRevoke = async () => {
  if (!contextMenu.value.msg) return
  const msg = contextMenu.value.msg
  // 优先使用 id (历史记录), 其次使用 msg_id (如果协议中有)
  const msgId = msg.id || msg.msg_id
  
  if (!msgId) {
    ElMessage.warning('无法撤回刚发送的消息(需刷新页面)')
    return
  }

  try {
    await revokeMessage({
       msg_id: msgId,
       chat_type: chatStore.currentChat.isGroup ? 2 : 1,
       target_id: chatStore.currentChat.id
    })
    // 乐观更新
    chatStore.handleRevokeMessage(msgId)
    hideContextMenu()
  } catch (e) {
    console.error('Revoke failed', e)
    ElMessage.error('撤回失败')
  }
}

// 全局点击关闭菜单
onMounted(() => {
  document.addEventListener('click', hideContextMenu)
})
onUnmounted(() => {
  document.removeEventListener('click', hideContextMenu)
})

// 监听当前会话变化，检查是否被禁言
watch(() => chatStore.currentChat?.id, async (newVal) => {
    if (!newVal || !chatStore.currentChat.isGroup) {
        isMuted.value = false
        return
    }
    
    // 如果是群聊，检查自己在群里的状态
    try {
        const members: any = await getGroupMembers(newVal)
        const myId = userStore.userInfo.id
        const list = Array.isArray(members) ? members : members.data || []
        const me = list.find((m: any) => m.user_id === myId)
        if (me && me.muted === 1) {
            isMuted.value = true
        } else {
            isMuted.value = false
        }
    } catch (e) {
        console.error('Check mute status failed', e)
        isMuted.value = false
    }
}, { immediate: true })

// 录音相关
const isRecording = ref(false)
const mediaRecorder = ref<MediaRecorder | null>(null)
const audioChunks = ref<Blob[]>([])

const triggerUpload = () => {
  fileInput.value.click()
}

const handleFileChange = async (e: Event) => {
  const files = (e.target as HTMLInputElement).files
  if (!files || files.length === 0) return
  
  const file = files[0]
  // 简单校验 (放宽到 50MB)
  if (!file || file.size > 50 * 1024 * 1024) {
    ElMessage.warning('文件大小不能超过 50MB')
    return
  }

  const formData = new FormData()
  formData.append('file', file)
  
  try {
    const res: any = await uploadFile(formData)
    // 假设后端返回 { url: '...' }
    if (res.url) {
       // 判断文件类型
       const isImage = file.type.startsWith('image/')
       // 发送消息
       emit('sendMessage', isImage ? '[图片]' : '[文件]', isImage ? 2 : 9, res.url, file.name, file.size)
    }
  } catch (e) {
    console.error('Upload failed', e)
  } finally {
    // 清空 input 即使再次选择同一文件也能触发 change
    (e.target as HTMLInputElement).value = ''
  }
}

// 开始录音
const startRecording = async () => {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    mediaRecorder.value = new MediaRecorder(stream)
    audioChunks.value = []
    
    mediaRecorder.value.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.value.push(event.data)
      }
    }

    mediaRecorder.value.onstop = async () => {
      const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })
      const file = new File([audioBlob], `voice_${Date.now()}.webm`, { type: 'audio/webm' })
      
      const formData = new FormData()
      formData.append('file', file)

      try {
        const res: any = await uploadFile(formData)
        if (res.url) {
          emit('sendMessage', '[语音]', 3, res.url)
        }
      } catch (e) {
        console.error('Voice upload failed', e)
        ElMessage.error('语音发送失败')
      }
    }

    mediaRecorder.value.start()
    isRecording.value = true
  } catch (err) {
    console.error('Error accessing microphone', err)
    ElMessage.error('无法访问麦克风')
  }
}

// 停止录音
const stopRecording = () => {
  if (mediaRecorder.value && isRecording.value) {
    mediaRecorder.value.stop()
    isRecording.value = false
    // 停止所有轨道
    mediaRecorder.value.stream.getTracks().forEach(track => track.stop())
  }
}

const send = () => {
  if(!content.value.trim()) return
  // 默认 type=1 (文本)
  emit('sendMessage', content.value, 1, '') 
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
.target-status { margin-left: 10px; font-size: 12px; color: #2ecc71; display: flex; align-items: center; gap: 5px; }
.status-dot-small { width: 8px; height: 8px; border-radius: 50%; background: #95a5a6; display: inline-block; }
.status-dot-small.online { background: #2ecc71; }
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
.toolbar-top { margin-bottom: 5px; display: flex; gap: 10px; }
.tool-icon { font-size: 20px; color: #666; cursor: pointer; transition: color 0.3s; }
.tool-icon:hover { color: #4facfe; }
.sao-textarea { flex: 1; border: none; resize: none; outline: none; font-size: 14px; font-family: inherit; margin-bottom: 10px; }
.toolbar-bottom { display: flex; justify-content: flex-end; }
.image-bubble { padding: 5px; }
.send-btn { background: #f5f5f5; color: #606266; border: 1px solid #dcdfe6; padding: 6px 20px; border-radius: 4px; cursor: pointer; }
.send-btn:hover { background: #ecf5ff; border-color: #c6e2ff; color: #409eff; }

.recording-tip {
    font-size: 12px;
    color: #e67e22;
    margin-left: 10px;
    animation: blink 1s infinite;
}
.recording-active {
    color: #e74c3c !important;
}

@keyframes blink {
    0% { opacity: 1; }
    50% { opacity: 0.5; }
    100% { opacity: 1; }
}

/* 消息撤回菜单 */
.msg-context-menu {
  position: fixed;
  background: white;
  border: 1px solid #eee;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
  border-radius: 4px;
  padding: 5px 0;
  z-index: 9999;
}
.menu-item {
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
  color: #333;
}
.menu-item:hover {
  background: #f5f7fa;
  color: #409eff;
}

/* 系统消息样式 */
.system-row {
  justify-content: center;
}
.system-msg-content {
  background: #f0f0f0;
  color: #999;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 4px;
}
.send-btn:hover { background: #4facfe; color: white; border-color: #4facfe; }

.recording-active { color: #f56c6c !important; transform: scale(1.2); }
.recording-tip { font-size: 12px; color: #f56c6c; margin-left: 10px; line-height: 20px; animation: pulse 1.5s infinite; }
@keyframes pulse { 0% { opacity: 1; } 50% { opacity: 0.5; } 100% { opacity: 1; } }

.system-row { justify-content: center !important; }
.system-msg-content { background-color: #f5f5f5; color: #999; padding: 4px 12px; border-radius: 4px; font-size: 12px; }

.tool-icon.disabled { cursor: not-allowed; opacity: 0.5; }
.tool-icon.disabled:hover { color: #666; }
.btn-disabled { background: #e0e0e0 !important; color: #999 !important; border-color: #dcdfe6 !important; cursor: not-allowed; }

/* 文件消息样式 */
.file-bubble {
  padding: 0 !important;
  overflow: hidden;
  width: 240px;
}
.file-card {
  display: flex;
  align-items: center;
  padding: 10px;
  background: #f9f9f9;
  border-radius: 8px;
}
.file-icon {
  font-size: 40px;
  color: #909399;
  margin-right: 10px;
}
.file-info {
  flex: 1;
  overflow: hidden;
}
.file-name {
  font-size: 14px;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}
.file-size {
  font-size: 12px;
  color: #999;
}
.download-btn {
  margin-left: 10px;
  color: #409eff;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  text-decoration: none;
}
.download-btn:hover {
  color: #66b1ff;
}
</style>