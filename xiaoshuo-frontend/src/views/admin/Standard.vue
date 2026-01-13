<template>
  <div class="admin-standard-container">
    <div class="admin-header">
      <h1>审核标准配置</h1>
      <p>管理小说审核标准和配置</p>
    </div>
    
    <div class="admin-content">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="审核标准" name="standards">
          <div class="standards-config">
            <div class="config-actions">
              <el-button type="primary" @click="addStandard">添加标准</el-button>
            </div>
            
            <el-table 
              :data="standards" 
              style="width: 100%"
              v-loading="loading.standards"
            >
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="title" label="标题" width="200" />
              <el-table-column prop="type" label="类型" width="120">
                <template #default="{ row }">
                  <el-tag :type="getTypeType(row.type)">
                    {{ getTypeText(row.type) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="description" label="描述" show-overflow-tooltip />
              <el-table-column prop="created_at" label="创建时间" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200">
                <template #default="{ row }">
                  <el-button size="small" @click="editStandard(row)">编辑</el-button>
                  <el-button size="small" type="danger" @click="deleteStandard(row.id)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="系统消息" name="messages">
          <div class="messages-config">
            <div class="config-actions">
              <el-button type="primary" @click="addMessage">添加消息</el-button>
            </div>
            
            <el-table 
              :data="messages" 
              style="width: 100%"
              v-loading="loading.messages"
            >
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="title" label="标题" width="200" />
              <el-table-column prop="type" label="类型" width="100">
                <template #default="{ row }">
                  <el-tag :type="getMessageType(row.type)">
                    {{ getMessageText(row.type) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="content" label="内容" show-overflow-tooltip />
              <el-table-column prop="is_read" label="状态" width="80">
                <template #default="{ row }">
                  <el-tag :type="row.is_read ? 'success' : 'info'">
                    {{ row.is_read ? '已读' : '未读' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="创建时间" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200">
                <template #default="{ row }">
                  <el-button size="small" @click="editMessage(row)">编辑</el-button>
                  <el-button size="small" type="success" @click="publishMessage(row.id)">发布</el-button>
                  <el-button size="small" type="danger" @click="deleteMessage(row.id)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
    
    <!-- 审核标准编辑弹窗 -->
    <el-dialog 
      :title="editingStandard.id ? '编辑审核标准' : '添加审核标准'" 
      v-model="standardDialogVisible"
      width="600px"
    >
      <el-form :model="editingStandard" :rules="standardRules" ref="standardForm" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="editingStandard.title" placeholder="请输入标准标题" />
        </el-form-item>
        
        <el-form-item label="类型" prop="type">
          <el-select v-model="editingStandard.type" placeholder="请选择标准类型">
            <el-option label="内容标准" value="content"></el-option>
            <el-option label="格式标准" value="format"></el-option>
            <el-option label="质量标准" value="quality"></el-option>
            <el-option label="合规标准" value="compliance"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="editingStandard.description" 
            type="textarea" 
            :rows="4"
            placeholder="请输入标准描述" 
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="standardDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveStandard">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 系统消息编辑弹窗 -->
    <el-dialog 
      :title="editingMessage.id ? '编辑系统消息' : '添加系统消息'" 
      v-model="messageDialogVisible"
      width="700px"
    >
      <el-form :model="editingMessage" :rules="messageRules" ref="messageForm" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="editingMessage.title" placeholder="请输入消息标题" />
        </el-form-item>
        
        <el-form-item label="类型" prop="type">
          <el-select v-model="editingMessage.type" placeholder="请选择消息类型">
            <el-option label="通知" value="notification"></el-option>
            <el-option label="警告" value="warning"></el-option>
            <el-option label="信息" value="info"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="内容" prop="content">
          <el-input 
            v-model="editingMessage.content" 
            type="textarea" 
            :rows="6"
            placeholder="请输入消息内容" 
          />
        </el-form-item>
        
        <el-form-item label="接收用户">
          <el-select v-model="editingMessage.target_user_id" placeholder="选择特定用户（可选）" clearable>
            <el-option 
              v-for="user in users" 
              :key="user.id" 
              :label="user.nickname || user.email" 
              :value="user.id"
            />
          </el-select>
          <div class="tip">留空则向所有用户发送</div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="messageDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveMessage">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import dayjs from 'dayjs'

export default {
  name: 'AdminStandard',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const activeTab = ref('standards')
    const standards = ref([])
    const messages = ref([])
    const users = ref([])
    
    const loading = ref({
      standards: false,
      messages: false,
      users: false
    })
    
    const standardDialogVisible = ref(false)
    const messageDialogVisible = ref(false)
    
    const editingStandard = ref({
      id: null,
      title: '',
      type: '',
      description: ''
    })
    
    const editingMessage = ref({
      id: null,
      title: '',
      type: '',
      content: '',
      target_user_id: null
    })
    
    // 表单验证规则
    const standardRules = {
      title: [
        { required: true, message: '请输入标准标题', trigger: 'blur' },
        { min: 2, max: 50, message: '标题长度在2-50个字符之间', trigger: 'blur' }
      ],
      type: [
        { required: true, message: '请选择标准类型', trigger: 'change' }
      ],
      description: [
        { required: true, message: '请输入标准描述', trigger: 'blur' },
        { max: 500, message: '描述长度不能超过500个字符', trigger: 'blur' }
      ]
    }
    
    const messageRules = {
      title: [
        { required: true, message: '请输入消息标题', trigger: 'blur' },
        { min: 2, max: 100, message: '标题长度在2-100个字符之间', trigger: 'blur' }
      ],
      type: [
        { required: true, message: '请选择消息类型', trigger: 'change' }
      ],
      content: [
        { required: true, message: '请输入消息内容', trigger: 'blur' },
        { max: 1000, message: '内容长度不能超过1000个字符', trigger: 'blur' }
      ]
    }
    
    // 检查用户权限
    if (!userStore.isAuthenticated || !userStore.isAdmin) {
      router.push('/')
      ElMessage.error('您没有权限访问此页面')
      return
    }
    
    // 获取审核标准
    const fetchStandards = async () => {
      loading.value.standards = true
      try {
        const response = await axios.get('/api/v1/admin/review-criteria', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        standards.value = response.data.data
      } catch (error) {
        console.error('获取审核标准失败:', error)
        ElMessage.error('获取审核标准失败')
      } finally {
        loading.value.standards = false
      }
    }
    
    // 获取系统消息
    const fetchMessages = async () => {
      loading.value.messages = true
      try {
        const response = await axios.get('/api/v1/admin/system-messages', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        messages.value = response.data.data
      } catch (error) {
        console.error('获取系统消息失败:', error)
        ElMessage.error('获取系统消息失败')
      } finally {
        loading.value.messages = false
      }
    }
    
    // 获取用户列表
    const fetchUsers = async () => {
      loading.value.users = true
      try {
        const response = await axios.get('/api/v1/admin/users', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        users.value = response.data.data
      } catch (error) {
        console.error('获取用户列表失败:', error)
        ElMessage.error('获取用户列表失败')
      } finally {
        loading.value.users = false
      }
    }
    
    // 添加标准
    const addStandard = () => {
      editingStandard.value = {
        id: null,
        title: '',
        type: '',
        description: ''
      }
      standardDialogVisible.value = true
    }
    
    // 编辑标准
    const editStandard = (standard) => {
      editingStandard.value = { ...standard }
      standardDialogVisible.value = true
    }
    
    // 保存标准
    const saveStandard = async () => {
      try {
        if (editingStandard.value.id) {
          // 更新标准
          await axios.put(`/api/v1/admin/review-criteria/${editingStandard.value.id}`, editingStandard.value, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
          ElMessage.success('标准更新成功')
        } else {
          // 创建标准
          await axios.post('/api/v1/admin/review-criteria', editingStandard.value, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
          ElMessage.success('标准创建成功')
        }
        
        standardDialogVisible.value = false
        fetchStandards()
      } catch (error) {
        console.error('保存标准失败:', error)
        ElMessage.error('保存标准失败')
      }
    }
    
    // 删除标准
    const deleteStandard = async (id) => {
      try {
        await ElMessageBox.confirm('确定要删除这个审核标准吗？', '删除确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'danger'
        })
        
        await axios.delete(`/api/v1/admin/review-criteria/${id}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('标准删除成功')
        fetchStandards()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除标准失败:', error)
          ElMessage.error('删除标准失败')
        }
      }
    }
    
    // 添加消息
    const addMessage = () => {
      editingMessage.value = {
        id: null,
        title: '',
        type: '',
        content: '',
        target_user_id: null
      }
      messageDialogVisible.value = true
    }
    
    // 编辑消息
    const editMessage = (message) => {
      editingMessage.value = { ...message }
      messageDialogVisible.value = true
    }
    
    // 保存消息
    const saveMessage = async () => {
      try {
        if (editingMessage.value.id) {
          // 更新消息
          await axios.put(`/api/v1/admin/system-messages/${editingMessage.value.id}`, editingMessage.value, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
          ElMessage.success('消息更新成功')
        } else {
          // 创建消息
          await axios.post('/api/v1/admin/system-messages', editingMessage.value, {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          })
          ElMessage.success('消息创建成功')
        }
        
        messageDialogVisible.value = false
        fetchMessages()
      } catch (error) {
        console.error('保存消息失败:', error)
        ElMessage.error('保存消息失败')
      }
    }
    
    // 发布消息
    const publishMessage = async (id) => {
      try {
        await ElMessageBox.confirm('确定要发布这个系统消息吗？', '发布确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'success'
        })
        
        await axios.post(`/api/v1/admin/system-messages/${id}/publish`, {}, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('消息发布成功')
        fetchMessages()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('发布消息失败:', error)
          ElMessage.error('发布消息失败')
        }
      }
    }
    
    // 删除消息
    const deleteMessage = async (id) => {
      try {
        await ElMessageBox.confirm('确定要删除这个系统消息吗？', '删除确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'danger'
        })
        
        await axios.delete(`/api/v1/admin/system-messages/${id}`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        ElMessage.success('消息删除成功')
        fetchMessages()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除消息失败:', error)
          ElMessage.error('删除消息失败')
        }
      }
    }
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
    
    // 获取类型类型
    const getTypeType = (type) => {
      switch (type) {
        case 'content': return 'primary'
        case 'format': return 'success'
        case 'quality': return 'warning'
        case 'compliance': return 'danger'
        default: return 'info'
      }
    }
    
    // 获取类型文本
    const getTypeText = (type) => {
      switch (type) {
        case 'content': return '内容标准'
        case 'format': return '格式标准'
        case 'quality': return '质量标准'
        case 'compliance': return '合规标准'
        default: return type
      }
    }
    
    // 获取消息类型
    const getMessageType = (type) => {
      switch (type) {
        case 'notification': return 'primary'
        case 'warning': return 'warning'
        case 'info': return 'info'
        default: return 'info'
      }
    }
    
    // 获取消息文本
    const getMessageText = (type) => {
      switch (type) {
        case 'notification': return '通知'
        case 'warning': return '警告'
        case 'info': return '信息'
        default: return type
      }
    }
    
    onMounted(() => {
      fetchStandards()
      fetchMessages()
      fetchUsers()
    })
    
    return {
      activeTab,
      standards,
      messages,
      users,
      loading,
      standardDialogVisible,
      messageDialogVisible,
      editingStandard,
      editingMessage,
      standardRules,
      messageRules,
      addStandard,
      editStandard,
      saveStandard,
      deleteStandard,
      addMessage,
      editMessage,
      saveMessage,
      publishMessage,
      deleteMessage,
      formatDate,
      getTypeType,
      getTypeText,
      getMessageType,
      getMessageText
    }
  }
}
</script>

<style scoped>
.admin-standard-container {
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

.config-actions {
  margin-bottom: 20px;
}

.tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

@media (max-width: 768px) {
  .admin-standard-container {
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