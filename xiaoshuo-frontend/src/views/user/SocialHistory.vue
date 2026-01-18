<template>
  <div class="social-container">
    <div class="page-header">
      <el-button 
        type="text" 
        @click="goBack"
        class="back-button"
      >
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>社交历史</h2>
    </div>
    
    <div class="content">
      <div class="social-summary">
        <h3>社交活动概览</h3>
        <div class="summary-stats">
          <div class="stat-item">
            <div class="stat-number">{{ socialStats.totalComments }}</div>
            <div class="stat-label">总评论数</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ socialStats.totalRatings }}</div>
            <div class="stat-label">总评分次</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ socialStats.totalLikes }}</div>
            <div class="stat-label">获赞总数</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ socialStats.totalInteractions }}</div>
            <div class="stat-label">互动数</div>
          </div>
        </div>
      </div>
      
      <div class="social-timeline">
        <h3>活动时间线</h3>
        
        <div class="timeline-filters">
          <el-radio-group v-model="activeFilter" @change="filterActivities">
            <el-radio-button label="all">全部</el-radio-button>
            <el-radio-button label="comments">评论</el-radio-button>
            <el-radio-button label="ratings">评分</el-radio-button>
            <el-radio-button label="likes">点赞</el-radio-button>
          </el-radio-group>
        </div>
        
        <div class="timeline">
          <div 
            v-for="activity in filteredActivities" 
            :key="activity.id" 
            class="timeline-item"
          >
            <div class="timeline-dot"></div>
            <div class="timeline-content">
              <div class="activity-header">
                <div class="activity-type">
                  <el-tag 
                    :type="getActivityType(activity.type).tagType" 
                    size="small"
                  >
                    {{ getActivityType(activity.type).label }}
                  </el-tag>
                  <span class="activity-time">{{ formatDate(activity.created_at) }}</span>
                </div>
                            <el-button 
                              size="small" 
                              @click="viewNovel(activity.novel_id)"
                              :disabled="!activity.novel_id"
                            >
                              <el-icon><View /></el-icon>
                              查看小说
                            </el-button>              </div>
              <div class="activity-content">
                <div v-if="activity.type === 'comment'">
                  <strong>评论内容:</strong> {{ activity.content }}
                </div>
                <div v-else-if="activity.type === 'rating'">
                  <strong>评分:</strong> 
                  <el-rate 
                    v-model="activity.score" 
                    disabled 
                    :max="10" 
                    show-text
                  />
                  <div v-if="activity.comment">评价: {{ activity.comment }}</div>
                </div>
                <div v-else-if="activity.type === 'like'">
                  <strong>点赞了{{ activity.target_type === 'comment' ? '评论' : '评分' }}</strong>
                </div>
              </div>
            </div>
          </div>
          
          <div v-if="filteredActivities.length === 0" class="no-activities">
            <el-empty description="暂无社交活动" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'
import { 
  ArrowLeft,
  View
} from '@element-plus/icons-vue'

