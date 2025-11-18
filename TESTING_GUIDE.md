# TikTok 云控系统 - 数据库测试验证指南

## 📋 测试概览

本文档提供了 TikTok 云控系统数据模型的完整测试方案和验证步骤。

---

## 🗂️ 数据库表结构概览

### 已创建的 8 个核心表

| 序号 | 表名 | 说明 | 主要字段数 | 索引数 | 外键 |
|------|------|------|-----------|--------|------|
| 1 | `tk_devices` | 设备管理（核心计费） | 15 | 4 | 0 |
| 2 | `tk_scripts` | 脚本管理 | 14 | 3 | 0 |
| 3 | `tk_tasks` | 任务管理 | 18 | 6 | 2 |
| 4 | `tk_accounts` | TikTok账号资产 | 15 | 7 | 1 |
| 5 | `tk_task_logs` | 任务日志 | 9 | 4 | 2 |
| 6 | `tk_task_results` | 任务结果保存 | 8 | 5 | 2 |
| 7 | `tk_device_performance` | 设备性能监控 | 11 | 3 | 1 |
| 8 | `tk_alerts` | 报警记录 | 11 | 6 | 1 |

**总计**: 8 个表，101 个字段，38 个索引，9 个外键约束

---

## 🚀 快速开始 - 数据库测试

### 方法 1: 使用生成的 SQL 脚本（推荐）

#### 步骤 1: 创建测试数据库

```bash
# 连接 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE tiktok_cloud_control_test DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 使用数据库
USE tiktok_cloud_control_test;
```

#### 步骤 2: 导入表结构

```bash
# 从项目根目录执行
mysql -u root -p tiktok_cloud_control_test < server/scripts/tiktok_tables.sql
```

#### 步骤 3: 验证表创建成功

```sql
-- 查看所有 TikTok 相关表
SHOW TABLES LIKE 'tk_%';

-- 期望输出：8个表
-- +-------------------------------+
-- | Tables_in_test (tk_%)         |
-- +-------------------------------+
-- | tk_accounts                   |
-- | tk_alerts                     |
-- | tk_device_performance         |
-- | tk_devices                    |
-- | tk_scripts                    |
-- | tk_task_logs                  |
-- | tk_task_results               |
-- | tk_tasks                      |
-- +-------------------------------+
```

#### 步骤 4: 查看表结构详情

```sql
-- 查看设备表结构
DESC tk_devices;

-- 查看所有表的详细信息
SELECT
    TABLE_NAME AS '表名',
    TABLE_COMMENT AS '说明',
    ENGINE AS '引擎',
    TABLE_COLLATION AS '字符集'
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'tiktok_cloud_control_test'
  AND TABLE_NAME LIKE 'tk_%'
ORDER BY TABLE_NAME;
```

---

### 方法 2: 使用 Go 测试（需要网络环境支持）

#### 运行单元测试

```bash
# 进入 server 目录
cd server

# 运行所有模型测试
go test -v ./model/tiktok/...

# 运行特定测试
go test -v ./model/tiktok/ -run TestTableCreation
go test -v ./model/tiktok/ -run TestDeviceModel
go test -v ./model/tiktok/ -run TestAllModelsIntegration
```

#### 测试文件说明

**位置**: `server/model/tiktok/tiktok_model_test.go`

**包含的测试**:
1. ✅ `TestTableCreation` - 测试所有表能否正确创建
2. ✅ `TestDeviceModel` - 测试设备模型 CRUD 和业务方法
3. ✅ `TestScriptModel` - 测试脚本模型
4. ✅ `TestTaskModel` - 测试任务模型和重试逻辑
5. ✅ `TestAccountModel` - 测试账号模型和状态管理
6. ✅ `TestTaskLogModel` - 测试日志模型和截图功能
7. ✅ `TestDevicePerformanceModel` - 测试性能监控和阈值判断
8. ✅ `TestAlertModel` - 测试报警模型和已读状态
9. ✅ `TestAllModelsIntegration` - 综合集成测试

---

## 📊 表结构详细验证

### 1. tk_devices（设备管理表）⭐ 核心计费

#### 关键字段验证

