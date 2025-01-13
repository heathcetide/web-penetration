package com.security.service;

/**
 * 通用爬虫服务接口，用于爬取学习平台的内容。
 */
public interface GenericCrawlerService {
    void crawlLearningContent(String url); // 爬取学习内容
    void crawlVideoContent(String url); // 爬取视频内容
    void updateProgress(String url, String contentType, int progress); // 更新进度
} 