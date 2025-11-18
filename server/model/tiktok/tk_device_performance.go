package tiktok

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TkDevicePerformance 设备性能监控模型
// 定时上报设备性能数据（CPU、内存、电量、网络等）
type TkDevicePerformance struct {
	global.GVA_MODEL
	DeviceID       uint      `json:"deviceId" gorm:"not null;index;comment:设备ID"`         // 设备ID
	CPUUsage       float64   `json:"cpuUsage" gorm:"type:decimal(5,2);comment:CPU使用率%"`   // CPU 使用率（%）
	MemoryUsage    float64   `json:"memoryUsage" gorm:"type:decimal(5,2);comment:内存使用率%"` // 内存使用率（%）
	BatteryLevel   int       `json:"batteryLevel" gorm:"comment:电量%"`                     // 电量（%）
	NetworkType    string    `json:"networkType" gorm:"size:20;comment:网络类型"`             // 网络类型（WiFi/4G/5G）
	StorageFree    int64     `json:"storageFree" gorm:"comment:可用存储空间MB"`                 // 可用存储空间（MB）
	ReportedAt     time.Time `json:"reportedAt" gorm:"index;comment:上报时间"`               // 上报时间

	// 关联数据（用于查询时加载）
	Device         *TkDevice `json:"device" gorm:"foreignKey:DeviceID"`                   // 关联设备
}

// TableName 指定表名
func (TkDevicePerformance) TableName() string {
	return "tk_device_performance"
}

// NetworkType 网络类型常量
const (
	NetworkTypeWiFi = "WiFi" // WiFi
	NetworkType4G   = "4G"   // 4G
	NetworkType5G   = "5G"   // 5G
	NetworkTypeOther = "Other" // 其他
)

// IsCPUHigh 判断 CPU 使用率是否过高（>80%）
func (p *TkDevicePerformance) IsCPUHigh() bool {
	return p.CPUUsage > 80.0
}

// IsMemoryHigh 判断内存使用率是否过高（>80%）
func (p *TkDevicePerformance) IsMemoryHigh() bool {
	return p.MemoryUsage > 80.0
}

// IsBatteryLow 判断电量是否过低（<20%）
func (p *TkDevicePerformance) IsBatteryLow() bool {
	return p.BatteryLevel < 20
}

// HasPerformanceIssue 判断是否存在性能问题
func (p *TkDevicePerformance) HasPerformanceIssue() bool {
	return p.IsCPUHigh() || p.IsMemoryHigh() || p.IsBatteryLow()
}