```sql
-- 1.1 验证设备唯一性约束
SELECT COLUMN_NAME, IS_NULLABLE, COLUMN_KEY, COLUMN_TYPE
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = 'tiktok_cloud_control_test'
  AND TABLE_NAME = 'tk_devices'
  AND COLUMN_NAME = 'device_id';

-- 期望: UNIQUE 约束

-- 1.2 验证核心计费字段
SELECT COLUMN_NAME, IS_NULLABLE, COLUMN_DEFAULT
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = 'tiktok_cloud_control_test'
  AND TABLE_NAME = 'tk_devices'
  AND COLUMN_NAME = 'expire_time';

-- 期望: NOT NULL

-- 1.3 验证索引
SHOW INDEX FROM tk_devices;

-- 期望索引:
-- - PRIMARY (id)
-- - idx_user_id
-- - idx_device_id
-- - idx_expire_time
-- - idx_deleted_at
```

#### 插入测试数据

```sql
INSERT INTO tk_devices (
    user_id, device_id, device_name, device_token,
    model, android_version, bind_time, expire_time,
    status, last_active_time, group_name,
    created_at, updated_at
) VALUES (
    1, 'TEST-DEVICE-001', '测试设备1', 'token-abc123',
    'Xiaomi 13', '13.0', NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY),
    'online', NOW(), '注册组',
    NOW(), NOW()
);

-- 查询验证
SELECT device_id, device_name, status, expire_time FROM tk_devices;
```

### 2. tk_scripts（脚本管理表）

#### 验证 JSON 字段支持

```sql
-- 2.1 插入带 JSON 参数的脚本
INSERT INTO tk_scripts (
    name, category, version, description, file_path,
    params_schema, default_params, is_enabled, author,
    created_at, updated_at
) VALUES (
    'TikTok注册',
    'register',
    '1.0.0',
    '自动注册TikTok账号',
    '/scripts/register.js',
    JSON_OBJECT(
        'type', 'object',
        'properties', JSON_OBJECT(
            'phoneSource', JSON_OBJECT('type', 'string', 'enum', JSON_ARRAY('random', 'pool'))
        )
    ),
    JSON_OBJECT('phoneSource', 'random'),
    TRUE,
    '系统管理员',
    NOW(),
    NOW()
);

-- 2.2 查询并验证 JSON 数据
SELECT name, params_schema, default_params FROM tk_scripts;

-- 2.3 使用 JSON 函数查询
SELECT name,
       JSON_EXTRACT(default_params, '$.phoneSource') AS phone_source
FROM tk_scripts;
```

### 3. tk_tasks（任务管理表）

#### 验证外键约束

```sql
-- 3.1 查看外键约束
SELECT
    CONSTRAINT_NAME,
    COLUMN_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM information_schema.KEY_COLUMN_USAGE
WHERE TABLE_SCHEMA = 'tiktok_cloud_control_test'
  AND TABLE_NAME = 'tk_tasks'
  AND REFERENCED_TABLE_NAME IS NOT NULL;

-- 期望:
-- script_id -> tk_scripts(id)
-- device_id -> tk_devices(id)

-- 3.2 测试外键约束（应该失败）
INSERT INTO tk_tasks (
    user_id, script_id, device_id, task_name, task_type,
    priority, status, created_at, updated_at
) VALUES (
    1, 99999, 99999, '测试任务', 'once',
    'high', 'pending', NOW(), NOW()
);
-- 期望: ERROR 1452 - Cannot add or update a child row: foreign key constraint fails

-- 3.3 正确插入任务（假设已有 script_id=1, device_id=1）
INSERT INTO tk_tasks (
    user_id, script_id, device_id, task_name, task_type,
    priority, status, retry_count, max_retry,
    created_at, updated_at
) VALUES (
    1, 1, 1, '批量注册任务', 'batch',
    'high', 'pending', 0, 3,
    NOW(), NOW()
);
```

### 4. tk_accounts（账号资产表）⭐ 核心资产

#### 验证账号管理功能

