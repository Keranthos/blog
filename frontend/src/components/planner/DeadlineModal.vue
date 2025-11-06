<template>
  <Teleport to="body">
    <div class="modal-overlay" @click="close">
      <div class="modal-container" @click.stop>
        <div class="modal-header">
          <h3>{{ deadline ? '编辑DDL' : '添加DDL' }}</h3>
          <button class="close-btn" @click="close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标题 *</label>
            <input v-model="form.title" type="text" placeholder="输入DDL标题" />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="form.description" placeholder="输入描述（可选）" rows="3"></textarea>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>截止日期 *</label>
              <input v-model="form.dueDate" type="date" />
            </div>
            <div class="form-group">
              <label>截止时间</label>
              <input v-model="form.dueTime" type="time" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>优先级</label>
              <select v-model="form.priority">
                <option :value="1">低</option>
                <option :value="2">中</option>
                <option :value="3">高</option>
              </select>
            </div>
            <div class="form-group">
              <label>分类</label>
              <input v-model="form.category" type="text" placeholder="工作、学习、生活等" />
            </div>
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
import { createDeadline, updateDeadline } from '@/api/planner'

const props = defineProps({
  deadline: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'save'])

const form = ref({
  title: '',
  description: '',
  dueDate: new Date().toISOString().split('T')[0],
  dueTime: '',
  priority: 2,
  category: ''
})

onMounted(() => {
  if (props.deadline) {
    form.value = {
      title: props.deadline.title || '',
      description: props.deadline.description || '',
      dueDate: props.deadline.dueDate ? new Date(props.deadline.dueDate).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
      dueTime: props.deadline.dueTime || '',
      priority: props.deadline.priority || 2,
      category: props.deadline.category || ''
    }
  }
})

const close = () => {
  emit('close')
}

const save = async () => {
  if (!form.value.title.trim()) {
    alert('请填写标题')
    return
  }
  if (!form.value.dueDate) {
    alert('请选择截止日期')
    return
  }

  try {
    if (props.deadline) {
      await updateDeadline(props.deadline.ID, form.value)
    } else {
      await createDeadline(form.value)
    }
    emit('save')
  } catch (error) {
    console.error('保存DDL失败:', error)
  }
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
