package com.security.service;

public interface IVulnScanService {
    /**
     * 创建漏洞扫描任务
     */
    Long createVulnScanTask(Long targetId);
    
    /**
     * SQL注入检测
     */
    void sqlInjectionScan(Long taskId, String url);
    
    /**
     * XSS漏洞检测
     */
    void xssScan(Long taskId, String url);
    
    /**
     * CSRF漏洞检测
     */
    void csrfScan(Long taskId, String url);
} 