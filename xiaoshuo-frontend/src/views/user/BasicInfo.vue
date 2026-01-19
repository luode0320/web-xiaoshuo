<template>
  <div class="basic-info-container">
    <div class="header">
      <el-button type="primary" link @click="goBack" class="back-button">
        <el-icon>
          <ArrowLeft />
        </el-icon>
      </el-button>
      <h2>基本信息</h2>
    </div>

    <div class="content">
      <el-card class="info-card">
        <!-- 添加头像部分 -->
        <div class="avatar-section">
          <div class="avatar-preview">
            <input type="file" ref="avatarInputRef" accept="image/*" @change="handleAvatarUpload" style="display: none;" />
            <div @click="selectAvatarFile" class="clickable">
              <img v-if="user?.avatar" :src="user.avatar" :alt="user.nickname" class="avatar-image clickable" />
              <div v-else class="avatar-placeholder clickable">{{ user?.nickname?.charAt(0) || 'U' }}</div>
            </div>
          </div>
        </div>

        <div class="info-item">
          <span class="label">昵称:</span>
          <span class="value">{{ user?.nickname || '未设置' }}</span>
        </div>
        <div class="info-item">
          <span class="label">邮箱:</span>
          <span class="value">{{ user?.email }}</span>
        </div>
        <div class="info-item">
          <span class="label">注册时间:</span>
          <span class="value">{{ formatDate(user?.created_at) }}</span>
        </div>
        <div class="info-item">
          <span class="label">账户状态:</span>
          <span class="value">{{ user?.is_active ? '正常' : '已冻结' }}</span>
        </div>
        <div class="info-item" v-if="userStore.isAdmin">
          <span class="label">管理员:</span>
          <span class="value">是</span>
        </div>
      </el-card>

      <div class="button-container">
        <el-button type="primary" @click="showEditDialog = true" class="edit-button">
          编辑昵称
        </el-button>
        <el-button type="primary" @click="showPasswordDialog = true" class="password-button">
          重置密码
        </el-button>
      </div>
    </div>

    <!-- 编辑昵称对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑昵称" width="400px">
      <el-form :model="editForm" :rules="nicknameRules" ref="editFormRef" label-width="80px">
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="editForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="updateNickname">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog v-model="showPasswordDialog" title="重置密码" width="400px">
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input v-model="passwordForm.currentPassword" type="password" placeholder="请输入当前密码" />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showPasswordDialog = false">取消</el-button>
          <el-button type="primary" @click="updatePassword">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'

