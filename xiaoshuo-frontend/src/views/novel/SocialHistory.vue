<template>
  <div class="social-history-container">
    <div class="social-history-header">
      <h1>社交活动历史</h1>
      <p>查看您的评论、评分和互动历史</p>
    </div>
    
    <div class="social-history-content">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="我的评论" name="comments">
          <div class="comments-history">
            <div 
              v-for="comment in comments" 
              :key="comment.id" 
              class="comment-item"
            >
              <div class="comment-header">
                <div class="novel-info">
                  <h4 @click="viewNovel(comment.novel_id)" class="novel-title">
                    {{ comment.novel?.title || '未知小说' }}
                  </h4>
                  <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                </div>
                <div class="comment-actions">
                  <el-button 
                    size="small" 
                    type="primary" 
                    @click="editComment(comment)"
                    v-if="canEditComment(comment)"
                  >
                    编辑
                  </el-button>
                  <el-button 
                    size="small" 
                    type="danger" 
                    @click="deleteComment(comment.id)"
                    v-if="canDeleteComment(comment)"
                  >
                    删除
                  </el-button>
                </div>
              </div>
              <div class="comment-content">{{ comment.content }}</div>
              <div class="comment-stats">
                <span class="like-count">
                  <i class="el-icon-thumb"></i> {{ comment.like_count || 0 }} 个赞
                </span>
                <span class="reply-count" v-if="comment.reply_count">
                  <i class="el-icon-chat-line-round"></i> {{ comment.reply_count }} 条回复
                </span>
              </div>
            </div>
            
            <!-- 分页 -->
            <div class="pagination" v-if="commentPagination.total > commentPagination.limit">
              <el-pagination
                v-model:current-page="commentPagination.page"
                :page-size="commentPagination.limit"
                :total="commentPagination.total"
                @current-change="handleCommentPageChange"
                layout="prev, pager, next, jumper"
              />
            </div>
            
            <div v-if="comments.length === 0" class="empty-state">
              <p>暂无评论记录</p>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="我的评分" name="ratings">
          <div class="ratings-history">
            <div 
              v-for="rating in ratings" 
              :key="rating.id" 
              class="rating-item"
            >
              <div class="rating-header">
                <div class="novel-info">
                  <h4 @click="viewNovel(rating.novel_id)" class="novel-title">
                    {{ rating.novel?.title || '未知小说' }}
                  </h4>
                  <div class="rating-score">
                    <el-rate 
                      v-model="rating.score" 
                      disabled 
                      :max="10" 
                      show-text
                    />
                  </div>
                  <span class="rating-time">{{ formatDate(rating.created_at) }}</span>
                </div>
                <div class="rating-actions">
                  <el-button 
                    size="small" 
                    type="primary" 
                    @click="editRating(rating)"
                    v-if="canEditRating(rating)"
                  >
                    编辑
                  </el-button>
                  <el-button 
                    size="small" 
                    type="danger" 
                    @click="deleteRating(rating.id)"
                    v-if="canDeleteRating(rating)"
                  >
                    删除
                  </el-button>
                </div>
              </div>
              <div class="rating-comment" v-if="rating.comment">{{ rating.comment }}</div>
              <div class="rating-stats">
                <span class="like-count">
                  <i class="el-icon-thumb"></i> {{ rating.like_count || 0 }} 个赞
                </span>
              </div>
            </div>
            
            <!-- 分页 -->
            <div class="pagination" v-if="ratingPagination.total > ratingPagination.limit">
              <el-pagination
                v-model:current-page="ratingPagination.page"
                :page-size="ratingPagination.limit"
                :total="ratingPagination.total"
                @current-change="handleRatingPageChange"
                layout="prev, pager, next, jumper"
              />
            </div>
            
            <div v-if="ratings.length === 0" class="empty-state">
              <p>暂无评分记录</p>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="我的点赞" name="likes">
          <el-tabs v-model="likeTab" class="sub-tabs">
            <el-tab-pane label="点赞的评论" name="liked-comments">
              <div class="liked-comments">
                <div 
                  v-for="like in likedComments" 
                  :key="like.id" 
                  class="like-item"
                >
                  <div class="like-header">
                    <div class="novel-info">
                      <h4 @click="viewNovel(like.comment?.novel_id)" class="novel-title">
                        {{ like.comment?.novel?.title || '未知小说' }}
                      </h4>
                      <span class="like-time">{{ formatDate(like.created_at) }}</span>
                    </div>
                  </div>
                  <div class="like-content">评论: {{ like.comment?.content }}</div>
                  <div class="like-author">作者: {{ like.comment?.user?.nickname }}</div>
                </div>
                
                <div v-if="likedComments.length === 0" class="empty-state">
                  <p>暂无点赞评论记录</p>
                </div>
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="点赞的评分" name="liked-ratings">
              <div class="liked-ratings">
                <div 
                  v-for="like in likedRatings" 
                  :key="like.id" 
                  class="like-item"
                >
                  <div class="like-header">
                    <div class="novel-info">
                      <h4 @click="viewNovel(like.rating?.novel_id)" class="novel-title">
                        {{ like.rating?.novel?.title || '未知小说' }}
                      </h4>
                      <div class="rating-preview">
                        <el-rate 
                          v-model="like.rating?.score" 
                          disabled 
                          :max="10" 
                        />
                        <span class="rating-text">{{ like.rating?.score }} 分</span>
                      </div>
                      <span class="like-time">{{ formatDate(like.created_at) }}</span>
                    </div>
                  </div>
                  <div class="like-content" v-if="like.rating?.comment">评价: {{ like.rating?.comment }}</div>
                  <div class="like-author">作者: {{ like.rating?.user?.nickname }}</div>
                </div>
                
                <div v-if="likedRatings.length === 0" class="empty-state">
                  <p>暂无点赞评分记录</p>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-tab-pane>
        
        <el-tab-pane label="互动通知" name="notifications">
          <div class="notifications">
            <div 
              v-for="notification in notifications" 
              :key="notification.id" 
              class="notification-item"
              :class="{ 'unread': !notification.is_read }"
            >
              <div class="notification-header">
                <div class="notification-type">
                  <el-tag 
                    :type="getNotificationType(notification.type)" 
                    size="small"
                  >
                    {{ getNotificationText(notification.type) }}
                  </el-tag>
                  <span class="notification-time">{{ formatDate(notification.created_at) }}</span>
                </div>
                <el-button 
                  size="small" 
                  @click="markAsRead(notification.id)"
                  v-if="!notification.is_read"
                >
                  标记已读
                </el-button>
              </div>
              <div class="notification-content" @click="handleNotificationClick(notification)">
                {{ notification.content }}
              </div>
            </div>
            
            <!-- 分页 -->
            <div class="pagination" v-if="notificationPagination.total > notificationPagination.limit">
              <el-pagination
                v-model:current-page="notificationPagination.page"
                :page-size="notificationPagination.limit"
                :total="notificationPagination.total"
                @current-change="handleNotificationPageChange"
                layout="prev, pager, next, jumper"
              />
            </div>
            
            <div v-if="notifications.length === 0" class="empty-state">
              <p>暂无互动通知</p>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
    
    <!-- 评论编辑弹窗 -->
    <el-dialog title="编辑评论" v-model="commentEditDialogVisible" width="600px">
      <el-form :model="editingComment" :rules="commentRules" ref="commentFormRef">
        <el-form-item label="评论内容" prop="content">
          <el-input 
            v-model="editingComment.content" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入评论内容"
            maxlength="1000"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="commentEditDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveComment" :loading="commentSaving">保存</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 评分编辑弹窗 -->
    <el-dialog title="编辑评分" v-model="ratingEditDialogVisible" width="600px">
      <el-form :model="editingRating" :rules="ratingRules" ref="ratingFormRef">
        <el-form-item label="评分" prop="score">
          <el-rate v-model="editingRating.score" :max="10" show-text allow-half />
        </el-form-item>
        <el-form-item label="评价" prop="comment">
          <el-input 
            v-model="editingRating.comment" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入您的评价..."
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="ratingEditDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveRating" :loading="ratingSaving">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import dayjs from 'dayjs'

