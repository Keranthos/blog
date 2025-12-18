<template>
  <div class="comment-section">
    <div class="comment-header">
      <h2 class="section-title">
        <svg class="comment-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
        </svg>
        评论
        <span v-if="comments.length" class="comment-count">({{ comments.length }})</span>
      </h2>
    </div>

    <!-- 发表评论 -->
    <div class="comment-form">
      <div class="form-content">
        <div class="comment-input-wrapper">
          <textarea
            ref="commentTextarea"
            v-model="newComment"
            placeholder="支持 Markdown 语法"
            rows="4"
            :maxlength="300"
            @focus="handleInputFocus"
          ></textarea>
          <button class="emoji-btn" @click="toggleEmojiPicker" type="button" title="插入表情">
            😊
          </button>
        </div>
        <div class="form-actions">
          <div v-if="replyingTo" class="reply-indicator">
            <span class="reply-label">回复 @{{ getReplyTargetName() }}</span>
            <button class="cancel-reply-btn" @click="cancelReply">
              <font-awesome-icon icon="times" />
            </button>
          </div>
          <div class="action-buttons">
            <span class="char-count">{{ newComment.length }}/300</span>
            <button class="preview-btn" @click="togglePreview">
              <font-awesome-icon icon="eye" />
              预览
            </button>
            <button class="submit-btn" :disabled="!newComment.trim()" @click="submitComment">
              <font-awesome-icon icon="paper-plane" />
              发布
            </button>
          </div>
        </div>
        <!-- 预览评论 -->
        <div v-if="previewVisible" class="preview-comment">
          <div class="preview-label">预览效果</div>
          <div class="comment-item preview-item">
            <div class="comment-avatar">
              <div class="avatar-circle">{{ (user.value?.name || user.value?.username) ? (user.value?.name || user.value?.username).charAt(0) : 'U' }}</div>
            </div>
            <div class="comment-content">
              <div class="comment-header-info">
                <span class="comment-author">{{ user.value?.name || user.value?.username || '当前用户' }}</span>
                <span v-if="replyingTo" class="reply-tag">@{{ getReplyTargetName() }}</span>
                <span class="comment-time">刚刚</span>
              </div>
              <div v-if="newComment.trim()" class="comment-bubble">
                <div class="comment-text" v-html="renderedPreview"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Emoji 选择器 -->
    <EmojiPicker
      :visible="emojiPickerVisible"
      @select="insertEmoji"
      @close="emojiPickerVisible = false"
    />

    <!-- 评论列表 -->
    <div v-if="Array.isArray(comments) && comments.length > 0" class="comments-list">
      <div v-for="comment in getAllCommentsInOrder()" :key="comment.ID" class="comment-item" :class="{ 'reply-comment': comment.parent_id }">
        <div class="comment-avatar">
          <div class="avatar-circle" :class="{ 'small': comment.parent_id }">{{ comment.username.charAt(0) }}</div>
        </div>
        <div class="comment-content" :class="{ 'selected': replyingTo === comment.ID }">
          <div class="comment-header-info">
            <span class="comment-author">{{ comment.username }}</span>
            <span v-if="comment.parent_id" class="reply-tag">@{{ getParentCommentUsername(comment.parent_id) }}</span>
            <span class="comment-time">{{ formatCommentTime(comment.CreatedAt) }}</span>
            <button class="hover-reply-btn" @click="startReply(comment.ID, comment.username)">
              <font-awesome-icon icon="reply" />
              回复
            </button>
          </div>
          <div class="comment-bubble">
            <div class="comment-text">{{ comment.content }}</div>
          </div>
          <div class="comment-actions">
            <button
              v-if="user.level >= 3"
              class="delete-btn"
              title="删除评论"
              @click="handleDeleteComment(comment.ID)"
            >
              <font-awesome-icon icon="trash" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, computed } from 'vue'
