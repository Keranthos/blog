<template>
  <div class="auth-container">
    <div class="auth-box">
      <!--切换面板容器-->
      <div class="panels-container">
        <!-- 登录界面 -->
        <div v-if="!isUserLoggedIn" class="panel login-panel">
          <h2>🔒 登入账号</h2>
          <input v-model="loginUsername" type="text" placeholder="名字" />
          <input v-model="loginPassword" type="password" placeholder="密码" />
          <button class="auth-button" @click="handleLogin">🔑 登录</button>
          <p v-if="loginErrorMessage" class="error-message">{{ loginErrorMessage }}</p>
        </div>

        <!-- 注册界面 -->
        <div v-if="!isUserLoggedIn" class="panel register-panel">
          <h2>📝 创建账号</h2>
          <input v-model="registerUsername" type="text" placeholder="名字" />
          <input v-model="registerPassword" type="password" placeholder="密码" />
          <input v-model="registerPasswordConfirm" type="password" placeholder="重复密码" />
          <input v-model="registerAvatar" type="text" placeholder="头像链接（可选）" />
          <button class="auth-button" @click="handleRegister">📝 注册</button>
          <p v-if="registerErrorMessage" class="error-message">{{ registerErrorMessage }}</p>
        </div>

        <!-- 用户信息更新界面 -->
        <div v-if="isUserLoggedIn" class="panel user-info-panel">
          <h2>🔧 更新信息</h2>
          <input v-model="username" type="text" placeholder="用户名" />
          <input v-model="password" type="password" placeholder="新密码" />
          <input v-model="passwordConfirm" type="password" placeholder="确认新密码" />
          <input v-model="avatar" type="text" placeholder="头像链接" />
          <button class="auth-button" @click="handleUpdateUserInfo">更新信息</button>
          <p v-if="updateErrorMessage" class="error-message">{{ updateErrorMessage }}</p>
        </div>
      </div>

      <!-- 切换按钮区域 -->
      <div v-if="!isUserLoggedIn" class="switch-panel" :class="{ 'shift-left': isRegisterMode }">
        <div v-if="!isRegisterMode">
          <h2>欢迎来到山角函兽的小窝</h2>
          <button class="switch-button" @click="togglePanel">📝 注册</button>
        </div>
        <div v-else>
          <h2>Stay Hungry, Stay Foolish</h2>
          <button class="switch-button" @click="togglePanel">🔑 登录</button>
        </div>
        <button class="back-button" @click="goBack">🔙 返回</button>
      </div>
      <div v-else class="switch-panel" :class="{ 'shift-left': isRegisterMode }">
        <div>
          <h2>Stay Hungry, Stay Foolish</h2>
          <button class="switch-button" @click="handleLogout">🔑 登出</button>
        </div>
        <button class="back-button" @click="goBack">🔙 返回</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

const store = useStore()
const router = useRouter()

const isRegisterMode = ref(false)
const isUserLoggedIn = computed(() => store.state.user.isLogged)
const loginUsername = ref('')
const loginPassword = ref('')
const registerUsername = ref('')
const registerPassword = ref('')
const registerPasswordConfirm = ref('')
const registerAvatar = ref('')
const username = ref('')
const password = ref('')
const passwordConfirm = ref('')
const avatar = ref('')
const loginErrorMessage = ref('')
const registerErrorMessage = ref('')
const loading = ref(false)

