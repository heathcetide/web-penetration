package com.security.service;

/**
 * 代码注入检测服务接口，用于检测 Web 应用是否存在代码注入漏洞。
 */
public interface CodeInjectionTestingService {
    boolean detectCodeInjection(String url, String payload);
    boolean validateInjectionPoints(String url); // 验证注入点
    boolean checkForCodeInjectionVulnerabilities(String url); // 检查代码注入漏洞
    boolean analyzeInjectionPatterns(String url); // 分析注入模式
    boolean testForBlindInjection(String url); // 测试盲注入
} 