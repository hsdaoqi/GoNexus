<template>
  <div class="group-panel">
    <!-- 1. 顶部群信息 (头像 + 名称) -->
    <div class="group-header">
       <el-upload
          class="avatar-uploader"
          action="#"
          :show-file-list="false"
          :auto-upload="false"
          :on-change="onAvatarChange"
          :disabled="!isOwner"
        >
          <div class="avatar-box">
              <el-avatar :size="48" :src="chatStore.currentChat?.avatar" shape="square" />
              <div v-if="isOwner" class="upload-mask">
                  <el-icon><Camera /></el-icon>
              </div>
          </div>
        </el-upload>
        <div class="group-basic-info">
           <div class="group-name text-truncate" :title="chatStore.currentChat?.nickname">
             {{ chatStore.currentChat?.nickname }}
           </div>
           <div class="group-id">ID: {{ chatStore.currentChat?.id }}</div>
        </div>
    </div>

    <!-- 2. 群公告 -->
    <div class="section-group">
      <div class="section-header">
        <span class="label">群公告</span>
        <el-button v-if="isOwner" type="primary" link size="small" @click="startEditNotice">编辑</el-button>
      </div>
      <div class="notice-container">
          <div v-if="isEditingNotice" class="edit-box">
            <el-input 
              v-model="editNoticeContent" 
              type="textarea" 
              :rows="3" 
              placeholder="请输入群公告"
              resize="none"
              class="notice-input"
            />
            <div class="edit-actions">
              <el-button type="primary" size="small" link @click="saveNotice">保存</el-button>
              <el-button type="info" size="small" link @click="cancelEditNotice">取消</el-button>
            </div>
          </div>
          <div 
              v-else 
              class="notice-content text-truncate-multiline"
              :class="{ 'is-owner': isOwner }"
              @click="isOwner && startEditNotice()"
              :title="isOwner ? '点击编辑公告' : ''"
          >
              {{ chatStore.currentChat?.notice || '暂无公告' }}
          </div>
      </div>
    </div>

    <!-- 3. 群成员列表 (仿QQ) -->
    <div class="section-group flex-1">
      <div class="section-header">
        <span class="label">群聊成员 {{ groupMembers.length }}</span>
        <div class="header-actions">
           <el-icon class="action-icon" title="搜索成员" @click="showSearch = !showSearch"><Search /></el-icon>
           <el-icon class="action-icon" title="邀请好友" @click="inviteDialogVisible = true"><Plus /></el-icon>
        </div>
      </div>
      
      <!-- 搜索框 -->
      <div v-if="showSearch" class="search-box">
          <el-input v-model="searchKeyword" placeholder="搜索群成员" size="small" prefix-icon="Search" clearable />
      </div>

      <!-- 成员列表 (垂直列表) -->
      <div class="member-list">
         <template v-for="m in filteredMembers" :key="m.ID">
            <!-- 只有群主或管理员能操作他人，且不能操作自己 -->
            <el-dropdown 
                v-if="(isOwner || m.role < 2) && m.user_id !== userStore.userInfo.id && isOwner" 
                trigger="contextmenu" 
                placement="bottom-start"
                @command="(cmd: string) => handleMemberAction(cmd, m)"
            >
                <div class="member-row">
                    <el-avatar :size="32" :src="m.user_avatar" />
                    <div class="member-info">
                        <div class="member-name text-truncate">
                            {{ m.user_name || m.nickname }}
                            <span v-if="m.muted" style="color: #f56c6c; font-size: 12px; transform: scale(0.8); display: inline-block;">(禁言)</span>
                        </div>
                    </div>
                    <div v-if="m.role === 3" class="role-tag owner">群主</div>
                    <div v-else-if="m.role === 2" class="role-tag admin">管理员</div>
                </div>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item command="mute">{{ m.muted ? '解除禁言' : '禁言' }}</el-dropdown-item>
                        <el-dropdown-item command="kick" style="color: #f56c6c;">移除成员</el-dropdown-item>
                        <el-dropdown-item v-if="isOwner" command="setAdmin">{{ m.role === 2 ? '取消管理员' : '设为管理员' }}</el-dropdown-item>
                        <el-dropdown-item v-if="isOwner" command="transferOwner" style="color: #e6a23c;">转让群主</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>

            <!-- 普通展示 (自己，或无权操作的人) -->
            <div v-else class="member-row">
                <el-avatar :size="32" :src="m.user_avatar" />
                <div class="member-info">
                    <div class="member-name text-truncate">
                        {{ m.user_name || m.nickname }}
                        <span v-if="m.muted" style="color: #f56c6c; font-size: 12px; transform: scale(0.8); display: inline-block;">(禁言)</span>
                    </div>
                </div>
                <div v-if="m.role === 3" class="role-tag owner">群主</div>
                <div v-else-if="m.role === 2" class="role-tag admin">管理员</div>
            </div>
         </template>
      </div>
    </div>
    
    <!-- 4. 底部操作 -->
    <div class="panel-footer">
        <el-button v-if="isOwner" type="danger" plain class="full-width" size="default">解散群聊</el-button>
        <el-button v-else type="danger" plain class="full-width" size="default">退出群聊</el-button>
    </div>

    <!-- 邀请好友弹窗 -->
    <el-dialog v-model="inviteDialogVisible" title="邀请好友" width="320px" class="custom-dialog" align-center append-to-body>
        <div class="invite-list">
          <div v-if="availableFriends.length === 0" class="empty-tip">暂无好友可邀请</div>
          <div v-for="f in availableFriends" :key="f.id" class="friend-item invite-item" @click="handleInviteFriend(f)">
             <el-avatar :size="32" :src="f.avatar" />
             <div class="friend-info">
               <div class="friend-name">{{ f.nickname }}</div>
             </div>
             <div class="action-btn"><el-icon><Plus /></el-icon></div>
          </div>
        </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Plus, Search, Camera, Edit } from '@element-plus/icons-vue'
