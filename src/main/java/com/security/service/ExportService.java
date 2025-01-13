package com.security.service;

import com.security.model.entity.CrawlResult;

/**
 * 导出服务接口，用于将爬虫结果导出为不同格式。 [蒋浩天]
 */
public interface ExportService {
    String exportToJson(CrawlResult crawlResult);
    String exportToCsv(CrawlResult crawlResult);
} 