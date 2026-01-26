<template>
    <div class="dashboard-container">
      
      <!-- 1. 顶部导航 -->
      <GlobalNavbar />
  
      <div class="main-layout">
        <!-- 2. 左侧：主要内容区 -->
        <div class="content-area">
          
          <!-- A. 欢迎 & AI 快捷入口 (亮点功能) -->
          <div class="welcome-card">
            <h1>早安, {{ userStore.userInfo.username }} 👋</h1>
            <p class="subtitle">GoNexus AI 大脑已就绪。你可以直接在这里搜索知识库，或回顾聊天记忆。</p>
            
            <div class="ai-search-bar">
              <input 
                v-model="quickQuery" 
                placeholder="问问 AI：昨天大家在群里讨论了什么？" 
                @keydown.enter="handleQuickAsk"
              />
              <button @click="handleQuickAsk">
                <el-icon><MagicStick /></el-icon> Ask AI
              </button>
            </div>
            
            <!-- AI 快速回答展示 -->
            <div v-if="aiResponse" class="ai-quick-result">
              <div class="ai-tag">AI Insight</div>
              <p>{{ aiResponse }}</p>
            </div>
          </div>
  
          <!-- B. 最近会话 (Recent Chats) -->
          <div class="section-header">
            <h3>最近会话</h3>
            <el-button text bg size="small" @click="router.push('/chat')">进入完整模式 <el-icon class="el-icon--right"><ArrowRight /></el-icon></el-button>
          </div>
  
          <div class="chat-grid">
            <div v-if="recentChats.length === 0" class="empty-chat-tip" style="text-align: center; color: #999; width: 100%; padding: 20px;">
                暂无最近会话
            </div>
            <div 
                v-for="(chat, index) in recentChats" 
                :key="index" 
                class="chat-card" 
                @click="handleChatClick(chat)"
            >
              <div :class="['card-icon', chat.isGroup ? 'group-bg' : 'private-bg']">
                <el-icon v-if="chat.isGroup"><ChatDotRound /></el-icon>
                <el-icon v-else><User /></el-icon>
              </div>
              <div class="card-info">
                <h4>{{ chat.display_name }}</h4>
                <p class="text-truncate" style="max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
                    {{ chat.last_msg }}
                </p>
              </div>
              <span class="time">{{ formatTimeAgo(chat.updated_at) }}</span>
            </div>
          </div>
  
        </div>
  
        <!-- 3. 右侧：侧边栏 (概览数据) -->
        <div class="sidebar">
          <!-- 统计卡片 -->
          <div class="stat-card">
            <div class="stat-item">
              <div class="num">
                {{ friendStore.isLoading ? '...' : friendStore.onlineCount }}
                <span v-if="friendStore.lastUpdateDate" class="update-time">
                  {{ formatTimeAgo(friendStore.lastUpdateDate) }}
                </span>
              </div>
              <div class="label">好友在线</div>
            </div>
            <div class="divider"></div>
            <div class="stat-item">
              <div class="num">85%</div>
              <div class="label">AI 记忆容量</div>
            </div>
          </div>
  
          <!-- 系统公告 -->
          <div class="widget-box">
            <div class="widget-title">📢 系统公告</div>
            <ul class="notice-list">
              <li>🚀 GoNexus v1.0 正式上线</li>
              <li>🤖 RAG 向量数据库已重置</li>
              <li>🔒 新增端到端加密支持</li>
            </ul>
          </div>
  
          <!-- 快捷操作 -->
          <div class="widget-box">
          <div class="widget-title">⚡ 快捷操作</div>
            <div class="quick-actions">
            <button class="action-btn" @click="createVisible = true">
                <el-icon><Plus /></el-icon> 创建群组
            </button>
              <button class="action-btn outline" @click="router.push('/chat?action=add_friend')">
                <el-icon><Search /></el-icon> 查找好友
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建群组弹窗 -->
<el-dialog v-model="createVisible" title="创建新群组" width="400px" align-center>
  <el-form :model="createForm" label-position="top">
    <el-form-item label="群名称">
      <el-input v-model="createForm.name" placeholder="例如：周末开黑小队" />
    </el-form-item>
    <el-form-item label="群公告 (可选)">
      <el-input v-model="createForm.notice" type="textarea" placeholder="介绍一下这个群..." />
    </el-form-item>
    <el-form-item label="群头像 (可选)">
      <input type="file" accept="image/*" @change="onCreateAvatarSelect" />
    </el-form-item>
  </el-form>
  <template #footer>
    <el-button @click="createVisible = false">取消</el-button>
    <el-button type="primary" @click="handleCreateGroup">立即创建</el-button>
  </template>
