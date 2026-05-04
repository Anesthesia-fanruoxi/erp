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
