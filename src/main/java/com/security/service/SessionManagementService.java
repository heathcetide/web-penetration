package com.security.service;

/**
 * 会话管理测试服务接口，用于测试 Web 应用的会话管理机制。
 */
public interface SessionManagementService {

    /**
     * 测试 Web 应用的会话管理机制是否安全。
     *
     * @param url 要测试的目标 URL
     * @return 如果会话管理存在安全问题返回 true，否则返回 false
     */
    boolean testSessionManagement(String url);
}
