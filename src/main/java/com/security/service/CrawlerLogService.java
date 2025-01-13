package com.security.service;

/**
 * 爬虫日志服务接口，用于记录和管理爬虫日志。
 */
public interface CrawlerLogService {
    void logInfo(String message); // 记录信息日志
    void logError(String message); // 记录错误日志
    String getLog(String taskId); // 获取特定任务的日志
} 