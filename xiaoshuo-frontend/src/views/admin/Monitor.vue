<template>
  <div class="admin-monitor-container">
    <div class="admin-header">
      <h1>用户行为监控</h1>
      <p>监控用户活动，分析行为模式</p>
    </div>
    
    <div class="admin-content">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="用户活动日志" name="activity">
          <div class="activity-logs">
            <div class="filter-section">
              <el-form :model="activityFilters" inline>
                <el-form-item label="用户">
                  <el-select 
                    v-model="activityFilters.user_id" 
                    placeholder="选择用户" 
                    clearable
                    @change="fetchActivityLogs"
                  >
                    <el-option 
                      v-for="user in users" 
                      :key="user.id" 
                      :label="user.nickname || user.email" 
                      :value="user.id"
                    />
                  </el-select>
                </el-form-item>
                
                <el-form-item label="操作类型">
                  <el-select 
                    v-model="activityFilters.action_type" 
                    placeholder="选择操作类型" 
                    clearable
                    @change="fetchActivityLogs"
                  >
                    <el-option label="登录" value="login"></el-option>
                    <el-option label="上传小说" value="upload_novel"></el-option>
                    <el-option label="发布评论" value="create_comment"></el-option>
                    <el-option label="提交评分" value="create_rating"></el-option>
                    <el-option label="搜索" value="search"></el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item label="时间范围">
                  <el-date-picker
                    v-model="activityFilters.date_range"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    @change="fetchActivityLogs"
                  />
                </el-form-item>
                
                <el-form-item>
                  <el-button type="primary" @click="fetchActivityLogs">查询</el-button>
                  <el-button @click="resetActivityFilters">重置</el-button>
                </el-form-item>
              </el-form>
            </div>
            
            <el-table 
              :data="activityLogs" 
              style="width: 100%"
              v-loading="loading.activity"
            >
              <el-table-column prop="user.nickname" label="用户" width="150" />
              <el-table-column prop="action_type" label="操作类型" width="120" />
              <el-table-column prop="description" label="描述" show-overflow-tooltip />
              <el-table-column prop="ip_address" label="IP地址" width="130" />
              <el-table-column prop="created_at" label="时间" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
            </el-table>
            
            <!-- 分页 -->
            <div class="pagination" v-if="activityPagination.total > activityPagination.limit">
              <el-pagination
                v-model:current-page="activityPagination.page"
                :page-size="activityPagination.limit"
                :total="activityPagination.total"
                @current-change="handleActivityPageChange"
                layout="prev, pager, next, jumper"
              />
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="用户统计" name="statistics">
          <div class="user-statistics">
            <div class="stats-cards">
              <el-card class="stat-card">
                <div class="stat-content">
                  <div class="stat-title">总用户数</div>
                  <div class="stat-value">{{ userStats.total_users }}</div>
                  <div class="stat-change">
                    <span :class="userStats.new_users_today >= 0 ? 'increase' : 'decrease'">
                      {{ userStats.new_users_today >= 0 ? '+' : '' }}{{ userStats.new_users_today }} 今日新增
                    </span>
                  </div>
                </div>
              </el-card>
              
              <el-card class="stat-card">
                <div class="stat-content">
                  <div class="stat-title">活跃用户</div>
                  <div class="stat-value">{{ userStats.active_users_today }}</div>
                  <div class="stat-change">
                    <span :class="userStats.active_users_week >= 0 ? 'increase' : 'decrease'">
                      {{ userStats.active_users_week >= 0 ? '+' : '' }}{{ userStats.active_users_week }} 本周新增
                    </span>
                  </div>
                </div>
              </el-card>
              
              <el-card class="stat-card">
                <div class="stat-content">
                  <div class="stat-title">上传小说数</div>
                  <div class="stat-value">{{ userStats.total_novels }}</div>
                  <div class="stat-change">
                    <span :class="userStats.novels_today >= 0 ? 'increase' : 'decrease'">
                      {{ userStats.novels_today >= 0 ? '+' : '' }}{{ userStats.novels_today }} 今日上传
                    </span>
                  </div>
                </div>
              </el-card>
              
              <el-card class="stat-card">
                <div class="stat-content">
                  <div class="stat-title">评论数</div>
                  <div class="stat-value">{{ userStats.total_comments }}</div>
                  <div class="stat-change">
                    <span :class="userStats.comments_today >= 0 ? 'increase' : 'decrease'">
                      {{ userStats.comments_today >= 0 ? '+' : '' }}{{ userStats.comments_today }} 今日评论
                    </span>
                  </div>
                </div>
              </el-card>
            </div>
            
            <div class="chart-section">
              <el-card>
                <template #header>
                  <div class="card-header">
                    <span>用户活跃趋势</span>
                  </div>
                </template>
                <div ref="trendChartRef" style="height: 300px;"></div>
              </el-card>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="系统日志" name="logs">
          <div class="system-logs">
            <div class="filter-section">
              <el-form :model="logFilters" inline>
                <el-form-item label="日志级别">
                  <el-select 
                    v-model="logFilters.level" 
                    placeholder="选择日志级别" 
                    clearable
                    @change="fetchSystemLogs"
                  >
                    <el-option label="信息" value="info"></el-option>
                    <el-option label="警告" value="warning"></el-option>
                    <el-option label="错误" value="error"></el-option>
                  </el-select>
                </el-form-item>
                
                <el-form-item label="时间范围">
                  <el-date-picker
                    v-model="logFilters.date_range"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    @change="fetchSystemLogs"
                  />
                </el-form-item>
                
                <el-form-item>
                  <el-button type="primary" @click="fetchSystemLogs">查询</el-button>
                  <el-button @click="resetLogFilters">重置</el-button>
                </el-form-item>
              </el-form>
            </div>
            
            <el-table 
              :data="systemLogs" 
              style="width: 100%"
              v-loading="loading.system"
            >
              <el-table-column prop="level" label="级别" width="100">
                <template #default="{ row }">
                  <el-tag :type="getLogLevelType(row.level)">
                    {{ row.level.toUpperCase() }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="message" label="消息" show-overflow-tooltip />
              <el-table-column prop="created_at" label="时间" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
            </el-table>
            
            <!-- 分页 -->
            <div class="pagination" v-if="logPagination.total > logPagination.limit">
              <el-pagination
                v-model:current-page="logPagination.page"
                :page-size="logPagination.limit"
                :total="logPagination.total"
                @current-change="handleLogPageChange"
                layout="prev, pager, next, jumper"
              />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import apiClient from '@/utils/api'
import dayjs from 'dayjs'
import * as echarts from 'echarts'

export default {
  name: 'AdminMonitor',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const activeTab = ref('activity')
    const users = ref([])
    const activityLogs = ref([])
    const systemLogs = ref([])
    
    const loading = ref({
      activity: false,
      users: false,
      system: false
    })
    
    const userStats = ref({
      total_users: 0,
      new_users_today: 0,
      active_users_today: 0,
      active_users_week: 0,
      total_novels: 0,
      novels_today: 0,
      total_comments: 0,
      comments_today: 0
    })
    
    // 分页
    const activityPagination = ref({
      page: 1,
      limit: 10,
      total: 0
    })
    
    const logPagination = ref({
      page: 1,
      limit: 10,
      total: 0
    })
    
    // 过滤器
    const activityFilters = ref({
      user_id: null,
      action_type: null,
      date_range: []
    })
    
    const logFilters = ref({
      level: null,
      date_range: []
    })
    
    const trendChartRef = ref(null)
    let trendChart = null
    
    // 检查用户权限
    if (!userStore.isAuthenticated || !userStore.isAdmin) {
      router.push('/')
      ElMessage.error('您没有权限访问此页面')
      return
    }
    
    // 获取用户列表
    const fetchUsers = async () => {
      loading.value.users = true
      try {
        const response = await apiClient.get('/api/v1/admin/users', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        users.value = response.data.data
      } catch (error) {
        console.error('获取用户列表失败:', error)
        ElMessage.error('获取用户列表失败')
      } finally {
        loading.value.users = false
      }
    }
    
    // 获取活动日志
    const fetchActivityLogs = async () => {
      loading.value.activity = true
      try {
        const params = {
          page: activityPagination.value.page,
          limit: activityPagination.value.limit
        }
        
        if (activityFilters.value.user_id) {
          params.user_id = activityFilters.value.user_id
        }
        
        if (activityFilters.value.action_type) {
          params.action_type = activityFilters.value.action_type
        }
        
        if (activityFilters.value.date_range && activityFilters.value.date_range.length === 2) {
          params.start_date = activityFilters.value.date_range[0]
          params.end_date = activityFilters.value.date_range[1]
        }
        
        const response = await apiClient.get('/api/v1/admin/user-activities', {
          params,
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        activityLogs.value = response.data.data
        activityPagination.value.total = response.data.total || 0
      } catch (error) {
        console.error('获取活动日志失败:', error)
        ElMessage.error('获取活动日志失败')
      } finally {
        loading.value.activity = false
      }
    }
    
    // 获取系统日志
    const fetchSystemLogs = async () => {
      loading.value.system = true
      try {
        const params = {
          page: logPagination.value.page,
          limit: logPagination.value.limit
        }
        
        if (logFilters.value.level) {
          params.level = logFilters.value.level
        }
        
        if (logFilters.value.date_range && logFilters.value.date_range.length === 2) {
          params.start_date = logFilters.value.date_range[0]
          params.end_date = logFilters.value.date_range[1]
        }
        
        const response = await apiClient.get('/api/v1/admin/system-logs', {
          params,
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        systemLogs.value = response.data.data
        logPagination.value.total = response.data.total || 0
      } catch (error) {
        console.error('获取系统日志失败:', error)
        ElMessage.error('获取系统日志失败')
      } finally {
        loading.value.system = false
      }
    }
    
    // 获取用户统计
    const fetchUserStats = async () => {
      try {
        const response = await apiClient.get('/api/v1/admin/user-statistics', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        userStats.value = response.data.data
      } catch (error) {
        console.error('获取用户统计失败:', error)
        ElMessage.error('获取用户统计失败')
      }
    }
    
    // 获取趋势数据并绘制图表
    const fetchTrendData = async () => {
      try {
        const response = await apiClient.get('/api/v1/admin/user-trend', {
          headers: {
            'Authorization': `Bearer ${userStore.token}`
          }
        })
        
        const data = response.data.data
        nextTick(() => {
          initTrendChart(data)
        })
      } catch (error) {
        console.error('获取趋势数据失败:', error)
        ElMessage.error('获取趋势数据失败')
      }
    }
    
    // 初始化趋势图表
    const initTrendChart = (data) => {
      if (!trendChartRef.value) return
      
      trendChart = echarts.init(trendChartRef.value)
      
      const option = {
        tooltip: {
          trigger: 'axis'
        },
        legend: {
          data: ['新增用户', '活跃用户', '上传小说']
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: data.dates
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '新增用户',
            type: 'line',
            stack: 'Total',
            data: data.new_users
          },
          {
            name: '活跃用户',
            type: 'line',
            stack: 'Total',
            data: data.active_users
          },
          {
            name: '上传小说',
            type: 'line',
            stack: 'Total',
            data: data.uploaded_novels
          }
        ]
      }
      
      trendChart.setOption(option)
    }
    
    // 分页处理
    const handleActivityPageChange = (page) => {
      activityPagination.value.page = page
      fetchActivityLogs()
    }
    
    const handleLogPageChange = (page) => {
      logPagination.value.page = page
      fetchSystemLogs()
    }
    
    // 重置过滤器
    const resetActivityFilters = () => {
      activityFilters.value = {
        user_id: null,
        action_type: null,
        date_range: []
      }
      activityPagination.value.page = 1
      fetchActivityLogs()
    }
    
    const resetLogFilters = () => {
      logFilters.value = {
        level: null,
        date_range: []
      }
      logPagination.value.page = 1
      fetchSystemLogs()
    }
    
    // 格式化日期
    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
    
    // 获取日志级别类型
    const getLogLevelType = (level) => {
      switch (level) {
        case 'error': return 'danger'
        case 'warning': return 'warning'
        case 'info': return 'primary'
        default: return 'primary'
      }
    }
    
    onMounted(() => {
      fetchUsers()
      fetchActivityLogs()
      fetchSystemLogs()
      fetchUserStats()
      fetchTrendData()
    })
    
    // 组件卸载时销毁图表
    onUnmounted(() => {
      if (trendChart) {
        trendChart.dispose()
      }
    })
    
    return {
      activeTab,
      users,
      activityLogs,
      systemLogs,
      userStats,
      loading,
      activityPagination,
      logPagination,
      activityFilters,
      logFilters,
      trendChartRef,
      fetchActivityLogs,
      fetchSystemLogs,
      handleActivityPageChange,
      handleLogPageChange,
      resetActivityFilters,
      resetLogFilters,
      formatDate,
      getLogLevelType
    }
  }
}
</script>

<style scoped>
.admin-monitor-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.admin-header {
  text-align: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.admin-header h1 {
  font-size: 1.8rem;
  color: #333;
  margin-bottom: 10px;
}

.admin-header p {
  color: #666;
}

.admin-content {
  min-height: 500px;
}

.filter-section {
  margin-bottom: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.stat-content {
  text-align: center;
  padding: 20px;
}

.stat-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin-bottom: 10px;
}

.stat-change {
  font-size: 12px;
  color: #999;
}

.stat-change .increase {
  color: #67c23a;
}

.stat-change .decrease {
  color: #f56c6c;
}

.chart-section {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  text-align: center;
  margin-top: 20px;
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

@media (max-width: 768px) {
  .admin-monitor-container {
    padding: 15px;
  }
  
  .admin-header h1 {
    font-size: 1.5rem;
  }
  
  .stats-cards {
    grid-template-columns: 1fr;
  }
  
  :deep(.el-table) {
    font-size: 14px;
  }
  
  :deep(.el-table .el-table__cell) {
    padding: 8px 0;
  }
  
  .filter-section {
    :deep(.el-form-item) {
      margin-bottom: 10px;
    }
  }
}