<!-- eslint-disable vue/no-v-html -->
<template>
  <div class="edit-article-view">
    <NavBar />
    <div class="edit-container">
      <div class="article-form">
        <!-- 内容类型选择 - 占据整行 -->
        <div class="form-group content-type-selector full-width">
          <label>内容类型：</label>
          <div class="type-buttons">
            <button
              type="button"
              :class="['type-btn', { active: contentType === 'article' }]"
              @click="contentType = 'article'"
            >
              📝 文章
            </button>
            <button
              type="button"
              :class="['type-btn', { active: contentType === 'media' }]"
              @click="contentType = 'media'"
            >
              🎬 媒体卡片
            </button>
          </div>
        </div>

        <!-- 文章类型和标题（仅文章模式显示） -->
        <div v-if="contentType === 'article'" class="form-row">
          <div class="form-group">
            <label>文章类型：</label>
            <select v-model="articleData.type">
              <option value="blog">博客</option>
              <option value="research">科研日记</option>
              <option value="project">项目</option>
              <option disabled>────────────</option>
              <option value="moment">随笔</option>
            </select>
          </div>
          <div class="form-group">
            <label>文章标题：</label>
            <input
              v-model="articleData.title"
              type="text"
              placeholder="请输入文章标题"
            />
          </div>
        </div>

        <!-- 媒体类型和名称（仅媒体模式显示） -->
        <div v-if="contentType === 'media'" class="form-row">
          <div class="form-group">
            <label>媒体类型：</label>
            <select v-model="mediaData.type">
              <option value="books">书单</option>
              <option value="novels">小说</option>
              <option value="movies">电影</option>
            </select>
          </div>
          <div class="form-group">
            <label>媒体名称：</label>
            <input
              v-model="mediaData.name"
              type="text"
              placeholder="请输入媒体名称"
            />
          </div>
        </div>

        <!-- 媒体卡片特有段 - 评分 -->
        <div v-if="contentType === 'media'" class="form-group">
          <label>评分(1-10):</label>
          <input v-model.number="mediaData.rating" type="number" min="1" max="10" placeholder="请输入评分" />
        </div>

        <!-- 文章特有段 - 标签和置顶选项 -->
        <div v-if="contentType === 'article'" class="form-row">
          <div class="form-group">
            <label>标签（用逗号分隔）：</label>
            <input v-model="tagsInput" type="text" placeholder="例如：技术前端,Vue" />
          </div>
          <div class="form-group">
            <div class="checkbox-container">
              <input id="isTop" v-model="articleData.isTop" type="checkbox" />
              <label for="isTop" class="checkbox-label-text">置顶文章</label>
            </div>
          </div>
        </div>

        <!-- 封面图片 - 占据整行 -->
        <div v-if="contentType === 'article'" class="form-group full-width">
          <label>封面图片：</label>
          <div class="image-manager">
            <!-- 图片预览和输入区域 -->
            <div class="image-preview-input-row">
              <!-- 图片预览区域 -->
              <div class="image-preview-section">
                <div
                  class="image-preview-container"
                  :class="{ 'drag-over': isDragOver }"
                  @dragover.prevent="handleDragOver"
                  @dragleave.prevent="handleDragLeave"
                  @drop.prevent="handleDrop"
                >
                  <div v-if="articleData.image" class="image-preview">
                    <!-- 加载状态指示器 -->
                    <div v-if="currentImageLoading" class="image-loading">
                      <div class="loading-spinner"></div>
                      <span class="loading-text">图片加载中...</span>
                    </div>
                    <!-- 错误状态指示器 -->
                    <div v-else-if="currentImageError" class="image-error">
                      <font-awesome-icon icon="exclamation-triangle" class="error-icon" />
                      <span class="error-text">图片加载失败</span>
                      <button type="button" class="retry-btn" @click="retryImageLoad">
                        <font-awesome-icon icon="redo" />
                        重试
                      </button>
                    </div>
                    <!-- 正常图片显示 -->
                    <img
                      v-else
                      :src="articleData.image"
                      alt="封面预览"
                      @error="handleImageError"
                      @load="handleImageLoad"
                    />
                    <div class="image-overlay">
                      <button type="button" class="clear-btn" @click="articleData.image = ''">
                        <font-awesome-icon icon="trash" />
                      </button>
                    </div>
                  </div>
                  <div v-else class="image-placeholder">
                    <font-awesome-icon icon="image" class="placeholder-icon" />
                    <span class="placeholder-text">暂无封面图片</span>
                    <div class="drag-hint">拖拽图片到此处上传</div>
                  </div>
                </div>
              </div>

              <!-- 图片链接输入区域 -->
              <div class="image-input-section">
                <div class="input-group">
                  <label class="input-label">
                    <font-awesome-icon icon="link" />
                    封面链接
                  </label>
                  <input
                    v-model="articleData.image"
                    type="text"
                    placeholder="请输入图片链接"
                    class="image-url-input"
                    @input="handleImageUrlInput"
                  />
                </div>
                <div class="image-upload-controls">
                  <button type="button" class="random-image-btn" @click="getRandomImage">
                    <font-awesome-icon icon="dice" />
                    随机图片
                  </button>
                  <button type="button" class="upload-image-btn" @click="triggerFileUpload">
                    <font-awesome-icon icon="upload" />
                    上传图片
                  </button>
                  <input
                    ref="fileInput"
                    type="file"
                    accept="image/*"
                    style="display: none"
                    @change="handleFileUpload"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="contentType === 'media'" class="form-group full-width">
          <label>媒体封面：</label>
          <div class="image-preview-input-row">
            <!-- 图片预览区域 -->
            <div class="image-preview-section">
              <div
                class="image-preview-container media-preview-container"
                :class="{ 'drag-over': isDragOver }"
                @dragover.prevent="handleDragOver"
                @dragleave.prevent="handleDragLeave"
                @drop.prevent="handleDrop"
              >
                <div v-if="mediaData.image" class="image-preview">
                  <!-- 加载状态指示器 -->
                  <div v-if="currentImageLoading" class="image-loading">
                    <div class="loading-spinner"></div>
                    <span class="loading-text">图片加载中...</span>
                  </div>
                  <!-- 错误状态指示器 -->
                  <div v-else-if="currentImageError" class="image-error">
                    <font-awesome-icon icon="exclamation-triangle" class="error-icon" />
                    <span class="error-text">图片加载失败</span>
                    <button type="button" class="retry-btn" @click="retryImageLoad">
                      <font-awesome-icon icon="redo" />
                      重试
                    </button>
                  </div>
                  <!-- 正常图片显示 -->
                  <img
                    v-else
                    :src="mediaData.image"
                    alt="封面预览"
                    @error="handleImageError"
                    @load="handleImageLoad"
                  />
                  <div class="image-overlay">
                    <button type="button" class="clear-btn" @click="mediaData.image = ''">
                      <font-awesome-icon icon="trash" />
                    </button>
                  </div>
                </div>
                <div v-else class="image-placeholder">
                  <font-awesome-icon icon="image" class="placeholder-icon" />
                  <span class="placeholder-text">暂无封面图片</span>
                  <div class="drag-hint">拖拽图片到此处上传</div>
                </div>
              </div>
            </div>

            <!-- 图片链接输入区域 -->
            <div class="image-input-section">
              <div class="input-group">
                <label class="input-label">
                  <font-awesome-icon icon="link" />
                  封面链接
                </label>
                <input
                  v-model="mediaData.image"
                  type="text"
                  placeholder="请输入图片链接"
                  class="image-url-input"
                  @input="handleImageUrlInput"
                />
              </div>
              <div class="image-upload-controls">
                <button type="button" class="random-image-btn" @click="getRandomImage">
                  <font-awesome-icon icon="dice" />
                  随机图片
                </button>
                <button type="button" class="upload-image-btn" @click="triggerFileUpload">
                  <font-awesome-icon icon="upload" />
                  上传图片
                </button>
                <input
                  ref="fileInput"
                  type="file"
                  accept="image/*"
                  style="display: none"
                  @change="handleFileUpload"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 图片选择器-->
        <div v-if="showImagePicker" class="image-picker">
          <h4>选择预设图片</h4>
          <div class="image-grid">
            <div
              v-for="img in presetImages"
              :key="img.id"
              class="image-option"
              :class="{ selected: (contentType === 'article' ? articleData.image : mediaData.image) === img.url }"
              @click="selectImage(img.url)"
            >
              <img :src="img.url" :alt="img.name" />
              <span class="image-name">{{ img.name }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Typora风格的所见即所得编辑器 -->
    <div v-if="contentType === 'article'" class="typora-editor">
      <div class="editor-content">
        <!-- 左右分屏编辑器 -->
        <div class="split-editor">
          <!-- 左侧：Markdown源码编辑器 -->
          <div class="editor-panel left-panel">
            <div class="panel-header">
              <span class="panel-title">Markdown 源码</span>
            </div>

            <!-- 左侧工具栏 -->
            <div class="panel-toolbar">
              <div class="toolbar-group">
                <button title="粗体 (Ctrl+B)" class="toolbar-btn" @click="execWysiwygCommand('bold')">
                  <font-awesome-icon icon="bold" />
                </button>
                <button title="斜体 (Ctrl+I)" class="toolbar-btn" @click="execWysiwygCommand('italic')">
                  <font-awesome-icon icon="italic" />
                </button>
                <button title="删除" class="toolbar-btn" @click="execWysiwygCommand('strikeThrough')">
                  <font-awesome-icon icon="strikethrough" />
                </button>
              </div>
              <div class="toolbar-divider"></div>
              <div class="toolbar-group">
                <button title="标题 1" class="toolbar-btn" @click="execHeadingCommand(1)">
                  H1
                </button>
                <button title="标题 2" class="toolbar-btn" @click="execHeadingCommand(2)">
                  H2
                </button>
                <button title="标题 3" class="toolbar-btn" @click="execHeadingCommand(3)">
                  H3
                </button>
              </div>
              <div class="toolbar-divider"></div>
              <div class="toolbar-group">
                <button title="链接 (Ctrl+K)" class="toolbar-btn" @click="execLinkCommand()">
                  <font-awesome-icon icon="link" />
                </button>
                <button title="图片" class="toolbar-btn" @click="execImageCommand()">
                  <font-awesome-icon icon="image" />
                </button>
                <button title="代码" class="toolbar-btn" @click="execCodeCommand()">
                  <font-awesome-icon icon="code" />
                </button>
              </div>
              <div class="toolbar-divider"></div>
              <div class="toolbar-group">
                <button title="引用" class="toolbar-btn" @click="insertMarkdown('> ', '')">
                  <font-awesome-icon icon="quote-left" />
                </button>
                <button title="无序列表" class="toolbar-btn" @click="insertMarkdown('- ', '')">
                  <font-awesome-icon icon="list-ul" />
                </button>
                <button title="有序列表" class="toolbar-btn" @click="insertMarkdown('1. ', '')">
                  <font-awesome-icon icon="list-ol" />
                </button>
              </div>
              <div class="toolbar-divider"></div>
              <div class="toolbar-group">
                <button title="分割线" class="toolbar-btn" @click="insertMarkdown('---\n', '')">
                  <font-awesome-icon icon="minus" />
                </button>
                <button title="表格" class="toolbar-btn" @click="insertTable">
                  <font-awesome-icon icon="table" />
                </button>
              </div>
            </div>

            <div class="panel-content">
              <textarea
                ref="editorTextarea"
                v-model="markdownContent"
                placeholder="# 开始写作吧...&#10;&#10;支持 Markdown 语法：&#10;- **粗体** *斜体*&#10;- # 标题&#10;- [链接](url)&#10;- `代码`&#10;- > 引用&#10;&#10;💡 提示：可以直接Ctrl+V 粘贴图片"
                @input="updatePreview"
                @keydown.tab.prevent="handleTab"
                @paste="handlePaste"
                @keyup="updateCursorPosition"
                @click="updateCursorPosition"
                @scroll="syncScroll"
              ></textarea>
            </div>

            <!-- 左侧统计信息 -->
            <div class="panel-footer">
              <div class="editor-stats">
                <span class="word-count">{{ wordCount }} 字</span>
                <span class="line-count">{{ lineCount }} 行</span>
                <span class="cursor-position">{{ cursorLine }} {{ cursorColumn }}</span>
              </div>
            </div>
          </div>

          <!-- 右侧：渲染预览 -->
          <div class="editor-panel right-panel">
            <div class="panel-header">
              <span class="panel-title">预览效果</span>
            </div>
            <div class="panel-content">
              <div
                class="preview-content markdown-body"
                v-html="sanitizedContent"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 隐藏的文件输入-->
      <input
        ref="imageInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="handleImageUpload"
      />
    </div>

    <!-- 媒体编辑器隐藏的文件输入 -->
    <input
      v-if="contentType === 'media'"
      ref="mediaImageInput"
      type="file"
      accept="image/*"
      style="display: none"
      @change="handleMediaImageUpload"
    />

    <!-- 媒体卡片编辑 - 使用与文章编辑器相同的分屏样式 -->
    <div v-if="contentType === 'media'" class="typora-editor">
      <div class="editor-content">
        <!-- 左右分屏编辑器 -->
        <div class="split-editor">
          <!-- 左侧：Markdown源码编辑器 -->
          <div class="editor-panel left-panel">
            <div class="panel-header">
              <span class="panel-title">媒体评价</span>
            </div>

            <!-- 左侧工具栏 -->
            <div class="panel-toolbar">
              <div class="toolbar-group">
                <button title="粗体" class="toolbar-btn" @click="insertMediaMarkdown('**', '**')">
                  <font-awesome-icon icon="bold" />
                </button>
                <button title="斜体" class="toolbar-btn" @click="insertMediaMarkdown('*', '*')">
                  <font-awesome-icon icon="italic" />
                </button>
                <button title="删除线" class="toolbar-btn" @click="insertMediaMarkdown('~~', '~~')">
                  <font-awesome-icon icon="strikethrough" />
                </button>
              </div>
              <div class="toolbar-divider"></div>
              <div class="toolbar-group">
                <button title="标题 1" class="toolbar-btn" @click="insertMediaMarkdown('# ', '')">
                  H1
                </button>
                <button title="标题 2" class="toolbar-btn" @click="insertMediaMarkdown('## ', '')">
                  H2
                </button>
                <button title="标题 3" class="toolbar-btn" @click="insertMediaMarkdown('### ', '')">
                  H3
                </button>
              </div>
              <div class="toolbar-divider"></div>
              <div class="toolbar-group">
                <button title="链接" class="toolbar-btn" @click="insertMediaLink()">
                  <font-awesome-icon icon="link" />
                </button>
                <button title="图片" class="toolbar-btn" @click="insertMediaImage()">
                  <font-awesome-icon icon="image" />
                </button>
                <button title="代码" class="toolbar-btn" @click="insertMediaMarkdown('`', '`')">
                  <font-awesome-icon icon="code" />
                </button>
              </div>
              <div class="toolbar-divider"></div>
              <div class="toolbar-group">
                <button title="引用" class="toolbar-btn" @click="insertMediaMarkdown('> ', '')">
                  <font-awesome-icon icon="quote-left" />
                </button>
                <button title="无序列表" class="toolbar-btn" @click="insertMediaMarkdown('- ', '')">
                  <font-awesome-icon icon="list-ul" />
                </button>
                <button title="有序列表" class="toolbar-btn" @click="insertMediaMarkdown('1. ', '')">
                  <font-awesome-icon icon="list-ol" />
                </button>
              </div>
              <div class="toolbar-divider"></div>
              <div class="toolbar-group">
                <button title="分割线" class="toolbar-btn" @click="insertMediaMarkdown('---\n', '')">
                  <font-awesome-icon icon="minus" />
                </button>
                <button title="表格" class="toolbar-btn" @click="insertMediaTable()">
                  <font-awesome-icon icon="table" />
                </button>
              </div>
            </div>

            <!-- 文本编辑器 -->
            <textarea
              ref="mediaTextarea"
              v-model="mediaData.description"
              placeholder="# 写下你的评价...&#10;&#10;支持 Markdown 语法：&#10;- **粗体** *斜体* ~~删除线~~&#10;- # 标题 1-3&#10;- [链接](url)&#10;- ![图片](url)&#10;- `代码`&#10;- > 引用&#10;- - 无序列表 1. 有序列表&#10;- --- 分割线&#10;- 表格&#10;&#10;💡 提示：可以直接Ctrl+V 粘贴图片"
              class="source-editor"
              @keydown.tab.prevent="handleMediaTab"
              @paste="handleMediaPaste"
              @input="updateMediaPreview"
              @keyup="updateMediaCursorPosition"
              @click="updateMediaCursorPosition"
              @scroll="syncMediaScroll"
            ></textarea>

            <!-- 左侧底部统计 -->
            <div class="panel-footer">
              <div class="editor-stats">
                <span class="word-count">{{ mediaWordCount }} 字</span>
                <span class="line-count">{{ mediaLineCount }} 行</span>
                <span class="cursor-position">{{ mediaCursorLine }} {{ mediaCursorColumn }}</span>
              </div>
            </div>
          </div>

          <!-- 右侧：实时预览 -->
          <div class="editor-panel right-panel">
            <div class="panel-header">
              <span class="panel-title">预览效果</span>
            </div>

            <!-- 预览内容 -->
            <div class="preview-content markdown-body" v-html="sanitizedMediaContent"></div>
          </div>
        </div>
      </div>
    </div>

    <div class="action-buttons">
      <button class="save-btn" :disabled="!canSave" @click="contentType === 'article' ? saveArticle() : saveMedia()">
        {{ isEditing ? (contentType === 'article' ? '更新文章' : '更新卡片') : (contentType === 'article' ? '发布文章' : '创建卡片') }}
      </button>
      <button v-if="hasChanges" class="discard-btn" @click="discardChanges">放弃更改</button>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, computed, onBeforeUnmount, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import axios from 'axios'
import DOMPurify from 'dompurify' // 防止XSS攻击
import { marked } from 'marked'
import hljs from 'highlight.js'
import '@/assets/styles/github-highlight.css'
import '@/assets/styles/github-markdown.css'
import NavBar from '@/components/NavBar.vue'
import { createArticle, updateArticle } from '@/api/Articles/edit'
import { getJWT, requestFunc } from '@/api/req'
import { getArticleByID } from '@/api/Articles/browse'
import { getMediaByID } from '@/api/media/browse'
// import { imageConfig } from '@/config/images'
import apiConfig from '@/config/api'
import { showErrorMessage, showSuccessMessage, showWarningMessage, showCustomMessage } from '@/utils/waifuMessage'

const route = useRoute()
const router = useRouter()
const store = useStore()

const contentType = ref('article') // 'article' 'media'
const markdownContent = ref('')
const renderedContent = ref('')
const sanitizedContent = ref('')
const tagsInput = ref('')
const showImagePicker = ref(false)
// const editorMode = ref('wysiwyg') // 已移除，现在使用分屏模式
const cursorLine = ref(1) // 光标行号
const cursorColumn = ref(1) // 光标列号
const editorTextarea = ref(null)
const imageInput = ref(null)
const wysiwygContent = ref(null)
const mediaTextarea = ref(null)
const mediaCursorLine = ref(1) // 媒体编辑器光标行号
const mediaCursorColumn = ref(1) // 媒体编辑器光标列号
const mediaImageInput = ref(null) // 媒体编辑器图片上传input

// 图片加载状态管理
const imageLoadingStates = ref(new Map())
const imageCache = ref(new Map())
const imageRetryCount = ref(new Map())

// 待上传文件队列管理
const pendingUploads = ref(new Map()) // key: localUrl, value: { file, localUrl, serverUrl }

// 封面图片上传队列
const pendingCoverUploads = ref(new Map()) // key: localUrl, value: { file, localUrl, serverUrl }

// 移除API选择相关代码，改为后端随机选择

// 文件上传相关
const fileInput = ref(null)
const isDragOver = ref(false)

const articleData = ref({
  title: '',
  type: 'blog',
  image: '',
  tags: [],
  isTop: false
})

const mediaData = ref({
  name: '',
  type: 'books',
  image: '',
  rating: 8,
  description: ''
})

// 从配置文件加载预设图片
const presetImages = ref([])

const isEditing = computed(() => !!route.params.id)

// 原始数据，用于检测是否有更改
const originalArticleData = ref(null)
const originalMediaData = ref(null)
const originalMarkdownContent = ref('')
const originalTagsInput = ref('')

const canSave = computed(() => {
  if (contentType.value === 'article') {
    return articleData.value.title.trim() && markdownContent.value.trim()
  } else {
    return mediaData.value.name.trim() && mediaData.value.rating >= 1 && mediaData.value.rating <= 10
  }
})

// 检测是否有更改
const hasChanges = computed(() => {
  if (contentType.value === 'article') {
    if (!originalArticleData.value) return false

    return (
      articleData.value.title !== originalArticleData.value.title ||
      articleData.value.image !== originalArticleData.value.image ||
      articleData.value.type !== originalArticleData.value.type ||
      articleData.value.isTop !== originalArticleData.value.isTop ||
      markdownContent.value !== originalMarkdownContent.value ||
      tagsInput.value !== originalTagsInput.value
    )
  } else {
    if (!originalMediaData.value) return false

    return (
      mediaData.value.name !== originalMediaData.value.name ||
      mediaData.value.image !== originalMediaData.value.image ||
      mediaData.value.type !== originalMediaData.value.type ||
      mediaData.value.rating !== originalMediaData.value.rating ||
      mediaData.value.description !== originalMediaData.value.description
    )
  }
})

// 字数统计
const wordCount = computed(() => {
  return markdownContent.value.replace(/\s/g, '').length
})

// 行数统计
const lineCount = computed(() => {
  return markdownContent.value.split('\n').length
})

// 媒体编辑器字数统计
const mediaWordCount = computed(() => {
  return mediaData.value.description.replace(/\s/g, '').length
})

// 媒体编辑器行数统计
const mediaLineCount = computed(() => {
  return mediaData.value.description.split('\n').length
})

// 媒体内容预览（Markdown 渲染）- 与文章编辑器保持一致的功能
const sanitizedMediaContent = computed(() => {
  if (!mediaData.value.description) return ''

  // 使用 marked 渲染 Markdown，确保代码块结构正确
  const rendered = marked(mediaData.value.description, {
    breaks: true, // 将换行符转换为 <br>
    gfm: true, // 启用 GitHub Flavored Markdown
    headerIds: false,
    mangle: false,
    sanitize: false // 不禁用HTML标签，让br标签通过
  })

  // 处理blob URL，尝试显示实际图片
  let processedContent = rendered
  const blobImageRegex = /<img[^>]*src="(blob:[^"]*)"[^>]*>/g
  processedContent = processedContent.replace(blobImageRegex, (match, blobUrl) => {
    const altMatch = match.match(/alt="([^"]*)"/)
    const altText = altMatch ? altMatch[1] : '粘贴的图片'
    return `<div class="blob-image-container">
      <img src="${blobUrl}" alt="${altText}" class="blob-image"
           onerror="this.style.display='none'; this.nextElementSibling.style.display='flex';" />
      <div class="image-placeholder-preview" style="display: none;">
        <div class="placeholder-icon">📷</div>
        <div class="placeholder-text">${altText}</div>
        <div class="placeholder-hint">将在保存时上传</div>
      </div>
    </div>`
  })

  // 使用与文章编辑器相同的 DOMPurify 配置
  const sanitized = DOMPurify.sanitize(processedContent, {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'ul', 'ol', 'li', 'blockquote', 'pre', 'code', 'span', 'div', 'table', 'thead', 'tbody', 'tr', 'th', 'td', 'img'],
    ALLOWED_ATTR: ['class', 'style', 'src', 'alt', 'title', 'width', 'height'],
    ALLOW_DATA_ATTR: true,
    ALLOW_UNKNOWN_PROTOCOLS: true
  })

  // 确保DOM更新后高亮代码块
  nextTick(() => {
    setTimeout(() => {
      const finalCodeBlocks = document.querySelectorAll('.preview-content.markdown-body pre code')
      finalCodeBlocks.forEach((block) => {
        if (block.textContent && block.textContent.trim()) {
          hljs.highlightElement(block)
        }
      })
    }, 100)
  })

  return sanitized
})