</el-dialog>
  </template>
  
  <script setup lang="ts">
  import { ref, reactive, onMounted, onUnmounted, computed } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import GlobalNavbar from '@/components/GlobalNavbar.vue'
  import { Connection, MagicStick, ArrowRight, ChatDotRound, User, Plus, Search } from '@element-plus/icons-vue'
  import { useUserStore } from '../../store/user'
  import { useFriendStore } from '../../store/friend'
  import { useChatStore } from '../../store/chat'
  import { askAI } from '../../api/chat'
import { semanticSearch } from '../../api/ai'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createGroup, getMyGroups } from '../../api/group'
import { uploadFile } from '../../api/file'

const createVisible = ref(false)
const createForm = reactive({ name: '', notice: '', avatar: '' })
const createAvatarFile = ref<File | null>(null)
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const friendStore = useFriendStore()
const chatStore = useChatStore()

const quickQuery = ref('')
const aiResponse = ref('')
const recentChats = ref<any[]>([])

onMounted(async () => {
  userStore.fetchUserInfo()
  // 获取好友列表并启动实时更新
  friendStore.startStatusPolling(30000) // 30秒更新一次
  
  // 加载群组和好友作为最近会话
  try {
      const groupsRes: any = await getMyGroups()
      const groups = groupsRes.data || groupsRes || []
      
      // 确保好友列表已加载
      if (friendStore.friends.length === 0) {
          await friendStore.fetchFriends()
      }
        
        // 简单的合并策略：取前2个群 + 前2个好友
        const recentGroups = groups.slice(0, 2).map((g: any) => ({
            ...g,
            isGroup: true,
            display_name: g.name,
            last_msg: '点击进入群聊...',
            updated_at: new Date().toISOString() // 模拟时间，实际应从后端获取
        }))
        
        const recentFriends = friendStore.friends.slice(0, 2).map((f: any) => ({
            ...f,
            isGroup: false,
            display_name: f.nickname || f.username,
            last_msg: f.lastMsg || f.signature || '点击开始聊天...',
            updated_at: new Date().toISOString()
        }))
        
        recentChats.value = [...recentGroups, ...recentFriends]
    } catch (e) {
        console.error('Failed to load chats', e)
    }
  })

  onUnmounted(() => {
    // 停止好友状态轮询
    friendStore.stopStatusPolling()
  })
  
  const handleQuickAsk = async () => {
      if (!quickQuery.value.trim()) return
      
      aiResponse.value = 'AI 正在思考中...'
      try {
          // 默认搜索 ID=1 的群组 (公共大厅)
          // 如果没有群组，则尝试搜索第一个好友
          // TODO: 后续应支持选择搜索范围
          let targetId = 1
          let chatType = 2
          
          const res: any = await semanticSearch({
              query: quickQuery.value,
              target_id: targetId,
              chat_type: chatType
          })
          
          aiResponse.value = res.answer || res.data?.answer || 'AI 未能找到相关答案'
      } catch (e) {
          aiResponse.value = 'AI 服务暂时不可用'
      }
  }
  
  const handleChatClick = (chat: any) => {
      // 设置当前会话并跳转
      chatStore.currentChat = {
          id: chat.ID || chat.id,
          nickname: chat.display_name || chat.name || chat.nickname,
          username: chat.username,
          avatar: chat.avatar,
          isGroup: chat.isGroup
      }
      router.push('/chat')
  }

  const onCreateAvatarSelect = (e: Event) => {
    const file = (e.target as HTMLInputElement).files?.[0] || null
    if (file && !file.type.startsWith('image/')) {
      ElMessage.error('请选择图片文件')
      return
    }
    createAvatarFile.value = file
  }
  
  const handleCreateGroup = async () => {
    if (!createForm.name) {
      ElMessage.warning('请输入群名称')
      return
    }
    try {
      if (createAvatarFile.value) {
        const fd = new FormData()
        fd.append('file', createAvatarFile.value)
        const res: any = await uploadFile(fd)
        createForm.avatar = res.url
      }
      await createGroup(createForm)
      ElMessage.success('群组创建成功！')
      createVisible.value = false
      const { name, notice, avatar } = createForm
      createForm.name = ''
      createForm.notice = ''
      createForm.avatar = ''
      createAvatarFile.value = null
      await ElMessageBox.confirm('创建成功，是否立即进入群聊？', '提示', {
        confirmButtonText: '进入群聊',
        cancelButtonText: '稍后再说',
        type: 'info'
      }).then(() => {
        router.push('/chat')
      }).catch(() => {
        createForm.name = name
        createForm.notice = notice
        createForm.avatar = avatar
      })
    } catch (e) {}
  }
  
  // 格式化时间显示
  const formatTimeAgo = (date: any) => {
    if (!date) return '未知'

    const now = new Date()
    // 确保date是Date对象，如果是字符串则转换为Date
    const dateObj = typeof date === 'string' ? new Date(date) : date

    // 检查转换后的日期是否有效
    if (isNaN(dateObj.getTime())) return '未知'

    const diffMs = now.getTime() - dateObj.getTime()
    const diffSec = Math.floor(diffMs / 1000)
    const diffMin = Math.floor(diffSec / 60)

    if (diffSec < 30) return '刚刚更新'
    if (diffMin < 1) return `${diffSec}秒前`
    if (diffMin < 60) return `${diffMin}分钟前`
    return '实时'
  }
  </script>
  
  <style scoped>
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');
  
  .dashboard-container {
    min-height: 100vh;
    background-color: #f8fafc; /* 极简灰白底色 */
    font-family: 'Inter', sans-serif;
    color: #1e293b;
  }
  
  /* Layout */
  .main-layout {
    max-width: 1200px;
    margin: 0 auto;
    padding: 40px 20px;
    display: grid;
    grid-template-columns: 1fr 320px; /* 左侧自适应，右侧固定 */
    gap: 40px;
  }
  
  /* Welcome Card (Hero) */
  .welcome-card {
    background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
    color: white;
    border-radius: 16px;
    padding: 40px;
    margin-bottom: 40px;
    box-shadow: 0 10px 25px -5px rgba(59, 130, 246, 0.3);
  }
  .welcome-card h1 { margin: 0 0 10px 0; font-size: 28px; }
  .subtitle { opacity: 0.9; margin-bottom: 25px; font-weight: 300; }
  
  .ai-search-bar {
    display: flex;
    background: rgba(255,255,255,0.15);
    padding: 5px;
    border-radius: 12px;
    backdrop-filter: blur(5px);
    border: 1px solid rgba(255,255,255,0.2);
  }
  .ai-search-bar input {
    flex: 1;
    background: transparent;
    border: none;
    color: white;
    padding: 0 15px;
    font-size: 16px;
    outline: none;
  }
  .ai-search-bar input::placeholder { color: rgba(255,255,255,0.6); }
  .ai-search-bar button {
    background: white;
    color: #2563eb;
    border: none;
    padding: 10px 20px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    display: flex; align-items: center; gap: 5px;
    transition: transform 0.2s;
  }
  .ai-search-bar button:hover { transform: scale(1.02); }
  
  .ai-quick-result {
    margin-top: 20px;
    background: rgba(0,0,0,0.2);
    padding: 15px;
    border-radius: 8px;
    border-left: 4px solid #fbbf24;
  }
  .ai-tag { font-size: 12px; font-weight: bold; color: #fbbf24; margin-bottom: 5px; text-transform: uppercase; }
  
  /* Section Header */
  .section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
  .section-header h3 { margin: 0; font-size: 18px; font-weight: 600; }
  
  /* Chat Grid */
  .chat-grid { display: grid; gap: 15px; }
  .chat-card {
    background: white;
    padding: 20px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    gap: 15px;
    cursor: pointer;
    transition: all 0.2s;
    border: 1px solid #f1f5f9;
  }
  .chat-card:hover { transform: translateY(-2px); box-shadow: 0 4px 6px -1px rgba(0,0,0,0.05); border-color: #e2e8f0; }
  .card-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 24px; }
  .group-bg { background: #e0f2fe; color: #0284c7; }
  .private-bg { background: #f3e8ff; color: #9333ea; }
  .card-info { flex: 1; }
  .card-info h4 { margin: 0 0 5px 0; font-size: 16px; }
  .card-info p { margin: 0; font-size: 14px; color: #64748b; }
  .time { font-size: 12px; color: #94a3b8; }
  
  /* Sidebar */
  .sidebar { display: flex; flex-direction: column; gap: 24px; }
  
  .stat-card {
    background: white; padding: 20px; border-radius: 12px; display: flex; justify-content: space-around; border: 1px solid #f1f5f9;
  }
  .stat-item { text-align: center; }
  .stat-item .num {
    font-size: 24px;
    font-weight: 700;
    color: #3b82f6;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
  }
  .update-time {
    font-size: 10px;
    font-weight: 400;
    color: #94a3b8;
    white-space: nowrap;
  }
  .stat-item .label { font-size: 12px; color: #64748b; margin-top: 5px; }
  .divider { width: 1px; background: #e2e8f0; }
  
  .widget-box { background: white; padding: 20px; border-radius: 12px; border: 1px solid #f1f5f9; }
  .widget-title { font-weight: 600; margin-bottom: 15px; font-size: 14px; color: #475569; }
  .notice-list { list-style: none; padding: 0; margin: 0; font-size: 14px; color: #334155; }
  .notice-list li { margin-bottom: 10px; padding-left: 10px; border-left: 2px solid #e2e8f0; }
  
  .quick-actions { display: flex; flex-direction: column; gap: 10px; }
  .action-btn {
    width: 100%; padding: 10px; border-radius: 8px; border: none; background: #3b82f6; color: white; font-weight: 600; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 8px;
  }
  .action-btn.outline { background: white; border: 1px solid #e2e8f0; color: #475569; }
  .action-btn:hover { opacity: 0.9; }
  .action-btn.outline:hover { background: #f8fafc; }
  </style>
