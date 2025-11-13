<template>
  <div class="day-view">
    <div class="day-header">
      <h2 class="day-title">{{ formatDate(currentDate) }}</h2>
      <div class="day-actions">
        <button class="btn btn-primary" @click="showAddTaskModal = true">
          <font-awesome-icon icon="plus" /> 添加任务
        </button>
        <button class="btn btn-secondary" @click="showAddEntertainmentModal = true">
          <font-awesome-icon icon="plus" /> 添加娱乐
        </button>
      </div>
    </div>

    <div class="day-content">
      <!-- 左侧：任务和娱乐 -->
      <div class="day-main">
        <!-- 今日计划 -->
        <div class="section">
          <div class="section-header">
            <h3 class="section-title">
              <font-awesome-icon icon="list" /> 今日计划
              <span class="progress-text">({{ completedCount }}/{{ tasks.length }})</span>
            </h3>
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: progressPercentage + '%' }"></div>
            </div>
          </div>
          <div class="tasks-list">
            <div
              v-for="task in tasks"
              :key="task.ID"
              class="task-item"
              :class="{ completed: task.completed }"
            >
              <input
                type="checkbox"
                :checked="task.completed"
                class="task-checkbox"
                @change="toggleTask(task)"
              />
              <div class="task-content">
                <div class="task-title">{{ task.title }}</div>
                <div v-if="task.description" class="task-description">{{ task.description }}</div>
                <div class="task-meta">
                  <span class="priority" :class="`priority-${task.priority}`">
                    {{ getPriorityText(task.priority) }}
                  </span>
                  <span v-if="task.estimatedTime" class="time">
                    ⏱️ {{ task.estimatedTime }}分钟
                  </span>
                </div>
              </div>
              <div class="task-actions">
                <button class="icon-btn" @click="editTask(task)">
                  <font-awesome-icon icon="pen-to-square" />
                </button>
                <button class="icon-btn" @click="deleteTask(task.ID)">
                  <font-awesome-icon icon="trash" />
                </button>
              </div>
            </div>
            <div v-if="tasks.length === 0" class="empty-state">
              暂无任务，点击"添加任务"开始计划今天吧！
            </div>
          </div>
        </div>

        <!-- 娱乐活动 -->
        <div class="section">
          <div class="section-header">
            <h3 class="section-title">
              <font-awesome-icon icon="gamepad" /> 娱乐活动
            </h3>
          </div>
          <div class="entertainments-list">
            <div
              v-for="ent in entertainments"
              :key="ent.ID"
              class="entertainment-item"
            >
              <div class="ent-icon">🎮</div>
              <div class="ent-content">
                <div class="ent-type">{{ ent.type }}</div>
                <div class="ent-meta">
                  <span>⏱️ {{ formatDuration(ent.duration) }}</span>
                  <span v-if="ent.cost > 0">💰 ¥{{ ent.cost.toFixed(2) }}</span>
                  <span v-if="ent.rating > 0">⭐ {{ ent.rating }}</span>
                </div>
                <div v-if="ent.notes" class="ent-notes">{{ ent.notes }}</div>
              </div>
              <div class="ent-actions">
                <button class="icon-btn" @click="editEntertainment(ent)">
                  <font-awesome-icon icon="pen-to-square" />
                </button>
                <button class="icon-btn" @click="deleteEntertainment(ent.ID)">
                  <font-awesome-icon icon="trash" />
                </button>
              </div>
            </div>
            <div v-if="entertainments.length === 0" class="empty-state">
              暂无娱乐活动
            </div>
          </div>
        </div>

        <!-- 今日反思 -->
        <div class="section">
          <div class="section-header">
            <h3 class="section-title">
              <font-awesome-icon icon="pen-to-square" /> 今日反思
            </h3>
          </div>
          <div class="reflection-editor">
            <textarea
              v-if="canEditReflection && isEditingReflection"
              v-model="reflectionContent"
              placeholder="记录今天的思考、感悟、收获..."
              class="reflection-textarea"
              rows="6"
            ></textarea>
            <div
              v-else
              class="reflection-readonly"
            >
              {{ reflectionContent || '暂无反思内容' }}
            </div>
            <button
              v-if="canEditReflection"
              class="btn btn-primary"
              @click="isEditingReflection ? saveReflection() : editReflection()"
            >
              {{ isEditingReflection ? '保存反思' : '编辑' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 右侧：侧边栏 -->
      <div class="day-sidebar">
        <!-- 日历 -->
        <div class="sidebar-section">
          <h4 class="sidebar-title">日历</h4>
          <div class="mini-calendar">
            <div class="calendar-header">
              <button class="calendar-nav" @click="prevMonth">‹</button>
              <span class="calendar-month">{{ currentMonthYear }}</span>
              <button class="calendar-nav" @click="nextMonth">›</button>
            </div>
            <div class="calendar-grid">
              <div v-for="day in weekDays" :key="day" class="calendar-weekday">{{ day }}</div>
              <div
                v-for="date in calendarDays"
                :key="date.key"
                class="calendar-day"
                :class="{ today: date.isToday, selected: date.isSelected, other: date.isOtherMonth }"
                @click="selectDate(date.date)"
              >
                {{ date.day }}
              </div>
            </div>
          </div>
        </div>

        <!-- 本周概览 -->
        <div class="sidebar-section">
          <h4 class="sidebar-title">本周任务</h4>
          <div class="week-overview">
            <div v-for="day in weekTasks" :key="day.date" class="week-day">
              <div class="week-day-name">{{ day.name }}</div>
              <div class="week-day-count">{{ day.count }}个任务</div>
            </div>
          </div>
        </div>

        <!-- 即将到来的DDL -->
        <div class="sidebar-section">
          <h4 class="sidebar-title">即将到来的DDL</h4>
          <div class="upcoming-deadlines">
            <div
              v-for="ddl in upcomingDeadlines"
              :key="ddl.ID"
              class="deadline-item"
            >
              <div class="deadline-title">{{ ddl.title }}</div>
              <div class="deadline-time">{{ formatDeadlineTime(ddl.dueDate, ddl.dueTime) }}</div>
            </div>
            <div v-if="upcomingDeadlines.length === 0" class="empty-state-small">
              暂无即将到来的DDL
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加任务模态框 -->
    <TaskModal
      v-if="showAddTaskModal"
      :task="editingTask"
      @close="closeTaskModal"
      @save="handleTaskSave"
    />

    <!-- 添加娱乐模态框 -->
    <EntertainmentModal
      v-if="showAddEntertainmentModal"
      :entertainment="editingEntertainment"
      @close="closeEntertainmentModal"
      @save="handleEntertainmentSave"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import {
  getTasks,
  createTask,
  updateTask,
  deleteTask as deleteTaskAPI,
  getEntertainments,
  createEntertainment,
  updateEntertainment,
  deleteEntertainment as deleteEntertainmentAPI,
  getReflection,
  createOrUpdateReflection,
  getDeadlines
} from '@/api/planner'
import TaskModal from './TaskModal.vue'
import EntertainmentModal from './EntertainmentModal.vue'

const props = defineProps({
  currentDate: {
    type: Date,
    default: () => new Date()
  }
})

const tasks = ref([])
const allWeekTasks = ref([])
const entertainments = ref([])
const reflectionContent = ref('')
const currentReflection = ref(null)
const isEditingReflection = ref(false)
const showAddTaskModal = ref(false)
const showAddEntertainmentModal = ref(false)
const editingTask = ref(null)
const editingEntertainment = ref(null)
const selectedDate = ref(new Date())
const calendarMonth = ref(new Date())
const upcomingDeadlines = ref([])

const weekDays = ['一', '二', '三', '四', '五', '六', '日']

const completedCount = computed(() => {
  return tasks.value.filter(t => t.completed).length
})

const progressPercentage = computed(() => {
  if (tasks.value.length === 0) return 0
  return Math.round((completedCount.value / tasks.value.length) * 100)
})

const currentMonthYear = computed(() => {
  return `${calendarMonth.value.getFullYear()}年${calendarMonth.value.getMonth() + 1}月`
})

const calendarDays = computed(() => {
  const year = calendarMonth.value.getFullYear()
  const month = calendarMonth.value.getMonth()
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const firstDayWeek = (firstDay.getDay() + 6) % 7
  const daysInMonth = lastDay.getDate()

  const days = []

  // 上个月的最后几天
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  for (let i = firstDayWeek - 1; i >= 0; i--) {
    days.push({
      day: prevMonthLastDay - i,
      date: new Date(year, month - 1, prevMonthLastDay - i),
      isOtherMonth: true,
      isToday: false,
      isSelected: false,
      key: `prev-${prevMonthLastDay - i}`
    })
  }

  // 当前月的日期
  for (let i = 1; i <= daysInMonth; i++) {
    const date = new Date(year, month, i)
    const today = new Date()
    const isToday = date.toDateString() === today.toDateString()
    const isSelected = date.toDateString() === selectedDate.value.toDateString()
    days.push({
      day: i,
      date,
      isOtherMonth: false,
      isToday,
      isSelected,
      key: `current-${i}`
    })
  }

  // 下个月的前几天（填满6行）
  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    days.push({
      day: i,
      date: new Date(year, month + 1, i),
      isOtherMonth: true,
      isToday: false,
      isSelected: false,
      key: `next-${i}`
    })
  }

  return days
})

const weekTasks = computed(() => {
  // 返回本周7天的任务数量
  const week = []
  const today = new Date()
  const startOfWeek = new Date(today)
  const dayOffset = (today.getDay() + 6) % 7
  startOfWeek.setDate(today.getDate() - dayOffset)
  startOfWeek.setHours(0, 0, 0, 0)

  for (let i = 0; i < 7; i++) {
    const date = new Date(startOfWeek)
    date.setDate(startOfWeek.getDate() + i)
    const dateStr = date.toISOString().split('T')[0]

    // 计算这一天的任务数量
    const dayTasks = allWeekTasks.value.filter(t => {
      const taskDate = new Date(t.date).toISOString().split('T')[0]
      return taskDate === dateStr
    })

    week.push({
      name: weekDays[i],
      date: dateStr,
      count: dayTasks.length
    })
  }
  return week
})

onMounted(() => {
  loadData()
  loadUpcomingDeadlines()
  loadWeekTasks()
})

watch(() => props.currentDate, () => {
  selectedDate.value = new Date(props.currentDate)
  loadData()
}, { deep: true })

watch(selectedDate, () => {
  isEditingReflection.value = false // 切换日期时重置编辑状态
  loadData()
})

const loadData = async () => {
  const dateStr = formatDateForAPI(selectedDate.value)
  await Promise.all([
    loadTasks(dateStr),
    loadEntertainments(dateStr),
    loadReflection(dateStr)
  ])
}

const loadTasks = async (dateStr) => {
  try {
    const response = await getTasks(dateStr)
    if (response?.data) {
      tasks.value = response.data
    } else {
      tasks.value = []
    }
  } catch (error) {
    // 404 或其他错误时返回空数组
    if (error.response?.status === 404) {
      tasks.value = []
    } else {
      console.error('加载任务失败:', error)
      tasks.value = []
    }
  }
}

const loadEntertainments = async (dateStr) => {
  try {
    const response = await getEntertainments(dateStr)
    if (response?.data) {
      entertainments.value = response.data
    } else {
      entertainments.value = []
    }
  } catch (error) {
    // 404 或其他错误时返回空数组
    if (error.response?.status === 404) {
      entertainments.value = []
    } else {
      console.error('加载娱乐活动失败:', error)
      entertainments.value = []
    }
  }
}

const loadReflection = async (dateStr) => {
  try {
    const response = await getReflection(dateStr)
    if (response?.data && response.data.content) {
      currentReflection.value = response.data
      reflectionContent.value = response.data.content
      isEditingReflection.value = false // 有内容时默认是查看模式
    } else {
      reflectionContent.value = ''
      currentReflection.value = null
      // 如果没有内容且可以编辑，默认进入编辑模式
      isEditingReflection.value = canEditReflection.value
    }
  } catch (error) {
    // 404 或其他错误时清空内容
    if (error.response?.status === 404) {
      reflectionContent.value = ''
      currentReflection.value = null
      // 如果没有内容且可以编辑，默认进入编辑模式
      isEditingReflection.value = canEditReflection.value
    } else {
      console.error('加载反思失败:', error)
    }
  }
}

const loadUpcomingDeadlines = async () => {
  try {
    const response = await getDeadlines(false)
    if (response?.data) {
      upcomingDeadlines.value = response.data.slice(0, 5)
    } else {
      upcomingDeadlines.value = []
    }
  } catch (error) {
    // 404 或其他错误时返回空数组
    if (error.response?.status === 404) {
      upcomingDeadlines.value = []
    } else {
      console.error('加载DDL失败:', error)
      upcomingDeadlines.value = []
    }
  }
}

const loadWeekTasks = async () => {
  try {
    // 获取本周的开始日期（周一）
    const today = new Date()
    const startOfWeek = new Date(today)
    const dayOffset = (today.getDay() + 6) % 7
    startOfWeek.setDate(today.getDate() - dayOffset)
    startOfWeek.setHours(0, 0, 0, 0)

    // 获取本周的所有任务（不传日期参数会返回所有任务，然后在前端过滤）
    const response = await getTasks()
    if (response?.data) {
      // 过滤出本周的任务
      const endOfWeek = new Date(startOfWeek)
      endOfWeek.setDate(startOfWeek.getDate() + 6)
      endOfWeek.setHours(23, 59, 59, 999)

      allWeekTasks.value = response.data.filter(task => {
        const taskDate = new Date(task.date)
        return taskDate >= startOfWeek && taskDate <= endOfWeek
      })
    } else {
      allWeekTasks.value = []
    }
  } catch (error) {
    // 404 或其他错误时返回空数组
    if (error.response?.status === 404) {
      allWeekTasks.value = []
    } else {
      console.error('加载周任务失败:', error)
      allWeekTasks.value = []
    }
  }
}

const toggleTask = async (task) => {
  try {
    task.completed = !task.completed
    await updateTask(task.ID, task)
    await loadWeekTasks() // 刷新本周任务
  } catch (error) {
    console.error('更新任务失败:', error)
    task.completed = !task.completed // 回滚
  }
}

const editTask = (task) => {
  editingTask.value = { ...task }
  showAddTaskModal.value = true
}

const deleteTask = async (id) => {
  if (!confirm('确定要删除这个任务吗？')) return
  try {
    await deleteTaskAPI(id)
    await loadData()
    await loadWeekTasks() // 刷新本周任务
  } catch (error) {
    console.error('删除任务失败:', error)
  }
}

const editEntertainment = (ent) => {
  editingEntertainment.value = { ...ent }
  showAddEntertainmentModal.value = true
}

const deleteEntertainment = async (id) => {
  if (!confirm('确定要删除这个娱乐活动吗？')) return
  try {
    await deleteEntertainmentAPI(id)
    await loadData()
  } catch (error) {
    console.error('删除娱乐活动失败:', error)
  }
}

const saveReflection = async () => {
  try {
    const reflectionData = {
      date: formatDateForAPI(selectedDate.value),
      content: reflectionContent.value,
      tags: [] // 不再使用标签
    }
    await createOrUpdateReflection(reflectionData)
    await loadReflection(formatDateForAPI(selectedDate.value))
    isEditingReflection.value = false // 保存后切换到查看模式
  } catch (error) {
    console.error('保存反思失败:', error)
  }
}

const editReflection = () => {
  // 确保内容已加载，然后切换到编辑模式
  if (!reflectionContent.value) {
    loadReflection(formatDateForAPI(selectedDate.value)).then(() => {
      isEditingReflection.value = true
    })
  } else {
    isEditingReflection.value = true
  }
}

const closeTaskModal = () => {
  showAddTaskModal.value = false
  editingTask.value = null
}

const handleTaskSave = async (taskData) => {
  try {
    if (editingTask.value) {
      await updateTask(editingTask.value.ID, taskData)
    } else {
      await createTask(taskData)
    }
    await loadData()
    await loadWeekTasks() // 刷新本周任务
    closeTaskModal()
  } catch (error) {
    console.error('保存任务失败:', error)
  }
}

const closeEntertainmentModal = () => {
  showAddEntertainmentModal.value = false
  editingEntertainment.value = null
}

const handleEntertainmentSave = async (entertainmentData) => {
  try {
    if (editingEntertainment.value) {
      await updateEntertainment(editingEntertainment.value.ID, entertainmentData)
    } else {
      await createEntertainment(entertainmentData)
    }
    await loadData()
    closeEntertainmentModal()
  } catch (error) {
    console.error('保存娱乐活动失败:', error)
  }
}

const selectDate = (date) => {
  selectedDate.value = date
}

const prevMonth = () => {
  calendarMonth.value = new Date(calendarMonth.value.getFullYear(), calendarMonth.value.getMonth() - 1, 1)
}

const nextMonth = () => {
  calendarMonth.value = new Date(calendarMonth.value.getFullYear(), calendarMonth.value.getMonth() + 1, 1)
}

// 判断是否可以编辑反思（仅今天和昨天可编辑）
const canEditReflection = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)

  const selected = new Date(selectedDate.value)
  selected.setHours(0, 0, 0, 0)

  return selected.getTime() === today.getTime() || selected.getTime() === yesterday.getTime()
})

