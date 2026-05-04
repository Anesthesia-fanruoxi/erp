import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const request = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000,
})

// 请求拦截器 - 自动携带 Token
request.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

// 响应拦截器 - 统一错误处理
request.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code !== 200) {
      // 可以做统一提示
      if (res.code === 401) {
        // Token失效，跳转登录
        const authStore = useAuthStore()
        authStore.clearAuth()
        router.push('/login')
      }
      return Promise.reject(new Error(res.message || 'Error'))
    }
    return res
  },
  (error) => {
    return Promise.reject(error)
  }
)

export default request