// 图片加载状态计算属性
const currentImageLoading = computed(() => {
  const currentImage = contentType.value === 'article' ? articleData.value.image : mediaData.value.image
  if (!currentImage) return false

  const loading = imageLoadingStates.value.get(currentImage) || false
  // console.log('当前图片加载状态', { currentImage, loading })
  return loading
})

const currentImageError = computed(() => {
  const currentImage = contentType.value === 'article' ? articleData.value.image : mediaData.value.image
  if (!currentImage) return false

  // 只有在重试次数超3次且图片确实加载失败时才显示错误
  const retryCount = imageRetryCount.value.get(currentImage) || 0
  const isLoading = imageLoadingStates.value.get(currentImage) || false
  const hasError = retryCount > 3 && !isLoading

  // console.log('当前图片错误状态', { currentImage, retryCount, isLoading, hasError })
  return hasError
})

// 更新预览并进行高亮代码块
const updatePreview = () => {
  // 使用 marked 渲染 Markdown，确保代码块结构正确
  renderedContent.value = marked(markdownContent.value, {
    breaks: true, // 将换行符转换为 <br>
    gfm: true, // 启用 GitHub Flavored Markdown
    headerIds: false,
    mangle: false,
    sanitize: false // 不禁用HTML标签，让br标签通过
  })

  // 在渲染后在HTML 中添加内联样式强制左对齐
  const tempDiv = document.createElement('div')
  tempDiv.innerHTML = renderedContent.value

  // 为所有块级元素添加text-align: left，但排除标题
  const blockElements = tempDiv.querySelectorAll('p, div, ul, ol, li, blockquote, pre')
  blockElements.forEach(el => {
    el.style.textAlign = 'left'
    el.style.textAlignLast = 'left'
  })

  // 只对标题设置左对齐，不设置其他样式
  const headingElements = tempDiv.querySelectorAll('h1, h2, h3, h4, h5, h6')
  headingElements.forEach(el => {
    el.style.textAlign = 'left'
    el.style.textAlignLast = 'left'
    // 确保不覆盖字体大小
    el.style.fontSize = ''
    el.style.fontWeight = ''
  })

  // 确保代码块有正确的结构
  const codeBlocks = tempDiv.querySelectorAll('pre code')
  // console.log('找到的代码块数量:', codeBlocks.length) // 调试
  // console.log('Marked 渲染的完整HTML:', renderedContent.value) // 调试
  codeBlocks.forEach(block => {
    // console.log('代码块内', block.textContent) // 调试
    // console.log('代码块HTML:', block.outerHTML) // 调试
    // 确保代码块是块级元素
    block.style.display = 'block'
    block.style.whiteSpace = 'pre'
    block.style.overflow = 'visible'
    block.style.wordBreak = 'normal'
  })

  renderedContent.value = tempDiv.innerHTML

  // 处理blob URL，尝试显示实际图片，失败时显示占位符
  let processedContent = renderedContent.value

  // 查找所有blob URL的图片标签
  const blobImageRegex = /<img[^>]*src="(blob:[^"]*)"[^>]*>/g
  processedContent = processedContent.replace(blobImageRegex, (match, blobUrl) => {
    // 提取alt文本
    const altMatch = match.match(/alt="([^"]*)"/)
    const altText = altMatch ? altMatch[1] : '粘贴的图片'

    // 尝试显示实际图片，添加错误处理
    return `<div class="blob-image-container">
      <img src="${blobUrl}" alt="${altText}" class="blob-image" 
           onerror="this.style.display='none'; this.nextElementSibling.style.display='flex';" />
      <div class="image-placeholder-preview" style="display: none;">
        <div class="placeholder-icon">📷</div>
        <div class="placeholder-text">${altText}</div>
        <div class="placeholder-hint">将在保存时上传</div>
      </div>
    </div>`
  })

  // 使用更宽松的 DOMPurify 配置，确保代码块、br标签和图片不被过滤
  sanitizedContent.value = DOMPurify.sanitize(processedContent, {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'u', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'ul', 'ol', 'li', 'blockquote', 'pre', 'code', 'span', 'div', 'table', 'thead', 'tbody', 'tr', 'th', 'td', 'img'],
    ALLOWED_ATTR: ['class', 'style', 'src', 'alt', 'title', 'width', 'height'],
    ALLOW_DATA_ATTR: true,
    ALLOW_UNKNOWN_PROTOCOLS: true
  })

  // 确保DOM 更新后高亮代码块（不添加复制按钮）
  nextTick(() => {
    // 延迟执行，确保代码块完全渲染
    setTimeout(() => {
      const finalCodeBlocks = document.querySelectorAll('pre code')
      // console.log('DOM 中的代码块数量', finalCodeBlocks.length) // 调试

      // 移除JavaScript样式处理，完全依赖CSS
      finalCodeBlocks.forEach((block) => {
        // console.log('DOM 代码块内容', block.textContent) // 调试
        // 确保代码块没有被破坏
        if (block.textContent && block.textContent.trim()) {
          hljs.highlightElement(block)
        }
      })
    }, 100) // 延迟100ms确保渲染完成
  })
}