export default {
  name: 'BasicInfo',
  components: {
    ArrowLeft
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const showEditDialog = ref(false)
    const showPasswordDialog = ref(false)
    const editForm = ref({
      nickname: ''
    })
    const passwordForm = ref({
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    const editFormRef = ref(null)
    const passwordFormRef = ref(null)
    const avatarInputRef = ref(null)


    const user = computed(() => userStore.user)

    // 昵称表单验证规则
    const nicknameRules = {
      nickname: [
        { required: true, message: '请输入昵称', trigger: 'blur' },
        { min: 1, max: 20, message: '昵称长度在1到20个字符之间', trigger: 'blur' }
      ]
    }

    // 密码表单验证规则
    const validatePassword = (rule, value, callback) => {
      if (value && passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }

    const passwordRules = {
      currentPassword: [
        { required: true, message: '请输入当前密码', trigger: 'blur' }
      ],
      newPassword: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度在6到20个字符之间', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认新密码', trigger: 'blur' },
        { validator: validatePassword, trigger: 'blur' }
      ]
    }

    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }

    // 返回上一页
    const goBack = () => {
      router.push('/profile')
    }

    // 选择头像文件
    const selectAvatarFile = () => {
      if (avatarInputRef.value) {
        avatarInputRef.value.click()
      }
    }

    // 处理头像上传
    const handleAvatarUpload = async (event) => {
      const file = event.target.files[0]
      if (!file) return

      // 验证文件
      const isImage = file.type.startsWith('image/')
      const isLt5M = file.size / 1024 / 1024 < 5

      if (!isImage) {
        ElMessage.error('上传头像图片只能是 JPG/PNG 格式!')
        return
      }
      if (!isLt5M) {
        ElMessage.error('上传头像大小不能超过 5MB!')
        return
      }

      // 将文件转换为base64
      const reader = new FileReader()
      reader.onload = async (e) => {
        try {
          const base64Data = e.target.result

          const response = await apiClient.put('/api/v1/users/profile', {
            avatar: base64Data
          })

          if (response.data.code === 200) {
            ElMessage.success('头像更新成功')
            // 更新用户信息
            await userStore.fetchProfile()
          } else {
            ElMessage.error('更新头像失败: ' + response.data.message)
          }
        } catch (error) {
          console.error('更新头像失败:', error)
          ElMessage.error('更新头像失败: ' + (error.response?.data?.message || error.message))
        }
      }
      reader.readAsDataURL(file)
    }

    // 更新昵称
    const updateNickname = async () => {
      try {
        // 验证表单
        if (!editForm.value.nickname.trim()) {
          ElMessage.error('请输入昵称')
          return
        }

        if (editForm.value.nickname.length < 1 || editForm.value.nickname.length > 20) {
          ElMessage.error('昵称长度必须在1到20个字符之间')
          return
        }

        // 检查昵称是否发生变化
        if (editForm.value.nickname === user.value?.nickname) {
          ElMessage.info('昵称未发生变化')
          showEditDialog.value = false
          return
        }

        const updateData = {
          nickname: editForm.value.nickname
        }

        const response = await apiClient.put('/api/v1/users/profile', updateData)

        if (response.data.code === 200) {
          ElMessage.success('昵称更新成功')
          await userStore.fetchProfile()
          showEditDialog.value = false
        } else {
          ElMessage.error('更新失败: ' + response.data.message)
        }
      } catch (error) {
        console.error('更新昵称失败:', error)
        ElMessage.error('更新失败: ' + (error.response?.data?.message || error.message))
      }
    }

    // 更新密码
    const updatePassword = async () => {
      try {
        // 验证表单
        if (!passwordForm.value.currentPassword) {
          ElMessage.error('请输入当前密码')
          return
        }
        if (!passwordForm.value.newPassword) {
          ElMessage.error('请输入新密码')
          return
        }
        if (!passwordForm.value.confirmPassword) {
          ElMessage.error('请确认新密码')
          return
        }
        if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
          ElMessage.error('两次输入的密码不一致')
          return
        }
        if (passwordForm.value.newPassword.length < 6) {
          ElMessage.error('密码长度至少为6位')
          return
        }

        const response = await apiClient.put('/api/v1/users/profile', {
          current_password: passwordForm.value.currentPassword,
          new_password: passwordForm.value.newPassword
        })

        if (response.data.code === 200) {
          ElMessage.success('密码更新成功，请使用新密码重新登录')
          showPasswordDialog.value = false
          passwordForm.value = {
            currentPassword: '',
            newPassword: '',
            confirmPassword: ''
          }
          // 由于密码已更改，用户需要重新登录
          userStore.logout()
          router.push('/login')
        } else {
          ElMessage.error('密码更新失败: ' + response.data.message)
        }
      } catch (error) {
        console.error('更新密码失败:', error)
        ElMessage.error('密码更新失败: ' + (error.response?.data?.message || error.message))
      }
    }

    onMounted(() => {
      if (user.value) {
        editForm.value.nickname = user.value.nickname
      }
    })

    return {
      user,
      userStore,
      showEditDialog,
      showPasswordDialog,
      editForm,
      passwordForm,
      editFormRef,
      passwordFormRef,
      avatarInputRef,

      nicknameRules,
      passwordRules,
      formatDate,
      goBack,
      selectAvatarFile,
      handleAvatarUpload,

      updateNickname,
      updatePassword
    }
  }
}
</script>

<style scoped>
.basic-info-container {
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

.info-card {
  width: 100%;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f5f5f5;
}

.avatar-preview {
  display: flex;
  align-items: center;
}

.avatar-image {
  width: 75px;
  height: 75px;
  /* border-radius: 50%;border: 2px solid #e4e7ed; */
  object-fit: contain;
  cursor: pointer;
  transition: opacity 0.3s;
}

.avatar-image:hover {
  opacity: 0.8;
}

.clickable {
  cursor: pointer;
  transition: opacity 0.3s;
}

.clickable:hover {
  opacity: 0.8;
}

.avatar-placeholder {
  width: 75px;
  height: 75px;
  border-radius: 50%;
  background: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
  cursor: pointer;
}

.avatar-upload {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.upload-btn {
  margin-top: 5px;
}

.info-item {
  display: flex;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f5f5f5;
}

.info-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.label {
  font-weight: bold;
  width: 100px;
  color: #666;
}

.value {
  flex: 1;
  color: #333;
}

.edit-button {
  align-self: flex-start;
}

.button-container {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .basic-info-container {
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

  .avatar-section {
    flex-direction: column;
    align-items: center;
    gap: 15px;
  }

  .label {
    width: 80px;
    font-size: 14px;
  }
}
</style>