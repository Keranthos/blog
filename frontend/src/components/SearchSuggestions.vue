<template>
  <div v-if="show" class="search-suggestions">
    <div class="suggestions-container">
      <!-- 搜索历史 -->
      <div v-if="searchHistory.length > 0" class="suggestions-section">
        <div class="section-header">
          <font-awesome-icon icon="clock" class="section-icon" />
          <span>搜索历史</span>
        </div>
        <div class="suggestions-list">
          <div
            v-for="(item, index) in filteredHistory"
            :key="`history-${index}`"
            class="suggestion-item history-item"
            @click="selectSuggestion(item)"
            @mouseenter="setActiveIndex(index)"
          >
            <font-awesome-icon icon="clock" class="item-icon" />
            <span class="item-text">{{ item }}</span>
            <button
              class="remove-btn"
              title="删除"
              @click.stop="removeFromHistory(item)"
            >
              <font-awesome-icon icon="times" />
            </button>
          </div>
        </div>
      </div>

      <!-- 热门搜索 -->
      <div v-if="hotKeywords.length > 0" class="suggestions-section">
        <div class="section-header">
          <font-awesome-icon icon="fire" class="section-icon" />
          <span>热门搜索</span>
        </div>
        <div class="suggestions-list">
          <div
            v-for="(keyword, index) in filteredHotKeywords"
            :key="`hot-${index}`"
            class="suggestion-item hot-item"
            @click="selectSuggestion(keyword)"
            @mouseenter="setActiveIndex(filteredHistory.length + index)"
          >
            <font-awesome-icon icon="fire" class="item-icon" />
            <span class="item-text">{{ keyword }}</span>
            <span class="hot-badge">{{ index + 1 }}</span>
          </div>
        </div>
      </div>

      <!-- 相关标签 -->
      <div v-if="relatedTags.length > 0" class="suggestions-section">
        <div class="section-header">
          <font-awesome-icon icon="tags" class="section-icon" />
          <span>相关标签</span>
        </div>
        <div class="suggestions-list">
          <div
            v-for="(tag, index) in filteredRelatedTags"
            :key="`tag-${index}`"
            class="suggestion-item tag-item"
            @click="selectSuggestion(tag)"
            @mouseenter="setActiveIndex(filteredHistory.length + filteredHotKeywords.length + index)"
          >
            <font-awesome-icon icon="tag" class="item-icon" />
            <span class="item-text">{{ tag }}</span>
          </div>
        </div>
      </div>

      <!-- 无建议时显示提示 -->
      <div v-if="!hasSuggestions" class="no-suggestions">
        <font-awesome-icon icon="lightbulb" class="no-suggestions-icon" />
        <p>输入关键词开始搜索</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  keyword: {
    type: String,
    default: ''
  },
  searchHistory: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['select', 'close', 'remove-from-history'])

// 热门搜索关键词
const hotKeywords = ref([
  'Vue.js',
  'JavaScript',
  '前端开发',
  'React',
  'Node.js',
  'Python',
  'Go语言',
  '数据库',
  '算法',
  '机器学习'
])

// 相关标签
const relatedTags = ref([
  '技术分享',
  '项目经验',
  '学习笔记',
  '生活感悟',
  '编程技巧',
  '前端框架',
  '后端开发',
  '数据库设计',
  '系统架构',
  '代码优化'
])

const activeIndex = ref(0)

// 过滤搜索历史
const filteredHistory = computed(() => {
  if (!props.keyword) return props.searchHistory.slice(0, 5)
  return props.searchHistory
    .filter(item => item.toLowerCase().includes(props.keyword.toLowerCase()))
    .slice(0, 5)
})

// 过滤热门关键词
const filteredHotKeywords = computed(() => {
  if (!props.keyword) return hotKeywords.value.slice(0, 5)
  return hotKeywords.value
    .filter(keyword => keyword.toLowerCase().includes(props.keyword.toLowerCase()))
    .slice(0, 5)
})

// 过滤相关标签
const filteredRelatedTags = computed(() => {
  if (!props.keyword) return relatedTags.value.slice(0, 5)
  return relatedTags.value
    .filter(tag => tag.toLowerCase().includes(props.keyword.toLowerCase()))
    .slice(0, 5)
})

// 是否有建议
const hasSuggestions = computed(() => {
  return filteredHistory.value.length > 0 ||
         filteredHotKeywords.value.length > 0 ||
         filteredRelatedTags.value.length > 0
})

// 总建议数量
const totalSuggestions = computed(() => {
  return filteredHistory.value.length +
         filteredHotKeywords.value.length +
         filteredRelatedTags.value.length
})

// 选择建议
const selectSuggestion = (suggestion) => {
  emit('select', suggestion)
}

// 从历史记录中删除
const removeFromHistory = (item) => {
  emit('remove-from-history', item)
}

// 设置激活索引
const setActiveIndex = (index) => {
  activeIndex.value = index
}

