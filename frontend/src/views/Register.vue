<template>
  <div class="register-container">
    <a-card class="register-card" title="注册">
      <a-form :model="form" @submit="handleSubmit">
        <a-form-item field="username" label="用户名" 
          :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model="form.username" placeholder="请输入用户名" />
        </a-form-item>
        
        <a-form-item field="email" label="邮箱" 
          :rules="[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效的邮箱地址' }
          ]">
          <a-input v-model="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        
        <a-form-item field="password" label="密码" 
          :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model="form.password" placeholder="请输入密码" />
        </a-form-item>
        
        <a-form-item field="confirmPassword" label="确认密码" 
          :rules="[
            { required: true, message: '请确认密码' },
            { validator: validateConfirmPassword }
          ]">
          <a-input-password v-model="form.confirmPassword" placeholder="请确认密码" />
        </a-form-item>
        
        <a-form-item>
          <a-button type="primary" html-type="submit" long :loading="loading">
            注册
          </a-button>
        </a-form-item>
        <a-link href="/login">已有账号？立即登录</a-link>
      </a-form>
    </a-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'
import { register } from '@/api/auth'

export default defineComponent({
  name: 'Register',
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const form = reactive({
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    })

    const validateConfirmPassword = (value: string) => {
      if (value !== form.password) {
        return '两次输入的密码不一致'
      }
      return true
    }

    const handleSubmit = async () => {
      try {
        loading.value = true
        await register({
          username: form.username,
          email: form.email,
          password: form.password
        })
        Message.success('注册成功')
        router.push('/login')
      } catch (error) {
        Message.error('注册失败：' + error.response?.data?.error || '未知错误')
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      loading,
      handleSubmit,
      validateConfirmPassword
    }
  }
})
</script>

<style scoped>
.register-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
}

.register-card {
  width: 400px;
}
</style> 