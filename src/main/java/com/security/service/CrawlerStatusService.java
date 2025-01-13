package com.security.service;

/**
 * 爬虫状态监控服务接口，用于监控爬虫的运行状态。
 */
public interface CrawlerStatusService {
    String getCrawlStatus(String taskId); // 获取爬取任务状态
    int getTotalPagesCrawled(String taskId); // 获取已爬取的总页面数
    int getTotalErrors(String taskId); // 获取爬取过程中发生的错误总数
} 