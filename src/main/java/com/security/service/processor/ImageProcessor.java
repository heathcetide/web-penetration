package com.security.service.processor;

import net.coobird.thumbnailator.Thumbnails;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.io.File;
import java.util.Map;

@Component
public class ImageProcessor extends ResourceProcessor {

    private static Logger log = LoggerFactory.getLogger(ImageProcessor.class);
    @Override
    protected void preProcess(File file) {
        // 检查图片格式
        // 检查图片大小
        // 检查图片分辨率
    }
    
    @Override
    protected void doProcess(File file) {
        // 生成缩略图
        generateThumbnail(file);
        // 压缩图片
        compressImage(file);
        // 去除EXIF信息
        removeExif(file);
    }
    
    @Override
    protected void postProcess(File file, Map<String, String> metadata) {
        // 更新元数据
        // 添加水印
        // 转换格式
    }
    
    @Override
    protected void handleError(File file, Exception e) {
        // 记录错误日志
        // 移动到错误目录
        // 发送通知
    }
    
    private void generateThumbnail(File file) {
        try {
            String thumbPath = file.getParent() + "/thumb_" + file.getName();
            Thumbnails.of(file)
                    .size(200, 200)
                    .keepAspectRatio(true)
                    .toFile(new File(thumbPath));
        } catch (Exception e) {
            log.error("生成缩略图失败", e);
        }
    }
    
    private void compressImage(File file) {
        try {
            String compressedPath = file.getParent() + "/compressed_" + file.getName();
            Thumbnails.of(file)
                    .scale(1.0)
                    .outputQuality(0.8)
                    .toFile(new File(compressedPath));
        } catch (Exception e) {
            log.error("压缩图片失败", e);
        }
    }
    
    private void removeExif(File file) {
        // 实现EXIF信息移除
    }
} 