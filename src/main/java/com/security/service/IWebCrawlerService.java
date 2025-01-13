package com.security.service;

import com.security.model.dto.response.CrawlerResultDTO;

import java.util.List;

public interface IWebCrawlerService {
    /**
     * 创建爬虫任务
     */
    Long createCrawlerTask(String url, Integer depth);
    
    /**
     * 基本信息收集
     * - 网站标题
     * - Meta信息
     * - 服务器信息
     */
    void basicInfoCrawl(Long taskId, String url);
    
    /**
     * 链接爬取
     * - 内部链接
     * - 外部链接
     * - 资源链接（JS/CSS等）
     */
    List<String> linksCrawl(String url, Integer depth);
    
    /**
     * 敏感信息收集
     * - 邮箱
     * - 电话号码
     * - API接口
     */
    void sensitiveInfoCrawl(Long taskId, String url);
    
    /**
     * JS文件分析
     * - API接口提取
     * - 敏感信息提取
     */
    void jsAnalysis(Long taskId, String url);
    
    /**
     * 获取爬虫结果
     */
    CrawlerResultDTO getCrawlerResult(Long taskId);
} 