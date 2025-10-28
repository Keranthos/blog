import { requestFunc } from '../req'

// 获取讲演列表
export const getPresentations = (page = 1, pageSize = 9) => {
  return requestFunc(`/presentations?page=${page}&page_size=${pageSize}`, {
    method: 'GET'
  }, false)
}

// 获取讲演详情
export const getPresentation = (id) => {
  return requestFunc(`/presentations/${id}`, {
    method: 'GET'
  }, false)
}

// 创建讲演
export const createPresentation = (formData) => {
  return requestFunc('/presentations', {
    method: 'POST',
    data: formData
  }, true)
}

// 更新讲演
export const updatePresentation = (id, data) => {
  return requestFunc(`/presentations/${id}`, {
    method: 'PUT',
    data
  }, true)
}

// 删除讲演
export const deletePresentation = (id) => {
  return requestFunc(`/presentations/${id}`, {
    method: 'DELETE'
  }, true)
}
