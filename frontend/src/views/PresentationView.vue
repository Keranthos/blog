<template>
  <div class="presentation-view">
    <!-- 错误边界 -->
    <div v-if="hasError" class="error-boundary">
      <h2>出现错误</h2>
      <p>{{ errorMessage }}</p>
      <button class="retry-btn" @click="retry">重试</button>
    </div>

    <!-- 主要内容 -->
    <div v-else>
      <NavBar />

      <!-- 英雄图片区域 -->
      <div class="hero-section">
        <div class="hero-image-wrapper">
          <img :src="presentationBackground" alt="Presentation Hero" class="hero-image-full" />
        </div>
      </div>

      <!-- 分隔文字 -->
      <div class="separator-section">
        <p class="separator-text">表达是人的本能，要如何满足和抑制它呢？</p>
      </div>

      <!-- 内容卡片区域 -->
      <div class="content-section">
        <div class="content-container">
          <!-- 加载状态 -->
          <div v-if="loading" class="loading-container">
            <div class="loading-spinner"></div>
            <p>加载中...</p>
          </div>

          <!-- 卡片网格 -->
          <div v-else class="cards-grid">
            <div
              v-for="presentation in currentPagePresentations"
              :key="presentation.id"
              class="presentation-card"
              @click="openPresentation(presentation)"
            >
              <div class="card-image">
                <img
                  :src="getThumbnailUrl(presentation.thumbnail)"
                  :alt="presentation.title"
                  @error="handleImageError"
                  @load="handleImageLoad"
                />
              </div>
              <div class="card-content">
                <h3 class="card-title">
                  <font-awesome-icon icon="file-powerpoint" class="title-icon" />
                  {{ presentation.title }}
                </h3>
                <div class="card-tags">
                  <span
                    v-for="tag in (presentation.tags && presentation.tags.length > 0 ? presentation.tags : ['讲演', '演示'])"
                    :key="tag"
                    class="tag"
                  >
                    {{ tag }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页控件 -->
          <div v-if="totalPages > 1" class="pagination">
            <button
              class="pagination-btn"
              :disabled="currentPage === 1"
              @click="goToPage(currentPage - 1)"
            >
              <font-awesome-icon icon="chevron-left" />
            </button>

            <div class="pagination-pages">
              <button
                v-for="page in visiblePages"
                :key="page"
                class="page-btn"
                :class="{ active: page === currentPage }"
                @click="goToPage(page)"
              >
                {{ page }}
              </button>
            </div>

            <button
              class="pagination-btn"
              :disabled="currentPage === totalPages"
              @click="goToPage(currentPage + 1)"
            >
              <font-awesome-icon icon="chevron-right" />
            </button>
          </div>
        </div>
      </div>

      <!-- PDF预览模态框 -->
      <div v-if="showPreview" class="preview-modal" @click="closePreview">
        <div class="preview-content" @click.stop>
          <!-- 右上角关闭按钮 -->
          <button class="close-btn" @click="closePreview">
            <font-awesome-icon icon="times" />
          </button>

          <!-- PDF预览区域 -->
          <div v-if="currentPresentation" class="preview-body">
            <div v-if="currentPresentation.file_type === 'pdf'" class="pdf-preview-container">
              <div class="pdf-viewer-container">
                <canvas ref="pdfCanvas" class="pdf-canvas"></canvas>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 密码验证弹窗 -->
      <div v-if="showPasswordModal" class="password-modal" @click="closePasswordModal">
        <div class="password-content" @click.stop>
          <h3>请输入密码</h3>
          <div class="password-input-group">
            <input
              v-model="passwordInput"
              type="password"
              placeholder="请输入密码"
              class="password-input"
              :disabled="pdfJsLoading"
              @keyup.enter="verifyPassword"
            />
          </div>
          <div class="password-buttons">
            <button
              class="confirm-button"
              :disabled="!passwordInput.trim() || pdfJsLoading"
              @click="verifyPassword"
            >
              <span v-if="pdfJsLoading">加载中...</span>
              <span v-else>确定</span>
            </button>
            <button class="cancel-button" @click="closePasswordModal">取消</button>
          </div>
          <div v-if="passwordError" class="password-error">{{ passwordError }}</div>
          <div v-if="pdfJsLoading" class="loading-hint">正在加载PDF预览组件...</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import NavBar from '@/components/NavBar.vue'
import { getPresentations } from '@/api/presentations'
import backgroundImage from '@/assets/background-image.jpg'
import presentationBackgroundImg from '@/assets/presentation_background.webp'

const presentationBackground = presentationBackgroundImg

// 响应式数据
const loading = ref(true)
const presentations = ref([])
const currentPage = ref(1)
const itemsPerPage = 9 // 每页9张卡片（3行×3列）
const showPreview = ref(false)
const currentPresentation = ref(null)
const totalPages = ref(1)
const totalCount = ref(0)

// 错误处理
const hasError = ref(false)
const errorMessage = ref('')

// 密码验证相关
const showPasswordModal = ref(false)
const passwordInput = ref('')
const passwordError = ref('')
const pendingPresentation = ref(null)

// PDF.js加载状态
const pdfJsLoading = ref(false)

// 翻页功能
// PDF相关状态
const pdfCurrentPage = ref(1)
const pdfTotalPages = ref(1)
const pdfCanvas = ref(null)

// 使用普通变量避免Vue响应式代理破坏PDF.js内部结构
let pdfDocument = null
let pdfLoadingTask = null

// 英雄图片
const defaultThumbnail = ref(backgroundImage)

// 计算属性
const currentPagePresentations = computed(() => {
  return presentations.value
})

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }

  return pages
})

