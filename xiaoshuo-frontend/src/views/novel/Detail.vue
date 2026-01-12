<template>
  <div class="novel-detail">
    <div v-if="novel" class="novel-header">
      <div class="novel-cover">
        <div class="cover-placeholder">小说封面</div>
      </div>
      <div class="novel-info">
        <h1>{{ novel.title }}</h1>
        <p class="author">作者: {{ novel.author }}</p>
        <p class="protagonist" v-if="novel.protagonist">主角: {{ novel.protagonist }}</p>
        <div class="stats">
          <span class="clicks">点击: {{ novel.click_count }}</span>
          <span class="word-count">字数: {{ novel.word_count }}万</span>
          <span class="upload-time">上传: {{ formatDate(novel.created_at) }}</span>
        </div>
        <div class="categories">
          <el-tag 
            v-for="category in novel.categories" 
            :key="category.id" 
            type="info" 
            size="small"
            class="category-tag"
          >
            {{ category.name }}
          </el-tag>
        </div>
        <div class="description">
          <h3>简介</h3>
          <p>{{ novel.description }}</p>
        </div>
        <div class="actions">
          <el-button type="primary" @click="startReading">开始阅读</el-button>
          <el-button v-if="isAuthenticated" @click="showRatingDialog">评分</el-button>
          <el-button @click="downloadNovel">下载</el-button>
        </div>
      </div>
    </div>
    
    <div class="novel-content">
      <!-- 评分区域 -->
      <div class="rating-section">
        <h3>评分与评价</h3>
        <div class="rating-summary">
          <div class="avg-rating">
            <span class="score">{{ avgRating }}</span>
            <span class="text">分</span>
          </div>
          <div class="rating-details">
            <p>共有 {{ ratingCount }} 人评分</p>
          </div>
        </div>
      </div>
      
      <!-- 评论区域 -->
      <div class="comments-section">
        <h3>评论区</h3>
        <div v-if="isAuthenticated" class="comment-input">
          <el-input 
            v-model="newComment" 
            type="textarea" 
            :rows="4" 
            placeholder="发表您的评论..."
            maxlength="1000"
            show-word-limit
          />
          <el-button type="primary" @click="submitComment" :loading="commentLoading" class="submit-btn">
            发表评论
          </el-button>
        </div>
        
        <div class="comments-list">
          <div 
            v-for="comment in comments" 
            :key="comment.id" 
            class="comment-item"
          >
            <div class="comment-header">
              <span class="username">{{ comment.user.nickname }}</span>
              <span class="time">{{ formatDate(comment.created_at) }}</span>
            </div>
            <div class="comment-content">{{ comment.content }}</div>
            <div class="comment-actions">
              <el-button type="text" size="small">回复</el-button>
              <el-button type="text" size="small" @click="likeComment(comment.id)">
                <i class="el-icon-thumb"></i> {{ comment.like_count }}
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 评分对话框 -->
    <el-dialog v-model="ratingDialogVisible" title="为小说评分" width="500px">
      <el-form :model="ratingForm" :rules="ratingRules" ref="ratingFormRef">
        <el-form-item label="评分" prop="score">
          <el-rate v-model="ratingForm.score" :max="10" show-text allow-half />
        </el-form-item>
        <el-form-item label="评价" prop="comment">
          <el-input 
            v-model="ratingForm.comment" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入您的评价..."
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ratingDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRating" :loading="ratingLoading">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import dayjs from 'dayjs'

