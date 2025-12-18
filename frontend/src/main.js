/* eslint-disable space-before-function-paren */
/* eslint-disable no-multi-spaces */
/* eslint-disable semi */
/* eslint-disable no-undef */
/* eslint-disable indent */
/* eslint-disable quotes */
import './assets/main.css'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import store from './store'
import './assets/styles/github-markdown.css'
import './assets/styles/github-highlight.css'
import headImage from './assets/my_headportrait.jpg'
import { initErrorReporter, handleVueError } from './utils/errorReporter'
// 首屏必需的图标（NavBar 中使用）
import {
  faHouse,
  faBlog,
  faMagnifyingGlass,
  faQuestion,
  faEllipsis,
  faBars,
  faCalendar,
  faImages,
  faGear,
  faLocationDot,
  faChalkboard,
  faEnvelope,
  faCommentDots
} from '@fortawesome/free-solid-svg-icons'

// 其他图标按需加载（通过 iconLoader 工具）
// 注意：Font Awesome 的 tree-shaking 会自动移除未使用的图标
// 但为了进一步优化首屏加载，我们将非首屏图标改为按需加载

// 早期错误处理：在应用初始化前捕获浏览器扩展错误（如 MetaMask）
if (typeof window !== 'undefined') {
  // 捕获未处理的 Promise 拒绝（来自 MetaMask 等扩展）
  window.addEventListener('unhandledrejection', (event) => {
    const reasonStr = String(event.reason || '')
    if (
      reasonStr.includes('MetaMask') ||
      reasonStr.includes('Failed to connect to MetaMask') ||
      reasonStr.includes('ethereum') ||
      event.reason?.message?.includes('MetaMask') ||
      event.reason?.message?.includes('Failed to connect to MetaMask')
    ) {
      event.preventDefault() // 阻止错误显示
      return false
    }
  }, { capture: true })

  // 捕获运行时错误（来自浏览器扩展）
  window.addEventListener('error', (event) => {
    if (event.filename && (
      event.filename.includes('chrome-extension://') ||
      event.filename.includes('moz-extension://') ||
      event.filename.includes('safari-extension://') ||
      event.filename.includes('extension://')
    )) {
      event.preventDefault() // 阻止错误显示
      return false
    }
    if (event.message && (
      event.message.includes('MetaMask') ||
      event.message.includes('Failed to connect to MetaMask')
    )) {
      event.preventDefault() // 阻止错误显示
      return false
    }
  }, { capture: true })
}

// 只添加首屏必需的图标
library.add(
  faHouse,
  faBlog,
  faMagnifyingGlass,
  faQuestion,
  faEllipsis,
  faBars,
  faCalendar,
  faImages,
  faGear,
  faLocationDot,
  faChalkboard,
  faEnvelope,
  faCommentDots
)

// 预加载其他常用图标（在页面加载完成后）
if (typeof window !== 'undefined') {
  window.addEventListener('load', async () => {
    // 延迟加载其他图标，避免阻塞首屏渲染
    setTimeout(async () => {
      try {
        const { loadIcons } = await import('@/utils/iconLoader')
        // 预加载常用图标
        await loadIcons([
          { name: 'diagramProject', style: 'solid' },
          { name: 'flask', style: 'solid' },
          { name: 'penToSquare', style: 'solid' },
          { name: 'comment', style: 'solid' },
          { name: 'reply', style: 'solid' },
          { name: 'clock', style: 'solid' },
          { name: 'share', style: 'solid' },
          { name: 'arrowUp', style: 'solid' },
          { name: 'tag', style: 'solid' },
          { name: 'heart', style: 'solid' },
          { name: 'sun', style: 'solid' },
          { name: 'moon', style: 'solid' },
          { name: 'times', style: 'solid' },
          { name: 'edit', style: 'solid' },
          { name: 'copy', style: 'solid' },
          { name: 'star', style: 'solid' },
          { name: 'star', style: 'regular' },
          { name: 'trash', style: 'solid' },
          { name: 'plus', style: 'solid' },
          { name: 'github', style: 'brands' },
          { name: 'weixin', style: 'brands' },
          { name: 'weibo', style: 'brands' },
          { name: 'bilibili', style: 'brands' }
        ])
      } catch (error) {
        console.warn('预加载图标失败:', error)
      }
    }, 1000) // 延迟1秒加载，确保首屏优先
  })
}