// 方法
const loadPresentations = async () => {
  try {
    loading.value = true
    hasError.value = false
    errorMessage.value = ''

    const response = await getPresentations(currentPage.value, itemsPerPage)

    if (response.data && response.data.data) {
      presentations.value = response.data.data
      totalPages.value = response.data.pagination.total_page
      totalCount.value = response.data.pagination.total

      // 调试：检查标签和缩略图数据
      console.log('讲演数据:', presentations.value)
      presentations.value.forEach(presentation => {
        console.log(`讲演 "${presentation.title}" 的标签:`, presentation.tags)
        console.log(`讲演 "${presentation.title}" 的缩略图字段:`, presentation.thumbnail)
        console.log(`讲演 "${presentation.title}" 的文件路径:`, presentation.file_path)
        console.log(`讲演 "${presentation.title}" 的文件类型:`, presentation.file_type)
        console.log(`讲演 "${presentation.title}" 的缩略图URL:`, getThumbnailUrl(presentation.thumbnail))
        console.log('---')
      })
    } else {
      // 如果没有数据，显示空状态
      presentations.value = []
      totalPages.value = 1
      totalCount.value = 0
    }
  } catch (error) {
    console.error('加载讲演数据失败:', error)
    hasError.value = true
    errorMessage.value = error.message || '加载数据失败，请稍后重试'
    presentations.value = []
    totalPages.value = 1
    totalCount.value = 0
  } finally {
    loading.value = false
  }
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadPresentations() // 重新加载数据
  }
}

const openPresentation = async (presentation) => {
  // 如果是私密讲演，显示密码验证弹窗
  if (presentation.is_private) {
    pendingPresentation.value = presentation
    showPasswordModal.value = true
    passwordInput.value = ''
    passwordError.value = ''
  } else {
    // 公开讲演：按需加载PDF.js然后打开预览
    try {
      // 检查PDF.js是否已加载，如果没有则加载
      if (typeof window.pdfjsLib === 'undefined') {
        console.log('按需加载PDF.js...')
        await loadPdfJs()
      }

      // 打开预览
      currentPresentation.value = presentation
      showPreview.value = true

      // 等待DOM更新后加载PDF
      nextTick(() => {
        loadPdfDocument(presentation)
      })

      // 添加键盘事件监听
      document.addEventListener('keydown', handleKeydown)
    } catch (error) {
      console.error('加载PDF.js失败:', error)
      // 可以显示错误提示给用户
    }
  }
}