// 媒体编辑器插入Markdown
const insertMediaMarkdown = (before, after) => {
  const textarea = mediaTextarea.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = mediaData.value.description.substring(start, end)
  const beforeText = mediaData.value.description.substring(0, start)
  const afterText = mediaData.value.description.substring(end)

  mediaData.value.description = beforeText + before + selectedText + after + afterText

  nextTick(() => {
    textarea.focus()
    textarea.setSelectionRange(
      start + before.length,
      start + before.length + selectedText.length
    )
  })
}

// 媒体编辑Tab 键处理
const handleMediaTab = (event) => {
  const textarea = mediaTextarea.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd

  // 插入两个空格
  mediaData.value.description =
    mediaData.value.description.substring(0, start) +
    '  ' +
    mediaData.value.description.substring(end)

  // 移动光标
  nextTick(() => {
    textarea.selectionStart = textarea.selectionEnd = start + 2
    updateMediaCursorPosition()
  })
}

// 媒体编辑器更新预览
const updateMediaPreview = () => {
  // sanitizedMediaContent 是 computed，会自动更新
  // 但需要确保代码高亮
  nextTick(() => {
    setTimeout(() => {
      const finalCodeBlocks = document.querySelectorAll('.preview-content.markdown-body pre code')
      finalCodeBlocks.forEach((block) => {
        if (block.textContent && block.textContent.trim()) {
          hljs.highlightElement(block)
        }
      })
    }, 100)
  })
}

// 媒体编辑器光标位置更新
const updateMediaCursorPosition = () => {
  if (!mediaTextarea.value) return

  const textarea = mediaTextarea.value
  const text = textarea.value
  const cursorPos = textarea.selectionStart

  // 计算行号和列号
  const textBeforeCursor = text.substring(0, cursorPos)
  const lines = textBeforeCursor.split('\n')
  mediaCursorLine.value = lines.length
  mediaCursorColumn.value = lines[lines.length - 1].length + 1
}

// 媒体编辑器滚动同步
const syncMediaScroll = () => {
  if (!mediaTextarea.value) return

  const textarea = mediaTextarea.value
  const previewContent = document.querySelector('.preview-content.markdown-body')

  if (!previewContent) return

  // 计算滚动比例
  const scrollTop = textarea.scrollTop
  const scrollHeight = textarea.scrollHeight
  const clientHeight = textarea.clientHeight
  const maxScroll = scrollHeight - clientHeight

  if (maxScroll > 0) {
    const scrollRatio = scrollTop / maxScroll
    const previewMaxScroll = previewContent.scrollHeight - previewContent.clientHeight
    const targetScrollTop = scrollRatio * previewMaxScroll

    previewContent.scrollTop = targetScrollTop
  }
}

// 媒体编辑器粘贴处理（Ctrl+V 粘贴图片）
const handleMediaPaste = async (event) => {
  const items = event.clipboardData.items
  if (!items) return

  // 查找图片
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item.type.indexOf('image') !== -1) {
      // 阻止默认粘贴行为
      event.preventDefault()

      // 获取图片文件
      const file = item.getAsFile()
      if (!file) continue

      const textarea = mediaTextarea.value
      const start = textarea.selectionStart

      try {
        // 创建本地预览URL
        const localUrl = URL.createObjectURL(file)

        // 生成文件名
        const timestamp = new Date().getTime()
        const filename = `paste_${timestamp}.png`

        // 插入本地预览的Markdown
        const imageMarkdown = `\n\n![${filename}](${localUrl})\n\n`
        const beforeText = mediaData.value.description.substring(0, start)
        const afterText = mediaData.value.description.substring(start)
        mediaData.value.description = beforeText + imageMarkdown + afterText

        // 添加到待上传队列
        pendingUploads.value.set(localUrl, {
          file,
          localUrl,
          serverUrl: null
        })

        // 更新预览
        updateMediaPreview()

        showCustomMessage('图片已粘贴，将在保存时上传', 3000)
      } catch (error) {
        showErrorMessage('paste_failed')
      }

      break
    }
  }
}

// 媒体编辑器插入链接
const insertMediaLink = () => {
  const url = prompt('请输入链接地址:')
  if (url) {
    insertMediaMarkdown('[', `](${url})`)
  }
}

// 媒体编辑器插入图片
const insertMediaImage = () => {
  if (mediaImageInput.value) {
    mediaImageInput.value.click()
  }
}

// 媒体编辑器插入表格
const insertMediaTable = () => {
  const tableMarkdown = '\n| 列1 | 列2 | 列3 |\n|-----|-----|-----|\n| 内容1 | 内容2 | 内容3 |\n| 内容4 | 内容5 | 内容6 |\n\n'
  insertMediaMarkdown(tableMarkdown, '')
}

// 媒体编辑器图片上传处理
const handleMediaImageUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const textarea = mediaTextarea.value
  const start = textarea.selectionStart

  try {
    // 显示上传中提示
    const uploadingText = '\n\n![上传中...]()\n\n'
    const beforeText = mediaData.value.description.substring(0, start)
    const afterText = mediaData.value.description.substring(start)
    mediaData.value.description = beforeText + uploadingText + afterText

    // 上传图片
    const imageUrl = await uploadImageFile(file)
    if (!imageUrl) {
      // 上传失败，移除上传中文本
      mediaData.value.description = mediaData.value.description.replace(uploadingText, '')
      return
    }

    // 替换上传中文本为实际图片
    mediaData.value.description = mediaData.value.description.replace(
      uploadingText,
      `\n\n![${file.name}](${imageUrl})\n\n`
    )

    // 更新预览
    updateMediaPreview()

    showSuccessMessage('upload')
  } catch (error) {
    // 移除上传中文件
    mediaData.value.description = mediaData.value.description.replace('\n\n![上传中...]()\n\n', '')
    showErrorMessage('upload_failed')
  }

  // 清空 input
  event.target.value = ''
}

