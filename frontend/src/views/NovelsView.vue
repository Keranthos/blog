<template>
  <div class="media-view">
    <NavBar />
    <div class="media-layout">
      <!-- 左侧：顶部图片 + 外部文案 -->
      <div class="left-col">
        <div class="header">
          <img :src="novelsBackground" alt="Background Image" class="header-image" loading="lazy" decoding="async" @error="onImgError($event)" />
        </div>
        <div class="intro">
          <div v-if="introTitle" class="intro-title">{{ introTitle }}</div>
          <pre v-if="introBody" class="intro-text">{{ introBody }}</pre>
        </div>
      </div>

      <!-- 右侧：长条列表卡片（无图片） -->
      <div class="right-col" @scroll="handleScroll">
        <div class="media-rows">
          <MediaCard
            v-for="media in mediaList"
            :key="media.ID + '_' + (media.__type || '')"
            :ref="el => { if (el) mediaCardRefs[media.ID] = el }"
            :media="media"
            :type="media.__type || 'novels'"
            variant="row"
            @media-deleted="removeMedia"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, defineAsyncComponent } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMediaList } from '@/api/media/browse'
import MediaCard from '@/components/MediaCard.vue'
import novelsBackgroundImg from '@/assets/novels_background.jpg'

const NavBar = defineAsyncComponent(() => import('@/components/NavBar.vue'))
const route = useRoute()
const router = useRouter()

const novelsBackground = novelsBackgroundImg
const mediaList = ref([])
const currentPage = ref(1)
const isLoading = ref(false)
const mediaCache = ref({})
const introText = ref('')
const introTitle = ref('')
const introBody = ref('')
const fallbackImg = '/images/sunset-mountains.jpg'
const mediaCardRefs = ref({}) // 存储MediaCard组件的引用

// 解析介绍文本，提取第一句话作为标题
const parseIntroText = (text) => {
  if (!text) {
    introTitle.value = ''
    introBody.value = ''
    return
  }

  // 找到第一个句号、问号或感叹号
  const firstSentenceEnd = text.match(/[。！？]/)
  if (firstSentenceEnd) {
    const endIndex = firstSentenceEnd.index + 1
    introTitle.value = text.substring(0, endIndex).trim()
    // 正文从标题结束位置开始，保留所有后续内容（包括换行和空格）
    introBody.value = text.substring(endIndex)
  } else {
    // 如果没有找到标点，则将第一行作为标题（遇到第一个换行符为止）
    const firstLineEnd = text.indexOf('\n')
    if (firstLineEnd > 0) {
      introTitle.value = text.substring(0, firstLineEnd).trim()
      // 正文从第一行结束后开始，保留所有后续内容（包括换行和空格）
      introBody.value = text.substring(firstLineEnd + 1) // +1 跳过换行符，但保留后续格式
    } else {
      introTitle.value = text.trim()
      introBody.value = ''
    }
  }
}

// 图片错误回退
const onImgError = (e) => {
  const img = e?.target
  if (img && img.src !== fallbackImg) {
    img.src = fallbackImg
  }
}

const loadMedia = async () => {
  if (isLoading.value) return
  isLoading.value = true
  try {
    if (mediaCache.value[currentPage.value]) {
      mediaList.value.push(...mediaCache.value[currentPage.value])
    } else {
      const [novelsRes, booksRes, moviesRes] = await Promise.all([
        getMediaList('novels', currentPage.value),
        getMediaList('books', currentPage.value),
        getMediaList('movies', currentPage.value)
      ])
      const pick = (res) => {
        if (!res) return []
        if (Array.isArray(res.data)) return res.data
        if (res.data && Array.isArray(res.data.medias)) return res.data.medias
        if (Array.isArray(res.medias)) return res.medias
        return []
      }
      const merged = []
      merged.push(...pick(novelsRes).map(m => ({ ...m, __type: 'novels' })))
      merged.push(...pick(booksRes).map(m => ({ ...m, __type: 'books' })))
      merged.push(...pick(moviesRes).map(m => ({ ...m, __type: 'movies' })))
      merged.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt))
      mediaList.value.push(...merged)
      mediaCache.value[currentPage.value] = merged
    }
    currentPage.value++
  } catch (error) {
    console.error('Failed to load media:', error)
  } finally {
    isLoading.value = false
  }
}

const handleScroll = (event) => {
  const { scrollTop, scrollHeight, clientHeight } = event.target
  if (scrollTop + clientHeight >= scrollHeight - 10) {
    loadMedia()
  }
}

const removeMedia = (mediaId) => {
  mediaList.value = mediaList.value.filter(media => media.ID !== mediaId)
  for (const page in mediaCache.value) {
    mediaCache.value[page] = mediaCache.value[page].filter(media => media.ID !== mediaId)
  }
}