import { marked } from 'marked'
import { useStore } from 'vuex'
import { getCommentsByID } from '@/api/Comments/browse'
import { createComment, deleteComment as deleteCommentAPI } from '@/api/Comments/edit'
import { showErrorMessage, showSuccessMessage } from '@/utils/waifuMessage'
import { protectLatex, restoreAndRenderLatex } from '@/utils/latex'
import EmojiPicker from '@/components/EmojiPicker.vue'

const props = defineProps({
  type: String, // 'blog', 'project', 'research', 'moment'
  articleId: String
})

const emit = defineEmits(['commentAdded', 'commentDeleted'])

const store = useStore()
const user = computed(() => store.state.user)

const comments = ref([])
const newComment = ref('')
const previewVisible = ref(false)
// 内嵌回复表单已移除，不再单独维护replyContent
const replyingTo = ref(null)
const replyTargetName = ref('')
const emojiPickerVisible = ref(false)
const commentTextarea = ref(null)

// 加载评论
const loadComments = async () => {
  try {
    const id = props.articleId
    comments.value = (await getCommentsByID(props.type, id)).data
  } catch (error) {
    showErrorMessage(error)
  }
}

// 提交评论
const submitComment = async () => {
  if (!newComment.value.trim()) return

  const id = props.articleId

  if (!user.value || !user.value.isLogged || !store.state.token) {
    showErrorMessage('401')
    return
  }

  try {
    // 如果是回复，使用parentId
    const parentId = replyingTo.value || null
    await createComment(user.value, id, props.type, newComment.value, parentId, store.state.token)

    newComment.value = ''
    cancelReply() // 清空回复状态
    await loadComments()
    showSuccessMessage('comment')
    emit('commentAdded')
  } catch (error) {
    showErrorMessage(error)
  }
}

// 开始回复
const startReply = (commentId, username) => {
  replyingTo.value = commentId
  replyTargetName.value = username
  // 聚焦到输入框
  nextTick(() => {
    const textarea = document.querySelector('.comment-form textarea')
    if (textarea) {
      textarea.focus()
    }
  })
}

// 取消回复
const cancelReply = () => {
  replyingTo.value = null
  replyTargetName.value = ''
}

// 获取回复目标名称
const getReplyTargetName = () => {
  return replyTargetName.value
}

// 处理输入框聚焦
const handleInputFocus = () => {
  // 如果当前没有在回复状态，清空回复状态
  if (!replyingTo.value) {
    cancelReply()
  }
}

// 切换 Emoji 选择器
const toggleEmojiPicker = () => {
  emojiPickerVisible.value = !emojiPickerVisible.value
}

// 插入 Emoji
const insertEmoji = (emoji) => {
  const textarea = commentTextarea.value
  if (!textarea) {
    // 如果没有 textarea ref，直接追加到末尾
    newComment.value += emoji
    return
  }

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = newComment.value

  // 在光标位置插入 Emoji
  newComment.value = text.substring(0, start) + emoji + text.substring(end)

  // 移动光标到插入位置之后
  nextTick(() => {
    textarea.focus()
    const newPosition = start + emoji.length
    textarea.setSelectionRange(newPosition, newPosition)
  })

  // 关闭选择器
  emojiPickerVisible.value = false
}

// 预览切换
const togglePreview = () => {
  previewVisible.value = !previewVisible.value
}
// 预览渲染（Markdown -> HTML）
const renderedPreview = computed(() => {
  if (!newComment.value) return ''

  // 预处理：保护 LaTeX 公式（在 marked 渲染前处理，避免公式中的特殊字符被 Markdown 解析）
  const { protected: protectedContent, placeholders: latexPlaceholders } = protectLatex(newComment.value)

  // 处理换行符，将 \n 转换为 <br>（在保护 LaTeX 之后，避免影响公式）
  const content = protectedContent.replace(/\n/g, '<br>')

  // 然后使用 marked 渲染 Markdown
  let html = marked(content, {
    breaks: false, // 我们已经手动处理了换行
    gfm: true,
    headerIds: false,
    mangle: false,
    pedantic: false,
    sanitize: false,
    smartLists: true,
    smartypants: false
  })

  // 恢复并渲染 LaTeX 公式
  html = restoreAndRenderLatex(html, latexPlaceholders)

  return html
})

