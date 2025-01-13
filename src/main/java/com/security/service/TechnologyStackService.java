package com.security.service;

/**
 * 技术栈服务接口，用于获取 Web 应用的技术栈信息。[刘铭昊]
 */
public interface TechnologyStackService {
    /**
     * 获取当前对象所使用的技术栈信息。
     * @param url Web 应用的 URL
     * @return 技术栈信息的字符串描述
     */
    String getTechnologyStack(String url);
} 