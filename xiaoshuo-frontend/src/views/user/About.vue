<template>
  <div class="about-container">
    <div class="header">
      <el-button type="primary" link @click="goBack" class="back-button">
        <el-icon>
          <ArrowLeft />
        </el-icon>
      </el-button>
    </div>

    <div class="about-content">
      <div class="section">
        <h2>联系我们</h2>
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
  components: {
    ArrowLeft
  },
  setup() {
    const router = useRouter()

    const goBack = () => {
      router.push('/profile')
    }

    return {
      goBack
    }
  }
}
</script>

<style scoped>
.about-container {
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  min-height: calc(100vh - 60px);
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

.about-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
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

/* 移动端适配 */
@media (max-width: 768px) {
  .about-container {
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

  .tech-stack {
    grid-template-columns: 1fr;
  }

  .features {
    grid-template-columns: 1fr;
  }
}
</style>