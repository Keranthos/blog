import axios from 'axios'
import apiConfig from '@/config/api'
import { showErrorMessage } from '@/utils/waifuMessage'
import store from '@/store'

const apiUrl = apiConfig.apiURL

function getJWT () {
  return localStorage.getItem('jwt')
}

function isJWTExpired (token) {
  if (!token) return true
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return Date.now() >= payload.exp * 1000
  } catch (error) {
    return true // 如果解析失败，认为token无效
  }
}

// 获取JWT过期时间信息
function getJWTExpirationInfo (token) {
  if (!token) return null
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    const expTime = payload.exp * 1000 // 转换为毫秒
    const currentTime = Date.now()
    const timeLeft = expTime - currentTime

    return {
      exp: expTime,
      expDate: new Date(expTime),
      timeLeft,
      isExpired: timeLeft <= 0,
      hoursLeft: Math.floor(timeLeft / (1000 * 60 * 60)),
      minutesLeft: Math.floor((timeLeft % (1000 * 60 * 60)) / (1000 * 60))
    }
  } catch (error) {
    return null
  }
}

// 创建axios实例，添加统一的错误处理
const apiClient = axios.create({
  baseURL: apiUrl,
  timeout: 10000 // 10秒超时
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 检查JWT过期时间并显示信息
    const token = config.headers?.Authorization?.replace('Bearer ', '') || getJWT()
    if (token) {
      // JWT 过期时间检查（生产环境可移除）
      // const expInfo = getJWTExpirationInfo(token)
    }

    // 不再显示"正在请求数据"消息
    // if (window.showMessage) {
    //   window.showMessage('正在请求数据...', 1000, 5)
    // }
    return config
  },
  (error) => {
    showErrorMessage('network')
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    // 请求成功，只对明确标记需要看板娘提示的响应才显示消息
    if (window.showMessage && response.data.message && response.data.userToast === true) {
      window.showMessage(response.data.message, 3000, 8)
    }
    return response
  },
  (error) => {
    // 统一错误处理
    if (error.response) {
      const status = error.response.status
      const message = error.response.data?.error || error.response.data?.message

      // 处理401未授权错误（JWT过期）
      if (status === 401) {
        console.log('服务器返回401，JWT可能已过期，自动退出登录')
        showErrorMessage(401) // 使用看板娘显示登录过期信息

        // 自动退出登录
        store.dispatch('logout')

        // 延迟跳转到登录页面
        setTimeout(() => {
          if (window.location.pathname !== '/login-register') {
            window.location.href = '/login-register'
          }
        }, 2000)
      } else {
        // 使用看板娘显示其他错误信息
        if (message) {
          showErrorMessage(message)
        } else {
          showErrorMessage(status)
        }
      }
    } else if (error.request) {
      // 网络错误
      showErrorMessage('network')
    } else {
      // 其他错误
      showErrorMessage('unknown')
    }

    return Promise.reject(error)
  }
)

async function requestFunc (Url, object, isTokenNeeded) {
  if (isTokenNeeded) {
    // 优先使用传入的token，否则从localStorage获取
    const token = object.headers?.Authorization?.replace('Bearer ', '') || getJWT()

    if (isJWTExpired(token)) {
      console.log('JWT已过期，自动退出登录')
      showErrorMessage(401) // 使用看板娘显示登录过期信息

      // 自动退出登录
      store.dispatch('logout')

      // 延迟跳转到登录页面
      setTimeout(() => {
        if (window.location.pathname !== '/login-register') {
          window.location.href = '/login-register'
        }
      }, 2000)

      return null
    }

    const headers = { ...object.headers }
    // 对于FormData，不设置Content-Type，让浏览器自动设置
    // 对于其他数据，设置application/json
    if (!(object.data instanceof FormData)) {
      headers['Content-Type'] = 'application/json'
    }
    headers.Authorization = `Bearer ${token}`

    return await apiClient({
      url: Url,
      method: object.method,
      headers,
      data: object?.data ? (object.data instanceof FormData ? object.data : JSON.stringify(object.data)) : null
    })
  } else {
    const headers = { ...object.headers }
    // 对于FormData，不设置Content-Type，让浏览器自动设置
    // 对于其他数据，设置application/json
    if (!(object.data instanceof FormData)) {
      headers['Content-Type'] = 'application/json'
    }

    return await apiClient({
      url: Url,
      method: object.method,
      headers,
      data: object?.data ? (object.data instanceof FormData ? object.data : JSON.stringify(object.data)) : null
    })
  }
}

export { requestFunc, getJWTExpirationInfo, getJWT }
