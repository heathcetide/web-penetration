package com.security.tools;

import com.security.model.entity.CrawlResult;

/**
 * 爬虫工具接口，负责实现具体的爬虫逻辑。
 */
public interface WebCrawlerTool {
    CrawlResult executeCrawl(String url);
    CrawlResult executeCrawlWithResources(String url); // 爬取资源（图片、视频等）
} 