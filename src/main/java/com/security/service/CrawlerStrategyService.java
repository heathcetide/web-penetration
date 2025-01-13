package com.security.service;

/**
 * 爬虫策略服务接口，用于定义爬虫的爬取策略。
 */
public interface CrawlerStrategyService {
    void setCrawlDelay(int delay); // 设置爬取延迟
    void setCrawlPolicy(String policy); // 设置爬取策略（如深度优先、广度优先）
    void setAllowedDomains(String[] domains); // 设置允许爬取的域名
} 