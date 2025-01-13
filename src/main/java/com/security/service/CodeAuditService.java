package com.security.service;

/**
 * 代码审计和静态分析服务接口，用于对 Web 应用的源代码进行审计和静态分析。
 */
public interface CodeAuditService {
    boolean auditCode(String sourceCode);
} 