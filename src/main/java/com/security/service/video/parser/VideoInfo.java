package com.security.service.video.parser;

import java.util.List;
import java.util.Map;

public class VideoInfo {
    private String title;
    private String description;
    private List<VideoSource> sources;
    private Map<String, Object> metadata;
    
    public String getTitle() {
        return title;
    }
    
    public void setTitle(String title) {
        this.title = title;
    }
    
    public String getDescription() {
        return description;
    }
    
    public void setDescription(String description) {
        this.description = description;
    }
    
    public List<VideoSource> getSources() {
        return sources;
    }
    
    public void setSources(List<VideoSource> sources) {
        this.sources = sources;
    }
    
    public Map<String, Object> getMetadata() {
        return metadata;
    }
    
    public void setMetadata(Map<String, Object> metadata) {
        this.metadata = metadata;
    }
}