const closePreview = () => {
  showPreview.value = false
  currentPresentation.value = null

  // 彻底清理PDF资源
  if (pdfLoadingTask) {
    pdfLoadingTask.destroy()
    pdfLoadingTask = null
  }

  if (pdfDocument) {
    pdfDocument.destroy()
    pdfDocument = null
  }

  // 重置状态
  pdfCurrentPage.value = 1
  pdfTotalPages.value = 1

  // 移除键盘事件监听
  document.removeEventListener('keydown', handleKeydown)
}

const handleKeydown = (event) => {
  // 只在PDF预览模式下处理键盘事件
  if (!showPreview.value) return

  switch (event.key) {
    case 'Escape':
      // 退出预览
      event.preventDefault()
      closePreview()
      break
    case 'ArrowLeft':
    case 'ArrowUp':
      // 上一页
      event.preventDefault()
      goToPreviousPage()
      break
    case 'ArrowRight':
    case 'ArrowDown':
      // 下一页
      event.preventDefault()
      goToNextPage()
      break
  }
}

const goToPreviousPage = () => {
  if (pdfCurrentPage.value > 1) {
    pdfCurrentPage.value--
    renderPdfPage()
  }
}

const goToNextPage = () => {
  if (pdfCurrentPage.value < pdfTotalPages.value) {
    pdfCurrentPage.value++
    renderPdfPage()
  }
}

const loadPdfDocument = async (presentation) => {
  try {
    console.log('开始加载PDF文档:', presentation)

    // 先清理之前的资源
    if (pdfLoadingTask) {
      try {
        pdfLoadingTask.destroy()
      } catch (e) {
        console.warn('销毁旧PDF加载任务时出错:', e)
      }
      pdfLoadingTask = null
    }
    if (pdfDocument) {
      try {
        pdfDocument.destroy()
      } catch (e) {
        console.warn('销毁旧PDF文档时出错:', e)
      }
      pdfDocument = null
    }

    // 检查PDF.js是否已加载
    if (typeof window.pdfjsLib === 'undefined') {
      console.error('PDF.js未加载，尝试重新加载')
      await loadPdfJs()
    }

    const pdfUrl = `http://localhost:3000/api/presentations/${presentation.id}/preview`
    console.log('PDF URL:', pdfUrl)

    // 设置PDF.js worker，利用浏览器缓存
    window.pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/2.16.105/pdf.worker.min.js'

    // 使用简单的配置加载PDF文档
    pdfLoadingTask = window.pdfjsLib.getDocument({
      url: pdfUrl,
      disableAutoFetch: true,
      disableStream: true
    })

    pdfDocument = await pdfLoadingTask.promise
    console.log('PDF文档加载成功，总页数:', pdfDocument.numPages)

    // 获取总页数
    pdfTotalPages.value = pdfDocument.numPages
    pdfCurrentPage.value = 1

    // 等待DOM更新后再渲染
    await nextTick()
    await renderPdfPage()
  } catch (error) {
    console.error('加载PDF失败:', error)
    console.error('错误详情:', error.message)

    // 清理资源
    if (pdfLoadingTask) {
      try {
        pdfLoadingTask.destroy()
      } catch (e) {
        console.warn('清理PDF加载任务时出错:', e)
      }
      pdfLoadingTask = null
    }
    if (pdfDocument) {
      try {
        pdfDocument.destroy()
      } catch (e) {
        console.warn('清理PDF文档时出错:', e)
      }
      pdfDocument = null
    }
  }
}

