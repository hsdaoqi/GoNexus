<template>
<div class="side-panel">
    <!-- 用户信息 -->
    <div class="user-profile">
    <el-avatar :size="40" :src="userStore.userInfo.avatar" class="my-avatar" />
    <div class="user-info">
        <span class="username">{{ userStore.userInfo.nickname || userStore.userInfo.username }}</span>
        <div class="status-row">
        <span class="status-dot"></span>
        <span class="status-text">Link Start</span>
        </div>
    </div>
    </div>

    <!-- 顶部导航 (返回大厅) -->
    <div class="nav-header">
        <button class="home-btn" @click="$router.push('/')">
        <el-icon><HomeFilled /></el-icon>
        <span>大厅</span>
        </button>
    </div>

    <!-- 工具栏 -->
    <div class="tool-bar">
    <div class="search-wrapper">
        <el-icon class="search-icon"><Search /></el-icon>
        <input type="text" placeholder="搜索好友..." class="search-input">
    </div>
    <button class="add-btn" @click="dialogVisible = true">
        <el-icon><Plus /></el-icon>
    </button>
    </div>


    <div class="menu-list">
        <!-- 🔥新的朋友入口 -->
    <div class="menu-title">NOTIFICATIONS</div>  
      <div class="friend-item system-item" @click="showRequestDialog">
        <div class="avatar-box">
           <el-icon><Bell /></el-icon>
           <!-- 小红点 -->
           <div v-if="pendingList.length > 0" class="red-dot">{{ pendingList.length }}</div>
        </div>
        <div class="friend-info">
          <div class="friend-name">新的朋友 / New Friends</div>
          <div class="friend-sig">
             {{ pendingList.length > 0 ? `${pendingList.length} 条申请待处理` : '暂无新消息' }}
          </div>
        </div>
      </div>
      <!-- 好友列表 -->
    <div class="menu-title">FRIENDS LIST</div>
    <div 
        v-for="friend in friendList" 
        :key="friend.id"
        :class="['friend-item', { active: chatStore.currentChat?.id === friend.id }]"
        @click="handleSelect(friend)"
    >
        <el-avatar :size="36" :src="friend.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
        <div class="friend-info">
            <div class="friend-name">{{ friend.nickname || friend.username }}</div>
            <div class="friend-sig text-truncate">{{ friend.signature || 'Link Start!' }}</div>
        </div>

        <!-- 🔥 新增：删除按钮 (使用 .stop 阻止冒泡，防止触发选中) -->
         <div class="delete-btn" @click.stop="handleDelete(friend)">
            <el-icon><Delete /></el-icon>
         </div>

    </div>
    </div>

    <!-- 添加好友弹窗 (直接内嵌在组件里) -->
    <el-dialog v-model="dialogVisible" title="System Alert" width="400px" class="sao-dialog" :show-close="false" align-center>
        <div class="dialog-content">
        <div class="dialog-icon"><el-icon><User /></el-icon></div>
        <h3>添加好友 / Add Friend</h3>
        <input v-model="addForm.id" class="sao-input-orange" placeholder="Target ID..." />
        <input v-model="addForm.msg" class="sao-input-orange" placeholder="Verification Message..." style="margin-top:10px" />
        </div>
        <template #footer>
        <div class="dialog-footer">
            <button class="sao-btn-cancel" @click="dialogVisible = false">CANCEL</button>
            <button class="sao-btn-confirm" @click="handleAddFriend">SEND</button>
        </div>
        </template>
    </el-dialog>

     <!-- 🔥 新增：处理好友申请弹窗 -->
     <el-dialog v-model="requestVisible" title="System Notifications" width="500px" class="sao-dialog" :show-close="false">
      <div class="dialog-content request-list">
        <div v-if="pendingList.length === 0" class="empty-tip">暂无新的申请</div>
        
        <div v-for="req in pendingList" :key="req.id" class="request-item">
          <el-avatar :size="40" :src="req.requester_avatar" />
          <div class="req-info">
            <div class="req-name">{{ req.requester_name }}</div>
            <div class="req-msg">留言: {{ req.verify_msg }}</div>
          </div>
          <div class="req-actions">
            <!-- Action 1: 同意, 2: 拒绝 -->
            <button class="sao-btn-mini accept" @click="processReq(req.id, 1)">✔ ACC</button>
            <button class="sao-btn-mini reject" @click="processReq(req.id, 2)">✘ REJ</button>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="sao-btn-cancel" @click="requestVisible = false">CLOSE</button>
      </template>
    </el-dialog>
