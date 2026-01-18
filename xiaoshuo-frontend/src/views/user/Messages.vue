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

    <div class="content" ref="contentRef">
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

      <!-- 加载更多提示 -->
      <div v-if="hasMore" class="loading-more">
        <el-skeleton v-if="loadingMore" :rows="2" />
      </div>

      <!-- 没有更多数据提示 -->
      <div v-if="!hasMore && messages.length > 0" class="no-more">
        <el-divider>没有更多数据了</el-divider>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, watch } from 'vue'
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
    const loading = ref(false)
    const loadingMore = ref(false)
    const hasMore = ref(true)
    const currentPage = ref(1)
    const pageSize = ref(10)
    const contentRef = ref(null) // 用于滚动监听

    // 获取系统消息
    const fetchMessages = async (isLoadMore = false) => {
      if (isLoadMore) {
        loadingMore.value = true
      } else {
        loading.value = true
      }

      try {
        // 尝试从API获取系统消息，如果失败则使用模拟数据
        const response = await apiClient.get('/api/v1/users/search-history', { // 使用一个可能存在的API作为示例
          params: {
            page: currentPage.value,
            limit: pageSize.value
          }
        })

        if (response.data.code === 200) {
          const newMessages = response.data.data.search_history || response.data.data.messages || []
          const total = response.data.data.pagination?.total || response.data.data.total || 0

          if (newMessages && newMessages.length > 0) {
            if (isLoadMore) {
              messages.value = [...messages.value, ...newMessages]
            } else {
              messages.value = newMessages
            }

            // 检查是否还有更多数据
            hasMore.value = messages.value.length < total
          } else {
            hasMore.value = false
          }
        } else {
          ElMessage.error('获取系统消息失败: ' + response.data.message)
          hasMore.value = false
        }
      } catch (error) {
        console.error('获取系统消息失败:', error)
        // 如果API不存在，使用模拟数据（仅首次加载）
        if (!isLoadMore && messages.value.length === 0) {
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
          hasMore.value = false
        } else {
          hasMore.value = false
        }
      } finally {
        loading.value = false
        loadingMore.value = false
      }
    }

    // 加载更多数据
    const loadMore = async () => {
      if (loading.value || loadingMore.value || !hasMore.value) return

      currentPage.value += 1
      await fetchMessages(true)
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

    onMounted(async () => {
      // 初始加载数据
      await fetchMessages()

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
      messages,
      loading,
      loadingMore,
      hasMore,
      contentRef,
      getMessageType,
      getMessageTypeText,
      formatDate,
      markAsRead,
      goBack
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