export default {
  name: 'SocialHistory',
  components: {
    ArrowLeft,
    View
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const socialStats = ref({
      totalComments: 0,
      totalRatings: 0,
      totalLikes: 0,
      totalInteractions: 0
    })
    
    const allActivities = ref([])
    const filteredActivities = ref([])
    const activeFilter = ref('all')
    
    const goBack = () => {
      router.go(-1) // 返回上一页
    }
    
    // 获取社交统计
    const fetchSocialStats = async () => {
      try {
        const response = await apiClient.get(`/api/v1/users/social-stats`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        socialStats.value = response.data.data || {
          totalComments: 0,
          totalRatings: 0,
          totalLikes: 0,
          totalInteractions: 0
        }
      } catch (error) {
        console.error('获取社交统计失败:', error)
        ElMessage.error('获取社交统计失败')
        // 设置默认值
        socialStats.value = {
          totalComments: 0,
          totalRatings: 0,
          totalLikes: 0,
          totalInteractions: 0
        }
      }
    }
    
    // 获取社交活动
    const fetchSocialActivities = async () => {
      try {
        // 这里需要从多个API获取不同类型的数据
        const [commentsRes, ratingsRes] = await Promise.allSettled([
          apiClient.get(`/api/v1/users/comments`, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          }),
          apiClient.get(`/api/v1/users/ratings`, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
        ]);
        
        const comments = commentsRes.status === 'fulfilled' ? 
          (commentsRes.value.data.data.data || commentsRes.value.data.data || []) : [];
        const ratings = ratingsRes.status === 'fulfilled' ? 
          (ratingsRes.value.data.data.data || ratingsRes.value.data.data || []) : [];
        
        // 将评论和评分合并为统一的活动列表
        const activities = [];
        
        comments.forEach(comment => {
          activities.push({
            id: `comment-${comment.id}`,
            type: 'comment',
            content: comment.content,
            novel_id: comment.novel_id,
            created_at: comment.created_at,
            novel: comment.novel
          });
        });
        
        ratings.forEach(rating => {
          activities.push({
            id: `rating-${rating.id}`,
            type: 'rating',
            score: rating.score,
            comment: rating.comment,
            novel_id: rating.novel_id,
            created_at: rating.created_at,
            novel: rating.novel
          });
        });
        
        // 按时间倒序排列
        allActivities.value = activities.sort((a, b) => 
          new Date(b.created_at) - new Date(a.created_at)
        );
        
        filteredActivities.value = [...allActivities.value];
      } catch (error) {
        console.error('获取社交活动失败:', error)
        ElMessage.error('获取社交活动失败')
        allActivities.value = []
        filteredActivities.value = []
      }
    }
    
    // 获取活动类型信息
    const getActivityType = (type) => {
      switch (type) {
        case 'comment':
          return { label: '评论', tagType: 'primary' };
        case 'rating':
          return { label: '评分', tagType: 'warning' };
        case 'like':
          return { label: '点赞', tagType: 'danger' };
        default:
          return { label: type, tagType: 'info' };
      }
    }
    
    // 过滤活动
    const filterActivities = () => {
      if (activeFilter.value === 'all') {
        filteredActivities.value = [...allActivities.value];
      } else {
        filteredActivities.value = allActivities.value.filter(activity => 
          activity.type === activeFilter.value
        );
      }
    }
    
    // 查看小说
    const viewNovel = (novelId) => {
      if (novelId) {
        router.push(`/novel/${novelId}`)
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
      
      fetchSocialStats()
      fetchSocialActivities()
    })
    
    return {
      socialStats,
      allActivities,
      filteredActivities,
      activeFilter,
      goBack,
      getActivityType,
      filterActivities,
      viewNovel,
      formatDate
    }
  }
}
</script>

<style scoped>
.social-container {
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

.social-summary {
  margin-bottom: 30px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.social-summary h3 {
  margin-top: 0;
  color: #333;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.summary-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.stat-item {
  text-align: center;
  padding: 15px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.social-timeline h3 {
  color: #333;
  margin-bottom: 15px;
}

.timeline-filters {
  margin-bottom: 20px;
}

.timeline {
  position: relative;
  padding-left: 30px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 15px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #e4e7ed;
}

.timeline-item {
  position: relative;
  margin-bottom: 25px;
}

.timeline-dot {
  position: absolute;
  left: -35px;
  top: 10px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #409eff;
  border: 2px solid white;
  box-shadow: 0 0 0 2px #409eff;
}

.timeline-content {
  padding: 15px;
  background: white;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  flex-wrap: wrap;
  gap: 10px;
}

.activity-type {
  display: flex;
  align-items: center;
  gap: 10px;
}

.activity-time {
  color: #999;
  font-size: 0.9rem;
}

.activity-content {
  color: #333;
  line-height: 1.6;
}

.activity-content :deep(.el-rate) {
  display: inline-block;
  margin-left: 10px;
}

.no-activities {
  text-align: center;
  padding: 40px 0;
  color: #999;
}

@media (max-width: 768px) {
  .social-container {
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
  
  .timeline {
    padding-left: 20px;
  }
  
  .timeline::before {
    left: 8px;
  }
  
  .timeline-dot {
    left: -22px;
  }
  
  .activity-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .summary-stats {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>