<template>
  <div
    v-if="visible"
    ref="menuRef"
    class="text-selection-menu"
    :style="menuStyle"
    @mousedown.prevent.stop
  >
    <button class="menu-btn" @click="handleCopy" @mousedown.prevent>
      <font-awesome-icon icon="copy" class="menu-icon" />
      <span class="menu-text">复制</span>
    </button>
    <button class="menu-btn" :class="{ active: isHighlighted }" @click="handleHighlight" @mousedown.prevent>
      <font-awesome-icon icon="highlighter" class="menu-icon" />
      <span class="menu-text">高亮</span>
    </button>
    <button class="menu-btn" @click="handleShare" @mousedown.prevent>
      <font-awesome-icon icon="share" class="menu-icon" />
      <span class="menu-text">分享</span>
    </button>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  position: {
    type: Object,
    default: () => ({ x: 0, y: 0 })
  },
  selectedText: {
    type: String,
    default: ''
  },
  isHighlighted: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['copy', 'highlight', 'share'])

const menuRef = ref(null)
const menuStyle = computed(() => {
  if (!props.visible) return {}

  // 确保菜单不会超出视口
  const menuWidth = 180 // 菜单宽度（增加了文字，需要更宽）
  const menuHeight = 50 // 菜单高度（增加了文字，需要更高）
  const padding = 10 // 边距

  let x = props.position.x
  let y = props.position.y - menuHeight - padding

  // 防止超出右边界
  if (x + menuWidth > window.innerWidth) {
    x = window.innerWidth - menuWidth - padding
  }

  // 防止超出左边界
  if (x < padding) {
    x = padding
  }

  // 如果上方空间不足，显示在下方
  if (y < padding) {
    y = props.position.y + padding
  }

  // 防止超出下边界
  if (y + menuHeight > window.innerHeight - padding) {
    y = window.innerHeight - menuHeight - padding
  }

  return {
    left: `${x}px`,
    top: `${y}px`,
    opacity: props.visible ? 1 : 0,
    transform: props.visible ? 'scale(1)' : 'scale(0.8)'
  }
})

const handleCopy = (e) => {
  e.preventDefault()
  e.stopPropagation()
  // 清除文本选择，避免显示选中特效
  window.getSelection().removeAllRanges()
  emit('copy', props.selectedText)
}

const handleHighlight = (e) => {
  e.preventDefault()
  e.stopPropagation()
  // 清除文本选择，避免显示选中特效
  window.getSelection().removeAllRanges()
  emit('highlight', props.selectedText)
}

const handleShare = (e) => {
  e.preventDefault()
  e.stopPropagation()
  // 清除文本选择，避免显示选中特效
  window.getSelection().removeAllRanges()
  emit('share', props.selectedText)
}

// 监听位置变化，确保菜单始终可见
watch(() => props.position, () => {
  if (props.visible) {
    nextTick(() => {
      // 菜单位置已通过 computed 自动计算
    })
  }
}, { deep: true })
</script>

<style scoped>
.text-selection-menu {
  position: fixed;
  z-index: 9999;
  display: flex;
  gap: 4px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 8px;
  padding: 4px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  transition: opacity 0.2s ease, transform 0.2s ease;
  pointer-events: auto;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.menu-btn {
  min-width: 50px;
  height: 50px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #666;
  transition: all 0.2s ease;
  padding: 6px 8px;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.menu-icon {
  font-size: 18px;
  transition: all 0.2s ease;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  pointer-events: none;
}

.menu-text {
  font-size: 12px;
  font-weight: 500;
  transition: all 0.2s ease;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  pointer-events: none;
}

.menu-btn:hover {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
  transform: scale(1.05);
}

.menu-btn:hover .menu-icon {
  color: #a855f7;
}

.menu-btn:hover .menu-text {
  color: #a855f7;
}

.menu-btn.active {
  /* 移除背景色，避免看起来像选中特效 */
  background: transparent;
  color: #666;
}

.menu-btn.active .menu-icon,
.menu-btn.active .menu-text {
  color: #666;
}

.menu-btn:active {
  transform: scale(0.95);
}
</style>
