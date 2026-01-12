import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createPersistedState } from 'pinia-plugin-persistedstate'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './assets/css/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 创建Vue应用
const app = createApp(App)

// 创建Pinia实例并配置持久化
const pinia = createPinia()
pinia.use(createPersistedState())
app.use(pinia)

// 使用路由
app.use(router)

// 使用Element Plus
app.use(ElementPlus)

// 全局注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 挂载应用
app.mount('#app')