const renderPdfPage = async () => {
  if (!pdfDocument || !pdfCanvas.value) {
    console.error('PDF文档或Canvas未准备好')
    return
  }

  try {
    console.log('开始渲染PDF页面:', pdfCurrentPage.value)

    const page = await pdfDocument.getPage(pdfCurrentPage.value)
    const canvas = pdfCanvas.value
    const context = canvas.getContext('2d')

    // 获取预览内容区域的实际可用尺寸（完全填充）
    const previewBody = document.querySelector('.preview-body')
    if (!previewBody) {
      console.error('无法找到预览内容区域')
      return
    }

    const availableWidth = previewBody.clientWidth
    const availableHeight = previewBody.clientHeight // 完全使用可用高度

    // 检查可用尺寸是否有效
    if (availableWidth <= 0 || availableHeight <= 0) {
      console.warn('可用区域尺寸无效，等待DOM更新')
      // 延迟重试
      setTimeout(() => {
        if (showPreview.value && pdfDocument && pdfCanvas.value) {
          renderPdfPage()
        }
      }, 100)
      return
    }

    console.log('可用区域尺寸:', { availableWidth, availableHeight })

    // 计算缩放比例，完全填充可用区域
    const viewport = page.getViewport({ scale: 1 })
    console.log('PDF页面原始尺寸:', { width: viewport.width, height: viewport.height })

    // 计算保持宽高比的缩放比例
    const scaleX = availableWidth / viewport.width
    const scaleY = availableHeight / viewport.height
    const scale = Math.min(scaleX, scaleY) // 使用较小的比例确保完全显示
    console.log('计算缩放比例:', { scaleX, scaleY, finalScale: scale })

    // 创建缩放后的视口
    const scaledViewport = page.getViewport({ scale })

    // 设置canvas的显示尺寸（精确匹配缩放后的尺寸）
    canvas.style.width = scaledViewport.width + 'px'
    canvas.style.height = scaledViewport.height + 'px'

    // 设置canvas的实际像素尺寸（考虑设备像素比）
    const devicePixelRatio = window.devicePixelRatio || 1
    canvas.width = Math.floor(scaledViewport.width * devicePixelRatio)
    canvas.height = Math.floor(scaledViewport.height * devicePixelRatio)

    // 缩放上下文以匹配设备像素比
    context.scale(devicePixelRatio, devicePixelRatio)

    // 渲染页面
    const renderContext = {
      canvasContext: context,
      viewport: scaledViewport
    }

    console.log('最终Canvas尺寸:', {
      显示尺寸: `${scaledViewport.width}px × ${scaledViewport.height}px`,
      实际像素: `${canvas.width} × ${canvas.height}`,
      缩放比例: scale
    })

    await page.render(renderContext).promise
    console.log('PDF页面渲染完成')
  } catch (error) {
    console.error('渲染PDF页面失败:', error)
    console.error('错误详情:', error.message)
  }
}

const closePasswordModal = () => {
  showPasswordModal.value = false
  pendingPresentation.value = null
  passwordInput.value = ''
  passwordError.value = ''
}

const verifyPassword = async () => {
  if (!passwordInput.value.trim()) {
    passwordError.value = '请输入密码'
    return
  }

  try {
    const response = await fetch(`http://localhost:3000/api/presentations/${pendingPresentation.value.id}/verify-password`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        password: passwordInput.value
      })
    })

    const result = await response.json()

    if (response.ok && result.success) {
      // 密码验证成功，按需加载PDF.js然后打开预览
      try {
        // 检查PDF.js是否已加载，如果没有则加载
        if (typeof window.pdfjsLib === 'undefined') {
          console.log('密码验证成功，按需加载PDF.js...')
          await loadPdfJs()
        }

        // 打开预览
        currentPresentation.value = pendingPresentation.value
        showPreview.value = true
        showPasswordModal.value = false
        passwordInput.value = ''
        passwordError.value = ''

        // 加载PDF文档
        loadPdfDocument(pendingPresentation.value)
        // 添加键盘事件监听
        document.addEventListener('keydown', handleKeydown)
      } catch (error) {
        console.error('加载PDF.js失败:', error)
        passwordError.value = 'PDF加载失败，请重试'
      }
    } else {
      // 密码错误
      passwordError.value = result.error || '密码错误'
    }
  } catch (error) {
    console.error('密码验证失败:', error)
    passwordError.value = '验证失败，请重试'
  }
}

