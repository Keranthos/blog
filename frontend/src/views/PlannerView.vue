<template>
  <div class="planner-view">
    <NavBar />
    <div class="planner-container">
      <!-- 标签页导航 -->
      <div ref="tabsContainer" class="planner-tabs">
        <div
          ref="dayTab"
          class="tab-item"
          :class="{ active: activeTab === 'day' }"
          @click="switchTab('day')"
        >
          <font-awesome-icon icon="calendar" class="tab-icon" />
          <span class="tab-text">日视图</span>
        </div>
        <div
          ref="weekTab"
          class="tab-item"
          :class="{ active: activeTab === 'week' }"
          @click="switchTab('week')"
        >
          <font-awesome-icon icon="calendar" class="tab-icon" />
          <span class="tab-text">周视图</span>
        </div>
        <div
          ref="deadlineTab"
          class="tab-item"
          :class="{ active: activeTab === 'deadline' }"
          @click="switchTab('deadline')"
        >
          <font-awesome-icon icon="clock" class="tab-icon" />
          <span class="tab-text">DDL视图</span>
        </div>
        <div
          class="tab-indicator"
          :style="indicatorStyle"
        ></div>
      </div>

      <!-- 内容区域 -->
      <div class="planner-content">
        <DayView v-if="activeTab === 'day'" :current-date="currentDate" />
        <WeekView v-if="activeTab === 'week'" :current-date="currentDate" />
        <DeadlineView v-if="activeTab === 'deadline'" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import NavBar from '@/components/NavBar.vue'
import DayView from '@/components/planner/DayView.vue'
import WeekView from '@/components/planner/WeekView.vue'
import DeadlineView from '@/components/planner/DeadlineView.vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

const activeTab = ref('day')
const currentDate = ref(new Date())
const tabsContainer = ref(null)
const dayTab = ref(null)
const weekTab = ref(null)
const deadlineTab = ref(null)
const windowWidth = ref(window.innerWidth)

const indicatorStyle = computed(() => {
  // 依赖 windowWidth 以确保窗口大小变化时重新计算
  // eslint-disable-next-line no-unused-vars
  const _width = windowWidth.value

  let activeTabElement = null
  if (activeTab.value === 'day') {
    activeTabElement = dayTab.value
  } else if (activeTab.value === 'week') {
    activeTabElement = weekTab.value
  } else if (activeTab.value === 'deadline') {
    activeTabElement = deadlineTab.value
  }

  if (!activeTabElement || !tabsContainer.value) {
    return {
      width: '0px',
      left: '0px',
      opacity: 0
    }
  }

  const containerRect = tabsContainer.value.getBoundingClientRect()
  const tabRect = activeTabElement.getBoundingClientRect()

  return {
    width: `${tabRect.width}px`,
    left: `${tabRect.left - containerRect.left}px`,
    opacity: 1
  }
})

const switchTab = (tab) => {
  activeTab.value = tab
}

// 监听窗口大小变化，重新计算指示器位置
const handleResize = () => {
  windowWidth.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.planner-view {
  min-height: 100vh;
  background-image: url('@/assets/background-image.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  padding-top: 80px; /* 为导航栏留出空间 */
  position: relative;
}

.planner-view::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0.7) 0%,
    rgba(255, 255, 255, 0.6) 30%,
    rgba(255, 255, 255, 0.5) 60%,
    rgba(255, 255, 255, 0.4) 100%
  );
  backdrop-filter: blur(6px);
  -webkit-backdrop-filter: blur(6px);
  z-index: 0;
  pointer-events: none;
}

.planner-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  position: relative;
  z-index: 1;
}

.planner-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  background: transparent; /* 移除背景 */
  padding: 10px;
  position: relative;
}

.tab-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 10px;
  cursor: pointer;
  transition: color 0.3s ease;
  color: #1a1a1a;
  font-weight: 500;
  position: relative;
  z-index: 1;
}

.tab-item:hover {
  color: #1a1a1a;
}

.tab-item.active {
  color: white; /* 白色字体 */
}

.tab-indicator {
  position: absolute;
  top: 10px;
  height: calc(100% - 20px);
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%); /* 亮紫色背景 */
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(168, 85, 247, 0.3);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 0;
  pointer-events: none;
}

.tab-icon {
  font-size: 18px;
}

.tab-text {
  font-size: 16px;
}

.planner-content {
  background: transparent;
  padding: 30px;
  min-height: 600px;
}

@media (max-width: 768px) {
  .planner-container {
    padding: 10px;
  }

  .planner-tabs {
    flex-direction: column;
    gap: 5px;
  }

  .tab-item {
    padding: 10px;
  }

  .planner-content {
    padding: 15px;
  }
}
</style>
