import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

export const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // Mengirim cookie HTTP-only otomatis
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface User {
  id: string
  name: string
  email: string
  username: string
  role: 'super_admin' | 'admin' | 'teacher' | 'student' | 'curriculum' | 'principal'
  isActive: boolean
  schoolID?: string
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  meta?: {
    page?: number
    limit?: number
    total?: number
  }
  error?: {
    code: string
    message: string
  }
}