const getThumbnailUrl = (thumbnail) => {
  console.log('getThumbnailUrl 输入:', thumbnail)

  if (!thumbnail) {
    console.log('缩略图为空，使用默认图片')
    return defaultThumbnail.value
  }

  // 检查是否是PDF文件路径（错误情况）
  if (thumbnail.includes('.pdf') || (thumbnail.includes('presentations/') && !thumbnail.includes('thumbnail'))) {
    console.warn('警告：缩略图字段似乎指向PDF文件而不是封面图片:', thumbnail)
    console.log('使用默认图片替代')
    return defaultThumbnail.value
  }

  // 检查是否是完整的URL
  if (thumbnail.startsWith('http://') || thumbnail.startsWith('https://')) {
    console.log('使用完整URL:', thumbnail)
    // 添加时间戳强制刷新图片
    const separator = thumbnail.includes('?') ? '&' : '?'
    return `${thumbnail}${separator}t=${Date.now()}`
  }

  // 如果是相对路径，构建完整URL
  const fullUrl = `http://localhost:3000${thumbnail}`
  console.log('构建完整URL:', fullUrl)
  // 添加时间戳强制刷新图片
  const separator = fullUrl.includes('?') ? '&' : '?'
  return `${fullUrl}${separator}t=${Date.now()}`
}

const handleImageError = (event) => {
  console.log('图片加载失败:', event.target.src)
  event.target.src = defaultThumbnail.value
}

const handleImageLoad = (event) => {
  console.log('图片加载成功:', event.target.src)
  console.log('图片尺寸:', event.target.naturalWidth, 'x', event.target.naturalHeight)
  console.log('显示尺寸:', event.target.width, 'x', event.target.height)
}

onMounted(async () => {
  try {
    // 添加全局错误处理
    window.addEventListener('error', (event) => {
      console.error('全局错误:', event.error)
      console.error('错误详情:', event.error?.message)
      console.error('错误堆栈:', event.error?.stack)
    })

    // 添加未捕获的Promise错误处理
    window.addEventListener('unhandledrejection', (event) => {
      console.error('未捕获的Promise错误:', event.reason)
      console.error('错误详情:', event.reason?.message)
    })

    // 移除组件挂载时的PDF.js加载，改为按需加载
    loadPresentations()

    // 添加窗口大小变化监听
    window.addEventListener('resize', handleResize)
  } catch (error) {
    console.error('组件初始化错误:', error)
  }
})

onUnmounted(() => {
  // 组件卸载时清理PDF资源
  if (pdfLoadingTask) {
    pdfLoadingTask.destroy()
    pdfLoadingTask = null
  }
  if (pdfDocument) {
    pdfDocument.destroy()
    pdfDocument = null
  }
  // 移除键盘事件监听
  document.removeEventListener('keydown', handleKeydown)
  // 移除窗口大小变化监听
  window.removeEventListener('resize', handleResize)
})

// 重试函数
const retry = () => {
  hasError.value = false
  errorMessage.value = ''
  loadPresentations()
}

// 窗口大小变化处理（防抖）
let resizeTimeout = null
const handleResize = () => {
  // 清除之前的定时器
  if (resizeTimeout) {
    clearTimeout(resizeTimeout)
  }

  // 设置新的定时器，延迟执行
  resizeTimeout = setTimeout(() => {
    // 如果正在预览PDF，重新渲染当前页面以适应新尺寸
    if (showPreview.value && pdfDocument && pdfCanvas.value) {
      renderPdfPage()
    }
  }, 100) // 100ms防抖
}

