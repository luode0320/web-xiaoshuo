<template>
  <div class="ratings-container">
    <div class="page-header">
      <el-button 
        type="text" 
        @click="goBack"
        class="back-button"
      >
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>我的评分</h2>
    </div>
    
    <div class="content">
      <div class="ratings-stats">
        <div class="stat-item">
          <div class="stat-number">{{ userRatings.length }}</div>
          <div class="stat-label">总评分次</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ averageScore }}</div>
          <div class="stat-label">平均分</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ novelsRated.length }}</div>
          <div class="stat-label">评分过的小说</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ highRatedCount }}</div>
          <div class="stat-label">高分(8分以上)</div>
        </div>
      </div>
      
      <div class="ratings-list">
        <div 
          v-for="rating in userRatings" 
          :key="rating.id" 
          class="rating-item"
        >
          <div class="rating-main">
            <div class="rating-header">
              <div class="novel-info">
                <el-tag type="info" size="small" class="novel-tag">
                  {{ rating.novel?.title || '未知小说' }}
                </el-tag>
                <el-rate 
                  v-model="rating.score" 
                  disabled 
                  :max="10" 
                  show-text
                  class="rating-stars"
                />
              </div>
              <div class="rating-actions">
                <el-button 
                  size="small" 
                  @click="viewNovel(rating.novel_id)"
                >
                  <el-icon><View /></el-icon>
                  查看小说
                </el-button>
                <el-button 
                  size="small" 
                  type="danger"
                  @click="deleteRating(rating.id)"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </div>
            <div v-if="rating.comment" class="rating-comment">
              <strong>评价:</strong> {{ rating.comment }}
            </div>
            <div class="rating-time">
              评分时间: {{ formatDate(rating.created_at) }}
            </div>
          </div>
        </div>
        
        <div v-if="userRatings.length === 0" class="no-ratings">
          <el-empty description="暂无评分" />
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
  name: 'Ratings',
  components: {
    ArrowLeft,
    View,
    Delete
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const userRatings = ref([])
    
    const goBack = () => {
      router.go(-1) // 返回上一页
    }
    
    // 获取用户评分
    const fetchUserRatings = async () => {
      try {
        const response = await apiClient.get(`/api/v1/users/ratings`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        userRatings.value = response.data.data.data || response.data.data || []
      } catch (error) {
        console.error('获取评分失败:', error)
        ElMessage.error('获取评分失败')
        userRatings.value = []
      }
    }
    
    // 删除评分
    const deleteRating = async (ratingId) => {
      try {
        await ElMessageBox.confirm(
          '确定要删除这条评分吗？此操作不可恢复。', 
          '删除评分', 
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await apiClient.delete(`/api/v1/ratings/${ratingId}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('评分删除成功')
        fetchUserRatings() // 刷新评分列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除评分失败:', error)
          ElMessage.error('删除评分失败')
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
    
    // 计算平均分
    const averageScore = computed(() => {
      if (userRatings.value.length === 0) return '0.0'
      const total = userRatings.value.reduce((sum, rating) => sum + rating.score, 0)
      const avg = total / userRatings.value.length
      return avg.toFixed(1)
    })
    
    // 计算评分过的小说数量
    const novelsRated = computed(() => {
      const novelIds = new Set()
      userRatings.value.forEach(rating => {
        if (rating.novel_id) {
          novelIds.add(rating.novel_id)
        }
      })
      return Array.from(novelIds)
    })
    
    // 计算高分(8分以上)数量
    const highRatedCount = computed(() => {
      return userRatings.value.filter(rating => rating.score >= 8).length
    })
    
    onMounted(() => {
      if (!userStore.isAuthenticated) {
        router.push('/login')
        return
      }
      
      fetchUserRatings()
    })
    
    return {
      userRatings,
      goBack,
      deleteRating,
      viewNovel,
      formatDate,
      averageScore,
      novelsRated,
      highRatedCount
    }
  }
}
</script>

<style scoped>
.ratings-container {
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

.ratings-stats {
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
  color: #67c23a;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.ratings-list {
  max-width: 100%;
}

.rating-item {
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
  transition: box-shadow 0.3s ease;
}

.rating-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.rating-main {
  width: 100%;
}

.rating-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
  flex-wrap: wrap;
  gap: 10px;
}

.novel-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
  flex: 1;
  min-width: 200px;
}

.novel-tag {
  font-size: 0.9em;
}

.rating-stars :deep(.el-rate__text) {
  margin-left: 5px;
  color: #606266;
  font-size: 14px;
}

.rating-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.rating-comment {
  padding: 10px 15px;
  background: #f0f9ff;
  border-radius: 6px;
  margin-bottom: 10px;
  border-left: 4px solid #409eff;
  color: #333;
}

.rating-time {
  font-size: 0.85rem;
  color: #999;
}

.no-ratings {
  text-align: center;
  padding: 40px 0;
  color: #999;
}

@media (max-width: 768px) {
  .ratings-container {
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
  
  .ratings-stats {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .rating-header {
    flex-direction: column;
    align-items: stretch;
    gap: 15px;
  }
  
  .novel-info {
    align-items: flex-start;
  }
  
  .rating-actions {
    align-self: flex-start;
  }
}
</style>