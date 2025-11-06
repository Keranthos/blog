import { requestFunc } from '@/api/req'

const baseURL = '/planner'

// ==================== 任务相关 ====================

// 获取任务列表
export async function getTasks (date = null) {
  const response = await requestFunc(
    `${baseURL}/tasks${date ? `?date=${date}` : ''}`,
    {
      method: 'GET',
      headers: {}
    },
    true
  )
  return response?.data
}

// 创建任务
export async function createTask (taskData) {
  const response = await requestFunc(
    `${baseURL}/tasks`,
    {
      method: 'POST',
      data: taskData
    },
    true
  )
  return response?.data
}

// 更新任务
export async function updateTask (id, taskData) {
  const response = await requestFunc(
    `${baseURL}/tasks/${id}`,
    {
      method: 'PUT',
      data: taskData
    },
    true
  )
  return response?.data
}

// 删除任务
export async function deleteTask (id) {
  const response = await requestFunc(
    `${baseURL}/tasks/${id}`,
    {
      method: 'DELETE'
    },
    true
  )
  return response?.data
}

// ==================== 娱乐活动相关 ====================

// 获取娱乐活动列表
export async function getEntertainments (date = null) {
  const response = await requestFunc(
    `${baseURL}/entertainments${date ? `?date=${date}` : ''}`,
    {
      method: 'GET',
      headers: {}
    },
    true
  )
  return response?.data
}

// 创建娱乐活动
export async function createEntertainment (entertainmentData) {
  const response = await requestFunc(
    `${baseURL}/entertainments`,
    {
      method: 'POST',
      data: entertainmentData
    },
    true
  )
  return response?.data
}

// 更新娱乐活动
export async function updateEntertainment (id, entertainmentData) {
  const response = await requestFunc(
    `${baseURL}/entertainments/${id}`,
    {
      method: 'PUT',
      data: entertainmentData
    },
    true
  )
  return response?.data
}

// 删除娱乐活动
export async function deleteEntertainment (id) {
  const response = await requestFunc(
    `${baseURL}/entertainments/${id}`,
    {
      method: 'DELETE'
    },
    true
  )
  return response?.data
}

// ==================== 每日反思相关 ====================

// 获取每日反思
export async function getReflection (date = null) {
  const response = await requestFunc(
    `${baseURL}/reflection${date ? `?date=${date}` : ''}`,
    {
      method: 'GET',
      headers: {}
    },
    true
  )
  return response?.data
}

// 创建或更新每日反思
export async function createOrUpdateReflection (reflectionData) {
  const response = await requestFunc(
    `${baseURL}/reflection`,
    {
      method: 'POST',
      data: reflectionData
    },
    true
  )
  return response?.data
}

// 获取反思列表
export async function getReflections (limit = 30) {
  const response = await requestFunc(
    `${baseURL}/reflections?limit=${limit}`,
    {
      method: 'GET',
      headers: {}
    },
    true
  )
  return response?.data
}

// ==================== 截止日期相关 ====================

// 获取截止日期列表
export async function getDeadlines (completed = null) {
  const params = completed !== null ? `?completed=${completed}` : ''
  const response = await requestFunc(
    `${baseURL}/deadlines${params}`,
    {
      method: 'GET',
      headers: {}
    },
    true
  )
  return response?.data
}

// 创建截止日期
export async function createDeadline (deadlineData) {
  const response = await requestFunc(
    `${baseURL}/deadlines`,
    {
      method: 'POST',
      data: deadlineData
    },
    true
  )
  return response?.data
}

// 更新截止日期
export async function updateDeadline (id, deadlineData) {
  const response = await requestFunc(
    `${baseURL}/deadlines/${id}`,
    {
      method: 'PUT',
      data: deadlineData
    },
    true
  )
  return response?.data
}

// 删除截止日期
export async function deleteDeadline (id) {
  const response = await requestFunc(
    `${baseURL}/deadlines/${id}`,
    {
      method: 'DELETE'
    },
    true
  )
  return response?.data
}
