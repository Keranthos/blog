<template>
  <div class="question-box-view">
    <NavBar />
    <div class="header">
      <img
        src="https://picsum.photos/id/1043/1920/1080"
        alt="Background Image" class="header-image"
      />
    </div>
    <div class="content">
      <div class="question-list">
        <div v-for="question in paginatedQuestions" :key="question.ID" class="question-item">
          <div class="question">
            <p class="author">
              {{ question.Author }}
            </p>
            <p class="question-content">{{ question.Content }}</p>
          </div>
          <div class="answer">
            <p v-if="question.Answer" class="answered">
              <span class="label">回答:</span>
              <span class="text">{{ question.Answer }}</span>
            </p>
            <p v-else class="unanswered">
              待回答
            </p>
            <div v-if="user.level >= 3&&!question.Answer" class="submit-question">
              <textarea v-model="answers[question.ID]" placeholder="输入你的回答..."></textarea>
              <button @click="handleSubmitAnswer(question.ID)">提交回答</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 翻页组件 -->
      <div v-if="questions.length > 0" class="pagination">
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
    <div v-if="user.level < 3" class="ask-question">
      <textarea v-model="newQuestion" placeholder="输入你的问题..."></textarea>
      <button @click="handleSubmitQuestion">提交问题</button>
    </div>
    <div v-if="isLoading" class="loading">
      <i class="fas fa-spinner fa-spin"></i> 正在加载更多问题...
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useStore } from 'vuex'
import { getQuestions } from '@/api/questionBox/browse'
import { submitQuestion, submitAnswer } from '@/api/questionBox/edit'
import { showErrorMessage } from '@/utils/waifuMessage'
import NavBar from '@/components/NavBar'

const store = useStore()
const user = store.state.user

const questions = ref([])
const newQuestion = ref('')
const answers = ref({})
const currentPage = ref(1)
const isLoading = ref(false)

// 翻页相关
const pageSize = ref(10) // 每页显示10个问题

// 翻页计算属性
const totalPages = computed(() => Math.ceil(questions.value.length / pageSize.value))
const paginatedQuestions = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return questions.value.slice(start, end)
})

// 翻页方法
const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

const loadQuestions = async () => {
  if (isLoading.value) return
  isLoading.value = true
  try {
    // 一次性加载所有问题
    const response = await getQuestions(1)
    questions.value = response.questions || []
  } catch (error) {
    console.error('Failed to load questions:', error)
  } finally {
    isLoading.value = false
  }
}

// 删除滚动处理函数，使用翻页替代

const handleSubmitQuestion = async () => {
  try {
    await submitQuestion(user, newQuestion.value)
    newQuestion.value = ''
    currentPage.value = 1
    questions.value = []
    loadQuestions()
  } catch (error) {
    showErrorMessage(error)
  }
}

const handleSubmitAnswer = async (questionId) => {
  try {
    console.log('jwt:', store.state.token)
    await submitAnswer(user, questionId, answers.value[questionId])
    answers.value[questionId] = ''
    currentPage.value = 1
    questions.value = []
    loadQuestions()
  } catch (error) {
    console.error('Failed to submit answer:', error)
  }
}

onMounted(() => {
  loadQuestions()
})
</script>

<style scoped>
.question-box-view {
  min-height: 100vh;
  padding-top: 100px;
  padding-bottom: 200px;
}

.header {
  position: relative;
  width: 100%;
  margin-bottom: 60px;
  overflow: hidden;
}

.header-image {
  width: 90%;
  max-width: 1200px;
  height: 350px;
  object-fit: cover;
  border-radius: 30px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  display: block;
  margin: 0 auto;
}

.content {
  width: 90%;
  max-width: 900px;
  margin: 0 auto;
}

.question-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 40px;
}

/* 问题卡片 */
.question-item {
  background: white;
  border-radius: 16px;
  padding: 25px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(102, 126, 234, 0.1);
}

.question-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
  border-color: rgba(102, 126, 234, 0.3);
}

.question {
  padding: 20px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 12px;
  border-left: 4px solid #667eea;
  margin-bottom: 20px;
}

.author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #667eea;
  font-size: 1rem;
  margin-bottom: 10px;
}

.author::before {
  content: '👤';
  font-size: 1.1rem;
}

.question-content {
  font-size: 1.05rem;
  color: #333;
  line-height: 1.7;
  word-wrap: break-word;
}

.answer {
  margin-top: 15px;
  word-wrap: break-word;
}

.answered {
  padding: 20px;
  background: linear-gradient(135deg, rgba(40, 167, 69, 0.05) 0%, rgba(52, 211, 153, 0.05) 100%);
  border-radius: 12px;
  border-left: 4px solid #28a745;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.answered .label {
  font-weight: 600;
  color: #28a745;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 6px;
}

.answered .label::before {
  content: '✓';
  font-size: 1.2rem;
}

.answered .text {
  font-size: 1.05rem;
  color: #333;
  line-height: 1.7;
  word-wrap: break-word;
}

.unanswered {
  padding: 15px 20px;
  background: linear-gradient(135deg, rgba(255, 193, 7, 0.1) 0%, rgba(255, 152, 0, 0.1) 100%);
  border-radius: 12px;
  border-left: 4px solid #ffc107;
  color: #f57c00;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}

.unanswered::before {
  content: '⏳';
  font-size: 1.2rem;
}

/* 提问区 */
.ask-question {
  position: fixed;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%);
  width: 90%;
  max-width: 900px;
  padding: 25px 30px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  box-shadow: 0 -8px 40px rgba(0, 0, 0, 0.15);
  border-radius: 20px;
  z-index: 100;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.ask-question textarea {
  width: 100%;
  padding: 15px;
  border: 2px solid #e8eaf0;
  border-radius: 12px;
  resize: vertical;
  min-height: 80px;
  font-size: 1rem;
  font-family: inherit;
  transition: all 0.3s ease;
}

.ask-question textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.ask-question button {
  margin-top: 15px;
  padding: 12px 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 25px;
  color: white;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  display: flex;
  align-items: center;
  gap: 8px;
}

.ask-question button::before {
  content: '✉️';
}

.ask-question button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(102, 126, 234, 0.5);
}

.submit-question textarea {
  width: 100%;
  padding: 15px;
  margin-top: 15px;
  border: 2px solid #e8eaf0;
  border-radius: 12px;
  resize: vertical;
  min-height: 100px;
  font-size: 1rem;
}

.loading {
  text-align: center;
  margin: 40px 0;
  color: #667eea;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

/* 响应式 */
@media (max-width: 768px) {
  .question-box-view {
    padding-top: 80px;
  }

  .header-image {
    height: 250px;
  }

  .content {
    width: 95%;
  }

  .question-item {
    padding: 20px;
  }

  .ask-question {
    width: 95%;
    padding: 20px;
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
</style>