// 保存文章
const saveArticle = async () => {
  if (!canSave.value) return

  try {
    const user = store.state.user
    const token = getJWT()
    if (!user || !user.isLogged || !token) {
      showErrorMessage('401')
      return
    }

    // 批量上传待上传的图片
    if (pendingUploads.value.size > 0) {
      // console.log(`开始批量上传 ${pendingUploads.value.size} 个图片...`)

      for (const [localUrl, uploadInfo] of pendingUploads.value) {
        try {
          const serverUrl = await uploadImageFile(uploadInfo.file)
          if (serverUrl) {
            // 替换Markdown中的本地URL为服务器URL
            markdownContent.value = markdownContent.value.replace(localUrl, serverUrl)
            uploadInfo.serverUrl = serverUrl
            // console.log(`图片上传成功: ${localUrl} -> ${serverUrl}`)
          } else {
            // console.error(`图片上传失败: ${localUrl}`)
          }
        } catch (error) {
          // console.error(`图片上传出错: ${localUrl}`, error)
        }
      }

      // 清理本地URL
      for (const [localUrl] of pendingUploads.value) {
        URL.revokeObjectURL(localUrl)
      }
      pendingUploads.value.clear()

      // console.log('批量上传完成')
    }

    // 批量上传封面图片
    if (pendingCoverUploads.value.size > 0) {
      // console.log(`开始批量上传 ${pendingCoverUploads.value.size} 个封面图片...`)

      for (const [, uploadInfo] of pendingCoverUploads.value) {
        try {
          const serverUrl = await uploadImageFile(uploadInfo.file)
          if (serverUrl) {
            // 替换封面图片的本地URL为服务器URL
            if (contentType.value === 'article') {
              articleData.value.image = serverUrl
            } else {
              mediaData.value.image = serverUrl
            }
            uploadInfo.serverUrl = serverUrl
            // console.log(`封面图片上传成功: ${localUrl} -> ${serverUrl}`)
          } else {
            // console.error(`封面图片上传失败: ${localUrl}`)
          }
        } catch (error) {
          // console.error(`封面图片上传出错: ${localUrl}`, error)
        }
      }

      // 清理本地URL
      for (const [localUrl] of pendingCoverUploads.value) {
        URL.revokeObjectURL(localUrl)
      }
      pendingCoverUploads.value.clear()

      // console.log('封面图片批量上传完成')
    }

    // 处理标签 - 确保发送正确的JSON格式
    const tagsArray = tagsInput.value.split(',').map(tag => tag.trim()).filter(tag => tag)
    articleData.value.tags = tagsArray

    // 碎碎念使用不同的 API
    if (articleData.value.type === 'moment') {
      const { createMoment, updateMoment } = await import('@/api/Moments/edit')
      if (isEditing.value) {
        await updateMoment(user, route.params.id, articleData.value.title, markdownContent.value, articleData.value.image)
        showSuccessMessage('update')
      } else {
        await createMoment(user, articleData.value.title, markdownContent.value, articleData.value.image, user.username)
        showSuccessMessage('submit')
      }
      router.push('/moments')
      return
    }

    // 根据文章类型设置内容
    const articlePayload = {
      title: articleData.value.title,
      image: articleData.value.image,
      tags: articleData.value.tags
    }

    // 所有类型都使用content字段
    articlePayload.content = markdownContent.value

    if (isEditing.value) {
      // console.log('正在更新文章...', { token, articlePayload, type: articleData.value.type, id: route.params.id })
      await updateArticle(token, articlePayload, articleData.value.type, route.params.id)
      showSuccessMessage('update')
    } else {
      // console.log('正在创建文章...', { token, articlePayload, type: articleData.value.type })
      await createArticle(token, articlePayload, articleData.value.type)
      showSuccessMessage('submit')
    }

    // 跳转到对应页面
    router.push(`/${articleData.value.type}`)
  } catch (error) {
    showErrorMessage(error)
  }
}

// 保存媒体卡片
const saveMedia = async () => {
  if (!canSave.value) {
    showErrorMessage('empty_input')
    return
  }

  try {
    const user = store.state.user
    const token = getJWT()
    if (!user || !user.isLogged || !token) {
      showErrorMessage('401')
      return
    }

    // 批量上传待上传的图片（正文中的图片）
    if (pendingUploads.value.size > 0) {
      for (const [localUrl, uploadInfo] of pendingUploads.value) {
        try {
          const serverUrl = await uploadImageFile(uploadInfo.file)
          if (serverUrl) {
            // 替换Markdown中的本地URL为服务器URL
            mediaData.value.description = mediaData.value.description.replace(localUrl, serverUrl)
            uploadInfo.serverUrl = serverUrl
          }
        } catch (error) {
          // 忽略单个图片上传失败，但记录警告
          console.warn('正文图片上传失败:', error)
        }
      }

      // 清理本地URL
      for (const [localUrl] of pendingUploads.value) {
        URL.revokeObjectURL(localUrl)
      }
      pendingUploads.value.clear()
    }

    // 批量上传封面图片
    if (pendingCoverUploads.value.size > 0) {
      for (const [, uploadInfo] of pendingCoverUploads.value) {
        try {
          const serverUrl = await uploadImageFile(uploadInfo.file)
          if (serverUrl) {
            // 替换封面图片的本地URL为服务器URL
            mediaData.value.image = serverUrl
            uploadInfo.serverUrl = serverUrl
          }
        } catch (error) {
          // 封面图片上传失败，记录错误
          console.error('封面图片上传失败:', error)
          showErrorMessage('upload_failed')
        }
      }

      // 清理本地URL
      for (const [localUrl] of pendingCoverUploads.value) {
        URL.revokeObjectURL(localUrl)
      }
      pendingCoverUploads.value.clear()
    }

    const { createMedia, updateMedia } = await import('@/api/media/edit')

    // 后端字段名：Poster, Name, Review, Rating, Type
    const mediaPayload = {
      Poster: mediaData.value.image || '', // 前端 image 后端 Poster，允许为空
      Name: mediaData.value.name,
      Review: mediaData.value.description || '', // 前端 description 后端 Review，允许为空
      Rating: mediaData.value.rating,
      Type: mediaData.value.type,
      Date: new Date().toISOString().split('T')[0] // 添加日期
    }

    if (isEditing.value) {
      await updateMedia(mediaData.value.type, route.params.id, mediaPayload)
      showSuccessMessage('update')
    } else {
      await createMedia(user, mediaPayload, mediaData.value.type) // (user, media, type)
      showSuccessMessage('submit')
    }

    // 跳转到统一媒体页面（所有媒体类型都使用 /fragments/novels）
    router.push('/fragments/novels')
  } catch (error) {
    console.error('保存媒体卡片失败:', error)
    showErrorMessage(error)
  }
}

// 选择图片
const selectImage = (url) => {
  if (contentType.value === 'article') {
    articleData.value.image = url
  } else {
    mediaData.value.image = url
  }
  showImagePicker.value = false
}

