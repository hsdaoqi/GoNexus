<template>
<div class="sao-container">
    <!-- 1. 粒子特效背景 -->
    <vue-particles id="tsparticles" :options="particlesOptions"/>

    <!-- 2. 内容区域 (左右布局) -->
    <div class="content-wrapper">
    
    <!-- 左侧：结衣精灵 (Yui) -->
    <div class="yui-area">
        <div class="speech-bubble">
        <p class="typing-text">欢迎来到 GoNexus 世界!</p>
        <p>我是结衣，请登录您的<br/>账号开始冒险吧 (Link Start)!</p>
        </div>
        <img src="../../assets/login1.png" class="yui-sprite" alt="Yui" />
    </div>

    <!-- 右侧：登录面板  -->
    <div class="login-panel">
        <!-- 装饰角标 (SAO 风格的四个角) -->
        <div class="corner top-left"></div>
        <div class="corner top-right"></div>
        <div class="corner bottom-left"></div>
        <div class="corner bottom-right"></div>

        <div class="panel-header">
        <h1 class="app-title">GoNexus</h1>
        <p class="sub-title">SWORD ART ONLINE</p>
        </div>

        <el-form :model="form" :rules="rules" ref="formRef" class="sao-form">
        <!-- 账号 -->
        <div class="input-group">
            <span class="input-label">玩家ID</span>
            <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="请输入玩家ID" class="sao-input"/>
            </el-form-item>
        </div>
        
        <!-- 密码 -->
        <div class="input-group">
            <span class="input-label">密码</span>
            <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password class="sao-input"/>
            </el-form-item>
        </div>

        <!-- 按钮组 -->
        <div class="btn-group">
            <button class="sao-btn primary-btn" @click.prevent="handleLogin">
            {{ loading ? '连接中...' : '开始冒险' }}
            </button>
            <button class="sao-btn secondary-btn" @click.prevent="handleRegister">
            创建角色
            </button>
        </div>
        </el-form>

        <!-- 底部状态条 -->
        <div class="system-status">
        <span class="status-dot"></span> NerveGear 已连接 - 准备就绪
        </div>
    </div>
    </div>
</div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login, register } from '../../api/auth'

const router = useRouter()
const loading = ref(false)
const formRef = ref()

const form = reactive({ username: '', password: '' })
const rules = {
username: [{ required: true, message: '请输入玩家ID', trigger: 'blur' }],
password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

// --- 核心业务逻辑 (保持不变) ---
const handleLogin = async () => {
if (!formRef.value) return
await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
    loading.value = true
    try {
        const data: any = await login(form)
        localStorage.setItem('token', data.token)
        localStorage.setItem('user_id', data.user_id)
        ElMessage.success('Link Start! 登录成功')
        router.push('/')
    } catch (error) {
    } finally {
        loading.value = false
    }
    }
})
}

const handleRegister = async () => {
    router.push('/register')
}



// --- 粒子特效配置 (模仿图4的氛围) ---
const particlesOptions = {
background: {
    color: { value: "transparent" }, // 背景透明，为了漏出后面的图片
},
fpsLimit: 120,
interactivity: {
    events: {
    onClick: { enable: true, mode: "push" },
    onHover: { enable: true, mode: "repulse" },
    },
    modes: {
    bubble: { distance: 400, duration: 2, opacity: 0.8, size: 40 },
    push: { quantity: 4 },
    repulse: { distance: 100, duration: 0.4 },
    },
},
particles: {
    color: { value: "#ffffff" }, // 粒子白色
    links: {
    color: "#ffffff",
    distance: 150,
    enable: true,
    opacity: 0.2, // 连线淡淡的
    width: 1,
    },
    move: {
    direction: "none",
    enable: true,
    outModes: "bounce",
    random: false,
    speed: 1, // 速度慢一点，唯美
    straight: false,
    },
    number: {
    density: { enable: true, area: 800 },
    value: 60, // 粒子数量
    },
    opacity: { value: 0.5 },
    shape: { type: "circle" },
    size: { value: { min: 1, max: 3 } },
},
detectRetina: true,
};
</script>

<style scoped>
/* 引入字体 (可选，如果有SAO字体更好，这里用系统字体模拟) */
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700&display=swap');

.sao-container {
height: 100vh;
width: 100vw;
/* 背景图设置 */
background-image: url('../../assets/background1.jpg');
background-size: cover;
background-position: center;
background-repeat: no-repeat;

display: flex;
justify-content: center;
align-items: center;
overflow: hidden;
position: relative;
padding: 20px; /* 添加内边距防止内容贴边 */
box-sizing: border-box;
}

/* 粒子层要覆盖在背景图之上，但在内容之下 */
#tsparticles {
position: absolute;
top: 0;
left: 0;
width: 100%;
height: 100%;
z-index: 1;
}

