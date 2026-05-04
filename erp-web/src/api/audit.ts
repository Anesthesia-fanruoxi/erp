import request from './request'

// 注册审核
export function getAuditList(params?: { page?: number; pageSize?: number }) {
  return request.get('/registrations', { params })
}

export function approveRegistration(id: number) {
  return request.put(`/registrations/${id}/approve`)
}

export function rejectRegistration(id: number) {
  return request.put(`/registrations/${id}/reject`)
}

// 操作审计日志
export interface AuditLogQuery {
  page?: number
  pageSize?: number
  userName?: string
  action?: string
  startTime?: number
  endTime?: number
}

export function getAuditLogList(params?: AuditLogQuery) {
  return request.get('/audit/logs', { params })
}
