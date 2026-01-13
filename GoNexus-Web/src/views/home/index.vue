<template>
    <div class="dashboard-container">
      
      <!-- 1. 顶部导航 -->
      <header class="navbar">
        <div class="brand">
          <div class="logo-icon"><el-icon><Connection /></el-icon></div>
          <span class="brand-name">GoNexus</span>
        </div>
        <div class="user-menu">
          <el-dropdown>
            <div class="avatar-box">
              <el-avatar :size="36" :src="userStore.userInfo.avatar" />
              <span class="name">{{ userStore.userInfo.nickname || userStore.userInfo.username }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="router.push('/profile')">个人信息</el-dropdown-item>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
  
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
            <!-- 模拟最近联系人，点击直接跳转聊天 -->
            <div class="chat-card" @click="router.push('/chat')">
              <div class="card-icon group-bg"><el-icon><ChatDotRound /></el-icon></div>
              <div class="card-info">
                <h4>公共大厅</h4>
                <p>点击进入沉浸式聊天...</p>
              </div>
              <span class="time">Just now</span>
            </div>
  
            <!-- 这里可以用 v-for 渲染 getPublicGroups 的前几个结果 -->
            <div class="chat-card" v-for="i in 2" :key="i" @click="router.push('/chat')">
              <div class="card-icon private-bg"><el-icon><User /></el-icon></div>
              <div class="card-info">
                <h4>示例好友 {{ i }}</h4>
                <p>[图片] 昨晚的文件发我一下</p>
              </div>
              <span class="time">10 min ago</span>
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
              <button class="action-btn" @click="router.push('/chat')">
                <el-icon><Plus /></el-icon> 创建群组
              </button>
              <button class="action-btn outline" @click="router.push('/chat')">
                <el-icon><Search /></el-icon> 查找好友
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, onUnmounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { Connection, MagicStick, ArrowRight, ChatDotRound, User, Plus, Search } from '@element-plus/icons-vue'
  import { useUserStore } from '../../store/user'
  import { useFriendStore } from '../../store/friend'
  import { askAI } from '../../api/chat'
  import { ElMessage } from 'element-plus'
  
  const router = useRouter()
  const userStore = useUserStore()
  const friendStore = useFriendStore()

  const quickQuery = ref('')
  const aiResponse = ref('')

  onMounted(() => {
    userStore.fetchUserInfo()
    // 获取好友列表并启动实时更新
    friendStore.startStatusPolling(30000) // 30秒更新一次
  })

  onUnmounted(() => {
    // 停止好友状态轮询
    friendStore.stopStatusPolling()
  })
  
  const handleQuickAsk = async () => {
    if (!quickQuery.value) return
    const loadingMsg = ElMessage.info('AI 正在思考中...')
    try {
      const res: any = await askAI(quickQuery.value)
      aiResponse.value = res.answer
      loadingMsg.close()
    } catch (e) {
      ElMessage.error('AI 服务暂时不可用')
    }
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

  const handleLogout = () => {
    localStorage.clear()
    userStore.clearUser()
    friendStore.clearFriends()
    router.push('/login')
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
  
  /* Navbar */
  .navbar {
    height: 64px;
    background: white;
    border-bottom: 1px solid #e2e8f0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 40px;
    position: sticky; top: 0; z-index: 100;
  }
  .brand { display: flex; align-items: center; gap: 10px; }
  .logo-icon { width: 32px; height: 32px; background: #3b82f6; color: white; border-radius: 8px; display: flex; align-items: center; justify-content: center; }
  .brand-name { font-weight: 700; font-size: 20px; letter-spacing: -0.5px; }
  .avatar-box { display: flex; align-items: center; gap: 10px; cursor: pointer; }
  .name { font-weight: 500; font-size: 14px; }
  
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