<template>
  <div class="comments-container">
    <div class="page-header">
      <el-button 
        type="text" 
        @click="goBack"
        class="back-button"
      >
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>我的评论</h2>
    </div>
    
    <div class="content">
      <div class="comments-stats">
        <div class="stat-item">
          <div class="stat-number">{{ userComments.length }}</div>
          <div class="stat-label">总评论数</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ recentComments.length }}</div>
          <div class="stat-label">近期评论</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ novelsCommentedOn.length }}</div>
          <div class="stat-label">评论过的小说</div>
        </div>
      </div>
      
      <div class="comments-list">
        <div 
          v-for="comment in userComments" 
          :key="comment.id" 
          class="comment-item"
        >
          <div class="comment-main">
            <div class="comment-header">
              <div class="novel-info">
                <el-tag type="info" size="small" class="novel-tag">
                  {{ comment.novel?.title || '未知小说' }}
                </el-tag>
                <span class="time">{{ formatDate(comment.created_at) }}</span>
              </div>
              <div class="comment-actions">
                <el-button 
                  size="small" 
                  @click="viewNovel(comment.novel_id)"
                >
                  <el-icon><View /></el-icon>
                  查看小说
                </el-button>
                <el-button 
                  size="small" 
                  type="danger"
                  @click="deleteComment(comment.id)"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </div>
            <div class="comment-content">
              {{ comment.content }}
            </div>
          </div>
        </div>
        
        <div v-if="userComments.length === 0" class="no-comments">
          <el-empty description="暂无评论" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'
import { 
  ArrowLeft,
  View,
  Delete
} from '@element-plus/icons-vue'

export default {
  name: 'Comments',
  components: {
    ArrowLeft,
    View,
    Delete
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const userComments = ref([])
    
    const goBack = () => {
      router.go(-1) // 返回上一页
    }
    
    // 获取用户评论
    const fetchUserComments = async () => {
      try {
        const response = await apiClient.get(`/api/v1/users/comments`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        userComments.value = response.data.data.data || response.data.data || []
      } catch (error) {
        console.error('获取评论失败:', error)
        ElMessage.error('获取评论失败')
        userComments.value = []
      }
    }
    
    // 删除评论
    const deleteComment = async (commentId) => {
      try {
        await ElMessageBox.confirm(
          '确定要删除这条评论吗？此操作不可恢复。', 
          '删除评论', 
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await apiClient.delete(`/api/v1/comments/${commentId}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('评论删除成功')
        fetchUserComments() // 刷新评论列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除评论失败:', error)
          ElMessage.error('删除评论失败')
        }
      }
    }
    
    // 查看小说
    const viewNovel = (novelId) => {
      router.push(`/novel/${novelId}`)
    }
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
    
    // 计算最近的评论（最近30天）
    const recentComments = computed(() => {
      const thirtyDaysAgo = dayjs().subtract(30, 'day')
      return userComments.value.filter(comment => 
        dayjs(comment.created_at).isAfter(thirtyDaysAgo)
      )
    })
    
    // 计算评论过的小说数量
    const novelsCommentedOn = computed(() => {
      const novelIds = new Set()
      userComments.value.forEach(comment => {
        if (comment.novel_id) {
          novelIds.add(comment.novel_id)
        }
      })
      return Array.from(novelIds)
    })
    
    onMounted(() => {
      if (!userStore.isAuthenticated) {
        router.push('/login')
        return
      }
      
      fetchUserComments()
    })
    
    return {
      userComments,
      goBack,
      deleteComment,
      viewNovel,
      formatDate,
      recentComments,
      novelsCommentedOn
    }
  }
}
</script>

<style scoped>
.comments-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.page-header {
  display: flex;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.back-button {
  margin-right: 15px;
  font-size: 16px;
}

.page-header h2 {
  margin: 0;
  color: #333;
}

.content {
  flex: 1;
}

.comments-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #e6a23c;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.comments-list {
  max-width: 100%;
}

.comment-item {
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
  transition: box-shadow 0.3s ease;
}

.comment-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.comment-main {
  width: 100%;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
  flex-wrap: wrap;
  gap: 10px;
}

.novel-info {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.novel-tag {
  font-size: 0.9em;
}

.time {
  color: #999;
  font-size: 0.9rem;
}

.comment-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.comment-content {
  padding: 15px;
  background: #f8f9fa;
  border-radius: 6px;
  line-height: 1.6;
  color: #333;
  border-left: 4px solid #409eff;
}

.no-comments {
  text-align: center;
  padding: 40px 0;
  color: #999;
}

@media (max-width: 768px) {
  .comments-container {
    padding: 15px;
    margin: 10px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .back-button {
    align-self: flex-start;
    margin-bottom: 15px;
  }
  
  .comments-stats {
    grid-template-columns: 1fr;
  }
  
  .comment-header {
    flex-direction: column;
    align-items: stretch;
    gap: 15px;
  }
  
  .novel-info {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .comment-actions {
    align-self: flex-start;
  }
}
</style>