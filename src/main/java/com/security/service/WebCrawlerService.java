package com.security.service;

import com.security.model.entity.CrawlResult;

/**
 * 爬虫服务接口，用于自动爬取目标网站并收集信息。
 */
public interface WebCrawlerService {

    /**
     * 爬取目标网站的基本信息。
     *
     * @param url 要爬取的目标网站 URL
     * @return 爬取结果，包括页面内容和相关信息
     */
    CrawlResult crawlWebsite(String url);

    /**
     * 爬取目标网站的资源（如图片、视频等）。
     *
     * @param url 要爬取的目标网站 URL
     * @return 爬取结果，包括页面内容和资源信息
     */
    CrawlResult crawlWebsiteWithResources(String url);
}
