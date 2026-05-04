import request from './request'

export interface MenuItem {
  id: number
  code: string
  name: string
  type: number // 1-菜单 2-目录
  parentId: number
  path: string
  icon: string
  sort: number
  visible: number
  children?: MenuItem[]
}

// 获取全部菜单/权限项 (平铺,前端自行组装树)
export function getMenuList() {
  return request.get<any, MenuItem[]>('/menus')
}

export function createMenu(data: Partial<MenuItem>) {
  return request.post('/menus', data)
}

export function updateMenu(id: number, data: Partial<MenuItem>) {
  return request.put(`/menus/${id}`, data)
}

export function deleteMenu(id: number) {
  return request.delete(`/menus/${id}`)
}
