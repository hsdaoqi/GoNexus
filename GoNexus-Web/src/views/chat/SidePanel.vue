<template>
  <div class="side-panel">
    <!-- 1. 用户信息 -->
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

    <!-- 2. 顶部导航 (返回大厅) -->
    <div class="nav-header">
      <button class="home-btn" @click="$router.push('/')">
        <el-icon><HomeFilled /></el-icon>
        <span>返回大厅</span>
      </button>
    </div>

    <!-- 3. 工具栏 -->
    <div class="tool-bar">
      <div class="search-wrapper">
        <el-icon class="search-icon"><Search /></el-icon>
        <input type="text" placeholder="搜索..." class="search-input">
      </div>
      <button class="add-btn" @click="dialogVisible = true" title="添加好友">
        <el-icon><Plus /></el-icon>
      </button>
    </div>

    <!-- 4. Tab 切换 -->
    <div class="panel-tabs">
      <div :class="['tab-item', { active: currentTab === 'friend' }]" @click="currentTab = 'friend'">
        <el-icon><User /></el-icon> 好友
      </div>
      <div :class="['tab-item', { active: currentTab === 'group' }]" @click="currentTab = 'group'">
        <el-icon><ChatLineSquare /></el-icon> 群组
      </div>
    </div>

    <!-- 5. 列表区域 -->
    <div class="menu-list">
      
      <!-- 通知栏 -->
      <div class="menu-title">NOTIFICATIONS</div>  
      <div class="friend-item system-item" @click="showRequestDialog">
        <div class="avatar-box">
           <el-icon><Bell /></el-icon>
           <div v-if="pendingList.length > 0" class="red-dot">{{ pendingList.length }}</div>
        </div>
        <div class="friend-info">
          <div class="friend-name">系统通知 / System</div>
          <div class="friend-sig">
             {{ pendingList.length > 0 ? `${pendingList.length} 条申请待处理` : '暂无新消息' }}
          </div>
        </div>
      </div>
     
      <!-- 好友列表 -->
      <template v-if="currentTab === 'friend'">
        <div class="menu-title">FRIENDS LIST</div>
        <div v-if="friendStore.friends.length === 0" class="empty-tip">暂无好友连接</div>
        <div 
            v-for="friend in friendStore.friends" 
            :key="friend.id"
            :class="['friend-item', { active: chatStore.currentChat?.id === friend.id && !chatStore.currentChat?.isGroup }]"
            @click="handleSelect(friend)"
        >
            <div style="position: relative; display: flex;">
              <el-avatar :size="36" :src="friend.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
              <div v-if="(chatStore.unreadMap?.['friend_' + friend.id] || 0) > 0" class="red-dot">
                {{ (chatStore.unreadMap?.['friend_' + friend.id] || 0) > 99 ? '99+' : chatStore.unreadMap?.['friend_' + friend.id] }}
              </div>
              <div :class="['online-dot', { online: friend.isOnline }]"></div>
            </div>
            <div class="friend-info">
                <div class="friend-name">{{ friend.nickname || friend.username }}</div>
                <div class="friend-sig text-truncate">{{ friend.signature || 'Link Start!' }}</div>
            </div>
            <div class="action-btn delete-btn" @click.stop="handleDelete(friend)">
                <el-icon><Delete /></el-icon>
            </div>
        </div>
      </template>

      <!-- 群组列表 -->
      <template v-if="currentTab === 'group'">
        <div class="menu-title">MY GROUPS</div>
        <div v-if="groupList.length === 0" class="empty-tip">暂无群组数据</div>
        <div 
           v-for="group in groupList" 
           :key="group.ID" 
           :class="['friend-item', { active: chatStore.currentChat?.id === group.ID && chatStore.currentChat?.isGroup }]"
           @click="handleSelectGroup(group)"
         >
           <div style="position: relative; display: flex;">
             <el-avatar shape="square" :size="36" :src="group.avatar" />
             <div v-if="(chatStore.unreadMap?.['group_' + group.ID] || 0) > 0" class="red-dot">
                {{ (chatStore.unreadMap?.['group_' + group.ID] || 0) > 99 ? '99+' : chatStore.unreadMap?.['group_' + group.ID] }}
             </div>
           </div>
           <div class="friend-info">
            <div class="friend-name">{{ group.name }}</div>
            <div class="friend-sig text-truncate">{{ group.notice || '暂无公告' }}</div>
          </div>
        </div>
     </template>
    </div>

    <!-- 弹窗：添加好友 -->
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

    <!-- 弹窗：好友申请处理 -->
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
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { Search, Plus, HomeFilled, User, Bell, Delete, ChatLineSquare } from '@element-plus/icons-vue'
import { useUserStore } from '../../store/user'
import { useChatStore } from '../../store/chat'
import { useFriendStore } from '../../store/friend'
import { getFriendList, addFriend, getPendingRequests, handleRequest, deleteFriend } from '../../api/friend' 
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMyGroups } from '../../api/group'

