create database web_pentest_db;

use web_pentest_db;

-- 用户表
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 角色表
CREATE TABLE `role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 用户角色关联表
CREATE TABLE `user_role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';



-- 爬虫任务表
CREATE TABLE `crawler_task` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '任务名称',
  `description` varchar(500) DEFAULT NULL COMMENT '任务描述',
  `start_url` varchar(500) NOT NULL COMMENT '起始URL',
  `max_depth` int NOT NULL DEFAULT '0' COMMENT '最大爬取深度',
  `concurrent_count` int NOT NULL DEFAULT '1' COMMENT '并发数',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-待执行，1-执行中，2-已完成，3-失败',
  `created_by` bigint NOT NULL COMMENT '创建人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='爬虫任务表';

-- 爬虫配置表
CREATE TABLE `crawler_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL COMMENT '任务ID',
  `config_key` varchar(50) NOT NULL COMMENT '配置键',
  `config_value` varchar(500) NOT NULL COMMENT '配置值',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='爬虫配置表';

-- 爬虫结果表
CREATE TABLE `crawler_result` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL COMMENT '任务ID',
  `url` varchar(500) NOT NULL COMMENT '爬取URL',
  `title` varchar(200) DEFAULT NULL COMMENT '页面标题',
  `content` text COMMENT '页面内容',
  `depth` int NOT NULL DEFAULT '0' COMMENT '爬取深度',
  `status_code` int DEFAULT NULL COMMENT 'HTTP状态码',
  `content_type` varchar(100) DEFAULT NULL COMMENT '内容类型',
  `download_time` bigint DEFAULT NULL COMMENT '下载时间(ms)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_url` (`url`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='爬虫结果表';


-- 系统配置表
CREATE TABLE `system_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `config_key` varchar(50) NOT NULL COMMENT '配置键',
  `config_value` varchar(500) NOT NULL COMMENT '配置值',
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 操作日志表
CREATE TABLE `operation_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `operation` varchar(50) NOT NULL COMMENT '操作类型',
  `description` varchar(500) DEFAULT NULL COMMENT '操作描述',
  `ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';



-- 爬虫统计表
CREATE TABLE `crawler_stats` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL COMMENT '任务ID',
  `total_urls` int NOT NULL DEFAULT '0' COMMENT '总URL数',
  `success_urls` int NOT NULL DEFAULT '0' COMMENT '成功URL数',
  `failed_urls` int NOT NULL DEFAULT '0' COMMENT '失败URL数',
  `total_time` bigint NOT NULL DEFAULT '0' COMMENT '总耗时(ms)',
  `avg_time` bigint NOT NULL DEFAULT '0' COMMENT '平均耗时(ms)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='爬虫统计表';


-- 权限表
CREATE TABLE `permission` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '权限名称',
  `code` varchar(50) NOT NULL COMMENT '权限编码',
  `type` varchar(20) NOT NULL COMMENT '权限类型',
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 角色权限关联表
CREATE TABLE `role_permission` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `permission_id` bigint NOT NULL COMMENT '权限ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 工作流表
CREATE TABLE `workflow` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '工作流名称',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_by` bigint NOT NULL COMMENT '创建人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流表';

-- 工作流实例表
CREATE TABLE `workflow_instance` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `workflow_id` bigint NOT NULL COMMENT '工作流ID',
  `status` tinyint NOT NULL COMMENT '状态',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `created_by` bigint NOT NULL COMMENT '创建人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_workflow_id` (`workflow_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流实例表';

-- 工作流任务表
CREATE TABLE `workflow_task` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `instance_id` bigint NOT NULL COMMENT '工作流实例ID',
  `task_type` varchar(50) NOT NULL COMMENT '任务类型',
  `task_config` text COMMENT '任务配置',
  `status` tinyint NOT NULL COMMENT '状态',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_instance_id` (`instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流任务表';

-- 任务结果表
CREATE TABLE `task_result` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL COMMENT '任务ID',
  `result_type` varchar(50) NOT NULL COMMENT '结果类型',
  `result_data` text COMMENT '结果数据',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务结果表';

-- 工作流变量表
CREATE TABLE `workflow_variable` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `instance_id` bigint NOT NULL COMMENT '工作流实例ID',
  `var_key` varchar(50) NOT NULL COMMENT '变量键',
  `var_value` text COMMENT '变量值',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_instance_id` (`instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流变量表';

-- 安全度量表
CREATE TABLE `security_metric` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '指标名称',
  `type` varchar(50) NOT NULL COMMENT '指标类型',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全度量表';

-- 度量历史表
CREATE TABLE `metric_history` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `metric_id` bigint NOT NULL COMMENT '指标ID',
  `value` decimal(10,2) NOT NULL COMMENT '指标值',
  `measure_time` datetime NOT NULL COMMENT '度量时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_metric_id` (`metric_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='度量历史表';

-- 漏洞表
CREATE TABLE `vulnerability` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL COMMENT '漏洞名称',
  `type` varchar(50) NOT NULL COMMENT '漏洞类型',
  `severity` varchar(20) NOT NULL COMMENT '严重程度',
  `description` text COMMENT '描述',
  `solution` text COMMENT '解决方案',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='漏洞表';

-- 漏洞详情表
CREATE TABLE `vuln_detail` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `vuln_id` bigint NOT NULL COMMENT '漏洞ID',
  `target_url` varchar(500) NOT NULL COMMENT '目标URL',
  `request` text COMMENT '请求数据',
  `response` text COMMENT '响应数据',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_vuln_id` (`vuln_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='漏洞详情表';

-- 报告表
CREATE TABLE `report` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '报告名称',
  `type` varchar(50) NOT NULL COMMENT '报告类型',
  `content` text COMMENT '报告内容',
  `created_by` bigint NOT NULL COMMENT '创建人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报告表';

-- 报告模板表
CREATE TABLE `report_template` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '模板名称',
  `content` text NOT NULL COMMENT '模板内容',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报告模板表';

-- 扫描目标表
CREATE TABLE `scan_target` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `url` varchar(500) NOT NULL COMMENT '目标URL',
  `type` varchar(50) NOT NULL COMMENT '目标类型',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态',
  `last_scan_time` datetime DEFAULT NULL COMMENT '最后扫描时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='扫描目标表';


-- 任务调度表
CREATE TABLE `task_schedule` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL COMMENT '任务ID',
  `schedule_type` varchar(20) NOT NULL COMMENT '调度类型：once-一次性，cron-定时',
  `cron_expression` varchar(100) DEFAULT NULL COMMENT 'cron表达式',
  `next_run_time` datetime DEFAULT NULL COMMENT '下次运行时间',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务调度表';

-- 任务执行表
CREATE TABLE `task_execution` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL COMMENT '任务ID',
  `status` tinyint NOT NULL COMMENT '执行状态',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `error_message` text COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务执行表';

-- 任务依赖表
CREATE TABLE `task_dependency` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL COMMENT '任务ID',
  `dependent_task_id` bigint NOT NULL COMMENT '依赖任务ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_dependency` (`task_id`,`dependent_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务依赖表';

-- 工作流审计表
CREATE TABLE `workflow_audit` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `workflow_id` bigint NOT NULL COMMENT '工作流ID',
  `instance_id` bigint DEFAULT NULL COMMENT '实例ID',
  `operation` varchar(50) NOT NULL COMMENT '操作类型',
  `operator_id` bigint NOT NULL COMMENT '操作人ID',
  `detail` text COMMENT '详细信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_workflow_id` (`workflow_id`),
  KEY `idx_instance_id` (`instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流审计表';

-- 工作流触发器表
CREATE TABLE `workflow_trigger` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `workflow_id` bigint NOT NULL COMMENT '工作流ID',
  `trigger_type` varchar(50) NOT NULL COMMENT '触发器类型',
  `trigger_config` text COMMENT '触发器配置',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_workflow_id` (`workflow_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作流触发器表';

-- 安全KPI表
CREATE TABLE `security_kpi` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT 'KPI名称',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `target_value` decimal(10,2) NOT NULL COMMENT '目标值',
  `weight` decimal(5,2) NOT NULL COMMENT '权重',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全KPI表';

-- KPI结果表
CREATE TABLE `kpi_result` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `kpi_id` bigint NOT NULL COMMENT 'KPI ID',
  `actual_value` decimal(10,2) NOT NULL COMMENT '实际值',
  `score` decimal(5,2) NOT NULL COMMENT '得分',
  `measure_time` datetime NOT NULL COMMENT '度量时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_kpi_id` (`kpi_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='KPI结果表';

-- 安全评分卡表
CREATE TABLE `security_scorecard` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '评分卡名称',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `total_score` decimal(5,2) NOT NULL COMMENT '总分',
  `measure_time` datetime NOT NULL COMMENT '度量时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全评分卡表';

-- 缓存配置表
CREATE TABLE `cache_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `cache_key` varchar(100) NOT NULL COMMENT '缓存键',
  `ttl` int NOT NULL COMMENT '过期时间(秒)',
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cache_key` (`cache_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='缓存配置表';

-- 缓存统计表
CREATE TABLE `cache_stats` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `cache_key` varchar(100) NOT NULL COMMENT '缓存键',
  `hits` bigint NOT NULL DEFAULT '0' COMMENT '命中次数',
  `misses` bigint NOT NULL DEFAULT '0' COMMENT '未命中次数',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cache_key` (`cache_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='缓存统计表';

-- 系统日志表
CREATE TABLE `system_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `log_type` varchar(50) NOT NULL COMMENT '日志类型',
  `log_level` varchar(20) NOT NULL COMMENT '日志级别',
  `message` text NOT NULL COMMENT '日志信息',
  `stack_trace` text COMMENT '堆栈信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_log_type` (`log_type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统日志表';

-- 安全日志表
CREATE TABLE `security_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `event_type` varchar(50) NOT NULL COMMENT '事件类型',
  `severity` varchar(20) NOT NULL COMMENT '严重程度',
  `source_ip` varchar(50) DEFAULT NULL COMMENT '来源IP',
  `target` varchar(200) DEFAULT NULL COMMENT '目标',
  `description` text COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_event_type` (`event_type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全日志表';

-- 配置变更表
CREATE TABLE `config_change` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `config_key` varchar(50) NOT NULL COMMENT '配置键',
  `old_value` text COMMENT '旧值',
  `new_value` text COMMENT '新值',
  `changed_by` bigint NOT NULL COMMENT '变更人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配置变更表';

-- 报告章节表
CREATE TABLE `report_section` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `report_id` bigint NOT NULL COMMENT '报告ID',
  `title` varchar(200) NOT NULL COMMENT '章节标题',
  `content` text COMMENT '章节内容',
  `order_num` int NOT NULL DEFAULT '0' COMMENT '排序号',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_report_id` (`report_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报告章节表';

-- 扫描结果表
CREATE TABLE `scan_result` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `target_id` bigint NOT NULL COMMENT '目标ID',
  `scan_type` varchar(50) NOT NULL COMMENT '扫描类型',
  `result_data` text COMMENT '扫描结果',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `status` tinyint NOT NULL COMMENT '状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_target_id` (`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='扫描结果表';

-- MFA认证表
CREATE TABLE `mfa_auth` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `secret_key` varchar(100) NOT NULL COMMENT '密钥',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-未启用，1-已启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='MFA认证表';

-- 会话表
CREATE TABLE `session` (
  `id` varchar(64) NOT NULL COMMENT '会话ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `ip_address` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` varchar(200) DEFAULT NULL COMMENT '用户代理',
  `expired_at` datetime NOT NULL COMMENT '过期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_expired_at` (`expired_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话表';

-- 邮件验证表
CREATE TABLE `email_verify` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `verify_code` varchar(64) NOT NULL COMMENT '验证码',
  `expired_at` datetime NOT NULL COMMENT '过期时间',
  `verified` tinyint NOT NULL DEFAULT '0' COMMENT '是否已验证：0-未验证，1-已验证',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_verify_code` (`verify_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮件验证表';

-- API访问记录表
CREATE TABLE `api_access_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `api_path` varchar(200) NOT NULL COMMENT 'API路径',
  `method` varchar(10) NOT NULL COMMENT '请求方法',
  `params` text COMMENT '请求参数',
  `response_code` int DEFAULT NULL COMMENT '响应码',
  `response_time` int DEFAULT NULL COMMENT '响应时间(ms)',
  `ip_address` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='API访问记录表';

-- 定时任务表
CREATE TABLE `scheduled_task` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '任务名称',
  `cron_expression` varchar(100) NOT NULL COMMENT 'cron表达式',
  `task_class` varchar(200) NOT NULL COMMENT '任务类名',
  `task_method` varchar(100) NOT NULL COMMENT '任务方法',
  `params` text COMMENT '任务参数',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='定时任务表';

-- 通知配置表
CREATE TABLE `notification_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `type` varchar(50) NOT NULL COMMENT '通知类型',
  `config` text NOT NULL COMMENT '配置内容',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知配置表';

-- 通知记录表
CREATE TABLE `notification_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `type` varchar(50) NOT NULL COMMENT '通知类型',
  `target` varchar(200) NOT NULL COMMENT '通知目标',
  `content` text NOT NULL COMMENT '通知内容',
  `status` tinyint NOT NULL COMMENT '状态',
  `error_msg` varchar(500) DEFAULT NULL COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知记录表';


-- 用户组表
CREATE TABLE `user_group` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '组名',
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
  `parent_id` bigint DEFAULT NULL COMMENT '父组ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户组表';

-- 用户组成员表
CREATE TABLE `user_group_member` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `group_id` bigint NOT NULL COMMENT '组ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_user` (`group_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户组成员表';

-- 用户登录历史表
CREATE TABLE `user_login_history` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `login_time` datetime NOT NULL COMMENT '登录时间',
  `login_ip` varchar(50) DEFAULT NULL COMMENT '登录IP',
  `login_type` varchar(20) NOT NULL COMMENT '登录类型',
  `device_info` varchar(200) DEFAULT NULL COMMENT '设备信息',
  `status` tinyint NOT NULL COMMENT '状态：0-失败，1-成功',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户登录历史表';

-- 密码重置表
CREATE TABLE `password_reset` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `reset_token` varchar(100) NOT NULL COMMENT '重置令牌',
  `expired_at` datetime NOT NULL COMMENT '过期时间',
  `used` tinyint NOT NULL DEFAULT '0' COMMENT '是否已使用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_reset_token` (`reset_token`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='密码重置表';

-- 资源表
CREATE TABLE `resource` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '资源名称',
  `type` varchar(50) NOT NULL COMMENT '资源类型',
  `url` varchar(200) NOT NULL COMMENT '资源URL',
  `method` varchar(10) DEFAULT NULL COMMENT 'HTTP方法',
  `permission_code` varchar(50) DEFAULT NULL COMMENT '权限编码',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源表';

-- 任务标签表
CREATE TABLE `task_tag` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '标签名',
  `color` varchar(20) DEFAULT NULL COMMENT '标签颜色',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务标签表';

-- 任务标签关联表
CREATE TABLE `task_tag_relation` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL COMMENT '任务ID',
  `tag_id` bigint NOT NULL COMMENT '标签ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_tag` (`task_id`,`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务标签关联表';

-- 系统公告表
CREATE TABLE `system_announcement` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(200) NOT NULL COMMENT '公告标题',
  `content` text NOT NULL COMMENT '公告内容',
  `type` varchar(50) NOT NULL COMMENT '公告类型',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-未发布，1-已发布',
  `publish_time` datetime DEFAULT NULL COMMENT '发布时间',
  `created_by` bigint NOT NULL COMMENT '创建人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统公告表';


-- 数据字典表
CREATE TABLE `dict` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `type` varchar(50) NOT NULL COMMENT '字典类型',
  `code` varchar(50) NOT NULL COMMENT '字典编码',
  `name` varchar(100) NOT NULL COMMENT '字典名称',
  `value` varchar(100) NOT NULL COMMENT '字典值',
  `order_num` int DEFAULT '0' COMMENT '排序号',
  `remark` varchar(200) DEFAULT NULL COMMENT '备注',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_type_code` (`type`,`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据字典表';

-- 区域表
CREATE TABLE `region` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `parent_id` bigint DEFAULT NULL COMMENT '父级ID',
  `name` varchar(100) NOT NULL COMMENT '区域名称',
  `code` varchar(20) NOT NULL COMMENT '区域编码',
  `level` tinyint NOT NULL COMMENT '层级',
  `order_num` int DEFAULT '0' COMMENT '排序号',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='区域表';

-- 文件上传表
CREATE TABLE `file_upload` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `file_name` varchar(200) NOT NULL COMMENT '文件名',
  `file_path` varchar(500) NOT NULL COMMENT '文件路径',
  `file_size` bigint NOT NULL COMMENT '文件大小',
  `file_type` varchar(100) NOT NULL COMMENT '文件类型',
  `md5` varchar(32) DEFAULT NULL COMMENT '文件MD5',
  `upload_by` bigint NOT NULL COMMENT '上传人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_upload_by` (`upload_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件上传表';

-- 系统参数表
CREATE TABLE `system_param` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `param_key` varchar(50) NOT NULL COMMENT '参数键',
  `param_value` varchar(500) NOT NULL COMMENT '参数值',
  `param_type` varchar(20) NOT NULL COMMENT '参数类型',
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_param_key` (`param_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统参数表';

-- 接口文档表
CREATE TABLE `api_doc` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `api_name` varchar(100) NOT NULL COMMENT '接口名称',
  `api_path` varchar(200) NOT NULL COMMENT '接口路径',
  `method` varchar(10) NOT NULL COMMENT '请求方法',
  `description` text COMMENT '接口描述',
  `request_params` text COMMENT '请求参数',
  `response_params` text COMMENT '响应参数',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_api_path_method` (`api_path`,`method`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='接口文档表';

-- 接口版本表
CREATE TABLE `api_version` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `api_id` bigint NOT NULL COMMENT '接口ID',
  `version` varchar(20) NOT NULL COMMENT '版本号',
  `change_log` text COMMENT '变更日志',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_api_id` (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='接口版本表';

-- 数据备份表
CREATE TABLE `data_backup` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `backup_name` varchar(100) NOT NULL COMMENT '备份名称',
  `backup_type` varchar(50) NOT NULL COMMENT '备份类型',
  `file_path` varchar(500) NOT NULL COMMENT '文件路径',
  `file_size` bigint NOT NULL COMMENT '文件大小',
  `status` tinyint NOT NULL COMMENT '状态',
  `created_by` bigint NOT NULL COMMENT '创建人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据备份表';

-- 数据恢复记录表
CREATE TABLE `data_restore` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `backup_id` bigint NOT NULL COMMENT '备份ID',
  `restore_time` datetime NOT NULL COMMENT '恢复时间',
  `status` tinyint NOT NULL COMMENT '状态',
  `error_msg` text COMMENT '错误信息',
  `created_by` bigint NOT NULL COMMENT '创建人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_backup_id` (`backup_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据恢复记录表';


-- 告警规则表
CREATE TABLE `alert_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '规则名称',
  `rule_type` varchar(50) NOT NULL COMMENT '规则类型',
  `rule_content` text NOT NULL COMMENT '规则内容',
  `severity` varchar(20) NOT NULL COMMENT '严重程度',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警规则表';

-- 告警记录表
CREATE TABLE `alert_record` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `rule_id` bigint NOT NULL COMMENT '规则ID',
  `alert_target` varchar(200) NOT NULL COMMENT '告警目标',
  `alert_content` text NOT NULL COMMENT '告警内容',
  `severity` varchar(20) NOT NULL COMMENT '严重程度',
  `status` tinyint NOT NULL COMMENT '状态',
  `processed_by` bigint DEFAULT NULL COMMENT '处理人ID',
  `process_time` datetime DEFAULT NULL COMMENT '处理时间',
  `process_result` varchar(500) DEFAULT NULL COMMENT '处理结果',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_rule_id` (`rule_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警记录表';

-- 告警通知配置表
CREATE TABLE `alert_notify_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `rule_id` bigint NOT NULL COMMENT '规则ID',
  `notify_type` varchar(50) NOT NULL COMMENT '通知类型',
  `notify_target` varchar(500) NOT NULL COMMENT '通知目标',
  `notify_template` text COMMENT '通知模板',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_rule_id` (`rule_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警通知配置表';

-- 告警抑制规则表
CREATE TABLE `alert_suppress_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `rule_id` bigint NOT NULL COMMENT '规则ID',
  `suppress_type` varchar(50) NOT NULL COMMENT '抑制类型',
  `suppress_condition` text NOT NULL COMMENT '抑制条件',
  `suppress_duration` int NOT NULL COMMENT '抑制时长(秒)',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_rule_id` (`rule_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警抑制规则表';