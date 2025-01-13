package com.security.service;

/**
 * 爬虫数据分析服务接口，用于分析爬取的数据。
 */
public interface CrawlDataAnalysisService {
    void analyzeCrawlData(String url); // 分析爬取的数据
    String generateCrawlReport(String url); // 生成爬取报告
} 