import { requestFunc } from '../req'

async function submitQuestion (user, content) { // 前端传递"content": "用户输入的问题内容"
  try {
    const res = await requestFunc('/questions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      data: {
        content,
        author: user.name || '匿名用户'
      }
    }, false)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function submitAnswer (user, questionId, answer) { // 前端传递"answer": "用户输入的回答内容"
  try {
    const res = await requestFunc(`/questions/${questionId}/answer`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      data: {
        answer
      }
    }, true) // requestFunc 会自动从 localStorage 获取最新的 token
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { submitQuestion, submitAnswer }
