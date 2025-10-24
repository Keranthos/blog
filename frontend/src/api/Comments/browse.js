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

export { getCommentsByID }
