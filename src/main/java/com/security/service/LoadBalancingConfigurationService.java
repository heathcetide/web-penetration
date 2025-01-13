package com.security.service;

/**
 * 负载均衡配置检查服务接口，用于检查 Web 应用的负载均衡配置。 【安】
 */
public interface LoadBalancingConfigurationService {
    boolean checkLoadBalancingConfiguration(String url);
    boolean testSessionPersistence(String url); // 测试会话持久性
    boolean validateLoadBalancingAlgorithm(String url); // 验证负载均衡算法
    boolean checkForLoadBalancingFailures(String url); // 检查负载均衡失败
    boolean analyzeTrafficDistribution(String url); // 分析流量分配
    boolean validateHealthCheckEndpoints(String url); // 验证健康检查端点
    boolean testLoadBalancing(String url);//负载均衡测试服务接口，用于测试 Web 应用的负载均衡配置。
} 