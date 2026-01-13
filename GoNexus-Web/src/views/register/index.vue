<template>
<!-- 整体布局与登录页一致 -->
<div class="sao-container">
    <vue-particles id="tsparticles" :options="particlesOptions"/>
    <div class="content-wrapper">

    <!-- 结衣引导 (可以加一段注册提示) -->
    <div class="yui-area">
        <div class="speech-bubble">
        <p class="typing-text">欢迎创建新角色！</p>
        <p>请填写以下信息， <br/> 即可踏入 Sword Art Online!</p>
        </div>
        <img src="../../assets/login1.png" class="yui-sprite" alt="Yui" />
    </div>

    <!-- 注册面板 -->
    <div class="login-panel">
        <!-- 装饰角标 -->
        <div class="corner top-left"></div>
        <div class="corner top-right"></div>
        <div class="corner bottom-left"></div>
        <div class="corner bottom-right"></div>

        <div class="panel-header">
        <h1 class="app-title">GoNexus</h1>
        <p class="sub-title">SWORD ART ONLINE</p>
        </div>

        <!-- 表单区域 -->
        <el-form :model="form" :rules="rules" ref="formRef" class="sao-form">
        <!-- 用户名/ID -->
        <div class="input-group">
            <span class="input-label">玩家用户名</span>
            <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="请输入唯一的玩家用户名" class="sao-input"/>
            </el-form-item>
        </div>
        
        <!-- 密码 -->
        <div class="input-group">
            <span class="input-label">密码</span>
            <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码 (至少6位)" show-password class="sao-input"/>
            </el-form-item>
        </div>

        <!-- 确认密码 (新增) -->
        <div class="input-group">
            <span class="input-label">确认密码</span>
            <el-form-item prop="passwordConfirm">
            <el-input 
                v-model="form.passwordConfirm" 
                type="password" 
                placeholder="请再次输入密码" 
                show-password
                class="sao-input"
            />
            </el-form-item>
        </div>

        <!-- 按钮组 -->
        <div class="btn-group">
            <!-- 注册按钮 -->
            <button class="sao-btn primary-btn" @click.prevent="handleRegister">
            {{ loading ? '创建中...' : '创建角色' }}
            </button>
            <!-- 返回登录 -->
            <button class="sao-btn secondary-btn" @click.prevent="handleToLogin">
            已有账号？去登录
            </button>
        </div>
        </el-form>
    </div>
    </div>
</div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register } from '../../api/auth' // 导入注册接口


// --- 页面常量 ---
const loading = ref(false)
const formRef = ref()
const router = useRouter()

// --- 表单数据 ---
const form = reactive({
username: '',
password: '',
passwordConfirm: ''
})


// --- 自定义密码校验函数 ---
// 为什么要把校验写成函数？这样才能在 rules 里引用，并且能拿到 form 里的值
const validatePassword = (rule: any, value: string, callback: (err?: Error) => void) => {
if (value === '') {
    callback(new Error('请再次输入密码'))
} else if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
} else {
    callback() // 校验通过
}
}
// --- 表单校验规则 ---
const rules = {
username: [
    { required: true, message: '请输入玩家ID', trigger: 'blur' },
    { min: 4, message: '玩家ID至少4位', trigger: 'blur' }
],
password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
],
passwordConfirm: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    // 自定义校验：两次密码必须一致
    { validator: validatePassword, trigger: 'blur' }
]
}



// --- 业务逻辑 ---

// 注册处理
const handleRegister = async () => {
if (!formRef.value) return
await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
    loading.value = true
    try {
        // 调用后端注册接口
        await register(form)
        ElMessage.success('角色创建成功！请前往登录')
        // 注册成功后，跳转到登录页
        router.push('/login') 
    } catch (error) {
        // 错误已经在 request.ts 中处理并弹出提示
    } finally {
        loading.value = false
    }
    }
})
}

// 跳转到登录页
const handleToLogin = () => {
router.push('/login')
}

