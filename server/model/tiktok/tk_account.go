package tiktok

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TkAccount TikTok 账号资产管理模型 - 核心资产
// 存储通过自动化脚本注册或批量导入的 TikTok 账号
type TkAccount struct {
	global.GVA_MODEL
	UserID            uint       `json:"userId" gorm:"not null;index;comment:所属用户ID"`                           // 所属用户ID
	DeviceID          *uint      `json:"deviceId" gorm:"index;comment:注册设备ID"`                                  // 注册设备ID（可选）
	Source            string     `json:"source" gorm:"type:enum('auto_register','manual_import');not null;comment:账号来源"` // 账号来源
	Username          string     `json:"username" gorm:"not null;size:100;index;comment:TikTok用户名"`             // TikTok 用户名
	Phone             string     `json:"phone" gorm:"size:20;index;comment:手机号"`                                // 手机号
	Email             string     `json:"email" gorm:"size:100;index;comment:邮箱"`                                // 邮箱
	PasswordEncrypted string     `json:"passwordEncrypted" gorm:"not null;size:255;comment:AES加密的密码"`           // AES 加密的密码
	RegisteredAt      *time.Time `json:"registeredAt" gorm:"comment:注册时间"`                                      // 注册时间
	LastLoginAt       *time.Time `json:"lastLoginAt" gorm:"comment:最后登录时间"`                                     // 最后登录时间
	Status            string     `json:"status" gorm:"type:enum('normal','banned','pending');default:'normal';index;comment:账号状态"` // 账号状态
	Remark            string     `json:"remark" gorm:"type:text;comment:备注"`                                     // 备注

	// 关联数据（用于查询时加载）
	Device            *TkDevice  `json:"device" gorm:"foreignKey:DeviceID"`                                      // 关联设备
}

// TableName 指定表名
func (TkAccount) TableName() string {
	return "tk_accounts"
}

// AccountSource 账号来源常量
const (
	AccountSourceAutoRegister  = "auto_register"  // 自动注册
	AccountSourceManualImport  = "manual_import"  // 手动导入
)

// AccountStatus 账号状态常量
const (
	AccountStatusNormal  = "normal"  // 正常
	AccountStatusBanned  = "banned"  // 封禁
	AccountStatusPending = "pending" // 待验证
)

// IsNormal 判断账号是否正常
func (a *TkAccount) IsNormal() bool {
	return a.Status == AccountStatusNormal
}

// IsBanned 判断账号是否被封禁
func (a *TkAccount) IsBanned() bool {
	return a.Status == AccountStatusBanned
}
