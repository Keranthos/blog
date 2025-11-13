# 讲演 PDF 预览：按需加载与安全策略

## 设计目标

讲演模块需要同时满足两个诉求：一方面让用户点开卡片就能即时预览 PDF，另一方面又要确保私密讲演不会被未授权访问。实现思路是“密码门槛 + 按需加载 PDF.js + 键盘友好的阅读体验”，并在关闭时彻底释放资源。

## 按需加载 PDF.js

PDF.js 体积不小，因此在首次需要时才动态加载脚本，并利用全局标记避免重复请求：

```715:792:frontend/src/views/PresentationView.vue
const loadPdfJs = () => {
  return new Promise((resolve, reject) => {
    try {
      if (window.pdfjsLib) {
        resolve()
        return
      }

      if (window.pdfjsLoading) {
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

      window.pdfjsLoading = true
      pdfJsLoading.value = true

      const existingScripts = document.querySelectorAll('script[src*="pdf.js"]')
      existingScripts.forEach(script => {
        try { script.remove() } catch (e) { console.warn('移除旧脚本时出错:', e) }
      })

      const script = document.createElement('script')
      script.src = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/2.16.105/pdf.min.js'

      script.onload = () => {
        if (!window.pdfjsLib) {
          reject(new Error('PDF.js未正确加载'))
          return
        }
        window.pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/2.16.105/pdf.worker.min.js'
        window.pdfjsLoading = false
        pdfJsLoading.value = false
        resolve()
      }

      script.onerror = (error) => {
        window.pdfjsLoading = false
        pdfJsLoading.value = false
        reject(new Error('PDF.js脚本加载失败'))
      }

      document.head.appendChild(script)
    } catch (error) {
      window.pdfjsLoading = false
      pdfJsLoading.value = false
      reject(error)
    }
  })
}
```

通过 `window.pdfjsLib` 与 `window.pdfjsLoading` 这两个全局变量，脚本只会被下载一次，同时在加载失败时能正确恢复状态。

## 密码保护的讲演

私密讲演在打开前会弹出密码框，只有验证通过才会继续加载 PDF：

```554:604:frontend/src/views/PresentationView.vue
const verifyPassword = async () => {
  if (!passwordInput.value.trim()) {
    passwordError.value = '请输入密码'
    return
  }

  try {
    const response = await fetch(`/api/presentations/${pendingPresentation.value.id}/verify-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password: passwordInput.value })
    })

    const result = await response.json()

    if (response.ok && result.success) {
      if (typeof window.pdfjsLib === 'undefined') {
        await loadPdfJs()
      }
      currentPresentation.value = pendingPresentation.value
      showPreview.value = true
      showPasswordModal.value = false
      passwordInput.value = ''
      passwordError.value = ''
      loadPdfDocument(pendingPresentation.value)
      document.addEventListener('keydown', handleKeydown)
    } else {
      passwordError.value = result.error || '密码错误'
    }
  } catch (error) {
    passwordError.value = '验证失败，请重试'
  }
}
```

验证成功后才会真正拉取 PDF.js 并渲染文档，从而避免私密文件被随意访问。

## 自适应渲染与键盘交互

渲染时会依据当前窗口宽度和 `devicePixelRatio` 计算缩放比例，保证在 Retina 屏上依旧清晰；同时绑定 `ArrowLeft/ArrowRight` 控制翻页，`Escape` 关闭预览。窗口缩放时通过防抖重新渲染当前页，避免画面模糊。

## 资源释放与异常兜底

关闭模态框时会销毁 `pdfLoadingTask` 与 `pdfDocument`，清除键盘事件监听，并重置页码状态。即便加载过程中发生异常，也能通过日志快速定位问题，保证应用在长时间运行后不泄漏内存。

## 小结

通过“按需加载 + 密码验证 + 键盘友好 + 彻底回收”这一套机制，讲演页兼顾了性能与安全。后续可以进一步扩展缩略图导航、页码跳转，或支持 PPTX、视频等更多格式，让知识分享更顺畅。
