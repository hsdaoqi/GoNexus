<template>
  <div class="page-container">
    <GlobalNavbar />
    <div class="moments-container">
      <div class="moments-header">
        <div class="tabs">
        <div 
          :class="['tab-item', activeTab === 'all' ? 'active' : '']" 
          @click="switchTab('all')"
        >
          全站广场
        </div>
        <div 
          :class="['tab-item', activeTab === 'friend' ? 'active' : '']" 
          @click="switchTab('friend')"
        >
          好友圈
        </div>
      </div>
      <el-button type="primary" round @click="showCreateDialog = true">
        <el-icon><EditPen /></el-icon> 发布动态
      </el-button>
    </div>

    <div class="moments-content" v-loading="loading">
      <div v-if="posts.length === 0" class="empty-state">
        <el-empty description="暂无动态，快来发布第一条吧！" />
      </div>

      <div v-else class="post-list">
        <div v-for="post in posts" :key="post.id" class="post-card">
          <div class="post-header">
            <el-avatar :src="post.user_avatar" :size="40" />
            <div class="user-info">
              <div class="nickname">{{ post.user_nickname }}</div>
              <div class="meta">
                <span class="time">{{ formatTime(post.created_at) }}</span>
                <span v-if="post.location" class="location">· {{ post.location }}</span>
              </div>
            </div>
            <div class="mood-tag" v-if="post.mood && post.mood !== 'Neutral'">
              {{ getMoodEmoji(post.mood) }}
            </div>
          </div>

          <div class="post-body">
            <div class="text-content">{{ post.content }}</div>
            <div v-if="post.media && post.media.length" class="media-grid">
              <el-image 
                v-for="(url, index) in post.media" 
                :key="index" 
                :src="url" 
                :preview-src-list="post.media"
                fit="cover"
                class="media-item"
              />
            </div>
          </div>

          <div class="post-footer">
            <div class="action-btn" @click="handleLike(post)" :class="{ liked: post.is_liked }">
              <el-icon><star-filled v-if="post.is_liked" /><star v-else /></el-icon>
              <span>{{ post.like_count || '赞' }}</span>
            </div>
            <div class="action-btn" @click="handleComment(post)">
              <el-icon><chat-dot-square /></el-icon>
              <span>{{ post.comment_count || '评论' }}</span>
            </div>
          </div>
          
          <!-- 评论区 -->
          <div class="comments-section" v-if="post.showComments">
            <div class="comment-list" v-loading="post.loadingComments">
              <div v-if="!post.comments || post.comments.length === 0" class="empty-comments">
                暂无评论
              </div>
              <div v-for="comment in post.comments" :key="comment.id" class="comment-item">
                <el-avatar :src="comment.user_avatar" :size="24" class="comment-avatar" />
                <div class="comment-content-box">
                  <span class="comment-user">{{ comment.user_nickname }}</span>
                  <span class="comment-text">{{ comment.content }}</span>
                  <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
                </div>
              </div>
            </div>
            <div class="comment-input-box">
              <el-input 
                v-model="post.newComment" 
                placeholder="写下你的评论..." 
                class="comment-input"
                @keyup.enter="submitComment(post)"
              >
                <template #append>
                  <el-button @click="submitComment(post)" :loading="post.submittingComment">发送</el-button>
                </template>
              </el-input>
            </div>
          </div>
        </div>
      </div>
      
      <div class="load-more" v-if="posts.length > 0">
         <el-button v-if="hasMore" text @click="loadMore">加载更多</el-button>
         <span v-else class="no-more">没有更多了</span>
      </div>
    </div>

    <!-- 发布动态弹窗 -->
    <el-dialog v-model="showCreateDialog" title="发布新动态" width="500px">
      <el-form :model="createForm">
        <el-form-item>
          <el-input 
            v-model="createForm.content" 
            type="textarea" 
            :rows="4" 
            placeholder="分享你的新鲜事..." 
          />
        </el-form-item>
        <el-form-item label="可见性">
          <el-radio-group v-model="createForm.visibility">
            <el-radio :label="0">公开</el-radio>
            <el-radio :label="1">仅好友</el-radio>
            <el-radio :label="2">私密</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="配图">
          <el-upload
            v-model:file-list="fileList"
            action="http://localhost:8080/api/v1/file/upload"
            :headers="uploadHeaders"
            list-type="picture-card"
            :on-success="handleUploadSuccess"
            :on-remove="handleRemove"
            :limit="9"
            name="file"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="submitPost" :loading="submitting">发布</el-button>
        </span>
      </template>
    </el-dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import GlobalNavbar from '@/components/GlobalNavbar.vue'
import { EditPen, Star, StarFilled, ChatDotSquare, Plus } from '@element-plus/icons-vue'
import { getMoments, createPost, toggleLike, createComment, getComments, type Post as ApiPost } from '@/api/moment'
import { ElMessage, type UploadProps, type UploadUserFile } from 'element-plus'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

// 扩展 Post 类型以支持 UI 状态
interface Post extends ApiPost {
  showComments?: boolean
  loadingComments?: boolean
  submittingComment?: boolean
  newComment?: string
}

const activeTab = ref('all')
const posts = ref<Post[]>([])
const loading = ref(false)
const page = ref(1)
const hasMore = ref(true)

const showCreateDialog = ref(false)
const submitting = ref(false)
const createForm = reactive({
  content: '',
  visibility: 0,
  media: [] as string[]
})

const fileList = ref<UploadUserFile[]>([])

