<template>
  <div class="deadline-view">
    <div class="deadline-header">
      <h2 class="deadline-title">截止日期管理</h2>
      <button class="btn btn-primary" @click="showAddModal = true">
        <font-awesome-icon icon="plus" /> 添加DDL
      </button>
    </div>

    <div class="deadline-content">
      <div
        v-for="group in groupedDeadlines"
        :key="group.category"
        class="deadline-group"
      >
        <div class="group-header" :class="`group-${group.category}`">
          <font-awesome-icon :icon="group.icon" class="group-icon" />
          <span class="group-title">{{ group.title }}</span>
          <span class="group-count">({{ group.deadlines.length }})</span>
        </div>
        <div class="deadline-list">
          <div
            v-for="deadline in group.deadlines"
            :key="deadline.ID"
            class="deadline-item"
            :class="`category-${group.category}`"
          >
            <div class="deadline-checkbox">
              <input
                type="checkbox"
                :checked="deadline.completed"
                @change="toggleDeadline(deadline)"
              />
            </div>
            <div class="deadline-content">
              <div class="deadline-title-text">{{ deadline.title }}</div>
              <div v-if="deadline.description" class="deadline-description">
                {{ deadline.description }}
              </div>
              <div class="deadline-meta">
                <span class="deadline-time">{{ formatDeadlineTime(deadline.dueDate, deadline.dueTime) }}</span>
                <span class="deadline-countdown">{{ getCountdown(deadline.dueDate, deadline.dueTime) }}</span>
              </div>
            </div>
            <div class="deadline-actions">
              <button class="icon-btn" @click="editDeadline(deadline)">
                <font-awesome-icon icon="pen-to-square" />
              </button>
              <button class="icon-btn" @click="deleteDeadline(deadline.ID)">
                <font-awesome-icon icon="trash" />
              </button>
            </div>
          </div>
          <div v-if="group.deadlines.length === 0" class="empty-state">
            暂无{{ group.title }}
          </div>
        </div>
      </div>
    </div>

    <DeadlineModal
      v-if="showAddModal"
      :deadline="editingDeadline"
      @close="closeModal"
      @save="handleSave"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { getDeadlines, updateDeadline, deleteDeadline as deleteDeadlineAPI } from '@/api/planner'
import DeadlineModal from './DeadlineModal.vue'

const deadlines = ref([])
const showAddModal = ref(false)
const editingDeadline = ref(null)

const groupedDeadlines = computed(() => {
  const now = new Date()
  now.setHours(0, 0, 0, 0)

  // 计算本周的开始和结束
  const today = new Date(now)
  const startOfWeek = new Date(today)
  const dayOffset = (today.getDay() + 6) % 7
  startOfWeek.setDate(today.getDate() - dayOffset)
  startOfWeek.setHours(0, 0, 0, 0)
  const endOfWeek = new Date(startOfWeek)
  endOfWeek.setDate(startOfWeek.getDate() + 6)
  endOfWeek.setHours(23, 59, 59, 999)

  const groups = [
    {
      category: 'today',
      title: '今日',
      icon: 'calendar',
      deadlines: []
    },
    {
      category: 'week',
      title: '本周',
      icon: 'calendar',
      deadlines: []
    },
    {
      category: 'future',
      title: '未来',
      icon: 'calendar',
      deadlines: []
    }
  ]

  deadlines.value.forEach(deadline => {
    // 计算截止时间
    const dueDate = new Date(deadline.dueDate)
    dueDate.setHours(0, 0, 0, 0)

    // 如果有具体时间，设置时间
    const dueDateTime = new Date(dueDate)
    if (deadline.dueTime) {
      const [hours, minutes] = deadline.dueTime.split(':')
      dueDateTime.setHours(parseInt(hours) || 23, parseInt(minutes) || 59, 59, 999)
    } else {
      dueDateTime.setHours(23, 59, 59, 999)
    }

    // 如果已完成，检查是否已经过了截止时间
    if (deadline.completed) {
      const now = new Date()
      // 如果还没过截止时间，仍然显示；如果已经过了，则隐藏
      if (now > dueDateTime) {
        return // 已完成的DDL且已过截止时间，不显示
      }
    }

    // 判断属于哪个分类
    if (dueDate.getTime() === today.getTime()) {
      groups[0].deadlines.push(deadline) // 今日
    } else if (dueDate >= startOfWeek && dueDate <= endOfWeek) {
      groups[1].deadlines.push(deadline) // 本周
    } else if (dueDate > endOfWeek) {
      groups[2].deadlines.push(deadline) // 未来
    }
  })

  // 按优先级和日期排序
  groups.forEach(group => {
    group.deadlines.sort((a, b) => {
      if (a.priority !== b.priority) return b.priority - a.priority
      return new Date(a.dueDate) - new Date(b.dueDate)
    })
  })

  return groups
})

