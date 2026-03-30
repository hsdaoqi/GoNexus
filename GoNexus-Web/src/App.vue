
<template>
  <!-- 路由出口：URL 变了，这里显示的内容就会变 -->
  <!-- 比如访问 /login，这里就显示 Login.vue 的内容 -->
  <router-view></router-view>
  
  <!-- 全局 AI 小精灵 (暂时禁用) -->
  <!-- <AiCompanion v-if="showCompanion" @action="handleAiAction" /> -->
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useWebSocket } from './composables/useWebSocket'
import { useRoute, useRouter } from 'vue-router'
// import AiCompanion from '@/components/AiCompanion.vue'
import { ElMessage } from 'element-plus'

// 在应用级别管理WebSocket连接
const { isConnected } = useWebSocket()
const route = useRoute()
const router = useRouter()

// 在登录/注册页面不显示小精灵
const showCompanion = computed(() => {
  return !['/login', '/register'].includes(route.path)
})

const handleAiAction = (type: string) => {
  if (type === 'chat') {
    // 跳转到聊天页面 (假设路径是 /chat)
    router.push('/chat')
    ElMessage.success('已为您跳转到聊天页面')
  } else if (type === 'analyze') {
    // 分析功能目前主要针对动态广场
    if (route.path === '/moments') {
      // 如果已经在动态广场，这里其实很难直接调用子组件的方法
      // 可以通过 EventBus 或 Store 通信，或者简单提示
      ElMessage.info('请在动态广场页面点击小精灵进行分析哦~')
    } else {
      router.push('/moments')
      ElMessage.success('正在前往动态广场进行分析...')
    }
  }
}
</script>


<style scoped>
</style>