</div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Plus, HomeFilled, User, Bell, Delete } from '@element-plus/icons-vue' // 引入 Bell 图标
import { useUserStore } from '../../store/user'
import { useChatStore } from '../../store/chat'
import { getFriendList, addFriend,getPendingRequests, handleRequest,deleteFriend } from '../../api/friend' 
import { ElMessage ,ElMessageBox} from 'element-plus'

const userStore = useUserStore()
const chatStore = useChatStore()

const friendList = ref<any[]>([])
const dialogVisible = ref(false)
const addForm = reactive({ id: '', msg: '' })
const pendingList = ref<any[]>([]) // 待处理列表
const requestVisible = ref(false)  // 申请弹窗控制

// 加载数据
const initData = async () => {
  try {
    friendList.value = await getFriendList() as any
    pendingList.value = await getPendingRequests() as any // 获取申请
  } catch(e){
    console.error('加载数据失败:', e)
  }
}

// 显示申请列表弹窗
const showRequestDialog = async () => {
  requestVisible.value = true
  // 打开时刷新一下
  pendingList.value = await getPendingRequests() as any
}

// 处理申请 (同意/拒绝)
const processReq = async (reqId: number, action: number) => {
  try {
    await handleRequest({ request_id: reqId, action })
    ElMessage.success(action === 1 ? '已添加好友' : '已拒绝')
    
    // 刷新列表
    initData()
    
    // 如果列表空了，自动关窗
    const newList = pendingList.value.filter((item: any) => item.ID !== reqId)
    pendingList.value = newList
    if (newList.length === 0) requestVisible.value = false
    
  } catch (e) {
    // error handled
  }
}

onMounted(() => {
  initData()
})

// 选中好友 -> 调用 Store Actions
const handleSelect = (friend: any) => {
chatStore.selectFriend(friend)
}

// 添加好友逻辑
const handleAddFriend = async () => {
if (!addForm.id) return
try {
    await addFriend({ receiver_id: parseInt(addForm.id), verify_msg: addForm.msg })
    ElMessage.success('申请发送成功')
    dialogVisible.value = false
    addForm.id = ''
    addForm.msg = ''
} catch(e) {}
}

// 添加刷新待处理申请列表的方法
const refreshPendingList = async () => {
  pendingList.value = await getPendingRequests() as any
}

// 删除好友逻辑
const handleDelete = (friend: any) => {
  ElMessageBox.confirm(
    `确定要解除与 [${friend.nickname || friend.username}] 的连接吗？`,
    'System Alert',
    {
      confirmButtonText: '断开连接',
      cancelButtonText: '取消',
      type: 'warning',
      customClass: 'sao-message-box' // 我们可以自定义样式
    }
  ).then(async () => {
    try {
      await deleteFriend({ friend_id: friend.id })
      ElMessage.success('连接已断开')
      
      // 1. 如果当前正打开着这个人的聊天窗口，清空它
      if (chatStore.currentChat?.id === friend.id) {
        chatStore.currentChat = null
        chatStore.messages = []
      }
      
      // 2. 刷新列表
      await initData()
    } catch (e) {
      // error
    }
  }).catch(() => {})
}



// 暴露方法给父组件
defineExpose({
  refreshPendingList,
  initData
})
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700&display=swap');

.side-panel {
/* 宽度由父容器控制，这里占满父容器 */
width: 100%; 
height: 100%;
background: #2c3e50; 
color: white;
display: flex;
flex-direction: column;
box-shadow: 2px 0 10px rgba(0,0,0,0.1);
}