const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${localStorage.getItem('token')}`
  }
})

const handleRemove: UploadProps['onRemove'] = (uploadFile) => {
  if (uploadFile.response && (uploadFile.response as any).data) {
    const url = (uploadFile.response as any).data.url
    const index = createForm.media.indexOf(url)
    if (index !== -1) {
      createForm.media.splice(index, 1)
    }
  }
}

const handleUploadSuccess: UploadProps['onSuccess'] = (response) => {
  if (response.code === 200) {
    createForm.media.push(response.data.url)
  } else {
    ElMessage.error(response.msg || '上传失败')
  }
}

const switchTab = (tab: string) => {
  if (activeTab.value === tab) return
  activeTab.value = tab
  page.value = 1
  posts.value = []
  hasMore.value = true
  fetchPosts()
}

const fetchPosts = async () => {
  if (loading.value) return
  loading.value = true
  try {
    const res: any = await getMoments({
      page: page.value,
      page_size: 10,
      type: activeTab.value
    })
    if (res.list && res.list.length > 0) {
      posts.value.push(...res.list)
      if (res.list.length < 10) {
        hasMore.value = false
      }
    } else {
      hasMore.value = false
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  page.value++
  fetchPosts()
}

const submitPost = async () => {
  if (submitting.value) return
  if (!createForm.content.trim()) {
    ElMessage.warning('请输入内容')
    return
  }
  submitting.value = true
  try {
    await createPost({
      content: createForm.content,
      visibility: createForm.visibility,
      media: createForm.media
    })
    ElMessage.success('发布成功')
    showCreateDialog.value = false
    createForm.content = ''
    createForm.media = []
    fileList.value = []
    // 刷新列表
    page.value = 1
    posts.value = []
    hasMore.value = true
    fetchPosts()
  } catch (error) {
    ElMessage.error('发布失败')
  } finally {
    submitting.value = false
  }
}

const getMoodEmoji = (mood: string) => {
  const emojis: Record<string, string> = {
    'Happy': '😊',
    'Sad': '😢',
    'Angry': '😠',
    'Surprised': '😲',
    'Neutral': '😐'
  }
  return emojis[mood] || '😐'
}

const formatTime = (time: string) => {
  return dayjs(time).fromNow()
}

const handleLike = async (post: Post) => {
  try {
    const res: any = await toggleLike(post.id)
    post.is_liked = res.is_liked
    post.like_count = res.like_count
  } catch (error) {
    ElMessage.error('操作失败')
  }
}


const handleComment = async (post: Post) => {
  post.showComments = !post.showComments
  if (post.showComments && (!post.comments || post.comments.length === 0)) {
    post.loadingComments = true
    try {
      const res: any = await getComments(post.id)
      post.comments = res || []
    } catch (error) {
      console.error(error)
      ElMessage.error('获取评论失败')
    } finally {
      post.loadingComments = false
    }
  }
}

const submitComment = async (post: Post) => {
  if (!post.newComment?.trim()) return
  
  post.submittingComment = true
  try {
    const res: any = await createComment({
      post_id: post.id,
      content: post.newComment
    })
    
    if (!post.comments) post.comments = []
    post.comments.push(res)
    post.comment_count = (post.comment_count || 0) + 1
    post.newComment = ''
    ElMessage.success('评论成功')
  } catch (error) {
    ElMessage.error('评论失败')
  } finally {
    post.submittingComment = false
  }
}

onMounted(() => {
  fetchPosts()
})
</script>

<style scoped>
.moments-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.moments-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  background: white;
  padding: 15px 20px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.tabs {
  display: flex;
  gap: 20px;
}

.tab-item {
  font-size: 16px;
  font-weight: 600;
  color: #606266;
  cursor: pointer;
  padding-bottom: 5px;
  transition: all 0.3s;
}

.tab-item.active {
  color: #409eff;
  border-bottom: 2px solid #409eff;
}

.post-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.post-header {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.user-info {
  margin-left: 12px;
  flex: 1;
}

.nickname {
  font-weight: 600;
  font-size: 15px;
  color: #303133;
}

.meta {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.post-body {
  margin-bottom: 15px;
}

.text-content {
  font-size: 15px;
  line-height: 1.6;
  color: #303133;
  margin-bottom: 10px;
  white-space: pre-wrap;
}

.post-footer {
  display: flex;
  border-top: 1px solid #ebeef5;
  padding-top: 12px;
}

.action-btn {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
  gap: 5px;
  transition: color 0.2s;
}

.action-btn:hover {
  color: #409eff;
}

.action-btn.liked {
  color: #f56c6c;
}

/* 评论区样式 */
.comments-section {
  margin-top: 15px;
  background-color: #f9fafc;
  border-radius: 8px;
  padding: 15px;
}

.comment-list {
  margin-bottom: 15px;
}

.empty-comments {
  text-align: center;
  color: #909399;
  font-size: 13px;
  padding: 10px 0;
}

.comment-item {
  display: flex;
  gap: 10px;
  margin-bottom: 12px;
}

.comment-content-box {
  flex: 1;
  font-size: 14px;
  line-height: 1.5;
}

.comment-user {
  font-weight: 600;
  color: #303133;
  margin-right: 8px;
}

.comment-text {
  color: #606266;
}

.comment-time {
  display: block;
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.comment-input-box {
  display: flex;
  align-items: center;
}

.load-more {
  text-align: center;
  padding: 20px 0;
}

.no-more {
  color: #909399;
  font-size: 13px;
}
</style>
