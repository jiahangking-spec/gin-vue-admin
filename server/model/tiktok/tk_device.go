package tiktok

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TkDevice 设备管理模型 - 核心计费模块
// 一个用户可以绑定多台设备，每台设备独立计费
type TkDevice struct {
	global.GVA_MODEL
	UserID         uint      `json:"userId" gorm:"not null;index;comment:所属用户ID"`                        // 所属用户ID
	DeviceID       string    `json:"deviceId" gorm:"unique;not null;size:100;comment:设备唯一标识(IMEI/设备ID)"` // 设备唯一标识
	DeviceName     string    `json:"deviceName" gorm:"size:100;comment:设备备注名称"`                          // 设备备注名称
	DeviceToken    string    `json:"deviceToken" gorm:"not null;size:255;comment:设备授权Token"`             // 设备授权Token
	Model          string    `json:"model" gorm:"size:100;comment:设备型号"`                                 // 设备型号
	AndroidVersion string    `json:"androidVersion" gorm:"size:50;comment:Android版本"`                    // Android版本
	BindTime       time.Time `json:"bindTime" gorm:"comment:绑定时间"`                                       // 绑定时间
	ExpireTime     time.Time `json:"expireTime" gorm:"not null;index;comment:到期时间-核心计费字段"`              // 到期时间（核心计费字段）
	Status         string    `json:"status" gorm:"type:enum('online','offline','expired');default:'offline';comment:设备状态"` // 设备状态
	LastActiveTime time.Time `json:"lastActiveTime" gorm:"comment:最后活跃时间"`                              // 最后活跃时间
	GroupName      string    `json:"groupName" gorm:"size:100;comment:设备分组名称"`                          // 设备分组（如"注册组"、"养号组"）
}

// TableName 指定表名
func (TkDevice) TableName() string {
	return "tk_devices"
}

// IsExpired 检查设备是否已过期
func (d *TkDevice) IsExpired() bool {
	return time.Now().After(d.ExpireTime)
}

// IsOnline 检查设备是否在线
func (d *TkDevice) IsOnline() bool {
	return d.Status == "online"
}

// CanReceiveTask 检查设备是否可以接收任务
func (d *TkDevice) CanReceiveTask() bool {
	return d.IsOnline() && !d.IsExpired()
}
