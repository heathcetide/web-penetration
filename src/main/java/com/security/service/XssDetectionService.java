package com.security.service;

/**
 * XSS 检测服务接口，用于检测 Web 应用是否存在跨站脚本漏洞。 [刘铭昊]
 */
public interface XssDetectionService {

    /**
     * 检测指定 URL 是否存在 XSS 漏洞。
     *
     * @param url     要检测的目标 URL
     * @param payload XSS 攻击载荷，用于测试漏洞
     * @return 如果存在 XSS 漏洞返回 true，否则返回 false
     */
    boolean detectXss(String url, String payload);

    /**
     * 测试指定 URL 是否存在存储型 XSS 漏洞。
     *
     * @param url 要测试的目标 URL
     * @return 如果存在存储型 XSS 漏洞返回 true，否则返回 false
     */
    boolean testForStoredXss(String url);
}