```sql
-- 4.1 插入自动注册账号
INSERT INTO tk_accounts (
    user_id, device_id, source, username, phone, email,
    password_encrypted, registered_at, status, remark,
    created_at, updated_at
) VALUES (
    1, 1, 'auto_register', 'tiktok_user_001', '13800138001', 'user001@test.com',
    'AES_ENCRYPTED_PASSWORD_HERE', NOW(), 'normal', '自动注册测试账号',
    NOW(), NOW()
);

-- 4.2 统计不同来源的账号数量
SELECT source, COUNT(*) AS count, status
FROM tk_accounts
GROUP BY source, status;

-- 4.3 查询特定设备注册的账号
SELECT device_id, COUNT(*) AS registered_count
FROM tk_accounts
WHERE source = 'auto_register'
GROUP BY device_id;
```

### 5. tk_task_logs（任务日志表）

#### 验证日志记录和级联删除

```sql
-- 5.1 插入日志
INSERT INTO tk_task_logs (
    task_id, device_id, log_level, log_content,
    screenshot_url, created_at, updated_at
) VALUES
(1, 1, 'info', '开始执行TikTok注册脚本', NULL, NOW(), NOW()),
(1, 1, 'info', '正在填写注册信息', NULL, NOW(), NOW()),
(1, 1, 'error', '脚本执行失败: 网络超时', '/screenshots/error_001.png', NOW(), NOW());

-- 5.2 按日志级别统计
SELECT log_level, COUNT(*) AS count
FROM tk_task_logs
GROUP BY log_level;

-- 5.3 测试级联删除（删除任务时，日志应该自动删除）
DELETE FROM tk_tasks WHERE id = 1;

-- 验证日志是否被级联删除
SELECT COUNT(*) FROM tk_task_logs WHERE task_id = 1;
-- 期望: 0
```

### 6. tk_task_results（任务结果表）

#### 验证结果数据保存

```sql
-- 6.1 插入注册成功结果
INSERT INTO tk_task_results (
    task_id, script_id, result_type, result_data,
    created_at, updated_at
) VALUES (
    1, 1, 'register',
    JSON_OBJECT(
        'account', 'tiktok_user_001',
        'password', 'encrypted_pass_123',
        'phone', '13800138001',
        'email', 'user001@test.com',
        'registeredAt', NOW()
    ),
    NOW(), NOW()
);

-- 6.2 查询并解析结果
SELECT
    result_type,
    JSON_EXTRACT(result_data, '$.account') AS account,
    JSON_EXTRACT(result_data, '$.phone') AS phone
FROM tk_task_results;

-- 6.3 统计不同类型的结果
SELECT result_type, COUNT(*) AS count
FROM tk_task_results
GROUP BY result_type;
```

### 7. tk_device_performance（设备性能监控表）

#### 验证性能数据和阈值判断

```sql
-- 7.1 插入性能数据
INSERT INTO tk_device_performance (
    device_id, cpu_usage, memory_usage, battery_level,
    network_type, storage_free, reported_at,
    created_at, updated_at
) VALUES
(1, 45.5, 60.2, 80, 'WiFi', 15360, NOW(), NOW(), NOW()),
(1, 85.2, 92.1, 15, '4G', 1024, NOW(), NOW(), NOW());

-- 7.2 查询性能异常的设备
SELECT
    device_id,
    cpu_usage,
    memory_usage,
    battery_level,
    reported_at
FROM tk_device_performance
WHERE cpu_usage > 80
   OR memory_usage > 80
   OR battery_level < 20
ORDER BY reported_at DESC;

-- 7.3 计算平均性能指标
SELECT
    device_id,
    AVG(cpu_usage) AS avg_cpu,
    AVG(memory_usage) AS avg_memory,
    AVG(battery_level) AS avg_battery
FROM tk_device_performance
WHERE reported_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
GROUP BY device_id;
```

### 8. tk_alerts（报警记录表）

#### 验证报警管理