import { useChatStore } from '../../store/chat'
import { useUserStore } from '../../store/user'
import { getGroupMembers, inviteMember, updateGroup, kickMember, muteMember, setAdmin, transferGroup } from '../../api/group'
import { getFriendList } from '../../api/friend'
import { uploadFile } from '../../api/file'
import { ElMessage, ElMessageBox } from 'element-plus'

const chatStore = useChatStore()
const userStore = useUserStore()

const groupMembers = ref<any[]>([])
const searchKeyword = ref('')
const showSearch = ref(false)
const inviteDialogVisible = ref(false)
const friendList = ref<any[]>([])

// 公告编辑状态
const isEditingNotice = ref(false)
const editNoticeContent = ref('')

// 是否是群主
const isOwner = computed(() => {
    if (!chatStore.currentChat || !userStore.userInfo) return false
    // 转换为字符串比较，防止类型不一致 (number vs string)
    return String(chatStore.currentChat.owner_id) === String(userStore.userInfo.id)
})

// 过滤后的成员列表
const filteredMembers = computed(() => {
    if (!searchKeyword.value) return groupMembers.value
    return groupMembers.value.filter(m => {
        const name = m.user_name || m.nickname || ''
        return name.toLowerCase().includes(searchKeyword.value.toLowerCase())
    })
})

// 可邀请的好友 (不在群内的好友)
const availableFriends = computed(() => {
  const memberIds = new Set(groupMembers.value.map(m => m.user_id))
  return friendList.value.filter(f => !memberIds.has(f.id))
})

// 获取数据
const fetchData = async () => {
    if (!chatStore.currentChat || !chatStore.currentChat.isGroup) return
    try {
        const [members, friends] = await Promise.all([
            getGroupMembers(chatStore.currentChat.id),
            getFriendList()
        ])
        groupMembers.value = Array.isArray(members) ? members : (members as any).data || []
        friendList.value = Array.isArray(friends) ? friends : []
    } catch (e) {
        console.error("加载群信息失败", e)
    }
}