.content-wrapper {
z-index: 10; /* 保证在粒子上面 */
display: flex;
align-items: center;
justify-content: center;
gap: 40px; /* 结衣和登录框的距离 */
flex-wrap: wrap; /* 允许换行 */
max-width: calc(100vw - 40px); /* 减去容器的padding */
width: 100%;
box-sizing: border-box;
}

/* --- 左侧：结衣精灵 --- */
.yui-area {
position: relative;
display: flex;
flex-direction: column;
align-items: center;
animation: float 3s ease-in-out infinite; /* 上下悬浮动画 */
flex: 0 0 auto; /* 防止被压缩 */
min-width: 0; /* 允许缩小到最小 */
}

.yui-sprite {
width: 180px; /* 根据图片实际大小调整 */
filter: drop-shadow(0 0 10px rgba(255, 255, 255, 0.6));
}

/* 气泡框 */
.speech-bubble {
background: rgba(255, 255, 255, 0.9);
padding: 15px 20px;
border-radius: 20px;
border-bottom-left-radius: 2px; /* 变成气泡形状 */
margin-bottom: 10px;
box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
max-width: 200px;
font-size: 14px;
color: #333;
position: relative;
}
.speech-bubble::after { /* 小三角 */
content: '';
position: absolute;
bottom: -8px;
left: 20px;
border-width: 8px 8px 0;
border-style: solid;
border-color: rgba(255, 255, 255, 0.9) transparent;
}

/* --- 右侧：登录面板 (图4风格核心) --- */
.login-panel {
width: 380px;
max-width: calc(100vw - 60px); /* 确保在小屏幕上不会超出 */
background: rgba(255, 255, 255, 0.75); /* 半透明白 */
backdrop-filter: blur(10px); /* 毛玻璃特效 */
border-radius: 12px;
padding: 40px 30px;
box-shadow: 0 8px 32px rgba(31, 38, 135, 0.2);
border: 1px solid rgba(255, 255, 255, 0.5);
position: relative;
overflow: hidden; /* 主要是为了那个蓝色顶条 */
box-sizing: border-box;
}

/* 顶部的蓝色装饰条 */
.login-panel::before {
content: '';
position: absolute;
top: 0;
left: 0;
width: 100%;
height: 6px;
background: linear-gradient(90deg, #6dd5fa, #2980b9);
}

/* 标题区 */
.panel-header {
text-align: center;
margin-bottom: 30px;
border-bottom: 1px solid rgba(0,0,0,0.1);
padding-bottom: 15px;
}

.app-title {
font-family: 'Orbitron', sans-serif; /* 科幻字体 */
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

/* 表单样式 */
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
border-left: 3px solid #2980b9; /* 左边的小蓝条 */
}

/* 深度选择器修改 Element Plus 样式 */
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

/* 自定义按钮 */
.btn-group {
margin-top: 30px;
display: flex;
flex-direction: column;
gap: 12px;
}

.sao-btn {
width: 100%;
height: 40px;
border: none;
border-radius: 20px; /* 圆角按钮 */
font-weight: bold;
cursor: pointer;
transition: all 0.3s;
font-family: 'Orbitron', sans-serif;
letter-spacing: 1px;
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

/* 底部状态条 */
.system-status {
margin-top: 25px;
font-size: 10px;
color: #7f8c8d;
text-align: center;
background: rgba(0,0,0,0.05);
padding: 8px;
border-radius: 4px;
}
.status-dot {
display: inline-block;
width: 6px;
height: 6px;
background-color: #2ecc71; /* 绿色在线 */
border-radius: 50%;
margin-right: 5px;
box-shadow: 0 0 5px #2ecc71;
}

/* --- 装饰性角落 (SAO UI 经典元素) --- */
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

/* 响应式布局 */
@media (max-width: 1024px) {
.sao-container {
    padding: 10px;
}

.content-wrapper {
    gap: 20px;
    flex-direction: column; /* 小屏幕垂直排列 */
}

.login-panel {
    width: 100%;
    max-width: 380px;
}

.yui-area {
    order: 1; /* 结衣精灵在上面 */
}

.login-panel {
    order: 2; /* 登录面板在下面 */
}
}

@media (max-width: 480px) {
.sao-container {
    padding: 5px;
}

.content-wrapper {
    gap: 15px;
}

.yui-sprite {
    width: 140px; /* 缩小精灵图片 */
}

.speech-bubble {
    max-width: 160px;
    font-size: 12px;
}

.login-panel {
    padding: 30px 20px;
    width: 100%;
    max-width: 320px;
}

.app-title {
    font-size: 28px;
}
}

/* 动画定义 */
@keyframes float {
0% { transform: translateY(0px); }
50% { transform: translateY(-15px); }
100% { transform: translateY(0px); }
}
</style>