```sql
-- 8.1 插入不同级别的报警
INSERT INTO tk_alerts (
    user_id, device_id, alert_type, alert_level,
    alert_message, is_read, created_at, updated_at
) VALUES
(1, 1, 'device_offline', 'critical', '设备已离线超过10分钟', FALSE, NOW(), NOW()),
(1, 1, 'cpu_high', 'warning', 'CPU使用率超过80%', FALSE, NOW(), NOW()),
(1, NULL, 'device_expiring', 'info', '您有1台设备将在3天后到期', FALSE, NOW(), NOW());

-- 8.2 查询未读的严重报警
SELECT alert_type, alert_message, created_at
FROM tk_alerts
WHERE is_read = FALSE
  AND alert_level = 'critical'
ORDER BY created_at DESC;

-- 8.3 标记报警为已读
UPDATE tk_alerts
SET is_read = TRUE
WHERE user_id = 1;

-- 8.4 统计报警数量
SELECT
    alert_level,
    COUNT(*) AS count,
    SUM(CASE WHEN is_read = FALSE THEN 1 ELSE 0 END) AS unread_count
FROM tk_alerts
GROUP BY alert_level;
```

---

## 🔍 性能测试查询

### 索引效率验证

```sql
-- 使用 EXPLAIN 分析查询性能

-- 1. 查询特定设备
EXPLAIN SELECT * FROM tk_devices WHERE device_id = 'TEST-DEVICE-001';
-- 期望: Using index

-- 2. 查询即将过期的设备
EXPLAIN SELECT * FROM tk_devices
WHERE expire_time BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 3 DAY);
-- 期望: Using index condition

-- 3. 查询用户的待执行任务
EXPLAIN SELECT * FROM tk_tasks
WHERE user_id = 1 AND status = 'pending'
ORDER BY priority DESC, created_at ASC;
-- 期望: Using where; Using index

-- 4. 查询任务的所有日志
EXPLAIN SELECT * FROM tk_task_logs WHERE task_id = 1;
-- 期望: Using index
```

### 复杂联表查询测试

```sql
-- 1. 查询任务详情（包含脚本和设备信息）
SELECT
    t.id AS task_id,
    t.task_name,
    s.name AS script_name,
    s.category AS script_category,
    d.device_name,
    d.status AS device_status,
    t.status AS task_status,
    t.created_at
FROM tk_tasks t
LEFT JOIN tk_scripts s ON t.script_id = s.id
LEFT JOIN tk_devices d ON t.device_id = d.id
WHERE t.user_id = 1
ORDER BY t.created_at DESC
LIMIT 10;

-- 2. 统计每个脚本的执行情况
SELECT
    s.name AS script_name,
    COUNT(t.id) AS total_tasks,
    SUM(CASE WHEN t.status = 'success' THEN 1 ELSE 0 END) AS success_count,
    SUM(CASE WHEN t.status = 'failed' THEN 1 ELSE 0 END) AS failed_count,
    ROUND(SUM(CASE WHEN t.status = 'success' THEN 1 ELSE 0 END) * 100.0 / COUNT(t.id), 2) AS success_rate
FROM tk_scripts s
LEFT JOIN tk_tasks t ON s.id = t.script_id
GROUP BY s.id, s.name;

-- 3. 查询设备的最新性能数据
SELECT
    d.device_name,
    p.cpu_usage,
    p.memory_usage,
    p.battery_level,
    p.reported_at
FROM tk_devices d
LEFT JOIN (
    SELECT device_id, cpu_usage, memory_usage, battery_level, reported_at
    FROM tk_device_performance
    WHERE (device_id, reported_at) IN (
        SELECT device_id, MAX(reported_at)
        FROM tk_device_performance
        GROUP BY device_id
    )
) p ON d.id = p.device_id
WHERE d.user_id = 1;
```

---

## ✅ 验证检查清单

### 表结构验证

- [ ] 所有 8 个表成功创建
- [ ] 所有字段类型正确（VARCHAR, INT, DATETIME, JSON, ENUM等）
- [ ] 所有必填字段（NOT NULL）配置正确
- [ ] 所有默认值设置正确
- [ ] 软删除字段（deleted_at）存在且有索引

### 索引验证

- [ ] 主键索引（id）存在
- [ ] 用户 ID 索引（user_id）存在
- [ ] 外键字段索引存在
- [ ] 高频查询字段索引存在
- [ ] 唯一性约束（device_id）生效