// 上传图片到服务器（通用函数）
const uploadImageFile = async (file) => {
  // 检查文件大小（5MB）
  if (file.size > 5 * 1024 * 1024) {
    showWarningMessage('file_too_large')
    return null
  }

  // 检查文件类型
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    showWarningMessage('invalid_file_type')
    return null
  }

  // 创建 FormData
  const formData = new FormData()
  formData.append('image', file)

  try {
    // 上传图片
    const token = getJWT()
    const response = await axios.post(`${apiConfig.apiURL}/upload/image`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${token}`
      }
    })

    // 返回图片 URL
    return apiConfig.baseURL + response.data.url
  } catch (error) {
    showErrorMessage('upload_failed')
    throw error
  }
}

// 处理图片上传（文件选择）
const handleImageUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  const textarea = editorTextarea.value
  const start = textarea.selectionStart

  try {
    // 显示上传中提示
    const uploadingText = '\n\n![上传中...]()\n\n'
    const beforeText = markdownContent.value.substring(0, start)
    const afterText = markdownContent.value.substring(start)
    markdownContent.value = beforeText + uploadingText + afterText

    // 上传图片
    const imageUrl = await uploadImageFile(file)
    if (!imageUrl) {
      // 上传失败，移除上传中文本
      markdownContent.value = markdownContent.value.replace(uploadingText, '')
      return
    }

    // 替换上传中文本为实际图片
    markdownContent.value = markdownContent.value.replace(
      uploadingText,
      `\n\n![${file.name}](${imageUrl})\n\n`
    )

    // 更新预览
    updatePreview()

    showSuccessMessage('upload')
  } catch (error) {
    // 移除上传中文件
    markdownContent.value = markdownContent.value.replace('\n\n![上传中...]()\n\n', '')
    showErrorMessage('upload_failed')
  }

  // 清空 input
  event.target.value = ''
}

// 处理粘贴事件（Ctrl+V 粘贴图片）
const handlePaste = async (event) => {
  const items = event.clipboardData.items
  if (!items) return

  // 查找图片
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item.type.indexOf('image') !== -1) {
      // 阻止默认粘贴行为
      event.preventDefault()

      // 获取图片文件
      const file = item.getAsFile()
      if (!file) continue

      const textarea = editorTextarea.value
      const start = textarea.selectionStart

      try {
        // 创建本地预览URL
        const localUrl = URL.createObjectURL(file)

        // 生成文件名
        const timestamp = new Date().getTime()
        const filename = `paste_${timestamp}.png`

        // 插入本地预览的Markdown
        const imageMarkdown = `\n\n![${filename}](${localUrl})\n\n`
        const beforeText = markdownContent.value.substring(0, start)
        const afterText = markdownContent.value.substring(start)
        markdownContent.value = beforeText + imageMarkdown + afterText

        // 添加到待上传队列
        pendingUploads.value.set(localUrl, {
          file,
          localUrl,
          serverUrl: null
        })

        // 更新预览
        updatePreview()

        // console.log('图片本地预览创建成功，等待保存时上传')
      } catch (error) {
        // console.error('创建本地预览失败:', error)
        showErrorMessage('paste_failed')
      }

      break
    }
  }
}

// 插入 Markdown 语法
const insertMarkdown = (before, after) => {
  const textarea = editorTextarea.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = markdownContent.value.substring(start, end)
  const beforeText = markdownContent.value.substring(0, start)
  const afterText = markdownContent.value.substring(end)

  markdownContent.value = beforeText + before + selectedText + after + afterText

  // 更新光标位置
  nextTick(() => {
    const newPos = start + before.length + selectedText.length
    textarea.focus()
    textarea.setSelectionRange(newPos, newPos)
  })

  updatePreview()
}

// 处理 Tab 键
const handleTab = (event) => {
  const textarea = editorTextarea.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd

  // 插入两个空格
  markdownContent.value = markdownContent.value.substring(0, start) + '  ' + markdownContent.value.substring(end)

  // 更新光标位置
  nextTick(() => {
    textarea.selectionStart = textarea.selectionEnd = start + 2
  })
}

// 移除API列表加载函数，改为后端随机选择

// 触发文件选择
const triggerFileUpload = () => {
  if (fileInput.value) {
    fileInput.value.click()
  }
}

// 处理文件上传
const handleFileUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    showErrorMessage('只支持图片文件')
    return
  }

  // 检查文件大小 (10MB)
  if (file.size > 10 * 1024 * 1024) {
    showErrorMessage('文件大小不能超过10MB')
    return
  }

  await uploadCustomImage(file)
}

// 拖拽处理
const handleDragOver = (event) => {
  event.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (event) => {
  event.preventDefault()
  isDragOver.value = false
}

const handleDrop = async (event) => {
  event.preventDefault()
  isDragOver.value = false

  const files = event.dataTransfer.files
  if (files.length === 0) return

  const file = files[0]
  if (!file.type.startsWith('image/')) {
    showErrorMessage('只支持图片文件')
    return
  }

  if (file.size > 10 * 1024 * 1024) {
    showErrorMessage('文件大小不能超过10MB')
    return
  }

  await uploadCustomImage(file)
}

// 上传自定义图片（延时上传版本）
const uploadCustomImage = async (file) => {
  try {
    // console.log('开始处理自定义图片:', file.name)

    // 创建本地预览URL
    const localUrl = URL.createObjectURL(file)
    // console.log('创建本地预览URL:', localUrl)

    // 设置图片URL为本地预览
    if (contentType.value === 'article') {
      articleData.value.image = localUrl
    } else {
      mediaData.value.image = localUrl
    }

    // 添加到封面图片上传队列
    pendingCoverUploads.value.set(localUrl, {
      file,
      localUrl,
      serverUrl: null
    })

    // 显示成功消息
    showCustomMessage(`图片已选择，将在保存时上传: ${file.name}`)

    // 设置加载状态
    imageLoadingStates.value.set(localUrl, true)
    imageRetryCount.value.set(localUrl, 0)

    // 异步清除加载状态
    setTimeout(() => {
      imageLoadingStates.value.set(localUrl, false)
    }, 100)

    // console.log('自定义图片已添加到上传队列:', localUrl)
  } catch (error) {
    // console.error('处理自定义图片失败:', error)
    showErrorMessage('处理失败: ' + error.message)
  }
}

// 获取随机图片
const getRandomImage = async () => {
  try {
    // console.log('正在获取随机图片...')

    const token = getJWT()
    // console.log('当前JWT token:', token ? '存在' : '不存在')
    if (!token) {
      showErrorMessage('请先登录后再获取随机图片')
      return
    }

    // 添加防缓存参数
    const cacheBuster = Date.now()
    // console.log('发送请求到后端...')
    const response = await requestFunc('/random-image', {
      method: 'POST',
      data: {
        cacheBuster
      }
    }, true)

    // console.log('后端响应:', response)
    if (!response) {
      // console.error('请求失败，可能是token过期')
      throw new Error('请求失败，可能是token过期')
    }

    const data = response.data
    // console.log('后端返回的数据', data)

    if (!data.success || !data.imageUrl) {
      throw new Error(data.message || '后端未能获取到图片URL')
    }

    // 显示API来源信息
    if (data.apiName) {
      showCustomMessage(`图片来源: ${data.apiName}`)
    }

    let imageUrl = data.imageUrl
    // console.log('后端返回的图片URL:', imageUrl)

    // 在图片URL后添加时间戳，强制刷新
    const separator = imageUrl.includes('?') ? '&' : '?'
    imageUrl = `${imageUrl}${separator}t=${Date.now()}`

    // 设置图片URL
    if (contentType.value === 'article') {
      articleData.value.image = imageUrl
    } else {
      mediaData.value.image = imageUrl
    }

    // 重置该图片的加载状态和重试计数
    imageLoadingStates.value.set(imageUrl, true)
    imageRetryCount.value.set(imageUrl, 0)

    // console.log('最终使用的图片URL（带防缓存）:', imageUrl)

    // 异步预加载图片并清除加载状态
    setTimeout(async () => {
      try {
        await preloadImage(imageUrl)
        // console.log('图片预加载成功', imageUrl)
        // 预加载成功后清除加载状态
        imageLoadingStates.value.set(imageUrl, false)
      } catch (error) {
        // console.log('图片预加载失败，但不影响显示:', error)
        // 即使预加载失败，也清除加载状态，让图片正常显示
        imageLoadingStates.value.set(imageUrl, false)
      }
    }, 100)
  } catch (error) {
    // console.error('获取随机图片失败:', error)
    // console.error('错误详情:', {
    //   message: error.message,
    //   response: error.response,
    //   status: error.response?.status,
    //   data: error.response?.data
    // })
    await useFallbackImage()
  }
}

// 备用图片方案
const useFallbackImage = async () => {
  try {
    // 使用后端API获取备用图片
    const token = getJWT()
    if (!token) {
      showErrorMessage('请先登录后再获取随机图片')
      return
    }

    const response = await requestFunc('/random-image', {
      method: 'POST',
      data: {
        apiUrl: 'https://api.r10086.com/img-api.php?type=动漫综合1',
        apiName: '樱道备用',
        cacheBuster: Date.now()
      }
    }, true)

    if (!response) {
      throw new Error('请求失败，可能是token过期')
    }

    const data = response.data
    if (!data.success || !data.imageUrl) {
      throw new Error(data.message || '后端未能获取到备用图片URL')
    }

    let imageUrl = data.imageUrl

    // 在图片URL后添加时间戳，强制刷新
    const separator = imageUrl.includes('?') ? '&' : '?'
    imageUrl = `${imageUrl}${separator}t=${Date.now()}`

    if (contentType.value === 'article') {
      articleData.value.image = imageUrl
    } else {
      mediaData.value.image = imageUrl
    }

    // 设置加载状态
    imageLoadingStates.value.set(imageUrl, true)
    imageRetryCount.value.set(imageUrl, 0)

    // 异步清除加载状态
    setTimeout(() => {
      imageLoadingStates.value.set(imageUrl, false)
    }, 100)

    showWarningMessage('使用稳定备用图片')
  } catch (fallbackError) {
    // console.error('备用方案失败:', fallbackError)
    await useFinalFallbackImage()
  }
}

// 最终备用图片方案
const useFinalFallbackImage = async () => {
  try {
    // 最终降级方案：使用樱道API
    const finalFallbacks = [
      'https://api.r10086.com/img-api.php?type=动漫综合1',
      'https://api.r10086.com/img-api.php?type=原神横屏系列1',
      'https://api.r10086.com/img-api.php?type=鬼灭之刃横屏系列1',
      'https://api.r10086.com/img-api.php?type=火影忍者横屏系列',
      'https://api.r10086.com/img-api.php?type=海贼王横屏系列',
      'https://api.r10086.com/img-api.php?type=进击的巨人横屏系列',
      'https://api.r10086.com/img-api.php?type=刀剑神域横屏系列',
      'https://api.r10086.com/img-api.php?type=Fate横屏系列1',
      'https://api.r10086.com/img-api.php?type=明日方舟1',
      'https://api.r10086.com/img-api.php?type=少女前线1',
      'https://api.r10086.com/img-api.php?type=东方project1',
      'https://api.r10086.com/img-api.php?type=P站系列',
      'https://api.r10086.com/img-api.php?type=CG系列1',
      'https://api.r10086.com/img-api.php?type=猫娘1'
    ]

    const token = getJWT()
    if (!token) {
      showErrorMessage('请先登录后再获取随机图片')
      return
    }

    // 尝试每个最终备用源
    for (const finalUrl of finalFallbacks) {
      // console.log('尝试最终备用图片源:', finalUrl)

      try {
        const response = await requestFunc('/random-image', {
          method: 'POST',
          data: {
            apiUrl: finalUrl,
            apiName: '最终备选',
            cacheBuster: Date.now()
          }
        }, true)

        if (response && response.data.success && response.data.imageUrl) {
          let imageUrl = response.data.imageUrl

          // 在图片URL后添加时间戳，强制刷新
          const separator = imageUrl.includes('?') ? '&' : '?'
          imageUrl = `${imageUrl}${separator}t=${Date.now()}`

          if (contentType.value === 'article') {
            articleData.value.image = imageUrl
          } else {
            mediaData.value.image = imageUrl
          }

          // 设置加载状态
          imageLoadingStates.value.set(imageUrl, true)
          imageRetryCount.value.set(imageUrl, 0)

          // 异步清除加载状态
          setTimeout(() => {
            imageLoadingStates.value.set(imageUrl, false)
          }, 100)

          showErrorMessage('使用最终备用图片源')
          return
        }
      } catch (error) {
        // console.warn('备用源失效', finalUrl, error)
        continue
      }
    }

    // 如果所有图片源都失败，显示错误信息
    showErrorMessage('所有图片源都无法访问，请稍后重试')
  } catch (error) {
    // console.error('最终备用方案失效', error)
    showErrorMessage('图片获取失败，请手动输入图片链接')
  }
}

// 图片预加载函数
const preloadImage = (url) => {
  return new Promise((resolve, reject) => {
    // 检查缓存
    if (imageCache.value.has(url)) {
      resolve(url)
      return
    }

    const img = new Image()
    img.onload = () => {
      imageCache.value.set(url, true)
      resolve(url)
    }
    img.onerror = () => {
      reject(new Error('图片加载失败: ' + url))
    }
    img.src = url
  })
}

// 图片错误处理
const handleImageError = (event) => {
  if (!event.target) return

  const img = event.target
  const imageUrl = img.src

  // console.log('图片加载失败:', imageUrl)

  // 获取当前重试次数
  const currentRetryCount = imageRetryCount.value.get(imageUrl) || 0

  if (currentRetryCount < 2) { // 减少重试次数
    // 增加重试次数
    imageRetryCount.value.set(imageUrl, currentRetryCount + 1)

    // 设置加载状态
    imageLoadingStates.value.set(imageUrl, true)

    // console.log('尝试重新加载图片:', imageUrl)

    // 使用时间戳强制刷新图片
    const separator = imageUrl.includes('?') ? '&' : '?'
    const newUrl = `${imageUrl}${separator}t=${Date.now()}`

    // 更新图片
    if (event.target) {
      event.target.src = newUrl
    }
  } else {
    // console.log('图片重试次数已达上限:', imageUrl)
    imageLoadingStates.value.set(imageUrl, false)
  }
}

// 图片加载成功处理
const handleImageLoad = (event) => {
  if (!event.target) return

  const img = event.target
  const imageUrl = img.src

  // console.log('图片加载成功:', imageUrl)

  // 清除加载状态和重试计数
  imageLoadingStates.value.set(imageUrl, false)
  imageRetryCount.value.set(imageUrl, 0)
}

// 清理函数
const cleanup = () => {
  imageCache.value.clear()
  imageLoadingStates.value.clear()
  imageRetryCount.value.clear()

  // 清理待上传队列中的本地URL
  for (const [localUrl] of pendingUploads.value) {
    URL.revokeObjectURL(localUrl)
  }
  pendingUploads.value.clear()

  // 清理封面图片上传队列中的本地URL
  for (const [localUrl] of pendingCoverUploads.value) {
    URL.revokeObjectURL(localUrl)
  }
  pendingCoverUploads.value.clear()
}

// 加载现有文章数据
const loadExistingArticle = async () => {
  if (!isEditing.value) return

  try {
    const id = route.params.id
    const articleType = articleData.value.type || route.query.articleType || 'blog'

    // 碎碎念使用不同的 API
    if (articleType === 'moment') {
      const { getMoment } = await import('@/api/Moments/browse')
      const res = await getMoment(id)
      const data = res.data

      // 填充文章数据
      articleData.value.title = data.Title || ''
      articleData.value.image = data.Image || ''
      articleData.value.type = 'moment'
      articleData.value.isTop = false

      // 填充内容
      markdownContent.value = data.Content || ''
      tagsInput.value = ''

      // 保存原始数据用于检测更改
      originalArticleData.value = JSON.parse(JSON.stringify(articleData.value))
      originalMarkdownContent.value = markdownContent.value
      originalTagsInput.value = tagsInput.value

      // 更新预览
      updatePreview()
    } else {
      // 其他类型使用标准 API
      const res = await getArticleByID(articleType, id)
      const data = res.data

      // 填充文章数据
      articleData.value.title = data.title || ''
      articleData.value.image = data.image || ''
      articleData.value.type = articleType
      articleData.value.isTop = data.isTop || false

      // 填充内容
      markdownContent.value = data.content || ''
      tagsInput.value = (data.tags || []).join(',')

      // 保存原始数据用于检测更改
      originalArticleData.value = JSON.parse(JSON.stringify(articleData.value))
      originalMarkdownContent.value = markdownContent.value
      originalTagsInput.value = tagsInput.value

      // 更新预览
      updatePreview()
    }
  } catch (error) {
    console.error('加载文章失败:', error)
    showErrorMessage('加载文章失败，请重试')
  }
}

// 加载现有媒体数据
const loadExistingMedia = async () => {
  if (!isEditing.value) return

  try {
    const id = route.params.id
    const mediaType = mediaData.value.type

    const res = await getMediaByID(id, mediaType)
    const data = res.data || res // 后端可能直接返回数据或包装在data中

    // 填充媒体数据（后端字段名：Poster, Name, Review, Rating, Type）
    mediaData.value.name = data.Name || ''
    mediaData.value.image = data.Poster || ''
    mediaData.value.rating = data.Rating || 8
    mediaData.value.description = data.Review || ''
    mediaData.value.type = data.Type || mediaType

    // 保存原始数据用于检测更改
    originalMediaData.value = JSON.parse(JSON.stringify(mediaData.value))
  } catch (error) {
    console.error('加载媒体数据失败:', error)
    showErrorMessage('加载媒体数据失败，请重试')
  }
}

// 组件挂载时初始化
onMounted(async () => {
  // 从路由query参数设置内容类型
  if (route.query.contentType) {
    contentType.value = route.query.contentType
  }

  // 从路由query参数设置类型
  if (route.query.articleType) {
    articleData.value.type = route.query.articleType
  }

  if (route.query.mediaType) {
    mediaData.value.type = route.query.mediaType
  }

  // 如果是编辑模式，加载现有数据
  if (isEditing.value) {
    if (contentType.value === 'media') {
      await loadExistingMedia()
      // 初始化媒体预览
      updateMediaPreview()
    } else {
      await loadExistingArticle()
      // 初始化预览
      updatePreview()
    }
  } else {
    // 非编辑模式也要初始化预览
    if (contentType.value === 'media') {
      updateMediaPreview()
    } else {
      updatePreview()
    }
  }
})

// 组件卸载时清理
onBeforeUnmount(() => {
  cleanup()
})

// 图片重试加载
const retryImageLoad = () => {
  const currentImage = contentType.value === 'article' ? articleData.value.image : mediaData.value.image
  if (currentImage) {
    imageRetryCount.value.set(currentImage, 0)
    imageLoadingStates.value.set(currentImage, true)
  }
}

// 处理图片URL输入
const handleImageUrlInput = () => {
  const currentImage = contentType.value === 'article' ? articleData.value.image : mediaData.value.image
  if (currentImage) {
    imageLoadingStates.value.set(currentImage, true)
    imageRetryCount.value.set(currentImage, 0)

    // 异步清除加载状态，避免一直显示加载中
    setTimeout(() => {
      imageLoadingStates.value.set(currentImage, false)
    }, 500) // 给图片一些加载时间
  }
}

// 所见即所得编辑器命令
const execWysiwygCommand = (command) => {
  if (!wysiwygContent.value) return

  // 使用现代的选择API替代已废弃的execCommand
  const selection = window.getSelection()
  if (selection.rangeCount > 0) {
    const range = selection.getRangeAt(0)
    const selectedText = selection.toString()

    if (selectedText) {
      // 有选中文本时，应用格式
      let formattedText = selectedText
      switch (command) {
        case 'bold':
          formattedText = `<strong>${selectedText}</strong>`
          break
        case 'italic':
          formattedText = `<em>${selectedText}</em>`
          break
        case 'strikeThrough':
          formattedText = `<del>${selectedText}</del>`
          break
      }

      // 替换选中文本
      range.deleteContents()
      const tempDiv = document.createElement('div')
      tempDiv.innerHTML = formattedText
      const fragment = document.createDocumentFragment()
      while (tempDiv.firstChild) {
        fragment.appendChild(tempDiv.firstChild)
      }
      range.insertNode(fragment)
    } else {
      // 没有选中文本时，插入格式标记
      let formatMark = ''
      switch (command) {
        case 'bold':
          formatMark = '**粗体文本**'
          break
        case 'italic':
          formatMark = '*斜体文本*'
          break
        case 'strikeThrough':
          formatMark = '~~删除线文本~~'
          break
      }

      range.insertNode(document.createTextNode(formatMark))
    }
  }

  // 触发输入事件以更新渲染
  nextTick(() => {
    handleWysiwygInput()
  })
}

// 标题命令
const execHeadingCommand = (level) => {
  const headingText = '#'.repeat(level) + ' '
  insertMarkdown(headingText, '')
}

// 链接命令
const execLinkCommand = () => {
  const url = prompt('请输入链接地址:')
  if (url) {
    insertMarkdown('[', `](${url})`)
  }
}

// 图片命令
const execImageCommand = () => {
  if (imageInput.value) {
    imageInput.value.click()
  }
}

// 代码命令
const execCodeCommand = () => {
  insertMarkdown('`', '`')
}

// 插入表格
const insertTable = () => {
  const tableMarkdown = '\n| 列1 | 列2 | 列3 |\n|-----|-----|-----|\n| 内容1 | 内容2 | 内容3 |\n| 内容4 | 内容5 | 内容6 |\n\n'
  insertMarkdown(tableMarkdown, '')
}

// 切换编辑器模式（已移除，现在使用分屏模式）

// 更新光标位置
const updateCursorPosition = () => {
  if (!editorTextarea.value) return

  const textarea = editorTextarea.value
  const text = textarea.value
  const cursorPos = textarea.selectionStart

  // 计算行号和列号
  const textBeforeCursor = text.substring(0, cursorPos)
  const lines = textBeforeCursor.split('\n')
  cursorLine.value = lines.length
  cursorColumn.value = lines[lines.length - 1].length + 1
}

// 滚动同步
const syncScroll = () => {
  if (!editorTextarea.value) return

  const textarea = editorTextarea.value
  const previewContent = document.querySelector('.preview-content')

  if (!previewContent) return

  // 计算滚动比例
  const scrollTop = textarea.scrollTop
  const scrollHeight = textarea.scrollHeight
  const clientHeight = textarea.clientHeight
  const maxScroll = scrollHeight - clientHeight

  if (maxScroll > 0) {
    const scrollRatio = scrollTop / maxScroll
    const previewMaxScroll = previewContent.scrollHeight - previewContent.clientHeight
    const targetScrollTop = scrollRatio * previewMaxScroll

    previewContent.scrollTop = targetScrollTop
  }
}

// 所见即所得输入处理
const handleWysiwygInput = () => {
  if (!wysiwygContent.value) return

  // 保存光标位置
  const selection = saveSelection(wysiwygContent.value)

  // 1. 获取当前编辑器的纯文本内容
  const rawText = wysiwygContent.value.innerText || wysiwygContent.value.textContent

  // 2. 将纯文本转换为Markdown格式保存
  markdownContent.value = rawText

  // 3. 立即将Markdown渲染为HTML显示在编辑器中
  nextTick(() => {
    const renderedHTML = marked(rawText)
    wysiwygContent.value.innerHTML = DOMPurify.sanitize(renderedHTML)

    // 恢复光标位置
    restoreSelection(wysiwygContent.value, selection)
  })
}

// 所见即所得相关函数（已移除，现在使用分屏模式）

// 保存光标位置
const saveSelection = (container) => {
  const selection = window.getSelection()
  if (selection.rangeCount === 0) return null

  const range = selection.getRangeAt(0)
  const preSelectionRange = range.cloneRange()
  preSelectionRange.selectNodeContents(container)
  preSelectionRange.setEnd(range.startContainer, range.startOffset)

  return {
    start: preSelectionRange.toString().length,
    end: preSelectionRange.toString().length + range.toString().length
  }
}

// 恢复光标位置
const restoreSelection = (container, savedSel) => {
  if (!savedSel) return

  const selection = window.getSelection()
  const range = document.createRange()

  let charIndex = 0
  const nodeStack = [container]
  let node
  let foundStart = false
  let stop = false

  while (!stop && (node = nodeStack.pop())) {
    if (node.nodeType === 3) {
      const nextCharIndex = charIndex + node.length
      if (!foundStart && savedSel.start >= charIndex && savedSel.start <= nextCharIndex) {
        range.setStart(node, savedSel.start - charIndex)
        foundStart = true
      }
      if (foundStart && savedSel.end >= charIndex && savedSel.end <= nextCharIndex) {
        range.setEnd(node, savedSel.end - charIndex)
        stop = true
      }
      charIndex = nextCharIndex
    } else {
      let i = node.childNodes.length
      while (i--) {
        nodeStack.push(node.childNodes[i])
      }
    }
  }

  selection.removeAllRanges()
  selection.addRange(range)
}

// 同步所见即所得内容到Markdown（已移除，现在使用分屏模式）

// 放弃更改
const discardChanges = () => {
  if (confirm('确定要放弃所有更改吗？')) {
    if (contentType.value === 'article') {
      articleData.value = { ...originalArticleData.value }
      markdownContent.value = originalMarkdownContent.value
      tagsInput.value = originalTagsInput.value
    } else {
      mediaData.value = { ...originalMediaData.value }
    }
  }
}

// 导出函数
defineExpose({
  getRandomImage,
  useFallbackImage,
  useFinalFallbackImage
})
</script>

  <style scoped>
  /* 主容器样式 */
  .edit-article-view {
    min-height: 100vh;
    background: transparent;
    padding: 20px;
    position: relative;
  }

  .edit-container {
    max-width: 1200px;
    margin: 80px auto 0 auto;
    background: rgba(255, 255, 255, 0.25);
    border-radius: 20px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    overflow: hidden;
  }

  .article-form {
    padding: 30px;
  }

  /* 表单组样式 */
  .form-group {
    margin-bottom: 25px;
  }

  .form-group.full-width {
    width: 100%;
  }

  /* 两列布局 */
  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
    margin-bottom: 25px;
  }

  .form-row .form-group {
    margin-bottom: 0;
  }

  /* 单列布局（用于需要全宽的元素） */
  .form-row-single {
    display: grid;
    grid-template-columns: 1fr;
    gap: 20px;
    margin-bottom: 25px;
  }

  .form-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 600;
    color: #333;
    font-size: 14px;
    text-align: left;
  }

  .form-group input,
  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 12px 16px;
    border: 2px solid #e1e5e9;
    border-radius: 12px;
    font-size: 14px;
    transition: all 0.3s ease;
    background: #fff;
    text-align: left;
  }

  .form-group input:focus,
  .form-group select:focus,
  .form-group textarea:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  /* 内容类型选择器 */
  .content-type-selector {
    margin-bottom: 30px;
  }

  .type-buttons {
    display: flex;
    gap: 15px;
    margin-top: 10px;
  }

  .type-btn {
    flex: 1;
    padding: 15px 20px;
    border: 2px solid #e1e5e9;
    border-radius: 12px;
    background: #fff;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 14px;
    font-weight: 500;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .type-btn:hover {
    border-color: #667eea;
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.2);
  }

  .type-btn.active {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border-color: #667eea;
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
  }

  /* 复选框容器样式 */
  .checkbox-container {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    cursor: pointer;
  }

  .checkbox-container input[type="checkbox"] {
    width: 24px;
    height: 24px;
    margin: 40px 0 0 0;
    cursor: pointer;
  }

  .checkbox-label-text {
    font-size: 16px;
    color: #333;
    cursor: pointer;
    margin: 45px 0 0 0;
    font-weight: 500;
  }

  /* 图片管理器样式 */
  .image-manager {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  /* 图片预览和输入框横向布局 */
  .image-preview-input-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
    align-items: start;
  }

  .image-preview-section {
    width: 100%;
  }

  .image-preview-container {
    position: relative;
    width: 100%;
    height: 200px;
    border: 2px dashed #e1e5e9;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    background: #f8f9fa;
    transition: all 0.3s ease;
  }

  .image-preview-container.drag-over {
    border-color: #28a745;
    background: rgba(40, 167, 69, 0.1);
    transform: scale(1.02);
  }

  /* 媒体预览容器 - 方形比例 */
  .media-preview-container {
    width: 200px;
    height: 200px;
    margin: 0 auto;
  }

  .image-preview {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .image-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center;
    border-radius: 8px;
  }

  .image-overlay {
    position: absolute;
    top: 10px;
    right: 10px;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .image-preview:hover .image-overlay {
    opacity: 1;
  }

  .clear-btn {
    background: rgba(255, 0, 0, 0.8);
    color: white;
    border: none;
    border-radius: 50%;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .clear-btn:hover {
    background: rgba(255, 0, 0, 1);
    transform: scale(1.1);
  }

  .image-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #999;
    gap: 10px;
  }

  .placeholder-icon {
    font-size: 48px;
    color: #ddd;
  }

  .placeholder-text {
    font-size: 14px;
    color: #999;
  }

  .drag-hint {
    font-size: 12px;
    color: #28a745;
    margin-top: 4px;
    font-weight: 500;
  }

  .image-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 15px;
    color: #667eea;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }

  .loading-text {
    font-size: 14px;
    color: #667eea;
  }

  .image-error {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 15px;
    color: #e74c3c;
  }

  .error-icon {
    font-size: 48px;
    color: #e74c3c;
  }

  .error-text {
    font-size: 14px;
    color: #e74c3c;
  }

  .retry-btn {
    background: #e74c3c;
    color: white;
    border: none;
    border-radius: 8px;
    padding: 8px 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    gap: 5px;
  }

  .retry-btn:hover {
    background: #c0392b;
    transform: translateY(-2px);
  }

  .image-input-section {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }

  .image-upload-controls {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .upload-image-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 12px;
    padding: 12px 20px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 500;
    white-space: nowrap;
  }

  .upload-image-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
  }

  .input-group {
    flex: 1;
  }

  .input-label {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
    font-weight: 600;
    color: #333;
    font-size: 14px;
  }

  .image-url-input {
    width: 100%;
    padding: 12px 16px;
    border: 2px solid #e1e5e9;
    border-radius: 12px;
    font-size: 14px;
    transition: all 0.3s ease;
  }

  .image-url-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .random-image-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 12px;
    padding: 12px 20px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 500;
    white-space: nowrap;
  }

  .random-image-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
  }

  /* 图片选择器样式 */
  .image-picker {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .image-picker h4 {
    color: white;
    margin-bottom: 20px;
    text-align: center;
  }

  .image-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 15px;
    max-width: 600px;
    max-height: 400px;
    overflow-y: auto;
    padding: 20px;
    background: white;
    border-radius: 12px;
  }

  .image-option {
    cursor: pointer;
    border: 2px solid #e1e5e9;
    border-radius: 8px;
    overflow: hidden;
    transition: all 0.3s ease;
  }

  .image-option:hover {
    border-color: #667eea;
    transform: scale(1.05);
  }

  .image-option.selected {
    border-color: #667eea;
    background: rgba(102, 126, 234, 0.1);
  }

  .image-option img {
    width: 100%;
    height: 100px;
    object-fit: cover;
  }

  .image-name {
    display: block;
    padding: 8px;
    text-align: center;
    font-size: 12px;
    color: #666;
  }

  /* Typora风格编辑器样式 */
  .typora-editor {
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
    /* 让编辑器距离侧边230px */
    margin: 40px 230px 20px 230px;
  }

  .editor-toolbar {
    display: flex;
    align-items: center;
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.1);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    gap: 8px;
    flex-wrap: wrap;
  }

  /* 左侧面板工具栏样式 */
  .panel-toolbar {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    background: rgba(255, 255, 255, 0.05);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    gap: 6px;
    flex-wrap: wrap;
    font-size: 12px;
  }

  .toolbar-group {
    display: flex;
    gap: 4px;
  }

  .toolbar-divider {
    width: 1px;
    height: 24px;
    background: rgba(255, 255, 255, 0.3);
    margin: 0 8px;
  }

  .toolbar-spacer {
    flex: 1;
  }

  .toolbar-btn {
    padding: 8px 12px;
    border: none;
    background: transparent;
    color: #666;
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.2s ease;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  /* 左侧面板工具栏按钮样式 */
  .panel-toolbar .toolbar-btn {
    padding: 6px 8px;
    font-size: 12px;
  }

  .toolbar-btn:hover {
    background: rgba(102, 126, 234, 0.1);
    color: #667eea;
  }

  .toolbar-btn:active {
    transform: scale(0.95);
  }

  .editor-content {
    min-height: 500px;
    position: relative;
  }

  /* 分屏编辑器样式 */
  .split-editor {
    display: flex;
    height: 600px;
    gap: 20px;
    background: transparent;
  }

  .editor-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 0;
    overflow: hidden;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  }

  /* 左侧面板圆角 */
  .left-panel {
    border-top-left-radius: 12px;
    border-bottom-left-radius: 12px;
  }

  /* 右侧面板圆角 */
  .right-panel {
    border-top-right-radius: 12px;
    border-bottom-right-radius: 12px;
  }

  .panel-header {
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.1);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    font-weight: 600;
    color: #333;
    font-size: 14px;
  }

  .panel-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .left-panel textarea {
    width: 100%;
    height: 100%;
    border: none;
    outline: none;
    padding: 16px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
    line-height: 1.6;
    background: transparent;
    resize: none;
    /* 自定义滚动条样式 - 灰色 */
    scrollbar-width: thin;
    scrollbar-color: #ccc #f5f5f5;
  }

  /* Webkit浏览器滚动条样式 */
  .left-panel textarea::-webkit-scrollbar {
    width: 8px;
  }

  .left-panel textarea::-webkit-scrollbar-track {
    background: #f5f5f5;
    border-radius: 4px;
  }

  .left-panel textarea::-webkit-scrollbar-thumb {
    background: #ccc;
    border-radius: 4px;
  }

  .left-panel textarea::-webkit-scrollbar-thumb:hover {
    background: #999;
  }

  .right-panel .preview-content {
    flex: 1;
    padding: 16px;
    overflow-y: auto;
    line-height: 1.6;
    font-size: 16px;
    /* 隐藏滚动条但保持滚动功能 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE and Edge */
  }

  /* 强制标题样式 - 使用更高优先级 */
  .right-panel .preview-content.markdown-body h1 {
    font-size: 2em !important;
    font-weight: 600 !important;
    margin: 0.67em 0 !important;
    border-bottom: 1px solid #eaecef !important;
    padding-bottom: 0.3em !important;
    line-height: 1.25 !important;
  }

  .right-panel .preview-content.markdown-body h2 {
    font-size: 1.5em !important;
    font-weight: 600 !important;
    margin: 0.83em 0 !important;
    border-bottom: 1px solid #eaecef !important;
    padding-bottom: 0.3em !important;
    line-height: 1.25 !important;
  }

  .right-panel .preview-content.markdown-body h3 {
    font-size: 1.25em !important;
    font-weight: 600 !important;
    margin: 1em 0 !important;
    line-height: 1.25 !important;
  }

  .right-panel .preview-content.markdown-body h4 {
    font-size: 1em !important;
    font-weight: 600 !important;
    margin: 1.33em 0 !important;
    line-height: 1.25 !important;
  }

  .right-panel .preview-content.markdown-body h5 {
    font-size: 0.875em !important;
    font-weight: 600 !important;
    margin: 1.67em 0 !important;
    line-height: 1.25 !important;
  }

  .right-panel .preview-content.markdown-body h6 {
    font-size: 0.85em !important;
    font-weight: 600 !important;
    margin: 2.33em 0 !important;
    color: #6a737d !important;
    line-height: 1.25 !important;
  }

  /* 强制表格样式 */
  .right-panel .preview-content.markdown-body table {
    border-collapse: collapse !important;
    width: 100% !important;
    margin: 16px 0 !important;
    border: 1px solid #d0d7de !important;
    display: table !important;
  }

  .right-panel .preview-content.markdown-body th,
  .right-panel .preview-content.markdown-body td {
    border: 1px solid #d0d7de !important;
    padding: 8px 12px !important;
    text-align: left !important;
    vertical-align: top !important;
  }

  .right-panel .preview-content.markdown-body th {
    background-color: #f6f8fa !important;
    font-weight: 600 !important;
  }

  .right-panel .preview-content.markdown-body tr:nth-child(even) {
    background-color: #f6f8fa !important;
  }

  /* 确保表格中的br标签正确换行 */
  .right-panel .preview-content.markdown-body td br,
  .right-panel .preview-content.markdown-body th br {
    display: block !important;
    content: "" !important;
    margin-top: 0.5em !important;
  }

  /* 确保表格单元格内容正确换行 */
  .right-panel .preview-content.markdown-body td,
  .right-panel .preview-content.markdown-body th {
    white-space: normal !important;
    word-wrap: break-word !important;
    line-height: 1.4 !important;
  }

  /* 预览区域图片样式 - 宽度不超过80%，居中显示 */
  .right-panel .preview-content.markdown-body img {
    max-width: 80% !important;
    width: auto !important;
    height: auto !important;
    border-radius: 8px !important;
    margin: 20px auto !important;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
    display: block !important;
    object-fit: contain !important;
  }

  /* 图片占位符样式 */
  .image-placeholder-preview {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px;
    margin: 10px auto;
    max-width: 50%;
    border: 2px dashed #667eea;
    border-radius: 8px;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
    text-align: center;
  }

  .placeholder-icon {
    font-size: 2rem;
    margin-bottom: 8px;
    opacity: 0.7;
  }

  .placeholder-text {
    font-size: 0.9rem;
    font-weight: 500;
    color: #667eea;
    margin-bottom: 4px;
  }

  .placeholder-hint {
    font-size: 0.8rem;
    color: #999;
    font-style: italic;
  }

  /* 媒体预览区域图片样式 */
  .media-preview.markdown-body img {
    max-width: 50% !important;
    max-height: 50vh !important;
    width: auto !important;
    height: auto !important;
    border-radius: 8px !important;
    margin: 10px auto !important;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
    display: block !important;
  }

  /* 隐藏Webkit浏览器滚动条 */
  .right-panel .preview-content::-webkit-scrollbar {
    display: none;
  }

  /* 媒体预览区域文本左对齐 */
  .right-panel .preview-content.markdown-body {
    text-align: left !important;
  }
  .right-panel .preview-content.markdown-body * {
    text-align: left !important;
  }
  .right-panel .preview-content.markdown-body p,
  .right-panel .preview-content.markdown-body div,
  .right-panel .preview-content.markdown-body span {
    text-align: left !important;
  }
  /* 确保图片居中但其他内容左对齐 */
  .right-panel .preview-content.markdown-body img {
    text-align: center !important;
  }

  /* 确保列表缩进正确 */
  .preview-content.markdown-body ul,
  .preview-content.markdown-body ol {
    padding-left: 2em;
  }

  .preview-content.markdown-body ul ul,
  .preview-content.markdown-body ul ol,
  .preview-content.markdown-body ol ul,
  .preview-content.markdown-body ol ol {
    margin-top: 0.25em;
    margin-bottom: 0.25em;
  }

  .wysiwyg-editor {
    height: 100%;
  }

  .wysiwyg-content {
    min-height: 500px;
    padding: 30px 16px;
    line-height: 1.6;
    font-size: 16px;
    outline: none;
    background: transparent;
    text-align: left;
  }

  /* 所见即所得编辑器渲染样式 */
  .wysiwyg-content h1,
  .wysiwyg-content h2,
  .wysiwyg-content h3,
  .wysiwyg-content h4,
  .wysiwyg-content h5,
  .wysiwyg-content h6 {
    margin: 20px 0 10px 0;
    font-weight: 600;
    color: #333;
  }

  .wysiwyg-content h1 {
    font-size: 28px;
    border-bottom: 2px solid #e1e5e9;
    padding-bottom: 10px;
  }

  .wysiwyg-content h2 {
    font-size: 24px;
    border-bottom: 1px solid #e1e5e9;
    padding-bottom: 8px;
  }

  .wysiwyg-content h3 {
    font-size: 20px;
  }

  .wysiwyg-content h4 {
    font-size: 18px;
  }

  .wysiwyg-content h5 {
    font-size: 16px;
  }

  .wysiwyg-content h6 {
    font-size: 14px;
  }

  .wysiwyg-content p {
    margin: 10px 0;
    line-height: 1.6;
  }

  .wysiwyg-content strong {
    font-weight: 600;
    color: #333;
  }

  .wysiwyg-content em {
    font-style: italic;
    color: #555;
  }

  .wysiwyg-content del {
    text-decoration: line-through;
    color: #999;
  }

  .wysiwyg-content code {
    background: #f5f5f5;
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
    color: #e74c3c;
  }

  .wysiwyg-content pre {
    background: #f8f9fa;
    border: 1px solid #e1e5e9;
    border-radius: 8px;
    padding: 16px;
    margin: 16px 0;
    overflow-x: auto;
  }

  .wysiwyg-content pre code {
    background: none;
    padding: 0;
    color: #333;
    font-size: 14px;
  }

  .wysiwyg-content blockquote {
    border-left: 4px solid #667eea;
    margin: 16px 0;
    padding: 0 16px;
    color: #666;
    font-style: italic;
  }

  .wysiwyg-content ul,
  .wysiwyg-content ol {
    margin: 10px 0;
    padding-left: 20px;
  }

  .wysiwyg-content li {
    margin: 5px 0;
  }

  .wysiwyg-content a {
    color: #667eea;
    text-decoration: none;
  }

  .wysiwyg-content a:hover {
    text-decoration: underline;
  }

  .wysiwyg-content img {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
    margin: 10px 0;
  }

  .wysiwyg-content table {
    border-collapse: collapse;
    width: 100%;
    margin: 16px 0;
  }

  .wysiwyg-content th,
  .wysiwyg-content td {
    border: 1px solid #e1e5e9;
    padding: 8px 12px;
    text-align: left;
  }

  .wysiwyg-content th {
    background: #f8f9fa;
    font-weight: 600;
  }

  .wysiwyg-content:empty:before {
    content: attr(data-placeholder);
    color: #999;
    font-style: italic;
    white-space: pre-wrap;
    text-align: left;
  }

  .source-editor {
    height: 100%;
  }

  .source-editor textarea {
    width: 100%;
    height: 500px;
    border: none;
    outline: none;
    padding: 30px 16px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
    line-height: 1.6;
    background: transparent;
    resize: none;
  }

  .editor-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.1);
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    font-size: 12px;
    color: #666;
  }

  /* 左侧面板底部样式 */
  .panel-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: rgba(255, 255, 255, 0.05);
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    font-size: 11px;
    color: #666;
  }

  .editor-stats {
    display: flex;
    gap: 16px;
  }

  .editor-stats span {
    padding: 4px 8px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
  }

  /* 左侧面板统计信息样式 */
  .panel-footer .editor-stats {
    gap: 8px;
  }

  .panel-footer .editor-stats span {
    padding: 2px 6px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 3px;
    font-size: 10px;
  }

  /* 媒体编辑器样式 */
  .media-editor-container {
    display: flex;
    gap: 20px;
    /* 使用与文章编辑器相同的边距 */
    margin: 40px 230px 20px 230px;
  }

  .media-editor {
    flex: 1;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  }

  .editor-header {
    padding: 16px 20px;
    background: rgba(255, 255, 255, 0.1);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }

  .editor-header h3 {
    margin: 0 0 12px 0;
    font-size: 16px;
    color: #333;
  }

  .media-textarea {
    width: 100%;
    height: 400px;
    border: none;
    outline: none;
    padding: 30px 16px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
    line-height: 1.6;
    background: transparent;
    resize: none;
  }

  .media-preview-panel {
    flex: 1;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  }

  .preview-header {
    padding: 16px 20px;
    background: rgba(255, 255, 255, 0.1);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }

  .preview-header h3 {
    margin: 0;
    font-size: 16px;
    color: #333;
  }

  .media-preview {
    padding: 30px 16px;
    max-height: 400px;
    overflow-y: auto;
  }

  /* 媒体预览标题层级样式 */
  .media-preview.markdown-body h1 {
    font-size: 2em !important;
    font-weight: 600 !important;
    margin: 0.67em 0 !important;
    border-bottom: 1px solid #eaecef !important;
    padding-bottom: 0.3em !important;
    line-height: 1.25 !important;
  }

  .media-preview.markdown-body h2 {
    font-size: 1.5em !important;
    font-weight: 600 !important;
    margin: 0.83em 0 !important;
    border-bottom: 1px solid #eaecef !important;
    padding-bottom: 0.3em !important;
    line-height: 1.25 !important;
  }

  .media-preview.markdown-body h3 {
    font-size: 1.25em !important;
    font-weight: 600 !important;
    margin: 1em 0 !important;
    line-height: 1.25 !important;
  }

  .media-preview.markdown-body h4 {
    font-size: 1em !important;
    font-weight: 600 !important;
    margin: 1.33em 0 !important;
    line-height: 1.25 !important;
  }

  .media-preview.markdown-body h5 {
    font-size: 0.875em !important;
    font-weight: 600 !important;
    margin: 1.67em 0 !important;
    line-height: 1.25 !important;
  }

  .media-preview.markdown-body h6 {
    font-size: 0.85em !important;
    font-weight: 600 !important;
    margin: 2.33em 0 !important;
    color: #6a737d !important;
    line-height: 1.25 !important;
  }

  /* 媒体预览表格样式 */
  .media-preview.markdown-body table {
    border-collapse: collapse !important;
    width: 100% !important;
    margin: 16px 0 !important;
    border: 1px solid #d0d7de !important;
    display: table !important;
  }

  .media-preview.markdown-body th,
  .media-preview.markdown-body td {
    border: 1px solid #d0d7de !important;
    padding: 8px 12px !important;
    text-align: left !important;
    vertical-align: top !important;
  }

  .media-preview.markdown-body th {
    background-color: #f6f8fa !important;
    font-weight: 600 !important;
  }

  .media-preview.markdown-body tr:nth-child(even) {
    background-color: #f6f8fa !important;
  }

  /* 确保媒体预览表格中的br标签正确换行 */
  .media-preview.markdown-body td br,
  .media-preview.markdown-body th br {
    display: block !important;
    content: "" !important;
    margin-top: 0.5em !important;
  }

  /* 确保媒体预览表格单元格内容正确换行 */
  .media-preview.markdown-body td,
  .media-preview.markdown-body th {
    white-space: normal !important;
    word-wrap: break-word !important;
    line-height: 1.4 !important;
  }

  .empty-preview {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: #999;
    gap: 12px;
  }

  .empty-icon {
    font-size: 48px;
    color: #ddd;
  }

  .preview-toggle-btn {
    padding: 6px 12px;
    border: 1px solid rgba(255, 255, 255, 0.3);
    background: transparent;
    color: #666;
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .preview-toggle-btn:hover {
    background: rgba(102, 126, 234, 0.1);
    color: #667eea;
  }

  /* 操作按钮样式 */
  .action-buttons {
    display: flex;
    gap: 15px;
    justify-content: center;
    padding: 30px;
    /* 移除外边距，使用容器的内边距 */
    margin: 40px 0 20px 0;
  }

  .save-btn {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 12px;
    padding: 15px 30px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
  }

  .save-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
  }

  .save-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }

  .discard-btn {
    background: transparent;
    color: #e74c3c;
    border: 2px solid #e74c3c;
    border-radius: 12px;
    padding: 13px 28px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .discard-btn:hover {
    background: #e74c3c;
    color: white;
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(231, 76, 60, 0.3);
  }

  /* 响应式设计 */
  @media (max-width: 768px) {
    .edit-article-view {
      padding: 10px;
    }

    .edit-container {
      margin: 0;
      border-radius: 12px;
    }

    .article-form {
      padding: 20px;
    }

    .form-row {
      grid-template-columns: 1fr;
      gap: 15px;
    }

    .image-preview-input-row {
      grid-template-columns: 1fr;
      gap: 15px;
    }

    .media-editor-container {
      flex-direction: column;
      margin: 30px 0 20px 0;
    }

    .typora-editor,
    .action-buttons {
      margin: 30px 0 20px 0;
    }

    .editor-toolbar {
      flex-wrap: wrap;
      gap: 4px;
    }

    .toolbar-btn {
      padding: 6px 8px;
      font-size: 12px;
    }

    /* 移动端左侧面板工具栏 */
    .panel-toolbar {
      padding: 6px 8px;
      gap: 4px;
    }

    .panel-toolbar .toolbar-btn {
      padding: 4px 6px;
      font-size: 10px;
    }

    .panel-footer {
      padding: 6px 8px;
      font-size: 10px;
    }

    .panel-footer .editor-stats {
      gap: 6px;
    }

    .panel-footer .editor-stats span {
      padding: 1px 4px;
      font-size: 9px;
    }

    .panel-footer .theme-btn {
      padding: 3px 6px;
      font-size: 9px;
    }

    .wysiwyg-content,
    .source-editor textarea,
    .media-textarea,
    .media-preview {
      padding: 20px 20px;
    }

    /* 移动端分屏编辑器 */
    .split-editor {
      flex-direction: column;
      gap: 10px;
    }

    .editor-panel {
      min-height: 300px;
    }

    .left-panel textarea,
    .right-panel .preview-content {
      padding: 12px;
    }

    .action-buttons {
      flex-direction: column;
      gap: 10px;
    }

    .save-btn,
    .discard-btn {
      width: 100%;
      padding: 12px 20px;
    }
  }

  @media (max-width: 480px) {
    .type-buttons {
      flex-direction: column;
    }

    .image-grid {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
      gap: 10px;
      padding: 15px;
    }

    .editor-toolbar {
      padding: 8px 12px;
    }

    .toolbar-group {
      gap: 2px;
    }

    .toolbar-btn {
      padding: 4px 6px;
      font-size: 11px;
    }

    /* 超小屏幕左侧面板工具栏 */
    .panel-toolbar {
      padding: 4px 6px;
      gap: 2px;
    }

    .panel-toolbar .toolbar-btn {
      padding: 3px 4px;
      font-size: 9px;
    }

    .panel-footer {
      padding: 4px 6px;
      font-size: 9px;
    }

    .panel-footer .editor-stats {
      gap: 4px;
    }

    .panel-footer .editor-stats span {
      padding: 1px 3px;
      font-size: 8px;
    }

    .panel-footer .theme-btn {
      padding: 2px 4px;
      font-size: 8px;
    }
  }
</style>