const formatDate = (date) => {
  const d = new Date(date)
  const weekdays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${weekdays[d.getDay()]}`
}

const formatDateForAPI = (date) => {
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const formatDuration = (minutes) => {
  if (minutes < 60) return `${minutes}分钟`
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return mins > 0 ? `${hours}小时${mins}分钟` : `${hours}小时`
}

const formatDeadlineTime = (dueDate, dueTime) => {
  const date = new Date(dueDate)
  const now = new Date()
  const diff = date - now
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days < 0) return '已过期'
  if (days === 0) return '今天' + (dueTime ? ` ${dueTime}` : '')
  if (days === 1) return '明天' + (dueTime ? ` ${dueTime}` : '')
  if (days <= 7) return `${days}天后` + (dueTime ? ` ${dueTime}` : '')
  return `${date.getMonth() + 1}月${date.getDate()}日` + (dueTime ? ` ${dueTime}` : '')
}

const getPriorityText = (priority) => {
  const map = { 1: '低', 2: '中', 3: '高' }
  return map[priority] || '低'
}
</script>

<style scoped>
.day-view {
  padding: 0;
}

.day-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.day-title {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
}

.day-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-primary {
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.4);
}

.btn-secondary {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
}

.btn-secondary:hover {
  background: rgba(168, 85, 247, 0.2);
}

.day-content {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 30px;
}

.day-main {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.section {
  background: transparent;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 25px;
}

.section-header {
  margin-bottom: 15px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.progress-text {
  font-size: 14px;
  font-weight: 400;
  color: #666;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #a855f7 0%, #7c3aed 100%);
  transition: width 0.3s ease;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.task-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  transition: all 0.2s ease;
  text-align: left; /* 确保文本左对齐 */
}

.task-item:hover {
  background: #f0f0f0;
}

.task-item.completed {
  opacity: 0.6;
}

.task-item.completed .task-title {
  text-decoration: line-through;
}

.task-checkbox {
  width: 20px;
  height: 20px;
  margin-top: 2px;
  cursor: pointer;
}

.task-content {
  flex: 1;
  text-align: left; /* 确保文本左对齐 */
}

.task-title {
  font-size: 15px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.task-description {
  font-size: 13px;
  color: #666;
  margin-bottom: 6px;
}

.task-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #999;
}

.priority {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
}

.priority-1 {
  background: #e8f5e9;
  color: #4caf50;
}

.priority-2 {
  background: #fff3e0;
  color: #ff9800;
}

.priority-3 {
  background: #ffebee;
  color: #f44336;
}

.task-actions {
  display: flex;
  gap: 6px;
}

.icon-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  transition: all 0.2s ease;
}

.icon-btn:hover {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
}

.entertainments-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.entertainment-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  text-align: left; /* 确保文本左对齐 */
}

.ent-icon {
  font-size: 24px;
}

.ent-content {
  flex: 1;
  text-align: left; /* 确保文本左对齐 */
}

.ent-type {
  font-size: 15px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.ent-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #666;
  margin-bottom: 4px;
}

.ent-notes {
  font-size: 13px;
  color: #666;
  margin-top: 4px;
  font-style: italic;
}

.ent-actions {
  display: flex;
  gap: 6px;
}

.reflection-editor {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.reflection-editor .btn-primary {
  align-self: flex-end;
  width: fit-content;
  padding: 10px 20px;
  font-size: 16px;
}

.reflection-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  resize: vertical;
  min-height: 120px;
}

.reflection-textarea:focus {
  outline: none;
  border-color: #a855f7;
}

.reflection-readonly {
  width: 100%;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  min-height: 120px;
  background-color: #f8f9fa;
  color: #666;
  white-space: pre-wrap;
  word-wrap: break-word;
  line-height: 1.6;
  text-align: left; /* 确保文本左对齐，不居中 */
}

.day-sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-section {
  background: transparent;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
}

.mini-calendar {
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0.7) 0%,
    rgba(255, 255, 255, 0.6) 30%,
    rgba(255, 255, 255, 0.5) 60%,
    rgba(255, 255, 255, 0.4) 100%
  );
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.sidebar-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 15px;
}

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.calendar-nav {
  width: 30px;
  height: 30px;
  border: none;
  background: #f0f0f0;
  border-radius: 6px;
  cursor: pointer;
  font-size: 18px;
  color: #666;
}

.calendar-nav:hover {
  background: #e0e0e0;
}

.calendar-month {
  font-size: 14px;
  font-weight: 500;
  color: #1a1a1a;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.calendar-weekday {
  font-size: 12px;
  color: #999;
  padding: 4px;
  text-align: center; /* 日历星期保留居中 */
}

.calendar-day {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.calendar-day:hover {
  background: #f0f0f0;
}

.calendar-day.today {
  background: rgba(168, 85, 247, 0.2);
  color: #1a1a1a;
  font-weight: 600;
  border: 2px solid #a855f7;
}

.calendar-day.selected {
  background: rgba(168, 85, 247, 0.2);
  color: #a855f7;
  font-weight: 600;
}

.calendar-day.other {
  color: #ccc;
}

.week-overview {
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0.7) 0%,
    rgba(255, 255, 255, 0.6) 30%,
    rgba(255, 255, 255, 0.5) 60%,
    rgba(255, 255, 255, 0.4) 100%
  );
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.week-day {
  display: flex;
  justify-content: space-between;
  padding: 8px;
  background: transparent;
  border-radius: 6px;
  font-size: 13px;
}

.week-day-name {
  font-weight: 500;
  color: #1a1a1a;
}

.week-day-count {
  color: #666;
}

.upcoming-deadlines {
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0.7) 0%,
    rgba(255, 255, 255, 0.6) 30%,
    rgba(255, 255, 255, 0.5) 60%,
    rgba(255, 255, 255, 0.4) 100%
  );
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.deadline-item {
  padding: 10px;
  background: transparent;
  border-radius: 6px;
  border-left: 3px solid #ff9800;
}

.deadline-title {
  font-size: 13px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.deadline-time {
  font-size: 12px;
  color: #666;
}

.empty-state {
  padding: 40px 20px;
  color: #999;
  font-size: 14px;
}

.empty-state-small {
  padding: 20px;
  color: #999;
  font-size: 12px;
}

@media (max-width: 1024px) {
  .day-content {
    grid-template-columns: 1fr;
  }

  .day-sidebar {
    order: -1;
  }
}
</style>
