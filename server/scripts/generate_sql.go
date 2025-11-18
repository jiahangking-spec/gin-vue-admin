package main

import (
	"fmt"
	"os"
)

// 生成 TikTok 云控系统的建表 SQL
func main() {
	sql := `-- TikTok 云控管理系统 - MySQL 数据库建表脚本
-- 生成时间: 2025-11-18
-- 基于 GORM 模型自动生成

-- ============================================
-- 1. 设备管理表 - 核心计费模块
-- ============================================
CREATE TABLE IF NOT EXISTS tk_devices (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    user_id BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
    device_id VARCHAR(100) NOT NULL UNIQUE COMMENT '设备唯一标识(IMEI/设备ID)',
    device_name VARCHAR(100) DEFAULT NULL COMMENT '设备备注名称',
    device_token VARCHAR(255) NOT NULL COMMENT '设备授权Token',
    model VARCHAR(100) DEFAULT NULL COMMENT '设备型号',
    android_version VARCHAR(50) DEFAULT NULL COMMENT 'Android版本',
    bind_time DATETIME(3) DEFAULT NULL COMMENT '绑定时间',
    expire_time DATETIME(3) NOT NULL COMMENT '到期时间-核心计费字段',
    status ENUM('online','offline','expired') DEFAULT 'offline' COMMENT '设备状态',
    last_active_time DATETIME(3) DEFAULT NULL COMMENT '最后活跃时间',
    group_name VARCHAR(100) DEFAULT NULL COMMENT '设备分组名称',

    INDEX idx_user_id (user_id),
    INDEX idx_device_id (device_id),
    INDEX idx_expire_time (expire_time),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备管理表-核心计费模块';

-- ============================================
-- 2. 脚本管理表
-- ============================================
CREATE TABLE IF NOT EXISTS tk_scripts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    name VARCHAR(100) NOT NULL COMMENT '脚本名称',
    category VARCHAR(50) NOT NULL COMMENT '脚本分类',
    version VARCHAR(20) DEFAULT '1.0.0' COMMENT '脚本版本',
    description TEXT COMMENT '脚本描述',
    file_path VARCHAR(255) NOT NULL COMMENT '脚本文件路径',
    params_schema JSON COMMENT '参数配置Schema',
    default_params JSON COMMENT '默认参数值',
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    author VARCHAR(50) DEFAULT NULL COMMENT '脚本作者',
    tags VARCHAR(255) DEFAULT NULL COMMENT '脚本标签(逗号分隔)',

    INDEX idx_category (category),
    INDEX idx_name (name),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='脚本管理表';

-- ============================================
-- 3. 任务管理表
-- ============================================
CREATE TABLE IF NOT EXISTS tk_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    user_id BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
    script_id BIGINT UNSIGNED NOT NULL COMMENT '脚本ID',
    device_id BIGINT UNSIGNED NOT NULL COMMENT '执行设备ID',
    task_name VARCHAR(100) DEFAULT NULL COMMENT '任务名称',
    task_type ENUM('once','scheduled','loop','batch') NOT NULL COMMENT '任务类型',
    params JSON COMMENT '任务参数',
    priority ENUM('high','medium','low') DEFAULT 'medium' COMMENT '优先级',
    status ENUM('pending','running','success','failed','cancelled') DEFAULT 'pending' COMMENT '任务状态',
    scheduled_time DATETIME(3) DEFAULT NULL COMMENT '定时执行时间',
    start_time DATETIME(3) DEFAULT NULL COMMENT '开始执行时间',
    end_time DATETIME(3) DEFAULT NULL COMMENT '结束时间',
    retry_count INT DEFAULT 0 COMMENT '已重试次数',
    max_retry INT DEFAULT 3 COMMENT '最大重试次数',
    error_message TEXT COMMENT '错误信息',
    duration INT DEFAULT NULL COMMENT '执行耗时(秒)',

    INDEX idx_user_id (user_id),
    INDEX idx_script_id (script_id),
    INDEX idx_device_id (device_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),

    FOREIGN KEY (script_id) REFERENCES tk_scripts(id) ON DELETE RESTRICT,
    FOREIGN KEY (device_id) REFERENCES tk_devices(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务管理表';

-- ============================================
-- 4. TikTok账号资产表 - 核心资产
-- ============================================
CREATE TABLE IF NOT EXISTS tk_accounts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    user_id BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
    device_id BIGINT UNSIGNED DEFAULT NULL COMMENT '注册设备ID',
    source ENUM('auto_register','manual_import') NOT NULL COMMENT '账号来源',
    username VARCHAR(100) NOT NULL COMMENT 'TikTok用户名',
    phone VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    email VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
    password_encrypted VARCHAR(255) NOT NULL COMMENT 'AES加密的密码',
    registered_at DATETIME(3) DEFAULT NULL COMMENT '注册时间',
    last_login_at DATETIME(3) DEFAULT NULL COMMENT '最后登录时间',
    status ENUM('normal','banned','pending') DEFAULT 'normal' COMMENT '账号状态',
    remark TEXT COMMENT '备注',

    INDEX idx_user_id (user_id),
    INDEX idx_device_id (device_id),
    INDEX idx_username (username),
    INDEX idx_phone (phone),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),

    FOREIGN KEY (device_id) REFERENCES tk_devices(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='TikTok账号资产表';

-- ============================================
-- 5. 任务日志表
-- ============================================
CREATE TABLE IF NOT EXISTS tk_task_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    task_id BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    device_id BIGINT UNSIGNED NOT NULL COMMENT '设备ID',
    log_level ENUM('info','warn','error','debug') DEFAULT 'info' COMMENT '日志级别',
    log_content TEXT NOT NULL COMMENT '日志内容',
    screenshot_url VARCHAR(255) DEFAULT NULL COMMENT '截图URL',

    INDEX idx_task_id (task_id),
    INDEX idx_device_id (device_id),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),

    FOREIGN KEY (task_id) REFERENCES tk_tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES tk_devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务日志表';

-- ============================================
-- 6. 任务结果表 - 核心结果保存
-- ============================================
CREATE TABLE IF NOT EXISTS tk_task_results (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    task_id BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    script_id BIGINT UNSIGNED NOT NULL COMMENT '脚本ID',
    result_type VARCHAR(50) NOT NULL COMMENT '结果类型',
    result_data JSON NOT NULL COMMENT '结果数据JSON格式',

    INDEX idx_task_id (task_id),
    INDEX idx_script_id (script_id),
    INDEX idx_result_type (result_type),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),

    FOREIGN KEY (task_id) REFERENCES tk_tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (script_id) REFERENCES tk_scripts(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务结果表';

-- ============================================
-- 7. 设备性能监控表
-- ============================================
CREATE TABLE IF NOT EXISTS tk_device_performance (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    device_id BIGINT UNSIGNED NOT NULL COMMENT '设备ID',
    cpu_usage DECIMAL(5,2) DEFAULT NULL COMMENT 'CPU使用率%',
    memory_usage DECIMAL(5,2) DEFAULT NULL COMMENT '内存使用率%',
    battery_level INT DEFAULT NULL COMMENT '电量%',
    network_type VARCHAR(20) DEFAULT NULL COMMENT '网络类型',
    storage_free BIGINT DEFAULT NULL COMMENT '可用存储空间MB',
    reported_at DATETIME(3) DEFAULT NULL COMMENT '上报时间',

    INDEX idx_device_id (device_id),
    INDEX idx_reported_at (reported_at),
    INDEX idx_deleted_at (deleted_at),

    FOREIGN KEY (device_id) REFERENCES tk_devices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备性能监控表';

-- ============================================
-- 8. 报警记录表
-- ============================================
CREATE TABLE IF NOT EXISTS tk_alerts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at DATETIME(3) DEFAULT NULL COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT NULL COMMENT '更新时间',
    deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

    user_id BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
    device_id BIGINT UNSIGNED DEFAULT NULL COMMENT '设备ID',
    alert_type VARCHAR(50) NOT NULL COMMENT '报警类型',
    alert_level ENUM('info','warning','error','critical') DEFAULT 'warning' COMMENT '报警级别',
    alert_message TEXT NOT NULL COMMENT '报警消息',
    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已读',

    INDEX idx_user_id (user_id),
    INDEX idx_device_id (device_id),
    INDEX idx_alert_level (alert_level),
    INDEX idx_is_read (is_read),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),

    FOREIGN KEY (device_id) REFERENCES tk_devices(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报警记录表';

-- ============================================
-- 初始化示例数据（可选）
-- ============================================

-- 插入示例脚本
INSERT INTO tk_scripts (name, category, version, description, file_path, is_enabled, author, created_at, updated_at) VALUES
('TikTok自动注册', 'register', '1.0.0', '自动注册TikTok账号', '/scripts/tiktok/register.js', TRUE, '系统管理员', NOW(), NOW()),
('TikTok自动登录', 'login', '1.0.0', '自动登录TikTok账号', '/scripts/tiktok/login.js', TRUE, '系统管理员', NOW(), NOW()),
('TikTok发布视频', 'post_video', '1.0.0', '自动发布TikTok视频', '/scripts/tiktok/post_video.js', TRUE, '系统管理员', NOW(), NOW()),
('TikTok点赞', 'like', '1.0.0', '自动点赞TikTok视频', '/scripts/tiktok/like.js', TRUE, '系统管理员', NOW(), NOW()),
('TikTok关注用户', 'follow', '1.0.0', '自动关注TikTok用户', '/scripts/tiktok/follow.js', TRUE, '系统管理员', NOW(), NOW()),
('TikTok浏览养号', 'browse', '1.0.0', 'TikTok浏览视频养号', '/scripts/tiktok/browse.js', TRUE, '系统管理员', NOW(), NOW());

-- ============================================
-- 数据库表结构验证查询
-- ============================================

-- 查看所有TikTok相关表
SELECT
    TABLE_NAME AS '表名',
    TABLE_COMMENT AS '说明',
    TABLE_ROWS AS '行数',
    ROUND(DATA_LENGTH/1024/1024, 2) AS '数据大小(MB)',
    ROUND(INDEX_LENGTH/1024/1024, 2) AS '索引大小(MB)',
    ENGINE AS '引擎',
    TABLE_COLLATION AS '字符集'
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME LIKE 'tk_%'
ORDER BY TABLE_NAME;

-- 查看表字段详情（以设备表为例）
SELECT
    COLUMN_NAME AS '字段名',
    COLUMN_TYPE AS '类型',
    IS_NULLABLE AS '可空',
    COLUMN_DEFAULT AS '默认值',
    COLUMN_KEY AS '键',
    COLUMN_COMMENT AS '注释'
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'tk_devices'
ORDER BY ORDINAL_POSITION;

-- 查看外键约束
SELECT
    CONSTRAINT_NAME AS '约束名',
    TABLE_NAME AS '表名',
    COLUMN_NAME AS '字段名',
    REFERENCED_TABLE_NAME AS '关联表',
    REFERENCED_COLUMN_NAME AS '关联字段'
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = DATABASE()
  AND REFERENCED_TABLE_NAME IS NOT NULL
  AND TABLE_NAME LIKE 'tk_%'
ORDER BY TABLE_NAME;

-- ============================================
-- 性能优化建议
-- ============================================

-- 1. 为高频查询字段添加复合索引
-- ALTER TABLE tk_tasks ADD INDEX idx_user_device_status (user_id, device_id, status);
-- ALTER TABLE tk_accounts ADD INDEX idx_user_status (user_id, status);

-- 2. 为日志表按月分区（可选，适用于海量日志）
-- ALTER TABLE tk_task_logs PARTITION BY RANGE (YEAR(created_at)*100 + MONTH(created_at)) (
--     PARTITION p202501 VALUES LESS THAN (202502),
--     PARTITION p202502 VALUES LESS THAN (202503)
-- );

-- 3. 定期清理过期数据（建议定时任务）
-- DELETE FROM tk_task_logs WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
-- DELETE FROM tk_device_performance WHERE reported_at < DATE_SUB(NOW(), INTERVAL 7 DAY);
`

	fmt.Println(sql)

	// 保存到文件
	filename := "tiktok_tables.sql"
	err := os.WriteFile(filename, []byte(sql), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n\n✅ SQL文件已生成: %s\n", filename)
	fmt.Println("========================================")
	fmt.Println("  使用说明")
	fmt.Println("========================================")
	fmt.Println("1. 创建数据库:")
	fmt.Println("   CREATE DATABASE tiktok_cloud_control DEFAULT CHARACTER SET utf8mb4;")
	fmt.Println("")
	fmt.Println("2. 导入表结构:")
	fmt.Println("   mysql -u root -p tiktok_cloud_control < tiktok_tables.sql")
	fmt.Println("")
	fmt.Println("3. 验证表结构:")
	fmt.Println("   SHOW TABLES LIKE 'tk_%';")
	fmt.Println("========================================")
}
