import axios from 'axios'

const API_URL = 'http://localhost:8080/api'

export interface LoginData {
  username: string
  password: string
}

export interface RegisterData {
  username: string
  password: string
  email: string
}

export const login = (data: LoginData) => {
  return axios.post(`${API_URL}/login`, data)
}

export const register = (data: RegisterData) => {
  return axios.post(`${API_URL}/register`, data)
} 