// 监听当前聊天变化
watch(() => chatStore.currentChat?.id, () => {
    fetchData()
    // 重置状态
    isEditingNotice.value = false
    showSearch.value = false
    searchKeyword.value = ''
}, { immediate: true })

// 头像上传
const onAvatarChange = async (file: any) => {
    if (!chatStore.currentChat) return
    const formData = new FormData()
    formData.append('file', file.raw)
    
    try {
        const res: any = await uploadFile(formData)
        // 根据后端返回结构调整
        const avatarUrl = res.url || res.data?.url || res // 假设后端直接返回 url 或 {url: ...}
        
        await updateGroup({
            id: chatStore.currentChat.id,
            avatar: avatarUrl
        })
        chatStore.currentChat.avatar = avatarUrl
        ElMessage.success('群头像更新成功')
    } catch (e) {
        ElMessage.error('头像上传失败')
        console.error(e)
    }
}

// 编辑公告逻辑
const startEditNotice = () => {
    editNoticeContent.value = chatStore.currentChat?.notice || ''
    isEditingNotice.value = true
}

const cancelEditNotice = () => {
    isEditingNotice.value = false
}

const saveNotice = async () => {
    if (!chatStore.currentChat) return
    try {
        await updateGroup({
            id: chatStore.currentChat.id,
            notice: editNoticeContent.value
        })
        chatStore.currentChat.notice = editNoticeContent.value
        isEditingNotice.value = false
        ElMessage.success('公告已更新')
    } catch (e) {
        ElMessage.error('更新失败')
    }
}

// 邀请逻辑
const handleInviteFriend = async (friend: any) => {
  if (!chatStore.currentChat) return
  try {
    await inviteMember({ group_id: chatStore.currentChat.id, friend_id: friend.id })
    ElMessage.success(`已邀请 ${friend.nickname}`)
    inviteDialogVisible.value = false
    // 刷新成员列表
    const members = await getGroupMembers(chatStore.currentChat.id)
    groupMembers.value = (Array.isArray(members) ? members : (members as any).data || []) || []
  } catch (e: any) {
    ElMessage.error(e.message || '邀请失败')
  }
}

// 成员管理操作
const handleMemberAction = async (command: string, member: any) => {
    if (!chatStore.currentChat) return
    const groupId = chatStore.currentChat.id
    
    if (command === 'kick') {
        try {
            await ElMessageBox.confirm(`确定要移除成员 ${member.user_name || member.nickname} 吗？`, '提示', {
                type: 'warning'
            })
            await kickMember({ group_id: groupId, member_id: member.user_id })
            ElMessage.success('已移除该成员')
            fetchData() // 刷新列表
        } catch (e) {
            if (e !== 'cancel') ElMessage.error('操作失败')
        }
    } else if (command === 'mute') {
        try {
            const newStatus = member.muted ? 0 : 1
            await muteMember({ group_id: groupId, member_id: member.user_id, mute: newStatus })
            ElMessage.success(newStatus ? '已禁言' : '已解除禁言')
            fetchData()
        } catch (e) {
            ElMessage.error('操作失败')
        }
    } else if (command === 'setAdmin') {
        try {
            const newIsAdmin = member.role !== 2
            await setAdmin({ group_id: groupId, member_id: member.user_id, is_admin: newIsAdmin })
            ElMessage.success(newIsAdmin ? '已设为管理员' : '已取消管理员')
            fetchData()
        } catch (e) {
            ElMessage.error('操作失败')
        }
    } else if (command === 'transferOwner') {
        try {
            await ElMessageBox.confirm(`确定要将群主转让给 ${member.user_name || member.nickname} 吗？转让后您将变为普通成员。`, '危险操作', {
                type: 'warning',
                confirmButtonText: '确认转让',
                cancelButtonText: '取消'
            })
            await transferGroup({ group_id: groupId, member_id: member.user_id })
            ElMessage.success('群主转让成功')
            fetchData()
        } catch (e) {
            if (e !== 'cancel') ElMessage.error('操作失败')
        }
    }
}
</script>

