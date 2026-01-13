<template>
<div class="sao-container">
    <!-- 粒子特效背景 -->
    <vue-particles id="tsparticles" :options="particlesOptions"/>

    <!-- 内容区域 -->
    <div class="content-wrapper">

    <!-- 个人信息面板 -->
    <div class="profile-panel">
        <!-- 装饰角标 -->
        <div class="corner top-left"></div>
        <div class="corner top-right"></div>
        <div class="corner bottom-left"></div>
        <div class="corner bottom-right"></div>

        <div class="panel-header">
        <div class="header-content">
            <div>
                <h1 class="app-title">GoNexus</h1>
                <p class="sub-title">USER PROFILE</p>
            </div>
            <button class="edit-profile-btn sao-btn secondary-btn" @click="showEditDialog = true">
                编辑资料
            </button>
        </div>
        </div>

        <!-- 标签页切换 -->
        <el-tabs v-model="activeTab" class="profile-tabs">
            <!-- 基本信息 -->
            <el-tab-pane label="基本信息" name="basic">
                <div class="profile-section">
                    <!-- 头像区域 -->
                    <div class="avatar-section">
                        <div class="avatar-wrapper">
                            <el-avatar
                                :size="100"
                                :src="userInfo.avatar"
                                class="user-avatar"
                            >
                                {{ userInfo.nickname ? userInfo.nickname.charAt(0).toUpperCase() : userInfo.username.charAt(0).toUpperCase() }}
                            </el-avatar>
                            <button class="avatar-edit-btn" @click="triggerFileSelect">
                                <el-icon><Edit /></el-icon>
                            </button>
                            <input
                                ref="fileInputRef"
                                type="file"
                                accept="image/*"
                                @change="handleAvatarChange"
                                style="display: none"
                            />
                        </div>
                        <div class="avatar-info">
                            <h3>{{ userInfo.nickname || userInfo.username }}</h3>
                            <p>@{{ userInfo.username }}</p>
                        </div>
                    </div>

                    <h3 class="section-title">个人信息</h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="info-label">用户名:</span>
                            <span class="info-value">{{ userInfo.username }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">昵称:</span>
                            <span class="info-value">{{ userInfo.nickname || '未设置' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">邮箱:</span>
                            <span class="info-value">{{ userInfo.email }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">性别:</span>
                            <span class="info-value">{{ userInfo.gender || '未设置' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">生日:</span>
                            <span class="info-value">{{ userInfo.birthday ? formatDate(userInfo.birthday) : '未设置' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">地区:</span>
                            <span class="info-value">{{ userInfo.location || '未设置' }}</span>
                        </div>
                        <div class="info-item full-width">
                            <span class="info-label">个性签名:</span>
                            <span class="info-value">{{ userInfo.signature || '这个人很懒，还没有写个性签名...' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">注册时间:</span>
                            <span class="info-value">{{ formatDate(userInfo.createdAt) }}</span>
                        </div>
                        <div class="info-item">
                            <span class="info-label">最后在线:</span>
                            <span class="info-value">{{ formatDate(userInfo.lastLogin) }}</span>
                        </div>
                    </div>
                </div>
            </el-tab-pane>

            <!-- 账号设置 -->
            <el-tab-pane label="账号设置" name="account">
                <div class="profile-section">
                    <h3 class="section-title">密码管理</h3>
                    <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" class="sao-form">
                        <div class="input-group">
                            <span class="input-label">当前密码</span>
                            <el-form-item prop="oldPassword">
                                <el-input
                                    v-model="passwordForm.oldPassword"
                                    type="password"
                                    placeholder="请输入当前密码"
                                    show-password
                                    class="sao-input"
                                />
                            </el-form-item>
                        </div>

                        <div class="input-group">
                            <span class="input-label">新密码</span>
                            <el-form-item prop="newPassword">
                                <el-input
                                    v-model="passwordForm.newPassword"
                                    type="password"
                                    placeholder="请输入新密码 (至少6位)"
                                    show-password
                                    class="sao-input"
                                />
                            </el-form-item>
                        </div>

                        <div class="input-group">
                            <span class="input-label">确认新密码</span>
                            <el-form-item prop="confirmPassword">
                                <el-input
                                    v-model="passwordForm.confirmPassword"
                                    type="password"
                                    placeholder="请再次输入新密码"
                                    show-password
                                    class="sao-input"
                                />
                            </el-form-item>
                        </div>

                        <div class="action-btns">
                            <button class="sao-btn primary-btn" @click="handleChangePassword" :disabled="passwordLoading">
                                {{ passwordLoading ? '修改中...' : '修改密码' }}
                            </button>
                        </div>
                    </el-form>
                </div>
            </el-tab-pane>

            <!-- 隐私设置 -->
            <el-tab-pane label="隐私设置" name="privacy">
                <div class="profile-section">
                    <h3 class="section-title">隐私选项</h3>
                    <div class="privacy-settings">
                        <div class="setting-item">
                            <div class="setting-info">
                                <h4>在线状态</h4>
                                <p>允许其他人看到您的在线状态</p>
                            </div>
                            <el-switch
                                v-model="privacySettings.showOnlineStatus"
                                @change="handlePrivacyChange"
                            />
                        </div>
                        <div class="setting-item">
                            <div class="setting-info">
                                <h4>最后在线时间</h4>
                                <p>显示您的最后在线时间</p>
                            </div>
                            <el-switch
                                v-model="privacySettings.showLastSeen"
                                @change="handlePrivacyChange"
                            />
                        </div>
                        <div class="setting-item">
                            <div class="setting-info">
                                <h4>个人资料可见性</h4>
                                <p>允许其他人查看您的个人资料</p>
                            </div>
                            <el-switch
                                v-model="privacySettings.profileVisible"
                                @change="handlePrivacyChange"
                            />
                        </div>
                    </div>
                </div>
            </el-tab-pane>
        </el-tabs>

        <!-- 返回按钮 -->
        <div class="back-btn">
            <button class="sao-btn secondary-btn" @click="handleBack">
                返回主页
            </button>
        </div>
    </div>
    </div>

    <!-- 编辑信息对话框 -->
    <el-dialog
        v-model="showEditDialog"
        title="编辑个人资料"
        width="500px"
        :close-on-click-modal="false"
    >
        <el-form :model="editForm" :rules="editRules" ref="editFormRef" class="sao-form">
            <div class="input-group">
                <span class="input-label">昵称</span>
                <el-form-item prop="nickname">
                    <el-input
                        v-model="editForm.nickname"
                        placeholder="请输入昵称"
                        class="sao-input"
                    />
                </el-form-item>
            </div>

            <div class="input-group">
                <span class="input-label">邮箱地址</span>
                <el-form-item prop="email">
                    <el-input
                        v-model="editForm.email"
                        placeholder="请输入邮箱地址"
                        class="sao-input"
                    />
                </el-form-item>
            </div>

            <div class="input-group">
                <span class="input-label">性别</span>
                <el-form-item prop="gender">
                    <el-select
                        v-model="editForm.gender"
                        placeholder="请选择性别"
                        class="sao-input"
                        clearable
                    >
                        <el-option label="男" value="male" />
                        <el-option label="女" value="female" />
                        <el-option label="其他" value="other" />
                    </el-select>
                </el-form-item>
            </div>

            <div class="input-group">
                <span class="input-label">生日</span>
                <el-form-item prop="birthday">
                    <el-date-picker
                        v-model="editForm.birthday"
                        type="date"
                        placeholder="选择生日"
                        class="sao-input"
                        format="YYYY-MM-DD"
                        value-format="YYYY-MM-DD"
                    />
                </el-form-item>
            </div>

            <div class="input-group">
                <span class="input-label">地区</span>
                <el-form-item prop="location">
                    <el-input
                        v-model="editForm.location"
                        placeholder="请输入地区"
                        class="sao-input"
                    />
                </el-form-item>
            </div>

            <div class="input-group">
                <span class="input-label">个性签名</span>
                <el-form-item prop="bio">
                    <el-input
                        v-model="editForm.signature"
                        type="textarea"
                        :rows="3"
                        placeholder="写点什么..."
                        class="sao-input"
                        maxlength="200"
                        show-word-limit
                    />
                </el-form-item>
            </div>
        </el-form>

        <template #footer>
            <div class="dialog-footer">
                <button class="sao-btn secondary-btn" @click="showEditDialog = false">
                    取消
                </button>
                <button class="sao-btn primary-btn" @click="handleUpdateProfile" :disabled="editLoading">
                    {{ editLoading ? '保存中...' : '保存' }}
                </button>
            </div>
        </template>
    </el-dialog>
</div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Edit } from '@element-plus/icons-vue'
import { useUserStore } from '../../store/user'
import { updateAvatar, updateUserInfo } from '../../api/user'

const userStore = useUserStore()

// --- 页面状态 ---
const router = useRouter()
const activeTab = ref('basic')
const showEditDialog = ref(false)
const passwordLoading = ref(false)
const editLoading = ref(false)
const fileInputRef = ref()

// --- 用户信息数据 ---
const userInfo = reactive({
    username: '',
    nickname: '',
    email: '',
    avatar: '',
    gender: '',
    birthday: '',
    location: '',
    signature: '',
    createdAt: '',
    lastLogin: ''
})

// --- 隐私设置 ---
const privacySettings = reactive({
    showOnlineStatus: true,
    showLastSeen: true,
    profileVisible: true
})

// --- 编辑表单 ---
const editForm = reactive({
    nickname: '',
    email: '',
    gender: '',
    birthday: '',
    location: '',
    signature: ''
})
const editFormRef = ref()

// --- 密码表单 ---
const passwordForm = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
})
const passwordFormRef = ref()

// --- 表单验证规则 ---
const editRules = {
    nickname: [
        { max: 20, message: '昵称不能超过20个字符', trigger: 'blur' }
    ],
    email: [
        { required: true, message: '请输入邮箱地址', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
    ],
    location: [
        { max: 50, message: '地区不能超过50个字符', trigger: 'blur' }
    ],
    bio: [
        { max: 200, message: '个性签名不能超过200个字符', trigger: 'blur' }
    ]
}

const passwordRules = {
    oldPassword: [
        { required: true, message: '请输入当前密码', trigger: 'blur' }
    ],
    newPassword: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, message: '密码至少6位', trigger: 'blur' }
    ],
    confirmPassword: [
        { required: true, message: '请确认新密码', trigger: 'blur' },
        {
            validator: (rule: any, value: string, callback: (err?: Error) => void) => {
                if (value !== passwordForm.newPassword) {
                    callback(new Error('两次输入的密码不一致'))
                } else {
                    callback()
                }
            },
            trigger: 'blur'
        }
    ]
}

// --- 业务逻辑 ---
// 加载用户信息
const loadUserInfo = async () => {
    try {
        // 从后端获取最新用户信息
        await userStore.fetchUserInfo()
        console.log(userStore.userInfo)
        // 从store中获取数据，确保数据是最新的
        userInfo.username = userStore.userInfo.username || 'Player001'
        userInfo.nickname = userStore.userInfo.nickname || '测试用户'
        userInfo.email = userStore.userInfo.email || 'player@example.com'
        userInfo.avatar = userStore.userInfo.avatar || ''
        userInfo.gender = userStore.userInfo.gender || 'male'
        userInfo.birthday = userStore.userInfo.birthday || '1990-01-01'
        userInfo.location = userStore.userInfo.location || '北京市'
        userInfo.signature = userStore.userInfo.signature || 'Hello, I am using GoNexus!'
        userInfo.createdAt = userStore.userInfo.createdAt || new Date().toISOString()
        userInfo.lastLogin = userStore.userInfo.lastLogin || new Date().toISOString()

        // 初始化编辑表单
        editForm.nickname = userInfo.nickname
        editForm.email = userInfo.email
        editForm.gender = userInfo.gender
        editForm.birthday = userInfo.birthday
        editForm.location = userInfo.location
        editForm.signature = userInfo.signature
    } catch (error) {
        ElMessage.error('加载用户信息失败')
        console.error('加载用户信息失败:', error)
    }
}

// 格式化日期
const formatDate = (dateString: string) => {
    if (!dateString) return '未知'
    return new Date(dateString).toLocaleString('zh-CN')
}

// 更新个人信息
const handleUpdateProfile = async () => {
    if (!editFormRef.value) return
    await editFormRef.value.validate(async (valid: boolean) => {
        if (valid) {
            editLoading.value = true
            try {
                // 准备更新数据，后端期望ISO 8601格式的日期
                const updateData = {
                    nickname: editForm.nickname,
                    email: editForm.email,
                    gender: editForm.gender,
                    birthday: editForm.birthday ? new Date(editForm.birthday + 'T00:00:00Z').toISOString() : null,
                    location: editForm.location,
                    signature: editForm.signature
                }

                // 调用更新用户信息的API
                await updateUserInfo(updateData)

                // 更新本地状态
                userInfo.nickname = editForm.nickname
                userInfo.email = editForm.email
                userInfo.gender = editForm.gender
                userInfo.birthday = editForm.birthday
                userInfo.location = editForm.location
                userInfo.signature = editForm.signature

                ElMessage.success('个人资料更新成功')
                showEditDialog.value = false
            } catch (error) {
                ElMessage.error('更新失败')
            } finally {
                editLoading.value = false
            }
        }
    })
}

// 头像上传处理
const triggerFileSelect = () => {
    fileInputRef.value?.click()
    
}

const handleAvatarChange = async (event: Event) => {
    const file = (event.target as HTMLInputElement).files?.[0]
    if (file) {
        // 检查文件类型
        if (!file.type.startsWith('image/')) {
            ElMessage.error('请选择图片文件')
            return
        }

        // 检查文件大小 (限制为5MB)
        if (file.size > 5 * 1024 * 1024) {
            ElMessage.error('图片大小不能超过5MB')
            return
        }

        // 先显示本地预览，让用户立即看到上传的图片
        const reader = new FileReader()
        reader.onload = (e) => {
            userInfo.avatar = e.target?.result as string
            // 同时更新store中的头像，确保其他页面也能看到预览
            userStore.setUserInfo({ avatar: e.target?.result as string })
        }
        reader.readAsDataURL(file)

        try {
            // 调用上传头像的API
            const formData = new FormData()
            formData.append('avatar', file)
            const res: any = await updateAvatar(formData)
            userInfo.avatar = res.avatar
            userStore.setUserInfo({ avatar: res.avatar })
            ElMessage.success('头像上传成功')

            // 重新从后端获取完整用户信息，确保数据同步
            await userStore.fetchUserInfo()
        } catch (error: any) {
            // 如果上传失败，显示错误信息
            ElMessage.error(error?.response?.data?.message || '头像上传失败')

            // 上传失败时，恢复原来的头像
            await userStore.fetchUserInfo()
        }
    }
}

// 隐私设置变更处理
const handlePrivacyChange = async () => {
    try {
        // 这里应该调用更新隐私设置的API
        ElMessage.success('隐私设置已更新')
    } catch (error) {
        ElMessage.error('设置更新失败')
    }
}

// 修改密码
const handleChangePassword = async () => {
    if (!passwordFormRef.value) return
    await passwordFormRef.value.validate(async (valid: boolean) => {
        if (valid) {
            passwordLoading.value = true
            try {
                // 调用修改密码的API
                ElMessage.success('密码修改成功')
                // 清空表单
                passwordForm.oldPassword = ''
                passwordForm.newPassword = ''
                passwordForm.confirmPassword = ''
            } catch (error) {
                ElMessage.error('密码修改失败')
            } finally {
                passwordLoading.value = false
            }
        }
    })
}

// 返回主页
const handleBack = () => {
    router.push('/')
}

// --- 生命周期 ---
onMounted(() => {
    loadUserInfo()
})

// --- 粒子特效配置 (和登录页一样) ---
const particlesOptions = {
    background: { color: { value: "transparent" } },
    fpsLimit: 120,
    interactivity: {
        events: { onClick: { enable: true, mode: "push" }, onHover: { enable: true, mode: "repulse" } },
        modes: { bubble: { distance: 400, duration: 2, opacity: 0.8, size: 40 }, push: { quantity: 4 }, repulse: { distance: 100, duration: 0.4 } }
    },
    particles: {
        color: { value: "#ffffff" },
        links: { color: "#ffffff", distance: 150, enable: true, opacity: 0.2, width: 1 },
        move: { direction: "none", enable: true, outModes: "bounce", random: false, speed: 1, straight: false },
        number: { density: { enable: true, area: 800 }, value: 60 },
        opacity: { value: 0.5 },
        shape: { type: "circle" },
        size: { value: { min: 1, max: 3 } }
    },
    detectRetina: true,
}
</script>

<style scoped>
/* 引入字体 */
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700&display=swap');

.sao-container {
    height: 100vh;
    width: 100vw;
    background-image: url('../../assets/background1.jpg');
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    position: relative;
    padding: 20px;
    box-sizing: border-box;
}

#tsparticles {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 1;
}

.content-wrapper {
    z-index: 10;
    display: flex;
    justify-content: center;
    align-items: center;
    max-width: 900px;
    width: 100%;
    box-sizing: border-box;
}


/* --- 个人信息面板 --- */
.profile-panel {
    width: 100%;
    max-width: 800px;
    background: rgba(255, 255, 255, 0.75);
    backdrop-filter: blur(10px);
    border-radius: 12px;
    padding: 40px 30px;
    box-shadow: 0 8px 32px rgba(31, 38, 135, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.5);
    position: relative;
    overflow: hidden;
    box-sizing: border-box;
}

.profile-panel::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 6px;
    background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
}

.panel-header {
    margin-bottom: 30px;
    border-bottom: 1px solid rgba(0,0,0,0.1);
    padding-bottom: 15px;
}

.header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.edit-profile-btn {
    font-size: 14px;
    padding: 8px 16px;
    height: auto;
}

.app-title {
    font-family: 'Orbitron', sans-serif;
    font-size: 32px;
    color: #2c3e50;
    margin: 0;
    letter-spacing: 2px;
}

.sub-title {
    font-size: 10px;
    color: #57606f;
    letter-spacing: 4px;
    margin-top: 5px;
    font-weight: bold;
}

/* 标签页样式 */
:deep(.profile-tabs .el-tabs__header) {
    margin-bottom: 20px;
}

:deep(.profile-tabs .el-tabs__nav-wrap::after) {
    display: none;
}

:deep(.profile-tabs .el-tabs__item) {
    color: #2c3e50;
    font-weight: bold;
    border-bottom: none;
    padding: 0 20px;
}

:deep(.profile-tabs .el-tabs__item.is-active) {
    color: #2980b9;
    border-bottom: 2px solid #2980b9;
}

:deep(.profile-tabs .el-tabs__active-bar) {
    display: none;
}

/* 信息展示样式 */
.profile-section {
    margin-bottom: 20px;
}

.section-title {
    font-size: 18px;
    color: #2c3e50;
    margin-bottom: 20px;
    padding-left: 5px;
    border-left: 3px solid #2980b9;
    font-weight: bold;
}

.info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 15px;
    margin-bottom: 20px;
}

.info-item {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.info-item.full-width {
    grid-column: 1 / -1;
}

.info-label {
    font-size: 12px;
    color: #57606f;
    font-weight: bold;
}

.info-value {
    font-size: 14px;
    color: #2c3e50;
    background: rgba(255, 255, 255, 0.8);
    padding: 8px 12px;
    border-radius: 6px;
    border: 1px solid #dcdfe6;
}

/* 头像区域样式 */
.avatar-section {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 30px;
    padding: 20px;
    background: rgba(255, 255, 255, 0.8);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.5);
}

.avatar-wrapper {
    position: relative;
    display: inline-block;
}

.user-avatar {
    border: 3px solid #2980b9;
    box-shadow: 0 4px 12px rgba(41, 128, 185, 0.3);
}

.avatar-edit-btn {
    position: absolute;
    bottom: 0;
    right: 0;
    background: #2980b9;
    color: white;
    border: none;
    border-radius: 50%;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.3s;
    box-shadow: 0 2px 8px rgba(41, 128, 185, 0.3);
}

.avatar-edit-btn:hover {
    background: #3498db;
    transform: scale(1.1);
}

.avatar-info h3 {
    margin: 0 0 5px 0;
    color: #2c3e50;
    font-size: 20px;
    font-weight: bold;
}

.avatar-info p {
    margin: 0;
    color: #57606f;
    font-size: 14px;
}

/* 隐私设置样式 */
.privacy-settings {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.setting-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px;
    background: rgba(255, 255, 255, 0.8);
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.5);
}

.setting-info h4 {
    margin: 0 0 5px 0;
    color: #2c3e50;
    font-size: 16px;
    font-weight: bold;
}

.setting-info p {
    margin: 0;
    color: #57606f;
    font-size: 14px;
}

/* 输入框样式 */
.input-group {
    margin-bottom: 15px;
}

.input-label {
    display: block;
    font-size: 12px;
    color: #2c3e50;
    font-weight: bold;
    margin-bottom: 5px;
    padding-left: 5px;
    border-left: 3px solid #2980b9;
}

/* Element Plus 组件样式覆盖 */
:deep(.el-input__wrapper) {
    background-color: rgba(255, 255, 255, 0.8) !important;
    box-shadow: none !important;
    border: 1px solid #dcdfe6;
    border-radius: 8px;
    transition: all 0.3s;
}

:deep(.el-input__wrapper.is-focus) {
    border-color: #2980b9 !important;
    box-shadow: 0 0 8px rgba(41, 128, 185, 0.2) !important;
}

/* 按钮样式 */
.action-btns {
    margin-top: 20px;
    display: flex;
    gap: 10px;
    justify-content: center;
}

.back-btn {
    margin-top: 30px;
    text-align: center;
}

.sao-btn {
    width: auto;
    padding: 0 20px;
    height: 36px;
    border: none;
    border-radius: 18px;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.3s;
    font-family: 'Orbitron', sans-serif;
    letter-spacing: 1px;
    font-size: 12px;
}

.primary-btn {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    color: white;
    box-shadow: 0 4px 15px rgba(0, 242, 254, 0.4);
}

.primary-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 242, 254, 0.6);
}

