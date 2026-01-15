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
          <el-tag 
            :type="getStatusType(novel.status)" 
            size="small"
            class="status-tag"
          >
            {{ getStatusText(novel.status) }}
          </el-tag>
          <!-- 上传频率提示（仅对当前用户上传的小说显示） -->
          <el-tag 
            v-if="isAuthenticated && novel.upload_user_id === userId && uploadFrequency.remaining_count !== undefined"
            :type="uploadFrequency.remaining_count > 2 ? 'success' : uploadFrequency.remaining_count > 0 ? 'warning' : 'danger'"
            size="small"
            class="upload-frequency-tag"
          >
            今日剩余上传: {{ uploadFrequency.remaining_count }}
          </el-tag>
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
        
        <!-- 分类和关键词设置（当用户阅读进度达到20%以上时显示） -->
        <div v-if="showClassificationSection" class="classification-section">
          <h3>设置分类和关键词</h3>
          <div class="current-classification">
            <p><strong>当前分类:</strong></p>
            <el-tag 
              v-for="category in novel.categories" 
              :key="category.id" 
              type="info"
              closable
              @close="removeCategory(category.id)"
            >
              {{ category.name }}
            </el-tag>
            <div v-if="novel.categories.length === 0" class="no-categories">
              尚未设置分类
            </div>
          </div>
          <div class="current-keywords">
            <p><strong>当前关键词:</strong></p>
            <el-tag 
              v-for="keyword in novel.keywords" 
              :key="keyword.id" 
              type="success"
              closable
              @close="removeKeyword(keyword.id)"
            >
              {{ keyword.keyword }}
            </el-tag>
            <div v-if="novel.keywords.length === 0" class="no-keywords">
              尚未设置关键词
            </div>
          </div>
          <div class="classification-inputs">
            <el-select 
              v-model="selectedCategory" 
              placeholder="选择分类" 
              style="width: 200px; margin-right: 10px;"
            >
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
            <el-input
              v-model="newKeyword"
              placeholder="输入关键词并按回车添加"
              style="width: 200px;"
              @keyup.enter="addKeyword"
            />
            <el-button @click="saveClassification" type="primary" :loading="savingClassification">
              保存设置
            </el-button>
          </div>
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
              
              <!-- 评分列表 -->
              <div class="ratings-list">
                <div 
                  v-for="rating in ratings" 
                  :key="rating.id" 
                  class="rating-item"
                >
                  <div class="rating-header">
                    <span class="username">{{ rating.user?.nickname || '匿名用户' }}</span>
                    <el-rate 
                      v-model="rating.score" 
                      disabled 
                      :max="10" 
                      show-text
                    />
                    <span class="time">{{ formatDate(rating.created_at) }}</span>
                  </div>
                  <div class="rating-comment">{{ rating.comment }}</div>
                  <div class="rating-actions">
                    <el-button 
                      size="small" 
                      @click="likeRating(rating.id)"
                      :class="{ 'liked': isRatingLiked(rating.id) }"
                    >
                      <i :class="isRatingLiked(rating.id) ? 'el-icon-thumb' : 'el-icon-thumb'"></i> 
                      {{ rating.like_count || 0 }}
                    </el-button>
                  </div>
                </div>
              </div>      </div>
      
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
              <el-button 
                type="text" 
                size="small" 
                @click="likeComment(comment.id)"
                :class="{ 'liked': isCommentLiked(comment.id) }"
              >
                <i :class="isCommentLiked(comment.id) ? 'el-icon-thumb' : 'el-icon-thumb'"></i> 
                {{ comment.like_count || 0 }}
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 小说操作历史区域（仅对上传者或管理员显示） -->
    <div v-if="isAuthenticated && novel && (novel.upload_user_id === userId || userStore.isAdmin)" class="history-section">
      <el-collapse>
        <el-collapse-item title="小说操作历史" name="history">
          <div class="history-content">
            <h4>管理日志</h4>
            <div v-if="novelHistory.activity_history && novelHistory.activity_history.admin_logs.length > 0" class="admin-logs">
              <div 
                v-for="log in novelHistory.activity_history.admin_logs" 
                :key="log.id" 
                class="log-item"
              >
                <div class="log-header">
                  <span class="log-action">{{ log.action }}</span>
                  <span class="log-time">{{ formatDate(log.created_at) }}</span>
                </div>
                <div class="log-details">{{ log.details }}</div>
                <div class="log-user" v-if="log.admin_user">
                  操作员: {{ log.admin_user.nickname }}
                </div>
              </div>
            </div>
            <div v-else class="no-logs">暂无管理日志</div>
            
            <h4 style="margin-top: 20px;">评分历史</h4>
            <div v-if="novelHistory.activity_history && novelHistory.activity_history.ratings.length > 0" class="ratings-history">
              <div 
                v-for="rating in novelHistory.activity_history.ratings" 
                :key="rating.id" 
                class="rating-item"
              >
                <div class="rating-header">
                  <span class="username">{{ rating.user?.nickname || '匿名用户' }}</span>
                  <el-rate 
                    v-model="rating.score" 
                    disabled 
                    :max="10" 
                    show-text
                  />
                  <span class="time">{{ formatDate(rating.created_at) }}</span>
                </div>
                <div class="rating-comment">{{ rating.comment }}</div>
              </div>
            </div>
            <div v-else class="no-ratings">暂无评分记录</div>
            
            <h4 style="margin-top: 20px;">评论历史</h4>
            <div v-if="novelHistory.activity_history && novelHistory.activity_history.comments.length > 0" class="comments-history">
              <div 
                v-for="comment in novelHistory.activity_history.comments" 
                :key="comment.id" 
                class="comment-item"
              >
                <div class="comment-header">
                  <span class="username">{{ comment.user?.nickname || '匿名用户' }}</span>
                  <span class="time">{{ formatDate(comment.created_at) }}</span>
                </div>
                <div class="comment-content">{{ comment.content }}</div>
              </div>
            </div>
            <div v-else class="no-comments">暂无评论记录</div>
          </div>
        </el-collapse-item>
      </el-collapse>
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
import apiClient from '@/utils/api'
import dayjs from 'dayjs'

