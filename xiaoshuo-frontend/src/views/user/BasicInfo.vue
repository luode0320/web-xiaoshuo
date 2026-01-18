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
          <span class="label">激活状态:</span>
          <span class="value">{{ user?.is_activated ? '已激活' : '未激活' }}</span>
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

      <el-button type="primary" @click="showEditDialog = true" class="edit-button">
        编辑信息
      </el-button>
    </div>

    <!-- 编辑信息对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑信息" width="400px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="昵称">
          <el-input v-model="editForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="updateUserInfo">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
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
    const editForm = ref({
      nickname: ''
    })

    const user = computed(() => userStore.user)

    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }

    // 返回上一页
    const goBack = () => {
      router.push('/profile')
    }

    // 编辑用户信息
    const updateUserInfo = async () => {
      try {
        const response = await apiClient.put('/api/v1/users/profile', {
          nickname: editForm.value.nickname
        })

        if (response.data.code === 200) {
          ElMessage.success('信息更新成功')
          await userStore.fetchUserInfo()
          showEditDialog.value = false
        } else {
          ElMessage.error('更新失败: ' + response.data.message)
        }
      } catch (error) {
        console.error('更新用户信息失败:', error)
        ElMessage.error('更新失败: ' + error.message)
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
      editForm,
      formatDate,
      goBack,
      updateUserInfo
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
  min-height: 0; /* 防止内容溢出父容器 */
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

  .label {
    width: 80px;
    font-size: 14px;
  }
}
</style>