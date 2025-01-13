package com.security.service;

/**
 * 代码注入检测服务接口，用于检测 Web 应用是否存在代码注入漏洞。
 */
public interface CodeInjectionService {
    boolean detectCodeInjection(String url, String payload);
} 