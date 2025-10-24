import { requestFunc } from '../req'

async function createMoment (user, title, content, image, author) {
  try {
    const res = await requestFunc('/moments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${user.token}`
      },
      data: {
        Title: title,
        Content: content,
        Image: image,
        Author: author
      }
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function updateMoment (user, momentId, title, content, image) {
  try {
    const res = await requestFunc(`/moments/${momentId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${user.token}`
      },
      data: {
        Title: title,
        Content: content,
        Image: image
      }
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function deleteMoment (user, momentId) {
  try {
    const res = await requestFunc(`/moments/${momentId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${user.token}`
      }
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { createMoment, updateMoment, deleteMoment }
