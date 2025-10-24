import { requestFunc } from '../req'

async function register (username, password, avatar) {
  try {
    const res = await requestFunc('/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      data: {
        username,
        password,
        avatar
      }
    }, false)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { register }
