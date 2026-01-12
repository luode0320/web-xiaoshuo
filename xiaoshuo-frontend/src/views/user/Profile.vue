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
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <div class="basic-info">
            <el-form :model="profileForm" :rules="profileRules" ref="profileFormRef" label-width="80px">
              <el-form-item label="昵称" prop="nickname">
                <el-input 
                  v-model="profileForm.nickname" 
                  placeholder="请输入昵称"
                  :disabled="!editing"
                />
              </el-form-item>
              <el-form-item label="邮箱">
                <el-input 
                  v-model="user.email" 
                  disabled
                />
              </el-form-item>
              <el-form-item>
                <el-button 
                  v-if="!editing" 
                  type="primary" 
                  @click="editing = true"
                >
                  编辑资料
                </el-button>
                <template v-else>
                  <el-button @click="cancelEdit">取消</el-button>
                  <el-button type="primary" @click="saveProfile">保存</el-button>
                </template>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="上传历史" name="uploads">
          <div class="novel-list">
            <div 
              v-for="novel in uploadHistory" 
              :key="novel.id" 
              class="novel-item"
            >
              <div class="novel-info">
                <h4>{{ novel.title }}</h4>
                <p>作者: {{ novel.author }}</p>
                <p>状态: 
                  <el-tag :type="getStatusType(novel.status)" size="small">
                    {{ getStatusText(novel.status) }}
                  </el-tag>
                </p>
                <p>上传时间: {{ formatDate(novel.created_at) }}</p>
              </div>
              <div class="novel-actions">
                <el-button 
                  size="small" 
                  @click="viewNovel(novel.id)"
                >
                  查看
                </el-button>
                <el-button 
                  v-if="isAllowedToDelete(novel)" 
                  size="small" 
                  type="danger" 
                  @click="deleteNovel(novel.id)"
                >
                  删除
                </el-button>
              </div>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="我的评论" name="comments">
          <div class="comments-list">
            <div 
              v-for="comment in userComments" 
              :key="comment.id" 
              class="comment-item"
            >
              <div class="comment-header">
                <span class="novel-title">{{ comment.novel?.title }}</span>
                <span class="time">{{ formatDate(comment.created_at) }}</span>
              </div>
              <div class="comment-content">{{ comment.content }}</div>
              <div class="comment-actions">
                <el-button size="small" @click="viewNovel(comment.novel_id)">查看小说</el-button>
              </div>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="我的评分" name="ratings">
          <div class="ratings-list">
            <div 
              v-for="rating in userRatings" 
              :key="rating.id" 
              class="rating-item"
            >
              <div class="rating-header">
                <span class="novel-title">{{ rating.novel?.title }}</span>
                <el-rate 
                  v-model="rating.score" 
                  disabled 
                  :max="10" 
                  show-text
                />
              </div>
              <div class="rating-comment">{{ rating.comment }}</div>
              <div class="rating-actions">
                <el-button size="small" @click="viewNovel(rating.novel_id)">查看小说</el-button>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import dayjs from 'dayjs'