// 获取所有评论按正确顺序排列（回复紧跟父评论）
const getAllCommentsInOrder = () => {
  const result = []

  // 遍历所有顶级评论
  for (const comment of comments.value) {
    // 添加顶级评论
    result.push(comment)

    // 添加该评论的所有回复（按时间倒序）
    if (comment.replies && comment.replies.length > 0) {
      // 对回复按时间排序（最新的在前）
      const sortedReplies = [...comment.replies].sort((a, b) =>
        new Date(b.CreatedAt) - new Date(a.CreatedAt)
      )
      result.push(...sortedReplies)
    }
  }

  return result
}

// 获取父评论用户名
const getParentCommentUsername = (parentId) => {
  // 在所有评论中查找父评论（包括顶级评论和回复）
  const allComments = getAllCommentsInOrder()
  for (const comment of allComments) {
    if (comment.ID === parentId) {
      return comment.username
    }
  }
  return '未知用户'
}

// 旧的内嵌回复表单已移除，此函数不再使用

// 删除评论
const handleDeleteComment = async (commentId) => {
  // 第一次确认
  const firstConfirm = confirm('确定要删除这条评论吗？\n\n⚠️ 此操作不可撤销！')
  if (!firstConfirm) return

  // 第二次确认（防误触）
  const secondConfirm = confirm('再次确认：真的要删除这条评论吗？\n\n删除后将无法恢复！')
  if (!secondConfirm) return

  try {
    const user = store.state.user
    await deleteCommentAPI(user, commentId, store.state.token)
    await loadComments()

    // 显示看板娘消息（如果可用）
    if (window.showMessage) {
      window.showMessage('评论删除成功～', 3000, 9)
    } else {
      showSuccessMessage('delete')
    }

    emit('commentDeleted')
  } catch (error) {
    // 显示看板娘错误消息（如果可用）
    if (window.showMessage) {
      window.showMessage('(｡•́︿•̀｡)<br>删除失败了…请重试吧～', 5000, 10)
    } else {
      showErrorMessage(error)
    }
  }
}

// 格式化评论时间
const formatCommentTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date

  if (diff < 60000) { // 1分钟内
    return '刚刚'
  } else if (diff < 3600000) { // 1小时内
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (diff < 86400000) { // 24小时内
    return `${Math.floor(diff / 3600000)}小时前`
  } else if (diff < 604800000) { // 7天内
    return `${Math.floor(diff / 86400000)}天前`
  } else {
    return date.toLocaleDateString('zh-CN')
  }
}

// 暴露方法给父组件
defineExpose({
  loadComments
})
</script>

<style scoped>
.comment-section {
  margin-top: 3rem;
  padding-top: 2rem;
  background: transparent;
}

.comment-header {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.5rem;
  color: #333;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
}

.comment-icon {
  width: 24px;
  height: 24px;
  color: #8b5cf6;
}

.comment-count {
  font-size: 1rem;
  color: #999;
  font-weight: 400;
}

/* 评论列表 */
.comments-list {
  margin-bottom: 2rem;
}

.comment-item {
  display: flex;
  gap: 0.8rem;
  margin-bottom: 2rem;
  padding: 0;
  position: relative;
  align-items: flex-start;
}

.comment-avatar {
  flex-shrink: 0;
  margin-top: 0.2rem;
}

.avatar-circle {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #2d4a6b;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1.1rem;
  box-shadow: 0 2px 8px rgba(45, 74, 107, 0.3);
  transition: transform 0.2s ease;
}

.avatar-circle:hover {
  transform: scale(1.05);
}

.avatar-circle.small {
  width: 36px;
  height: 36px;
  font-size: 0.9rem;
}

