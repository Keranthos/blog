import { requestFunc } from '../req'

async function createComment (user, blogID, type, content) {
  try {
    const res = await requestFunc('/comments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${user.token}`
      },
      data: {
        blog_id: blogID,
        type,
        content,
        username: user.username
      }
    }, true)
    return res.ok
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function deleteComment (user, commentId) {
  try {
    const res = await requestFunc(`/comments/${commentId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${user.token}`
      }
    }, true)
    return res.ok
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { createComment, deleteComment }
