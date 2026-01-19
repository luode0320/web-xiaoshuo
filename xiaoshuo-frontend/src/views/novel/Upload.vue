<template>
  <div class="upload-container">
    <div class="upload-section">
      <h2>上传小说</h2>
      
      <!-- 上传区域移到上方 -->
      <div class="upload-area">
        <el-upload
          ref="uploadRef"
          :auto-upload="false"
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
          :before-upload="beforeUpload"
          :file-list="fileList"
          :multiple="true"
          accept=".txt,.epub"
          drag
        >
          <div class="upload-content">
            <el-icon class="upload-icon"><Upload /></el-icon>
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
          <div 
            v-for="(file, index) in fileList" 
            :key="index" 
            class="file-item"
          >
            <div class="file-info">
              <el-icon><Document /></el-icon>
              <span class="file-name">{{ file.name }}</span>
              <span class="file-size">({{ formatFileSize(file.size) }})</span>
            </div>
            <el-button 
              type="danger" 
              size="small" 
              @click="removeFile(index)"
            >
              取消
            </el-button>
          </div>
        </div>
      </div>
      
      <!-- 批量上传按钮 -->
      <div class="upload-actions">
        <el-button 
          type="primary" 
          @click="submitUpload" 
          :loading="uploading"
          :disabled="fileList.length === 0"
          size="large"
        >
          {{ uploading ? '上传中...' : `批量上传 (${fileList.length} 个文件)` }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Document } from '@element-plus/icons-vue'
import apiClient from '@/utils/api'

export default {
  name: 'NovelUpload',
  components: {
    Upload,
    Document
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const uploadRef = ref(null)
    
    const uploading = ref(false)
    const fileList = ref([]) // 存储所有选中的文件
    
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
        
        // 跳转到用户个人页面查看上传历史
        router.push('/profile')
      } catch (error) {
        console.error('批量上传失败:', error)
        ElMessage.error('批量上传失败')
      } finally {
        uploading.value = false
      }
    }
    
    // 检查用户是否已登录
    if (!userStore.isAuthenticated) {
      router.push('/login')
    }
    
    return {
      uploadRef,
      uploading,
      fileList,
      handleFileChange,
      handleFileRemove,
      beforeUpload,
      removeFile,
      formatFileSize,
      submitUpload
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
</style>