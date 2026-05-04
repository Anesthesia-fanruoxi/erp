import request from './request'

// 登录
export function login(data: { userName: string; password: string }) {
  return request.post('/auth/login', data)
}

// 注册
export function register(data: { userName: string; password: string; realName: string; email?: string; phone?: string }) {
  return request.post('/auth/register', data)
}

// 登出
export function logout() {
  return request.post('/auth/logout')
}

// 获取当前用户信息
export function getCurrentUser() {
  return request.get('/auth/current')
}

// 获取动态菜单
export function getMenus() {
  return request.get('/auth/menus')
}

// 获取个人信息
export function getProfile() {
  return request.get('/auth/profile')
}

// 修改个人信息
export function updateProfile(data: { realName: string; email: string; phone: string; password?: string }) {
  return request.put('/auth/profile', data)
}
