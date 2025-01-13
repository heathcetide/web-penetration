package com.security.service.video.parser;

import org.springframework.stereotype.Component;

import java.util.Map;

@Component
public class DefaultVideoParser implements VideoParser {
    @Override
    public boolean supports(String url) {
        return true;  // 作为默认解析器支持所有URL
    }
    
    @Override
    public VideoInfo parse(String url, Map<String, Object> context) {
        VideoInfo info = new VideoInfo();
        // 使用之前的通用解析逻辑
        return info;
    }
    
    @Override
    public int getOrder() {
        return Integer.MAX_VALUE;  // 最低优先级
    }
    
    @Override
    public String getName() {
        return "Default";
    }
} 