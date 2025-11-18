package tiktok

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
)

// TkTaskResult 任务结果模型 - 核心结果保存
// 根据不同脚本类型保存不同的结果数据
type TkTaskResult struct {
	global.GVA_MODEL
	TaskID     uint           `json:"taskId" gorm:"not null;index;comment:任务ID"`       // 任务ID
	ScriptID   uint           `json:"scriptId" gorm:"not null;index;comment:脚本ID"`     // 脚本ID
	ResultType string         `json:"resultType" gorm:"not null;size:50;index;comment:结果类型"` // 结果类型（register, post_video, login等）
	ResultData common.JSONMap `json:"resultData" gorm:"type:json;not null;comment:结果数据JSON格式"` // 结果数据（JSON 格式）

	// 关联数据（用于查询时加载）
	Task       *TkTask        `json:"task" gorm:"foreignKey:TaskID"`                    // 关联任务
	Script     *TkScript      `json:"script" gorm:"foreignKey:ScriptID"`                // 关联脚本
}

// TableName 指定表名
func (TkTaskResult) TableName() string {
	return "tk_task_results"
}

// ResultType 结果类型常量
const (
	ResultTypeRegister  = "register"   // 注册结果
	ResultTypePostVideo = "post_video" // 发布视频结果
	ResultTypeLogin     = "login"      // 登录结果
	ResultTypeLike      = "like"       // 点赞结果
	ResultTypeComment   = "comment"    // 评论结果
	ResultTypeFollow    = "follow"     // 关注结果
	ResultTypeCustom    = "custom"     // 自定义结果
)

// RegisterResult 注册脚本结果数据结构示例
// 实际使用时，ResultData 字段会存储此结构的 JSON 格式
type RegisterResult struct {
	Account      string `json:"account"`      // 账号
	Password     string `json:"password"`     // 密码
	Phone        string `json:"phone"`        // 手机号
	Email        string `json:"email"`        // 邮箱
	RegisteredAt string `json:"registeredAt"` // 注册时间
}

// PostVideoResult 发布视频脚本结果数据结构示例
type PostVideoResult struct {
	VideoURL    string   `json:"videoUrl"`    // 视频 URL
	Title       string   `json:"title"`       // 标题
	Description string   `json:"description"` // 描述
	Tags        []string `json:"tags"`        // 标签
	PostedAt    string   `json:"postedAt"`    // 发布时间
}

// LoginResult 登录脚本结果数据结构示例
type LoginResult struct {
	Account   string `json:"account"`   // 账号
	LoginAt   string `json:"loginAt"`   // 登录时间
	IPAddress string `json:"ipAddress"` // IP 地址
	Success   bool   `json:"success"`   // 是否成功
}
