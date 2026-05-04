package system

import "erp-service/common"

// LoginReq 登录请求
type LoginReq struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterReq 注册请求
type RegisterReq struct {
	UserName string `json:"userName" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"realName" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// LoginResp 登录响应
type LoginResp struct {
	Token string `json:"token"`
}

// DeviceBindReq 设备绑定请求
type DeviceBindReq struct {
	MachineCode string `json:"machineCode" binding:"required"`
	DeviceName  string `json:"deviceName"`
}

// DeviceLoginReq 设备免登录请求
type DeviceLoginReq struct {
	MachineCode string `json:"machineCode" binding:"required"`
}

// DeviceBinding 设备绑定表模型
type DeviceBinding struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	UserID      uint             `json:"userId" gorm:"column:user_id;uniqueIndex"`
	MachineHash string           `json:"-" gorm:"column:machine_hash;size:255"`
	DeviceName  string           `json:"deviceName" gorm:"column:device_name;size:100"`
	BoundAt     common.MilliTime `json:"boundAt" gorm:"column:bound_at;autoCreateTime"`
}

func (DeviceBinding) TableName() string {
	return common.TableSysDeviceBinding
}
