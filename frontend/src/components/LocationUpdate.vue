<template>
  <div class="location-update">
    <div class="update-header">
      <h3>📍 更新我的位置</h3>
      <p>设置你当前所在的城市，访客可以看到你那里的天气</p>
    </div>

    <div class="update-form">
      <div class="form-group">
        <label for="city">城市名称</label>
        <input
          id="city"
          v-model="form.city"
          type="text"
          placeholder="例如：北京"
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
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const form = ref({
  city: '',
  country: '',
  latitude: null,
  longitude: null
})

const loading = ref(false)
const message = ref('')
const messageType = ref('')

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
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(form.value)
    })

    const data = await response.json()

    if (data.success) {
      showMessage('位置更新成功！天气信息已刷新', 'success')
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

// 显示消息
const showMessage = (text, type) => {
  message.value = text
  messageType.value = type
  setTimeout(() => {
    message.value = ''
  }, 3000)
}
</script>

<style scoped>
.location-update {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.update-header {
  text-align: center;
  margin-bottom: 25px;
}

.update-header h3 {
  color: white;
  margin-bottom: 10px;
  font-size: 1.3rem;
}

.update-header p {
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
}

.update-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  color: white;
  font-weight: 500;
  font-size: 0.9rem;
}

.form-group input {
  padding: 12px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  font-size: 0.9rem;
}

.form-group input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.coordinates {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 10px;
}

.btn-location,
.btn-update {
  flex: 1;
  padding: 12px 20px;
  border: none;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-location {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.btn-location:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
}

.btn-update {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-update:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
}

.btn-location:disabled,
.btn-update:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.message {
  margin-top: 15px;
  padding: 10px 15px;
  border-radius: 8px;
  font-size: 0.9rem;
  text-align: center;
}

.message.success {
  background: rgba(76, 175, 80, 0.2);
  color: #4CAF50;
  border: 1px solid rgba(76, 175, 80, 0.3);
}

.message.error {
  background: rgba(244, 67, 54, 0.2);
  color: #f44336;
  border: 1px solid rgba(244, 67, 54, 0.3);
}

@media (max-width: 480px) {
  .coordinates {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
