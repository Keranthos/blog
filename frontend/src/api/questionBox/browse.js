import { requestFunc } from '../req'

async function getQuestions (page = 1, limit = 10) { // 返回问题的列表，每个问题包含"id": 1,"author": "提问者用户名","content": "问题内容","answer": "回答内容"（如果没有回答，可以是 null 或空字符串）
  try {
    const res = await requestFunc(`/questions?page=${page}&limit=${limit}`, {
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

export { getQuestions }