export default {
  name: 'NovelDetail',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const userStore = useUserStore()
    
    const novel = ref(null)
    const comments = ref([])
    const ratings = ref([]) // 添加评分列表
    const newComment = ref('')
    const commentLoading = ref(false)
    const ratingDialogVisible = ref(false)
    const ratingLoading = ref(false)
    const ratingFormRef = ref(null)
    
    // 点赞状态管理
    const likedComments = ref(new Set())
    const likedRatings = ref(new Set())
    
    const ratingForm = reactive({
      score: 0,
      comment: ''
    })
    
    // 分类和关键词设置相关
    const showClassificationSection = ref(false)
    const categories = ref([])
    const selectedCategory = ref(null)
    const newKeyword = ref('')
    const savingClassification = ref(false)
    
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
    const userId = computed(() => userStore.userId)
    const avgRating = computed(() => {
      // 简单计算平均分，实际应该从API获取
      return novel.value?.avg_rating || 0
    })
    const ratingCount = computed(() => {
      return novel.value?.rating_count || 0
    })
    
    // 检查是否显示分类设置区域
    const checkShowClassificationSection = async () => {
      if (!isAuthenticated.value || !novel.value) {
        showClassificationSection.value = false
        return
      }
      
      try {
        // 获取用户阅读进度
        const response = await apiClient.get(`/api/v1/novels/${route.params.id}/progress`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        // 检查进度是否达到20%以上
        const progress = response.data.data
        if (progress && progress.progress && progress.progress > 20) {
          showClassificationSection.value = true
        } else {
          showClassificationSection.value = false
        }
      } catch (error) {
        // 如果获取进度失败，仍然可以不显示（或者按其他逻辑处理）
        showClassificationSection.value = false
      }
    }
    
    // 检查评论是否已点赞
    const isCommentLiked = (commentId) => {
      return likedComments.value.has(commentId)
    }
    
    // 检查评分是否已点赞
    const isRatingLiked = (ratingId) => {
      return likedRatings.value.has(ratingId)
    }
    
    // 获取分类列表
    const fetchCategories = async () => {
      try {
        const response = await apiClient.get('/api/v1/categories')
        categories.value = response.data.data
      } catch (error) {
        console.error('获取分类列表失败:', error)
      }
    }
    
    // 获取小说详情
    const fetchNovelDetail = async () => {
      try {
        const response = await apiClient.get(`/api/v1/novels/${route.params.id}`)
        novel.value = response.data.data
        // 如果响应中包含评分信息，更新avgRating和ratingCount
        if (response.data.data.avg_rating !== undefined) {
          // 这里avgRating通过计算属性获取，不需要额外设置
        }
        
        // 检查是否显示分类设置区域
        await checkShowClassificationSection()
      } catch (error) {
        console.error('获取小说详情失败:', error)
        ElMessage.error('获取小说详情失败')
      }
    }
    
    // 获取评论列表
    const fetchComments = async () => {
      try {
        const response = await apiClient.get(`/api/v1/comments?novel_id=${route.params.id}`)
        comments.value = response.data.data.comments
      } catch (error) {
        console.error('获取评论失败:', error)
      }
    }
    
    // 获取评分列表
    const fetchRatings = async () => {
      try {
        const response = await apiClient.get(`/api/v1/ratings/novel/${route.params.id}`)
        ratings.value = response.data.data.ratings
      } catch (error) {
        console.error('获取评分失败:', error)
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
        
        await apiClient.post('/api/v1/comments', {
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
        
        await apiClient.post('/api/v1/ratings', {
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
        if (isCommentLiked(commentId)) {
          // 取消点赞
          await apiClient.delete(`/api/v1/comments/${commentId}/like`, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
          likedComments.value.delete(commentId)
        } else {
          // 点赞
          await apiClient.post(`/api/v1/comments/${commentId}/like`, {}, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
          likedComments.value.add(commentId)
        }
        
        // 刷新评论列表以获取最新的点赞数
        fetchComments()
      } catch (error) {
        console.error('操作点赞失败:', error)
        ElMessage.error('操作点赞失败')
      }
    }
    
    // 点赞评分
    const likeRating = async (ratingId) => {
      if (!isAuthenticated.value) {
        router.push('/login')
        return
      }
      
      try {
        if (isRatingLiked(ratingId)) {
          // 取消点赞
          await apiClient.delete(`/api/v1/ratings/${ratingId}/like`, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
          likedRatings.value.delete(ratingId)
        } else {
          // 点赞
          await apiClient.post(`/api/v1/ratings/${ratingId}/like`, {}, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
          likedRatings.value.add(ratingId)
        }
        
        // 刷新小说信息以获取最新的评分
        fetchNovelDetail()
      } catch (error) {
        console.error('操作点赞失败:', error)
        ElMessage.error('操作点赞失败')
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
    
    // 添加关键词
    const addKeyword = () => {
      if (!newKeyword.value.trim()) {
        ElMessage.warning('请输入关键词')
        return
      }
      
      // 检查是否已存在
      const exists = novel.value.keywords.some(k => k.keyword === newKeyword.value.trim())
      if (!exists) {
        novel.value.keywords.push({
          id: Date.now(), // 临时ID，保存时会更新
          keyword: newKeyword.value.trim()
        })
      }
      newKeyword.value = ''
    }
    
    // 移除分类
    const removeCategory = (categoryId) => {
      novel.value.categories = novel.value.categories.filter(cat => cat.id !== categoryId)
    }
    
    // 移除关键词
    const removeKeyword = (keywordId) => {
      novel.value.keywords = novel.value.keywords.filter(kw => kw.id !== keywordId)
    }
    
    // 保存分类和关键词设置
    const saveClassification = async () => {
      if (!selectedCategory.value && novel.value.keywords.length === 0) {
        ElMessage.warning('请至少设置一个分类或关键词')
        return
      }
      
      try {
        savingClassification.value = true
        
        const payload = {
          category_id: selectedCategory.value,
          keywords: novel.value.keywords.map(kw => kw.keyword)
        }
        
        await apiClient.post(`/api/v1/novels/${route.params.id}/classify`, payload, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('分类和关键词设置成功')
        selectedCategory.value = null
      } catch (error) {
        console.error('保存分类和关键词失败:', error)
        ElMessage.error('保存分类和关键词失败')
      } finally {
        savingClassification.value = false
      }
    }
    
    // 获取上传频率信息
    const fetchUploadFrequency = async () => {
      if (!isAuthenticated.value) return
      
      try {
        const response = await apiClient.get('/api/v1/novels/upload-frequency', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        uploadFrequency.value = response.data.data
      } catch (error) {
        console.error('获取上传频率失败:', error)
      }
    }
    
    // 获取小说状态
    const fetchNovelStatus = async () => {
      if (!isAuthenticated.value || !novel.value) return
      
      try {
        const response = await apiClient.get(`/api/v1/novels/${route.params.id}/status`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        novelStatus.value = response.data.data.status
      } catch (error) {
        console.error('获取小说状态失败:', error)
      }
    }
    
    // 获取小说操作历史
    const fetchNovelHistory = async () => {
      if (!isAuthenticated.value || !novel.value) return
      
      try {
        const response = await apiClient.get(`/api/v1/novels/${route.params.id}/history`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        novelHistory.value = response.data.data
      } catch (error) {
        console.error('获取小说历史失败:', error)
      }
    }
    
    // 获取状态类型
    const getStatusType = (status) => {
      switch (status) {
        case 'pending':
          return 'warning'
        case 'approved':
          return 'success'
        case 'rejected':
          return 'danger'
        default:
          return 'info'
      }
    }
    
    // 获取状态文本
    const getStatusText = (status) => {
      switch (status) {
        case 'pending':
          return '待审核'
        case 'approved':
          return '已通过'
        case 'rejected':
          return '已拒绝'
        default:
          return status
      }
    }
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD')
    }
    
    // 上传频率信息
    const uploadFrequency = ref({})
    const novelStatus = ref({})
    const novelHistory = ref({})
    
    onMounted(async () => {
      await fetchCategories()
      await fetchNovelDetail()
      await fetchComments()
      await fetchRatings()
      await fetchUploadFrequency()
      await fetchNovelStatus()
      // 操作历史只对特定用户显示，按需获取
    })
    
    return {
      novel,
      comments,
      ratings,
      newComment,
      commentLoading,
      ratingDialogVisible,
      ratingLoading,
      ratingForm,
      ratingFormRef,
      ratingRules,
      likedComments,
      likedRatings,
      showClassificationSection,
      categories,
      selectedCategory,
      newKeyword,
      savingClassification,
      isAuthenticated,
      avgRating,
      ratingCount,
      isCommentLiked,
      isRatingLiked,
      checkShowClassificationSection,
      addKeyword,
      removeCategory,
      removeKeyword,
      saveClassification,
      submitComment,
      showRatingDialog,
      submitRating,
      likeComment,
      likeRating,
      startReading,
      downloadNovel,
      formatDate,
      uploadFrequency,
      novelStatus,
      novelHistory,
      getStatusType,
      getStatusText,
      fetchUploadFrequency,
      fetchNovelStatus,
      fetchNovelHistory
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

.rating-item {
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
}

.rating-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.username {
  font-weight: 500;
  color: #333;
}

.time {
  color: #999;
  font-size: 0.9rem;
}

.rating-comment {
  margin-bottom: 15px;
  line-height: 1.6;
  color: #333;
}

.rating-actions {
  display: flex;
  justify-content: flex-end;
}

.classification-section {
  margin-top: 25px;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
  border: 1px solid #eee;
}

.classification-section h3 {
  margin-bottom: 15px;
  color: #333;
}

.current-classification,
.current-keywords {
  margin-bottom: 15px;
}

.current-classification .el-tag,
.current-keywords .el-tag {
  margin-right: 8px;
  margin-bottom: 5px;
}

.no-categories,
.no-keywords {
  color: #999;
  font-size: 0.9rem;
  margin-top: 5px;
}

.classification-inputs {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.stats .status-tag {
  margin-left: 10px;
}

.stats .upload-frequency-tag {
  margin-left: 10px;
}

.history-section {
  margin-top: 30px;
  padding: 20px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #eee;
}

.history-content h4 {
  margin-bottom: 10px;
  color: #333;
  border-bottom: 1px solid #eee;
  padding-bottom: 5px;
}

.log-item {
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 4px;
  margin-bottom: 10px;
  background: white;
}

.log-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}

.log-action {
  font-weight: 500;
  color: #409eff;
}

.log-time {
  color: #999;
  font-size: 0.9rem;
}

.log-details {
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 5px;
}

.log-user {
  color: #999;
  font-size: 0.8rem;
}

.no-logs, .no-ratings, .no-comments {
  color: #999;
  font-style: italic;
  text-align: center;
  padding: 20px;
}

.rating-item, .comment-item {
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 4px;
  margin-bottom: 10px;
  background: white;
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
  
  .stats .status-tag,
  .stats .upload-frequency-tag {
    margin-left: 0;
    margin-top: 5px;
  }
  
  .classification-inputs {
    flex-direction: column;
    align-items: stretch;
  }
  
  .classification-inputs .el-select,
  .classification-inputs .el-input {
    width: 100% !important;
  }
}
</style>