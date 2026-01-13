
<template>
  <!-- 路由出口：URL 变了，这里显示的内容就会变 -->
  <!-- 比如访问 /login，这里就显示 Login.vue 的内容 -->
    <router-view></router-view>
  </template>
<script setup lang="ts">
import { useWebSocket } from './composables/useWebSocket'

// 在应用级别管理WebSocket连接
const { isConnected } = useWebSocket()
</script>


<style scoped>
</style>
