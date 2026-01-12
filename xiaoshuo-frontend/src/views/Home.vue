<template>
  <div class="home">
    <div class="banner">
      <h1>小说阅读系统</h1>
      <p>发现和阅读精彩小说</p>
    </div>
    
    <div class="container">
      <!-- 推荐小说 -->
      <div class="section">
        <h2>推荐小说</h2>
        <div class="novel-grid">
          <div 
            v-for="novel in recommendedNovels" 
            :key="novel.id" 
            class="novel-card"
            @click="goToNovelDetail(novel.id)"
          >
            <div class="novel-cover">
              <div class="cover-placeholder">封面</div>
            </div>
            <div class="novel-info">
              <h3>{{ novel.title }}</h3>
              <p class="author">作者: {{ novel.author }}</p>
              <p class="description">{{ novel.description }}</p>
              <div class="stats">
                <span class="clicks">点击: {{ novel.click_count }}</span>
                <span class="rating">评分: {{ novel.avg_rating || '暂无' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 分类导航 -->
      <div class="section">
        <h2>小说分类</h2>
        <div class="category-grid">
          <div 
            v-for="category in categories" 
            :key="category.id" 
            class="category-item"
            @click="goToCategory(category.id)"
          >
            {{ category.name }}
          </div>
        </div>
      </div>
      
      <!-- 排行榜 -->
      <div class="section">
        <h2>热门小说</h2>
        <div class="rank-list">
          <div 
            v-for="(novel, index) in hotNovels" 
            :key="novel.id" 
            class="rank-item"
            @click="goToNovelDetail(novel.id)"
          >
            <span class="rank-number">{{ index + 1 }}</span>
            <span class="novel-title">{{ novel.title }}</span>
            <span class="novel-author">{{ novel.author }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

export default {
  name: 'Home',
  setup() {
    const router = useRouter()
    const recommendedNovels = ref([])
    const categories = ref([])
    const hotNovels = ref([])
    
    // 获取推荐小说
    const fetchRecommendedNovels = async () => {
      try {
        const response = await axios.get('/api/v1/novels?limit=8')
        recommendedNovels.value = response.data.data.novels
      } catch (error) {
        console.error('获取推荐小说失败:', error)
      }
    }
    
    // 获取分类
    const fetchCategories = async () => {
      try {
        const response = await axios.get('/api/v1/categories')
        categories.value = response.data.data.categories
      } catch (error) {
        console.error('获取分类失败:', error)
      }
    }
    
    // 获取热门小说
    const fetchHotNovels = async () => {
      try {
        const response = await axios.get('/api/v1/rankings?type=total&limit=10')
        hotNovels.value = response.data.data.novels
      } catch (error) {
        console.error('获取热门小说失败:', error)
      }
    }
    
    const goToNovelDetail = (id) => {
      router.push(`/novel/${id}`)
    }
    
    const goToCategory = (id) => {
      router.push(`/category/${id}`)
    }
    
    onMounted(() => {
      fetchRecommendedNovels()
      fetchCategories()
      fetchHotNovels()
    })
    
    return {
      recommendedNovels,
      categories,
      hotNovels,
      goToNovelDetail,
      goToCategory
    }
  }
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-align: center;
  padding: 60px 20px;
}

.banner h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
}

.banner p {
  font-size: 1.2rem;
  opacity: 0.9;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.section {
  margin-bottom: 40px;
}

.section h2 {
  font-size: 1.5rem;
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #409eff;
  padding-bottom: 10px;
}

.novel-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
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

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 10px;
}

.category-item {
  background: white;
  padding: 15px;
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.1);
}

.category-item:hover {
  background: #409eff;
  color: white;
  transform: translateY(-2px);
}

.rank-list {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.rank-item {
  display: flex;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: background 0.3s ease;
}

.rank-item:last-child {
  border-bottom: none;
}

.rank-item:hover {
  background: #f5f5f5;
}

.rank-number {
  width: 30px;
  height: 30px;
  background: #409eff;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  font-size: 0.9rem;
}

.novel-title {
  flex: 1;
  font-weight: 500;
}

.novel-author {
  color: #999;
  font-size: 0.9rem;
}

@media (max-width: 768px) {
  .banner h1 {
    font-size: 2rem;
  }
  
  .novel-grid {
    grid-template-columns: 1fr;
  }
  
  .category-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}