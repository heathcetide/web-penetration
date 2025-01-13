package com.security.service;

/**
 * 爬虫调度服务接口，用于调度和管理爬虫任务。
 */
public interface CrawlerSchedulerService {
    void scheduleCrawl(String url); // 调度爬取任务
    void cancelCrawl(String taskId); // 取消爬取任务
    void pauseCrawl(String taskId); // 暂停爬取任务
    void resumeCrawl(String taskId); // 恢复爬取任务
} 