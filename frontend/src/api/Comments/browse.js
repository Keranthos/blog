import { requestFunc } from '../req'

async function getCommentsByID (type, blogID) {
  try {
    const res = await requestFunc(`/comments/${blogID}?type=${type}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }, false)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function getAllComments (page = 1, limit = 10) {
  try {
    const res = await requestFunc(`/comments?page=${page}&limit=${limit}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }, false)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { getCommentsByID, getAllComments }
