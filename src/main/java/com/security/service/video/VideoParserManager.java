package com.security.service.video;

import com.security.exception.BusinessException;
import com.security.service.video.parser.VideoInfo;
import com.security.service.video.parser.VideoParser;
import org.springframework.stereotype.Component;
import java.util.List;
import java.util.Comparator;
import java.util.Map;
import java.util.stream.Collectors;

@Component
public class VideoParserManager {
    private final List<VideoParser> parsers;
    
    public VideoParserManager(List<VideoParser> parsers) {
        this.parsers = parsers.stream()
            .sorted(Comparator.comparingInt(VideoParser::getOrder))
            .collect(Collectors.toList());
    }
    
    public VideoInfo parseVideo(String url, Map<String, Object> context) {
        // 找到第一个支持该URL的解析器
        VideoParser parser = parsers.stream()
            .filter(p -> p.supports(url))
            .findFirst()
            .orElseThrow(() -> new BusinessException("没有找到支持的视频解析器"));
            
        return parser.parse(url, context);
    }
} 