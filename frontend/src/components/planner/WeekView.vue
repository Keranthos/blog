<template>
  <div class="week-view">
    <div class="week-header">
      <h2 class="week-title">周视图</h2>
      <div class="week-nav">
        <button class="btn btn-secondary" @click="prevWeek">‹ 上一周</button>
        <span class="week-range">{{ weekRange }}</span>
        <button class="btn btn-secondary" @click="nextWeek">下一周 ›</button>
      </div>
    </div>

    <div class="week-grid">
      <div
        v-for="day in weekDays"
        :key="day.date"
        class="week-day-column"
        :class="{ today: day.isToday }"
      >
        <div class="week-day-header">
          <div class="day-name">{{ day.name }}</div>
          <div class="day-date">{{ day.dateStr }}</div>
          <div class="day-task-count">{{ day.taskCount }}个任务</div>
        </div>
        <div class="week-day-tasks">
          <div
            v-for="task in day.tasks"
            :key="task.ID"
            class="week-task-item"
            :class="{ completed: task.completed, [`priority-${task.priority}`]: true }"
          >
            <input
              type="checkbox"
              :checked="task.completed"
              @click.stop="toggleTask(task)"
              class="task-checkbox"
            />
            <div class="task-title">{{ task.title }}</div>
          </div>
          <div v-if="day.tasks.length === 0" class="empty-day">
            暂无任务
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getTasks, updateTask } from '@/api/planner'

const weekStart = ref(new Date())
const allTasks = ref([])

const weekDays = computed(() => {
  const days = []
  const today = new Date()

  for (let i = 0; i < 7; i++) {
    const date = new Date(weekStart.value)
    date.setDate(weekStart.value.getDate() + i)
    const dateStr = `${date.getMonth() + 1}/${date.getDate()}`
    const isToday = date.toDateString() === today.toDateString()
    const dayName = ['日', '一', '二', '三', '四', '五', '六'][date.getDay()]
    const dateKey = date.toISOString().split('T')[0]

    const dayTasks = allTasks.value.filter(t => {
      const taskDate = new Date(t.date).toISOString().split('T')[0]
      return taskDate === dateKey
    })

    days.push({
      date: dateKey,
      dateStr,
      name: dayName,
      isToday,
      taskCount: dayTasks.length,
      tasks: dayTasks
    })
  }

  return days
})

const weekRange = computed(() => {
  const start = weekStart.value
  const end = new Date(start)
  end.setDate(start.getDate() + 6)
  return `${start.getMonth() + 1}/${start.getDate()} - ${end.getMonth() + 1}/${end.getDate()}`
})

onMounted(() => {
  initWeek()
  loadWeekTasks()
})

const initWeek = () => {
  const today = new Date()
  const day = today.getDay()
  weekStart.value = new Date(today)
  weekStart.value.setDate(today.getDate() - day)
  weekStart.value.setHours(0, 0, 0, 0)
}

const loadWeekTasks = async () => {
  try {
    const startDate = weekStart.value.toISOString().split('T')[0]
    const endDate = new Date(weekStart.value)
    endDate.setDate(weekStart.value.getDate() + 6)
    const endDateStr = endDate.toISOString().split('T')[0]

    // 获取整周的任务
    const response = await getTasks()
    if (response?.data) {
      allTasks.value = response.data.filter(task => {
        const taskDate = new Date(task.date).toISOString().split('T')[0]
        return taskDate >= startDate && taskDate <= endDateStr
      })
    } else {
      allTasks.value = []
    }
  } catch (error) {
    // 404 或其他错误时返回空数组
    if (error.response?.status === 404) {
      allTasks.value = []
    } else {
      console.error('加载周任务失败:', error)
      allTasks.value = []
    }
  }
}

const prevWeek = () => {
  weekStart.value = new Date(weekStart.value)
  weekStart.value.setDate(weekStart.value.getDate() - 7)
  loadWeekTasks()
}

const nextWeek = () => {
  weekStart.value = new Date(weekStart.value)
  weekStart.value.setDate(weekStart.value.getDate() + 7)
  loadWeekTasks()
}

const toggleTask = async (task) => {
  try {
    task.completed = !task.completed
    await updateTask(task.ID, task)
    await loadWeekTasks()
  } catch (error) {
    console.error('更新任务失败:', error)
    task.completed = !task.completed
  }
}

</script>

<style scoped>
.week-view {
  padding: 0;
}

.week-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.week-title {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
}

.week-nav {
  display: flex;
  align-items: center;
  gap: 15px;
}

.week-range {
  font-size: 16px;
  font-weight: 500;
  color: #666;
  min-width: 200px;
  text-align: center;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-secondary {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
}

.btn-secondary:hover {
  background: rgba(168, 85, 247, 0.2);
}

.week-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 15px;
}

.week-day-column {
  background: transparent;
  border-radius: 12px;
  padding: 15px;
  min-height: 400px;
}

.week-day-column.today {
  border: 2px solid #a855f7;
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.2);
}

.week-day-header {
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e0e0e0;
}

.day-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.day-date {
  font-size: 13px;
  color: #666;
  margin-bottom: 4px;
}

.day-task-count {
  font-size: 12px;
  color: #999;
}

.week-day-tasks {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.week-task-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 6px;
  transition: all 0.2s ease;
  border-left: 3px solid transparent;
}

.week-task-item.completed {
  opacity: 0.6;
}

.week-task-item.completed .task-title {
  text-decoration: line-through;
}

.week-task-item.priority-1 {
  border-left-color: #4caf50;
}

.week-task-item.priority-2 {
  border-left-color: #ff9800;
}

.week-task-item.priority-3 {
  border-left-color: #f44336;
}

.task-checkbox {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.task-title {
  flex: 1;
  font-size: 13px;
  color: #1a1a1a;
  text-align: left;
}

.empty-day {
  text-align: center;
  padding: 20px;
  color: #999;
  font-size: 12px;
}

@media (max-width: 1200px) {
  .week-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 768px) {
  .week-header {
    flex-direction: column;
    gap: 15px;
  }

  .week-grid {
    grid-template-columns: 1fr;
  }
}
</style>