.comment-content {
  flex: 1;
  min-width: 0;
  max-width: calc(100% - 60px);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.comment-bubble {
  background: white;
  border-radius: 12px;
  padding: 0.85rem 1.05rem;
  border-top: 1.5px solid rgba(139, 92, 246, 0.4);
  border-left: 1.5px solid rgba(139, 92, 246, 0.4);
  border-bottom: 1px solid #f0f0f0;
  border-right: 1px solid #f0f0f0;
  box-shadow:
    0 4px 12px rgba(139, 92, 246, 0.2),
    0 2px 6px rgba(139, 92, 246, 0.15),
    0 1px 2px rgba(0, 0, 0, 0.05);
  position: relative;
  word-wrap: break-word;
  display: inline-block;
  max-width: 100%;
  margin-top: 0.3rem;
  transition: box-shadow 0.2s ease, border-color 0.2s ease;
}

/* 箭头外层边框（包含顶部紫色边框） */
.comment-bubble::after {
  content: '';
  position: absolute;
  left: -9px;
  top: 11px;
  width: 0;
  height: 0;
  border-top: 9px solid transparent;
  border-bottom: 9px solid transparent;
  border-right: 9px solid rgba(139, 92, 246, 0.4);
  z-index: 0;
}

/* 箭头内层白色填充 */
.comment-bubble::before {
  content: '';
  position: absolute;
  left: -8px;
  top: 12px;
  width: 0;
  height: 0;
  border-top: 8px solid transparent;
  border-bottom: 8px solid transparent;
  border-right: 8px solid white;
  z-index: 1;
}

.comment-header-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
  flex-wrap: wrap;
  position: relative;
}

.hover-reply-btn {
  opacity: 0;
  padding: 0.2rem 0.5rem;
  background: #8b5cf6;
  border: none;
  border-radius: 6px;
  color: white;
  cursor: pointer;
  font-size: 0.75rem;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  margin-left: auto;
}

.comment-item:hover .hover-reply-btn {
  opacity: 1;
}

.comment-content.selected {
  background: rgba(139, 92, 246, 0.1);
  border-radius: 8px;
  padding: 0.5rem;
  margin: -0.5rem;
}

.comment-author {
  font-weight: 600;
  color: #333;
  font-size: 0.9rem;
}

.comment-time {
  color: #999;
  font-size: 0.8rem;
}

.reply-tag {
  background: #8b5cf6;
  color: white;
  padding: 0.15rem 0.5rem;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 500;
  margin: 0 0.3rem;
  display: inline-block;
  white-space: nowrap;
}

