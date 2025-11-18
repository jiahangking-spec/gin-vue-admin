package tiktok

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
)

// TkTask 任务管理模型
// 支持单次、定时、循环、批量任务
type TkTask struct {
	global.GVA_MODEL
	UserID        uint           `json:"userId" gorm:"not null;index;comment:所属用户ID"`                                            // 所属用户ID
	ScriptID      uint           `json:"scriptId" gorm:"not null;index;comment:脚本ID"`                                            // 脚本ID
	DeviceID      uint           `json:"deviceId" gorm:"not null;index;comment:执行设备ID"`                                          // 执行设备ID
	TaskName      string         `json:"taskName" gorm:"size:100;comment:任务名称"`                                                   // 任务名称
	TaskType      string         `json:"taskType" gorm:"type:enum('once','scheduled','loop','batch');not null;comment:任务类型"`     // 任务类型
	Params        common.JSONMap `json:"params" gorm:"type:json;comment:任务参数"`                                                    // 任务参数（脚本执行参数）
	Priority      string         `json:"priority" gorm:"type:enum('high','medium','low');default:'medium';comment:优先级"`          // 优先级
	Status        string         `json:"status" gorm:"type:enum('pending','running','success','failed','cancelled');default:'pending';index;comment:任务状态"` // 任务状态
	ScheduledTime *time.Time     `json:"scheduledTime" gorm:"comment:定时执行时间"`                                                    // 定时执行时间
	StartTime     *time.Time     `json:"startTime" gorm:"comment:开始执行时间"`                                                        // 开始执行时间
	EndTime       *time.Time     `json:"endTime" gorm:"comment:结束时间"`                                                            // 结束时间
	RetryCount    int            `json:"retryCount" gorm:"default:0;comment:已重试次数"`                                              // 已重试次数
	MaxRetry      int            `json:"maxRetry" gorm:"default:3;comment:最大重试次数"`                                               // 最大重试次数
	ErrorMessage  string         `json:"errorMessage" gorm:"type:text;comment:错误信息"`                                             // 错误信息
	Duration      int            `json:"duration" gorm:"comment:执行耗时(秒)"`                                                        // 执行耗时（秒）

	// 关联数据（用于查询时加载）
	Script        *TkScript      `json:"script" gorm:"foreignKey:ScriptID"`                                                      // 关联脚本
	Device        *TkDevice      `json:"device" gorm:"foreignKey:DeviceID"`                                                      // 关联设备
}

// TableName 指定表名
func (TkTask) TableName() string {
	return "tk_tasks"
}

// TaskType 任务类型常量
const (
	TaskTypeOnce      = "once"      // 单次任务
	TaskTypeScheduled = "scheduled" // 定时任务
	TaskTypeLoop      = "loop"      // 循环任务
	TaskTypeBatch     = "batch"     // 批量任务
)

// TaskStatus 任务状态常量
const (
	TaskStatusPending   = "pending"   // 等待中
	TaskStatusRunning   = "running"   // 执行中
	TaskStatusSuccess   = "success"   // 成功
	TaskStatusFailed    = "failed"    // 失败
	TaskStatusCancelled = "cancelled" // 已取消
)

// TaskPriority 任务优先级常量
const (
	TaskPriorityHigh   = "high"   // 高优先级
	TaskPriorityMedium = "medium" // 中优先级
	TaskPriorityLow    = "low"    // 低优先级
)

// CanRetry 判断任务是否可以重试
func (t *TkTask) CanRetry() bool {
	return t.Status == TaskStatusFailed && t.RetryCount < t.MaxRetry
}

// IsPending 判断任务是否待执行
func (t *TkTask) IsPending() bool {
	return t.Status == TaskStatusPending
}

// IsRunning 判断任务是否执行中
func (t *TkTask) IsRunning() bool {
	return t.Status == TaskStatusRunning
}

// IsFinished 判断任务是否已完成（成功或失败）
func (t *TkTask) IsFinished() bool {
	return t.Status == TaskStatusSuccess || t.Status == TaskStatusFailed || t.Status == TaskStatusCancelled
}
