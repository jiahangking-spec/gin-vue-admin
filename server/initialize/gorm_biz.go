package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/tiktok"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		// TikTok 云控系统核心表
		&tiktok.TkDevice{},             // 设备管理表
		&tiktok.TkScript{},             // 脚本管理表
		&tiktok.TkTask{},               // 任务管理表
		&tiktok.TkAccount{},            // TikTok 账号资产表
		&tiktok.TkTaskLog{},            // 任务日志表
		&tiktok.TkTaskResult{},         // 任务结果表
		&tiktok.TkDevicePerformance{},  // 设备性能监控表
		&tiktok.TkAlert{},              // 报警记录表
	)
	if err != nil {
		return err
	}
	return nil
}
