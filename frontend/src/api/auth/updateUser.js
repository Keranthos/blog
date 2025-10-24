import { requestFunc } from '../req'

async function updateUser (username, password, avatar) {
  try {
    const res = await requestFunc('/auth/updateUser', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      data: {
        username,
        password,
        avatar
      }
    }, true)
    return res.data
  } catch (error) {
    console.error(error)
    throw error
  }
}

export { updateUser }