// 打开指定媒体的详情框
const openMediaById = async (mediaId) => {
  // 确保媒体列表已加载
  if (mediaList.value.length === 0) {
    await loadMedia()
  }

  // 查找对应的媒体卡片
  const targetMedia = mediaList.value.find(m => m.ID === parseInt(mediaId))
  if (targetMedia && mediaCardRefs.value[targetMedia.ID]) {
    // 调用MediaCard的打开方法（如果MediaCard暴露了openMediaModal方法）
    const cardComponent = mediaCardRefs.value[targetMedia.ID]
    if (cardComponent && typeof cardComponent.openMediaModal === 'function') {
      cardComponent.openMediaModal()
    }
  }
}

// 监听路由query变化，打开指定的媒体详情框
watch(() => route.query.mediaId, async (mediaId) => {
  if (mediaId) {
    // 等待DOM更新
    await nextTick()
    // 延迟一点确保MediaCard组件已完全渲染
    setTimeout(() => {
      openMediaById(mediaId)
      // 清除query参数，避免刷新时重复打开
      router.replace({ query: {} })
    }, 300)
  }
}, { immediate: true })

onMounted(() => {
  loadMedia()
  // 加载左侧文案（来自 public/data/media-intro.txt）
  fetch('/data/media-intro.txt')
    .then(res => res.ok ? res.text() : '')
    .then(text => {
      introText.value = text || ''
      parseIntroText(text || '')
    })
    .catch(() => {
      introText.value = ''
      parseIntroText('')
    })

  // 如果有mediaId参数，打开对应的详情框
  if (route.query.mediaId) {
    nextTick(() => {
      setTimeout(() => {
        openMediaById(route.query.mediaId)
        // 清除query参数
        router.replace({ query: {} })
      }, 500)
    })
  }
})
</script>

<style scoped>
.media-view {
  min-height: 100vh;
  padding-top: 56px;
  padding-bottom: 100px;
  overflow-x: hidden; /* 防止横向滚动 */
  box-sizing: border-box;
}

.media-layout { width: 80%; max-width: 1280px; margin: 20px auto 0; display: grid; grid-template-columns: 1fr 1fr; gap: 36px; }

.left-col { display: flex; flex-direction: column; gap: 20px; }
.right-col { overflow-y: auto; max-height: calc(100vh - 180px); padding-right: 4px; }
.right-col { scrollbar-width: none; -ms-overflow-style: none; }
.right-col::-webkit-scrollbar { display: none; }

.header { position: relative; width: 100%; margin-bottom: 0; overflow: hidden; }

.header-image {
  width: 100%;
  height: auto;
  max-height: 420px;
  object-fit: cover;
  border-radius: 8px;
  display: block;
  margin: 0 auto;
}

.intro { background: transparent; border: none; border-radius: 0; padding: 8px 0; }

.intro-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #a855f7;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid rgba(168, 85, 247, 0.2);
  line-height: 1.6;
  letter-spacing: 0.3px;
  text-align: left;
}

.intro-text { white-space: pre-wrap; text-align: left; margin: 0; color: #333; font-size: 0.95rem; line-height: 1.7; }

.media-rows { display: flex; flex-direction: column; gap: 12px; padding: 4px; }

.create-media {
  width: 90%;
  max-width: 800px;
  margin: 40px auto 0;
  padding: 30px;
  background: transparent;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
}

.create-media h3 {
  color: #333;
  margin-bottom: 20px;
  font-size: 1.5rem;
}

.create-media input,
.create-media textarea {
  width: 100%;
  padding: 15px;
  margin-bottom: 15px;
  border: 2px solid #e8eaf0;
  border-radius: 12px;
  font-size: 1rem;
  font-family: inherit;
  transition: all 0.3s ease;
}

.create-media input:focus,
.create-media textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.create-media textarea {
  min-height: 120px;
  resize: vertical;
}

.create-media button {
  padding: 12px 30px;
  border: none;
  border-radius: 25px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.create-media button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(102, 126, 234, 0.5);
}

@media (max-width: 1200px) { .media-layout { width: 86%; grid-template-columns: 1fr; gap: 24px; } .right-col { max-height: none; } }

/* 紧凑模式：头图与卡片区域等宽（2/3 居中），高度不限制 */
@media (max-width: 1330px) {
  .media-view { padding-top: 40px; }
  .media-layout { width: 66.666%; margin: 0 auto; gap: 20px; }
  .right-col { padding-right: 0; }
}

@media (max-width: 1024px) {
  .media-list {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .media-view {
    padding-top: 80px;
  }

  /* 手机端：移除最小宽度限制，使用 100% 宽度 */
  .media-layout {
    width: 100% !important;
    margin: 0 auto;
    padding: 0 15px;
    min-width: 0 !important;
  }

  /* 手机端隐藏图片下方文案，仅展示卡片 */
  .intro { display: none; }

  .create-media {
    width: 95%;
    padding: 20px;
  }
}
</style>
