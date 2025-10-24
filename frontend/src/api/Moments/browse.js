import { requestFunc } from '../req'

async function getMomentsList (page = 1, limit = 10) {
  try {
    const res = await requestFunc(`/moments?page=${page}&limit=${limit}`, {
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

async function getMoment (id) {
  try {
    const res = await requestFunc(`/moments/${id}`, {
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

async function getMomentsNum () {
  try {
    const res = await requestFunc('/moments/count', {
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

export { getMomentsList, getMoment, getMomentsNum }
