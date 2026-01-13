import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

// 创建 axios 实例
const service = axios.create({
    // 你的 Go 后端地址
    baseURL: 'http://localhost:8080/api/v1', 
    timeout: 5000 // 请求超时时间
})

// 1. 请求拦截器：每次发请求前自动做的事
service.interceptors.request.use(
    (config) => {
        // 从浏览器缓存拿到 Token
        const token = localStorage.getItem('token')
        if (token) {
            // 给请求头加上 Token，格式要和后端中间件匹配 (Bearer )
            config.headers['Authorization'] = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// 2. 响应拦截器：收到后端回复后自动做的事
service.interceptors.response.use(
    (response) => {
        const res = response.data
        // 如果后端返回的 code 不是 200 (我们在 errcode.go 里定义的 CodeSuccess)
        if (res.code !== 200) {
            ElMessage.error(res.msg || '系统错误')
            
            // 401 代表 Token 没带或者过期了，强制登出
            // 注意：这里的 401 是后端 response.Error 返回的 http status，
            // 或者是 code 业务码，视你后端实现而定，这里做个双重保险
            if (res.code === 401 || res.code === 2004) {
                localStorage.clear() // 清空缓存
                // 使用vue-router进行跳转，避免页面刷新
                router.push('/login')
            }
            return Promise.reject(new Error(res.msg || 'Error'))
        } else {
            return res.data // 直接把 data 剥离出来返回给页面
        }
    },
    (error) => {
        // 处理 HTTP 状态码错误 (404, 500 等)
        ElMessage.error(error.message)
        return Promise.reject(error)
    }
)

export default service