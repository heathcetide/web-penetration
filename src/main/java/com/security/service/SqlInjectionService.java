package com.security.service;

/**
 * SQL 注入检测服务接口，用于检测 Web 应用是否存在 SQL 注入漏洞。
 */
public interface SqlInjectionService {

    /**
     * 检测指定 URL 是否存在 SQL 注入漏洞。
     *
     * @param url     要检测的目标 URL
     * @param payload 测试 SQL 注入的载荷
     * @return 如果存在 SQL 注入漏洞返回 true，否则返回 false
     */
    boolean detectSqlInjection(String url, String payload);

    /**
     * 测试指定 URL 是否存在盲 SQL 注入漏洞。
     *
     * @param url 要测试的目标 URL
     * @return 如果存在盲 SQL 注入漏洞返回 true，否则返回 false
     */
    boolean testForBlindSqlInjection(String url);
}
