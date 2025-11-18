package tiktok

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TkAlert 报警记录模型
// 记录系统各类异常报警信息
type TkAlert struct {
	global.GVA_MODEL
	UserID       uint   `json:"userId" gorm:"not null;index;comment:所属用户ID"`                                                   // 所属用户ID
	DeviceID     *uint  `json:"deviceId" gorm:"index;comment:设备ID"`                                                            // 设备ID（可选）
	AlertType    string `json:"alertType" gorm:"not null;size:50;index;comment:报警类型"`                                          // 报警类型
	AlertLevel   string `json:"alertLevel" gorm:"type:enum('info','warning','error','critical');default:'warning';index;comment:报警级别"` // 报警级别
	AlertMessage string `json:"alertMessage" gorm:"type:text;not null;comment:报警消息"`                                           // 报警消息
	IsRead       *bool  `json:"isRead" gorm:"default:false;index;comment:是否已读"`                                                // 是否已读

	// 关联数据（用于查询时加载）
	Device       *TkDevice `json:"device" gorm:"foreignKey:DeviceID"`                                                          // 关联设备
}

// TableName 指定表名
func (TkAlert) TableName() string {
	return "tk_alerts"
}

// AlertType 报警类型常量
const (
	AlertTypeDeviceOffline     = "device_offline"      // 设备离线
	AlertTypeCPUHigh           = "cpu_high"            // CPU 占用过高
	AlertTypeMemoryHigh        = "memory_high"         // 内存占用过高
	AlertTypeBatteryLow        = "battery_low"         // 电量过低
	AlertTypeTaskFailed        = "task_failed"         // 任务执行失败
	AlertTypeDeviceExpiring    = "device_expiring"     // 设备即将到期
	AlertTypeDeviceExpired     = "device_expired"      // 设备已过期
	AlertTypeScriptError       = "script_error"        // 脚本执行错误
	AlertTypeCustom            = "custom"              // 自定义报警
)

// AlertLevel 报警级别常量
const (
	AlertLevelInfo     = "info"     // 信息
	AlertLevelWarning  = "warning"  // 警告
	AlertLevelError    = "error"    // 错误
	AlertLevelCritical = "critical" // 严重
)

// MarkAsRead 标记为已读
func (a *TkAlert) MarkAsRead() {
	isRead := true
	a.IsRead = &isRead
}

// IsCritical 判断是否为严重报警
func (a *TkAlert) IsCritical() bool {
	return a.AlertLevel == AlertLevelCritical
}
