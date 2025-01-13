package com.security.service;

/**
 * 进度监控服务接口，用于监控学习和视频播放的进度。
 */
public interface ProgressMonitoringService {
    void monitorLearningProgress(String activityId); // 监控学习进度
    void monitorVideoProgress(String videoId); // 监控视频播放进度
    void reportProgress(String activityId, int progress); // 报告学习进度
    void reportVideoProgress(String videoId, int progress); // 报告视频播放进度
} 