const app = createApp(App)

// 配置 Vue 全局错误处理
app.config.errorHandler = handleVueError

app.use(router)
app.use(store)
app.component('FontAwesomeIcon', FontAwesomeIcon)
app.mount('#app')

// 初始化前端错误监控
initErrorReporter()

// 全局处理外链：所有外链都在新标签页打开
if (typeof window !== 'undefined') {
  // 判断是否是外链
  const isExternalLink = (url) => {
    if (!url) return false

    // 排除空链接、锚点链接、JavaScript链接、邮件链接等
    if (url.startsWith('#') ||
        url.startsWith('javascript:') ||
        url.startsWith('mailto:') ||
        url.startsWith('tel:') ||
        url.trim() === '') {
      return false
    }

    // 相对路径（以 / 开头）不是外链
    if (url.startsWith('/')) {
      return false
    }

    // 检查是否是完整的 URL
    try {
      const linkUrl = new URL(url, window.location.href)
      const currentUrl = new URL(window.location.href)

      // 如果协议、主机名、端口都相同，则是内链
      return linkUrl.origin !== currentUrl.origin
    } catch (e) {
      // 如果 URL 解析失败，可能是相对路径，不是外链
      return false
    }
  }

  // 全局点击事件监听器
  document.addEventListener('click', (event) => {
    // 查找最近的 <a> 标签
    let target = event.target
    let link = null

    // 向上查找 <a> 标签
    while (target && target !== document) {
      if (target.tagName === 'A' && target.href) {
        link = target
        break
      }
      target = target.parentElement
    }

    // 如果找到了链接
    if (link) {
      const href = link.getAttribute('href') || link.href

      // 如果已经是 target="_blank"，不需要处理
      if (link.target === '_blank') {
        return
      }

      // 排除 Vue Router 的 router-link（它们会阻止默认行为）
      if (link.hasAttribute('data-router-link') ||
          link.classList.contains('router-link-active') ||
          link.classList.contains('router-link-exact-active')) {
        return
      }

      // 如果是外链
      if (isExternalLink(href)) {
        // 阻止默认行为
        event.preventDefault()
        event.stopPropagation()

        // 在新标签页打开
        window.open(href, '_blank', 'noopener,noreferrer')

        return false
      }
    }
  }, true) // 使用捕获阶段，确保优先处理
}

// 设置浏览器标签栏标题与图标
try {
  document.title = '山角函兽的小窝'
  const setFavicon = (href) => {
    let link = document.querySelector("link[rel='icon']")
    if (!link) {
      link = document.createElement('link')
      link.rel = 'icon'
      document.head.appendChild(link)
    }
    link.href = href
  }
  // 生成带圆角的 favicon
  const createRoundedFavicon = (src, size = 64, radius = 12) => new Promise((resolve, reject) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => {
      try {
        const canvas = document.createElement('canvas')
        canvas.width = size
        canvas.height = size
        const ctx = canvas.getContext('2d')

        const r = Math.min(radius, size / 2)
        const w = size
        const h = size
        ctx.beginPath()
        ctx.moveTo(r, 0)
        ctx.arcTo(w, 0, w, h, r)
        ctx.arcTo(w, h, 0, h, r)
        ctx.arcTo(0, h, 0, 0, r)
        ctx.arcTo(0, 0, w, 0, r)
        ctx.closePath()
        ctx.clip()

        // 等比裁剪
        const minSide = Math.min(img.width, img.height)
        const sx = (img.width - minSide) / 2
        const sy = (img.height - minSide) / 2
        ctx.drawImage(img, sx, sy, minSide, minSide, 0, 0, size, size)
        resolve(canvas.toDataURL('image/png'))
      } catch (e) {
        reject(e)
      }
    }
    img.onerror = () => reject(new Error('favicon image load error'))
    img.src = src
  })

  createRoundedFavicon(headImage, 64, 12)
    .then(url => setFavicon(url))
    .catch(() => setFavicon(headImage))
} catch (e) {
  // 忽略在非浏览器环境下的错误
}
