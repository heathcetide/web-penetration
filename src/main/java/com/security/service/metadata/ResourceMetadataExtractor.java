package com.security.service.metadata;

import com.drew.imaging.ImageMetadataReader;
import com.drew.metadata.Metadata;
import com.drew.metadata.Directory;
import com.drew.metadata.Tag;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import ws.schild.jave.MultimediaObject;
import ws.schild.jave.info.MultimediaInfo;

import java.io.File;
import java.util.HashMap;
import java.util.Map;
import java.nio.file.Files;
import java.util.Date;

@Component
public class ResourceMetadataExtractor {

    private static Logger log = LoggerFactory.getLogger(ResourceMetadataExtractor.class);
    public Map<String, String> extractMetadata(File file) {
        Map<String, String> metadata = new HashMap<>();
        
        try {
            String mimeType = getMimeType(file);
            metadata.put("mimeType", mimeType);
            
            switch (mimeType.split("/")[0]) {
                case "image":
                    extractImageMetadata(file, metadata);
                    break;
                case "video":
                    extractVideoMetadata(file, metadata);
                    break;
                case "audio":
                    extractAudioMetadata(file, metadata);
                    break;
                default:
                    extractBasicMetadata(file, metadata);
            }
        } catch (Exception e) {
            log.error("提取元数据失败: {}", e.getMessage(), e);
        }
        
        return metadata;
    }
    
    private void extractImageMetadata(File file, Map<String, String> metadata) {
        try {
            Metadata imageMetadata = ImageMetadataReader.readMetadata(file);
            for (Directory directory : imageMetadata.getDirectories()) {
                for (Tag tag : directory.getTags()) {
                    metadata.put(tag.getTagName(), tag.getDescription());
                }
            }
        } catch (Exception e) {
            log.error("提取图片元数据失败", e);
        }
    }
    
    private String getMimeType(File file) {
        try {
            return Files.probeContentType(file.toPath());
        } catch (Exception e) {
            log.error("获取MIME类型失败", e);
            return null;
        }
    }
    
    private void extractVideoMetadata(File file, Map<String, String> metadata) {
        try {
            MultimediaObject mediaObject = new MultimediaObject(file);
            MultimediaInfo info = mediaObject.getInfo();
            
            metadata.put("duration", String.valueOf(info.getDuration()));
            metadata.put("video.codec", info.getVideo().getDecoder());
            metadata.put("video.bitRate", String.valueOf(info.getVideo().getBitRate()));
            metadata.put("video.frameRate", String.valueOf(info.getVideo().getFrameRate()));
        } catch (Exception e) {
            log.error("提取视频元数据失败", e);
        }
    }
    
    private void extractAudioMetadata(File file, Map<String, String> metadata) {
        try {
            MultimediaObject mediaObject = new MultimediaObject(file);
            MultimediaInfo info = mediaObject.getInfo();
            
            metadata.put("duration", String.valueOf(info.getDuration()));
            metadata.put("audio.codec", info.getAudio().getDecoder());
            metadata.put("audio.bitRate", String.valueOf(info.getAudio().getBitRate()));
            metadata.put("audio.channels", String.valueOf(info.getAudio().getChannels()));
        } catch (Exception e) {
            log.error("提取音频元数据失败", e);
        }
    }
    
    private void extractBasicMetadata(File file, Map<String, String> metadata) {
        metadata.put("fileName", file.getName());
        metadata.put("fileSize", String.valueOf(file.length()));
        metadata.put("lastModified", new Date(file.lastModified()).toString());
    }
} 