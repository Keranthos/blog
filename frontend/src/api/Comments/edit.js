import { requestFunc } from '../req'

async function createComment (user, blogID, type, content, parentId = null, token = null) {
  try {
    const data = {
      blogID: parseInt(blogID),
      type,
      content,
      username: user.name || user.username
    }

    // 如果有父评论ID，添加parent_id字段
    if (parentId) {
      data.parent_id = parseInt(parentId)
    }

    console.log('发送评论数据:', data)
    const res = await requestFunc('/comments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token || user.token}`
      },
      data
    }, true)
    return res.ok
  } catch (error) {
    console.error('评论提交错误详情:', error)
    console.error('错误响应:', error.response?.data)
    console.error('错误状态码:', error.response?.status)
    throw error
  }
}

async function deleteComment (user, commentId, token = null) {
  try {
    const res = await requestFunc(`/comments/${commentId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token || user.token}`
      }
    }, true)
    return res.ok
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { createComment, deleteComment }
