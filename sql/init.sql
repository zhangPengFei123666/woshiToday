-- ============================================
-- 分布式任务跑批系统 - 数据库初始化脚本
-- ============================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `distributed_scheduler` 
DEFAULT CHARACTER SET utf8mb4 
DEFAULT COLLATE utf8mb4_unicode_ci;

USE `distributed_scheduler`;

-- ============================================
-- 用户相关表
-- ============================================

-- 用户表
CREATE TABLE IF NOT EXISTS `sys_user` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(64) NOT NULL COMMENT '用户名',
    `password` VARCHAR(128) NOT NULL COMMENT '密码(加密)',
    `nickname` VARCHAR(64) DEFAULT '' COMMENT '昵称',
    `email` VARCHAR(128) DEFAULT '' COMMENT '邮箱',
    `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
    `avatar` VARCHAR(256) DEFAULT '' COMMENT '头像URL',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0-禁用 1-启用',
    `last_login_time` DATETIME DEFAULT NULL COMMENT '最后登录时间',·
    `last_login_ip` VARCHAR(64) DEFAULT '' COMMENT '最后登录IP',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    UNIQUE KEY `uk_username` (`username`),
    INDEX `idx_status` (`status`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS `sys_role` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '角色ID',
    `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
    `code` VARCHAR(64) NOT NULL COMMENT '角色编码',
    `description` VARCHAR(256) DEFAULT '' COMMENT '角色描述',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0-禁用 1-启用',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    UNIQUE KEY `uk_code` (`code`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统角色表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS `sys_user_role` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    INDEX `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ============================================
-- 任务组相关表
-- ============================================

