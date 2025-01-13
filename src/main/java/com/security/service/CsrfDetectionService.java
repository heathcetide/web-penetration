package com.security.service;

/**
 * CSRF 检测服务接口，用于检测 Web 应用是否存在跨站请求伪造漏洞。
 */
public interface CsrfDetectionService {
    boolean detectCsrf(String url); // 检测 CSRF
    boolean validateCsrfTokens(String url); // 验证 CSRF 令牌
} 