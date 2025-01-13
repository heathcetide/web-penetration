package com.security.service;

/**
 * 日志监控服务接口，用于检查 Web 应用的日志记录和监控机制。
 */
public interface LoggingMonitoringService {
    boolean checkLogging(String url);
} 