-- 任务组表
CREATE TABLE IF NOT EXISTS `task_group` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '任务组ID',
    `name` VARCHAR(128) NOT NULL COMMENT '任务组名称',
    `description` VARCHAR(512) DEFAULT '' COMMENT '描述',
    `app_name` VARCHAR(64) NOT NULL COMMENT '应用名称(用于执行器注册)',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0-禁用 1-启用',
    `created_by` BIGINT UNSIGNED DEFAULT 0 COMMENT '创建人ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    UNIQUE KEY `uk_app_name` (`app_name`),
    INDEX `idx_status` (`status`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务组表';

-- ============================================
-- 任务相关表
-- ============================================

-- 任务定义表
CREATE TABLE IF NOT EXISTS `task` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '任务ID',
    `group_id` BIGINT UNSIGNED NOT NULL COMMENT '任务组ID',
    `name` VARCHAR(128) NOT NULL COMMENT '任务名称',
    `description` VARCHAR(512) DEFAULT '' COMMENT '任务描述',
    `cron` VARCHAR(64) NOT NULL COMMENT 'Cron表达式',
    `executor_type` VARCHAR(32) NOT NULL DEFAULT 'HTTP' COMMENT '执行器类型 HTTP/GRPC/SCRIPT',
    `executor_handler` VARCHAR(256) NOT NULL COMMENT '执行器Handler',
    `executor_param` TEXT COMMENT '执行参数(JSON格式)',
    `route_strategy` VARCHAR(32) DEFAULT 'ROUND_ROBIN' COMMENT '路由策略 ROUND_ROBIN/RANDOM/CONSISTENT_HASH/LEAST_FREQUENTLY_USED/LEAST_RECENTLY_USED/FAILOVER/SHARDING_BROADCAST',
    `block_strategy` VARCHAR(32) DEFAULT 'SERIAL_EXECUTION' COMMENT '阻塞策略 SERIAL_EXECUTION/DISCARD_LATER/COVER_EARLY',
    `shard_num` INT UNSIGNED DEFAULT 1 COMMENT '分片数量',
    `retry_count` INT UNSIGNED DEFAULT 0 COMMENT '失败重试次数',
    `retry_interval` INT UNSIGNED DEFAULT 0 COMMENT '重试间隔(秒)',
    `timeout` INT UNSIGNED DEFAULT 0 COMMENT '任务超时时间(秒) 0-无限制',
    `alarm_email` VARCHAR(512) DEFAULT '' COMMENT '告警邮箱(多个用逗号分隔)',
    `priority` INT DEFAULT 0 COMMENT '优先级 数值越大优先级越高',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0-禁用 1-启用',
    `version` INT UNSIGNED DEFAULT 0 COMMENT '版本号(乐观锁)',
    `next_trigger_time` DATETIME DEFAULT NULL COMMENT '下次触发时间',
    `last_trigger_time` DATETIME DEFAULT NULL COMMENT '上次触发时间',
    `created_by` BIGINT UNSIGNED DEFAULT 0 COMMENT '创建人ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    INDEX `idx_group_id` (`group_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_next_trigger_time` (`next_trigger_time`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务定义表';

-- 任务依赖关系表(DAG)
CREATE TABLE IF NOT EXISTS `task_dependency` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    `depend_task_id` BIGINT UNSIGNED NOT NULL COMMENT '依赖的任务ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uk_task_depend` (`task_id`, `depend_task_id`),
    INDEX `idx_depend_task_id` (`depend_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务依赖关系表';

-- ============================================
-- 任务执行相关表
-- ============================================

-- 任务实例表(执行记录)
CREATE TABLE IF NOT EXISTS `task_instance` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '实例ID',
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    `group_id` BIGINT UNSIGNED NOT NULL COMMENT '任务组ID',
    `executor_id` VARCHAR(128) DEFAULT '' COMMENT '执行节点ID',
    `executor_address` VARCHAR(256) DEFAULT '' COMMENT '执行节点地址',
    `executor_handler` VARCHAR(256) DEFAULT '' COMMENT '执行器Handler',
    `executor_param` TEXT COMMENT '执行参数',
    `shard_index` INT UNSIGNED DEFAULT 0 COMMENT '分片索引',
    `shard_total` INT UNSIGNED DEFAULT 1 COMMENT '分片总数',
    `trigger_type` VARCHAR(32) DEFAULT 'CRON' COMMENT '触发类型 CRON/MANUAL/PARENT/API/RETRY',
    `trigger_time` DATETIME NOT NULL COMMENT '触发时间',
    `schedule_time` DATETIME DEFAULT NULL COMMENT '调度时间',
    `start_time` DATETIME DEFAULT NULL COMMENT '开始执行时间',
    `end_time` DATETIME DEFAULT NULL COMMENT '结束时间',
    `status` TINYINT DEFAULT 0 COMMENT '状态 0-待调度 1-调度中 2-执行中 3-执行成功 4-执行失败 5-已取消',
    `result_code` INT DEFAULT 0 COMMENT '结果码 0-成功 其他-失败',
    `result_msg` TEXT COMMENT '执行结果消息',
    `retry_count` INT UNSIGNED DEFAULT 0 COMMENT '已重试次数',
    `alarm_status` TINYINT DEFAULT 0 COMMENT '告警状态 0-默认 1-已告警',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_task_id` (`task_id`),
    INDEX `idx_group_id` (`group_id`),
    INDEX `idx_trigger_time` (`trigger_time`),
    INDEX `idx_status` (`status`),
    INDEX `idx_executor_id` (`executor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务实例表';

-- 任务执行日志表
CREATE TABLE IF NOT EXISTS `task_log` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    `instance_id` BIGINT UNSIGNED NOT NULL COMMENT '任务实例ID',
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    `log_time` DATETIME(3) NOT NULL COMMENT '日志时间(毫秒精度)',
    `log_level` VARCHAR(16) DEFAULT 'INFO' COMMENT '日志级别 DEBUG/INFO/WARN/ERROR',
    `log_content` TEXT COMMENT '日志内容',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_instance_id` (`instance_id`),
    INDEX `idx_task_id` (`task_id`),
    INDEX `idx_log_time` (`log_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务执行日志表';

-- ============================================
-- 执行节点相关表
-- ============================================

-- 执行器节点表
CREATE TABLE IF NOT EXISTS `executor_node` (
    `id` VARCHAR(128) PRIMARY KEY COMMENT '节点ID(UUID)',
    `group_id` BIGINT UNSIGNED NOT NULL COMMENT '任务组ID',
    `app_name` VARCHAR(64) NOT NULL COMMENT '应用名称',
    `host` VARCHAR(128) NOT NULL COMMENT '节点IP',
    `port` INT UNSIGNED NOT NULL COMMENT '节点端口',
    `weight` INT UNSIGNED DEFAULT 100 COMMENT '权重(用于负载均衡)',
    `max_concurrent` INT UNSIGNED DEFAULT 100 COMMENT '最大并发任务数',
    `current_load` INT UNSIGNED DEFAULT 0 COMMENT '当前负载(执行中任务数)',
    `cpu_usage` DECIMAL(5,2) DEFAULT 0 COMMENT 'CPU使用率',
    `memory_usage` DECIMAL(5,2) DEFAULT 0 COMMENT '内存使用率',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0-离线 1-在线',
    `last_heartbeat` DATETIME DEFAULT NULL COMMENT '最后心跳时间',
    `registered_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_group_id` (`group_id`),
    INDEX `idx_app_name` (`app_name`),
    INDEX `idx_status` (`status`),
    INDEX `idx_last_heartbeat` (`last_heartbeat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='执行器节点表';

-- ============================================
-- 告警相关表
-- ============================================

-- 告警规则表
CREATE TABLE IF NOT EXISTS `alarm_rule` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '规则ID',
    `name` VARCHAR(128) NOT NULL COMMENT '规则名称',
    `rule_type` VARCHAR(32) NOT NULL COMMENT '规则类型 TASK_FAIL/TASK_TIMEOUT/EXECUTOR_OFFLINE',
    `group_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '任务组ID(0表示全局)',
    `task_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '任务ID(0表示所有任务)',
    `threshold` INT UNSIGNED DEFAULT 1 COMMENT '阈值',
    `alarm_level` VARCHAR(16) DEFAULT 'WARNING' COMMENT '告警级别 INFO/WARNING/ERROR/CRITICAL',
    `notify_type` VARCHAR(64) DEFAULT 'EMAIL' COMMENT '通知方式 EMAIL/SMS/WEBHOOK(多个用逗号分隔)',
    `notify_target` TEXT COMMENT '通知目标(JSON格式)',
    `status` TINYINT DEFAULT 1 COMMENT '状态 0-禁用 1-启用',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
    INDEX `idx_rule_type` (`rule_type`),
    INDEX `idx_group_id` (`group_id`),
    INDEX `idx_task_id` (`task_id`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警规则表';

-- 告警记录表
CREATE TABLE IF NOT EXISTS `alarm_record` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    `rule_id` BIGINT UNSIGNED NOT NULL COMMENT '规则ID',
    `task_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '任务ID',
    `instance_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '实例ID',
    `alarm_type` VARCHAR(32) NOT NULL COMMENT '告警类型',
    `alarm_level` VARCHAR(16) NOT NULL COMMENT '告警级别',
    `alarm_title` VARCHAR(256) NOT NULL COMMENT '告警标题',
    `alarm_content` TEXT COMMENT '告警内容',
    `notify_status` TINYINT DEFAULT 0 COMMENT '通知状态 0-待发送 1-已发送 2-发送失败',
    `notify_time` DATETIME DEFAULT NULL COMMENT '通知时间',
    `notify_result` TEXT COMMENT '通知结果',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_rule_id` (`rule_id`),
    INDEX `idx_task_id` (`task_id`),
    INDEX `idx_alarm_type` (`alarm_type`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警记录表';

-- ============================================
-- 系统配置表
-- ============================================

-- 系统配置表
CREATE TABLE IF NOT EXISTS `sys_config` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    `config_key` VARCHAR(128) NOT NULL COMMENT '配置键',
    `config_value` TEXT COMMENT '配置值',
    `config_type` VARCHAR(32) DEFAULT 'STRING' COMMENT '配置类型 STRING/NUMBER/BOOLEAN/JSON',
    `description` VARCHAR(256) DEFAULT '' COMMENT '描述',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `operation_log` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    `user_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '操作用户ID',
    `username` VARCHAR(64) DEFAULT '' COMMENT '操作用户名',
    `module` VARCHAR(64) DEFAULT '' COMMENT '操作模块',
    `action` VARCHAR(64) DEFAULT '' COMMENT '操作类型',
    `target_type` VARCHAR(64) DEFAULT '' COMMENT '操作对象类型',
    `target_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '操作对象ID',
    `request_method` VARCHAR(16) DEFAULT '' COMMENT '请求方法',
    `request_url` VARCHAR(512) DEFAULT '' COMMENT '请求URL',
    `request_param` TEXT COMMENT '请求参数',
    `request_ip` VARCHAR(64) DEFAULT '' COMMENT '请求IP',
    `user_agent` VARCHAR(512) DEFAULT '' COMMENT 'User-Agent',
    `result_code` INT DEFAULT 0 COMMENT '结果码',
    `result_msg` VARCHAR(512) DEFAULT '' COMMENT '结果消息',
    `duration` INT UNSIGNED DEFAULT 0 COMMENT '耗时(毫秒)',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_module` (`module`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- ============================================
-- 初始化数据
-- ============================================

-- 插入默认管理员用户 (密码: admin123, 使用bcrypt加密)
INSERT INTO `sys_user` (`username`, `password`, `nickname`, `status`) 
VALUES ('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOhQDfVaAabD6', '超级管理员', 1)
ON DUPLICATE KEY UPDATE `nickname` = '超级管理员';

-- 插入默认角色
INSERT INTO `sys_role` (`name`, `code`, `description`) VALUES 
('超级管理员', 'SUPER_ADMIN', '拥有系统所有权限'),
('运维管理员', 'ADMIN', '可以管理任务和执行器'),
('普通用户', 'USER', '只能查看和执行自己的任务')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- 关联管理员角色
INSERT INTO `sys_user_role` (`user_id`, `role_id`) 
SELECT u.id, r.id FROM `sys_user` u, `sys_role` r 
WHERE u.username = 'admin' AND r.code = 'SUPER_ADMIN'
ON DUPLICATE KEY UPDATE `user_id` = `user_id`;

-- 插入系统配置
INSERT INTO `sys_config` (`config_key`, `config_value`, `config_type`, `description`) VALUES
('scheduler.thread_pool_size', '20', 'NUMBER', '调度器线程池大小'),
('scheduler.trigger_pool_size', '200', 'NUMBER', '触发器线程池大小'),
('executor.heartbeat_interval', '30', 'NUMBER', '执行器心跳间隔(秒)'),
('executor.dead_timeout', '90', 'NUMBER', '执行器离线判定时间(秒)'),
('log.retention_days', '30', 'NUMBER', '日志保留天数'),
('alarm.email_enabled', 'true', 'BOOLEAN', '是否启用邮件告警')
ON DUPLICATE KEY UPDATE `config_value` = VALUES(`config_value`);

-- 插入示例任务组
INSERT INTO `task_group` (`name`, `description`, `app_name`, `status`) VALUES
('默认执行器组', '系统默认的执行器组', 'default-executor', 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ============================================
-- 完成
-- ============================================
SELECT '数据库初始化完成!' AS message;

