import { requestFunc } from '../req'

// 返回一篇博客的简介，返回 id: number, title: string, tags: 数组，date: string， img: string
// 返回一篇博客的详情，返回 id: number, title: string, tags: 数组，date: string， img: string, content: string
// 返回一篇博客的评论，返回 id: number, comments: 数组{userinfo: string, content: string}

async function getArticlesList (type = 'blog', page = 1, limit = 9) { // 返回blog数组blogs: Array<{id: number, title: string, tags: string, date: string, img: string}>
  const res = await requestFunc(`/articles?page=${page}&limit=${limit}&type=${type}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  }, false)
  return res.data
}// get方法中不能有函数体，而是用查询参数

async function getArticlesNum (type) { // 返回键值对num: number
  const res = await requestFunc(`/articles/count?type=${type}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  }, false)
  return res.data
}

async function getArticleByID (type = 'blog', id) {
  const res = await requestFunc(`/articles/${id}?type=${type}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  }, false)
  return res.data
}

async function getTopArticles (limit = 3) { // 获取置顶文章
  const res = await requestFunc(`/articles/top?limit=${limit}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  }, false)
  return res.data
}

export { getArticlesList, getArticlesNum, getArticleByID, getTopArticles }
