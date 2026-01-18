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
          <div class="tab-content">
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

            <el-pagination v-model:current-page="commentsPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" :total="commentsTotal" layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleCommentsSizeChange" @current-change="handleCommentsCurrentChange" class="pagination" />
          </div>
        </el-tab-pane>

        <el-tab-pane label="我的评分" name="ratings">
          <div class="tab-content">
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

            <el-pagination v-model:current-page="ratingsPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" :total="ratingsTotal" layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleRatingsSizeChange" @current-change="handleRatingsCurrentChange" class="pagination" />
          </div>
        </el-tab-pane>
      </el-tabs>
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
  name: 'SocialHistory',
  components: {
    ArrowLeft
  },
  setup() {
    const router = useRouter()
    const activeTab = ref('comments')
    const comments = ref([])
    const ratings = ref([])
    const loading = ref({
      comments: false,
      ratings: false
    })
    const commentsPage = ref(1)
    const ratingsPage = ref(1)
    const pageSize = ref(10)
    const commentsTotal = ref(0)
    const ratingsTotal = ref(0)

    // 获取评论历史
    const fetchComments = async () => {
      loading.value.comments = true
      try {
        const response = await apiClient.get('/api/v1/comments', {
          params: {
            page: commentsPage.value,
            limit: pageSize.value
          }
        })

        if (response.data.code === 200) {
          comments.value = response.data.data.comments
          commentsTotal.value = response.data.data.pagination.total
        } else {
          ElMessage.error('获取评论历史失败: ' + response.data.message)
        }
      } catch (error) {
        console.error('获取评论历史失败:', error)
        ElMessage.error('获取评论历史失败: ' + error.message)
      } finally {
        loading.value.comments = false
      }
    }

    // 获取评分历史
    const fetchRatings = async () => {
      loading.value.ratings = true
      try {
        const response = await apiClient.get('/api/v1/users/ratings', {
          params: {
            page: ratingsPage.value,
            limit: pageSize.value
          }
        })

        if (response.data.code === 200) {
          ratings.value = response.data.data.ratings
          ratingsTotal.value = response.data.data.pagination.total
        } else {
          ElMessage.error('获取评分历史失败: ' + response.data.message)
        }
      } catch (error) {
        console.error('获取评分历史失败:', error)
        ElMessage.error('获取评分历史失败: ' + error.message)
      } finally {
        loading.value.ratings = false
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

    // 处理评论页面大小变化
    const handleCommentsSizeChange = (size) => {
      pageSize.value = size
      commentsPage.value = 1
      fetchComments()
    }

    // 处理评论当前页变化
    const handleCommentsCurrentChange = (page) => {
      commentsPage.value = page
      fetchComments()
    }

    // 处理评分页面大小变化
    const handleRatingsSizeChange = (size) => {
      pageSize.value = size
      ratingsPage.value = 1
      fetchRatings()
    }

    // 处理评分当前页变化
    const handleRatingsCurrentChange = (page) => {
      ratingsPage.value = page
      fetchRatings()
    }

    onMounted(() => {
      fetchComments()
      fetchRatings()
    })

    return {
      activeTab,
      comments,
      ratings,
      loading,
      commentsPage,
      ratingsPage,
      pageSize,
      commentsTotal,
      ratingsTotal,
      formatDate,
      viewNovel,
      goBack,
      handleCommentsSizeChange,
      handleCommentsCurrentChange,
      handleRatingsSizeChange,
      handleRatingsCurrentChange
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

.pagination {
  margin-top: 20px;
  justify-content: center;
  flex-shrink: 0;
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