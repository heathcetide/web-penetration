-- 目标资产表
CREATE TABLE target_asset (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    domain VARCHAR(255) COMMENT '域名',
    ip VARCHAR(50) COMMENT 'IP地址',
    status TINYINT COMMENT '状态',
    create_time DATETIME,
    update_time DATETIME
);

-- 扫描任务表
CREATE TABLE scan_task (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_name VARCHAR(100) COMMENT '任务名称',
    task_type TINYINT COMMENT '任务类型：1-资产收集，2-漏洞扫描',
    target_id BIGINT COMMENT '目标ID',
    status TINYINT COMMENT '状态：0-待执行，1-执行中，2-已完成',
    create_time DATETIME,
    update_time DATETIME
);

-- 漏洞信息表
CREATE TABLE vulnerability (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT COMMENT '关联任务ID',
    vuln_type VARCHAR(50) COMMENT '漏洞类型',
    vuln_level TINYINT COMMENT '漏洞等级',
    vuln_desc TEXT COMMENT '漏洞描述',
    create_time DATETIME
);

-- 网站信息表
CREATE TABLE website_info (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT COMMENT '关联任务ID',
    url VARCHAR(255) COMMENT '网站URL',
    title VARCHAR(255) COMMENT '网站标题',
    description TEXT COMMENT '网站描述',
    keywords VARCHAR(500) COMMENT '关键词',
    server VARCHAR(100) COMMENT '服务器信息',
    framework VARCHAR(100) COMMENT '框架信息',
    links TEXT COMMENT '页面链接',
    emails TEXT COMMENT '邮箱信息',
    phones TEXT COMMENT '电话信息',
    images TEXT COMMENT '图片链接',
    videos TEXT COMMENT '视频链接',
    files TEXT COMMENT '文件链接',
    comics TEXT COMMENT '漫画图片链接',
    create_time DATETIME COMMENT '创建时间',
    KEY idx_task_id (task_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='网站信息表';

-- 资源下载表
CREATE TABLE resource_download (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT COMMENT '关联任务ID',
    resource_url VARCHAR(1000) COMMENT '资源URL',
    resource_type VARCHAR(20) COMMENT '资源类型',
    local_path VARCHAR(255) COMMENT '本地存储路径',
    status TINYINT COMMENT '状态：0-待下载 1-下载中 2-已完成 3-失败',
    error_msg VARCHAR(500) COMMENT '错误信息',
    create_time DATETIME,
    update_time DATETIME,
    KEY idx_task_id (task_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源下载表';

ALTER TABLE website_info ADD COLUMN download_stats TEXT; 