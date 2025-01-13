package com.security.service;

/**
 * 视频播放服务接口，用于管理视频的播放状态和进度。
 */
public interface VideoPlaybackService {

    /**
     * 开始播放指定的视频。
     *
     * @param videoId 视频的唯一标识符
     */
    void startVideo(String videoId);

    /**
     * 暂停正在播放的视频。
     *
     * @param videoId 视频的唯一标识符
     */
    void pauseVideo(String videoId);

    /**
     * 恢复暂停的视频播放。
     *
     * @param videoId 视频的唯一标识符
     */
    void resumeVideo(String videoId);

    /**
     * 更新视频的播放进度。
     *
     * @param videoId  视频的唯一标识符
     * @param progress 当前播放进度（以秒为单位）
     */
    void updateVideoProgress(String videoId, int progress);

    /**
     * 获取视频的播放状态。
     *
     * @param videoId 视频的唯一标识符
     * @return 视频播放状态（如 "播放中"、"暂停"）
     */
    String getVideoPlaybackStatus(String videoId);
}
