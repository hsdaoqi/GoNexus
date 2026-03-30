<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="emit('update:visible', $event)"
    title="Recommended Friends"
    width="500px"
    class="sao-dialog"
    :show-close="false"
  >
    <div class="dialog-content request-list">
      <div v-if="recommendList.length === 0" class="empty-tip">暂无推荐好友</div>
      <div v-for="rec in recommendList" :key="rec.id" class="request-item recommend-item-row">
        <el-avatar :size="40" :src="rec.avatar" />
        <div class="req-info">
          <div class="req-name">{{ rec.nickname || rec.username }}</div>
          <div class="req-msg" style="color: #e67e22">{{ rec.reason }}</div>
          <div class="req-tags" v-if="rec.tags && rec.tags.length > 0">
            <span v-for="tag in rec.tags" :key="tag" class="tag-pill">{{ tag }}</span>
          </div>
        </div>
        <button class="sao-btn-mini add" @click="emit('add-friend', rec)">
          <el-icon><Plus /></el-icon> ADD
        </button>
      </div>
    </div>
    <template #footer>
      <button class="sao-btn-cancel" @click="emit('update:visible', false)">CLOSE</button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'

defineProps<{
  visible: boolean
  recommendList: any[]
}>()

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'add-friend', user: any): void
}>()
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700&display=swap');

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
.sao-btn-mini.add { background: #e67e22; }

/* 推荐好友样式 */
.recommend-item-row { transition: all 0.3s; cursor: pointer; }
.recommend-item-row:hover { background: rgba(230, 126, 34, 0.05); }
.tag-pill { display: inline-block; background: rgba(230, 126, 34, 0.1); color: #e67e22; padding: 2px 6px; border-radius: 4px; font-size: 10px; margin-right: 4px; margin-top: 4px; }
</style>