import { requestFunc } from '../req'

async function submitMedia (type, media) {
  try {
    const res = await requestFunc(`/media?type=${type}`, {
      method: 'POST',
      header: {
        'Content-Type': 'application/json'
      },
      data: media // 直接发送数据，不嵌套
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

// 创建媒体的别名函数
async function createMedia (user, media, type) {
  return submitMedia(type, media)
}

async function updateMedia (type, mediaId, media) {
  try {
    const res = await requestFunc(`/media/${mediaId}?type=${type}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      data: media // 直接发送数据，不嵌套
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function deleteMedia (type, mediaId) {
  try {
    const res = await requestFunc(`/media/${mediaId}?type=${type}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      }
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { submitMedia, createMedia, updateMedia, deleteMedia }
