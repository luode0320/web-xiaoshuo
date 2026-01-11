import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import persistedstate from 'pinia-plugin-persistedstate'
import App from './App.vue'
import routes from './router'
import './assets/css/index.css'

const app = createApp(App)

const router = createRouter({
  history: createWebHistory(),
  routes
})

const pinia = createPinia()
pinia.use(persistedstate)

app.use(router)
app.use(pinia)
app.use(ElementPlus)

app.mount('#app')