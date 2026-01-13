import { defineStore } from 'pinia'
import { getUserInfo } from '../api/user'

export const useUserStore = defineStore('user', {
  // 1. State: 相当于组件里的 data，存数据
  state: () => ({
    userInfo: {
      id: 0,
      username: 'Player',
      nickname: '',
      avatar: '',
      email: '',
      signature: '',
      gender: '',
      birthday: '',
      location: '',
      createdAt: '',
      lastLogin: ''
    }
  }),

  // 2. Actions: 相当于组件里的 methods，写业务逻辑
  actions: {
    // 异步获取用户信息并存入 state
    async fetchUserInfo() {
      try {
        const res: any = await getUserInfo()
        // 处理字段名映射，支持snake_case和camelCase
        const mappedData = {
          id: res.id || res.ID,
          username: res.username,
          nickname: res.nickname,
          avatar: res.avatar,
          email: res.email,
          signature: res.signature || res.bio,
          gender: res.gender,
          birthday: res.birthday,
          location: res.location,
          createdAt: res.createdAt || res.created_at,
          lastLogin: res.lastLogin || res.last_login
        }

        // 合并现有数据和新数据，确保所有字段都被正确设置
        this.userInfo = {
          ...this.userInfo,
          ...mappedData
        }

        // 处理默认头像
        if (!this.userInfo.avatar) {
          this.userInfo.avatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
        }

        console.log('处理后的用户信息:', this.userInfo)
      } catch (error) {
        console.error('获取用户信息失败', error)
        // 如果获取用户信息失败，说明token可能过期，清空登录状态
        this.clearUser()
        // 清空localStorage中的token
        localStorage.removeItem('token')
        throw error
      }
    },

    // 更新用户信息（用于头像等局部更新）
    setUserInfo(newInfo: Partial<typeof this.userInfo>) {
      this.userInfo = {
        ...this.userInfo,
        ...newInfo
      }
    },

    // 登出时清空数据
    clearUser() {
      this.$reset()
    }
  },

  // 3. 开启持久化：刷新页面数据不丢失
  persist: true
})