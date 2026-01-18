import { defineStore } from 'pinia'
import apiClient from '@/utils/api'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || '',
    isAuthenticated: false,
  }),

  getters: {
    isAdmin: (state) => {
      return state.user && state.user.is_admin === true
    }
  },

  actions: {
    async login(email, password) {
      try {
        const response = await apiClient.post('/api/v1/users/login', {
          email,
          password
        })
        
        if (response.data.code === 200) {
          const { token, user } = response.data.data
          this.user = user
          this.token = token
          this.isAuthenticated = true
          
          // 保存token到localStorage
          localStorage.setItem('token', token)
          
          return {
            ...response.data,
            success: true
          }
        } else {
          return {
            ...response.data,
            success: false
          }
        }
      } catch (error) {
        console.error('Login error:', error)
        throw error
      }
    },

    async register(email, password, nickname) {
      try {
        const response = await apiClient.post('/api/v1/users/register', {
          email,
          password,
          nickname
        })
        
        if (response.data.code === 200) {
          // 检查是否返回了token（有些实现可能在注册后自动登录，有些则需要激活）
          const { token, user, message } = response.data.data
          
          // 如果有token，更新用户状态
          if (token) {
            this.user = user
            this.token = token
            this.isAuthenticated = true
            // 保存token到localStorage
            localStorage.setItem('token', token)
          } else {
            // 如果没有token，但用户信息存在，也更新用户信息（可能需要激活）
            if (user) {
              this.user = user
            }
          }
          
          // 返回响应数据，包括可能的额外信息
          return {
            ...response.data,
            success: true,
            message: message,
            user: user // 确保user信息也被返回
          }
        }
        
        return {
          ...response.data,
          success: false
        }
      } catch (error) {
        console.error('Register error:', error)
        throw error
      }
    },

    async fetchProfile() {
      try {
        const response = await apiClient.get('/api/v1/users/profile')
        
        if (response.data.code === 200) {
          this.user = response.data.data.user
          this.isAuthenticated = true
        }
        
        return response.data
      } catch (error) {
        console.error('Fetch profile error:', error)
        throw error
      }
    },

    async updateProfile(profileData) {
      try {
        const response = await apiClient.put('/api/v1/users/profile', profileData)
        
        if (response.data.code === 200) {
          this.user = response.data.data.user
        }
        
        return response.data
      } catch (error) {
        console.error('Update profile error:', error)
        throw error
      }
    },

    logout() {
      this.user = null
      this.token = ''
      this.isAuthenticated = false
      
      // 清除localStorage中的token
      localStorage.removeItem('token')
    }
  }
})