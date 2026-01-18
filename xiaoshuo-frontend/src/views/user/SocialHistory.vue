<template>
  <div class="social-container">
    <div class="header">
      <el-button type="primary" link @click="goBack" class="back-button">
        <el-icon>
          <ArrowLeft />
        </el-icon>
      </el-button>
      <h2>社交历史</h2>
    </div>

    <div class="content">
      <el-tabs v-model="activeTab" class="tabs">
        <el-tab-pane label="我的评论" name="comments">
          <div class="tab-content" ref="commentsContentRef">
            <el-table :data="comments" style="width: 100%" v-loading="loading.comments">
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
            <div v-if="commentsHasMore" class="loading-more">
              <el-skeleton v-if="loading.comments" :rows="2" />
            </div>

            <!-- 没有更多数据提示 -->
            <div v-if="!commentsHasMore && comments.length > 0" class="no-more">
              <el-divider>没有更多数据了</el-divider>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="我的评分" name="ratings">
          <div class="tab-content" ref="ratingsContentRef">
            <el-table :data="ratings" style="width: 100%" v-loading="loading.ratings">
              <el-table-column prop="novel.title" label="小说标题" />
              <el-table-column prop="score" label="评分" width="80">
                <template #default="{ row }">
                  <el-rate v-model="row.score" disabled :max="5" allow-half />
                </template>
              </el-table-column>
              <el-table-column prop="comment" label="评分说明" show-overflow-tooltip />
              <el-table-column prop="like_count" label="点赞数" width="80" />
              <el-table-column prop="created_at" label="评分时间" width="180">
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
            <div v-if="ratingsHasMore" class="loading-more">
              <el-skeleton v-if="loading.ratings" :rows="2" />
            </div>

            <!-- 没有更多数据提示 -->
            <div v-if="!ratingsHasMore && ratings.length > 0" class="no-more">
              <el-divider>没有更多数据了</el-divider>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'

export default {
  name: 'SocialHistory',
  components: {
    ArrowLeft
  },
  setup() {
    const router = useRouter()
    const activeTab = ref('comments')

    // 评论相关
    const comments = ref([])
    const commentsPage = ref(1)
    const commentsHasMore = ref(true)
    const commentsContentRef = ref(null)

    // 评分相关
    const ratings = ref([])
    const ratingsPage = ref(1)
    const ratingsHasMore = ref(true)
    const ratingsContentRef = ref(null)

    // 全局加载状态
    const loading = ref({
      comments: false,
      ratings: false
    })

    const pageSize = ref(10)

    // 获取评论历史
    const fetchComments = async (isLoadMore = false) => {
      if (isLoadMore) {
        loading.value.comments = true
      } else {
        commentsPage.value = 1
        loading.value.comments = true
      }

      try {
        const response = await apiClient.get('/api/v1/comments', {
          params: {
            page: commentsPage.value,
            limit: pageSize.value
          }
        })

        if (response.data.code === 200) {
          const newComments = response.data.data.comments
          const total = response.data.data.pagination.total || 0

          if (newComments && newComments.length > 0) {
            if (isLoadMore && commentsPage.value > 1) {
              comments.value = [...comments.value, ...newComments]
            } else {
              comments.value = newComments
            }

            // 检查是否还有更多数据
            commentsHasMore.value = comments.value.length < total
          } else {
            commentsHasMore.value = false
          }
        } else {
          ElMessage.error('获取评论历史失败: ' + response.data.message)
          commentsHasMore.value = false
        }
      } catch (error) {
        console.error('获取评论历史失败:', error)
        ElMessage.error('获取评论历史失败: ' + error.message)
        commentsHasMore.value = false
      } finally {
        loading.value.comments = false
      }
    }

    // 获取评分历史
    const fetchRatings = async (isLoadMore = false) => {
      if (isLoadMore) {
        loading.value.ratings = true
      } else {
        ratingsPage.value = 1
        loading.value.ratings = true
      }

      try {
        const response = await apiClient.get('/api/v1/users/ratings', {
          params: {
            page: ratingsPage.value,
            limit: pageSize.value
          }
        })

        if (response.data.code === 200) {
          const newRatings = response.data.data.ratings
          const total = response.data.data.pagination.total || 0

          if (newRatings && newRatings.length > 0) {
            if (isLoadMore && ratingsPage.value > 1) {
              ratings.value = [...ratings.value, ...newRatings]
            } else {
              ratings.value = newRatings
            }

            // 检查是否还有更多数据
            ratingsHasMore.value = ratings.value.length < total
          } else {
            ratingsHasMore.value = false
          }
        } else {
          ElMessage.error('获取评分历史失败: ' + response.data.message)
          ratingsHasMore.value = false
        }
      } catch (error) {
        console.error('获取评分历史失败:', error)
        ElMessage.error('获取评分历史失败: ' + error.message)
        ratingsHasMore.value = false
      } finally {
        loading.value.ratings = false
      }
    }

    // 评论加载更多数据
    const loadMoreComments = async () => {
      if (loading.value.comments || !commentsHasMore.value) return

      commentsPage.value += 1
      await fetchComments(true)
    }

    // 评分加载更多数据
    const loadMoreRatings = async () => {
      if (loading.value.ratings || !ratingsHasMore.value) return

      ratingsPage.value += 1
      await fetchRatings(true)
    }

    // 滚动事件监听
    const handleCommentsScroll = () => {
      if (!commentsContentRef.value || loading.value.comments || !commentsHasMore.value) return

      const element = commentsContentRef.value
      // 检查是否滚动到底部
      if (element.scrollHeight - element.scrollTop <= element.clientHeight + 100) {
        loadMoreComments()
      }
    }

    const handleRatingsScroll = () => {
      if (!ratingsContentRef.value || loading.value.ratings || !ratingsHasMore.value) return

      const element = ratingsContentRef.value
      // 检查是否滚动到底部
      if (element.scrollHeight - element.scrollTop <= element.clientHeight + 100) {
        loadMoreRatings()
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

    // 标签切换时的处理
    const handleTabChange = (tabName) => {
      if (tabName === 'comments' && comments.value.length === 0) {
        fetchComments()
      } else if (tabName === 'ratings' && ratings.value.length === 0) {
        fetchRatings()
      }
    }

    onMounted(() => {
      // 初始加载评论数据
      fetchComments()

      // 为两个内容区域添加滚动事件监听
      if (commentsContentRef.value) {
        commentsContentRef.value.addEventListener('scroll', handleCommentsScroll)
      }

      if (ratingsContentRef.value) {
        ratingsContentRef.value.addEventListener('scroll', handleRatingsScroll)
      }
    })

    // 监听标签变化
    watch(activeTab, (newTab) => {
      handleTabChange(newTab)
    })

    // 组件卸载时移除事件监听
    onUnmounted(() => {
      if (commentsContentRef.value) {
        commentsContentRef.value.removeEventListener('scroll', handleCommentsScroll)
      }

      if (ratingsContentRef.value) {
        ratingsContentRef.value.removeEventListener('scroll', handleRatingsScroll)
      }
    })

    return {
      activeTab,
      comments,
      ratings,
      loading,
      commentsHasMore,
      ratingsHasMore,
      commentsContentRef,
      ratingsContentRef,
      formatDate,
      viewNovel,
      goBack
    }
  }
}
</script>

<style scoped>
.social-container {
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

.tabs {
  width: 100%;
}

.tab-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
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
  .social-container {
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