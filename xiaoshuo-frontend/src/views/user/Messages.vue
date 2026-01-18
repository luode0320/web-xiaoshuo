<template>
  <div class="messages-container">
    <div class="header">
      <el-button type="primary" link @click="goBack" class="back-button">
        <el-icon>
          <ArrowLeft />
        </el-icon>
      </el-button>
      <h2>系统消息</h2>
    </div>

    <div class="content">
      <el-card v-for="message in messages" :key="message.id" class="message-card" :class="{ unread: !message.is_read }">
        <div class="message-header">
          <h3>{{ message.title }}</h3>
          <el-tag :type="getMessageType(message.type)" size="small">
            {{ getMessageTypeText(message.type) }}
          </el-tag>
        </div>
        <div class="message-content">
          {{ message.content }}
        </div>
        <div class="message-footer">
          <span class="time">{{ formatDate(message.created_at) }}</span>
          <el-button size="small" @click="markAsRead(message.id)" :disabled="message.is_read" type="primary">
            {{ message.is_read ? '已读' : '标记为已读' }}
          </el-button>
        </div>
      </el-card>

      <div v-if="messages.length === 0" class="empty-state">
        <el-empty description="暂无系统消息" />
      </div>

      <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" :total="total" layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange" @current-change="handleCurrentChange" class="pagination" />
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
  name: 'Messages',
  components: {
    ArrowLeft
  },
  setup() {
    const router = useRouter()
    const messages = ref([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)

    // 获取系统消息
    const fetchMessages = async () => {
      try {
        // 这里需要根据实际API调整，假设有一个获取系统消息的接口
        // 目前使用一个模拟数据，实际需要后端API支持
        const response = await apiClient.get('/api/v1/admin/system-messages', {
          params: {
            page: currentPage.value,
            limit: pageSize.value
          }
        })

        if (response.data.code === 200) {
          messages.value = response.data.data.messages
          total.value = response.data.data.pagination.total
        } else {
          ElMessage.error('获取系统消息失败: ' + response.data.message)
        }
      } catch (error) {
        console.error('获取系统消息失败:', error)
        // 如果API不存在，使用模拟数据
        messages.value = [
          {
            id: 1,
            title: '系统维护通知',
            content: '系统将于今晚22:00-24:00进行维护，届时可能无法访问。',
            type: 'notification',
            is_read: false,
            created_at: new Date().toISOString()
          },
          {
            id: 2,
            title: '新功能上线',
            content: '我们新增了全文搜索功能，欢迎体验！',
            type: 'announcement',
            is_read: true,
            created_at: new Date(Date.now() - 86400000).toISOString() // 一天前
          }
        ]
        total.value = 2
      }
    }

    // 获取消息类型
    const getMessageType = (type) => {
      switch (type) {
        case 'notification': return 'info'
        case 'announcement': return 'success'
        case 'warning': return 'warning'
        default: return 'info'
      }
    }

    // 获取消息类型文本
    const getMessageTypeText = (type) => {
      switch (type) {
        case 'notification': return '通知'
        case 'announcement': return '公告'
        case 'warning': return '警告'
        default: return type
      }
    }

    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }

    // 标记为已读
    const markAsRead = async (id) => {
      try {
        // 这里需要根据实际API调整
        // 假设有一个标记消息为已读的API
        const response = await apiClient.put(`/api/v1/admin/system-messages/${id}/read`, {
          is_read: true
        })

        if (response.data.code === 200) {
          const message = messages.value.find(m => m.id === id)
          if (message) {
            message.is_read = true
          }
          ElMessage.success('已标记为已读')
        } else {
          ElMessage.error('标记失败: ' + response.data.message)
        }
      } catch (error) {
        console.error('标记消息为已读失败:', error)
        ElMessage.error('标记失败: ' + error.message)
      }
    }

    // 返回上一页
    const goBack = () => {
      router.push('/profile')
    }

    // 处理页面大小变化
    const handleSizeChange = (size) => {
      pageSize.value = size
      currentPage.value = 1
      fetchMessages()
    }

    // 处理当前页变化
    const handleCurrentChange = (page) => {
      currentPage.value = page
      fetchMessages()
    }

    onMounted(() => {
      fetchMessages()
    })

    return {
      messages,
      currentPage,
      pageSize,
      total,
      getMessageType,
      getMessageTypeText,
      formatDate,
      markAsRead,
      goBack,
      handleSizeChange,
      handleCurrentChange
    }
  }
}
</script>

<style scoped>
.messages-container {
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

.message-card {
  border-left: 4px solid #409eff;
  margin-bottom: 15px;
}

.message-card.unread {
  border-left: 4px solid #f56c6c;
  background-color: #fef0f0;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.message-header h3 {
  margin: 0;
  color: #333;
}

.message-content {
  color: #666;
  margin-bottom: 15px;
  line-height: 1.6;
}

.message-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid #f5f5f5;
  padding-top: 15px;
}

.time {
  color: #999;
  font-size: 12px;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
}

.pagination {
  margin-top: 20px;
  justify-content: center;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .messages-container {
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

  .message-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .message-footer {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }
}
</style>