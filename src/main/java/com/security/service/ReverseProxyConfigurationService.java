package com.security.service;

/**
 * 反向代理配置检查服务接口，用于检查 Web 应用的反向代理配置。
 */
public interface ReverseProxyConfigurationService {
    boolean checkReverseProxyConfiguration(String url);
    boolean validateProxyHeaders(String url); // 验证代理请求头
    boolean testProxyTimeout(String url); // 测试代理超时设置
    boolean checkProxySSL(String url); // 检查代理的 SSL 配置
    boolean validateProxyAuthentication(String url); // 验证代理认证
    boolean checkProxyRedirects(String url); // 检查代理重定向
    boolean testProxyLoadBalancing(String url); // 测试代理负载均衡
    boolean checkReverseProxy(String url);//反向代理和负载均衡配置检查服务接口，用于检查 Web 应用的反向代理和负载均衡配置。

    boolean checkProxyConfiguration(String url);//反向代理配置检查服务接口，用于检查 Web 应用的反向代理配置。

    boolean testProxyConfiguration(String url); //反向代理和负载均衡测试服务接口，用于测试 Web 应用的反向代理和负载均衡配置。
} 