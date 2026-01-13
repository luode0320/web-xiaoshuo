<template>
  <div class="search-container">
    <div class="search-header">
      <h1>小说搜索</h1>
      <div class="search-box">
        <el-autocomplete
          v-model="searchKeyword"
          :fetch-suggestions="querySearch"
          placeholder="搜索小说标题、作者、主角..."
          @keyup.enter="performSearch"
          @select="handleSelect"
          size="large"
          clearable
          :trigger-on-focus="false"
        >
          <template #default="{ item }">
            <div class="suggestion-item">
              <span class="suggestion-text">{{ item.text }}</span>
              <span class="suggestion-count" v-if="item.count">({{ item.count }})</span>
            </div>
          </template>
          <template #append>
            <el-button @click="performSearch" :icon="Search" :loading="searching" />
          </template>
        </el-autocomplete>
      </div>
      
      <!-- 搜索历史管理 -->
      <div v-if="isAuthenticated && searchHistory.length > 0" class="search-history">
        <div class="history-header">
          <h3>搜索历史</h3>
          <el-button 
            size="small" 
            type="danger" 
            @click="clearSearchHistory"
            :loading="clearingHistory"
          >
            清空历史
          </el-button>
        </div>
        <div class="history-tags">
          <el-tag
            v-for="history in searchHistory"
            :key="history.id"
            type="info"
            closable
            @close="removeSearchHistory(history.id)"
            @click="searchByKeyword(history.keyword)"
            class="history-tag"
          >
            {{ history.keyword }} ({{ history.count }})
          </el-tag>
        </div>
      </div>
    </div>
    
    <!-- 搜索统计展示区域 -->
    <div v-if="isAuthenticated && searchStats.total_searches" class="search-stats">
      <el-collapse>
        <el-collapse-item title="搜索统计" name="stats">
          <div class="stats-content">
            <div class="stat-item">
              <span class="stat-label">总搜索次数:</span>
              <span class="stat-value">{{ searchStats.total_searches }}</span>
            </div>
            <div class="stat-item" v-if="searchStats.top_keywords && searchStats.top_keywords.length > 0">
              <span class="stat-label">热门搜索:</span>
              <div class="top-keywords">
                <el-tag 
                  v-for="keyword in searchStats.top_keywords.slice(0, 5)" 
                  :key="keyword.id || keyword.keyword"
                  type="warning"
                  size="small"
                  class="keyword-tag"
                  @click="searchByKeyword(keyword.keyword || keyword.text)"
                >
                  {{ keyword.keyword || keyword.text }} ({{ keyword.count || keyword.Count }})
                </el-tag>
              </div>
            </div>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
    
    <div class="search-filters">
      <el-collapse v-model="activeFilter">
        <el-collapse-item title="筛选条件" name="filters">
          <div class="filter-options">
            <div class="filter-group">
              <label>分类:</label>
              <el-select 
                v-model="searchFilters.categoryId" 
                placeholder="选择分类" 
                clearable
                @change="performSearch"
              >
                <el-option
                  v-for="category in categories"
                  :key="category.id"
                  :label="category.name"
                  :value="category.id"
                />
              </el-select>
            </div>
            
            <div class="filter-group">
              <label>评分范围:</label>
              <el-slider
                v-model="searchFilters.ratingRange"
                range
                :max="10"
                :step="0.5"
                @change="performSearch"
                style="width: 200px;"
              />
              <div class="rating-values">
                <span>{{ searchFilters.ratingRange[0] }} - {{ searchFilters.ratingRange[1] }}</span>
              </div>
            </div>
            
            <div class="filter-group">
              <label>排序方式:</label>
              <el-select 
                v-model="searchFilters.orderBy" 
                placeholder="排序方式" 
                @change="performSearch"
              >
                <el-option label="相关度" value="relevance" />
                <el-option label="点击量" value="clicks" />
                <el-option label="评分" value="rating" />
                <el-option label="上传时间" value="time" />
              </el-select>
            </div>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
    
    <div class="search-results">
      <div v-if="searchResults.length > 0" class="results-header">
        <p>找到 {{ totalResults }} 个结果，用时 {{ searchTime }} 秒</p>
      </div>
      
      <div v-if="searchResults.length > 0" class="novel-grid">
        <div 
          v-for="novel in searchResults" 
          :key="novel.id" 
          class="novel-card"
          @click="viewNovel(novel.id)"
        >
          <div class="novel-cover">
            <div class="cover-placeholder">封面</div>
          </div>
          <div class="novel-info">
            <h3>{{ highlightText(novel.title) }}</h3>
            <p class="author">作者: {{ highlightText(novel.author) }}</p>
            <p class="protagonist" v-if="novel.protagonist">主角: {{ highlightText(novel.protagonist) }}</p>
            <p class="description">{{ highlightText(truncateText(novel.description, 150)) }}</p>
            <div class="stats">
              <span class="clicks">点击: {{ novel.click_count }}</span>
              <span class="rating">评分: {{ novel.avg_rating || '暂无' }}</span>
              <span class="category" v-if="novel.categories && novel.categories.length > 0">
                {{ novel.categories[0].name }}
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else-if="!searchKeyword && !hasSearched" class="search-tips">
        <h3>热门搜索</h3>
        <div class="hot-keywords">
          <el-tag 
            v-for="keyword in hotKeywords" 
            :key="keyword" 
            type="info" 
            size="large" 
            class="keyword-tag"
            @click="searchByKeyword(keyword)"
          >
            {{ keyword }}
          </el-tag>
        </div>
        <h3 style="margin-top: 30px;">推荐小说</h3>
        <div class="recommended-novels">
          <div 
            v-for="novel in recommendedNovels" 
            :key="novel.id" 
            class="recommended-item"
            @click="viewNovel(novel.id)"
          >
            <div class="recommended-cover">
              <div class="cover-placeholder">封面</div>
            </div>
            <div class="recommended-info">
              <h4>{{ novel.title }}</h4>
              <p>{{ novel.author }}</p>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else-if="hasSearched && searchResults.length === 0" class="no-results">
        <p>没有找到相关小说</p>
        <p>请尝试使用其他关键词搜索</p>
      </div>
    </div>
    
    <!-- 分页 -->
    <div class="pagination" v-if="totalResults > searchFilters.limit">
      <el-pagination
        v-model:current-page="searchFilters.page"
        :page-size="searchFilters.limit"
        :total="totalResults"
        @current-change="handlePageChange"
        layout="prev, pager, next, jumper"
      />
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import axios from 'axios'
import { useUserStore } from '@/stores/user'

