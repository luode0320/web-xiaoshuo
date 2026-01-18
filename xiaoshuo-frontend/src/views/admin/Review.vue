<template>
  <div class="admin-review-container">
    <div class="admin-header">
      <h1>小说审核管理</h1>
      <p>审核上传的小说，管理内容质量</p>
    </div>
    
    <div class="admin-content">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="待审核小说" name="pending">
          <div class="pending-novels">
            <el-table 
              :data="pendingNovels" 
              style="width: 100%"
              v-loading="loading.pending"
            >
              <el-table-column prop="title" label="标题" width="200" />
              <el-table-column prop="author" label="作者" width="120" />
              <el-table-column prop="upload_user.nickname" label="上传者" width="120" />
              <el-table-column prop="description" label="简介" show-overflow-tooltip />
              <el-table-column prop="word_count" label="字数" width="80" />
              <el-table-column prop="created_at" label="上传时间" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200">
                <template #default="{ row }">
                  <el-button size="small" @click="viewNovel(row.id)">查看</el-button>
                  <el-button size="small" type="success" @click="approveNovel(row.id)">通过</el-button>
                  <el-button size="small" type="danger" @click="rejectNovel(row.id)">拒绝</el-button>
                </template>
              </el-table-column>
            </el-table>
            
            <!-- 分页 -->
            <div class="pagination" v-if="pendingPagination.total > pendingPagination.limit">
              <el-pagination
                v-model:current-page="pendingPagination.page"
                :page-size="pendingPagination.limit"
                :total="pendingPagination.total"
                @current-change="handlePendingPageChange"
                layout="prev, pager, next, jumper"
              />
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="已审核小说" name="approved">
          <div class="approved-novels">
            <el-table 
              :data="approvedNovels" 
              style="width: 100%"
              v-loading="loading.approved"
            >
              <el-table-column prop="title" label="标题" width="200" />
              <el-table-column prop="author" label="作者" width="120" />
              <el-table-column prop="upload_user.nickname" label="上传者" width="120" />
              <el-table-column prop="click_count" label="点击量" width="100" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="getStatusType(row.status)">
                    {{ getStatusText(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="上传时间" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150">
                <template #default="{ row }">
                  <el-button size="small" @click="viewNovel(row.id)">查看</el-button>
                  <el-button size="small" type="danger" @click="deleteNovel(row.id)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
            
            <!-- 分页 -->
            <div class="pagination" v-if="approvedPagination.total > approvedPagination.limit">
              <el-pagination
                v-model:current-page="approvedPagination.page"
                :page-size="approvedPagination.limit"
                :total="approvedPagination.total"
                @current-change="handleApprovedPageChange"
                layout="prev, pager, next, jumper"
              />
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="审核日志" name="logs">
          <div class="admin-logs">
            <el-table 
              :data="adminLogs" 
              style="width: 100%"
              v-loading="loading.logs"
            >
              <el-table-column prop="admin_user.nickname" label="操作员" width="120" />
              <el-table-column prop="action" label="操作" width="120" />
              <el-table-column prop="target_type" label="目标类型" width="100" />
              <el-table-column prop="details" label="详情" show-overflow-tooltip />
              <el-table-column prop="created_at" label="操作时间" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
            </el-table>
            
            <!-- 分页 -->
            <div class="pagination" v-if="logsPagination.total > logsPagination.limit">
              <el-pagination
                v-model:current-page="logsPagination.page"
                :page-size="logsPagination.limit"
                :total="logsPagination.total"
                @current-change="handleLogsPageChange"
                layout="prev, pager, next, jumper"
              />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'

export default {
  name: 'AdminReview',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const activeTab = ref('pending')
    const pendingNovels = ref([])
    const approvedNovels = ref([])
    const adminLogs = ref([])
    
    const loading = ref({
      pending: false,
      approved: false,
      logs: false
    })
    
    const pendingPagination = ref({
      page: 1,
      limit: 10,
      total: 0
    })
    
    const approvedPagination = ref({
      page: 1,
      limit: 10,
      total: 0
    })
    
    const logsPagination = ref({
      page: 1,
      limit: 10,
      total: 0
    })
    
    // 检查用户权限
    if (!userStore.isAuthenticated || !userStore.isAdmin) {
      router.push('/')
      ElMessage.error('您没有权限访问此页面')
      return
    }
    
    // 获取待审核小说
    const fetchPendingNovels = async () => {
      loading.value.pending = true
      try {
        const response = await apiClient.get(`/api/v1/novels/pending?page=${pendingPagination.value.page}&limit=${pendingPagination.value.limit}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        pendingNovels.value = response.data.data.novels
        pendingPagination.value.total = response.data.data.pagination.total
      } catch (error) {
        console.error('获取待审核小说失败:', error)
        ElMessage.error('获取待审核小说失败')
      } finally {
        loading.value.pending = false
      }
    }
    
    // 获取已审核小说
    const fetchApprovedNovels = async () => {
      loading.value.approved = true
      try {
        const response = await apiClient.get(`/api/v1/novels?page=${approvedPagination.value.page}&limit=${approvedPagination.value.limit}&status=approved`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        approvedNovels.value = response.data.data.novels
        approvedPagination.value.total = response.data.data.pagination.total
      } catch (error) {
        console.error('获取已审核小说失败:', error)
        ElMessage.error('获取已审核小说失败')
      } finally {
        loading.value.approved = false
      }
    }
    
    // 获取管理员日志
    const fetchAdminLogs = async () => {
      loading.value.logs = true
      try {
        const response = await apiClient.get(`/api/v1/admin/logs?page=${logsPagination.value.page}&limit=${logsPagination.value.limit}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        adminLogs.value = response.data.data.logs
        logsPagination.value.total = response.data.data.pagination.total
      } catch (error) {
        console.error('获取管理员日志失败:', error)
        ElMessage.error('获取管理员日志失败')
      } finally {
        loading.value.logs = false
      }
    }
    
    // 审核通过小说
    const approveNovel = async (novelId) => {
      try {
        await ElMessageBox.confirm('确定要通过这个小说的审核吗？', '审核确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'success'
        })
        
        await apiClient.post(`/api/v1/novels/${novelId}/approve`, {}, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('审核通过成功')
        fetchPendingNovels() // 刷新待审核列表
        fetchApprovedNovels() // 刷新已审核列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('审核通过失败:', error)
          ElMessage.error('审核通过失败')
        }
      }
    }
    
    // 拒绝小说
    const rejectNovel = async (novelId) => {
      try {
        await ElMessageBox.confirm('确定要拒绝这个小说吗？', '审核确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        // 实际应用中可能需要实现拒绝功能
        // 这里暂时使用删除功能模拟拒绝
        await apiClient.delete(`/api/v1/novels/${novelId}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('拒绝成功')
        fetchPendingNovels() // 刷新待审核列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('拒绝失败:', error)
          ElMessage.error('拒绝失败')
        }
      }
    }
    
    // 删除小说
    const deleteNovel = async (novelId) => {
      try {
        await ElMessageBox.confirm('确定要删除这个小说吗？', '删除确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'danger'
        })
        
        await apiClient.delete(`/api/v1/novels/${novelId}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('删除成功')
        fetchApprovedNovels() // 刷新已审核列表
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
          ElMessage.error('删除失败')
        }
      }
    }
    
    // 查看小说
    const viewNovel = (novelId) => {
      router.push(`/novel/${novelId}`)
    }
    
    // 处理分页变化
    const handlePendingPageChange = (page) => {
      pendingPagination.value.page = page
      fetchPendingNovels()
    }
    
    const handleApprovedPageChange = (page) => {
      approvedPagination.value.page = page
      fetchApprovedNovels()
    }
    
    const handleLogsPageChange = (page) => {
      logsPagination.value.page = page
      fetchAdminLogs()
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
    
    onMounted(() => {
      fetchPendingNovels()
      fetchApprovedNovels()
      fetchAdminLogs()
    })
    
    return {
      activeTab,
      pendingNovels,
      approvedNovels,
      adminLogs,
      loading,
      pendingPagination,
      approvedPagination,
      logsPagination,
      approveNovel,
      rejectNovel,
      deleteNovel,
      viewNovel,
      handlePendingPageChange,
      handleApprovedPageChange,
      handleLogsPageChange,
      formatDate,
      getStatusType,
      getStatusText
    }
  }
}
</script>

<style scoped>
.admin-review-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.admin-header {
  text-align: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.admin-header h1 {
  font-size: 1.8rem;
  color: #333;
  margin-bottom: 10px;
}

.admin-header p {
  color: #666;
}

.admin-content {
  min-height: 500px;
}

.pagination {
  text-align: center;
  margin-top: 20px;
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

@media (max-width: 768px) {
  .admin-review-container {
    padding: 15px;
  }
  
  .admin-header h1 {
    font-size: 1.5rem;
  }
  
  :deep(.el-table) {
    font-size: 14px;
  }
  
  :deep(.el-table .el-table__cell) {
    padding: 8px 0;
  }
}
</style>