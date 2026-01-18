<template>
  <div class="basic-info-container">
    <div class="page-header">
      <el-button 
        type="text" 
        @click="goBack"
        class="back-button"
      >
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h2>基本信息</h2>
    </div>
    
    <div class="content">
      <el-form 
        :model="profileForm" 
        :rules="profileRules" 
        ref="profileFormRef" 
        label-width="100px"
        class="info-form"
      >
        <el-form-item label="头像">
          <div class="avatar-section">
            <div class="avatar-placeholder">{{ user?.nickname?.charAt(0) || 'U' }}</div>
            <p class="avatar-desc">当前头像基于您的昵称首字母生成</p>
          </div>
        </el-form-item>
        
        <el-form-item label="昵称" prop="nickname">
          <el-input 
            v-model="profileForm.nickname" 
            placeholder="请输入昵称"
            :disabled="!editing"
            class="editable-field"
          />
        </el-form-item>
        
        <el-form-item label="邮箱">
          <el-input 
            v-model="user.email" 
            disabled
            class="disabled-field"
          />
        </el-form-item>
        
        <el-form-item label="注册时间">
          <el-input 
            :value="formatDate(user?.created_at)" 
            disabled
            class="disabled-field"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button 
            v-if="!editing" 
            type="primary" 
            @click="editing = true"
            class="action-button edit-button"
          >
            <el-icon><Edit /></el-icon>
            编辑资料
          </el-button>
          <template v-else>
            <el-button 
              @click="cancelEdit"
              class="action-button cancel-button"
            >
              <el-icon><Close /></el-icon>
              取消
            </el-button>
            <el-button 
              type="primary" 
              @click="saveProfile"
              class="action-button save-button"
            >
              <el-icon><Check /></el-icon>
              保存
            </el-button>
          </template>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'
import { 
  ArrowLeft,
  Edit,
  Check,
  Close
} from '@element-plus/icons-vue'

export default {
  name: 'BasicInfo',
  components: {
    ArrowLeft,
    Edit,
    Check,
    Close
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const editing = ref(false)
    const profileFormRef = ref(null)
    
    const user = computed(() => userStore.user)
    
    const profileForm = reactive({
      nickname: userStore.user?.nickname || ''
    })
    
    const profileRules = {
      nickname: [
        { required: true, message: '请输入昵称', trigger: 'blur' },
        { max: 20, message: '昵称长度不能超过20个字符', trigger: 'blur' }
      ]
    }
    
    const goBack = () => {
      router.go(-1) // 返回上一页
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
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
    
    onMounted(() => {
      if (!userStore.isAuthenticated) {
        router.push('/login')
        return
      }
      
      // 初始化表单
      profileForm.nickname = userStore.user?.nickname || ''
    })
    
    return {
      editing,
      profileForm,
      profileFormRef,
      profileRules,
      user,
      goBack,
      saveProfile,
      cancelEdit,
      formatDate
    }
  }
}
</script>

<style scoped>
.basic-info-container {
  max-width: 800px;
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

.info-form {
  max-width: 600px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 15px;
}

.avatar-placeholder {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
}

.avatar-desc {
  color: #999;
  font-size: 14px;
  margin: 0;
}

.editable-field :deep(.el-input__wrapper),
.disabled-field :deep(.el-input__wrapper) {
  background-color: #fafafa;
}

.action-button {
  margin-right: 10px;
}

.edit-button {
  background: #409eff;
  color: white;
  border-color: #409eff;
}

.cancel-button {
  border-color: #dcdfe6;
}

.save-button {
  background: #67c23a;
  color: white;
  border-color: #67c23a;
}

@media (max-width: 768px) {
  .basic-info-container {
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
  
  .avatar-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>