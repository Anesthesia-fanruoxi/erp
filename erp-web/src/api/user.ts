import request from './request'

export function getUserList(params?: { page?: number; pageSize?: number }) {
  return request.get('/users', { params })
}

export function getUserDetail(id: number) {
  return request.get(`/users/${id}`)
}

export function createUser(data: any) {
  return request.post('/users', data)
}

export function updateUser(id: number, data: any) {
  return request.put(`/users/${id}`, data)
}

export function deleteUser(id: number) {
  return request.delete(`/users/${id}`)
}

export function updateUserStatus(id: number, status: number) {
  return request.put(`/users/${id}/status`, { status })
}
