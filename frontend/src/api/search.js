import { requestFunc } from './req'

async function searchArticles (keyword, type = 'all') {
  try {
    const res = await requestFunc(`/search?keyword=${encodeURIComponent(keyword)}&type=${type}`, {
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

export { searchArticles }
