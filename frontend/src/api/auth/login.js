import { requestFunc } from '../req'

async function login (username, password) { // 返回用户信息（详见vuex）和'jwt'
  try {
    const res = await requestFunc('/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      data: {
        username,
        password
      }
    }, false)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

async function logout () {
  localStorage.removeItem('jwt')
}

export { login, logout }