const loadPdfJs = () => {
  return new Promise((resolve, reject) => {
    try {
      // 检查PDF.js是否已经加载
      if (window.pdfjsLib) {
        console.log('PDF.js已加载，版本:', window.pdfjsLib.version)
        resolve()
        return
      }

      // 检查是否正在加载中
      if (window.pdfjsLoading) {
        console.log('PDF.js正在加载中，等待完成...')
        // 等待加载完成
        const checkLoading = () => {
          if (window.pdfjsLib) {
            resolve()
          } else if (window.pdfjsLoading) {
            setTimeout(checkLoading, 100)
          } else {
            reject(new Error('PDF.js加载失败'))
          }
        }
        checkLoading()
        return
      }

      // 设置加载状态
      window.pdfjsLoading = true
      pdfJsLoading.value = true

      // 移除可能存在的旧脚本
      const existingScripts = document.querySelectorAll('script[src*="pdf.js"]')
      existingScripts.forEach(script => {
        try {
          script.remove()
        } catch (e) {
          console.warn('移除旧脚本时出错:', e)
        }
      })

      // 使用稳定的PDF.js 2.x版本，移除时间戳以利用浏览器缓存
      const script = document.createElement('script')
      script.src = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/2.16.105/pdf.min.js'

      script.onload = () => {
        try {
          // 检查PDF.js是否正确加载
          if (!window.pdfjsLib) {
            throw new Error('PDF.js未正确加载')
          }

          // 设置worker路径，利用浏览器缓存
          window.pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/2.16.105/pdf.worker.min.js'

          // 清除加载状态
          window.pdfjsLoading = false
          pdfJsLoading.value = false

          console.log('PDF.js加载成功，版本:', window.pdfjsLib.version)
          resolve()
        } catch (error) {
          window.pdfjsLoading = false
          pdfJsLoading.value = false
          console.error('PDF.js初始化失败:', error)
          reject(error)
        }
      }

      script.onerror = (error) => {
        window.pdfjsLoading = false
        pdfJsLoading.value = false
        console.error('PDF.js脚本加载失败:', error)
        reject(new Error('PDF.js脚本加载失败'))
      }

      document.head.appendChild(script)
    } catch (error) {
      window.pdfjsLoading = false
      pdfJsLoading.value = false
      console.error('loadPdfJs函数执行失败:', error)
      reject(error)
    }
  })
}
</script>

<style scoped>
.presentation-view {
  min-height: 100vh;
  position: relative;
  padding-top: 40px;
}

/* 英雄图片区域 */
.hero-section {
  position: relative;
  padding: 30px 350px 40px 350px;
}

.hero-image-wrapper {
  position: relative;
  width: 100%;
}

.hero-image-wrapper .hero-image-full {
  width: 100%;
  height: auto;
  object-fit: contain;
  border-radius: 8px;
  display: block;
}

.hero-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(0, 0, 0, 0.3) 0%, rgba(0, 0, 0, 0.1) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-content {
  text-align: center;
  color: white;
}

.hero-title {
  font-size: 3rem;
  font-weight: 700;
  margin: 0 0 10px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.hero-subtitle {
  font-size: 1.2rem;
  margin: 0;
  opacity: 0.9;
}

/* 分隔区域 */
.separator-section {
  padding: 40px 0;
  text-align: left;
}

.separator-section .separator-text {
  max-width: 1200px;
  margin: 0 auto;
  padding: 10px 20px;
}

.separator-text {
  font-size: 1.1rem;
  color: #333;
  margin: 0;
  font-style: italic;
  padding: 10px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  border-left: 4px solid #667eea;
}

/* 内容区域 */
.content-section {
  padding: 60px 0;
  background: transparent;
}

.content-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 加载状态 */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 0;
  color: #666;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 卡片网格 */
.cards-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 25px; /* 从30px减少到25px，让卡片更宽 */
  margin-bottom: 60px;
}

.presentation-card {
  position: relative;
  background: white;
  border-radius: 5px;
  overflow: visible;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.presentation-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.card-image {
  width: 100%;
  height: 220px; /* 从280px减少到220px，让卡片更矮 */
  overflow: hidden;
  position: relative;
  background: white;
  z-index: 1;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: contain; /* 确保完整显示图片 */
  object-position: center top; /* 图片顶部对齐 */
  transition: transform 0.3s ease;
  background: white;
  display: block; /* 移除默认的inline间距 */
}

.presentation-card:hover .card-image img {
  transform: scale(1.05);
}

.card-content {
  position: relative;
  z-index: 2;
  padding: 20px; /* 从25px减少到20px，让卡片更紧凑 */
  text-align: left;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.card-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
  margin: 0 0 18px 0;
  line-height: 1.4;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  text-align: left;
  flex-shrink: 0; /* 防止标题被压缩 */
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.title-icon {
  color: #667eea;
  font-size: 1.1rem;
  margin-top: 2px;
  flex-shrink: 0;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-start;
  margin-top: auto; /* 将标签推到底部 */
  max-height: 60px; /* 限制标签区域最大高度 */
  overflow: hidden; /* 超出部分隐藏 */
}

.tag {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 6px 14px; /* 从4px 12px增加到6px 14px */
  border-radius: 15px;
  font-size: 0.85rem; /* 从0.8rem增加到0.85rem */
  font-weight: 500;
}

/* 分页控件 */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 40px;
}

.pagination-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #666;
}

