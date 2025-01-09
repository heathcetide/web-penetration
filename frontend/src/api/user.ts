import request from '@/utils/request'

export interface UserInfo {
  id: number
  username: string
  email: string
  role: string
  createdAt: string
  updatedAt: string
}

export interface UpdateUserData {
  id: number
  username?: string
  email?: string
  role?: string
  password?: string
}

// 获取用户列表
export const getUserList = (params: {
  page: number
  pageSize: number
  keyword?: string
}) => {
  return request.get('/users', { params })
}

// 获取单个用户信息
export const getUserInfo = (id: number) => {
  return request.get(`/users/${id}`)
}

// 更新用户信息
export const updateUser = (data: UpdateUserData) => {
  return request.put(`/users/${data.id}`, data)
}

// 删除用户
export const deleteUser = (id: number) => {
  return request.delete(`/users/${id}`)
}

// 修改密码
export const changePassword = (data: {
  oldPassword: string
  newPassword: string
}) => {
  return request.post('/users/change-password', data)
} 