<template>
  <div class="uploads-container">
    <div class="page-header">
      <el-button 
        type="text" 
        @click="goBack"
        class="back-button"
      >
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>上传历史</h2>
    </div>
    
    <div class="content">
      <div class="upload-stats">
        <div class="stat-item">
          <div class="stat-number">{{ uploadHistory.length }}</div>
          <div class="stat-label">总共上传</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ approvedCount }}</div>
          <div class="stat-label">已通过</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ pendingCount }}</div>
          <div class="stat-label">审核中</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ rejectedCount }}</div>
          <div class="stat-label">已拒绝</div>
        </div>
      </div>
      
      <div class="novel-list">
        <div 
          v-for="novel in uploadHistory" 
          :key="novel.id" 
          class="novel-item"
        >
          <div class="novel-info">
            <h4>{{ novel.title }}</h4>
            <p>作者: {{ novel.author }}</p>
            <p>字数: {{ novel.word_count }} 字</p>
            <p>文件大小: {{ formatFileSize(novel.file_size) }}</p>
            <p>状态: 
              <el-tag :type="getStatusType(novel.status)" size="small">
                {{ getStatusText(novel.status) }}
              </el-tag>
            </p>
            <p>上传时间: {{ formatDate(novel.created_at) }}</p>
            <p>点击量: {{ novel.click_count }}</p>
          </div>
          <div class="novel-actions">
            <el-button 
              size="small" 
              @click="viewNovel(novel.id)"
            >
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button 
              size="small" 
              type="primary"
              @click="editNovel(novel.id)"
            >
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button 
              v-if="isAllowedToDelete(novel)" 
              size="small" 
              type="danger" 
              @click="deleteNovel(novel.id)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
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
  Edit,
  Delete
} from '@element-plus/icons-vue'

export default {
  name: 'Uploads',
  components: {
    ArrowLeft,
    View,
    Edit,
    Delete
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const uploadHistory = ref([])
    
    const user = computed(() => userStore.user)
    
    const goBack = () => {
      router.go(-1) // 返回上一页
    }
    
    // 获取用户上传的小说
    const fetchUploadHistory = async () => {
      try {
        const response = await apiClient.get(`/api/v1/novels?upload_user_id=${userStore.user?.id}`)
        uploadHistory.value = response.data.data.novels
      } catch (error) {
        console.error('获取上传历史失败:', error)
        ElMessage.error('获取上传历史失败')
      }
    }
    
    // 删除小说
    const deleteNovel = async (novelId) => {
      try {
        await ElMessageBox.confirm(
          '确定要删除这本小说吗？此操作不可恢复。', 
          '删除小说', 
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await apiClient.delete(`/api/v1/novels/${novelId}`, {
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
    
    // 编辑小说
    const editNovel = (novelId) => {
      // 目前没有专门的编辑小说页面，暂时跳转到小说详情
      router.push(`/novel/${novelId}`)
    }
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
    
    // 格式化文件大小
    const formatFileSize = (bytes) => {
      if (bytes === 0) return '0 Bytes';
      const k = 1024;
      const sizes = ['Bytes', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
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
      return novel.upload_user_id === user.value?.id || userStore.isAdmin
    }
    
    // 计算统计信息
    const approvedCount = computed(() => {
      return uploadHistory.value.filter(novel => novel.status === 'approved').length
    })
    
    const pendingCount = computed(() => {
      return uploadHistory.value.filter(novel => novel.status === 'pending').length
    })
    
    const rejectedCount = computed(() => {
      return uploadHistory.value.filter(novel => novel.status === 'rejected').length
    })
    
    onMounted(() => {
      if (!userStore.isAuthenticated) {
        router.push('/login')
        return
      }
      
      fetchUploadHistory()
    })
    
    return {
      uploadHistory,
      user,
      goBack,
      deleteNovel,
      viewNovel,
      editNovel,
      formatDate,
      formatFileSize,
      getStatusType,
      getStatusText,
      isAllowedToDelete,
      approvedCount,
      pendingCount,
      rejectedCount
    }
  }
}
</script>

<style scoped>
.uploads-container {
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

.upload-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
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
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.novel-list {
  max-width: 100%;
}

.novel-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 15px;
  transition: box-shadow 0.3s ease;
}

.novel-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.novel-info h4 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 1.1em;
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
  flex-shrink: 0;
  margin-left: 15px;
}

@media (max-width: 768px) {
  .uploads-container {
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
  
  .upload-stats {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .novel-item {
    flex-direction: column;
    gap: 15px;
  }
  
  .novel-actions {
    flex-direction: row;
    margin-left: 0;
  }
}
</style>