.pagination-btn:hover:not(:disabled) {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-pages {
  display: flex;
  gap: 5px;
}

.page-btn {
  width: 40px;
  height: 40px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #666;
  font-weight: 500;
}

.page-btn:hover {
  background: #f5f5f5;
}

.page-btn.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

/* 错误边界样式 */
.error-boundary {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 50vh;
  padding: 2rem;
  text-align: center;
}

.error-boundary h2 {
  color: #e74c3c;
  margin-bottom: 1rem;
}

.error-boundary p {
  color: #666;
  margin-bottom: 2rem;
  max-width: 500px;
}

.retry-btn {
  background: #3498db;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.retry-btn:hover {
  background: #2980b9;
}

/* 预览模态框 */
.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(5px);
}

.preview-content {
  position: relative;
  width: 90vw;    /* 使用视口宽度 */
  height: 87vh;   /* 使用视口高度 */
  background: white;
  border-radius: 10px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  /* 确保内容完全填充，没有内边距 */
  padding: 0;
  margin: 0;
}

/* 右上角关闭按钮 */
.close-btn {
  position: absolute;
  top: 15px;
  right: 15px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  font-size: 18px;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.9);
  transform: scale(1.1);
}

/* PDF预览区域 */
.preview-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden; /* 确保没有滚动条 */
  position: relative;
  /* 确保没有边距和内边距 */
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}

.pdf-preview-container {
  width: 100%;
  height: 100%;
  background: white;
  display: flex;
  flex-direction: column;
  /* 确保没有边距和内边距 */
  margin: 0;
  padding: 0;
}

.pdf-viewer-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  /* 完全填充容器 */
  width: 100%;
  height: 100%;
  overflow: hidden; /* 隐藏任何可能的溢出 */
  /* 确保没有边距和内边距 */
  margin: 0;
  padding: 0;
}

.pdf-canvas {
  /* 完全填充容器 */
  display: block;
  background: white;
  width: 100%;
  height: 100%;
  object-fit: contain; /* 保持比例并适应容器 */
  /* 确保没有边距和内边距 */
  margin: 0;
  padding: 0;
  border: none;
}

/* 密码验证弹窗样式 */
.preview-info-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.preview-info {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.file-info {
  flex: 1;
  display: flex;
  gap: 15px;
  align-items: flex-start;
}

.file-icon {
  font-size: 3rem;
  color: #d63384;
  margin-top: 5px;
}

/* PDF文件图标颜色 */
.file-icon svg[data-icon="file-pdf"] {
  color: #dc3545;
}

.file-details h4 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 1.3rem;
}

.file-details p {
  margin: 5px 0;
  color: #666;
}

.file-type {
  font-weight: 600;
  color: #d63384 !important;
}

.file-size {
  font-size: 0.9rem;
}

.file-description {
  font-style: italic;
  margin-top: 10px !important;
}

.preview-thumbnail {
  width: 200px;
  height: 150px;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.preview-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-notice {
  background: #e3f2fd;
  border: 1px solid #bbdefb;
  border-radius: 8px;
  padding: 15px;
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.preview-notice svg {
  color: #1976d2;
  margin-top: 2px;
  flex-shrink: 0;
}

.preview-notice p {
  margin: 0;
  color: #1565c0;
  font-size: 0.9rem;
  line-height: 1.4;
}

.preview-actions {
  padding: 20px;
  border-top: 1px solid #eee;
  background: #f8f9fa;
  display: flex;
  gap: 15px;
  justify-content: center;
}

.download-btn,
.preview-btn,
.fullscreen-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
}

.download-btn {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
}

.download-btn:hover {
  background: linear-gradient(135deg, #5a6fd8, #6a4190);
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.preview-btn {
  background: #17a2b8;
  color: white;
}

.preview-btn:hover {
  background: #138496;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(23, 162, 184, 0.4);
}

.fullscreen-btn {
  background: #28a745;
  color: white;
}

.fullscreen-btn:hover {
  background: #218838;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(40, 167, 69, 0.4);
}

.close-btn {
  position: absolute;
  top: 15px;
  right: 15px;
  width: 40px;
  height: 40px;
  border: none;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  border-radius: 50%;
  cursor: pointer;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.7);
}

.preview-container {
  width: 100%;
  height: 100%;
}

.preview-iframe {
  width: 100%;
  height: 100%;
  border: none;
}

.preview-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
  font-size: 1.2rem;
}

.preview-placeholder svg {
  font-size: 4rem;
  margin-bottom: 20px;
  color: #667eea;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .cards-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
  }

  .hero-title {
    font-size: 2.5rem;
  }
}

