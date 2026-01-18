<template>
  <div class="uploads-container">
    <div class="header">
      <el-button type="primary" link @click="goBack" class="back-button">
        <el-icon>
          <ArrowLeft />
        </el-icon>
      </el-button>
      <h2>上传历史</h2>
    </div>

    <div class="content" ref="contentRef">
      <el-table :data="uploads" style="width: 100%" v-loading="loading">
        <el-table-column prop="title" label="小说标题" />
        <el-table-column prop="author" label="作者" width="120" />
        <el-table-column prop="word_count" label="字数" width="100" />
        <el-table-column prop="click_count" label="点击量" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
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
            <el-button size="small" @click="viewNovel(row.id)" type="primary">
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 加载更多提示 -->
      <div v-if="hasMore" class="loading-more">
        <el-skeleton v-if="loadingMore" :rows="2" />
      </div>

      <!-- 没有更多数据提示 -->
      <div v-if="!hasMore && uploads.length > 0" class="no-more">
        <el-divider>没有更多数据了</el-divider>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
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
    const loadingMore = ref(false)
    const hasMore = ref(true)
    const currentPage = ref(1)
    const pageSize = ref(10)
    const contentRef = ref(null) // 用于滚动监听

    // 获取上传历史
    const fetchUploads = async (isLoadMore = false) => {
      if (isLoadMore) {
        loadingMore.value = true
      } else {
        loading.value = true
      }

      try {
        const response = await apiClient.get('/api/v1/novels', {
          params: {
            page: currentPage.value,
            limit: pageSize.value,
            upload_user_id: 'current' // 这里需要根据实际API调整
          }
        })

        if (response.data.code === 200) {
          const newUploads = response.data.data.novels
          const total = response.data.data.pagination.total || 0

          if (newUploads && newUploads.length > 0) {
            if (isLoadMore) {
              uploads.value = [...uploads.value, ...newUploads]
            } else {
              uploads.value = newUploads
            }

            // 检查是否还有更多数据
            hasMore.value = uploads.value.length < total
          } else {
            hasMore.value = false
          }
        } else {
          ElMessage.error('获取上传历史失败: ' + response.data.message)
          hasMore.value = false
        }
      } catch (error) {
        console.error('获取上传历史失败:', error)
        ElMessage.error('获取上传历史失败: ' + error.message)
        hasMore.value = false
      } finally {
        loading.value = false
        loadingMore.value = false
      }
    }

    // 加载更多数据
    const loadMore = async () => {
      if (loading.value || loadingMore.value || !hasMore.value) return

      currentPage.value += 1
      await fetchUploads(true)
    }

    // 滚动事件监听
    const handleScroll = () => {
      if (!contentRef.value || loading.value || loadingMore.value || !hasMore.value) return

      const element = contentRef.value
      // 检查是否滚动到底部
      if (element.scrollHeight - element.scrollTop <= element.clientHeight + 100) {
        loadMore()
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

    onMounted(async () => {
      // 初始加载数据
      await fetchUploads()

      // 添加滚动事件监听
      if (contentRef.value) {
        contentRef.value.addEventListener('scroll', handleScroll)
      }
    })

    // 组件卸载时移除事件监听
    onUnmounted(() => {
      if (contentRef.value) {
        contentRef.value.removeEventListener('scroll', handleScroll)
      }
    })

    return {
      uploads,
      loading,
      loadingMore,
      hasMore,
      contentRef,
      formatDate,
      getStatusType,
      getStatusText,
      viewNovel,
      goBack
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
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
  /* 防止内容溢出父容器 */
}

.header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
  flex-shrink: 0;
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
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.loading-more {
  text-align: center;
  padding: 20px 0;
}

.no-more {
  text-align: center;
  padding: 20px 0;
  color: #999;
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