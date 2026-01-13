import { createRouter, createWebHistory } from 'vue-router'
import { getUserInfo } from '../api/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/index.vue')
  },
  {
    path:'/register',
    name:'Register',
    component: () => import('../views/register/index.vue')
  },
  {
    path: '/chat',
    name: 'Chat',
    // 还没写 Home，先放个空的占位，或者暂时指向 Login
    // component: () => import('../views/home/index.vue') 
    // ↓↓↓ 临时测试用：登录成功后跳转到一个简单文字页 ↓↓↓
    component:()=>import('../views/chat/index.vue')
  },
  {
    path: '/',
    name: 'Home',
    // 这里指向新的大厅页面
    component: () => import('../views/home/index.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/profile/index.vue')
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
    const token = localStorage.getItem('token')
    const whiteList = ['/login', '/register']

    // 白名单页面直接放行
    if (whiteList.includes(to.path)) {
      next()
      return
    }

    // 没有token，跳转登录
    if (!token) {
      next('/login')
      return
    }

    // 有token但需要验证有效性
    try {
      // 这里可以调用验证token的API，比如getUserInfo
      // 如果token无效，会在响应拦截器中处理
      await getUserInfo()
      next()
    } catch (error) {
      // token无效，跳转登录
      localStorage.clear()
      next('/login')
    }
  })

export default router