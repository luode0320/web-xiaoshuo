<template>
  <div class="profile-container">
    <div class="profile-header">
      <div class="avatar">
        <div class="avatar-placeholder">{{ user?.nickname?.charAt(0) || 'U' }}</div>
      </div>
      <div class="user-info">
        <h2>{{ user?.nickname || '未登录用户' }}</h2>
        <p>{{ user?.email }}</p>
        <p>注册时间: {{ formatDate(user?.created_at) }}</p>
      </div>
    </div>
    
    <div class="profile-content">
      <div class="sidebar">
        <el-menu
          :default-active="activeRoute"
          class="sidebar-menu"
          @select="handleMenuSelect"
          :router="false"
        >
          <el-menu-item index="/profile/basic">
            <el-icon><User /></el-icon>
            <span>基本信息</span>
          </el-menu-item>
          <el-menu-item index="/profile/uploads">
            <el-icon><Upload /></el-icon>
            <span>上传历史</span>
          </el-menu-item>
          <el-menu-item index="/profile/comments">
            <el-icon><ChatLineRound /></el-icon>
            <span>我的评论</span>
          </el-menu-item>
          <el-menu-item index="/profile/ratings">
            <el-icon><Star /></el-icon>
            <span>我的评分</span>
          </el-menu-item>
          <el-menu-item index="/profile/social">
            <el-icon><ChatDotSquare /></el-icon>
            <span>社交历史</span>
          </el-menu-item>
          <el-menu-item index="/profile/messages">
            <el-icon><Message /></el-icon>
            <span>系统消息</span>
          </el-menu-item>
          <el-menu-item 
            v-if="userStore.isAdmin" 
            index="/admin/review"
          >
            <el-icon><Setting /></el-icon>
            <span>管理后台</span>
          </el-menu-item>
          <el-menu-item index="/about">
            <el-icon><InfoFilled /></el-icon>
            <span>关于我们</span>
          </el-menu-item>
        </el-menu>
      </div>
      
      <div class="main-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { 
  User, 
  Upload, 
  ChatLineRound, 
  Star, 
  Message, 
  InfoFilled,
  Setting,
  ChatDotSquare
} from '@element-plus/icons-vue'

export default {
  name: 'Profile',
  components: {
    User,
    Upload,
    ChatLineRound,
    Star,
    Message,
    InfoFilled,
    Setting,
    ChatDotSquare
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const userStore = useUserStore()
    
    const activeRoute = ref('')
    
    const user = computed(() => userStore.user)
    
    // 处理菜单选择
    const handleMenuSelect = (index) => {
      if (index === '/profile/messages') {
        // 暂时显示提示，因为消息页面还没创建
        ElMessage.info('系统消息功能正在开发中')
      } else {
        router.push(index)
      }
    }
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
    
    onMounted(() => {
      if (!userStore.isAuthenticated) {
        router.push('/login')
        return
      }
      
      // 设置当前激活的路由
      activeRoute.value = route.path
    })
    
    return {
      activeRoute,
      user,
      handleMenuSelect,
      formatDate,
      userStore
    }
  }
}
</script>

<style scoped>
.profile-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.profile-header {
  display: flex;
  align-items: center;
  padding-bottom: 30px;
  border-bottom: 1px solid #eee;
  margin-bottom: 30px;
}

.avatar {
  margin-right: 20px;
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
}

.user-info h2 {
  margin: 0 0 10px 0;
  color: #333;
}

.user-info p {
  color: #666;
  margin: 5px 0;
}

.profile-content {
  display: flex;
  flex-direction: row;
  gap: 30px;
}

.sidebar {
  width: 220px;
  flex-shrink: 0;
}

.sidebar-menu {
  border-right: none;
}

.main-content {
  flex: 1;
}

@media (max-width: 768px) {
  .profile-container {
    padding: 15px;
  }
  
  .profile-content {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    margin-bottom: 20px;
  }
  
  .profile-header {
    flex-direction: column;
    text-align: center;
  }
  
  .avatar {
    margin-right: 0;
    margin-bottom: 15px;
  }
}
</style>