<template>
  <div class="about-container">
    <div class="about-header">
      <el-button 
        class="back-button" 
        @click="goBack"
        icon="ArrowLeft"
        :plain="true"
      >
        返回
      </el-button>
      <h1>关于小说阅读系统</h1>
      <p>一个专注于提供优质阅读体验的平台</p>
    </div>
    
    <div class="about-content">
      <div class="section">
        <h2>项目简介</h2>
        <p>
          小说阅读系统是一个全栈Web应用，旨在为用户提供一个方便、高效的小说阅读和分享平台。
          系统支持多种格式的小说上传、在线阅读、评论评分、分类管理等功能。
        </p>
      </div>
      
      <div class="section">
        <h2>技术栈</h2>
        <div class="tech-stack">
          <div class="tech-item">
            <h4>后端技术</h4>
            <ul>
              <li>Go语言</li>
              <li>Gin Web框架</li>
              <li>GORM ORM</li>
              <li>MySQL数据库</li>
              <li>Redis缓存</li>
              <li>JWT认证</li>
            </ul>
          </div>
          <div class="tech-item">
            <h4>前端技术</h4>
            <ul>
              <li>Vue.js 3</li>
              <li>Vue Router</li>
              <li>Pinia状态管理</li>
              <li>Element Plus UI组件库</li>
              <li>Vite构建工具</li>
              <li>Axios HTTP客户端</li>
            </ul>
          </div>
        </div>
      </div>
      
      <div class="section">
        <h2>主要功能</h2>
        <div class="features">
          <div class="feature-item">
            <div class="feature-icon">📚</div>
            <div class="feature-content">
              <h3>小说管理</h3>
              <p>支持多种格式小说上传，分类管理，审核机制</p>
            </div>
          </div>
          <div class="feature-item">
            <div class="feature-icon">📖</div>
            <div class="feature-content">
              <h3>在线阅读</h3>
              <p>支持在线阅读，个性化阅读设置，进度保存</p>
            </div>
          </div>
          <div class="feature-item">
            <div class="feature-icon">💬</div>
            <div class="feature-content">
              <h3>社交互动</h3>
              <p>评论系统，评分功能，点赞机制</p>
            </div>
          </div>
          <div class="feature-item">
            <div class="feature-icon">🔍</div>
            <div class="feature-content">
              <h3>智能搜索</h3>
              <p>全文搜索，分类筛选，高级搜索</p>
            </div>
          </div>
          <div class="feature-item">
            <div class="feature-icon">🏆</div>
            <div class="feature-content">
              <h3>排行榜</h3>
              <p>多种排行榜，发现热门小说</p>
            </div>
          </div>
          <div class="feature-item">
            <div class="feature-icon">🔒</div>
            <div class="feature-content">
              <h3>用户系统</h3>
              <p>用户注册登录，个人中心，权限管理</p>
            </div>
          </div>
        </div>
      </div>
      
      <div class="section">
        <h2>系统架构</h2>
        <p>
          本系统采用前后端分离架构，后端提供RESTful API接口，前端通过Vue.js构建单页面应用。
          系统包含用户管理、小说管理、评论评分、分类管理、搜索推荐等多个模块。
        </p>
      </div>
      
      <div class="section">
        <h2>联系我们</h2>
        <p>如果您有任何问题或建议，欢迎通过以下方式联系我们：</p>
        <ul>
          <li>邮箱: support@example.com</li>
          <li>GitHub: <a href="https://github.com/your-repo" target="_blank">项目地址</a></li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'

export default {
  name: 'About',
  setup() {
    const router = useRouter()
    
    const goBack = () => {
      // 检查是否存在历史记录，并判断来源页面
      if (window.history.length > 1) {
        // 如果来自 profile 相关页面，则返回 profile 页面，否则返回上一页
        const prevRoute = history.state?.back
        if (prevRoute && (prevRoute.includes('/profile') || prevRoute === '/profile')) {
          router.push('/profile')
        } else {
          router.go(-1)
        }
      } else {
        // 如果没有历史记录，则返回首页
        router.push('/')
      }
    }
    
    return {
      goBack
    }
  },
  components: {
    ArrowLeft
  }
}
</script>

<style scoped>
.about-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.about-header {
  position: relative;
  text-align: center;
  padding-bottom: 30px;
  border-bottom: 1px solid #eee;
  margin-bottom: 30px;
}

.back-button {
  position: absolute;
  left: 0;
  top: 0;
  margin: 0;
}

.about-header {
  text-align: center;
  padding-bottom: 30px;
  border-bottom: 1px solid #eee;
  margin-bottom: 30px;
}

.about-header h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 10px;
}

.about-header p {
  color: #666;
}

.section {
  margin-bottom: 40px;
}

.section h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 20px;
  border-bottom: 2px solid #409eff;
  padding-bottom: 10px;
}

.section p {
  line-height: 1.8;
  color: #333;
  margin-bottom: 15px;
}

.tech-stack {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
}

.tech-item h4 {
  color: #333;
  margin-bottom: 15px;
  font-size: 1.1rem;
}

.tech-item ul {
  list-style: none;
  padding: 0;
}

.tech-item li {
  padding: 5px 0;
  color: #666;
  border-left: 3px solid #409eff;
  padding-left: 15px;
  margin-bottom: 8px;
}

.features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 25px;
}

.feature-item {
  display: flex;
  gap: 15px;
  padding: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.feature-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.feature-icon {
  font-size: 2rem;
  flex-shrink: 0;
}

.feature-content h3 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 1.1rem;
}

.feature-content p {
  color: #666;
  margin: 0;
  line-height: 1.6;
}

.section ul {
  padding-left: 20px;
  color: #666;
  line-height: 1.8;
}

.section a {
  color: #409eff;
  text-decoration: none;
}

.section a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .about-container {
    padding: 15px;
  }
  
  .about-header h1 {
    font-size: 1.5rem;
  }
  
  .tech-stack {
    grid-template-columns: 1fr;
  }
  
  .features {
    grid-template-columns: 1fr;
  }
}
</style>