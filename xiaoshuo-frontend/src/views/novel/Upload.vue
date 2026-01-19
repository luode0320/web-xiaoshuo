<template>
  <div class="upload-container">
    <div class="upload-form">
      <h2>上传小说</h2>
      
      <el-form 
        :model="uploadForm" 
        :rules="uploadRules" 
        ref="uploadFormRef"
        label-width="100px"
      >
        <el-form-item label="小说标题" prop="title">
          <el-input 
            v-model="uploadForm.title" 
            placeholder="请输入小说标题"
          />
        </el-form-item>
        
        <el-form-item label="作者" prop="author">
          <el-input 
            v-model="uploadForm.author" 
            placeholder="请输入作者名"
          />
        </el-form-item>
        
        <el-form-item label="主角" prop="protagonist">
          <el-input 
            v-model="uploadForm.protagonist" 
            placeholder="请输入主角名"
          />
        </el-form-item>
        
        <el-form-item label="分类">
          <el-select 
            v-model="uploadForm.categoryIds" 
            multiple 
            placeholder="请选择分类"
            style="width: 100%"
          >
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="关键词">
          <el-input 
            v-model="uploadForm.keywords" 
            placeholder="请输入关键词，用逗号分隔"
          />
        </el-form-item>
        
        <el-form-item label="简介" prop="description">
          <el-input 
            v-model="uploadForm.description" 
            type="textarea" 
            :rows="4"
            placeholder="请输入小说简介"
            maxlength="1000"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="小说文件" prop="file">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :on-change="handleFileChange"
            :before-upload="beforeUpload"
            :show-file-list="true"
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
        </el-form-item>
        
        <el-form-item>
          <el-button 
            type="primary" 
            @click="submitUpload" 
            :loading="uploading"
            :disabled="!uploadForm.file"
          >
            {{ uploading ? '上传中...' : '上传小说' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import apiClient from '@/utils/api'

export default {
  name: 'NovelUpload',
  components: {
    UploadIcon: Upload
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const uploadRef = ref(null)
    const uploadFormRef = ref(null)
    
    const uploading = ref(false)
    const categories = ref([])
    
    const uploadForm = reactive({
      title: '',
      author: '',
      protagonist: '',
      categoryIds: [],
      keywords: '',
      description: '',
      file: null
    })
    
    const uploadRules = {
      title: [
        { required: true, message: '请输入小说标题', trigger: 'blur' },
        { max: 200, message: '标题长度不能超过200个字符', trigger: 'blur' }
      ],
      author: [
        { required: true, message: '请输入作者名', trigger: 'blur' },
        { max: 100, message: '作者名长度不能超过100个字符', trigger: 'blur' }
      ],
      description: [
        { max: 1000, message: '简介长度不能超过1000个字符', trigger: 'blur' }
      ],
      file: [
        { required: true, message: '请选择小说文件', trigger: 'change' }
      ]
    }
    
    // 获取分类列表
    const fetchCategories = async () => {
      try {
        const response = await apiClient.get('/api/v1/categories')
        categories.value = response.data.data.categories
      } catch (error) {
        console.error('获取分类失败:', error)
        ElMessage.error('获取分类失败')
      }
    }
    
    // 处理文件选择变化
    const handleFileChange = (file, fileList) => {
      if (fileList.length > 1) {
        fileList.shift() // 只保留最后一个文件
      }
      uploadForm.file = file.raw
    }
    
    // 上传前验证
    const beforeUpload = (file) => {
      const isTxtOrEpub = file.type === 'text/plain' || file.type === 'application/epub+zip' || file.name.toLowerCase().endsWith('.txt') || file.name.toLowerCase().endsWith('.epub')
      const isLt20M = file.size / 1024 / 1024 < 20
      
      if (!isTxtOrEpub) {
        ElMessage.error('只能上传 txt 或 epub 格式的文件!')
      }
      if (!isLt20M) {
        ElMessage.error('文件大小不能超过 20MB!')
      }
      
      return isTxtOrEpub && isLt20M
    }
    
    // 提交上传
    const submitUpload = async () => {
      if (!uploadFormRef.value) return
      
      try {
        await uploadFormRef.value.validate()
        
        if (!uploadForm.file) {
          ElMessage.error('请选择小说文件')
          return
        }
        
        uploading.value = true
        
        const formData = new FormData()
        formData.append('file', uploadForm.file)
        formData.append('title', uploadForm.title)
        formData.append('author', uploadForm.author)
        formData.append('protagonist', uploadForm.protagonist)
        formData.append('description', uploadForm.description)
        
        // 添加分类ID（如果有选择）
        if (uploadForm.categoryIds.length > 0) {
          uploadForm.categoryIds.forEach(id => {
            formData.append('category_ids', id)
          })
        }
        
        // 添加关键词（如果有输入）
        if (uploadForm.keywords.trim()) {
          const keywords = uploadForm.keywords.split(',').map(k => k.trim()).filter(k => k)
          keywords.forEach(keyword => {
            formData.append('keywords', keyword)
          })
        }
        
        const response = await apiClient.post('/api/v1/novels/upload', formData, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`,
            'Content-Type': 'multipart/form-data'
          }
        })
        
        ElMessage.success('小说上传成功！等待管理员审核。')
        
        // 重置表单
        resetForm()
        
        // 跳转到用户个人页面查看上传历史
        router.push('/profile')
      } catch (error) {
        console.error('上传失败:', error)
        if (error.response?.data?.message) {
          ElMessage.error(error.response.data.message)
        } else {
          ElMessage.error('上传失败')
        }
      } finally {
        uploading.value = false
      }
    }
    
    // 重置表单
    const resetForm = () => {
      uploadForm.title = ''
      uploadForm.author = ''
      uploadForm.protagonist = ''
      uploadForm.categoryIds = []
      uploadForm.keywords = ''
      uploadForm.description = ''
      uploadForm.file = null
      
      if (uploadRef.value) {
        uploadRef.value.clearFiles()
      }
    }
    
    // 检查用户是否已登录
    if (!userStore.isAuthenticated) {
      router.push('/login')
    }
    
    // 初始化时获取分类
    fetchCategories()
    
    return {
      uploadRef,
      uploadFormRef,
      uploading,
      categories,
      uploadForm,
      uploadRules,
      handleFileChange,
      beforeUpload,
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

.upload-form h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
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
</style>