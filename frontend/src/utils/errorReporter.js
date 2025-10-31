/**
 * 前端错误上报工具
 * 捕获并上报运行时错误、未处理的 Promise 拒绝等
 */

// 错误上报配置（可以后续接入 Sentry、LogRocket 等服务）
const ERROR_REPORT_CONFIG = {
  enabled: true, // 是否启用错误上报
  endpoint: null, // 错误上报 API 端点（可选，暂时只记录到 console）
  logToConsole: true, // 开发环境下记录到控制台
  collectUserInfo: false // 是否收集用户信息（隐私考虑，默认关闭）
}

/**
 * 格式化错误信息
 */
function formatError (error, context = {}) {
  const errorInfo = {
    timestamp: new Date().toISOString(),
    url: window.location.href,
    userAgent: navigator.userAgent,
    message: error?.message || 'Unknown error',
    stack: error?.stack || '',
    type: error?.name || 'Error',
    context: {
      route: window.location.pathname,
      ...context
    }
  }

  // 可选：收集用户信息（需要用户同意）
  if (ERROR_REPORT_CONFIG.collectUserInfo) {
    try {
      const store = require('@/store').default
      if (store && store.state && store.state.user) {
        errorInfo.user = {
          id: store.state.user.id,
          name: store.state.user.name,
          level: store.state.user.level
        }
      }
    } catch (e) {
      // 忽略收集用户信息失败
    }
  }

  return errorInfo
}

/**
 * 上报错误到后端（可选）
 */
async function reportError (errorInfo) {
  if (!ERROR_REPORT_CONFIG.enabled) {
    return
  }

  // 开发环境：记录到控制台
  if (ERROR_REPORT_CONFIG.logToConsole && process.env.NODE_ENV === 'development') {
    // eslint-disable-next-line no-console
    console.group('🚨 前端错误捕获')
    // eslint-disable-next-line no-console
    console.error('错误信息:', errorInfo)
    console.groupEnd()
  }

  // 生产环境：上报到后端（如果配置了 endpoint）
  if (ERROR_REPORT_CONFIG.endpoint && process.env.NODE_ENV === 'production') {
    try {
      const axios = require('axios').default
      await axios.post(ERROR_REPORT_CONFIG.endpoint, errorInfo, {
        timeout: 5000
      })
    } catch (e) {
      // 上报失败不阻塞用户操作
      // eslint-disable-next-line no-console
      console.error('错误上报失败:', e)
    }
  }
}

/**
 * 处理 JavaScript 运行时错误
 */
function handleError (event) {
  const error = event.error || new Error(event.message)
  const errorInfo = formatError(error, {
    filename: event.filename,
    lineno: event.lineno,
    colno: event.colno,
    errorType: 'runtime'
  })
  reportError(errorInfo)
}

/**
 * 处理未捕获的 Promise 拒绝
 */
function handleUnhandledRejection (event) {
  const error = event.reason instanceof Error
    ? event.reason
    : new Error(String(event.reason))
  const errorInfo = formatError(error, {
    errorType: 'unhandledRejection',
    reason: event.reason
  })
  reportError(errorInfo)
}

/**
 * 处理 Vue 组件错误（通过 Vue 全局错误处理）
 */
export function handleVueError (error, instance, info) {
  const errorInfo = formatError(error, {
    errorType: 'vue',
    componentName: instance?.$options?.name || instance?.$options?.__name || 'Unknown',
    componentFile: instance?.$options?.__file || 'Unknown',
    errorInfo: info
  })
  reportError(errorInfo)
}

/**
 * 初始化错误监控
 */
export function initErrorReporter () {
  if (!ERROR_REPORT_CONFIG.enabled) {
    return
  }

  // 捕获全局 JavaScript 错误
  window.addEventListener('error', handleError, true)

  // 捕获未处理的 Promise 拒绝
  window.addEventListener('unhandledrejection', handleUnhandledRejection, true)

  // 可选：捕获资源加载错误（图片、脚本等）
  window.addEventListener('error', (event) => {
    if (event.target !== window) {
      const errorInfo = formatError(
        new Error(`资源加载失败: ${event.target.tagName} - ${event.target.src || event.target.href}`),
        {
          errorType: 'resource',
          tagName: event.target.tagName,
          src: event.target.src || event.target.href
        }
      )
      reportError(errorInfo)
    }
  }, true)

  // 前端错误监控已初始化（生产环境不输出日志）
  if (process.env.NODE_ENV === 'development') {
    // eslint-disable-next-line no-console
    console.log('✅ 前端错误监控已初始化')
  }
}

/**
 * 手动上报错误（用于业务代码中的 try-catch）
 */
export function reportManualError (error, context = {}) {
  const errorInfo = formatError(error, {
    errorType: 'manual',
    ...context
  })
  reportError(errorInfo)
}