export default {
  name: 'SearchList',
  components: {
    Search
  },
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const searchKeyword = ref('')
    const searching = ref(false)
    const searchResults = ref([])
    const categories = ref([])
    const hotKeywords = ref([])
    const recommendedNovels = ref([])
    const totalResults = ref(0)
    const searchTime = ref(0)
    const hasSearched = ref(false)
    const activeFilter = ref([])
    const searchSuggestions = ref([])
    
    // 搜索历史相关
    const searchHistory = ref([])
    const showHistoryPanel = ref(false)
    const searchStats = ref({})
    const clearingHistory = ref(false)
    
    const searchFilters = reactive({
      page: 1,
      limit: 20,
      categoryId: '',
      ratingRange: [0, 10],
      orderBy: 'relevance'
    })
    
    // 计算属性
    const isAuthenticated = computed(() => userStore.isAuthenticated)
    
    // 查询搜索建议
    const querySearch = async (queryString, callback) => {
      if (!queryString) {
        callback([])
        return
      }

      try {
        const response = await axios.get(`/api/v1/search/suggestions?q=${queryString}`)
        const suggestions = response.data.data.suggestions.map(item => ({
          value: item.text,
          text: item.text,
          count: item.count
        }))
        callback(suggestions)
      } catch (error) {
        // 如果API失败，返回空结果
        callback([])
      }
    }
    
    // 选择搜索建议
    const handleSelect = (item) => {
      searchKeyword.value = item.value
      performSearch()
    }
    
    // 执行搜索
    const performSearch = async () => {
      if (!searchKeyword.value.trim()) return
      
      try {
        searching.value = true
        hasSearched.value = true
        
        const startTime = Date.now()
        
        const params = {
          q: searchKeyword.value,
          page: searchFilters.page,
          limit: searchFilters.limit
        }
        
        if (searchFilters.categoryId) {
          params.category_id = searchFilters.categoryId
        }
        
        if (searchFilters.ratingRange[0] > 0 || searchFilters.ratingRange[1] < 10) {
          params.min_score = searchFilters.ratingRange[0]
          params.max_score = searchFilters.ratingRange[1]
        }
        
        const queryString = new URLSearchParams(params).toString()
        const response = await axios.get(`/api/v1/search/novels?${queryString}`)
        
        searchResults.value = response.data.data.novels
        totalResults.value = response.data.data.pagination.total
        searchTime.value = ((Date.now() - startTime) / 1000).toFixed(2)
      } catch (error) {
        console.error('搜索失败:', error)
        ElMessage.error('搜索失败')
      } finally {
        searching.value = false
      }
    }
    
    // 搜索指定关键词
    const searchByKeyword = (keyword) => {
      searchKeyword.value = keyword
      searchFilters.page = 1
      performSearch()
    }
    
    // 获取分类列表
    const fetchCategories = async () => {
      try {
        const response = await axios.get('/api/v1/categories')
        categories.value = response.data.data.categories
      } catch (error) {
        console.error('获取分类失败:', error)
      }
    }
    
    // 获取热门搜索关键词
    const fetchHotKeywords = async () => {
      try {
        const response = await axios.get('/api/v1/search/hot-words')
        hotKeywords.value = response.data.data.keywords
      } catch (error) {
        console.error('获取热门关键词失败:', error)
        // 如果API失败，使用默认关键词
        hotKeywords.value = ['玄幻', '都市', '科幻', '言情', '武侠', '历史', '军事', '悬疑']
      }
    }
    
    // 获取推荐小说
    const fetchRecommendedNovels = async () => {
      try {
        const response = await axios.get('/api/v1/novels?limit=8')
        recommendedNovels.value = response.data.data.novels
      } catch (error) {
        console.error('获取推荐小说失败:', error)
      }
    }
    
    // 查看小说详情
    const viewNovel = (novelId) => {
      router.push(`/novel/${novelId}`)
    }
    
    // 处理分页变化
    const handlePageChange = async (page) => {
      searchFilters.page = page
      await performSearch()
    }
    
    // 截断文本
    const truncateText = (text, length) => {
      if (!text) return ''
      return text.length > length ? text.substring(0, length) + '...' : text
    }
    
    // 高亮搜索关键词
    const highlightText = (text) => {
      if (!text || !searchKeyword.value) return text
      
      const keyword = searchKeyword.value.toLowerCase()
      const lowerText = text.toLowerCase()
      const index = lowerText.indexOf(keyword)
      
      if (index === -1) return text
      
      const before = text.substring(0, index)
      const matched = text.substring(index, index + keyword.length)
      const after = text.substring(index + keyword.length)
      
      return `${before}<mark>${matched}</mark>${after}`
    }
    
    // 获取搜索历史
    const fetchSearchHistory = async () => {
      if (!isAuthenticated.value) return
      
      try {
        const response = await axios.get('/api/v1/users/search-history', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        searchHistory.value = response.data.data.search_history
      } catch (error) {
        console.error('获取搜索历史失败:', error)
      }
    }
    
    // 清空搜索历史
    const clearSearchHistory = async () => {
      if (!isAuthenticated.value) return
      
      try {
        await ElMessageBox.confirm(
          '确定要清空搜索历史吗？此操作不可恢复。',
          '确认清空',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        clearingHistory.value = true
        
        await axios.delete('/api/v1/users/search-history', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        searchHistory.value = []
        ElMessage.success('搜索历史已清空')
      } catch (error) {
        if (error !== 'cancel') {
          console.error('清空搜索历史失败:', error)
          ElMessage.error('清空搜索历史失败')
        }
      } finally {
        clearingHistory.value = false
      }
    }
    
    // 删除单个搜索历史
    const removeSearchHistory = async (historyId) => {
      if (!isAuthenticated.value) return
      
      try {
        // 从数组中移除，实际API中没有提供删除单个历史记录的端点
        // 可能需要通过后端API添加此功能，这里先在前端移除
        searchHistory.value = searchHistory.value.filter(h => h.id !== historyId)
      } catch (error) {
        console.error('删除搜索历史失败:', error)
        // 重新获取历史
        await fetchSearchHistory()
      }
    }
    
    onMounted(async () => {
      await fetchCategories()
      await fetchHotKeywords()
      await fetchRecommendedNovels()
      if (isAuthenticated.value) {
        await fetchSearchHistory()
      }
    })
    
    return {
      searchKeyword,
      searching,
      searchResults,
      categories,
      hotKeywords,
      recommendedNovels,
      totalResults,
      searchTime,
      hasSearched,
      activeFilter,
      searchFilters,
      searchHistory,
      searchStats,
      clearingHistory,
      performSearch,
      searchByKeyword,
      viewNovel,
      handlePageChange,
      truncateText,
      highlightText,
      isAuthenticated,
      fetchSearchHistory,
      clearSearchHistory,
      removeSearchHistory,
      fetchSearchStats
    }
  }
}
</script>

<style scoped>
.search-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-header {
  text-align: center;
  margin-bottom: 30px;
}

.search-header h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 20px;
}

.search-box {
  max-width: 600px;
  margin: 0 auto 20px;
}

.search-history {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 15px;
  margin-top: 15px;
  border: 1px solid #eee;
}

/* 搜索统计样式 */
.search-stats {
  margin: 20px 0;
  padding: 15px;
  background: var(--el-bg-color-overlay);
  border-radius: 8px;
  border: 1px solid var(--el-border-color);
}

.stats-content {
  padding: 10px 0;
}

.stat-item {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.stat-label {
  font-weight: bold;
  margin-right: 10px;
  min-width: 100px;
}

.stat-value {
  font-size: 16px;
  color: var(--el-color-primary);
  font-weight: bold;
}

.top-keywords {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.keyword-tag {
  cursor: pointer;
  transition: all 0.3s;
}

.keyword-tag:hover {
  transform: scale(1.05);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.history-header h3 {
  margin: 0;
  color: #333;
  font-size: 1rem;
}

.history-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.history-tag {
  cursor: pointer;
  transition: all 0.3s ease;
}

.history-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.suggestion-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px 0;
}

.suggestion-text {
  flex: 1;
}

.suggestion-count {
  color: #999;
  font-size: 0.9em;
}

.search-filters {
  margin-bottom: 30px;
}

.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  padding: 20px 0;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-group label {
  font-weight: 500;
  color: #333;
  white-space: nowrap;
}

.rating-values {
  font-size: 0.9rem;
  color: #666;
}

.search-results {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
  padding: 20px;
}

.results-header {
  padding: 10px 0 20px 0;
  border-bottom: 1px solid #eee;
  color: #666;
}

.novel-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin: 20px 0;
}

.novel-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.novel-card:hover {
  transform: translateY(-5px);
}

.novel-cover {
  height: 200px;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cover-placeholder {
  font-size: 0.8rem;
  color: #999;
}

.novel-info {
  padding: 15px;
}

.novel-info h3 {
  margin: 0 0 10px 0;
  font-size: 1.1rem;
  color: #333;
}

.author, .protagonist, .description {
  margin: 5px 0;
  font-size: 0.9rem;
  color: #666;
}

.description {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.stats {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 0.8rem;
  color: #999;
}

:deep(mark) {
  background-color: #ffe66d;
  padding: 0 2px;
}

.search-tips {
  text-align: center;
  padding: 40px 20px;
}

.hot-keywords {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 10px;
  margin: 20px 0;
}

.keyword-tag {
  cursor: pointer;
  transition: all 0.3s ease;
}

.keyword-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.recommended-novels {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
  margin-top: 20px;
}

.recommended-item {
  display: flex;
  gap: 10px;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.recommended-item:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.recommended-cover {
  width: 40px;
  height: 60px;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.recommended-info {
  flex: 1;
  text-align: left;
}

.recommended-info h4 {
  margin: 0 0 5px 0;
  font-size: 0.9rem;
  color: #333;
}

.recommended-info p {
  margin: 0;
  font-size: 0.8rem;
  color: #666;
}

.no-results {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.no-results p:first-child {
  font-size: 1.2rem;
  margin-bottom: 10px;
}

.pagination {
  text-align: center;
  margin-top: 30px;
}

@media (max-width: 768px) {
  .search-header h1 {
    font-size: 1.5rem;
  }
  
  .filter-options {
    flex-direction: column;
  }
  
  .novel-grid {
    grid-template-columns: 1fr;
  }
  
  .recommended-novels {
    grid-template-columns: 1fr;
  }
}