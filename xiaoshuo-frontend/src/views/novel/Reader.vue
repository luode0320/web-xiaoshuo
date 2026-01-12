<template>
  <div class="reader-container">
    <div v-if="loading" class="loading">
      <el-skeleton :rows="6" animated />
    </div>
    <div v-else class="reader-content">
      <!-- 阅读器头部 -->
      <div class="reader-header">
        <div class="header-left">
          <el-button @click="goBack" icon="ArrowLeft" plain>返回</el-button>
          <span class="novel-title">{{ novel?.title }}</span>
        </div>
        <div class="header-right">
          <el-button @click="showSettings = true" icon="Setting" plain>设置</el-button>
        </div>
      </div>
      
      <!-- 阅读区域 -->
      <div class="reading-area" :style="readerStyle">
        <div class="chapter-title">{{ currentChapter }}</div>
        <div class="content-text" v-html="processedContent"></div>
      </div>
      
      <!-- 阅读器底部 -->
      <div class="reader-footer">
        <div class="progress-info">
          <el-slider 
            v-model="currentProgress" 
            @change="onProgressChange"
            :max="totalChapters"
            :format-tooltip="formatChapter"
          />
          <div class="progress-text">
            第 {{ currentChapterIndex + 1 }} 章 / 共 {{ totalChapters }} 章
          </div>
        </div>
        <div class="navigation">
          <el-button @click="prevChapter" :disabled="currentChapterIndex <= 0" plain>
            上一章
          </el-button>
          <el-button @click="nextChapter" :disabled="currentChapterIndex >= totalChapters - 1" plain>
            下一章
          </el-button>
        </div>
      </div>
      
      <!-- 阅读设置对话框 -->
      <el-dialog v-model="showSettings" title="阅读设置" width="400px" destroy-on-close>
        <div class="settings-option">
          <label>字体大小</label>
          <el-slider 
            v-model="fontSize" 
            :min="14" 
            :max="24" 
            @change="saveSettings"
          />
        </div>
        <div class="settings-option">
          <label>背景颜色</label>
          <el-radio-group v-model="backgroundColor" @change="saveSettings">
            <el-radio label="white">白色</el-radio>
            <el-radio label="yellow">护眼</el-radio>
            <el-radio label="gray">灰色</el-radio>
            <el-radio label="black">黑色</el-radio>
          </el-radio-group>
        </div>
        <div class="settings-option">
          <label>字体类型</label>
          <el-radio-group v-model="fontFamily" @change="saveSettings">
            <el-radio label="serif">衬线字体</el-radio>
            <el-radio label="sans-serif">无衬线字体</el-radio>
            <el-radio label="cursive">手写体</el-radio>
          </el-radio-group>
        </div>
        <div class="settings-option">
          <label>行间距</label>
          <el-slider 
            v-model="lineHeight" 
            :min="1.2" 
            :max="2.0" 
            :step="0.1"
            @change="saveSettings"
          />
        </div>
        <template #footer>
          <el-button @click="showSettings = false">关闭</el-button>
          <el-button type="primary" @click="saveSettings(); showSettings = false;">保存</el-button>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export default {
  name: 'Reader',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const userStore = useUserStore()
    
    // 响应式数据
    const loading = ref(true)
    const novel = ref(null)
    const content = ref('')
    const showSettings = ref(false)
    
    // 阅读设置
    const settings = reactive({
      fontSize: 16,
      backgroundColor: 'white',
      fontFamily: 'serif',
      lineHeight: 1.6
    })
    
    // 章节信息
    const currentChapterIndex = ref(0)
    const chapters = ref([])
    
    // 从localStorage获取保存的设置
    const savedSettings = localStorage.getItem('readerSettings')
    if (savedSettings) {
      Object.assign(settings, JSON.parse(savedSettings))
    }
    
    // 计算属性
    const currentProgress = computed({
      get: () => currentChapterIndex.value + 1,
      set: (value) => {
        currentChapterIndex.value = value - 1
        updateChapterContent()
      }
    })
    
    const totalChapters = computed(() => chapters.value.length)
    
    const currentChapter = computed(() => {
      if (chapters.value[currentChapterIndex.value]) {
        return chapters.value[currentChapterIndex.value].title || `第${currentChapterIndex.value + 1}章`
      }
      return '第一章'
    })
    
    const processedContent = computed(() => {
      if (chapters.value[currentChapterIndex.value]) {
        return chapters.value[currentChapterIndex.value].content || ''
      }
      return content.value
    })
    
    const readerStyle = computed(() => {
      let bgColor = '#ffffff'
      switch (settings.backgroundColor) {
        case 'yellow':
          bgColor = '#f5f5dc'  // 浅黄色
          break
        case 'gray':
          bgColor = '#f0f0f0'  // 浅灰色
          break
        case 'black':
          bgColor = '#000000'  // 黑色
          break
      }
      
      return {
        fontSize: `${settings.fontSize}px`,
        backgroundColor: bgColor,
        fontFamily: settings.fontFamily,
        lineHeight: settings.lineHeight
      }
    })
    
    // 方法
    const loadNovel = async () => {
      try {
        const response = await axios.get(`/api/v1/novels/${route.params.id}`)
        novel.value = response.data.data
      } catch (error) {
        console.error('获取小说信息失败:', error)
        ElMessage.error('获取小说信息失败')
      }
    }
    
    const loadContent = async () => {
      try {
        const response = await axios.get(`/api/v1/novels/${route.params.id}/content`)
        content.value = response.data.data.content
        
        // 简单解析章节（实际应用中可能需要更复杂的解析）
        parseChapters(content.value)
      } catch (error) {
        console.error('获取小说内容失败:', error)
        ElMessage.error('获取小说内容失败')
      } finally {
        loading.value = false
      }
    }
    
    const parseChapters = (content) => {
      // 简单按换行分隔，实际应用中需要根据小说格式（如EPUB）进行解析
      const lines = content.split('\n')
      const chapterPattern = /^(第[一二三四五六七八九十零\d]+章|Chapter\s+\d+)/i
      
      let currentChapter = { title: '第一章', content: '', startLine: 0 }
      const parsedChapters = []
      
      for (let i = 0; i < lines.length; i++) {
        const line = lines[i].trim()
        if (chapterPattern.test(line) && line.length < 100) { // 简单判断是否为章节标题
          if (currentChapter.content) {
            parsedChapters.push({...currentChapter})
          }
          currentChapter = { title: line, content: '', startLine: i }
        } else {
          currentChapter.content += line + '\n'
        }
      }
      
      // 添加最后一章
      if (currentChapter.content) {
        parsedChapters.push({...currentChapter})
      }
      
      // 如果没有找到章节，将整个内容作为一章
      if (parsedChapters.length === 0) {
        parsedChapters.push({
          title: '正文',
          content: content,
          startLine: 0
        })
      }
      
      chapters.value = parsedChapters
    }
    
    const updateChapterContent = () => {
      // 在章节变化时更新阅读进度
      saveReadingProgress()
    }
    
    const saveReadingProgress = async () => {
      if (!userStore.isAuthenticated) return
      
      try {
        await axios.post(
          `/api/v1/novels/${route.params.id}/progress`, 
          {
            chapter_id: currentChapterIndex.value + 1,
            chapter_name: currentChapter.value,
            position: 0, // 实际应用中可以记录阅读位置
            progress: Math.round(((currentChapterIndex.value + 1) / totalChapters.value) * 100)
          },
          {
            headers: {
              'Authorization': `Bearer ${userStore.token}`
            }
          }
        )
      } catch (error) {
        console.error('保存阅读进度失败:', error)
      }
    }
    
    const loadReadingProgress = async () => {
      if (!userStore.isAuthenticated) return
      
      try {
        const response = await axios.get(`/api/v1/novels/${route.params.id}/progress`, {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        const progress = response.data.data
        if (progress.chapter_id) {
          currentChapterIndex.value = Math.max(0, progress.chapter_id - 1)
        }
      } catch (error) {
        // 如果获取阅读进度失败，使用默认第一章
        currentChapterIndex.value = 0
      }
    }
    
    const prevChapter = () => {
      if (currentChapterIndex.value > 0) {
        currentChapterIndex.value--
        updateChapterContent()
      }
    }
    
    const nextChapter = () => {
      if (currentChapterIndex.value < totalChapters.value - 1) {
        currentChapterIndex.value++
        updateChapterContent()
      }
    }
    
    const onProgressChange = (value) => {
      currentChapterIndex.value = value - 1
      updateChapterContent()
    }
    
    const formatChapter = (value) => {
      return `第${value}章`
    }
    
    const saveSettings = () => {
      localStorage.setItem('readerSettings', JSON.stringify(settings))
    }
    
    const goBack = () => {
      router.go(-1)
    }
    
    // 监听章节变化，保存阅读进度
    watch(currentChapterIndex, () => {
      saveReadingProgress()
    })
    
    onMounted(async () => {
      await loadNovel()
      await loadContent()
      await loadReadingProgress()
    })
    
    return {
      loading,
      novel,
      content,
      showSettings,
      settings,
      currentChapterIndex,
      chapters,
      currentProgress,
      totalChapters,
      currentChapter,
      processedContent,
      readerStyle,
      fontSize: computed({
        get: () => settings.fontSize,
        set: (value) => { settings.fontSize = value }
      }),
      backgroundColor: computed({
        get: () => settings.backgroundColor,
        set: (value) => { settings.backgroundColor = value }
      }),
      fontFamily: computed({
        get: () => settings.fontFamily,
        set: (value) => { settings.fontFamily = value }
      }),
      lineHeight: computed({
        get: () => settings.lineHeight,
        set: (value) => { settings.lineHeight = value }
      }),
      prevChapter,
      nextChapter,
      onProgressChange,
      formatChapter,
      saveSettings,
      goBack
    }
  }
}
</script>

<style scoped>
.reader-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
}

