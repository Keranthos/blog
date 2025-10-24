import { requestFunc } from '../req'

async function createArticle (token, article, type) {
  try {
    const res = await requestFunc(`/articles?type=${type}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      data: article
    }, true)

    // 检查响应是否为null（token过期等情况）
    if (!res) {
      console.error('API: 请求失败，可能token已过期')
      throw new Error('请求失败，请重新登录')
    }

    return res.ok
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function updateArticle (token, article, type, articleId) {
  try {
    console.log('API: 更新文章请求', { url: `/articles/${articleId}?type=${type}`, token: token?.substring(0, 20) + '...', article })
    const res = await requestFunc(`/articles/${articleId}?type=${type}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      data: article
    }, true)
    console.log('API: 更新文章响应', res)

    // 检查响应是否为null（token过期等情况）
    if (!res) {
      console.error('API: 请求失败，可能token已过期')
      throw new Error('请求失败，请重新登录')
    }

    return res.ok
  } catch (error) {
    console.error('API: 更新文章错误', error)
    throw error
  }
}

async function delArticle (token, type, postID) {
  try {
    const res = await requestFunc(`/articles/${postID}?type=${type}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      }
    }, true)
    return res.ok
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { createArticle, updateArticle, delArticle }
