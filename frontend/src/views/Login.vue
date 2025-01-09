<template>
  <div class="login-container">
    <a-card class="login-card" title="登录">
      <a-form :model="form" @submit="handleSubmit">
        <a-form-item field="username" label="用户名" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model="form.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item field="password" label="密码" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model="form.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" long :loading="loading">
            登录
          </a-button>
        </a-form-item>
        <a-link href="/register">还没有账号？立即注册</a-link>
      </a-form>
    </a-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'
import { login } from '@/api/auth'

export default defineComponent({
  name: 'Login',
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const form = reactive({
      username: '',
      password: ''
    })

    const handleSubmit = async () => {
      try {
        loading.value = true
        const response = await login(form)
        const { token, user } = response.data
        
        // 存储token和用户信息
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
        
        Message.success('登录成功')
        router.push('/dashboard')
      } catch (error) {
        Message.error('登录失败：' + error.response?.data?.error || '未知错误')
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      loading,
      handleSubmit
    }
  }
})
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
}

.login-card {
  width: 400px;
}
</style> 