package common

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm/schema"
)

// MilliTime 自定义时间类型
// - 数据库存储: DATETIME (标准时间格式)
// - JSON 序列化: 毫秒时间戳 (int64)
type MilliTime struct {
	time.Time
}

// MarshalJSON 序列化为毫秒时间戳
func (t MilliTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("0"), nil
	}
	return []byte(fmt.Sprintf("%d", t.Time.UnixMilli())), nil
}

// UnmarshalJSON 从毫秒时间戳反序列化
func (t *MilliTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == "0" {
		t.Time = time.Time{}
		return nil
	}
	var ms int64
	if _, err := fmt.Sscanf(s, "%d", &ms); err != nil {
		return fmt.Errorf("MilliTime: 无法解析 %s", s)
	}
	t.Time = time.UnixMilli(ms)
	return nil
}

// Scan 实现 sql.Scanner 接口 (从数据库读取)
func (t *MilliTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	default:
		return fmt.Errorf("MilliTime: 无法从 %T 扫描", value)
	}
}

// Value 实现 driver.Valuer 接口 (写入数据库，始终写入标准 time.Time)
func (t MilliTime) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

// GormDataType 告知 GORM 列类型为 datetime
func (t MilliTime) GormDataType() string {
	return "datetime"
}

// GormDBDataType 告知 GORM 迁移时使用 DATETIME
func (t MilliTime) GormDBDataType(_ interface{}, _ *schema.Field) string {
	return "DATETIME"
}

// MilliTimeNow 返回当前时间的 MilliTime
func MilliTimeNow() MilliTime {
	return MilliTime{Time: time.Now()}
}
