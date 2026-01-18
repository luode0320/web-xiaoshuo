<template>
  <div class="profile-container">
    <div class="profile-header">
      <div class="avatar">
        <div class="avatar-placeholder">{{ user?.nickname?.charAt(0) || 'U' }}</div>
      </div>
      <div class="user-info">
        <h2>{{ user?.nickname || '未登录用户' }}</h2>
        <!-- <p>{{ user?.email }}</p> -->
        <p>注册时间: {{ formatDate(user?.created_at) }}</p>
      </div>
    </div>

    <div class="profile-content">
      <div class="sidebar">
        <el-menu :default-active="activeRoute" class="sidebar-menu" @select="handleMenuSelect" :router="true">
          <el-menu-item index="/profile/basic" :route="true">
            <el-icon>
              <User />
            </el-icon>
            <span>基本信息</span>
          </el-menu-item>
          <el-menu-item index="/profile/uploads" :route="true">
            <el-icon>
              <Upload />
            </el-icon>
            <span>上传历史</span>
          </el-menu-item>
          <el-menu-item index="/profile/comments" :route="true">
            <el-icon>
              <ChatLineRound />
            </el-icon>
            <span>我的评论</span>
          </el-menu-item>
          <el-menu-item index="/profile/ratings" :route="true">
            <el-icon>
              <Star />
            </el-icon>
            <span>我的评分</span>
          </el-menu-item>
          <el-menu-item index="/profile/social" :route="true">
            <el-icon>
              <ChatDotSquare />
            </el-icon>
            <span>社交历史</span>
          </el-menu-item>
          <el-menu-item index="/profile/messages" :route="true">
            <el-icon>
              <Message />
            </el-icon>
            <span>系统消息</span>
          </el-menu-item>
          <el-menu-item v-if="userStore.isAdmin" index="/admin/review" :route="true">
            <el-icon>
              <Setting />
            </el-icon>
            <span>管理后台</span>
          </el-menu-item>
          <el-menu-item index="/about" :route="true">
            <el-icon>
              <InfoFilled />
            </el-icon>
            <span>关于我们</span>
          </el-menu-item>
        </el-menu>
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
      router.push(index)
    }

    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }

    onMounted(async () => {
      // 初始化用户状态，检查本地存储的token并获取用户信息
      await userStore.initializeUser()

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
  flex: 1;
  min-height: calc(100vh - 140px);
  /* 减去顶部可能的导航栏和底部导航栏的高度 */
  height: 100%;
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
  flex: 1;
  flex-direction: row;
  gap: 30px;
  min-height: 0;
  min-width: 0;
  /* 防止内容溢出 */
}

.sidebar {
  width: 220px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}

.sidebar-menu {
  border-right: none;
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .profile-container {
    padding: 15px;
    min-height: calc(100vh - 100px);
    /* 移动端底部导航栏高度 */
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