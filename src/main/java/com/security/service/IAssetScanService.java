package com.security.service;

public interface IAssetScanService {
    /**
     * 创建资产扫描任务
     */
    Long createAssetScanTask(String domain);
    
    /**
     * 执行端口扫描
     */
    void portScan(Long taskId, String target);
    
    /**
     * 执行子域名枚举
     */
    void subdomainScan(Long taskId, String domain);
    
    /**
     * 执行目录扫描
     */
    void directoryScan(Long taskId, String url);
} 