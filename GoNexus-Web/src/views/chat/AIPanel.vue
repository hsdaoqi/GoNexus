<template>
    <div class="ai-panel">
      <div class="ai-header">
        <span>YUI Navigation</span>
        <el-icon class="close-icon" @click="$emit('close')"><Close /></el-icon>
      </div>
      <div class="ai-content">
        <div class="ai-bubble">
          <p>我是结衣。正在监控聊天流...</p>
          <p>您可以查询聊天记录摘要。</p>
        </div>
        <div class="ai-actions">
          <input v-model="aiQuery" class="sao-input-mini" placeholder="输入问题..." />
          <button class="sao-btn-mini" @click="handleAskAI">Search</button>
        </div>
        <div class="ai-result" v-if="aiAnswer">
          <div class="result-title">Result:</div>
          <div class="result-text">{{ aiAnswer }}</div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue'
  import { Close } from '@element-plus/icons-vue'
  import { askAI } from '../../api/chat'
  import { useChatStore } from '../../store/chat'
  const emit = defineEmits(['close'])
  const aiQuery = ref('')
  const aiAnswer = ref('')
  const chatStore = useChatStore()

  const handleAskAI = async () => {
    if (!aiQuery.value) return
    aiAnswer.value = 'Analyzing...'
    try {
      const target_id = chatStore.currentChat?.id
      const chat_type = chatStore.currentChat?.isGroup ? 2 : 1
      const res: any = await askAI(aiQuery.value,target_id,chat_type)
      aiAnswer.value = res.answer
    } catch (e) {
      aiAnswer.value = 'System Error.'
    }
  }
  </script>
  
  <style scoped>
  @import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700&display=swap');
  
  .ai-panel {
    width: 100%; height: 100%; background: white; border-left: 1px solid #e0e0e0; display: flex; flex-direction: column; z-index: 10; box-shadow: -2px 0 10px rgba(0,0,0,0.05);
  }
  .ai-header { padding: 15px; border-bottom: 1px solid #eee; display: flex; justify-content: space-between; font-weight: bold; color: #e67e22; font-family: 'Orbitron'; }
  .close-icon { cursor: pointer; }
  .ai-content { padding: 20px; flex: 1; overflow-y: auto; background: #fdf6ec; }
  .ai-bubble { background: white; padding: 12px; border-radius: 8px; margin-bottom: 15px; font-size: 13px; color: #666; border: 1px solid #faecd8; }
  .sao-input-mini { width: 70%; padding: 6px; border: 1px solid #ddd; border-radius: 4px; }
  .sao-btn-mini { width: 25%; background: #e67e22; color: white; border: none; border-radius: 4px; cursor: pointer; }
  .result-title { font-weight: bold; font-size: 12px; margin-bottom: 5px; color: #333; }
  .result-text { font-size: 13px; line-height: 1.6; color: #333; background: white; padding: 10px; border-radius: 6px; border: 1px solid #eee; }
  </style>
