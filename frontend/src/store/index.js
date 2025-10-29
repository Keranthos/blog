import { createStore } from 'vuex'
import { login } from '@/api/auth/login'
import { register } from '@/api/auth/register'
import { updateUser } from '@/api/auth/updateUser'
import createPersistedState from 'vuex-persistedstate'

const store = createStore({
  state: {
    user: {
      isLogged: false,
      name: '',
      level: 1,
      avatar: ''
    },
    token: localStorage.getItem('jwt') || ''
  },
  mutations: {
    setUser (state, user) {
      state.user.isLogged = true
      state.user.name = user.name
      state.user.level = user.level
      state.user.avatar = user.avatar
    },
    setToken (state, token) {
      state.token = token
      localStorage.setItem('jwt', token)
    },
    logout (state) {
      state.user = {
        isLogged: false,
        name: '',
        level: 0,
        avatar: ''
      }
      state.token = ''
      localStorage.removeItem('jwt')
      localStorage.removeItem('vuex') // 清除 vuex 持久化数据
    }
  },
  plugins: [createPersistedState()], // 启用持久化插件
  actions: { // action 中调用 commit commit 会执行一个定义好的 mutation mutation 是唯一允许直接修改 state 的地方
    async login ({ commit }, { username, password }) { // commit为了调用mutations中的方法，后面是传入的参数
      const data = await login(username, password)// commit是Vuex提供的参数，由Vuex自动传入，不需要显示传入
      commit('setUser', data.user)// ...user是展开运算符，将user对象中的所有属性展开放到新对象中,然后再将isLogged属性覆盖
      commit('setToken', data.jwt)
    },
    async register ({ commit }, { username, password, avatar }) {
      const data = await register(username, password, avatar)
      commit('setUser', data.user)
      commit('setToken', data.jwt)
    },
    async updateUser ({ commit }, { username, password, avatar }) {
      const data = await updateUser(username, password, avatar)
      commit('setUser', data.user)
      commit('setToken', data.jwt)
    },
    logout ({ commit }) {
      commit('logout')
    }
  }
})

export default store