/* 紧凑模式：头图与内容等宽（2/3 居中），高度不限制 */
@media (max-width: 1330px) {
  .presentation-view { padding-top: 40px; }
  .hero-section { padding-left: 0; padding-right: 0; }
  .hero-image-wrapper, .content-container, .separator-section .separator-text { width: 66.666%; margin: 0 auto; min-width: 480px; }
}

@media (max-width: 768px) {
  .cards-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .hero-section {
    height: 300px;
  }

  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .content-container {
    padding: 0 15px;
  }

  .pagination {
    flex-wrap: wrap;
    gap: 5px;
  }

  .page-btn, .pagination-btn {
    width: 35px;
    height: 35px;
    font-size: 0.9rem;
  }
}

@media (max-width: 480px) {
  .hero-title {
    font-size: 1.8rem;
  }

  .separator-text {
    font-size: 1rem;
  }

  .card-title {
    font-size: 1rem;
  }

  .tag {
    font-size: 0.75rem;
    padding: 3px 10px;
  }
}

/* 预览模态框响应式设计 */
@media (max-width: 768px) {
  .preview-content {
    width: 95%;
    height: 95%;
  }

  .preview-info {
    flex-direction: column;
    gap: 15px;
  }

  .preview-thumbnail {
    width: 100%;
    height: 120px;
  }

  .file-icon {
    font-size: 2rem;
  }

  .preview-actions {
    flex-direction: column;
    gap: 10px;
  }

  .download-btn,
  .preview-btn,
  .fullscreen-btn {
    width: 100%;
    justify-content: center;
  }
}

/* 密码验证弹窗样式 */
.password-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3000;
  backdrop-filter: blur(5px);
}

.password-content {
  background: white;
  border-radius: 16px;
  padding: 40px;
  width: 380px;
  max-width: 90vw;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.3);
  animation: modalSlideIn 0.3s ease-out;
  text-align: center;
}

.password-content h3 {
  margin: 0 0 30px 0;
  color: #333;
  font-size: 1.4rem;
  font-weight: 600;
}

.password-input-group {
  margin-bottom: 25px;
}

.password-input {
  width: 100%;
  padding: 16px 20px;
  border: 2px solid #e1e5e9;
  border-radius: 12px;
  font-size: 1rem;
  transition: all 0.3s ease;
  background: #fafafa;
  box-sizing: border-box;
}

.password-input:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.password-input:disabled {
  background: #f0f0f0;
  color: #999;
  cursor: not-allowed;
}

.password-buttons {
  display: flex;
  gap: 15px;
  justify-content: center;
}

.confirm-button {
  padding: 14px 28px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  min-width: 80px;
}

.confirm-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.confirm-button:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.cancel-button {
  padding: 14px 28px;
  background: white;
  border: 2px solid #e1e5e9;
  color: #666;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  min-width: 80px;
}

.cancel-button:hover {
  border-color: #999;
  color: #333;
  background: #f9f9f9;
  transform: translateY(-1px);
}

.password-error {
  color: #e74c3c;
  font-size: 0.9rem;
  margin-top: 15px;
  text-align: center;
}

.loading-hint {
  margin: 15px 0 0 0;
  color: #667eea;
  font-size: 0.9rem;
  font-style: italic;
  padding: 8px;
  background: #f0f4ff;
  border-radius: 8px;
  border-left: 4px solid #667eea;
}

/* 密码弹窗响应式设计 */
@media (max-width: 480px) {
  .password-content {
    width: 95%;
  }

  .password-input-group {
    flex-direction: column;
  }
  .verify-btn {
    width: 100%;
  }
}
</style>
