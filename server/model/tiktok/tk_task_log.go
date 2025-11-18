package tiktok

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TkTaskLog 任务日志模型
// 记录任务执行的详细日志，支持实时推送
type TkTaskLog struct {
	global.GVA_MODEL
	TaskID        uint   `json:"taskId" gorm:"not null;index;comment:任务ID"`                              // 任务ID
	DeviceID      uint   `json:"deviceId" gorm:"not null;index;comment:设备ID"`                            // 设备ID
	LogLevel      string `json:"logLevel" gorm:"type:enum('info','warn','error','debug');default:'info';comment:日志级别"` // 日志级别
	LogContent    string `json:"logContent" gorm:"type:text;not null;comment:日志内容"`                       // 日志内容
	ScreenshotURL string `json:"screenshotUrl" gorm:"size:255;comment:截图URL"`                             // 截图 URL（关键步骤和异常截图）

	// 关联数据（用于查询时加载）
	Task          *TkTask   `json:"task" gorm:"foreignKey:TaskID"`                                        // 关联任务
	Device        *TkDevice `json:"device" gorm:"foreignKey:DeviceID"`                                    // 关联设备
}

// TableName 指定表名
func (TkTaskLog) TableName() string {
	return "tk_task_logs"
}

// LogLevel 日志级别常量
const (
	LogLevelInfo  = "info"  // 信息
	LogLevelWarn  = "warn"  // 警告
	LogLevelError = "error" // 错误
	LogLevelDebug = "debug" // 调试
)

// IsError 判断是否为错误日志
func (l *TkTaskLog) IsError() bool {
	return l.LogLevel == LogLevelError
}

// HasScreenshot 判断是否有截图
func (l *TkTaskLog) HasScreenshot() bool {
	return l.ScreenshotURL != ""
}
