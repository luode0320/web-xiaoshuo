<template>
  <div class="upload-container">
    <div class="header">
      <h2>{{ isUploadView ? '上传小说' : '上传历史' }}</h2>
      <el-switch v-model="isUploadView" class="toggle-switch" :active-value="true" :inactive-value="false" />
    </div>

    <!-- 上传视图 -->
    <div v-if="isUploadView" class="upload-section">
      <!-- 上传区域移到上方 -->
      <div class="upload-area">
        <el-upload ref="uploadRef" :auto-upload="false" :on-change="handleFileChange" :on-remove="handleFileRemove" :before-upload="beforeUpload" :file-list="fileList" :multiple="true"
          accept=".txt,.epub" drag>
          <div class="upload-content">
            <el-icon class="upload-icon">
              <Upload />
            </el-icon>
            <div class="upload-text">
              <div class="upload-drag-text">将文件拖到此处，或<em>点击上传</em></div>
              <div class="upload-hint">只能上传 txt 或 epub 文件，且不超过 20MB</div>
            </div>
          </div>
        </el-upload>
      </div>

      <!-- 显示已选中的文件列表 -->
      <div class="selected-files" v-if="fileList.length > 0">
        <h3>已选择文件 ({{ fileList.length }})</h3>
        <div class="file-list">
          <div v-for="(file, index) in fileList" :key="index" class="file-item">
            <div class="file-info">
              <el-icon>
                <Document />
              </el-icon>
              <span class="file-name">{{ file.name }}</span>
              <span class="file-size">({{ formatFileSize(file.size) }})</span>
            </div>
            <el-button type="danger" size="small" @click="removeFile(index)">
              取消
            </el-button>
          </div>
        </div>
      </div>

      <!-- 批量上传按钮 -->
      <div class="upload-actions">
        <el-button type="primary" @click="submitUpload" :loading="uploading" :disabled="fileList.length === 0" size="large">
          {{ uploading ? '上传中...' : `批量上传 (${fileList.length} 个文件)` }}
        </el-button>
      </div>
    </div>

    <!-- 上传历史视图 -->
    <div v-else class="history-section" ref="contentRef">
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
      <div v-if="hasMore && uploads.length > 0" class="loading-more">
        <el-skeleton v-if="loadingMore" :rows="2" />
      </div>

      <!-- 没有更多数据提示 -->
      <div v-if="!hasMore && uploads.length > 0" class="no-more">
        <el-divider>没有更多数据了</el-divider>
      </div>

      <!-- 没有数据提示 -->
      <div v-if="uploads.length === 0 && !loading" class="no-data">
        <p>暂无上传记录</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Document, ArrowLeft } from '@element-plus/icons-vue'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'

