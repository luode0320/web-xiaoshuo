<template>
  <div class="login-container">
    <div class="login-form">
      <h2>用户登录</h2>
      <el-form 
        :model="loginForm" 
        :rules="loginRules" 
        ref="loginFormRef"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="email">
          <el-input 
            v-model="loginForm.email" 
            placeholder="邮箱"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input 
            v-model="loginForm.password" 
            type="password" 
            placeholder="密码"
            prefix-icon="Lock"
          />
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleLogin" 
            :loading="loading"
            style="width: 100%"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
      <div class="register-link">
        <p>还没有账号？ <router-link to="/register">立即注册</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const loading = ref(false)
    const loginFormRef = ref(null)
    
    const loginForm = reactive({
      email: '',
      password: ''
    })
    
    const loginRules = {
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ]
    }
    
    const handleLogin = async () => {
      if (!loginFormRef.value) return
      
      try {
        await loginFormRef.value.validate()
        
        loading.value = true
        
        const result = await userStore.login(loginForm.email, loginForm.password)
        
        if (result.success) {
          ElMessage.success('登录成功')
          // 跳转到首页或之前的页面
          const redirect = router.currentRoute.value.query.redirect || '/'
          router.push(redirect)
        } else {
          ElMessage.error(result.error?.message || '登录失败')
        }
      } catch (error) {
        console.error('登录错误:', error)
        ElMessage.error('登录失败')
      } finally {
        loading.value = false
      }
    }
    
    return {
      loginForm,
      loginRules,
      loginFormRef,
      loading,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-form {
  background: white;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.login-form h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
  font-size: 24px;
}

.register-link {
  text-align: center;
  margin-top: 20px;
}

.register-link a {
  color: #409eff;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}