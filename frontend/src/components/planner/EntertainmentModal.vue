<template>
  <Teleport to="body">
    <div class="modal-overlay" @click="close">
      <div class="modal-container" @click.stop>
        <div class="modal-header">
          <h3>{{ entertainment ? '编辑娱乐活动' : '添加娱乐活动' }}</h3>
          <button class="close-btn" @click="close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>活动类型 *</label>
            <select v-model="form.type">
              <option value="">选择类型</option>
              <option value="游戏">游戏</option>
              <option value="追剧">追剧</option>
              <option value="阅读">阅读</option>
              <option value="电影">电影</option>
              <option value="运动">运动</option>
              <option value="其他">其他</option>
            </select>
          </div>
          <div class="form-group">
            <label>时长（分钟） *</label>
            <input v-model.number="form.duration" type="number" min="1" placeholder="输入时长" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>花费（元）</label>
              <input v-model.number="form.cost" type="number" min="0" step="0.01" placeholder="0" />
            </div>
            <div class="form-group">
              <label>评价（1-5星）</label>
              <input v-model.number="form.rating" type="number" min="0" max="5" placeholder="0" />
            </div>
          </div>
          <div class="form-group">
            <label>备注</label>
            <textarea v-model="form.notes" placeholder="输入备注（可选）" rows="3"></textarea>
          </div>
          <div class="form-group">
            <label>日期</label>
            <input v-model="form.date" type="date" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="close">取消</button>
          <button class="btn btn-primary" @click="save">保存</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  entertainment: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'save'])

const form = ref({
  type: '',
  duration: 0,
  cost: 0,
  rating: 0,
  notes: '',
  date: new Date().toISOString().split('T')[0]
})

onMounted(() => {
  if (props.entertainment) {
    form.value = {
      type: props.entertainment.type || '',
      duration: props.entertainment.duration || 0,
      cost: props.entertainment.cost || 0,
      rating: props.entertainment.rating || 0,
      notes: props.entertainment.notes || '',
      date: props.entertainment.date ? new Date(props.entertainment.date).toISOString().split('T')[0] : new Date().toISOString().split('T')[0]
    }
  }
})

const close = () => {
  emit('close')
}

const save = () => {
  if (!form.value.type.trim()) {
    alert('请选择活动类型')
    return
  }
  if (!form.value.duration || form.value.duration <= 0) {
    alert('请输入有效的时长')
    return
  }
  emit('save', {
    ...form.value,
    date: form.value.date
  })
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
}

.modal-container {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  font-size: 24px;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: #f0f0f0;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #1a1a1a;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #a855f7;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #e0e0e0;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.4);
}

.btn-secondary {
  background: #f0f0f0;
  color: #666;
}

.btn-secondary:hover {
  background: #e0e0e0;
}
</style>
