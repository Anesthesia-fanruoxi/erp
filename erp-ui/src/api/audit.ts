import request from './request'

export function getAuditList(params?: { page?: number; pageSize?: number }) {
  return request.get('/registrations', { params })
}

export function approveRegistration(id: number) {
  return request.put(`/registrations/${id}/approve`)
}

export function rejectRegistration(id: number) {
  return request.put(`/registrations/${id}/reject`)
}
