import storage from 'store'
import expirePlugin from 'store/plugins/expire'
import { login, getInfo, logout } from '@/api/login'
import { ACCESS_TOKEN } from '@/store/mutation-types'
import { welcome } from '@/utils/util'

storage.addPlugin(expirePlugin)
const user = {
  state: {
    token: '',
    name: '',
    welcome: '',
    avatar: '',
    roles: [],
    info: {}
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_NAME: (state, { name, welcome }) => {
      state.name = name
      state.welcome = welcome
    },
    SET_AVATAR: (state, avatar) => {
      state.avatar = avatar
    },
    SET_ROLES: (state, roles) => {
      state.roles = roles
    },
    SET_INFO: (state, info) => {
      state.info = info
    }
  },

  actions: {
    // 登录
    Login ({ commit }, userInfo) {
      return new Promise((resolve, reject) => {
        login(userInfo).then(response => {
          if (response.code !== 0) {
            reject(response)
            return
          }
          const result = response.info
          storage.set(ACCESS_TOKEN, result.token, new Date().getTime() + 7 * 24 * 60 * 60 * 1000)
          commit('SET_TOKEN', result.token)
          resolve()
        }).catch(error => {
          reject(error)
        })
      })
    },

    // 获取用户信息
    GetInfo ({ commit }) {
      return new Promise((resolve, reject) => {
        // 请求后端获取用户信息 /api/accounts/info
        getInfo().then(response => {
          const result = response.info
          // if (result.role && result.role.permissions.length > 0) {
          //   const role = { ...result.role }
          //   role.permissions = result.role.permissions.map(permission => {
          //     const per = {
          //       ...permission,
          //       actionList: (permission.actionEntitySet || {}).map(item => item.action)
          //      }
          //     return per
          //   })
          //   role.permissionList = role.permissions.map(permission => { return permission.permissionId })
          //   // 覆盖响应体的 role, 供下游使用
          //   result.role = role

          const permissions = {}
          result.actions.forEach(element => {
            if (permissions[element.cate_en] === undefined) {
              permissions[element.cate_en] = {
                permissionId: element.cate_en,
                needResource: element.need_resource,
                actionList: []
              }
            }
            permissions[element.cate_en].actionList.push(element.name)
          })
          commit('SET_ROLES', Object.values(permissions))
          commit('SET_INFO', result)
          commit('SET_NAME', { name: result.name, welcome: welcome() })
          commit('SET_AVATAR', '/avatar2.jpg')
          // 下游
          resolve(result)
          // } else {
          //   reject(new Error('getInfo: roles must be a non-null array !'))
          // }
        }).catch(error => {
          reject(error)
        })
      })
    },

    // 登出
    Logout ({ commit, state }) {
      return new Promise((resolve) => {
        logout(state.token).then(() => {
          commit('SET_TOKEN', '')
          commit('SET_ROLES', [])
          storage.remove(ACCESS_TOKEN)
          resolve()
        }).catch((err) => {
          console.log('logout fail:', err)
          // resolve()
        }).finally(() => {
        })
      })
    }

  }
}

export default user