.loading {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.reader-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.header-left, .header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.novel-title {
  font-weight: 500;
  color: #333;
}

.reading-area {
  flex: 1;
  overflow-y: auto;
  padding: 40px 20px;
  text-align: justify;
  color: #333;
}

.chapter-title {
  text-align: center;
  font-size: 1.5rem;
  margin-bottom: 30px;
  color: #333;
  font-weight: bold;
}

.content-text {
  line-height: 1.6;
  font-size: 16px;
}

.reader-footer {
  padding: 20px;
  background: white;
  border-top: 1px solid #eee;
  z-index: 10;
}

.progress-info {
  margin-bottom: 15px;
}

.progress-text {
  text-align: center;
  font-size: 0.9rem;
  color: #666;
  margin-top: 10px;
}

.navigation {
  display: flex;
  justify-content: center;
  gap: 20px;
}

.settings-option {
  margin-bottom: 20px;
}

.settings-option label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

@media (max-width: 768px) {
  .reader-header {
    padding: 10px 15px;
  }
  
  .reading-area {
    padding: 20px 15px;
  }
  
  .chapter-title {
    font-size: 1.2rem;
  }
  
  .content-text {
    font-size: 16px;
  }
  
  .reader-footer {
    padding: 15px;
  }
  
  .navigation {
    flex-direction: column;
    gap: 10px;
  }
  
  .navigation .el-button {
    margin: 0;
  }
}
</style>