export default {
  name: 'Profile',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const activeTab = ref('basic')
    const editing = ref(false)
    const profileFormRef = ref(null)
    
    const user = computed(() => userStore.user)
    const uploadHistory = ref([])
    const userComments = ref([])
    const userRatings = ref([])
    
    const profileForm = reactive({
      nickname: userStore.user?.nickname || ''
    })
    
    const profileRules = {
      nickname: [
        { required: true, message: '请输入昵称', trigger: 'blur' },
        { max: 20, message: '昵称长度不能超过20个字符', trigger: 'blur' }
      ]
    }
    
    // 获取用户上传的小说
    const fetchUploadHistory = async () => {
      try {
        const response = await axios.get(`/api/v1/novels?upload_user_id=${userStore.user?.id}`)
        uploadHistory.value = response.data.data.novels
      } catch (error) {
        console.error('获取上传历史失败:', error)
      }
    }
    
    // 获取用户评论
    const fetchUserComments = async () => {
      try {
        const response = await axios.get(`/api/v1/comments?user_id=${userStore.user?.id}`)
        userComments.value = response.data.data.comments
      } catch (error) {
        console.error('获取评论失败:', error)
      }
    }
    
    // 获取用户评分
    const fetchUserRatings = async () => {
      try {
        const response = await axios.get(`/api/v1/ratings?user_id=${userStore.user?.id}`)
        userRatings.value = response.data.data.ratings
      } catch (error) {
        console.error('获取评分失败:', error)
      }
    }
    
    // 保存用户资料
    const saveProfile = async () => {
      if (!profileFormRef.value) return
      
      try {
        await profileFormRef.value.validate()
        
        const result = await userStore.updateProfile(profileForm.nickname)
        
        if (result.success) {
          ElMessage.success('资料更新成功')
          editing.value = false
        } else {
          ElMessage.error('资料更新失败')
        }
      } catch (error) {
        console.error('保存资料失败:', error)
        ElMessage.error('保存资料失败')
      }
    }
    
    // 取消编辑
    const cancelEdit = () => {
      editing.value = false
      profileForm.nickname = userStore.user?.nickname || ''
    }
    
    // 删除小说
    const deleteNovel = async (novelId) => {
      try {
        await ElMessageBox.confirm('确定要删除这本小说吗？', '删除小说', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await axios.delete(`/api/v1/novels/${novelId}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('小说删除成功')
        fetchUploadHistory() // 刷新上传历史
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除小说失败:', error)
          ElMessage.error('删除小说失败')
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
    
    // 获取状态类型
    const getStatusType = (status) => {
      switch (status) {
        case 'approved': return 'success'
        case 'pending': return 'warning'
        case 'rejected': return 'danger'
        default: return 'info'
      }
    }
    
    // 获取状态文本
    const getStatusText = (status) => {
      switch (status) {
        case 'approved': return '已通过'
        case 'pending': return '审核中'
        case 'rejected': return '已拒绝'
        default: return status
      }
    }
    
    // 检查是否允许删除小说
    const isAllowedToDelete = (novel) => {
      return novel.upload_user_id === userStore.user?.id || userStore.isAdmin
    }
    
    onMounted(() => {
      if (!userStore.isAuthenticated) {
        router.push('/login')
        return
      }
      
      fetchUploadHistory()
      fetchUserComments()
      fetchUserRatings()
      
      // 初始化表单
      profileForm.nickname = userStore.user?.nickname || ''
    })
    
    return {
      activeTab,
      editing,
      profileForm,
      profileFormRef,
      profileRules,
      user,
      uploadHistory,
      userComments,
      userRatings,
      saveProfile,
      cancelEdit,
      deleteNovel,
      viewNovel,
      formatDate,
      getStatusType,
      getStatusText,
      isAllowedToDelete
    }
  }
}
</script>

<style scoped>
.profile-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
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

.basic-info {
  max-width: 500px;
}

.novel-list {
  max-width: 800px;
}

.novel-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
}

.novel-info h4 {
  margin: 0 0 10px 0;
  color: #333;
}

.novel-info p {
  margin: 5px 0;
  color: #666;
  font-size: 0.9rem;
}

.novel-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.comments-list, .ratings-list {
  max-width: 800px;
}

.comment-item, .rating-item {
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
}

.comment-header, .rating-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.novel-title {
  font-weight: 500;
  color: #333;
}

.time {
  color: #999;
  font-size: 0.9rem;
}

.comment-content, .rating-comment {
  margin-bottom: 15px;
  line-height: 1.6;
  color: #333;
}

.comment-actions, .rating-actions {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }
  
  .avatar {
    margin-right: 0;
    margin-bottom: 15px;
  }
  
  .novel-item {
    flex-direction: column;
    gap: 15px;
  }
  
  .comment-header, .rating-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
}