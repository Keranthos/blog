<template>
  <div class="edit-presentation-view">
    <NavBar />

    <div class="edit-container">
      <div class="edit-header">
        <h1>创建讲演</h1>
        <p>上传你的PDF文件并添加相关信息</p>
      </div>

      <div class="edit-form">
        <form @submit.prevent="savePresentation">
          <!-- 基本信息 -->
          <div class="form-section">
            <h3>基本信息</h3>

            <div class="form-group">
              <label for="title">标题 *</label>
              <input
                id="title"
                v-model="presentation.title"
                type="text"
                placeholder="请输入讲演标题"
                required
              />
            </div>

            <div class="form-group">
              <label for="tags">标签</label>
              <input
                id="tags"
                v-model="tagsInput"
                type="text"
                placeholder="请输入标签，用逗号分隔"
                @blur="updateTags"
              />
              <div class="tags-preview">
                <span
                  v-for="tag in presentation.tags"
                  :key="tag"
                  class="tag"
                >
                  {{ tag }}
                </span>
              </div>
            </div>

            <div class="form-group">
              <label class="checkbox-label">
                <input
                  v-model="presentation.isPrivate"
                  type="checkbox"
                  class="privacy-checkbox"
                />
                <span class="checkbox-text">设置为私密讲演</span>
              </label>
            </div>
          </div>

          <!-- 文件上传 -->
          <div class="form-section">
            <h3>PDF文件</h3>

            <div class="upload-area" :class="{ 'drag-over': isDragOver }">
              <input
                ref="fileInput"
                type="file"
                accept=".pdf"
                style="display: none"
                @change="handleFileSelect"
              />

              <div v-if="!presentation.file" class="upload-placeholder" @click="triggerFileSelect">
                <font-awesome-icon icon="cloud-upload-alt" />
                <p>点击选择文件或拖拽文件到此处</p>
                <p class="upload-hint">支持 .pdf 格式</p>
              </div>

              <div v-else class="file-info">
                <div class="file-details">
                  <font-awesome-icon icon="file-pdf" />
                  <div class="file-text">
                    <p class="file-name">{{ presentation.file.name }}</p>
                    <p class="file-size">{{ formatFileSize(presentation.file.size) }}</p>
                  </div>
                </div>
                <button type="button" class="remove-file" @click="removeFile">
                  <font-awesome-icon icon="times" />
                </button>
              </div>
            </div>
          </div>

          <!-- 缩略图 -->
          <div class="form-section">
            <h3>缩略图</h3>

            <div class="thumbnail-manager">
              <!-- 图片预览区域 -->
              <div class="thumbnail-preview-section">
                <div
                  class="thumbnail-preview-container"
                  :class="{ 'drag-over': isThumbnailDragOver }"
                  @dragover.prevent="handleThumbnailDragOver"
                  @dragleave.prevent="handleThumbnailDragLeave"
                  @drop.prevent="handleThumbnailDrop"
                >
                  <div v-if="presentation.thumbnail" class="thumbnail-preview">
                    <!-- 加载状态指示器 -->
                    <div v-if="thumbnailLoading" class="thumbnail-loading">
                      <div class="loading-spinner"></div>
                      <span class="loading-text">图片加载中...</span>
                    </div>
                    <!-- 错误状态指示器 -->
                    <div v-else-if="thumbnailError" class="thumbnail-error">
                      <font-awesome-icon icon="exclamation-triangle" class="error-icon" />
                      <span class="error-text">图片加载失败</span>
                      <button type="button" class="retry-btn" @click="retryThumbnailLoad">
                        <font-awesome-icon icon="redo" />
                        重试
                      </button>
                    </div>
                    <!-- 正常图片显示 -->
                    <img
                      v-else
                      :src="presentation.thumbnail"
                      alt="缩略图预览"
                      @error="handleThumbnailError"
                      @load="handleThumbnailLoad"
                    />
                    <div class="thumbnail-overlay">
                      <button type="button" class="clear-btn" @click="removeThumbnail">
                        <font-awesome-icon icon="trash" />
                      </button>
                    </div>
                  </div>
                  <div v-else class="thumbnail-placeholder">
                    <font-awesome-icon icon="image" class="placeholder-icon" />
                    <span class="placeholder-text">暂无缩略图</span>
                    <div class="drag-hint">拖拽图片到此处上传</div>
                  </div>
                </div>
              </div>

              <!-- 图片选择控制区域 -->
              <div class="thumbnail-controls-section">
                <div class="input-group">
                  <label class="input-label">
                    <font-awesome-icon icon="link" />
                    图片链接
                  </label>
                  <input
                    v-model="presentation.thumbnail"
                    type="text"
                    placeholder="请输入图片链接"
                    class="thumbnail-url-input"
                    @input="handleThumbnailUrlInput"
                  />
                </div>
                <div class="thumbnail-controls">
                  <button type="button" class="random-image-btn" @click="getRandomImage">
                    <font-awesome-icon icon="dice" />
                    随机图片
                  </button>
                  <button type="button" class="upload-image-btn" @click="triggerThumbnailUpload">
                    <font-awesome-icon icon="upload" />
                    上传图片
                  </button>
                  <input
                    ref="thumbnailInput"
                    type="file"
                    accept="image/*"
                    style="display: none"
                    @change="handleThumbnailUpload"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="goBack">
              取消
            </button>
            <button type="submit" class="btn btn-primary" :disabled="!isFormValid">
              {{ isUploading ? '上传中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import { createPresentation } from '@/api/presentations'
import { showErrorMessage, showSuccessMessage, showWarningMessage, showCustomMessage } from '@/utils/waifuMessage'

const router = useRouter()
const route = useRoute()

// 响应式数据
const presentation = ref({
  title: '',
  tags: [],
  file: null,
  thumbnail: null,
  thumbnailFile: null, // 添加缩略图文件对象
  isPrivate: false // 是否私密
})

const tagsInput = ref('')
const isDragOver = ref(false)
const isUploading = ref(false)

// 文件输入引用
const fileInput = ref(null)
const thumbnailInput = ref(null)

// 缩略图相关状态
const isThumbnailDragOver = ref(false)
const thumbnailLoading = ref(false)
const thumbnailError = ref(false)

// 计算属性
const isFormValid = computed(() => {
  return presentation.value.title.trim() && presentation.value.file
})

// 方法
const triggerFileSelect = () => {
  fileInput.value.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    presentation.value.file = file
  }
}

const handleThumbnailUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    showWarningMessage('invalid_file_type')
    return
  }

  // 检查文件大小 (10MB)
  if (file.size > 10 * 1024 * 1024) {
    showWarningMessage('file_too_large')
    return
  }

  await uploadCustomThumbnail(file)
}

