package com.security.service;

/**
 * 文件包含测试服务接口，用于测试 Web 应用是否存在文件包含漏洞。 [蒋浩天]
 */
public interface FileInclusionService {

    boolean testFileInclusion(String url);
} 