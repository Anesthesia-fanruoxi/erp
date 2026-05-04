import request from './request'

export function getRoleList() {
  return request.get('/roles')
}

export function createRole(data: any) {
  return request.post('/roles', data)
}

export function updateRole(id: number, data: any) {
  return request.put(`/roles/${id}`, data)
}

export function deleteRole(id: number) {
  return request.delete(`/roles/${id}`)
}

// 获取全部权限(菜单+按钮)
export function getPermissionList() {
  return request.get('/permissions')
}