.user-profile {
padding: 20px;
display: flex; align-items: center; gap: 12px;
border-bottom: 1px solid rgba(255,255,255,0.1);
}
.username { font-weight: bold; font-size: 16px; }
.status-text { font-size: 12px; color: #2ecc71; }
.status-dot { display:inline-block; width:8px; height:8px; background:#2ecc71; border-radius:50%; margin-right:5px;}

.tool-bar {
padding: 10px; display: flex; gap: 8px; border-bottom: 1px solid rgba(255,255,255,0.1);
}
.search-wrapper {
flex: 1; background: rgba(0,0,0,0.2); border-radius: 4px; display: flex; align-items: center; padding: 0 5px;
}
.search-input { background: transparent; border: none; color: white; width: 100%; font-size: 12px; outline: none;}
.add-btn { background: #4facfe; border: none; color: white; width: 28px; height: 28px; border-radius: 4px; cursor: pointer; display: flex; align-items: center; justify-content: center; }

.nav-header { padding: 10px; }
.home-btn {
width: 100%; background: #34495e; border: none; color: #fff; padding: 6px; border-radius: 6px; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 8px; font-family: 'Orbitron'; font-size: 12px;
}
.home-btn:hover { background: #4facfe; }

.menu-list { flex: 1; overflow-y: auto; padding: 10px; }
.menu-title { font-size: 10px; color: #7f8c8d; margin-bottom: 10px; font-family: 'Orbitron'; }

.friend-item {
display: flex; align-items: center; padding: 8px; border-radius: 6px; cursor: pointer; margin-bottom: 2px; color: #bdc3c7;
}
.friend-item:hover { background: rgba(255,255,255,0.1); color: white; }
.friend-item.active { background: #4facfe; color: white; }
.friend-info { margin-left: 10px; overflow: hidden; }
.friend-name { font-size: 14px; font-weight: 500; }
.friend-sig { font-size: 11px; opacity: 0.7; }

/* 弹窗样式 */
:deep(.sao-dialog) { background: rgba(255, 255, 255, 0.95); border-radius: 8px; border: 1px solid #ff9966; }
:deep(.el-dialog__header) { background: #ff9966; padding: 10px 20px; margin-right: 0; }
:deep(.el-dialog__body) { padding: 30px; text-align: center; }
.dialog-icon { font-size: 40px; color: #ff9966; margin-bottom: 10px; }
.sao-input-orange { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; outline: none; transition: border 0.3s; }
.sao-input-orange:focus { border-color: #ff9966; }
.dialog-footer { display: flex; justify-content: center; gap: 20px; padding-bottom: 20px; }
.sao-btn-confirm { background: #ff9966; color: white; border: none; padding: 8px 25px; border-radius: 20px; font-family: 'Orbitron'; cursor: pointer; }
.sao-btn-cancel { background: #999; color: white; border: none; padding: 8px 25px; border-radius: 20px; font-family: 'Orbitron'; cursor: pointer; }


/* 新的朋友入口样式 */
.system-item {
  background: rgba(255, 165, 0, 0.1);
  border-left: 3px solid orange;
}
.avatar-box {
  width: 36px; height: 36px; background: orange; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: white; position: relative;
}
.red-dot {
  position: absolute; top: -2px; right: -2px; background: red; color: white; font-size: 10px; width: 16px; height: 16px; border-radius: 50%; display: flex; align-items: center; justify-content: center;
}

/* 处理申请弹窗样式 */
.request-list { max-height: 300px; overflow-y: auto; }
.request-item {
  display: flex; align-items: center; padding: 10px; border-bottom: 1px solid #eee; gap: 10px;
}
.req-info { flex: 1; text-align: left; }
.req-name { font-weight: bold; font-size: 14px; color: #333; }
.req-msg { font-size: 12px; color: #666; }
.req-actions { display: flex; gap: 5px; }

.sao-btn-mini { border: none; padding: 4px 8px; border-radius: 4px; cursor: pointer; font-size: 10px; font-family: 'Orbitron'; color: white; }
.sao-btn-mini.accept { background: #2ecc71; }
.sao-btn-mini.reject { background: #e74c3c; }
.empty-tip { text-align: center; color: #999; padding: 20px; }


/* 🔥 好友列表项样式调整 */
.friend-item {
  display: flex; 
  align-items: center; 
  padding: 8px; 
  border-radius: 6px; 
  cursor: pointer; 
  margin-bottom: 2px;
  color: #bdc3c7;
  position: relative; /* 为了定位删除按钮 */
}

.friend-item:hover { 
  background: rgba(255,255,255,0.1); 
  color: white; 
}

/* 选中状态 */
.friend-item.active { 
  background: #4facfe; 
  color: white; 
}

/* 🔥 删除按钮样式 */
.delete-btn {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  display: none; /* 默认隐藏 */
  padding: 4px;
  border-radius: 4px;
  color: #e74c3c; /* 红色 */
  background: rgba(255, 255, 255, 0.2);
}

.delete-btn:hover {
  background: #e74c3c;
  color: white;
}

/* 只有当鼠标悬停在 item 上时，才显示删除按钮 */
.friend-item:hover .delete-btn {
  display: block;
}

/* 如果当前项被选中，删除按钮颜色改为白色以适配背景 */
.friend-item.active .delete-btn {
  color: white;
  background: rgba(0,0,0,0.2);
}
.friend-item.active .delete-btn:hover {
  background: #e74c3c;
}

/* ... 其他 CSS 保持不变 ... */
</style>

<!-- 全局样式覆盖 Element Plus MessageBox (可选，为了更好看) -->
<style>
.sao-message-box {
  border-radius: 8px !important;
  border: 1px solid #ff9966 !important;
  font-family: 'Segoe UI', sans-serif;
}
.sao-message-box .el-button--primary {
  background-color: #ff9966 !important;
  border-color: #ff9966 !important;
}
</style>