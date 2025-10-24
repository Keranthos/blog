import { requestFunc } from '../req'

async function getMediaList (type, page = 1, limit = 9) {
  try {
    const res = await requestFunc(`/media?type=${type}&page=${page}&limit=${limit}`, {
      method: 'GET',
      header: {
        'Content-Type': 'application/json'
      }
    }, false)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function getMediaByID (id, type) {
  try {
    const res = await requestFunc(`/media/${id}?type=${type}`, {
      method: 'GET',
      header: {
        'Content-Type': 'application/json'
      }
    }, false)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { getMediaList, getMediaByID }
