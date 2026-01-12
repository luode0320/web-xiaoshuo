<template>
  <div class="ranking-container">
    <div class="ranking-header">
      <h1>小说排行榜</h1>
      <p>发现最受欢迎的小说</p>
    </div>
    
    <div class="ranking-filters">
      <el-radio-group v-model="rankingType" @change="fetchRankings">
        <el-radio-button label="total">总榜</el-radio-button>
        <el-radio-button label="today">日榜</el-radio-button>
        <el-radio-button label="week">周榜</el-radio-button>
        <el-radio-button label="month">月榜</el-radio-button>
      </el-radio-group>
      
      <el-select 
        v-model="selectedCategory" 
        placeholder="选择分类" 
        @change="fetchRankings"
        clearable
        style="margin-left: 20px; width: 200px;"
      >
        <el-option
          v-for="category in categories"
          :key="category.id"
          :label="category.name"
          :value="category.id"
        />
      </el-select>
    </div>
    
    <div class="ranking-content">
      <div class="ranking-list">
        <div 
          v-for="(novel, index) in rankedNovels" 
          :key="novel.id" 
          class="ranking-item"
          :class="{ 'top-3': index < 3 }"
          @click="viewNovel(novel.id)"
        >
          <div class="rank-number" :class="{ 'top-3': index < 3 }">
            <span class="rank">{{ index + 1 }}</span>
          </div>
          <div class="novel-info">
            <h3>{{ novel.title }}</h3>
            <p class="author">作者: {{ novel.author }}</p>
            <p class="stats">
              <span class="clicks">点击: {{ getClickCount(novel) }}</span>
              <span class="rating">评分: {{ novel.avg_rating || '暂无' }}</span>
            </p>
          </div>
          <div class="novel-cover">
            <div class="cover-placeholder">封面</div>
          </div>
        </div>
      </div>
      
      <div v-if="rankedNovels.length === 0" class="no-data">
        <p>暂无排行榜数据</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
  name: 'RankingList',
  setup() {
    const router = useRouter()
    
    const rankingType = ref('total')
    const selectedCategory = ref('')
    const rankedNovels = ref([])
    const categories = ref([])
    
    // 获取排行榜数据
    const fetchRankings = async () => {
      try {
        const params = {
          type: rankingType.value,
          limit: 50
        }
        
        if (selectedCategory.value) {
          params.category_id = selectedCategory.value
        }
        
        const queryString = new URLSearchParams(params).toString()
        const response = await axios.get(`/api/v1/rankings?${queryString}`)
        rankedNovels.value = response.data.data.novels
      } catch (error) {
        console.error('获取排行榜失败:', error)
        ElMessage.error('获取排行榜失败')
      }
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
    
    // 获取点击数（根据排行榜类型）
    const getClickCount = (novel) => {
      switch (rankingType.value) {
        case 'today':
          return novel.today_clicks
        case 'week':
          return novel.week_clicks
        case 'month':
          return novel.month_clicks
        default:
          return novel.click_count
      }
    }
    
    // 查看小说详情
    const viewNovel = (novelId) => {
      router.push(`/novel/${novelId}`)
    }
    
    onMounted(async () => {
      await fetchCategories()
      await fetchRankings()
    })
    
    return {
      rankingType,
      selectedCategory,
      rankedNovels,
      categories,
      fetchRankings,
      getClickCount,
      viewNovel
    }
  }
}
</script>

<style scoped>
.ranking-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.ranking-header {
  text-align: center;
  margin-bottom: 30px;
}

.ranking-header h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 10px;
}

.ranking-header p {
  color: #666;
}

.ranking-filters {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 15px;
}

.ranking-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.ranking-list {
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}

.ranking-item {
  display: flex;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: background 0.3s ease;
}

.ranking-item:last-child {
  border-bottom: none;
}

.ranking-item:hover {
  background: #f5f5f5;
}

.ranking-item.top-3 {
  background: linear-gradient(90deg, #fff9e6 0%, #ffffff 100%);
}

.rank-number {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  font-size: 1.2rem;
  font-weight: bold;
}

.rank-number.top-3 {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  color: white;
}

.rank-number .rank {
  color: #999;
}

.ranking-item.top-3 .rank-number .rank {
  color: white;
}

.ranking-item.top-3:nth-child(1) .rank-number {
  background: linear-gradient(135deg, #ff6b6b, #ff8e53);
}

.ranking-item.top-3:nth-child(2) .rank-number {
  background: linear-gradient(135deg, #6bffb8, #53c2ff);
}

.ranking-item.top-3:nth-child(3) .rank-number {
  background: linear-gradient(135deg, #c28bff, #8e53ff);
}

.novel-info {
  flex: 1;
}

.novel-info h3 {
  margin: 0 0 8px 0;
  font-size: 1.1rem;
  color: #333;
}

.author, .stats {
  margin: 4px 0;
  font-size: 0.9rem;
  color: #666;
}

.stats {
  display: flex;
  gap: 15px;
}

.novel-cover {
  width: 60px;
  height: 80px;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.cover-placeholder {
  font-size: 0.7rem;
  color: #999;
}

.no-data {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

@media (max-width: 768px) {
  .ranking-filters {
    flex-direction: column;
    align-items: stretch;
  }
  
  .ranking-filters .el-select {
    margin-left: 0 !important;
    margin-top: 10px;
  }
  
  .ranking-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .rank-number {
    margin-right: 0;
    margin-bottom: 10px;
  }
  
  .novel-cover {
    align-self: flex-end;
    margin-top: 10px;
  }
}