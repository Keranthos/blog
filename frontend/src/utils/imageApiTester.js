/**
 * 图片API测试工具
 * 用于测试各个图片API的稳定性和响应速度
 */

// 图片API配置
export const imageApiConfig = {
  // 风景类API
  scenery: [
    {
      name: '随机风景',
      url: 'https://api.btstu.cn/sjbz/api.php?lx=fengjing',
      type: 'direct',
      reliable: true,
      category: 'scenery'
    },
    {
      name: '随机图片',
      url: 'https://api.btstu.cn/sjbz/api.php?lx=suiji',
      type: 'direct',
      reliable: true,
      category: 'scenery'
    },
    {
      name: '随机壁纸',
      url: 'https://api.wetab.link/api/wallpaper',
      type: 'direct',
      reliable: true,
      category: 'scenery'
    }
  ],
  // 动漫类API
  anime: [
    {
      name: '随机动漫',
      url: 'https://api.btstu.cn/sjbz/api.php?lx=dongman',
      type: 'direct',
      reliable: true,
      category: 'anime'
    },
    {
      name: '萌图',
      url: 'https://api.yimian.xyz/img?type=moe',
      type: 'direct',
      reliable: true,
      category: 'anime'
    },
    {
      name: '动漫图片',
      url: 'https://api.gmit.vip/Api/DmImg?format=image',
      type: 'direct',
      reliable: true,
      category: 'anime'
    },
    {
      name: 'Lolimi动漫',
      url: 'https://api.lolimi.cn/API/acg/',
      type: 'direct',
      reliable: true,
      category: 'anime'
    }
  ],
  // 综合类API
  mixed: [
    {
      name: '综合随机',
      url: 'https://api.btstu.cn/sjbz/api.php',
      type: 'direct',
      reliable: true,
      category: 'mixed'
    },
    {
      name: '小歪API',
      url: 'https://api.ixiaowai.cn/api/api.php',
      type: 'json',
      reliable: true,
      category: 'mixed'
    },
    {
      name: '墨天逸',
      url: 'https://api.mtyqx.cn/api/random.php',
      type: 'direct',
      reliable: true,
      category: 'mixed'
    },
    {
      name: '国际源随机',
      url: 'https://picsum.photos/800/600?random=' + Date.now(),
      type: 'direct',
      reliable: true,
      category: 'mixed'
    }
  ]
}

// 获取所有API
export const getAllApis = () => {
  return [
    ...imageApiConfig.scenery,
    ...imageApiConfig.anime,
    ...imageApiConfig.mixed
  ]
}

// 获取推荐API（可靠的API）
export const getRecommendedApis = () => {
  return getAllApis().filter(api => api.reliable)
}

// 测试单个API
export const testSingleApi = async (api, timeout = 10000) => {
  const startTime = Date.now()
  const result = {
    name: api.name,
    url: api.url,
    type: api.type,
    category: api.category,
    reliable: api.reliable,
    status: 'pending',
    responseTime: null,
    imageUrl: null,
    imageSize: null,
    errorMessage: null,
    testTime: new Date().toISOString()
  }

  try {
    let imageUrl = ''

    // 根据API类型获取图片URL
    if (api.type === 'direct') {
      imageUrl = api.url
    } else if (api.type === 'json') {
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), timeout)

      const response = await fetch(api.url, {
        method: 'GET',
        headers: {
          Accept: 'application/json'
        },
        signal: controller.signal
      })

      clearTimeout(timeoutId)

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      imageUrl = data.picurl || data.url || data.imgurl || data.data

      if (!imageUrl) {
        throw new Error('API返回的数据中没有找到图片URL')
      }
    }

    // 验证图片URL
    if (!imageUrl) {
      throw new Error('获取到的图片URL为空')
    }

    // 预加载图片
    await preloadImage(imageUrl, timeout)

    const responseTime = Date.now() - startTime
    result.status = responseTime > 3000 ? 'slow' : 'success'
    result.responseTime = responseTime
    result.imageUrl = imageUrl

    return result
  } catch (error) {
    const responseTime = Date.now() - startTime
    result.status = 'error'
    result.responseTime = responseTime
    result.errorMessage = error.message

    return result
  }
}

// 预加载图片
const preloadImage = (url, timeout = 10000) => {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'

    const timeoutId = setTimeout(() => {
      reject(new Error(`图片加载超时: ${url}`))
    }, timeout)

    img.onload = () => {
      clearTimeout(timeoutId)
      resolve(url)
    }

    img.onerror = () => {
      clearTimeout(timeoutId)
      reject(new Error(`图片加载失败: ${url}`))
    }

    img.src = url
  })
}

// 批量测试API
export const testAllApis = async (options = {}) => {
  const {
    timeout = 10000,
    autoNext = true,
    delay = 1000,
    onProgress = null,
    onComplete = null
  } = options

  const apis = getAllApis()
  const results = []
  let completed = 0

  for (let i = 0; i < apis.length; i++) {
    const api = apis[i]
    console.log(`测试 ${i + 1}/${apis.length}: ${api.name}`)

    const result = await testSingleApi(api, timeout)
    results.push(result)
    completed++

    // 进度回调
    if (onProgress) {
      onProgress({
        completed,
        total: apis.length,
        current: api,
        result
      })
    }

    // 自动下一个
    if (autoNext && i < apis.length - 1) {
      await new Promise(resolve => setTimeout(resolve, delay))
    }
  }

  // 完成回调
  if (onComplete) {
    onComplete(results)
  }

  return results
}

// 生成测试报告
export const generateTestReport = (results) => {
  const total = results.length
  const success = results.filter(r => r.status === 'success').length
  const slow = results.filter(r => r.status === 'slow').length
  const error = results.filter(r => r.status === 'error').length

  const successfulResults = results.filter(r => r.status === 'success' || r.status === 'slow')
  const averageResponseTime = successfulResults.length > 0
    ? Math.round(successfulResults.reduce((sum, r) => sum + r.responseTime, 0) / successfulResults.length)
    : 0

  const recommendedApis = results
    .filter(r => r.status === 'success' && r.responseTime < 2000)
    .sort((a, b) => a.responseTime - b.responseTime)

  const problemApis = results.filter(r => r.status === 'error')

  return {
    summary: {
      total,
      success,
      slow,
      error,
      successRate: Math.round((success / total) * 100),
      averageResponseTime
    },
    recommendedApis,
    problemApis,
    results
  }
}

// 导出测试结果到JSON
export const exportTestResults = (results, filename = 'image-api-test-results.json') => {
  const report = generateTestReport(results)
  const dataStr = JSON.stringify(report, null, 2)
  const dataBlob = new Blob([dataStr], { type: 'application/json' })

  const link = document.createElement('a')
  link.href = URL.createObjectURL(dataBlob)
  link.download = filename
  link.click()
}

// 快速测试（只测试推荐API）
export const quickTest = async (options = {}) => {
  const recommendedApis = getRecommendedApis()
  const results = []

  for (const api of recommendedApis) {
    const result = await testSingleApi(api, options.timeout || 5000)
    results.push(result)

    if (options.delay) {
      await new Promise(resolve => setTimeout(resolve, options.delay))
    }
  }

  return results
}
