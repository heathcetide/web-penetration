package com.security.service;

import com.security.service.video.EncryptionInfo;

import java.util.List;
import java.io.File;

public interface IResourceProcessService {
    /**
     * 提取图片资源
     */
    List<String> extractImages(String html, String baseUrl);
    
    /**
     * 提取视频资源
     */
    List<String> extractVideos(String html, String baseUrl);
    
    /**
     * 提取文件资源
     */
    List<String> extractFiles(String html, String baseUrl);
    
    /**
     * 提取漫画图片
     */
    List<String> extractComics(String html, String baseUrl);
    
    /**
     * 下载资源
     */
    void downloadResource(Long taskId, String url, String type);
    
    /**
     * 合并视频片段
     */
    void mergeVideoSegments(Long taskId, List<String> segments, String outputPath);
    
    /**
     * 处理Token视频
     */
    void handleTokenVideo(String url, String token);
    
    /**
     * 处理HLS加密视频
     */
    void handleHLSEncryptedVideo(String m3u8Url, String keyUrl);
    
    /**
     * 处理加密视频
     */
    void handleEncryptedVideo(String url, EncryptionInfo encryptionInfo);
    
    /**
     * 处理视频后处理
     */
    void postProcessVideo(File videoFile);

    List<String> extractSubtitles(String html, String baseUrl);

    List<String> extractAudios(String html, String baseUrl);
} 