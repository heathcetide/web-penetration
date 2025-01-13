package com.security.service.video.parser;

import java.util.List;
import java.util.Map;

public interface VideoParser {
    /**
     * 判断是否支持该URL
     */
    boolean supports(String url);
    
    /**
     * 解析视频信息
     */
    VideoInfo parse(String url, Map<String, Object> context);
    
    /**
     * 获取解析器优先级
     */
    int getOrder();
    
    /**
     * 获取解析器名称
     */
    String getName();
} 