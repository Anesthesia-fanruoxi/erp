package common

// 状态码常量
const (
	CodeSuccess       = 200
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeInternalError = 500
)

// 用户状态
const (
	UserStatusPending  = 0 // 待审核
	UserStatusActive   = 1 // 已激活
	UserStatusDisabled = 2 // 已禁用
)