onMounted(() => {
  loadDeadlines()
})

const loadDeadlines = async () => {
  try {
    const response = await getDeadlines()
    if (response?.data) {
      deadlines.value = response.data
    } else {
      deadlines.value = []
    }
  } catch (error) {
    // 404 或其他错误时返回空数组
    if (error.response?.status === 404) {
      deadlines.value = []
    } else {
      console.error('加载DDL失败:', error)
      deadlines.value = []
    }
  }
}

const formatDateForAPI = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const toggleDeadline = async (deadline) => {
  try {
    deadline.completed = !deadline.completed
    // 格式化日期为字符串
    const deadlineData = {
      title: deadline.title,
      description: deadline.description || '',
      dueDate: formatDateForAPI(deadline.dueDate),
      dueTime: deadline.dueTime || '',
      priority: deadline.priority,
      completed: deadline.completed
    }
    await updateDeadline(deadline.ID, deadlineData)
    await loadDeadlines()
  } catch (error) {
    console.error('更新DDL失败:', error)
    deadline.completed = !deadline.completed
  }
}

const editDeadline = (deadline) => {
  editingDeadline.value = { ...deadline }
  showAddModal.value = true
}

const deleteDeadline = async (id) => {
  if (!confirm('确定要删除这个DDL吗？')) return
  try {
    await deleteDeadlineAPI(id)
    await loadDeadlines()
  } catch (error) {
    console.error('删除DDL失败:', error)
  }
}

const closeModal = () => {
  showAddModal.value = false
  editingDeadline.value = null
}

const handleSave = async () => {
  await loadDeadlines()
  closeModal()
}

const formatDeadlineTime = (dueDate, dueTime) => {
  const date = new Date(dueDate)
  const dateStr = `${date.getMonth() + 1}月${date.getDate()}日`
  return dueTime ? `${dateStr} ${dueTime}` : dateStr
}

const getCountdown = (dueDate, dueTime) => {
  const now = new Date()
  const due = new Date(dueDate)
  if (dueTime) {
    const [hours, minutes] = dueTime.split(':')
    due.setHours(parseInt(hours), parseInt(minutes), 0, 0)
  }

  const diff = due - now
  if (diff < 0) return '已过期'

  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))

  if (days === 0) {
    if (hours === 0) return '即将到期'
    return `剩余${hours}小时`
  }
  if (days === 1) return '明天'
  if (days <= 7) return `剩余${days}天`
  return `剩余${days}天`
}
</script>

<style scoped>
.deadline-view {
  padding: 0;
}

.deadline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.deadline-title {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
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
  color: white; /* 按钮文字保持白色，因为背景是紫色 */
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.4);
}

.deadline-content {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  align-items: start;
}

.deadline-item .deadline-content {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.deadline-group {
  background: transparent;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.group-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px 20px;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
}

.group-today {
  background: rgba(244, 67, 54, 0.1);
  border-left: 4px solid #f44336;
}

.group-week {
  background: rgba(255, 152, 0, 0.1);
  border-left: 4px solid #ff9800;
}

.group-future {
  background: rgba(76, 175, 80, 0.1);
  border-left: 4px solid #4caf50;
}

.group-icon {
  font-size: 18px;
}

.group-count {
  font-size: 14px;
  opacity: 0.9;
}

.deadline-list {
  padding: 15px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
  overflow-y: auto;
  background: rgba(255, 255, 255, 0.05);
}

.deadline-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 15px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(5px);
  border-radius: 8px;
  border-left: 4px solid transparent;
  transition: all 0.2s ease;
}

.deadline-item:hover {
  background: rgba(255, 255, 255, 0.3);
}

.deadline-item.category-today {
  border-left-color: #f44336;
}

.deadline-item.category-week {
  border-left-color: #ff9800;
}

.deadline-item.category-future {
  border-left-color: #4caf50;
}

.deadline-checkbox {
  margin-top: 2px;
}

.deadline-checkbox input {
  width: 20px;
  height: 20px;
  cursor: pointer;
}

.deadline-title-text {
  font-size: 16px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 6px;
}

.deadline-description {
  font-size: 13px;
  color: #333;
  margin-bottom: 8px;
}

.deadline-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #666;
}

.deadline-time {
  font-weight: 500;
  color: #333;
}

.deadline-countdown {
  color: #ffd700;
  font-weight: 500;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

.deadline-actions {
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

.empty-state {
  text-align: center;
  padding: 30px;
  color: #666;
  font-size: 14px;
}

@media (max-width: 1024px) {
  .deadline-content {
    grid-template-columns: 1fr;
  }
}
</style>