.comment-actions {
  margin-top: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.delete-btn {
  padding: 0.2rem 0.4rem;
  background: #ff4757;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.7rem;
  transition: background 0.3s ease;
}

.delete-btn:hover {
  background: #ff3742;
}

.comment-text {
  color: #2b2b2b;
  line-height: 1.6;
  font-size: 0.95rem;
  margin: 0;
  text-align: left;
}

/* 内联代码样式 */
.comment-text :deep(code) {
  background: #f1f5f9;
  color: #475569;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 0.85em;
  border: 1px solid #e2e8f0;
  font-weight: 500;
}

/* 代码块样式 */
.comment-text :deep(pre) {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 1rem;
  margin: 0.5rem 0;
  overflow-x: auto;
}

.comment-text :deep(pre code) {
  background: none;
  border: none;
  padding: 0;
  color: #1e293b;
}

.comment-actions {
  margin-top: 0.3rem;
}

.reply-btn {
  padding: 0.25rem 0.6rem;
  background: #8b5cf6;
  border: none;
  border-radius: 6px;
  color: white;
  cursor: pointer;
  font-size: 0.75rem;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  margin-top: 0.3rem;
}

.reply-btn:hover {
  background: #7c3aed;
  transform: translateY(-1px);
}

/* 回复表单 */
.reply-form {
  margin-top: 0.8rem;
  padding: 0.8rem;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  border-left: 3px solid #8b5cf6;
}

.reply-form textarea {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  resize: vertical;
  font-family: inherit;
  margin-bottom: 1rem;
}

.reply-actions {
  display: flex;
  gap: 0.8rem;
  justify-content: flex-end;
}

.submit-reply-btn {
  padding: 0.6rem 1.2rem;
  background: #a855f7;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.submit-reply-btn:hover:not(:disabled) {
  background: #5a6fd8;
}

.submit-reply-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.cancel-reply-btn {
  padding: 0.6rem 1.2rem;
  background: #f5f5f5;
  color: #666;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.cancel-reply-btn:hover {
  background: #e0e0e0;
}

/* 回复评论样式 - 缩进一个头像宽度 */
.reply-comment {
  margin-left: 60px; /* 一个头像宽度 + 间距 */
}

.reply-comment .comment-avatar {
  margin-top: 0.1rem; /* 稍微调整头像位置 */
}

/* 发表评论表单 */
.comment-form {
  background: transparent;
  border-radius: 0;
  padding: 2rem 0;
  border: none;
  box-shadow: none;
}

.form-title {
  font-size: 1.2rem;
  color: #333;
  margin: 0 0 1.5rem 0;
}

.comment-input-wrapper {
  position: relative;
  width: 100%;
  margin-bottom: 1rem;
}

.form-content textarea {
  width: 100%;
  padding: 1rem 50px 1rem 1rem;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  resize: vertical;
  font-family: inherit;
  background: white;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
  font-size: 0.95rem;
  line-height: 1.5;
}

.emoji-btn {
  position: absolute;
  right: 12px;
  bottom: 12px;
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s;
  line-height: 1;
  font-family: "Apple Color Emoji", "Segoe UI Emoji", "Noto Color Emoji", "Segoe UI Symbol", "Android Emoji", "EmojiSymbols", "EmojiOne Mozilla", "Twemoji Mozilla", "Segoe UI", sans-serif;
}

.emoji-btn:hover {
  background: #f3f4f6;
  transform: scale(1.1);
}

.emoji-btn:active {
  transform: scale(0.95);
}

.form-content textarea:focus {
  outline: none;
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
}

.form-actions {
  display: flex;
  align-items: center;
  margin-top: 1rem;
}

.reply-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: #8b5cf6;
  border-radius: 8px;
  color: white;
}

.reply-label {
  font-size: 0.9rem;
  font-weight: 500;
}

.cancel-reply-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 4px;
  color: white;
  cursor: pointer;
  padding: 0.2rem 0.4rem;
  font-size: 0.8rem;
  transition: background 0.3s ease;
}

.cancel-reply-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-left: auto;
}

.char-count {
  color: #999;
  font-size: 0.9rem;
}

.preview-btn {
  padding: 0.8rem 2rem;
  background: #8b5cf6;
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  box-shadow: 0 2px 4px rgba(139, 92, 246, 0.3);
}

.preview-btn:hover {
  background: #7c3aed;
}

.submit-btn {
  padding: 0.8rem 2rem;
  background: linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%);
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  box-shadow: 0 2px 4px rgba(139, 92, 246, 0.3);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(139, 92, 246, 0.4);
}

.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

/* 预览评论样式 */
.preview-comment {
  margin-top: 1rem;
  padding: 1rem;
  background: rgba(139, 92, 246, 0.05);
  border: 2px dashed rgba(139, 92, 246, 0.3);
  border-radius: 12px;
}

.preview-label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #8b5cf6;
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.preview-item {
  opacity: 0.8;
  transform: scale(0.98);
  transition: all 0.3s ease;
}

.preview-item:hover {
  opacity: 1;
  transform: scale(1);
}

@media (max-width: 768px) {
  .comment-item {
    flex-direction: column;
    gap: 1rem;
  }

  .comment-avatar {
    align-self: flex-start;
  }

  .replies {
    padding-left: 1rem;
  }

  .reply-item {
    flex-direction: column;
    gap: 0.8rem;
  }

  .reply-actions {
    flex-direction: column;
  }

  .form-actions {
    justify-content: center;
  }
}
</style>
