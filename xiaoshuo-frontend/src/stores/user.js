// stores/user.js
import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null,
    isAuthenticated: false
  }),
  
  getters: {
    isAdmin: (state) => {
      return state.user && state.user.is_admin
    }
  },
  
  actions: {
    async login(email, password) {
      try {
        const response = await axios.post('/api/v1/users/login', {
          email,
          password
        })
        
        const { user, token } = response.data.data
        this.user = user
        this.token = token
        this.isAuthenticated = true
        
        // 保存token到本地存储
        localStorage.setItem('token', token)
        
        // 设置axios默认请求头
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        
        return { success: true, data: response.data }
      } catch (error) {
        return { success: false, error: error.response?.data || error.message }
      }
    },
    
    async register(email, password, nickname = '') {
      try {
        const response = await axios.post('/api/v1/users/register', {
          email,
          password,
          nickname
        })
        
        const { user, token } = response.data.data
        this.user = user
        this.token = token
        this.isAuthenticated = true
        
        // 保存token到本地存储
        localStorage.setItem('token', token)
        
        // 设置axios默认请求头
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        
        return { success: true, data: response.data }
      } catch (error) {
        return { success: false, error: error.response?.data || error.message }
      }
    },
    
    async fetchProfile() {
      if (!this.token) return
      
      try {
        const response = await axios.get('/api/v1/users/profile', {
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        })
        
        this.user = response.data.data
        this.isAuthenticated = true
        
        return { success: true, data: response.data }
      } catch (error) {
        // 如果获取用户信息失败，清除认证状态
        this.logout()
        return { success: false, error: error.response?.data || error.message }
      }
    },
    
    async updateProfile(nickname) {
      try {
        const response = await axios.put('/api/v1/users/profile', {
          nickname
        }, {
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        })
        
        this.user = response.data.data
        return { success: true, data: response.data }
      } catch (error) {
        return { success: false, error: error.response?.data || error.message }
      }
    },
    
    logout() {
      this.user = null
      this.token = null
      this.isAuthenticated = false
      
      // 清除本地存储的token
      localStorage.removeItem('token')
      
      // 清除axios默认请求头
      delete axios.defaults.headers.common['Authorization']
    }
  }
})