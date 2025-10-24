<template>
  <div class="location-update-view">
    <NavBar />
    <div class="update-container">
      <div class="update-header">
        <h1>📍 更新我的位置</h1>
        <p>设置你当前所在的城市，访客可以看到你那里的天气</p>
      </div>

      <div class="update-form">
        <div class="form-group">
          <label for="city">城市名称</label>
          <input
            id="city"
            v-model="form.city"
            type="text"
            placeholder="例如：北京、上海、广州"
            required
          />
        </div>

        <div class="form-group">
          <label for="country">国家（可选）</label>
          <input
            id="country"
            v-model="form.country"
            type="text"
            placeholder="例如：中国"
          />
        </div>

        <div class="coordinates">
          <div class="form-group">
            <label for="latitude">纬度</label>
            <input
              id="latitude"
              v-model.number="form.latitude"
              type="number"
              step="0.000001"
              placeholder="例如：39.9042"
              required
            />
          </div>

          <div class="form-group">
            <label for="longitude">经度</label>
            <input
              id="longitude"
              v-model.number="form.longitude"
              type="number"
              step="0.000001"
              placeholder="例如：116.4074"
              required
            />
          </div>
        </div>

        <div class="form-actions">
          <button
            :disabled="loading"
            class="btn-location"
            @click="getCurrentLocation"
          >
            📍 获取当前位置
          </button>

          <button
            :disabled="loading || !isFormValid"
            class="btn-update"
            @click="updateLocation"
          >
            {{ loading ? '更新中...' : '更新位置' }}
          </button>
        </div>
      </div>

      <div v-if="message" class="message" :class="messageType">
        {{ message }}
      </div>

      <!-- 当前天气预览 -->
      <div v-if="currentWeather" class="weather-preview">
        <h3>当前天气预览</h3>
        <div class="weather-card">
          <div class="weather-info">
            <div class="weather-main">
              <span class="temperature">{{ currentWeather.temperature }}</span>
              <span class="weather">{{ currentWeather.weather }}</span>
            </div>
            <div class="weather-location">{{ currentWeather.location }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import NavBar from '@/components/NavBar.vue'

const form = ref({
  city: '',
  country: '',
  latitude: null,
  longitude: null
})

const loading = ref(false)
const message = ref('')
const messageType = ref('')
const currentWeather = ref(null)

const isFormValid = computed(() => {
  return form.value.city &&
         form.value.latitude !== null &&
         form.value.longitude !== null
})

// 获取当前位置
const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    showMessage('浏览器不支持地理位置功能', 'error')
    return
  }

  loading.value = true
  navigator.geolocation.getCurrentPosition(
    (position) => {
      form.value.latitude = position.coords.latitude
      form.value.longitude = position.coords.longitude
      showMessage('位置获取成功！请填写城市名称', 'success')
      loading.value = false
    },
    (error) => {
      showMessage('获取位置失败：' + error.message, 'error')
      loading.value = false
    }
  )
}

// 更新位置
const updateLocation = async () => {
  if (!isFormValid.value) {
    showMessage('请填写完整信息', 'error')
    return
  }

  loading.value = true
  try {
    const response = await fetch('/api/weather/update-location', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('jwt')}`
      },
      body: JSON.stringify(form.value)
    })

    const data = await response.json()

    if (data.success) {
      showMessage('位置更新成功！天气信息已刷新', 'success')
      // 获取更新后的天气信息
      await getCurrentWeather()
      // 清空表单
      form.value = {
        city: '',
        country: '',
        latitude: null,
        longitude: null
      }
    } else {
      showMessage('更新失败：' + (data.error || '未知错误'), 'error')
    }
  } catch (error) {
    showMessage('网络错误：' + error.message, 'error')
  } finally {
    loading.value = false
  }
}

// 获取当前天气信息
const getCurrentWeather = async () => {
  try {
    const response = await fetch('/api/weather')
    const data = await response.json()

    if (data.success && data.data) {
      currentWeather.value = {
        location: data.data.location,
        temperature: data.data.temperature,
        weather: data.data.weather
      }
    }
  } catch (error) {
    console.warn('获取天气信息失败:', error)
  }
}

// 显示消息
const showMessage = (text, type) => {
  message.value = text
  messageType.value = type
  setTimeout(() => {
    message.value = ''
  }, 3000)
}

onMounted(async () => {
  await getCurrentWeather()
})
</script>

<style scoped>
.location-update-view {
  min-height: 100vh;
  padding-top: 100px;
}

.update-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 40px 20px;
}

.update-header {
  text-align: center;
  margin-bottom: 40px;
}

.update-header h1 {
  font-size: 2.5rem;
  color: #333;
  margin-bottom: 15px;
  font-weight: 700;
}

.update-header p {
  font-size: 1.1rem;
  color: #666;
  line-height: 1.6;
}

.update-form {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  margin-bottom: 30px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 25px;
}

.form-group label {
  color: #333;
  font-weight: 600;
  font-size: 1rem;
}

.form-group input {
  padding: 15px 20px;
  border: 2px solid rgba(102, 126, 234, 0.2);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.8);
  color: #333;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  background: white;
}

.form-group input::placeholder {
  color: #999;
}

.coordinates {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-actions {
  display: flex;
  gap: 20px;
  margin-top: 30px;
}

.btn-location,
.btn-update {
  flex: 1;
  padding: 15px 25px;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-location {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border: 2px solid rgba(102, 126, 234, 0.3);
}

.btn-location:hover:not(:disabled) {
  background: rgba(102, 126, 234, 0.2);
  transform: translateY(-2px);
}

.btn-update {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-update:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.btn-location:disabled,
.btn-update:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.message {
  padding: 15px 20px;
  border-radius: 12px;
  font-size: 1rem;
  text-align: center;
  margin-bottom: 20px;
}

.message.success {
  background: rgba(76, 175, 80, 0.1);
  color: #4CAF50;
  border: 1px solid rgba(76, 175, 80, 0.3);
}

.message.error {
  background: rgba(244, 67, 54, 0.1);
  color: #f44336;
  border: 1px solid rgba(244, 67, 54, 0.3);
}

.weather-preview {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.weather-preview h3 {
  color: #333;
  margin-bottom: 20px;
  font-size: 1.3rem;
  font-weight: 600;
}

.weather-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 15px;
  padding: 25px;
  color: white;
}

.weather-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.weather-main {
  display: flex;
  align-items: center;
  gap: 15px;
}

.temperature {
  font-size: 2.5rem;
  font-weight: 700;
}

.weather {
  font-size: 1.2rem;
  opacity: 0.9;
}

.weather-location {
  font-size: 1.1rem;
  opacity: 0.8;
}

@media (max-width: 768px) {
  .coordinates {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }

  .update-container {
    padding: 20px 15px;
  }

  .update-form {
    padding: 25px;
  }

  .weather-info {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }
}
</style>
