<template>
  <div class="comment-section">
    <div class="comment-header">
      <h2 class="section-title">
        评论区
        <span v-if="comments.length" class="comment-count">({{ comments.length }})</span>
      </h2>
    </div>

    <!-- 评论列表 -->
    <div v-if="Array.isArray(comments) && comments.length > 0" class="comments-list">
      <div v-for="comment in comments" :key="comment.ID" class="comment-item">
        <div class="comment-avatar">
          <div class="avatar-circle">{{ comment.username.charAt(0) }}</div>
        </div>
        <div class="comment-content">
          <div class="comment-header-info">
            <span class="comment-author">{{ comment.username }}</span>
            <span class="comment-time">{{ formatCommentTime(comment.CreatedAt) }}</span>
            <button
              v-if="user.level >= 3"
              class="delete-btn"
              title="删除评论"
              @click="handleDeleteComment(comment.ID)"
            >
              <font-awesome-icon icon="trash" />
            </button>
          </div>
          <div class="comment-text">{{ comment.content }}</div>
          <div class="comment-actions">
            <button class="reply-btn" @click="toggleReply(comment.ID)">
              <font-awesome-icon icon="reply" />
              回复
            </button>
          </div>

          <!-- 回复输入框 -->
          <div v-if="replyingTo === comment.ID" class="reply-form">
            <textarea
              v-model="replyContent"
              placeholder="写下你的回复..."
              rows="3"
            ></textarea>
            <div class="reply-actions">
              <button class="submit-reply-btn" :disabled="!replyContent.trim()" @click="submitReply(comment.ID)">
                发表回复
              </button>
              <button class="cancel-reply-btn" @click="cancelReply">
                取消
              </button>
            </div>
          </div>

          <!-- 子评论（回复） -->
          <div v-if="comment.replies && comment.replies.length > 0" class="replies">
            <div v-for="reply in comment.replies" :key="reply.ID" class="reply-item">
              <div class="reply-avatar">
                <div class="avatar-circle small">{{ reply.username.charAt(0) }}</div>
              </div>
              <div class="reply-content">
                <div class="reply-header-info">
                  <span class="reply-author">{{ reply.username }}</span>
                  <span class="reply-time">{{ formatCommentTime(reply.CreatedAt) }}</span>
                  <button
                    v-if="user.level >= 3"
                    class="delete-btn"
                    title="删除回复"
                    @click="handleDeleteComment(reply.ID)"
                  >
                    <font-awesome-icon icon="trash" />
                  </button>
                </div>
                <div class="reply-text">{{ reply.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 发表评论 -->
    <div class="comment-form">
      <h3 class="form-title">发表评论</h3>
      <div class="form-content">
        <textarea
          v-model="newComment"
          placeholder="写下你的想法..."
          rows="4"
        ></textarea>
        <div class="form-actions">
          <button class="submit-btn" :disabled="!newComment.trim()" @click="submitComment">
            <font-awesome-icon icon="paper-plane" />
            发表评论
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useStore } from 'vuex'
import { getCommentsByID } from '@/api/Comments/browse'
import { createComment, deleteComment as deleteCommentAPI } from '@/api/Comments/edit'
import { showErrorMessage, showSuccessMessage } from '@/utils/waifuMessage'

const props = defineProps({
  type: String, // 'blog', 'project', 'research', 'moment'
  articleId: String
})

const emit = defineEmits(['commentAdded', 'commentDeleted'])

const store = useStore()
const user = store.state.user

const comments = ref([])
const newComment = ref('')
const replyContent = ref('')
const replyingTo = ref(null)

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

  if (!user || !user.token) {
    showErrorMessage('401')
    return
  }

  try {
    await createComment(user, id, props.type, newComment.value)
    newComment.value = ''
    await loadComments()
    showSuccessMessage('comment')
    emit('commentAdded')
  } catch (error) {
    showErrorMessage(error)
  }
}

// 切换回复状态
const toggleReply = (commentId) => {
  if (replyingTo.value === commentId) {
    replyingTo.value = null
    replyContent.value = ''
  } else {
    replyingTo.value = commentId
    replyContent.value = ''
  }
}

// 取消回复
const cancelReply = () => {
  replyingTo.value = null
  replyContent.value = ''
}

// 提交回复
const submitReply = async (parentId) => {
  if (!replyContent.value.trim()) return

  const id = props.articleId

  if (!user || !user.token) {
    showErrorMessage('401')
    return
  }

  try {
    await createComment(user, id, props.type, replyContent.value, parentId)
    replyContent.value = ''
    replyingTo.value = null
    await loadComments()
    showSuccessMessage('comment')
    emit('commentAdded')
  } catch (error) {
    showErrorMessage(error)
  }
}

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
    await deleteCommentAPI(user, commentId)
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
  border-top: 2px solid #f0f0f0;
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
  gap: 1rem;
  margin-bottom: 2rem;
  padding: 1.5rem;
  background: #f9f9f9;
  border-radius: 12px;
  border: 1px solid #e0e0e0;
}

.comment-avatar {
  flex-shrink: 0;
}

.avatar-circle {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1.2rem;
}

.avatar-circle.small {
  width: 35px;
  height: 35px;
  font-size: 1rem;
}

.comment-content {
  flex: 1;
}

.comment-header-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 0.8rem;
}

.comment-author {
  font-weight: 600;
  color: #333;
}

.comment-time {
  color: #999;
  font-size: 0.9rem;
}

.delete-btn {
  margin-left: auto;
  padding: 0.3rem 0.6rem;
  background: #ff4757;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.8rem;
  transition: background 0.3s ease;
}

.delete-btn:hover {
  background: #ff3742;
}

.comment-text {
  color: #333;
  line-height: 1.6;
  margin-bottom: 1rem;
}

.comment-actions {
  margin-bottom: 1rem;
}

.reply-btn {
  padding: 0.4rem 0.8rem;
  background: #f5f5f5;
  border: none;
  border-radius: 6px;
  color: #666;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.3rem;
}

.reply-btn:hover {
  background: #e0e0e0;
  color: #333;
}

/* 回复表单 */
.reply-form {
  margin-top: 1rem;
  padding: 1rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #e0e0e0;
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
  background: #667eea;
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

/* 回复列表 */
.replies {
  margin-top: 1.5rem;
  padding-left: 2rem;
  border-left: 2px solid #e0e0e0;
}

.reply-item {
  display: flex;
  gap: 0.8rem;
  margin-bottom: 1rem;
  padding: 1rem;
  background: white;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
}

.reply-content {
  flex: 1;
}

.reply-header-info {
  display: flex;
  align-items: center;
  gap: 0.8rem;
  margin-bottom: 0.5rem;
}

.reply-author {
  font-weight: 600;
  color: #333;
  font-size: 0.9rem;
}

.reply-time {
  color: #999;
  font-size: 0.8rem;
}

.reply-text {
  color: #333;
  line-height: 1.5;
  font-size: 0.9rem;
}

/* 发表评论表单 */
.comment-form {
  background: #f9f9f9;
  border-radius: 12px;
  padding: 2rem;
  border: 1px solid #e0e0e0;
}

.form-title {
  font-size: 1.2rem;
  color: #333;
  margin: 0 0 1.5rem 0;
}

.form-content textarea {
  width: 100%;
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  resize: vertical;
  font-family: inherit;
  margin-bottom: 1rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.submit-btn {
  padding: 0.8rem 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: transform 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}

.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
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
