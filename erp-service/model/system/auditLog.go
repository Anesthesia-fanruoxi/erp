package system

import "erp-service/common"

// AuditLog 操作审计日志表模型
type AuditLog struct {
	ID         uint             `json:"id" gorm:"primaryKey"`
	UserID     uint             `json:"userId" gorm:"column:user_id;index"`
	UserName   string           `json:"userName" gorm:"column:user_name;size:50"`
	Action     string           `json:"action" gorm:"column:action;size:100"` // 操作描述, 如"创建用户"
	Method     string           `json:"method" gorm:"column:method;size:10"`  // HTTP方法
	Path       string           `json:"path" gorm:"column:path;size:200"`     // 请求路径
	Query      string           `json:"query" gorm:"column:query;type:text"`  // URL查询参数
	Body       string           `json:"body" gorm:"column:body;type:text"`    // 请求体(已脱敏)
	StatusCode int              `json:"statusCode" gorm:"column:status_code"` // 业务响应码
	IP         string           `json:"ip" gorm:"column:ip;size:50"`          // 来源IP
	Duration   int64            `json:"duration" gorm:"column:duration"`      // 耗时(ms)
	CreatedAt  common.MilliTime `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

func (AuditLog) TableName() string {
	return common.TableSysAuditLog
}

// AuditLogListReq 审计日志查询请求
type AuditLogListReq struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
	UserName  string `form:"userName"`  // 按操作人筛选
	Action    string `form:"action"`    // 按操作描述筛选
	StartTime int64  `form:"startTime"` // 开始时间(毫秒时间戳)
	EndTime   int64  `form:"endTime"`   // 结束时间(毫秒时间戳)
}