// --- 粒子特效配置 (和登录页一样) ---
const particlesOptions = {
background: { color: { value: "transparent" } },
fpsLimit: 120,
interactivity: { events: { onClick: { enable: true, mode: "push" }, onHover: { enable: true, mode: "repulse" } }, modes: { bubble: { distance: 400, duration: 2, opacity: 0.8, size: 40 }, push: { quantity: 4 }, repulse: { distance: 100, duration: 0.4 } } },
particles: { color: { value: "#ffffff" }, links: { color: "#ffffff", distance: 150, enable: true, opacity: 0.2, width: 1 }, move: { direction: "none", enable: true, outModes: "bounce", random: false, speed: 1, straight: false }, number: { density: { enable: true, area: 800 }, value: 60 }, opacity: { value: 0.5 }, shape: { type: "circle" }, size: { value: { min: 1, max: 3 } } },
detectRetina: true,
};
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
align-items: center;
gap: 40px;
}

/* --- 结衣精灵 --- */
.yui-area {
position: relative;
display: flex;
flex-direction: column;
align-items: center;
animation: float 3s ease-in-out infinite;
}

.yui-sprite {
width: 180px;
filter: drop-shadow(0 0 10px rgba(255, 255, 255, 0.6));
}

.speech-bubble {
background: rgba(255, 255, 255, 0.9);
padding: 15px 20px;
border-radius: 20px;
border-bottom-left-radius: 2px;
margin-bottom: 10px;
box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
max-width: 200px;
font-size: 14px;
color: #333;
position: relative;
}
.speech-bubble::after {
content: '';
position: absolute;
bottom: -8px;
left: 20px;
border-width: 8px 8px 0;
border-style: solid;
border-color: rgba(255, 255, 255, 0.9) transparent;
}

/* --- 登录面板 (复用) --- */
.login-panel {
width: 380px;
background: rgba(255, 255, 255, 0.75);
backdrop-filter: blur(10px);
border-radius: 12px;
padding: 40px 30px;
box-shadow: 0 8px 32px rgba(31, 38, 135, 0.2);
border: 1px solid rgba(255, 255, 255, 0.5);
position: relative;
overflow: hidden;
}
.login-panel::before { /* 顶部蓝色条 */
content: '';
position: absolute;
top: 0; left: 0; width: 100%; height: 6px;
background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
}

.panel-header {
text-align: center;
margin-bottom: 30px;
border-bottom: 1px solid rgba(0,0,0,0.1);
padding-bottom: 15px;
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
border-radius: 20px;
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
background-color: #2ecc71;
border-radius: 50%;
margin-right: 5px;
box-shadow: 0 0 5px #2ecc71;
}

/* SAO 装饰性角落 */
.corner {
position: absolute;
width: 15px;
height: 15px;
border-color: #2980b9; /* SAO 标志性蓝色 */
border-style: solid;
opacity: 0.6;
}
.top-left { top: 10px; left: 10px; border-width: 2px 0 0 2px; }
.top-right { top: 10px; right: 10px; border-width: 2px 2px 0 0; }
.bottom-left { bottom: 10px; left: 10px; border-width: 0 0 2px 2px; }
.bottom-right { bottom: 10px; right: 10px; border-width: 0 2px 2px 0; }

/* 结衣悬浮动画 */
@keyframes float {
0% { transform: translateY(0px); }
50% { transform: translateY(-15px); } /* 稍微向上浮 */
100% { transform: translateY(0px); }
}
.yui-area {
position: relative;
display: flex;
flex-direction: column;
align-items: center;
animation: float 3s ease-in-out infinite;
}
.yui-sprite {
width: 180px; /* 根据图片调整 */
filter: drop-shadow(0 0 10px rgba(255, 255, 255, 0.6)); /* 边缘发光 */
}
.speech-bubble {
background: rgba(255, 255, 255, 0.9);
padding: 15px 20px;
border-radius: 20px;
border-bottom-left-radius: 2px; /* 气泡尾巴 */
margin-bottom: 10px;
box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
max-width: 200px;
font-size: 14px;
color: #333;
position: relative;
}
.speech-bubble::after {
content: '';
position: absolute;
bottom: -8px;
left: 20px;
border-width: 8px 8px 0;
border-style: solid;
border-color: rgba(255, 255, 255, 0.9) transparent;
}
.typing-text::after { /* 模拟打字效果 */
content: '|';
animation: blink 0.7s infinite;
}
@keyframes blink {
0%, 100% { opacity: 1; }
50% { opacity: 0; }
}
</style>