// 键盘导航
const handleKeydown = (event) => {
  if (!props.show || !hasSuggestions.value) return

  switch (event.key) {
    case 'ArrowDown':
      event.preventDefault()
      activeIndex.value = (activeIndex.value + 1) % totalSuggestions.value
      break
    case 'ArrowUp':
      event.preventDefault()
      activeIndex.value = activeIndex.value > 0
        ? activeIndex.value - 1
        : totalSuggestions.value - 1
      break
    case 'Enter': {
      event.preventDefault()
      const suggestion = getSuggestionByIndex(activeIndex.value)
      if (suggestion) {
        selectSuggestion(suggestion)
      }
      break
    }
    case 'Escape':
      emit('close')
      break
  }
}

// 根据索引获取建议
const getSuggestionByIndex = (index) => {
  let currentIndex = 0

  // 搜索历史
  if (index < filteredHistory.value.length) {
    return filteredHistory.value[index]
  }
  currentIndex += filteredHistory.value.length

  // 热门关键词
  if (index < currentIndex + filteredHotKeywords.value.length) {
    return filteredHotKeywords.value[index - currentIndex]
  }
  currentIndex += filteredHotKeywords.value.length

  // 相关标签
  if (index < currentIndex + filteredRelatedTags.value.length) {
    return filteredRelatedTags.value[index - currentIndex]
  }

  return null
}

// 监听关键词变化，重置激活索引
watch(() => props.keyword, () => {
  activeIndex.value = 0
})

// 监听显示状态
watch(() => props.show, (show) => {
  if (show) {
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})

onMounted(() => {
  if (props.show) {
    document.addEventListener('keydown', handleKeydown)
  }
})
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.search-suggestions {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(102, 126, 234, 0.1);
  z-index: 1000;
  max-height: 400px;
  overflow-y: auto;
  backdrop-filter: blur(10px);
}

.suggestions-container {
  padding: 8px;
}

.suggestions-section {
  margin-bottom: 8px;
}

.suggestions-section:last-child {
  margin-bottom: 0;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  font-size: 0.85rem;
  font-weight: 600;
  color: #666;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
  margin-bottom: 4px;
}

.section-icon {
  font-size: 0.8rem;
  color: #667eea;
}

.suggestions-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.suggestion-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.suggestion-item:hover,
.suggestion-item.active {
  background: rgba(102, 126, 234, 0.1);
  transform: translateX(4px);
}

.item-icon {
  font-size: 0.9rem;
  color: #999;
  width: 16px;
  flex-shrink: 0;
}

.item-text {
  flex: 1;
  font-size: 0.95rem;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.remove-btn {
  padding: 4px;
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s ease;
  opacity: 0;
  flex-shrink: 0;
}

.suggestion-item:hover .remove-btn {
  opacity: 1;
}

.remove-btn:hover {
  background: rgba(255, 0, 0, 0.1);
  color: #ff4444;
}

.hot-badge {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
  flex-shrink: 0;
}

/* 不同类型的样式 */
.history-item .item-icon {
  color: #999;
}

.hot-item .item-icon {
  color: #ff6b6b;
}

.tag-item .item-icon {
  color: #667eea;
}

.no-suggestions {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.no-suggestions-icon {
  font-size: 2rem;
  color: #ddd;
  margin-bottom: 12px;
}

.no-suggestions p {
  font-size: 0.95rem;
  margin: 0;
}

/* 暗色模式适配 */
[data-theme="dark"] .search-suggestions {
  background: rgba(26, 33, 62, 0.95);
  border-color: rgba(255, 255, 255, 0.1);
}

[data-theme="dark"] .section-header {
  background: rgba(102, 126, 234, 0.1);
  color: #b0b0b0;
}

[data-theme="dark"] .suggestion-item:hover,
[data-theme="dark"] .suggestion-item.active {
  background: rgba(102, 126, 234, 0.15);
}

[data-theme="dark"] .item-text {
  color: #e0e0e0;
}

[data-theme="dark"] .item-icon {
  color: #999;
}

[data-theme="dark"] .remove-btn {
  color: #999;
}

[data-theme="dark"] .remove-btn:hover {
  color: #ff6b6b;
}

[data-theme="dark"] .no-suggestions {
  color: #b0b0b0;
}

[data-theme="dark"] .no-suggestions-icon {
  color: #666;
}

/* 滚动条样式 */
.search-suggestions::-webkit-scrollbar {
  width: 6px;
}

.search-suggestions::-webkit-scrollbar-track {
  background: rgba(102, 126, 234, 0.1);
  border-radius: 3px;
}

.search-suggestions::-webkit-scrollbar-thumb {
  background: rgba(102, 126, 234, 0.3);
  border-radius: 3px;
}

.search-suggestions::-webkit-scrollbar-thumb:hover {
  background: rgba(102, 126, 234, 0.5);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .suggestion-item {
    padding: 12px;
  }

  .item-text {
    font-size: 1rem;
  }

  .remove-btn {
    opacity: 1;
  }
}
</style>
