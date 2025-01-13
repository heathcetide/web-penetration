package com.security.service;

/**
 * 爬虫结果处理服务接口，用于处理爬取结果。
 */
public interface CrawlResultHandlerService {
    void saveCrawlResults(String url, String results); // 保存爬取结果
    void exportResultsToJson(String url); // 导出结果为 JSON 格式
    void exportResultsToCsv(String url); // 导出结果为 CSV 格式
} 