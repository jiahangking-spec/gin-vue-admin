package tiktok

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
)

// TkScript 脚本管理模型
// 所有脚本由管理员统一上传和维护，基于 EasyClick 框架
type TkScript struct {
	global.GVA_MODEL
	Name         string           `json:"name" gorm:"not null;size:100;index;comment:脚本名称"`            // 脚本名称
	Category     string           `json:"category" gorm:"not null;size:50;index;comment:脚本分类"`          // 脚本分类（注册、登录、发布视频等）
	Version      string           `json:"version" gorm:"size:20;default:'1.0.0';comment:脚本版本"`         // 脚本版本
	Description  string           `json:"description" gorm:"type:text;comment:脚本描述"`                   // 脚本描述
	FilePath     string           `json:"filePath" gorm:"not null;size:255;comment:脚本文件路径"`            // 脚本文件路径（云端存储）
	ParamsSchema common.JSONMap   `json:"paramsSchema" gorm:"type:json;comment:参数配置Schema"`            // 参数配置 Schema (JSON Schema)
	DefaultParams common.JSONMap  `json:"defaultParams" gorm:"type:json;comment:默认参数值"`                // 默认参数值
	IsEnabled    *bool            `json:"isEnabled" gorm:"default:true;comment:是否启用"`                  // 是否启用
	Author       string           `json:"author" gorm:"size:50;comment:脚本作者"`                          // 脚本作者
	Tags         string           `json:"tags" gorm:"size:255;comment:脚本标签(逗号分隔)"`                    // 脚本标签
}

// TableName 指定表名
func (TkScript) TableName() string {
	return "tk_scripts"
}

// ScriptCategory 脚本分类常量
const (
	CategoryRegister    = "register"     // TikTok 注册
	CategoryLogin       = "login"        // TikTok 登录
	CategoryPostVideo   = "post_video"   // 发布视频
	CategoryLike        = "like"         // 点赞
	CategoryComment     = "comment"      // 评论
	CategoryFollow      = "follow"       // 关注
	CategoryBrowse      = "browse"       // 浏览视频（养号）
	CategoryInteraction = "interaction"  // 互动（养号）
	CategoryCustom      = "custom"       // 自定义分类
)
