<template>
  <div id="app">
    <BottomNavigation />
  </div>
</template>

<script>
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import BottomNavigation from './components/BottomNavigation.vue'

export default {
  name: 'App',
  components: {
    BottomNavigation
  },
  setup() {
    const userStore = useUserStore()
    
    onMounted(async () => {
      // 在应用启动时初始化用户状态
      await userStore.initializeUser()
    })
  }
}
</script>

<style>
html, body {
  height: 100%;
  margin: 0;
  padding: 0;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  height: 100%;
  padding-bottom: 60px; /* 为底部导航栏预留空间 */
}

@media (min-width: 769px) {
  #app {
    padding-bottom: 0;
  }
}
</style>