const handleLogin = async () => {
  loginErrorMessage.value = ''
  try {
    loading.value = true
    await store.dispatch('login', {
      username: loginUsername.value,
      password: loginPassword.value
    })
    loginErrorMessage.value = '登录成功'
    router.push('/')
  } catch (error) {
    loginErrorMessage.value = '登录失败'
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  if (registerPassword.value.length < 6) {
    registerErrorMessage.value = '密码至少六位'
    return
  }
  if (registerPassword.value !== registerPasswordConfirm.value) {
    registerErrorMessage.value = '两次密码不一致'
    return
  }
  try {
    await store.dispatch('register', {
      username: registerUsername.value,
      password: registerPassword.value,
      avatar: registerAvatar.value || ''
    })
    router.push('/')
  } catch (error) {
    if (error.message === 'Username already exists') {
      registerErrorMessage.value = '用户名已存在'
    } else {
      registerErrorMessage.value = '注册失败'
    }
    console.error('Registration failed:', error)
  }
}

const handleLogout = () => {
  store.dispatch('logout')
  router.push('/')
}

const handleUpdateUserInfo = async () => {
  if (password.value && password.value !== passwordConfirm.value) {
    registerErrorMessage.value = '新密码和确认密码不一致'
    return
  }
  try {
    await store.dispatch('updateUserInfo', {
      username: username.value,
      password: password.value,
      avatar: avatar.value
    })
    registerErrorMessage.value = '信息更新成功'
  } catch (error) {
    console.error('Update failed:', error)
  }
}

const goBack = () => {
  router.push('/')
}

const togglePanel = () => {
  isRegisterMode.value = !isRegisterMode.value
}
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  width: 100%;
  padding: 20px;
}

.auth-box {
  position: relative;
  width: 100%;
  max-width: 1000px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 30px;
  overflow: hidden;
  display: flex;
  box-shadow: 0 20px 80px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  min-height: 600px;
}

.panels-container {
  display: flex;
  width: 200%;
  height: 100%;
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.panel {
  width: 50%;
  padding: 3rem 2rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  box-sizing: border-box;
  gap: 15px;
}

.panel h2 {
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 20px;
}

input {
  width: 85%;
  max-width: 400px;
  padding: 15px 20px;
  margin: 8px 0;
  border-radius: 25px;
  border: 2px solid #e8eaf0;
  font-size: 1rem;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.8);
}

input:focus {
  outline: none;
  border-color: #a855f7;
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.1);
  background: white;
}

input::placeholder {
  color: #aaa;
}

.auth-button {
  width: 85%;
  max-width: 400px;
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  color: white;
  padding: 15px 2rem;
  border: none;
  border-radius: 25px;
  cursor: pointer;
  margin-top: 1.5rem;
  font-size: 1.1rem;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 6px 25px rgba(168, 85, 247, 0.4);
}

.auth-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 10px 35px rgba(168, 85, 247, 0.6);
}

.switch-panel {
  position: absolute;
  top: 0;
  right: 0;
  width: 50%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.1) 0%, rgba(124, 58, 237, 0.1) 100%);
  backdrop-filter: blur(10px);
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
  padding: 3rem;
  gap: 30px;
}

.switch-panel.shift-left {
  transform: translateX(-100%);
}

.switch-panel h2 {
  font-size: 1.8rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 20px;
}

.switch-button {
  background: white;
  color: #a855f7;
  padding: 12px 40px;
  border: 2px solid #a855f7;
  border-radius: 25px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 600;
  transition: all 0.3s ease;
}

.switch-button:hover {
  background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%);
  color: white;
  border-color: transparent;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(168, 85, 247, 0.4);
}

.back-button {
  background: rgba(187, 187, 187, 0.2);
  color: #666;
  padding: 10px 30px;
  border: 2px solid #ddd;
  border-radius: 25px;
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 500;
  transition: all 0.3s ease;
}

.back-button:hover {
  background: #ddd;
  color: #333;
  transform: translateY(-2px);
}

.error-message {
  color: #ff4d4f;
  margin-top: 1rem;
  font-size: 0.9rem;
  padding: 10px 20px;
  background: rgba(255, 77, 79, 0.1);
  border-radius: 15px;
  border-left: 3px solid #ff4d4f;
}

/* 响应式 */
@media (max-width: 768px) {
  .auth-box {
    flex-direction: column;
    max-width: 90%;
  }

  .panels-container {
    width: 100%;
  }

  .panel {
    width: 100%;
    padding: 2rem 1rem;
  }

  .switch-panel {
    position: relative;
    width: 100%;
    padding: 2rem;
  }

  .switch-panel.shift-left {
    transform: none;
  }
}
</style>