const userStore = useUserStore()
const chatStore = useChatStore()
const friendStore = useFriendStore()

// 状态定义
const currentTab = ref('friend')
const groupList = ref<any[]>([])
const pendingList = ref<any[]>([])

// 弹窗状态
const dialogVisible = ref(false)
const requestVisible = ref(false)
const addForm = reactive({ id: '', msg: '' })

// 初始化数据
const initData = async () => {
  try {
    // 启动好友状态轮询 (内部会获取好友列表)
    friendStore.startStatusPolling()

    const [groups, pending] = await Promise.all([
        getMyGroups(),
        getPendingRequests()
    ])
    groupList.value = Array.isArray(groups) ? groups : []
    
    // 同步群组未读数
    groupList.value.forEach((group: any) => {
        if (group.unread_count && group.unread_count > 0) {
            chatStore.unreadMap[`group_${group.ID}`] = group.unread_count
        }
    })

    pendingList.value = Array.isArray(pending) ? pending : []
  } catch (e) {
    console.error("加载列表失败", e)
  }
}

onMounted(() => { initData() })

// 监听 Tab 切换
watch(currentTab, () => { initData() })

// 选中好友
const handleSelect = (friend: any) => {
  chatStore.selectFriend({ ...friend, isGroup: false })
}

// 选中群组
const handleSelectGroup = (group: any) => {
  chatStore.selectFriend({
    id: group.ID,
    nickname: group.name,
    avatar: group.avatar,
    isGroup: true,
    owner_id: group.owner_id,
    notice: group.notice
  })
}

// 其他业务逻辑 (保持不变)
const showRequestDialog = async () => {
  requestVisible.value = true
  pendingList.value = await getPendingRequests() as any
}

const processReq = async (reqId: number, action: number) => {
  try {
    await handleRequest({ request_id: reqId, action })
    ElMessage.success(action === 1 ? '已添加好友' : '已拒绝')
    initData()
    const newList = pendingList.value.filter((item: any) => item.ID !== reqId)
    pendingList.value = newList
    if (newList.length === 0) requestVisible.value = false
  } catch (e) {}
}

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

const handleDelete = (friend: any) => {
  ElMessageBox.confirm(
    `确定断开与 [${friend.nickname || friend.username}] 的连接？`, 'System Alert',
    { confirmButtonText: '断开', cancelButtonText: '取消', type: 'warning', customClass: 'sao-message-box' }
  ).then(async () => {
    try {
      await deleteFriend({ friend_id: friend.id })
      ElMessage.success('连接已断开')
      if (chatStore.currentChat?.id === friend.id) {
        chatStore.currentChat = null
        chatStore.messages = []
      }
      initData()
    } catch (e) {}
  }).catch(() => {})
}

// 暴露方法
defineExpose({ refreshPendingList: initData, initData })
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700&display=swap');

.side-panel {
  width: 100%; height: 100%; background: #2c3e50; color: white; display: flex; flex-direction: column; box-shadow: 2px 0 10px rgba(0,0,0,0.1);
}

