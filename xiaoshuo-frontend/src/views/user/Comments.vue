<template>
  <div class="comments-container">
    <div class="header">
      <el-button type="primary" link @click="goBack" class="back-button">
        <el-icon>
          <ArrowLeft />
        </el-icon>
      </el-button>
      <h2>我的评论</h2>
    </div>

    <div class="content" ref="contentRef">
      <el-table :data="comments" style="width: 100%" v-loading="loading">
        <el-table-column prop="novel.title" label="小说标题" />
        <el-table-column prop="content" label="评论内容" show-overflow-tooltip />
        <el-table-column prop="like_count" label="点赞数" width="80" />
        <el-table-column prop="created_at" label="评论时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="viewNovel(row.novel_id)" type="primary">
              查看小说
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 加载更多提示 -->
      <div v-if="hasMore" class="loading-more">
        <el-skeleton v-if="loadingMore" :rows="2" />
      </div>

      <!-- 没有更多数据提示 -->
      <div v-if="!hasMore && comments.length > 0" class="no-more">
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
  name: 'Comments',
  components: {
    ArrowLeft
  },
  setup() {
    const router = useRouter()
    const comments = ref([])
    const loading = ref(false)
    const loadingMore = ref(false)
    const hasMore = ref(true)
    const currentPage = ref(1)
    const pageSize = ref(10)
    const contentRef = ref(null) // 用于滚动监听

    // 获取评论历史
    const fetchComments = async (isLoadMore = false) => {
      if (isLoadMore) {
        loadingMore.value = true
      } else {
        loading.value = true
      }

      try {
        const response = await apiClient.get('/api/v1/comments', {
          params: {
            page: currentPage.value,
            limit: pageSize.value
          }
        })

        if (response.data.code === 200) {
          const newComments = response.data.data.comments
          const total = response.data.data.pagination.total || 0

          if (newComments && newComments.length > 0) {
            if (isLoadMore) {
              comments.value = [...comments.value, ...newComments]
            } else {
              comments.value = newComments
            }

            // 检查是否还有更多数据
            hasMore.value = comments.value.length < total
          } else {
            hasMore.value = false
          }
        } else {
          ElMessage.error('获取评论历史失败: ' + response.data.message)
          hasMore.value = false
        }
      } catch (error) {
        console.error('获取评论历史失败:', error)
        ElMessage.error('获取评论历史失败: ' + error.message)
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
      await fetchComments(true)
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
      await fetchComments()

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
      comments,
      loading,
      loadingMore,
      hasMore,
      contentRef,
      formatDate,
      viewNovel,
      goBack
    }
  }
}
</script>

<style scoped>
.comments-container {
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
  .comments-container {
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