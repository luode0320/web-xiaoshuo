<template>
  <div class="category-container">
    <div class="category-header">
      <h1>小说分类</h1>
      <p>浏览不同分类下的精彩小说</p>
    </div>
    
    <div class="category-content">
      <div class="category-grid">
        <div 
          v-for="category in categories" 
          :key="category.id" 
          class="category-card"
          @click="viewCategory(category.id)"
        >
          <div class="category-info">
            <h3>{{ category.name }}</h3>
            <p v-if="category.description">{{ category.description }}</p>
            <p class="novel-count">小说数量: {{ category.novel_count || 0 }}</p>
          </div>
          <div class="category-actions">
            <el-button type="primary" size="small">查看</el-button>
          </div>
        </div>
      </div>
      
      <!-- 子分类 -->
      <div v-if="subCategories.length > 0" class="sub-category-section">
        <h2>子分类</h2>
        <div class="sub-category-grid">
          <div 
            v-for="subCategory in subCategories" 
            :key="subCategory.id" 
            class="sub-category-item"
            @click="viewCategory(subCategory.id)"
          >
            <h4>{{ subCategory.name }}</h4>
            <p class="novel-count">小说: {{ subCategory.novel_count || 0 }}</p>
          </div>
        </div>
      </div>
      
      <!-- 分类下的小说列表 -->
      <div v-if="selectedCategory" class="novel-section">
        <h2>{{ selectedCategory.name }}下的小说</h2>
        <div class="novel-grid">
          <div 
            v-for="novel in categoryNovels" 
            :key="novel.id" 
            class="novel-card"
            @click="viewNovel(novel.id)"
          >
            <div class="novel-cover">
              <div class="cover-placeholder">封面</div>
            </div>
            <div class="novel-info">
              <h3>{{ novel.title }}</h3>
              <p class="author">作者: {{ novel.author }}</p>
              <p class="description">{{ truncateText(novel.description, 100) }}</p>
              <div class="stats">
                <span class="clicks">点击: {{ novel.click_count }}</span>
                <span class="rating">评分: {{ novel.avg_rating || '暂无' }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 分页 -->
        <div class="pagination" v-if="novelPagination.total > novelPagination.limit">
          <el-pagination
            v-model:current-page="novelPagination.page"
            :page-size="novelPagination.limit"
            :total="novelPagination.total"
            @current-change="handlePageChange"
            layout="prev, pager, next"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import apiClient from '@/utils/api'

export default {
  name: 'CategoryList',
  setup() {
    const route = useRouter()
    const router = useRouter()
    
    const categories = ref([])
    const subCategories = ref([])
    const categoryNovels = ref([])
    const selectedCategory = ref(null)
    
    const novelPagination = ref({
      page: 1,
      limit: 12,
      total: 0
    })
    
    // 获取所有分类
    const fetchCategories = async () => {
      try {
        const response = await apiClient.get('/api/v1/categories?with_children=true')
        categories.value = response.data.data.categories.map(cat => ({
          ...cat,
          novel_count: cat.novels ? cat.novels.length : 0
        }))
      } catch (error) {
        console.error('获取分类失败:', error)
        ElMessage.error('获取分类失败')
      }
    }
    
    // 获取分类下的小说
    const fetchCategoryNovels = async (categoryId) => {
      try {
        const response = await apiClient.get(`/api/v1/novels?category_id=${categoryId}&page=${novelPagination.value.page}&limit=${novelPagination.value.limit}`)
        categoryNovels.value = response.data.data.novels
        novelPagination.value.total = response.data.data.pagination.total
      } catch (error) {
        console.error('获取分类小说失败:', error)
        ElMessage.error('获取分类小说失败')
      }
    }
    
    // 查看分类详情
    const viewCategory = async (categoryId) => {
      try {
        const response = await apiClient.get(`/api/v1/categories/${categoryId}`)
        selectedCategory.value = response.data.data.category
        
        // 获取该分类下的小说
        novelPagination.value.page = 1
        await fetchCategoryNovels(categoryId)
      } catch (error) {
        console.error('获取分类详情失败:', error)
        ElMessage.error('获取分类详情失败')
      }
    }
    
    // 查看小说详情
    const viewNovel = (novelId) => {
      router.push(`/novel/${novelId}`)
    }
    
    // 处理分页变化
    const handlePageChange = async (page) => {
      novelPagination.value.page = page
      if (selectedCategory.value) {
        await fetchCategoryNovels(selectedCategory.value.id)
      }
    }
    
    // 截断文本
    const truncateText = (text, length) => {
      if (!text) return ''
      return text.length > length ? text.substring(0, length) + '...' : text
    }
    
    onMounted(async () => {
      await fetchCategories()
    })
    
    return {
      categories,
      subCategories,
      categoryNovels,
      selectedCategory,
      novelPagination,
      viewCategory,
      viewNovel,
      handlePageChange,
      truncateText
    }
  }
}
</script>

<style scoped>
.category-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.category-header {
  text-align: center;
  margin-bottom: 40px;
}

.category-header h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 10px;
}

.category-header p {
  color: #666;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.category-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
  cursor: pointer;
  transition: transform 0.3s ease;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.category-card:hover {
  transform: translateY(-5px);
}

.category-info h3 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 1.2rem;
}

.category-info p {
  margin: 5px 0;
  color: #666;
  font-size: 0.9rem;
}

.novel-count {
  color: #999;
  font-size: 0.8rem;
  margin-top: auto;
}

.category-actions {
  margin-top: 15px;
  text-align: right;
}

.sub-category-section {
  margin-bottom: 40px;
}

.sub-category-section h2 {
  font-size: 1.5rem;
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #409eff;
  padding-bottom: 10px;
}

.sub-category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.sub-category-item {
  background: white;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
}

.sub-category-item:hover {
  background: #409eff;
  color: white;
  transform: translateY(-2px);
}

.sub-category-item h4 {
  margin: 0 0 5px 0;
}

.novel-section h2 {
  font-size: 1.5rem;
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #409eff;
  padding-bottom: 10px;
}

.novel-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
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
  height: 150px;
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

.author, .description {
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

.pagination {
  text-align: center;
  margin-top: 30px;
}

@media (max-width: 768px) {
  .category-grid {
    grid-template-columns: 1fr;
  }
  
  .sub-category-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .novel-grid {
    grid-template-columns: 1fr;
  }
}
</style>