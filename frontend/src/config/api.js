// API 配置文件
// 统一管理 API 地址

// 基础 URL（从环境变量读取）
export const API_BASE_URL = process.env.VUE_APP_API_BASE_URL || 'http://localhost:3000'

// API 前缀
export const API_PREFIX = process.env.VUE_APP_API_PREFIX || '/api'

// 完整的 API URL
export const API_URL = `${API_BASE_URL}${API_PREFIX}`

// 导出配置对象
export default {
  baseURL: API_BASE_URL,
  apiPrefix: API_PREFIX,
  apiURL: API_URL
}
