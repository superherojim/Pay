import axios from 'axios'
import router from '../router'
import { ElMessage } from 'element-plus'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8100/v1',
    timeout: 10000
})

// 请求拦截器
api.interceptors.request.use(config => {
    const token = localStorage.getItem('token')
    if (config.url === '/login') return config

    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    } else {
        // 可以在此处跳转登录页或进行其他处理
    }
    return config
}, error => {
    console.error('请求拦截器错误:', error)
    return Promise.reject(error)
})

// 响应拦截器
api.interceptors.response.use(
    response => response.data,
    error => {
        if (!error.response) {
            ElMessage.error('网络连接异常，请检查网络状态')
            return Promise.reject(error)
        }

        const { status, data } = error.response
        const errorMessage = data?.message || '请求异常'

        // 自动处理认证失败
        if ([401, 403].includes(status)) {
            localStorage.removeItem('token')
            router.push('/login')
            ElMessage.warning('登录状态已过期，请重新登录')
            return Promise.reject(error)
        }

        // 统一错误提示
        if (status >= 500) {
            ElMessage.error(errorMessage)
        } else if (status >= 400) {
            ElMessage.warning(errorMessage)
        }

        // 返回格式化错误对象
        return Promise.reject({
            code: status,
            message: errorMessage,
            data: data
        })
    }
)

export default api

// 在现有API配置后添加统计接口
export const statsAPI = {
    getDashboardStats: () => api.get('/dashboard/stats')
} 