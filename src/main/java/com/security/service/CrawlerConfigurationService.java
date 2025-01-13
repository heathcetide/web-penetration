package com.security.service;

/**
 * 爬虫配置服务接口，用于配置爬虫的参数和选项。
 */
public interface CrawlerConfigurationService {
    void setUserAgent(String userAgent); // 设置用户代理
    void setMaxDepth(int maxDepth); // 设置最大爬取深度
    void setTimeout(int timeout); // 设置请求超时
    void setMaxPages(int maxPages); // 设置最大爬取页面数
}

