<template>
  <header class="navbar">
    <div class="brand">
      <div class="logo-icon"><el-icon><Connection /></el-icon></div>
      <span class="brand-name">GoNexus</span>
    </div>

    <div class="nav-menu">
      <el-button link :class="{ active: route.path === '/' }" @click="router.push('/')">首页</el-button>
      <el-button link :class="{ active: route.path === '/moments' }" @click="router.push('/moments')">动态广场</el-button>
      <el-button link :class="{ active: route.path === '/chat' }" @click="router.push('/chat')">消息</el-button>
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
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { Connection } from '@element-plus/icons-vue'
import { useUserStore } from '../store/user'
import { useFriendStore } from '../store/friend'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const friendStore = useFriendStore()

const handleLogout = () => {
  localStorage.clear()
  userStore.clearUser()
  friendStore.clearFriends()
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  height: 64px;
  background: white;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.02);
  position: sticky;
  top: 0;
  z-index: 100;
}

.brand { display: flex; align-items: center; gap: 10px; }
.logo-icon { width: 32px; height: 32px; background: #3b82f6; color: white; border-radius: 8px; display: flex; align-items: center; justify-content: center; }
.brand-name { font-size: 20px; font-weight: 700; color: #0f172a; letter-spacing: -0.5px; }

.nav-menu {
  flex: 1;
  display: flex;
  justify-content: center;
  gap: 10px;
}

.nav-menu .el-button {
  font-size: 15px;
  color: #64748b;
  padding: 8px 16px;
  border-radius: 8px;
  transition: all 0.2s;
}

.nav-menu .el-button:hover {
  color: #0f172a;
  background-color: #f1f5f9;
}

.nav-menu .el-button.active {
  color: #2563eb;
  background-color: #eff6ff;
  font-weight: 600;
}

.user-menu .avatar-box {
  display: flex; align-items: center; gap: 10px; cursor: pointer; padding: 4px 8px; border-radius: 6px; transition: background 0.2s;
}
.user-menu .avatar-box:hover { background: #f1f5f9; }
.user-menu .name { font-weight: 500; font-size: 14px; color: #334155; }
</style>