const removeFile = () => {
  presentation.value.file = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const removeThumbnail = () => {
  presentation.value.thumbnail = null
  presentation.value.thumbnailFile = null
  thumbnailLoading.value = false
  thumbnailError.value = false
  if (thumbnailInput.value) {
    thumbnailInput.value.value = ''
  }
}

// 触发缩略图上传
const triggerThumbnailUpload = () => {
  if (thumbnailInput.value) {
    thumbnailInput.value.click()
  }
}

// 上传自定义缩略图
const uploadCustomThumbnail = async (file) => {
  try {
    // 直接保存文件对象，而不是创建URL
    presentation.value.thumbnailFile = file

    // 创建本地预览URL用于显示
    const localUrl = URL.createObjectURL(file)
    presentation.value.thumbnail = localUrl

    // 设置加载状态
    thumbnailLoading.value = true
    thumbnailError.value = false

    // 异步清除加载状态
    setTimeout(() => {
      thumbnailLoading.value = false
    }, 100)

    showCustomMessage(`图片已选择: ${file.name}`, 3000)
  } catch (error) {
    showErrorMessage('处理失败: ' + error.message)
  }
}

// 获取随机图片

const getRandomImage = async () => {
  try {
    // 这里调用后端API获取随机图片
    // 暂时使用占位符
    const randomImages = [
      'https://picsum.photos/400/300?random=1',
      'https://picsum.photos/400/300?random=2',
      'https://picsum.photos/400/300?random=3',
      'https://picsum.photos/400/300?random=4',
      'https://picsum.photos/400/300?random=5'
    ]

    const randomIndex = Math.floor(Math.random() * randomImages.length)
    const imageUrl = randomImages[randomIndex]

    presentation.value.thumbnail = imageUrl
    thumbnailLoading.value = true
    thumbnailError.value = false

    // 异步清除加载状态
    setTimeout(() => {
      thumbnailLoading.value = false
    }, 100)

    showCustomMessage('已获取随机图片', 3000)
  } catch (error) {
    showErrorMessage('获取随机图片失败: ' + error.message)
  }
}

// 缩略图拖拽处理
const handleThumbnailDragOver = (event) => {
  event.preventDefault()
  isThumbnailDragOver.value = true
}

const handleThumbnailDragLeave = (event) => {
  event.preventDefault()
  isThumbnailDragOver.value = false
}

const handleThumbnailDrop = async (event) => {
  event.preventDefault()
  isThumbnailDragOver.value = false

  const files = event.dataTransfer.files
  if (files.length === 0) return

  const file = files[0]
  if (!file.type.startsWith('image/')) {
    showWarningMessage('invalid_file_type')
    return
  }

  if (file.size > 10 * 1024 * 1024) {
    showWarningMessage('file_too_large')
    return
  }

  await uploadCustomThumbnail(file)
}

// 缩略图URL输入处理
const handleThumbnailUrlInput = () => {
  thumbnailLoading.value = true
  thumbnailError.value = false

  // 异步清除加载状态
  setTimeout(() => {
    thumbnailLoading.value = false
  }, 100)
}

// 缩略图加载错误处理
const handleThumbnailError = () => {
  thumbnailLoading.value = false
  thumbnailError.value = true
}

// 缩略图加载成功处理
const handleThumbnailLoad = () => {
  thumbnailLoading.value = false
  thumbnailError.value = false
}

// 重试缩略图加载
const retryThumbnailLoad = () => {
  thumbnailError.value = false
  thumbnailLoading.value = true

  // 强制重新加载图片
  const img = presentation.value.thumbnail
  presentation.value.thumbnail = ''

  setTimeout(() => {
    presentation.value.thumbnail = img
    thumbnailLoading.value = false
  }, 100)
}

const updateTags = () => {
  if (tagsInput.value.trim()) {
    presentation.value.tags = tagsInput.value
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0)
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const savePresentation = async () => {
  if (!isFormValid.value) return

  try {
    isUploading.value = true

    // 调试：检查用户登录状态
    const token = localStorage.getItem('jwt')
    console.log('当前JWT token:', token ? '存在' : '不存在')
    console.log('Token内容:', token)

    // 创建FormData对象
    const formData = new FormData()
    formData.append('title', presentation.value.title)
    formData.append('tags', JSON.stringify(presentation.value.tags))
    formData.append('is_private', presentation.value.isPrivate.toString())
    formData.append('file', presentation.value.file)

    // 处理缩略图
    if (presentation.value.thumbnailFile) {
      // 如果有上传的缩略图文件，直接添加
      formData.append('thumbnail', presentation.value.thumbnailFile)
    } else if (presentation.value.thumbnail) {
      try {
        // 检查缩略图是否为dataURL格式
        if (presentation.value.thumbnail.startsWith('data:')) {
          // 将base64转换为blob
          const thumbnailBlob = dataURLtoBlob(presentation.value.thumbnail)
          formData.append('thumbnail', thumbnailBlob, 'thumbnail.jpg')
        } else if (presentation.value.thumbnail.startsWith('http')) {
          // 如果是外部URL，直接添加到FormData
          formData.append('thumbnail_url', presentation.value.thumbnail)
        }
      } catch (thumbnailError) {
        console.warn('缩略图处理失败:', thumbnailError)
        // 缩略图处理失败不影响主文件上传
      }
    }

    console.log('准备发送FormData:', formData)

    // 调用后端API上传文件
    const response = await createPresentation(formData)

    if (response.data && response.data.message) {
      showSuccessMessage('upload')
      router.push('/presentation')
    } else {
      throw new Error('上传失败')
    }
  } catch (error) {
    console.error('上传失败:', error)
    console.error('错误详情:', error.response?.data)
    showErrorMessage('upload_failed')
  } finally {
    isUploading.value = false
  }
}

const dataURLtoBlob = (dataURL) => {
  if (!dataURL || typeof dataURL !== 'string') {
    throw new Error('无效的dataURL')
  }

  const arr = dataURL.split(',')
  if (arr.length < 2) {
    throw new Error('dataURL格式不正确')
  }

  const mimeMatch = arr[0].match(/:(.*?);/)
  if (!mimeMatch) {
    throw new Error('无法解析MIME类型')
  }

  const mime = mimeMatch[1]
  const bstr = atob(arr[1])
  let n = bstr.length
  const u8arr = new Uint8Array(n)
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n)
  }
  return new Blob([u8arr], { type: mime })
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  // 检查是否是从"写点什么"菜单跳转过来的
  if (route.query.contentType === 'presentation') {
    console.log('创建新讲演')
  }
})
</script>

