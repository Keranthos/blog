import { requestFunc } from '../req'

/**
 * 获取文章的点赞状态和数量
 * @param {string} articleType - 文章类型 (blog/research/project/moment)
 * @param {number} articleID - 文章ID
 * @returns {Promise<{isLiked: boolean, likeCount: number}>}
 */
async function getLikeStatus (articleType, articleID) {
  try {
    const res = await requestFunc(`/likes/${articleType}/${articleID}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }, false)
    return res.data
  } catch (error) {
    console.error('获取点赞状态失败:', error)
    throw error
  }
}

/**
 * 切换点赞状态（点赞/取消点赞）
 * @param {string} articleType - 文章类型 (blog/research/project/moment)
 * @param {number} articleID - 文章ID
 * @param {Object} user - 用户对象
 * @param {string} token - 用户token
 * @returns {Promise<{isLiked: boolean, likeCount: number}>}
 */
async function toggleLike (articleType, articleID, user, token) {
  try {
    const res = await requestFunc(`/likes/${articleType}/${articleID}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token || user.token}`
      }
    }, true)
    return res.data
  } catch (error) {
    console.error('切换点赞状态失败:', error)
    if (error.response) {
      console.error('错误响应数据:', error.response.data)
      console.error('错误状态码:', error.response.status)
    }
    throw error
  }
}

export { getLikeStatus, toggleLike }
