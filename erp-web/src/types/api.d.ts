// 统一响应体
interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页响应
interface PageData<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 审计日志
interface AuditLog {
  id: number
  userId: number
  userName: string
  action: string
  method: string
  path: string
  query: string
  body: string
  statusCode: number
  ip: string
  duration: number
  createdAt: number
}
