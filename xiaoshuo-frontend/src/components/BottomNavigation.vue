<template>
  <div class="app-layout">
    <div class="main-content">
      <router-view />
    </div>
    <div class="bottom-nav">
      <div 
        v-for="item in navItems" 
        :key="item.path"
        class="nav-item"
        :class="{ active: isActive(item.path) }"
        @click="navigateTo(item)"
      >
        <div class="nav-icon">{{ item.icon }}</div>
        <div class="nav-text">{{ item.text }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

export default {
  name: 'AppLayout',
  setup() {
    const router = useRouter()
    const route = useRoute()
    const userStore = useUserStore()

    const navItems = computed(() => {
      const items = [
        {
          path: '/',
          text: '首页',
          icon: '🏠'
        },
        {
          path: '/category',
          text: '分类',
          icon: '📚'
        },
        {
          path: '/ranking',
          text: '排行榜',
          icon: '🏆'
        },
        {
          path: '/upload',
          text: '上传',
          icon: '📤',
          requiresAuth: true
        }
      ];
      
      // 根据用户登录状态添加用户相关导航项
      if (userStore.isAuthenticated) {
        items.push({
          path: '/profile',
          text: '用户',
          icon: '👤'
        });
      } else {
        items.push({
          path: '/login',
          text: '登录',
          icon: '🔒'
        });
      }
      
      // 只显示不需要认证或用户已认证的项目
      return items.filter(item => !item.requiresAuth || userStore.isAuthenticated)
    })

    const isActive = (path) => {
      if (path === '/') {
        return route.path === '/'
      }
      return route.path.startsWith(path)
    }

    const navigateTo = (item) => {
      if (item.path === route.path) return
      
      // 如果是需要认证的页面但用户未登录，跳转到登录页
      if (item.requiresAuth && !userStore.isAuthenticated) {
        router.push('/login')
        return
      }
      
      router.push(item.path)
    }

    return {
      navItems,
      isActive,
      navigateTo
    }
  }
}
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
}

.main-content {
  flex: 1;
  padding-bottom: 60px; /* 为底部导航栏留出空间 */
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: #ffffff;
  border-top: 1px solid #e0e0e0;
  display: flex;
  z-index: 1000;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #999999;
  font-size: 12px;
}

.nav-item:hover {
  background-color: #f5f5f5;
}

.nav-item.active {
  color: #1890ff;
}

.nav-icon {
  font-size: 20px;
  margin-bottom: 2px;
}

.nav-text {
  font-size: 12px;
}

@media (min-width: 769px) {
  .bottom-nav {
    display: none;
  }
  
  .main-content {
    padding-bottom: 0;
  }
}
</style>