.secondary-btn {
    background: #fff;
    border: 1px solid #4facfe;
    color: #4facfe;
}

.secondary-btn:hover {
    background: #f0f9ff;
}

/* 对话框样式 */
:deep(.el-dialog) {
    border-radius: 12px;
    overflow: hidden;
}

:deep(.el-dialog__header) {
    background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
    color: white;
    margin: 0;
    padding: 15px 20px;
    font-family: 'Orbitron', sans-serif;
    letter-spacing: 1px;
}

:deep(.el-dialog__body) {
    padding: 20px;
}

.dialog-footer {
    display: flex;
    gap: 10px;
    justify-content: flex-end;
}

/* 装饰性角落 */
.corner {
    position: absolute;
    width: 15px;
    height: 15px;
    border-color: #2980b9;
    border-style: solid;
    opacity: 0.6;
}
.top-left { top: 10px; left: 10px; border-width: 2px 0 0 2px; }
.top-right { top: 10px; right: 10px; border-width: 2px 2px 0 0; }
.bottom-left { bottom: 10px; left: 10px; border-width: 0 0 2px 2px; }
.bottom-right { bottom: 10px; right: 10px; border-width: 0 2px 2px 0; }

/* 悬浮动画 */
@keyframes float {
    0% { transform: translateY(0px); }
    50% { transform: translateY(-15px); }
    100% { transform: translateY(0px); }
}


/* 响应式设计 */
@media (max-width: 1024px) {
    .content-wrapper {
        max-width: 700px;
    }

    .profile-panel {
        padding: 30px 20px;
    }

    .info-grid {
        grid-template-columns: 1fr;
    }

    .avatar-section {
        flex-direction: column;
        text-align: center;
        gap: 15px;
    }

    .header-content {
        flex-direction: column;
        gap: 15px;
        align-items: flex-start;
    }
}

@media (max-width: 480px) {
    .sao-container {
        padding: 10px;
    }

    .header-content {
        flex-direction: column;
        gap: 10px;
        align-items: stretch;
    }

    .edit-profile-btn {
        align-self: flex-end;
        font-size: 12px;
        padding: 6px 12px;
    }

    .profile-panel {
        padding: 25px 15px;
    }

    .app-title {
        font-size: 28px;
    }

    .stat-card {
        padding: 15px;
    }

    .stat-value {
        font-size: 20px;
    }
}
</style>
