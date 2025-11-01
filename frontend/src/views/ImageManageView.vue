<template>
  <div class="image-manage-view">
    <NavBar />
    <div class="manage-container">
      <div class="header">
        <h1>📸 图片管理</h1>
        <div class="stats">
          <span>总计: {{ images.length }} 张</span>
          <span>大小: {{ totalSize }}</span>
        </div>
      </div>

      <div v-if="loading" class="loading">加载中...</div>

      <div v-else-if="images.length === 0" class="empty">
        <p>还没有上传任何图片哦～</p>
      </div>

      <div v-else class="image-grid">
        <div v-for="image in paginatedImages" :key="image.path" class="image-card">
          <div class="image-preview" @click="previewImage(image)">
            <img :src="getImageUrl(image.path)" :alt="image.name" />
          </div>
          <div class="image-info">
            <div class="image-name" :title="image.name">{{ image.name }}</div>
            <div class="image-meta">
              <span>{{ formatSize(image.size) }}</span>
              <span>{{ formatDate(image.uploadTime) }}</span>
            </div>
          </div>
          <div class="image-actions">
            <button class="btn-view" title="查看引用" @click="viewReferences(image)">
              <font-awesome-icon icon="eye" />
            </button>
            <button class="btn-delete" title="删除" @click="deleteImage(image)">
              <font-awesome-icon icon="trash" />
            </button>
          </div>
        </div>
      </div>

      <!-- 删除确认对话框 -->
      <Teleport to="body">
        <transition name="modal-fade">
          <div v-if="showDeleteConfirm" class="delete-confirm-overlay" @click="cancelDelete">
            <div class="delete-confirm-container" @click.stop>
              <div class="delete-confirm-header">
                <h3>确认删除</h3>
              </div>
              <div class="delete-confirm-body">
                <p>确定要删除 <strong>{{ deleteTargetImage?.name }}</strong> 吗？</p>
                <p class="warning-text">⚠️ 此操作不可撤销！</p>
              </div>
              <div class="delete-confirm-actions">
                <button class="btn-cancel" @click="cancelDelete">取消</button>
                <button class="btn-confirm-delete" @click="confirmDelete">删除</button>
              </div>
            </div>
          </div>
        </transition>
      </Teleport>

      <!-- 翻页组件 -->
      <div v-if="images.length > 0" class="pagination">
        <button
          class="page-btn"
          :disabled="currentPage === 1"
          @click="goToPage(currentPage - 1)"
        >
          上一页
        </button>
        <span class="page-info">
          第 {{ currentPage }} 页 / 共 {{ totalPages }} 页
        </span>
        <button
          class="page-btn"
          :disabled="currentPage === totalPages"
          @click="goToPage(currentPage + 1)"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- 图片预览弹窗 -->
    <div v-if="previewingImage" class="preview-modal" @click="closePreview">
      <div class="preview-content" @click.stop>
        <button class="close-btn" @click="closePreview">×</button>
        <img :src="getImageUrl(previewingImage.path)" :alt="previewingImage.name" />
        <div class="preview-info">
          <p><strong>文件名：</strong>{{ previewingImage.name }}</p>
          <p><strong>大小：</strong>{{ formatSize(previewingImage.size) }}</p>
          <p><strong>上传时间：</strong>{{ formatDate(previewingImage.uploadTime) }}</p>
          <p><strong>URL：</strong><code>{{ getImageUrl(previewingImage.path) }}</code></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeMount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import apiConfig from '@/config/api'
import { showSuccessMessage, showErrorMessage, showWarningMessage, showCustomMessage } from '@/utils/waifuMessage'
import { requestFunc } from '@/api/req'

const router = useRouter()
const route = useRoute()

const images = ref([])
const loading = ref(true)
const previewingImage = ref(null)
const showDeleteConfirm = ref(false)
const deleteTargetImage = ref(null)

// 翻页相关
const currentPage = ref(1)
const pageSize = ref(12) // 每页显示12张图片

// 翻页计算属性
const totalPages = computed(() => Math.ceil(images.value.length / pageSize.value))
const paginatedImages = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return images.value.slice(start, end)
})

// 从路由恢复状态（返回时恢复页码）
onBeforeMount(() => {
  if (route.query.page) {
    const page = parseInt(route.query.page)
    if (page >= 1) {
      currentPage.value = page
    }
  }
})

// 翻页方法
const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    // 更新URL以便返回时恢复
    router.replace({ query: { ...route.query, page } })
  }
}

// 获取图片列表
const loadImages = async () => {
  try {
    loading.value = true

    const response = await requestFunc('/images', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }, true) // 需要token

    if (!response) {
      // requestFunc返回null表示JWT过期，已由拦截器处理
      return
    }

    images.value = response.data.images || []
  } catch (error) {
    console.error('加载图片列表失败:', error)
    showErrorMessage(error)
  } finally {
    loading.value = false
  }
}

// 计算总大小
const totalSize = computed(() => {
  const total = images.value.reduce((sum, img) => sum + (img.size || 0), 0)
  return formatSize(total)
})

