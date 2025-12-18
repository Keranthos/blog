<template>
  <div v-if="visible" class="emoji-picker-wrapper" @click.self="close">
    <div class="emoji-picker" @click.stop>
      <div class="emoji-picker-header">
        <span class="emoji-picker-title">选择表情</span>
        <button class="emoji-picker-close" @click="close">
          <font-awesome-icon icon="times" />
        </button>
      </div>
      <div class="emoji-picker-content">
        <div
          v-for="category in emojiCategories"
          :key="category.name"
          class="emoji-category"
        >
          <div class="emoji-category-title">{{ category.title }}</div>
          <div class="emoji-grid">
            <button
              v-for="emoji in category.emojis"
              :key="emoji"
              class="emoji-item"
              :title="emoji"
              @click="selectEmoji(emoji)"
            >
              {{ emoji }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['select', 'close'])

// Emoji 分类数据
const emojiCategories = ref([
  {
    name: 'smileys',
    title: '表情',
    emojis: [
      '😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂', '🙂', '🙃',
      '😉', '😊', '😇', '🥰', '😍', '🤩', '😘', '😗', '😚', '😙',
      '😋', '😛', '😜', '🤪', '😝', '🤑', '🤗', '🤭', '🤫', '🤔',
      '🤐', '🤨', '😐', '😑', '😶', '😏', '😒', '🙄', '😬', '🤥',
      '😌', '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕', '🤢', '🤮',
      '🤧', '🥵', '🥶', '😶‍🌫️', '😵', '😵‍💫', '🤯', '🤠', '🥳', '🥸'
    ]
  },
  {
    name: 'gestures',
    title: '手势',
    emojis: [
      '👋', '🤚', '🖐', '✋', '🖖', '👌', '🤌', '🤏', '✌️', '🤞',
      '🤟', '🤘', '🤙', '👈', '👉', '👆', '🖕', '👇', '☝️', '👍',
      '👎', '✊', '👊', '🤛', '🤜', '👏', '🙌', '👐', '🤲', '🤝',
      '🙏', '✍️', '💪', '🦾', '🦿', '🦵', '🦶', '👂', '🦻', '👃'
    ]
  },
  {
    name: 'people',
    title: '人物',
    emojis: [
      '👶', '👧', '🧒', '👦', '👩', '🧑', '👨', '👩‍🦱', '👨‍🦱', '👩‍🦰',
      '👨‍🦰', '👱‍♀️', '👱', '👩‍🦳', '👨‍🦳', '👩‍🦲', '👨‍🦲', '🧔', '👵', '🧓',
      '👴', '👲', '👳‍♀️', '👳', '👮‍♀️', '👮', '👷‍♀️', '👷', '💂‍♀️', '💂',
      '🕵️‍♀️', '🕵️', '👩‍⚕️', '👨‍⚕️', '👩‍🌾', '👨‍🌾', '👩‍🍳', '👨‍🍳', '👩‍🎓', '👨‍🎓',
      '👩‍🎤', '👨‍🎤', '👩‍🏫', '👨‍🏫', '👩‍🏭', '👨‍🏭', '👩‍💻', '👨‍💻', '👩‍💼', '👨‍💼'
    ]
  },
  {
    name: 'animals',
    title: '动物',
    emojis: [
      '🐶', '🐱', '🐭', '🐹', '🐰', '🦊', '🐻', '🐼', '🐨', '🐯',
      '🦁', '🐮', '🐷', '🐽', '🐸', '🐵', '🙈', '🙉', '🙊', '🐒',
      '🐔', '🐧', '🐦', '🐤', '🐣', '🐥', '🦆', '🦅', '🦉', '🦇',
      '🐺', '🐗', '🐴', '🦄', '🐝', '🐛', '🦋', '🐌', '🐞', '🐜',
      '🦟', '🦗', '🕷', '🦂', '🐢', '🐍', '🦎', '🦖', '🦕', '🐙'
    ]
  },
  {
    name: 'food',
    title: '食物',
    emojis: [
      '🍏', '🍎', '🍐', '🍊', '🍋', '🍌', '🍉', '🍇', '🍓', '🍈',
      '🍒', '🍑', '🥭', '🍍', '🥥', '🥝', '🍅', '🍆', '🥑', '🥦',
      '🥬', '🥒', '🌶', '🌽', '🥕', '🥔', '🍠', '🥐', '🥯', '🍞',
      '🥖', '🥨', '🧀', '🥚', '🍳', '🥞', '🥓', '🥩', '🍗', '🍖',
      '🦴', '🌭', '🍔', '🍟', '🍕', '🥪', '🥙', '🌮', '🌯', '🥗'
    ]
  },
  {
    name: 'activities',
    title: '活动',
    emojis: [
      '⚽', '🏀', '🏈', '⚾', '🥎', '🎾', '🏐', '🏉', '🥏', '🎱',
      '🏓', '🏸', '🥅', '🏒', '🏑', '🏏', '🥍', '🏹', '🎣', '🥊',
      '🥋', '🎽', '🛹', '🛷', '⛸', '🥌', '🎿', '⛷', '🏂', '🏋️‍♀️',
      '🏋️', '🤼‍♀️', '🤼', '🤸‍♀️', '🤸', '🤺', '🤾‍♀️', '🤾', '🏌️‍♀️', '🏌️',
      '🏇', '🧘‍♀️', '🧘', '🏄‍♀️', '🏄', '🏊‍♀️', '🏊', '🤽‍♀️', '🤽', '🚣‍♀️'
    ]
  },
  {
    name: 'objects',
    title: '物品',
    emojis: [
      '⌚', '📱', '📲', '💻', '⌨️', '🖥', '🖨', '🖱', '🖲', '🕹',
      '🗜', '💾', '💿', '📀', '📼', '📷', '📸', '📹', '🎥', '📽',
      '🎞', '📞', '☎️', '📟', '📠', '📺', '📻', '🎙', '🎚', '🎛',
      '⏱', '⏲', '⏰', '🕰', '⌛', '⏳', '📡', '🔋', '🔌', '💡',
      '🔦', '🕯', '🧯', '🛢', '💸', '💵', '💴', '💶', '💷', '💰'
    ]
  },
  {
    name: 'symbols',
    title: '符号',
    emojis: [
      '❤️', '🧡', '💛', '💚', '💙', '💜', '🖤', '🤍', '🤎', '💔',
      '❣️', '💕', '💞', '💓', '💗', '💖', '💘', '💝', '💟', '☮️',
      '✝️', '☪️', '🕉', '☸️', '✡️', '🔯', '🕎', '☯️', '☦️', '🛐',
      '⛎', '♈', '♉', '♊', '♋', '♌', '♍', '♎', '♏', '♐',
      '♑', '♒', '♓', '🆔', '⚛️', '🉑', '☢️', '☣️', '📴', '📳'
    ]
  }
])

