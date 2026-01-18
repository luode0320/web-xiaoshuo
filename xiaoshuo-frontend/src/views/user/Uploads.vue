<template>
  <div class="uploads-container">
    <div class="header">
      <el-button 
        type="primary" 
        link 
        @click="goBack"
        class="back-button"
      >
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>上传历史</h2>
    </div>
    
    <div class="content">
      <el-table 
        :data="uploads" 
        style="width: 100%"
        v-loading="loading"
      >
        <el-table-column prop="title" label="小说标题" />
        <el-table-column prop="author" label="作者" width="120" />
        <el-table-column prop="word_count" label="字数" width="100" />
        <el-table-column prop="click_count" label="点击量" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag 
              :type="getStatusType(row.status)"
              size="small"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button 
              size="small" 
              @click="viewNovel(row.id)"
              type="primary"
            >
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'

export default {
  name: 'Uploads',
  components: {
    ArrowLeft
  },
  setup() {
    const router = useRouter()
    const uploads = ref([])
    const loading = ref(false)
    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)
    
    // 获取上传历史
    const fetchUploads = async () => {
      loading.value = true
      try {
        const response = await apiClient.get('/api/v1/novels', {
          params: {
            page: currentPage.value,
            limit: pageSize.value,
            upload_user_id: 'current' // 这里需要根据实际API调整
          }
        })
        
        if (response.data.code === 200) {
          uploads.value = response.data.data.novels
          total.value = response.data.data.pagination.total
        } else {
          ElMessage.error('获取上传历史失败: ' + response.data.message)
        }
      } catch (error) {
        console.error('获取上传历史失败:', error)
        ElMessage.error('获取上传历史失败: ' + error.message)
      } finally {
        loading.value = false
      }
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
        case 'pending': return '待审核'
        case 'rejected': return '已拒绝'
        default: return status
      }
    }
    
    // 查看小说
    const viewNovel = (id) => {
      router.push(`/novel/${id}`)
    }
    
    // 返回上一页
    const goBack = () => {
      router.push('/profile')
    }
    
    // 处理页面大小变化
    const handleSizeChange = (size) => {
      pageSize.value = size
      currentPage.value = 1
      fetchUploads()
    }
    
    // 处理当前页变化
    const handleCurrentChange = (page) => {
      currentPage.value = page
      fetchUploads()
    }
    
    onMounted(() => {
      fetchUploads()
    })
    
    return {
      uploads,
      loading,
      currentPage,
      pageSize,
      total,
      formatDate,
      getStatusType,
      getStatusText,
      viewNovel,
      goBack,
      handleSizeChange,
      handleCurrentChange
    }
  }
}
</script>

<style scoped>
.uploads-container {
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  min-height: 100%;
}

.header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.header h2 {
  margin: 0;
  margin-left: 15px;
  color: #333;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.pagination {
  margin-top: 20px;
  justify-content: center;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .uploads-container {
    padding: 15px;
  }
  
  .header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header h2 {
    margin-left: 0;
    margin-top: 10px;
  }
  
  .el-table {
    font-size: 12px;
  }
  
  .el-table .el-table__cell {
    padding: 6px 0;
  }
}
</style>