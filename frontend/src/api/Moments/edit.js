import { requestFunc, getJWT } from '../req'

async function createMoment (_user, title, content, image, author, tags = '') {
  try {
    const token = getJWT()
    if (!token) {
      throw new Error('JWT is missing, please re-login')
    }
    const res = await requestFunc('/moments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      data: {
        Title: title,
        Content: content,
        Image: image,
        Author: author,
        Tags: tags
      }
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function updateMoment (_user, momentId, title, content, image, tags = '') {
  try {
    const token = getJWT()
    if (!token) {
      throw new Error('JWT is missing, please re-login')
    }
    const res = await requestFunc(`/moments/${momentId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      data: {
        Title: title,
        Content: content,
        Image: image,
        Tags: tags
      }
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function deleteMoment (_user, momentId) {
  try {
    const token = getJWT()
    if (!token) {
      throw new Error('JWT is missing, please re-login')
    }
    const res = await requestFunc(`/moments/${momentId}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      }
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { createMoment, updateMoment, deleteMoment }
