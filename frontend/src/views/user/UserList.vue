<template>
  <div class="user-list">
    <div class="header">
      <h2>用户管理</h2>
      <div class="actions">
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索用户名或邮箱"
          style="width: 300px"
          @search="handleSearch"
        />
      </div>
    </div>

    <a-table
      :data="users"
      :loading="loading"
      :pagination="pagination"
      @page-change="onPageChange"
    >
      <template #columns>
        <a-table-column title="ID" data-index="id" />
        <a-table-column title="用户名" data-index="username" />
        <a-table-column title="邮箱" data-index="email" />
        <a-table-column title="角色" data-index="role">
          <template #cell="{ record }">
            <a-tag :color="record.role === 'admin' ? 'red' : 'blue'">
              {{ record.role }}
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="创建时间" data-index="createdAt">
          <template #cell="{ record }">
            {{ formatDate(record.createdAt) }}
          </template>
        </a-table-column>
        <a-table-column title="操作">
          <template #cell="{ record }">
            <a-space>
              <a-button type="text" @click="handleEdit(record)">
                编辑
              </a-button>
              <a-popconfirm
                content="确定要删除该用户吗？"
                @ok="handleDelete(record.id)"
              >
                <a-button type="text" status="danger">
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>

    <!-- 编辑用户对话框 -->
    <a-modal
      v-model:visible="editModalVisible"
      title="编辑用户"
      @ok="handleEditSubmit"
      @cancel="editModalVisible = false"
    >
      <a-form :model="editForm" ref="editFormRef">
        <a-form-item field="username" label="用户名" 
          :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model="editForm.username" />
        </a-form-item>
        <a-form-item field="email" label="邮箱"
          :rules="[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效的邮箱地址' }
          ]">
          <a-input v-model="editForm.email" />
        </a-form-item>
        <a-form-item field="role" label="角色">
          <a-select v-model="editForm.role">
            <a-option value="user">普通用户</a-option>
            <a-option value="admin">管理员</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="editForm.password" placeholder="不修改请留空" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import { getUserList, updateUser, deleteUser, UserInfo } from '@/api/user'
import { formatDate } from '@/utils/date'

export default defineComponent({
  name: 'UserList',
  setup() {
    const loading = ref(false)
    const users = ref<UserInfo[]>([])
    const searchKeyword = ref('')
    const pagination = reactive({
      total: 0,
      current: 1,
      pageSize: 10
    })

    // 编辑相关
    const editModalVisible = ref(false)
    const editFormRef = ref()
    const editForm = reactive({
      id: 0,
      username: '',
      email: '',
      role: '',
      password: ''
    })

    // 获取用户列表
    const fetchUsers = async () => {
      try {
        loading.value = true
        const res = await getUserList({
          page: pagination.current,
          pageSize: pagination.pageSize,
          keyword: searchKeyword.value
        })
        users.value = res.data.items
        pagination.total = res.data.total
      } catch (error) {
        Message.error('获取用户列表失败')
      } finally {
        loading.value = false
      }
    }

    // 搜索
    const handleSearch = () => {
      pagination.current = 1
      fetchUsers()
    }

    // 分页
    const onPageChange = (page: number) => {
      pagination.current = page
      fetchUsers()
    }

    // 编辑用户
    const handleEdit = (user: UserInfo) => {
      Object.assign(editForm, {
        id: user.id,
        username: user.username,
        email: user.email,
        role: user.role,
        password: ''
      })
      editModalVisible.value = true
    }

    // 提交编辑
    const handleEditSubmit = async () => {
      try {
        await editFormRef.value?.validate()
        const updateData = { ...editForm }
        if (!updateData.password) {
          delete updateData.password
        }
        await updateUser(updateData)
        Message.success('更新成功')
        editModalVisible.value = false
        fetchUsers()
      } catch (error) {
        Message.error('更新失败：' + error.response?.data?.error || '未知错误')
      }
    }

    // 删除用户
    const handleDelete = async (id: number) => {
      try {
        await deleteUser(id)
        Message.success('删除成功')
        fetchUsers()
      } catch (error) {
        Message.error('删除失败：' + error.response?.data?.error || '未知错误')
      }
    }

    onMounted(() => {
      fetchUsers()
    })

    return {
      loading,
      users,
      searchKeyword,
      pagination,
      editModalVisible,
      editForm,
      editFormRef,
      formatDate,
      handleSearch,
      onPageChange,
      handleEdit,
      handleEditSubmit,
      handleDelete
    }
  }
})
</script>

<style scoped>
.user-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.actions {
  display: flex;
  gap: 16px;
}
</style> 