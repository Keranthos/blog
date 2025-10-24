<template>
  <div class="basic-view" @click="goToDetail">
    <img v-if="props.image" :src="props.image" alt="Article Image" class="image" loading="lazy" />
    <div class="content" :class="{ 'no-image': !props.image }">
      <div v-if="showTooltipFlag" class="tooltip">{{ props.title }}</div>
      <div v-if="showTooltipFlag" class="tooltip-arrow"></div>
      <h2 @mouseenter="showTooltip" @mouseleave="hideTooltip">
        <font-awesome-icon :icon="getTypeIcon()" class="type-icon" />
        {{ props.title }}
      </h2>
      <div v-if="props.tags && props.tags.length > 0" class="tags">
        <span v-for="tag in props.tags" :key="tag" class="tag">{{ tag }}</span>
      </div>
      <div class="meta-info">
        <p class="time">Modify {{ formattedTime }}</p>
        <p v-if="readingTime" class="reading-time">
          <span class="reading-icon">⏱️</span>
          {{ readingTime.display }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
// 用于展示一篇博客（科研日记、项目笔记）的基础介绍卡片
import { useRouter } from 'vue-router'

const props = defineProps({
  id: Number,
  image: String,
  title: String,
  tags: Array,
  time: String,
  type: String, // 新增类型属性，用于区分不同类型的内容
  readingTime: Object, // 阅读时间对象
  articleType: String // 文章类型（blog, project, research）
})

const router = useRouter()

// Tooltip控制
const showTooltipFlag = ref(false)

const showTooltip = (event) => {
  // console.log('showTooltip called')
  event.stopPropagation() // 阻止事件冒泡
  showTooltipFlag.value = true
}

const hideTooltip = (event) => {
  // console.log('hideTooltip called')
  event.stopPropagation() // 阻止事件冒泡
  showTooltipFlag.value = false
}

const goToDetail = () => {
  let routeName
  const articleType = props.articleType || props.type
  switch (articleType) {
    case 'blog':
      routeName = 'BlogDetail'
      break
    case 'project':
      routeName = 'ProjectDetail'
      break
    case 'research':
      routeName = 'ResearchDetail'
      break
    default:
      routeName = 'BlogDetail'
  }
  router.push({ name: routeName, params: { id: props.id } })
}

// 获取类型图标
const getTypeIcon = () => {
  const articleType = props.articleType || props.type
  const iconMap = {
    blog: 'blog',
    project: 'diagram-project',
    research: 'flask'
  }
  return iconMap[articleType] || 'blog'
}

// 格式化 time 为 'YYYY-MM-DD' 格式
const formattedTime = computed(() => {
  const date = new Date(props.time)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0') // 月份从0开始，需加1
  const day = date.getDate().toString().padStart(2, '0')
  return `${year}-${month}-${day}`
})

/* const image = computed(() => {
  return props.image || 'https://via.placeholder.com/300'
}) */
</script>

<style scoped>
.basic-view {
  position: relative;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
  transition: all 0s;
  border: 1px solid rgba(255, 255, 255, 0.3);
  height: 350px;
}

.basic-view::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(106, 27, 154, 0.1) 0%, rgba(142, 45, 226, 0.1) 100%);
  opacity: 0;
  transition: opacity 0s;
  z-index: 1;
}

.basic-view:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.basic-view:hover::before {
  opacity: 1;
}

.image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  transition: transform 0s;
  position: relative;
  z-index: 1;
  transform-origin: center top;
  display: block;
  margin: 0;
  padding: 0;
}

.basic-view:hover .image {
  transform: scale(1.05) translateY(-10px);
}

.content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px;
  flex: 1;
  position: relative;
  z-index: 2;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

.content.no-image {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 8px;
}

h2 {
  margin: 0 0 8px 0;
  font-size: 1.4em;
  font-weight: 700;
  color: #2c3e50;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0s;
}

.basic-view:hover h2 {
  color: #667eea;
}

/* Tooltip样式 - 淡黄色版本 */
.tooltip {
  position: absolute;
  bottom: 100%;
  left: 15px;
  width: max-content;
  background: rgba(255, 248, 220, 0.95);
  color: #8B4513;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 0.9rem;
  font-weight: 600;
  z-index: 9999999;
  white-space: nowrap;
  text-align: left;
  box-shadow: 0 8px 32px rgba(139, 69, 19, 0.15);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 248, 220, 0.3);
  pointer-events: none;
  margin-bottom: 8px;
}

/* Tooltip箭头 */
.tooltip-arrow {
  position: absolute;
  top: 100%;
  left: 35px;
  width: 0;
  height: 0;
  border-left: 8px solid transparent;
  border-right: 8px solid transparent;
  border-top: 8px solid rgba(255, 248, 220, 0.95);
  z-index: 9999999;
  pointer-events: none;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.tag {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  border: 1px solid rgba(102, 126, 234, 0.2);
  transition: all 0s;
}

.tag:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: translateY(-1px);
}

.meta-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.08);
}

.time {
  color: #666;
  font-size: 0.85em;
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
  font-weight: 500;
}

.time::before {
  content: '📅';
  font-size: 0.9em;
}

.reading-time {
  color: #667eea;
  font-size: 0.8em;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 4px;
  margin: 0;
  background: rgba(102, 126, 234, 0.12);
  padding: 4px 8px;
  border-radius: 12px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.reading-icon {
  font-size: 0.8em;
}

/* 添加标签样式 */
.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.tag {
  padding: 4px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 20px;
  font-size: 0.75em;
  font-weight: 500;
  transition: transform 0.2s ease;
}

.tag:hover {
  transform: scale(1.05);
}

/* 类型图标样式 */
.type-icon {
  font-size: 1.1em;
  color: #667eea;
  flex-shrink: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .basic-view {
    height: 350px;
  }

  .image {
    height: 200px;
  }

  h2 {
    font-size: 1.1em;
  }
}
</style>
