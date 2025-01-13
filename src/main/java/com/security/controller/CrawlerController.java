package com.security.controller;

import com.security.model.entity.CrawlResult;
import com.security.service.WebCrawlerService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

/**
 * 爬虫控制器，用于处理爬虫相关的 HTTP 请求。
 */
@RestController
@RequestMapping("/api/crawler")
public class CrawlerController {

//    @Autowired
//    private WebCrawlerService webCrawlerService;

    @GetMapping("/crawl")
    public CrawlResult crawlWebsite(@RequestParam String url) {
//        return webCrawlerService.crawlWebsite(url);
        return null;
    }

    @GetMapping("/crawl/resources")
    public CrawlResult crawlWebsiteWithResources(@RequestParam String url) {
//        return webCrawlerService.crawlWebsiteWithResources(url);
        return null;
    }
} 