export default {
  name: 'SocialHistory',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const activeTab = ref('comments')
    const likeTab = ref('liked-comments')
    
    // 数据
    const comments = ref([])
    const ratings = ref([])
    const likedComments = ref([])
    const likedRatings = ref([])
    const notifications = ref([])
    
    // 分页
    const commentPagination = ref({
      page: 1,
      limit: 10,
      total: 0
    })
    
    const ratingPagination = ref({
      page: 1,
      limit: 10,
      total: 0
    })
    
    const notificationPagination = ref({
      page: 1,
      limit: 10,
      total: 0
    })
    
    // 编辑相关
    const commentEditDialogVisible = ref(false)
    const ratingEditDialogVisible = ref(false)
    const editingComment = ref({})
    const editingRating = ref({})
    const commentFormRef = ref(null)
    const ratingFormRef = ref(null)
    const commentSaving = ref(false)
    const ratingSaving = ref(false)
    
    // 验证规则
    const commentRules = {
      content: [
        { required: true, message: '请输入评论内容', trigger: 'blur' },
        { max: 1000, message: '评论内容不能超过1000个字符', trigger: 'blur' }
      ]
    }
    
    const ratingRules = {
      score: [
        { required: true, message: '请评分', trigger: 'change' }
      ],
      comment: [
        { max: 500, message: '评价不能超过500个字符', trigger: 'blur' }
      ]
    }
    
    // 检查用户认证
    if (!userStore.isAuthenticated) {
      router.push('/login')
      ElMessage.error('请先登录')
      return
    }
    
    // 获取我的评论
    const fetchComments = async () => {
      try {
        const response = await axios.get(`/api/v1/users/comments?page=${commentPagination.value.page}&limit=${commentPagination.value.limit}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        comments.value = response.data.data.data
        commentPagination.value.total = response.data.data.total
      } catch (error) {
        console.error('获取评论失败:', error)
        ElMessage.error('获取评论失败')
      }
    }
    
    // 获取我的评分
    const fetchRatings = async () => {
      try {
        const response = await axios.get(`/api/v1/users/ratings?page=${ratingPagination.value.page}&limit=${ratingPagination.value.limit}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        ratings.value = response.data.data.data
        ratingPagination.value.total = response.data.data.total
      } catch (error) {
        console.error('获取评分失败:', error)
        ElMessage.error('获取评分失败')
      }
    }
    
    // 获取点赞的评论
    const fetchLikedComments = async () => {
      try {
        const response = await axios.get('/api/v1/users/liked-comments', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        likedComments.value = response.data.data
      } catch (error) {
        console.error('获取点赞评论失败:', error)
        ElMessage.error('获取点赞评论失败')
      }
    }
    
    // 获取点赞的评分
    const fetchLikedRatings = async () => {
      try {
        const response = await axios.get('/api/v1/users/liked-ratings', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        likedRatings.value = response.data.data
      } catch (error) {
        console.error('获取点赞评分失败:', error)
        ElMessage.error('获取点赞评分失败')
      }
    }
    
    // 获取互动通知
    const fetchNotifications = async () => {
      try {
        const response = await axios.get(`/api/v1/users/notifications?page=${notificationPagination.value.page}&limit=${notificationPagination.value.limit}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        notifications.value = response.data.data.data
        notificationPagination.value.total = response.data.data.total
      } catch (error) {
        console.error('获取通知失败:', error)
        ElMessage.error('获取通知失败')
      }
    }
    
    // 删除评论
    const deleteComment = async (commentId) => {
      try {
        await ElMessageBox.confirm('确定要删除这条评论吗？', '删除确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await axios.delete(`/api/v1/comments/${commentId}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('评论删除成功')
        fetchComments() // 刷新评论列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除评论失败:', error)
          ElMessage.error('删除评论失败')
        }
      }
    }
    
    // 编辑评论
    const editComment = (comment) => {
      editingComment.value = { ...comment }
      commentEditDialogVisible.value = true
    }
    
    // 保存评论
    const saveComment = async () => {
      if (!commentFormRef.value) return
      
      try {
        await commentFormRef.value.validate()
        
        commentSaving.value = true
        
        await axios.put(`/api/v1/comments/${editingComment.value.id}`, {
          content: editingComment.value.content
        }, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('评论更新成功')
        commentEditDialogVisible.value = false
        fetchComments() // 刷新评论列表
      } catch (error) {
        console.error('更新评论失败:', error)
        ElMessage.error('更新评论失败')
      } finally {
        commentSaving.value = false
      }
    }
    
    // 删除评分
    const deleteRating = async (ratingId) => {
      try {
        await ElMessageBox.confirm('确定要删除这条评分吗？', '删除确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await axios.delete(`/api/v1/ratings/${ratingId}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('评分删除成功')
        fetchRatings() // 刷新评分列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除评分失败:', error)
          ElMessage.error('删除评分失败')
        }
      }
    }
    
    // 编辑评分
    const editRating = (rating) => {
      editingRating.value = { ...rating }
      ratingEditDialogVisible.value = true
    }
    
    // 保存评分
    const saveRating = async () => {
      if (!ratingFormRef.value) return
      
      try {
        await ratingFormRef.value.validate()
        
        ratingSaving.value = true
        
        await axios.put(`/api/v1/ratings/${editingRating.value.id}`, {
          score: editingRating.value.score,
          comment: editingRating.value.comment
        }, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('评分更新成功')
        ratingEditDialogVisible.value = false
        fetchRatings() // 刷新评分列表
      } catch (error) {
        console.error('更新评分失败:', error)
        ElMessage.error('更新评分失败')
      } finally {
        ratingSaving.value = false
      }
    }
    
    // 标记通知为已读
    const markAsRead = async (notificationId) => {
      try {
        await axios.post(`/api/v1/users/notifications/${notificationId}/read`, {}, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        // 更新本地数据
        const notification = notifications.value.find(n => n.id === notificationId)
        if (notification) {
          notification.is_read = true
        }
      } catch (error) {
        console.error('标记已读失败:', error)
        ElMessage.error('标记已读失败')
      }
    }
    
    // 处理通知点击
    const handleNotificationClick = (notification) => {
      // 根据通知类型跳转到相应页面
      if (notification.target_type === 'comment' && notification.target_id) {
        // 可能需要跳转到小说详情页并定位到评论
        router.push(`/novel/${notification.novel_id}`)
      } else if (notification.target_type === 'rating' && notification.target_id) {
        router.push(`/novel/${notification.novel_id}`)
      }
    }
    
    // 查看小说
    const viewNovel = (novelId) => {
      router.push(`/novel/${novelId}`)
    }
    
    // 分页处理
    const handleCommentPageChange = (page) => {
      commentPagination.value.page = page
      fetchComments()
    }
    
    const handleRatingPageChange = (page) => {
      ratingPagination.value.page = page
      fetchRatings()
    }
    
    const handleNotificationPageChange = (page) => {
      notificationPagination.value.page = page
      fetchNotifications()
    }
    
    // 标签切换处理
    const handleTabChange = (tabName) => {
      switch (tabName) {
        case 'comments':
          fetchComments()
          break
        case 'ratings':
          fetchRatings()
          break
        case 'likes':
          fetchLikedComments()
          fetchLikedRatings()
          break
        case 'notifications':
          fetchNotifications()
          break
      }
    }
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
    
    // 获取通知类型
    const getNotificationType = (type) => {
      switch (type) {
        case 'comment_reply': return 'primary'
        case 'comment_like': return 'success'
        case 'rating_like': return 'warning'
        case 'mention': return 'danger'
        default: return 'info'
      }
    }
    
    // 获取通知文本
    const getNotificationText = (type) => {
      switch (type) {
        case 'comment_reply': return '评论回复'
        case 'comment_like': return '评论点赞'
        case 'rating_like': return '评分点赞'
        case 'mention': return '提及'
        default: return type
      }
    }
    
    // 检查是否可以编辑评论
    const canEditComment = (comment) => {
      return userStore.isAuthenticated && comment.user_id === userStore.userId
    }
    
    // 检查是否可以删除评论
    const canDeleteComment = (comment) => {
      return userStore.isAuthenticated && (comment.user_id === userStore.userId || userStore.isAdmin)
    }
    
    // 检查是否可以编辑评分
    const canEditRating = (rating) => {
      return userStore.isAuthenticated && rating.user_id === userStore.userId
    }
    
    // 检查是否可以删除评分
    const canDeleteRating = (rating) => {
      return userStore.isAuthenticated && (rating.user_id === userStore.userId || userStore.isAdmin)
    }
    
    onMounted(() => {
      // 初始加载评论数据
      fetchComments()
      fetchRatings()
      fetchLikedComments()
      fetchLikedRatings()
      fetchNotifications()
    })
    
    return {
      activeTab,
      likeTab,
      comments,
      ratings,
      likedComments,
      likedRatings,
      notifications,
      commentPagination,
      ratingPagination,
      notificationPagination,
      commentEditDialogVisible,
      ratingEditDialogVisible,
      editingComment,
      editingRating,
      commentFormRef,
      ratingFormRef,
      commentSaving,
      ratingSaving,
      commentRules,
      ratingRules,
      fetchComments,
      fetchRatings,
      fetchLikedComments,
      fetchLikedRatings,
      fetchNotifications,
      deleteComment,
      editComment,
      saveComment,
      deleteRating,
      editRating,
      saveRating,
      markAsRead,
      handleNotificationClick,
      viewNovel,
      handleCommentPageChange,
      handleRatingPageChange,
      handleNotificationPageChange,
      handleTabChange,
      formatDate,
      getNotificationType,
      getNotificationText,
      canEditComment,
      canDeleteComment,
      canEditRating,
      canDeleteRating
    }
  }
}
</script>

<style scoped>
.social-history-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.social-history-header {
  text-align: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.social-history-header h1 {
  font-size: 1.8rem;
  color: #333;
  margin-bottom: 10px;
}

.social-history-header p {
  color: #666;
}

.comment-item, .rating-item, .like-item, .notification-item {
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
  background: white;
}

.comment-header, .rating-header, .like-header, .notification-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.novel-info {
  flex: 1;
}

.novel-title {
  font-weight: 500;
  color: #409eff;
  cursor: pointer;
  margin: 0 0 5px 0;
}

.novel-title:hover {
  text-decoration: underline;
}

.comment-time, .rating-time, .like-time, .notification-time {
  color: #999;
  font-size: 0.9rem;
}

.comment-content, .rating-comment, .like-content, .notification-content {
  margin-bottom: 10px;
  line-height: 1.6;
  color: #333;
}

.comment-stats, .rating-stats {
  display: flex;
  gap: 15px;
  font-size: 0.9rem;
  color: #666;
}

.like-count, .reply-count {
  display: flex;
  align-items: center;
  gap: 3px;
}

.rating-score {
  margin: 5px 0;
}

.rating-preview {
  display: flex;
  align-items: center;
  gap: 5px;
  margin: 5px 0;
}

.rating-text {
  font-size: 0.9rem;
  color: #666;
}

.like-author, .notification-type {
  font-size: 0.9rem;
  color: #999;
}

.notification-item.unread {
  border-left: 4px solid #409eff;
  background-color: #f0f9ff;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sub-tabs {
  margin-top: 15px;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #999;
}

.pagination {
  text-align: center;
  margin-top: 20px;
}

.dialog-footer {
  text-align: right;
}

@media (max-width: 768px) {
  .social-history-container {
    padding: 15px;
  }
  
  .comment-header, .rating-header, .like-header, .notification-header {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }
  
  .comment-actions, .rating-actions {
    align-self: flex-start;
  }
  
  .social-history-header h1 {
    font-size: 1.5rem;
  }
}
</style>