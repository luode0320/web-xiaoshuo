import { defineStore } from 'pinia'
import apiClient from '@/utils/api'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || '',
    isAuthenticated: false,
  }),

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
        }
        
        return response.data
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
          const { token, user } = response.data.data
          this.user = user
          this.token = token
          this.isAuthenticated = true
          
          // 保存token到localStorage
          localStorage.setItem('token', token)
        }
        
        return response.data
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