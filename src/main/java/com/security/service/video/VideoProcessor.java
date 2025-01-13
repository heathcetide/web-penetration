package com.security.service.video;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import ws.schild.jave.Encoder;
import ws.schild.jave.MultimediaObject;
import ws.schild.jave.encode.AudioAttributes;
import ws.schild.jave.encode.EncodingAttributes;
import ws.schild.jave.encode.VideoAttributes;
import ws.schild.jave.info.VideoSize;
import ws.schild.jave.filters.VideoFilter;

import java.io.File;
import java.util.ArrayList;
import java.util.List;

@Service
public class VideoProcessor {

    private static Logger log = LoggerFactory.getLogger(VideoProcessor.class);
    /**
     * 视频格式转换
     */
    public void convert(File source, String targetFormat) {
        try {
            File target = new File(source.getParent(), 
                    source.getName().replaceFirst("[.][^.]+$", "." + targetFormat));
            
            AudioAttributes audio = new AudioAttributes();
            audio.setCodec("aac");
            audio.setBitRate(128000);
            audio.setChannels(2);
            audio.setSamplingRate(44100);
            
            VideoAttributes video = new VideoAttributes();
            video.setCodec("h264");
            video.setBitRate(1000000);
            video.setFrameRate(30);
            
            EncodingAttributes attrs = new EncodingAttributes();
            attrs.setOutputFormat("mp4");
            attrs.setAudioAttributes(audio);
            attrs.setVideoAttributes(video);
            
            Encoder encoder = new Encoder();
            encoder.encode(new MultimediaObject(source), target, attrs);
            
        } catch (Exception e) {
            log.error("视频格式转换失败: {}", e.getMessage(), e);
        }
    }
    
    /**
     * 视频压缩
     */
    public void compress(File source, int targetSize) {
        try {
            File target = new File(source.getParent(), "compressed_" + source.getName());
            
            // 计算目标比特率
            MultimediaObject mediaObject = new MultimediaObject(source);
            long duration = mediaObject.getInfo().getDuration();
            int bitRate = (int) (targetSize * 8 * 1024 * 1024 / (duration / 1000));
            
            VideoAttributes video = new VideoAttributes();
            video.setCodec("h264");
            video.setBitRate(bitRate);
            
            AudioAttributes audio = new AudioAttributes();
            audio.setCodec("aac");
            audio.setBitRate(128000);
            
            EncodingAttributes attrs = new EncodingAttributes();
            attrs.setOutputFormat("mp4");
            attrs.setAudioAttributes(audio);
            attrs.setVideoAttributes(video);
            
            new Encoder().encode(mediaObject, target, attrs);
            
        } catch (Exception e) {
            log.error("视频压缩失败: {}", e.getMessage(), e);
        }
    }
    
    /**
     * 视频截图
     */
    public List<File> takeScreenshots(File source, int count) {
        List<File> screenshots = new ArrayList<>();
        try {
            MultimediaObject mediaObject = new MultimediaObject(source);
            long duration = mediaObject.getInfo().getDuration();
            long interval = duration / (count + 1);
            
            for (int i = 1; i <= count; i++) {
                File output = new File(source.getParent(), 
                        source.getName() + "_screenshot_" + i + ".jpg");
                
                VideoAttributes video = new VideoAttributes();
                video.setCodec("mjpeg");
                video.setSize(new VideoSize(1280, 720));
                
                EncodingAttributes attrs = new EncodingAttributes();
                attrs.setOutputFormat("image2");
                attrs.setOffset(Float.valueOf(interval * i));
                attrs.setVideoAttributes(video);
                
                new Encoder().encode(mediaObject, output, attrs);
                screenshots.add(output);
            }
        } catch (Exception e) {
            log.error("视频截图失败: {}", e.getMessage(), e);
        }
        return screenshots;
    }
    
    /**
     * 添加水印
     */
    public void addWatermark(File source, File watermark, WatermarkPosition position) {
        try {
            File target = new File(source.getParent(), "watermarked_" + source.getName());
            
            VideoAttributes video = new VideoAttributes();
            video.setCodec("h264");
            video.addFilter(new VideoFilter() {
                @Override
                public String getExpression() {
                    return "movie=" + watermark.getAbsolutePath() + " [watermark]; " +
                            "[in][watermark] overlay=" + position.getCoordinates() + " [out]";
                }
            });
            
            AudioAttributes audio = new AudioAttributes();
            audio.setCodec("copy");
            
            EncodingAttributes attrs = new EncodingAttributes();
            attrs.setOutputFormat("mp4");
            attrs.setAudioAttributes(audio);
            attrs.setVideoAttributes(video);
            
            new Encoder().encode(new MultimediaObject(source), target, attrs);
            
        } catch (Exception e) {
            log.error("添加水印失败: {}", e.getMessage(), e);
        }
    }
    
    /**
     * 检查视频完整性
     */
    public boolean checkIntegrity(File videoFile) {
        try {
            MultimediaObject mediaObject = new MultimediaObject(videoFile);
            // 尝试获取视频信息，如果能成功获取说明文件基本完整
            mediaObject.getInfo();
            return true;
        } catch (Exception e) {
            log.error("视频完整性检查失败: {}", e.getMessage());
            return false;
        }
    }
} 