<style scoped>
.edit-presentation-view {
  min-height: 100vh;
  background: transparent;
  padding-top: 100px; /* 防止被导航栏遮挡 */
}

.edit-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
}

.edit-header {
  text-align: center;
  margin-bottom: 40px;
}

.edit-header h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: #333;
  margin: 0 0 10px 0;
}

.edit-header p {
  font-size: 1.1rem;
  color: #666;
  margin: 0;
}

.edit-form {
  background: rgba(255, 255, 255, 0.5);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 40px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.form-section {
  margin-bottom: 40px;
}

.form-section h3 {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
  margin: 0 0 20px 0;
  padding-bottom: 10px;
  border-bottom: 2px solid #667eea;
}

.form-group {
  margin-bottom: 25px;
}

.form-group label {
  display: block;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s ease;
  background: white;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.tags-preview {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 0.9rem;
  font-weight: 500;
}

.upload-area {
  border: 2px dashed #ddd;
  border-radius: 12px;
  padding: 40px;
  text-align: center;
  transition: all 0.3s ease;
  background: white;
}

.upload-area.drag-over {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
}

.upload-placeholder {
  cursor: pointer;
}

.upload-placeholder svg {
  font-size: 3rem;
  color: #667eea;
  margin-bottom: 15px;
}

.upload-placeholder p {
  margin: 5px 0;
  color: #666;
}

.upload-hint {
  font-size: 0.9rem;
  color: #999;
}

.file-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.file-details {
  display: flex;
  align-items: center;
  gap: 15px;
}

.file-details svg {
  font-size: 2rem;
  color: #667eea;
}

.file-text p {
  margin: 0;
}

.file-name {
  font-weight: 600;
  color: #333;
}

.file-size {
  font-size: 0.9rem;
  color: #666;
}

.remove-file {
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 缩略图管理器样式 */
.thumbnail-manager {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.thumbnail-preview-section {
  flex: 1;
}

.thumbnail-preview-container {
  border: 2px dashed #ddd;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
  transition: all 0.3s ease;
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumbnail-preview-container.drag-over {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
}

.thumbnail-preview {
  position: relative;
  display: inline-block;
}

.thumbnail-preview img {
  max-width: 300px;
  max-height: 200px;
  border-radius: 8px;
  object-fit: cover;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.thumbnail-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.thumbnail-preview:hover .thumbnail-overlay {
  opacity: 1;
}

.clear-btn {
  background: #ff4757;
  color: white;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.clear-btn:hover {
  background: #ff3742;
  transform: scale(1.1);
}

.thumbnail-placeholder {
  color: #666;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.placeholder-icon {
  font-size: 3rem;
  color: #999;
}

.placeholder-text {
  font-size: 1.1rem;
  font-weight: 500;
}

.drag-hint {
  font-size: 0.9rem;
  color: #999;
}

/* 加载和错误状态 */
.thumbnail-loading,
.thumbnail-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 20px;
}

.loading-spinner {
  width: 30px;
  height: 30px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text,
.error-text {
  font-size: 0.9rem;
  color: #666;
}

.error-icon {
  color: #ff4757;
  font-size: 2rem;
}

.retry-btn {
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.retry-btn:hover {
  background: #5a6fd8;
}

/* 控制区域样式 */
.thumbnail-controls-section {
  flex: 1;
  min-width: 300px;
}

.input-group {
  margin-bottom: 20px;
}

.input-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.thumbnail-url-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s ease;
  background: white;
}

.thumbnail-url-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.thumbnail-controls {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.random-image-btn,
.upload-image-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  background: white;
  color: #333;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.random-image-btn:hover,
.upload-image-btn:hover {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
  color: #667eea;
}

.form-actions {
  display: flex;
  gap: 15px;
  justify-content: flex-end;
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #e0e0e0;
}

.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #5a6268;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .edit-container {
    padding: 20px 15px;
  }

  .edit-form {
    padding: 25px;
  }

  .edit-header h1 {
    font-size: 2rem;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}

/* 私密设置样式 */
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  margin-bottom: 10px;
}

.privacy-checkbox {
  width: 18px;
  height: 18px;
  accent-color: #667eea;
}

.checkbox-text {
  font-weight: 500;
  color: #333;
}

.privacy-hint {
  margin: 0;
  font-size: 0.9rem;
  color: #666;
  font-style: italic;
  margin-left: 28px;
}
</style>