.user-profile { padding: 20px; display: flex; align-items: center; gap: 12px; border-bottom: 1px solid rgba(255,255,255,0.1); }
.username { font-weight: bold; font-size: 16px; }
.status-dot { display:inline-block; width:8px; height:8px; background:#2ecc71; border-radius:50%; margin-right:5px;}
.status-text { font-size: 12px; color: #2ecc71; }

.nav-header { padding: 10px; }
.home-btn { width: 100%; background: #34495e; border: none; color: #fff; padding: 6px; border-radius: 6px; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 8px; font-family: 'Orbitron'; font-size: 12px; transition: background 0.3s; }
.home-btn:hover { background: #4facfe; }

.tool-bar { padding: 10px; display: flex; gap: 8px; border-bottom: 1px solid rgba(255,255,255,0.1); }
.search-wrapper { flex: 1; background: rgba(0,0,0,0.2); border-radius: 4px; display: flex; align-items: center; padding: 0 5px; }
.search-input { background: transparent; border: none; color: white; width: 100%; font-size: 12px; outline: none;}
.add-btn { background: #4facfe; border: none; color: white; width: 28px; height: 28px; border-radius: 4px; cursor: pointer; display: flex; align-items: center; justify-content: center; transition: transform 0.2s; }
.add-btn:hover { transform: scale(1.1); }

.panel-tabs { display: flex; margin: 10px; background: rgba(0,0,0,0.2); padding: 4px; border-radius: 4px; }
.tab-item { flex: 1; text-align: center; font-size: 12px; padding: 6px; cursor: pointer; color: #bdc3c7; border-radius: 4px; display: flex; align-items: center; justify-content: center; gap: 5px; transition: all 0.3s; }
.tab-item.active { background: #4facfe; color: white; font-weight: bold; }

.menu-list { flex: 1; overflow-y: auto; padding: 10px; }
.menu-title { font-size: 10px; color: #7f8c8d; margin-bottom: 10px; font-family: 'Orbitron'; letter-spacing: 1px; }
.empty-tip { text-align: center; color: #666; font-size: 12px; padding: 20px; }

/* 列表项 */
.friend-item { display: flex; align-items: center; padding: 8px; border-radius: 6px; cursor: pointer; margin-bottom: 2px; color: #bdc3c7; position: relative; transition: background 0.2s; }
.friend-item:hover { background: rgba(255,255,255,0.1); color: white; }
.friend-item.active { background: #4facfe; color: white; }
.friend-info { margin-left: 10px; overflow: hidden; flex: 1; }
.friend-name { font-size: 14px; font-weight: 500; }
.friend-sig { font-size: 11px; opacity: 0.7; }

/* 悬浮操作按钮 */
.action-btn { position: absolute; right: 10px; top: 50%; transform: translateY(-50%); display: none; padding: 4px; border-radius: 4px; font-size: 14px; }
.friend-item:hover .action-btn { display: block; }
.delete-btn { color: #e74c3c; background: rgba(255,255,255,0.2); }
.delete-btn:hover { background: #e74c3c; color: white; }
.setting-btn { color: #f1c40f; background: rgba(255,255,255,0.2); }
.setting-btn:hover { background: #f1c40f; color: white; }

/* 系统通知项 */
.system-item { background: rgba(255, 165, 0, 0.1); border-left: 3px solid orange; }
.avatar-box { width: 36px; height: 36px; background: orange; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: white; position: relative; }
.red-dot { position: absolute; top: -2px; right: -2px; background: red; color: white; font-size: 10px; width: 16px; height: 16px; border-radius: 50%; display: flex; align-items: center; justify-content: center; }

/* 在线状态点 */
.online-dot { position: absolute; bottom: 0; right: 0; width: 10px; height: 10px; background: #95a5a6; border-radius: 50%; border: 2px solid #2c3e50; transition: background-color 0.3s; }
.online-dot.online { background: #2ecc71; }

/* 通用按钮 */
.sao-btn-confirm { background: #4facfe; color: white; border: none; padding: 8px 25px; border-radius: 4px; font-family: 'Orbitron'; cursor: pointer; font-weight: bold; transition: background 0.3s; }
.sao-btn-confirm:hover { background: #2980b9; }
.sao-btn-cancel { background: #95a5a6; color: white; border: none; padding: 8px 25px; border-radius: 4px; font-family: 'Orbitron'; cursor: pointer; }

/* 弹窗覆盖 */
:deep(.sao-dialog) { background: rgba(255, 255, 255, 0.95); border-radius: 8px; border: 1px solid #ff9966; }
:deep(.el-dialog__header) { background: #ff9966; padding: 10px 20px; margin-right: 0; }
:deep(.el-dialog__body) { padding: 30px; text-align: center; }
.dialog-icon { font-size: 40px; color: #ff9966; margin-bottom: 10px; }
.sao-input-orange { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; outline: none; transition: border 0.3s; }
.sao-input-orange:focus { border-color: #ff9966; }
.request-list { max-height: 300px; overflow-y: auto; }
.request-item { display: flex; align-items: center; padding: 10px; border-bottom: 1px solid #eee; gap: 10px; }
.req-info { flex: 1; text-align: left; }
.req-name { font-weight: bold; font-size: 14px; color: #333; }
.req-msg { font-size: 12px; color: #666; }
.req-actions { display: flex; gap: 5px; }
.sao-btn-mini { border: none; padding: 4px 8px; border-radius: 4px; cursor: pointer; font-size: 10px; font-family: 'Orbitron'; color: white; }
.sao-btn-mini.accept { background: #2ecc71; }
.sao-btn-mini.reject { background: #e74c3c; }
</style>
