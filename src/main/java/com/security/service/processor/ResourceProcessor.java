package com.security.service.processor;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.io.File;
import java.util.Map;

@Component
public abstract class ResourceProcessor {

    private static Logger log = LoggerFactory.getLogger(ResourceProcessor.class);

    public void process(File file, Map<String, String> metadata) {
        try {
            // 预处理
            preProcess(file);
            
            // 处理
            doProcess(file);
            
            // 后处理
            postProcess(file, metadata);
            
        } catch (Exception e) {
            log.error("资源处理失败: {}", e.getMessage(), e);
            handleError(file, e);
        }
    }
    
    protected abstract void preProcess(File file);
    protected abstract void doProcess(File file);
    protected abstract void postProcess(File file, Map<String, String> metadata);
    protected abstract void handleError(File file, Exception e);
} 