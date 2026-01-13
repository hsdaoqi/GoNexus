import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// 引入 Element Plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// 引入图标库 (登录页要用图标)
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
// 🔥 引入粒子插件
import Particles from "@tsparticles/vue3";
import { loadSlim } from "@tsparticles/slim"; // 加载轻量版引擎
//引入pinia
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)//启用持久化插件

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
// 🔥 注册粒子
app.use(Particles, {
    init: async (engine) => {
      await loadSlim(engine);
    },
  });

app.use(pinia)
app.use(router)
app.use(ElementPlus) // 挂载 UI 库

app.mount('#app')