export default {
  name: 'NovelUpload',
  components: {
    Upload,
    Document,
    ArrowLeft
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const uploadRef = ref(null)

    // 上传相关
    const uploading = ref(false)
    const fileList = ref([]) // 存储所有选中的文件
    // 视图切换
    const isUploadView = ref(false) // 默认显示上传历史

    // 上传历史相关
    const uploads = ref([])
    const loading = ref(false)
    const loadingMore = ref(false)
    const hasMore = ref(true)
    const currentPage = ref(1)
    const pageSize = ref(10)
    const contentRef = ref(null) // 用于滚动监听

    // 切换视图
    const toggleView = () => {
      isUploadView.value = !isUploadView.value
      // 如果切换到历史视图且还没有加载数据，则加载数据
      if (!isUploadView.value && uploads.value.length === 0) {
        currentPage.value = 1
        fetchUploads()
      }
    }

    // 处理文件选择变化（添加文件）
    const handleFileChange = (file, uploadFileList) => {
      // 提取文件名（不包含扩展名）作为小说标题
      const fileName = file.name.replace(/\.(txt|epub)$/i, '')

      // 添加自定义属性到文件对象
      file.novelTitle = fileName

      // 确保 fileList.value 与 Element Plus Upload 组件的文件列表保持同步
      fileList.value = [...uploadFileList]
    }

    // 处理文件移除
    const handleFileRemove = (file, uploadFileList) => {
      fileList.value = [...uploadFileList]
    }

    // 从fileList中移除指定文件
    const removeFile = (index) => {
      fileList.value.splice(index, 1)
    }

    // 上传前验证
    const beforeUpload = (file) => {
      const isTxtOrEpub = file.type === 'text/plain' || file.type === 'application/epub+zip' || file.name.toLowerCase().endsWith('.txt') || file.name.toLowerCase().endsWith('.epub')
      const isLt20M = file.size / 1024 / 1024 < 20

      if (!isTxtOrEpub) {
        ElMessage.error('只能上传 txt 或 epub 格式的文件!')
        return false
      }
      if (!isLt20M) {
        ElMessage.error('文件大小不能超过 20MB!')
        return false
      }

      return true
    }

    // 格式化文件大小
    const formatFileSize = (size) => {
      if (size < 1024) {
        return size + ' B'
      } else if (size < 1024 * 1024) {
        return (size / 1024).toFixed(2) + ' KB'
      } else {
        return (size / (1024 * 1024)).toFixed(2) + ' MB'
      }
    }

    // 提交批量上传
    const submitUpload = async () => {
      if (fileList.value.length === 0) {
        ElMessage.error('请至少选择一个小说文件')
        return
      }

      try {
        uploading.value = true

        let successCount = 0
        let failCount = 0

        // 逐个上传文件
        for (const file of fileList.value) {
          try {
            const formData = new FormData()
            // 使用 file.raw 而不是 file，因为 file.raw 才是实际的文件对象
            formData.append('file', file.raw || file)
            // 使用文件名（去掉扩展名）作为标题
            const fileName = file.name ? file.name.replace(/\.(txt|epub)$/i, '') : (file.raw ? file.raw.name.replace(/\.(txt|epub)$/i, '') : '未知小说')
            formData.append('title', fileName)
            // 使用空的作者、主角等信息，后端会自动提取
            formData.append('author', '')
            formData.append('protagonist', '')
            formData.append('description', '')

            await apiClient.post('/api/v1/novels/upload', formData, {
              headers: {
                'Authorization': `Bearer ${userStore.token}`,
                'Content-Type': 'multipart/form-data'
              },
              timeout: 120000 // 设置120秒超时，因为文件上传可能需要较长时间
            })

            successCount++
          } catch (error) {
            console.error('上传失败:', error)
            const errorMessage = error.response?.data?.message || error.message || '上传失败'
            const fileName = file.name || (file.raw ? file.raw.name : '未知文件')
            ElMessage.error(`文件 "${fileName}" 上传失败: ${errorMessage}`)
            failCount++
          }
        }

        if (failCount === 0) {
          ElMessage.success(`成功上传 ${successCount} 本小说！等待管理员审核。`)
        } else if (successCount === 0) {
          ElMessage.error(`全部 ${failCount} 个文件上传失败。`)
        } else {
          ElMessage.warning(`成功上传 ${successCount} 本小说，${failCount} 个文件上传失败。`)
        }

        // 重置文件列表
        fileList.value = []
        if (uploadRef.value) {
          uploadRef.value.clearFiles()
        }

        // 切换到上传历史视图以显示新上传的小说
        isUploadView.value = false
        currentPage.value = 1
        await fetchUploads()
      } catch (error) {
        console.error('批量上传失败:', error)
        ElMessage.error('批量上传失败')
      } finally {
        uploading.value = false
      }
    }

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
            upload_user_id: userStore.user?.id || '' // 使用当前用户ID过滤上传历史
          }
        })

        if (response.data.code === 200) {
          const newUploads = response.data.data.novels || []
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

    // 检查用户是否已登录
    if (!userStore.isAuthenticated) {
      router.push('/login')
    }

    onMounted(async () => {
      // 初始加载上传历史数据（因为默认显示历史）
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
      // 上传相关
      uploadRef,
      uploading,
      fileList,
      handleFileChange,
      handleFileRemove,
      beforeUpload,
      removeFile,
      formatFileSize,
      submitUpload,
      // 视图切换
      isUploadView,
      toggleView,
      // 上传历史相关
      uploads,
      loading,
      loadingMore,
      hasMore,
      contentRef,
      formatDate,
      getStatusType,
      getStatusText,
      viewNovel
    }
  }
}
</script>

<style scoped>
.upload-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: calc(100vh - 60px);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
  flex-shrink: 0;
}

.header h2 {
  margin: 0;
  color: #333;
  flex: 1;
  /* 让标题占据可用空间 */
}

.toggle-switch {
  margin-left: auto;
}

.upload-section h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
}

.upload-area {
  margin-bottom: 30px;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.upload-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 16px;
}

.upload-text {
  text-align: center;
}

.upload-drag-text {
  color: #909399;
  font-size: 14px;
}

.upload-drag-text em {
  color: #409eff;
  font-style: normal;
}

.upload-hint {
  color: #c0c4cc;
  font-size: 12px;
  margin-top: 8px;
}

:deep(.el-upload-dragger) {
  width: 100%;
  height: 200px;
}

.selected-files {
  margin-bottom: 30px;
}

.selected-files h3 {
  margin-bottom: 15px;
  color: #333;
}

.file-list {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 10px;
}

.file-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.file-item:last-child {
  border-bottom: none;
}

.file-info {
  display: flex;
  align-items: center;
  flex: 1;
}

.file-info .el-icon {
  margin-right: 8px;
}

.file-name {
  flex: 1;
  word-break: break-all;
  margin-right: 10px;
}

.file-size {
  color: #909399;
  font-size: 12px;
}

.upload-actions {
  text-align: center;
  margin-top: 20px;
}

.history-section {
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

.no-data {
  text-align: center;
  padding: 40px 0;
  color: #999;
  font-size: 16px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .upload-container {
    padding: 15px;
  }

  .header {
    flex-wrap: nowrap;
    /* 确保标题和开关在同一行 */
  }

  .header h2 {
    margin: 0;
    flex: 1;
    font-size: 16px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    /* 如果标题太长则显示省略号 */
    max-width: calc(100% - 60px);
    /* 为开关预留空间 */
  }

  .toggle-switch {
    margin-left: 10px;
    flex-shrink: 0;
    /* 防止开关被压缩 */
  }

  .el-table {
    font-size: 12px;
  }

  .el-table .el-table__cell {
    padding: 6px 0;
  }

  :deep(.el-switch__core) {
    width: 40px !important;
    height: 20px !important;
  }

  :deep(.el-switch__core .el-switch__action) {
    width: 16px !important;
    height: 16px !important;
  }
}
</style>