// 格式化文件大小
const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

// 格式化日期
const formatDate = (timestamp) => {
  if (!timestamp) return '未知'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

// 获取图片URL
const getImageUrl = (path) => {
  return apiConfig.baseURL + path
}

// 查看引用
const viewReferences = async (image) => {
  try {
    loading.value = true
    const response = await requestFunc(`/images/usages?path=${encodeURIComponent(image.path)}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }, true) // 需要token

    if (!response) {
      return
    }

    const usages = response.data.usages || []
    if (usages.length === 0) {
      showWarningMessage('该图片未被使用')
      return
    }

    // 如果有多个引用，显示第一个；如果有多个，可以后续改进为选择列表
    const usage = usages[0]

    // 保存当前状态到sessionStorage，以便返回时恢复
    sessionStorage.setItem('imageManageState', JSON.stringify({
      page: currentPage.value,
      path: route.path
    }))

    // 根据类型跳转
    if (usage.type === 'blog') {
      router.push(`/blog/${usage.id}`)
    } else if (usage.type === 'research') {
      router.push(`/research/${usage.id}`)
    } else if (usage.type === 'project') {
      router.push(`/project/${usage.id}`)
    } else if (usage.type === 'moment') {
      router.push(`/moments/${usage.id}`)
    } else if (usage.type === 'media') {
      // 书影集：统一跳转到书影集界面（/fragments/novels），并传递mediaId以便打开详情框
      router.push({
        path: '/fragments/novels',
        query: { mediaId: usage.id }
      })
    }
  } catch (error) {
    console.error('获取图片引用失败:', error)
    showErrorMessage('获取图片引用失败')
  } finally {
    loading.value = false
  }
}

// 删除图片
const deleteImage = (image) => {
  deleteTargetImage.value = image
  showDeleteConfirm.value = true
}

// 取消删除
const cancelDelete = () => {
  showDeleteConfirm.value = false
  deleteTargetImage.value = null
}

// 确认删除
const confirmDelete = async () => {
  if (!deleteTargetImage.value) return

  const image = deleteTargetImage.value
  showDeleteConfirm.value = false

  try {
    const response = await requestFunc('/images', {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      },
      data: {
        path: image.path
      }
    }, true) // 需要token

    if (!response) {
      // requestFunc返回null表示JWT过期，已由拦截器处理
      return
    }

    // 检查是否删除成功
    if (response.error === 'Image is in use') {
      // 图片正在被使用，显示详细的使用情况
      let usageMsg = '该图片正在被使用，无法删除：\n'
      if (response.usageList && Array.isArray(response.usageList)) {
        usageMsg += response.usageList.join('、')
      } else if (response.message) {
        usageMsg += response.message
      }
      showCustomMessage(usageMsg, 6000)
      showErrorMessage('该图片正在被文章或媒体使用，无法删除。请先移除相关引用后再试。')
      return
    }

    // 删除成功
    images.value = images.value.filter(img => img.path !== image.path)
    showSuccessMessage('delete')
  } catch (error) {
    console.error('删除图片失败:', error)
    // 检查是否是图片正在使用的错误
    if (error.response && error.response.data && error.response.data.error === 'Image is in use') {
      let usageMsg = '该图片正在被使用，无法删除：\n'
      if (error.response.data.usageList && Array.isArray(error.response.data.usageList)) {
        usageMsg += error.response.data.usageList.join('、')
      } else if (error.response.data.message) {
        usageMsg += error.response.data.message
      }
      showCustomMessage(usageMsg, 6000)
      showErrorMessage('该图片正在被文章或媒体使用，无法删除。请先移除相关引用后再试。')
    } else {
      showErrorMessage(error)
    }
  } finally {
    deleteTargetImage.value = null
  }
}

// 预览图片
const previewImage = (image) => {
  previewingImage.value = image
}

// 关闭预览
const closePreview = () => {
  previewingImage.value = null
}

// 恢复页面状态（从其他页面返回时）
onMounted(() => {
  loadImages()

  // 检查是否从其他页面返回，恢复页码
  const savedState = sessionStorage.getItem('imageManageState')
  if (savedState) {
    try {
      const state = JSON.parse(savedState)
      if (state.page) {
        currentPage.value = parseInt(state.page)
      }
      // 清除保存的状态
      sessionStorage.removeItem('imageManageState')
    } catch (e) {
      console.error('恢复页面状态失败:', e)
    }
  }
})
</script>

<style scoped>
.image-manage-view {
  min-height: 100vh;
  background: transparent;
  position: relative;
  padding-top: 80px; /* 防止被导航栏遮挡 */
}

.manage-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-top: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid rgba(102, 126, 234, 0.2);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-radius: 15px;
  padding: 1.5rem;
  backdrop-filter: blur(10px);
}

.header h1 {
  font-size: 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  font-weight: 700;
}

.stats {
  display: flex;
  gap: 2rem;
  font-size: 0.9rem;
  color: #667eea;
  font-weight: 600;
}

.loading,
.empty {
  text-align: center;
  padding: 4rem 2rem;
  color: #667eea;
  font-size: 1.1rem;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  border-radius: 15px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(102, 126, 234, 0.2);
  font-weight: 600;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.image-card {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(20px);
}

.image-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 20px 60px rgba(102, 126, 234, 0.2);
  background: rgba(255, 255, 255, 0.95);
}

.image-preview {
  width: 100%;
  height: 200px;
  overflow: hidden;
  cursor: pointer;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.image-preview:hover img {
  transform: scale(1.05);
}

.image-info {
  padding: 1rem;
}

.image-name {
  font-weight: 600;
  color: #333;
  margin-bottom: 0.5rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-meta {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
  color: #999;
}

.image-actions {
  display: flex;
  gap: 0.5rem;
  padding: 0 1rem 1rem;
}

.image-actions button {
  flex: 1;
  padding: 0.6rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.btn-view {
  background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(33, 150, 243, 0.3);
}

.btn-view:hover {
  background: linear-gradient(135deg, #1976D2 0%, #1565C0 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(33, 150, 243, 0.4);
}

.btn-delete {
  background: linear-gradient(135deg, #f44336 0%, #da190b 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(244, 67, 54, 0.3);
}

.btn-delete:hover {
  background: linear-gradient(135deg, #da190b 0%, #c62828 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(244, 67, 54, 0.4);
}

/* 预览弹窗 */
.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 2rem;
}

.preview-content {
  position: relative;
  max-width: 90%;
  max-height: 90%;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.close-btn {
  position: absolute;
  top: 1rem;
  right: 1rem;
  width: 40px;
  height: 40px;
  border: none;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  font-size: 2rem;
  border-radius: 50%;
  cursor: pointer;
  z-index: 10;
  transition: background 0.3s ease;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.8);
}

.preview-content img {
  max-width: 100%;
  max-height: 60vh;
  object-fit: contain;
}

.preview-info {
  padding: 1.5rem;
  background: #f9f9f9;
}

.preview-info p {
  margin: 0.5rem 0;
  color: #333;
}

.preview-info code {
  background: #e0e0e0;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
  word-break: break-all;
}

@media (max-width: 768px) {
  .manage-container {
    padding: 1rem;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .stats {
    flex-direction: column;
    gap: 0.5rem;
  }

  .image-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 1rem;
  }
}

/* 翻页组件样式 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 2rem;
  margin: 3rem 0;
  padding: 2rem;
  background: transparent;
  border-radius: 20px;
}

.page-btn {
  padding: 1rem 2rem;
  background: linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%);
  color: white;
  border: none;
  border-radius: 15px;
  cursor: pointer;
  font-size: 1.1rem;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(139, 92, 246, 0.4);
  min-width: 120px;
  position: relative;
  overflow: hidden;
}

.page-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.page-btn:hover::before {
  left: 100%;
}

.page-btn:hover:not(:disabled) {
  transform: translateY(-3px) scale(1.05);
  box-shadow: 0 8px 25px rgba(139, 92, 246, 0.6);
}

.page-btn:active:not(:disabled) {
  transform: translateY(-1px) scale(1.02);
}

.page-btn:disabled {
  background: linear-gradient(135deg, #ccc 0%, #aaa 100%);
  cursor: not-allowed;
  transform: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.page-info {
  font-size: 1.2rem;
  color: #333;
  font-weight: 600;
  padding: 1rem 2rem;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 15px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

@media (max-width: 768px) {
  .pagination {
    flex-direction: column;
    gap: 1rem;
    margin: 2rem 0;
    padding: 1.5rem;
  }

  .page-btn {
    width: 100%;
    max-width: 250px;
    padding: 1.2rem 2rem;
    font-size: 1rem;
  }

  .page-info {
    font-size: 1rem;
    padding: 0.8rem 1.5rem;
  }
}

/* 删除确认对话框样式 */
.delete-confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  padding: 2rem;
}

.delete-confirm-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 2rem;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.delete-confirm-header {
  margin-bottom: 1.5rem;
}

.delete-confirm-header h3 {
  margin: 0;
  font-size: 1.5rem;
  color: #333;
  font-weight: 600;
}

.delete-confirm-body {
  margin-bottom: 2rem;
}

.delete-confirm-body p {
  margin: 0.5rem 0;
  color: #666;
  line-height: 1.6;
}

.delete-confirm-body strong {
  color: #333;
  font-weight: 600;
}

.warning-text {
  color: #f44336 !important;
  font-weight: 600;
}

.delete-confirm-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.btn-cancel,
.btn-confirm-delete {
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-cancel {
  background: #e0e0e0;
  color: #333;
}

.btn-cancel:hover {
  background: #d0d0d0;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.btn-confirm-delete {
  background: linear-gradient(135deg, #f44336 0%, #da190b 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(244, 67, 54, 0.3);
}

.btn-confirm-delete:hover {
  background: linear-gradient(135deg, #da190b 0%, #c62828 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(244, 67, 54, 0.4);
}

/* 模态框动画 */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