<style scoped>
.group-panel {
    height: 100%;
    background: #fff;
    border-left: 1px solid #e7e7e7;
    display: flex;
    flex-direction: column;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

/* 1. 顶部群信息 */
.group-header {
    padding: 20px 15px;
    display: flex;
    align-items: center;
    gap: 12px;
    border-bottom: 1px solid #f0f0f0;
}
.avatar-box {
    position: relative;
    cursor: pointer;
}
.upload-mask {
    position: absolute; top: 0; left: 0; width: 100%; height: 100%;
    background: rgba(0,0,0,0.5); border-radius: 4px;
    display: flex; align-items: center; justify-content: center;
    color: white; opacity: 0; transition: opacity 0.2s;
}
.avatar-box:hover .upload-mask { opacity: 1; }

.group-basic-info {
    flex: 1;
    overflow: hidden;
}
.group-name {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    margin-bottom: 4px;
}
.group-id {
    font-size: 12px;
    color: #999;
}

/* 2. 通用区块 */
.section-group {
    padding: 15px;
    border-bottom: 1px solid #f0f0f0;
}
.section-group.flex-1 {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    padding-bottom: 0;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
    height: 24px;
}
.section-header .label {
    font-size: 13px;
    color: #666;
}
.action-icon {
    cursor: pointer;
    color: #888;
    font-size: 16px;
    transition: color 0.2s;
}
.action-icon:hover { color: #4facfe; }
.header-actions {
    display: flex;
    gap: 10px;
}

/* 公告内容 */
.notice-content {
    font-size: 13px;
    color: #333;
    line-height: 1.5;
    background: #f9f9f9;
    padding: 8px;
    border-radius: 4px;
    min-height: 40px;
}
.notice-content.is-owner {
    cursor: pointer;
    transition: background-color 0.2s;
}
.notice-content.is-owner:hover {
    background-color: #eee;
}
.edit-box {
    display: flex;
    flex-direction: column;
    gap: 5px;
}
.edit-actions {
    display: flex;
    justify-content: flex-end;
}

/* 3. 成员列表 */
.search-box {
    margin-bottom: 10px;
}
.member-list {
    flex: 1;
    overflow-y: auto;
    /* 隐藏滚动条但保留功能 */
    scrollbar-width: thin;
}
.member-row {
    display: flex;
    align-items: center;
    padding: 8px 5px;
    gap: 10px;
    cursor: pointer;
    border-radius: 4px;
    transition: background 0.2s;
}
.member-row:hover {
    background: #f5f5f5;
}
.member-info {
    flex: 1;
    overflow: hidden;
}
.member-name {
    font-size: 14px;
    color: #333;
}
.role-tag {
    font-size: 10px;
    padding: 1px 4px;
    border-radius: 2px;
    color: white;
    white-space: nowrap;
}
.role-tag.owner { background: #f1c40f; }
.role-tag.admin { background: #2ecc71; }

/* 4. 底部 */
.panel-footer {
    padding: 15px;
    border-top: 1px solid #f0f0f0;
}

/* 工具类 */
.text-truncate {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.text-truncate-multiline {
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.full-width { width: 100%; }

/* 邀请弹窗 */
.invite-list { max-height: 300px; overflow-y: auto; }
.invite-item { display: flex; align-items: center; padding: 10px; border-bottom: 1px solid #f0f0f0; cursor: pointer; transition: background 0.2s; }
.invite-item:hover { background: #f5f7fa; }
.invite-item .friend-name { margin-left: 10px; font-size: 14px; flex: 1; }
.invite-item .action-btn { color: #4facfe; }
.empty-tip { text-align: center; color: #999; padding: 20px; font-size: 13px; }
</style>