export default {
  name: 'NovelDetail',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const userStore = useUserStore()
    
    const novel = ref(null)
    const comments = ref([])
    const newComment = ref('')
    const commentLoading = ref(false)
    const ratingDialogVisible = ref(false)
    const ratingLoading = ref(false)
    const ratingFormRef = ref(null)
    
    const ratingForm = reactive({
      score: 0,
      comment: ''
    })
    
    const ratingRules = {
      score: [
        { required: true, message: '请评分', trigger: 'change' }
      ],
      comment: [
        { max: 500, message: '评价不能超过500个字符', trigger: 'blur' }
      ]
    }
    
    // 计算属性
    const isAuthenticated = computed(() => userStore.isAuthenticated)
    const avgRating = computed(() => {
      // 简单计算平均分，实际应该从API获取
      return novel.value?.avg_rating || 0
    })
    const ratingCount = computed(() => {
      return novel.value?.rating_count || 0
    })
    
    // 获取小说详情
    const fetchNovelDetail = async () => {
      try {
        const response = await axios.get(`/api/v1/novels/${route.params.id}`)
        novel.value = response.data.data
      } catch (error) {
        console.error('获取小说详情失败:', error)
        ElMessage.error('获取小说详情失败')
      }
    }
    
    // 获取评论列表
    const fetchComments = async () => {
      try {
        const response = await axios.get(`/api/v1/comments?novel_id=${route.params.id}`)
        comments.value = response.data.data.comments
      } catch (error) {
        console.error('获取评论失败:', error)
      }
    }
    
    // 提交评论
    const submitComment = async () => {
      if (!newComment.value.trim()) {
        ElMessage.warning('评论内容不能为空')
        return
      }
      
      try {
        commentLoading.value = true
        
        await axios.post('/api/v1/comments', {
          content: newComment.value,
          novel_id: parseInt(route.params.id)
        }, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('评论发表成功')
        newComment.value = ''
        fetchComments() // 刷新评论列表
      } catch (error) {
        console.error('发表评论失败:', error)
        ElMessage.error('发表评论失败')
      } finally {
        commentLoading.value = false
      }
    }
    
    // 显示评分对话框
    const showRatingDialog = () => {
      if (!isAuthenticated.value) {
        router.push('/login')
        return
      }
      ratingDialogVisible.value = true
    }
    
    // 提交评分
    const submitRating = async () => {
      if (!ratingFormRef.value) return
      
      try {
        await ratingFormRef.value.validate()
        
        ratingLoading.value = true
        
        await axios.post('/api/v1/ratings', {
          score: ratingForm.score,
          comment: ratingForm.comment,
          novel_id: parseInt(route.params.id)
        }, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('评分成功')
        ratingDialogVisible.value = false
        // 重置表单
        ratingForm.score = 0
        ratingForm.comment = ''
        fetchNovelDetail() // 刷新小说信息
      } catch (error) {
        console.error('评分失败:', error)
        ElMessage.error('评分失败')
      } finally {
        ratingLoading.value = false
      }
    }
    
    // 点赞评论
    const likeComment = async (commentId) => {
      if (!isAuthenticated.value) {
        router.push('/login')
        return
      }
      
      try {
        await axios.post(`/api/v1/comments/${commentId}/like`, {}, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('点赞成功')
        fetchComments() // 刷新评论列表
      } catch (error) {
        console.error('点赞失败:', error)
        ElMessage.error('点赞失败')
      }
    }
    
    // 开始阅读
    const startReading = () => {
      if (!isAuthenticated.value) {
        router.push('/login')
        return
      }
      router.push(`/read/${route.params.id}`)
    }
    
    // 下载小说
    const downloadNovel = () => {
      // 实现下载功能
      window.open(`/api/v1/novels/${route.params.id}/download`, '_blank')
    }
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD')
    }
    
    onMounted(() => {
      fetchNovelDetail()
      fetchComments()
    })
    
    return {
      novel,
      comments,
      newComment,
      commentLoading,
      ratingDialogVisible,
      ratingLoading,
      ratingForm,
      ratingFormRef,
      ratingRules,
      isAuthenticated,
      avgRating,
      ratingCount,
      submitComment,
      showRatingDialog,
      submitRating,
      likeComment,
      startReading,
      downloadNovel,
      formatDate
    }
  }
}
</script>

<style scoped>
.novel-detail {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.novel-header {
  display: flex;
  gap: 30px;
  padding-bottom: 30px;
  border-bottom: 1px solid #eee;
  margin-bottom: 30px;
}

.novel-cover {
  flex-shrink: 0;
  width: 200px;
  height: 280px;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
}

.cover-placeholder {
  font-size: 0.9rem;
  color: #999;
}

.novel-info {
  flex: 1;
}

.novel-info h1 {
  font-size: 1.8rem;
  margin-bottom: 15px;
  color: #333;
}

.author, .protagonist {
  font-size: 1rem;
  color: #666;
  margin-bottom: 10px;
}

.stats {
  display: flex;
  gap: 20px;
  margin-bottom: 15px;
  font-size: 0.9rem;
  color: #999;
}

.categories {
  margin-bottom: 20px;
}

.category-tag {
  margin-right: 10px;
}

.description {
  margin-bottom: 25px;
}

.description h3 {
  margin-bottom: 10px;
  color: #333;
}

.description p {
  color: #666;
  line-height: 1.6;
}

.actions {
  display: flex;
  gap: 15px;
  margin-top: 20px;
}

.rating-section {
  margin-bottom: 40px;
}

.rating-section h3 {
  margin-bottom: 20px;
  color: #333;
  font-size: 1.2rem;
}

.rating-summary {
  display: flex;
  align-items: center;
  gap: 30px;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
}

.avg-rating {
  display: flex;
  align-items: baseline;
  gap: 5px;
}

.avg-rating .score {
  font-size: 2rem;
  font-weight: bold;
  color: #409eff;
}

.avg-rating .text {
  font-size: 1rem;
  color: #666;
}

.comments-section h3 {
  margin-bottom: 20px;
  color: #333;
  font-size: 1.2rem;
}

.comment-input {
  margin-bottom: 30px;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
}

.submit-btn {
  margin-top: 15px;
}

.comments-list {
  margin-top: 20px;
}

.comment-item {
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  font-size: 0.9rem;
  color: #999;
}

.comment-content {
  margin-bottom: 15px;
  line-height: 1.6;
  color: #333;
}

.comment-actions {
  display: flex;
  gap: 15px;
}

@media (max-width: 768px) {
  .novel-header {
    flex-direction: column;
  }
  
  .novel-cover {
    width: 150px;
    height: 210px;
    margin: 0 auto;
  }
  
  .rating-summary {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .stats {
    flex-direction: column;
    gap: 5px;
  }
}