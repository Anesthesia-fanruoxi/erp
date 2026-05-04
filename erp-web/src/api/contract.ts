import request from './request'

export interface ContractItem {
  id?: number
  seq: number
  name: string
  spec: string
  brand: string
  qty: string
  unit: string
  unitPrice: string
  amount: string
  operator: string
  location: string
  remark: string
}

export interface ContractForm {
  projectName: string
  orderNo: string
  orderDate: string
  fromCompany: string
  toCompany: string
  buyer: string
  attn: string
  buyerEmail: string
  buyerTel: string
  attnTel: string
  totalAmount: string
  deliveryAddr: string
  remark: string
  items: ContractItem[]
}

export function getContractList(params?: any) {
  return request.get('/contracts', { params })
}

export function getContractDetail(id: number) {
  return request.get(`/contracts/${id}`)
}

export function createContract(data: ContractForm) {
  return request.post('/contracts', data)
}

export function updateContract(id: number, data: ContractForm) {
  return request.put(`/contracts/${id}`, data)
}

export function deleteContract(id: number) {
  return request.delete(`/contracts/${id}`)
}
