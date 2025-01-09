<template>
  <div class="profile-container">
    <a-card class="profile-card">
      <template #title>
        <h3>个人信息</h3>
      </template>
      
      <a-descriptions :data="profileData" layout="vertical" bordered />
      
      <a-divider />
      
      <h4>修改密码</h4>
      <a-form
        :model="passwordForm"
        @submit="handlePasswordChange"
        ref="passwordFormRef"
      >
        <a-form-item
          field="oldPassword"
          label="当前密码"
          :rules="[{ required: true, message: '请输入当前密码' }]"
        >
          <a-input-password v-model="passwordForm.oldPassword" />
        </a-form-item>
        
        <a-form-item
          field="newPassword"
          label="新密码"
          :rules="[
            { required: true, message: '请输入新密码' },
            { minLength: 6, message: '密码长度不能小于6位' }
          ]"
        >
          <a-input-password v-model="passwordForm.newPassword" />
        </a-form-item>
        
        <a-form-item
          field="confirmPassword"
          label="确认新密码"
          :rules="[
            { required: true, message: '请确认新密码' },
            { validator: validateConfirmPassword }
          ]"
        >
          <a-input-password v-model="passwordForm.confirmPassword" />
        </a-form-item>
        
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">
            修改密码
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive, computed } from 'vue'
import { Message } from '@arco-design/web-vue'
import { changePassword } from '@/api/user'
import { formatDate } from '@/utils/date'

export default defineComponent({
  name: 'Profile',
  setup() {
    const loading = ref(false)
    const passwordFormRef = ref()
    
    // 从localStorage获取用户信息
    const userInfo = JSON.parse(localStorage.getItem('user') || '{}')
    
    const profileData = computed(() => [
      {
        label: '用户名',
        value: userInfo.username
      },
      {
        label: '邮箱',
        value: userInfo.email
      },
      {
        label: '角色',
        value: userInfo.role
      },
      {
        label: '注册时间',
        value: formatDate(userInfo.createdAt)
      }
    ])
    
    const passwordForm = reactive({
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    
    const validateConfirmPassword = (value: string) => {
      if (value !== passwordForm.newPassword) {
        return '两次输入的密码不一致'
      }
      return true
    }
    
    const handlePasswordChange = async () => {
      try {
        await passwordFormRef.value?.validate()
        loading.value = true
        
        await changePassword({
          oldPassword: passwordForm.oldPassword,
          newPassword: passwordForm.newPassword
        })
        
        Message.success('密码修改成功')
        passwordFormRef.value?.resetFields()
      } catch (error) {
        Message.error('密码修改失败：' + error.response?.data?.error || '未知错误')
      } finally {
        loading.value = false
      }
    }
    
    return {
      loading,
      profileData,
      passwordForm,
      passwordFormRef,
      validateConfirmPassword,
      handlePasswordChange
    }
  }
})
</script>

<style scoped>
.profile-container {
  padding: 20px;
  display: flex;
  justify-content: center;
}

.profile-card {
  width: 600px;
}
</style> 