### 约束验证

- [ ] 外键约束正确配置
- [ ] 级联删除（CASCADE）测试通过
- [ ] 级联设置NULL（SET NULL）测试通过
- [ ] 限制删除（RESTRICT）测试通过
- [ ] 唯一性约束测试通过

### 功能验证

- [ ] INSERT 操作成功
- [ ] SELECT 查询正确
- [ ] UPDATE 更新正常
- [ ] DELETE 删除（软删除）正常
- [ ] JSON 字段读写正常
- [ ] ENUM 枚举值限制生效

### 性能验证

- [ ] 单表查询使用索引
- [ ] 联表查询执行计划合理
- [ ] 无全表扫描（除必要情况）
- [ ] 复杂查询响应时间 < 100ms
- [ ] 插入操作响应时间 < 10ms

---

## 📈 测试数据生成脚本

### 批量生成测试数据

```sql
-- 生成多个测试设备
DELIMITER $$
CREATE PROCEDURE generate_test_devices(IN count INT)
BEGIN
    DECLARE i INT DEFAULT 1;
    WHILE i <= count DO
        INSERT INTO tk_devices (
            user_id, device_id, device_name, device_token,
            model, android_version, bind_time, expire_time,
            status, last_active_time, created_at, updated_at
        ) VALUES (
            1,
            CONCAT('DEVICE-', LPAD(i, 5, '0')),
            CONCAT('测试设备', i),
            MD5(CONCAT('token-', i)),
            'Xiaomi 13',
            '13.0',
            NOW(),
            DATE_ADD(NOW(), INTERVAL 30 DAY),
            'online',
            NOW(),
            NOW(),
            NOW()
        );
        SET i = i + 1;
    END WHILE;
END$$
DELIMITER ;

-- 生成 100 个测试设备
CALL generate_test_devices(100);

-- 验证
SELECT COUNT(*) FROM tk_devices;
```

---

## 🐛 常见问题排查

### 1. 外键约束创建失败

**问题**: `ERROR 1215: Cannot add foreign key constraint`

**解决方案**:
```sql
-- 检查父表是否存在
SHOW TABLES LIKE 'tk_scripts';

-- 检查父表字段类型
DESC tk_scripts;

-- 确保子表字段类型与父表一致
SELECT COLUMN_NAME, COLUMN_TYPE
FROM information_schema.COLUMNS
WHERE TABLE_NAME IN ('tk_tasks', 'tk_scripts')
  AND COLUMN_NAME IN ('id', 'script_id');
```

### 2. JSON 字段无法插入

**问题**: `Invalid JSON text`

**解决方案**:
```sql
-- 使用 JSON_VALID 验证
SELECT JSON_VALID('{"key": "value"}');

-- 使用 JSON_OBJECT 构建
SELECT JSON_OBJECT('key', 'value');
```

### 3. ENUM 值错误

**问题**: `Data truncated for column 'status'`

**解决方案**:
```sql
-- 查看 ENUM 允许的值
SHOW COLUMNS FROM tk_devices LIKE 'status';

-- 只能使用: 'online', 'offline', 'expired'
```

---

## 📚 参考文档

1. **GORM 文档**: https://gorm.io/zh_CN/docs/
2. **MySQL 8.0 文档**: https://dev.mysql.com/doc/refman/8.0/en/
3. **gin-vue-admin**: https://www.gin-vue-admin.com/
4. **项目需求文档**: 见项目根目录

---

## ✨ 测试总结模板

### 测试报告

**测试日期**: ___________

**测试人员**: ___________

**测试环境**:
- MySQL 版本: ___________
- 数据库名称: ___________

**测试结果**:

| 测试项 | 状态 | 备注 |
|--------|------|------|
| 表结构创建 | ✅/❌ | |
| 索引配置 | ✅/❌ | |
| 外键约束 | ✅/❌ | |
| CRUD 操作 | ✅/❌ | |
| JSON 字段 | ✅/❌ | |
| 性能测试 | ✅/❌ | |

**发现的问题**:
1. ___________
2. ___________

**优化建议**:
1. ___________
2. ___________

---

**最后更新**: 2025-11-18