const selectEmoji = (emoji) => {
  emit('select', emoji)
}

const close = () => {
  emit('close')
}
</script>

<style scoped>
.emoji-picker-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.emoji-picker {
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: slideUp 0.2s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.emoji-picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
}

.emoji-picker-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.emoji-picker-close {
  background: none;
  border: none;
  cursor: pointer;
  color: #666;
  font-size: 18px;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.emoji-picker-close:hover {
  background: #e5e7eb;
  color: #333;
}

.emoji-picker-content {
  padding: 16px;
  overflow-y: auto;
  max-height: calc(80vh - 60px);
}

.emoji-category {
  margin-bottom: 24px;
}

.emoji-category:last-child {
  margin-bottom: 0;
}

.emoji-category-title {
  font-size: 14px;
  font-weight: 600;
  color: #666;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e5e7eb;
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
}

.emoji-item {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 28px;
  padding: 8px;
  border-radius: 8px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  font-family: "Apple Color Emoji", "Segoe UI Emoji", "Noto Color Emoji", "Segoe UI Symbol", "Android Emoji", "EmojiSymbols", "EmojiOne Mozilla", "Twemoji Mozilla", "Segoe UI", sans-serif;
}

.emoji-item:hover {
  background: #f3f4f6;
  transform: scale(1.2);
}

.emoji-item:active {
  transform: scale(1.1);
}

/* 滚动条样式 */
.emoji-picker-content::-webkit-scrollbar {
  width: 6px;
}

.emoji-picker-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.emoji-picker-content::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.emoji-picker-content::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .emoji-picker {
    width: 95%;
    max-height: 70vh;
  }

  .emoji-grid {
    grid-template-columns: repeat(6, 1fr);
    gap: 6px;
  }

  .emoji-item {
    font-size: 24px;
    padding: 6px;
  }
}
</style>
