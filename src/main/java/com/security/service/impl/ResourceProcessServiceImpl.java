package com.security.service.impl;

import cn.hutool.core.io.FileUtil;
import cn.hutool.core.util.StrUtil;
import cn.hutool.http.HttpUtil;
import com.security.model.entity.ResourceDownload;
import com.security.mapper.ResourceDownloadMapper;
import com.security.service.IResourceProcessService;
import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.select.Elements;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import cn.hutool.http.HttpRequest;
import cn.hutool.http.HttpResponse;
//import cn.hutool.http.HttpProgress;
import com.security.config.ResourceConfig;
import com.security.service.download.DownloadProgressListener;
import com.security.util.RetryUtil;
import com.security.service.video.VideoProcessor;
import com.security.service.video.WatermarkPosition;
import com.security.service.video.VideoEncryptionDetector;
import com.security.service.video.EncryptionInfo;
import cn.hutool.json.JSONUtil;

public class ResourceProcessServiceImpl implements IResourceProcessService {
    
    @Autowired
    private ResourceDownloadMapper resourceDownloadMapper;
    
    @Autowired
    private ResourceConfig resourceConfig;
    
    @Autowired
    private VideoProcessor videoProcessor;
    
    @Autowired
    private VideoEncryptionDetector encryptionDetector;

    private static Logger log = LoggerFactory.getLogger(ResourceProcessServiceImpl.class);
    
    // 视频文件后缀正则
    private static final Pattern VIDEO_PATTERN = Pattern.compile("\\.(mp4|flv|m3u8|ts)$", Pattern.CASE_INSENSITIVE);
    // 图片文件后缀正则
    private static final Pattern IMAGE_PATTERN = Pattern.compile("\\.(jpg|jpeg|png|gif|webp)$", Pattern.CASE_INSENSITIVE);
    // 常见文件后缀正则
    private static final Pattern FILE_PATTERN = Pattern.compile("\\.(pdf|doc|docx|xls|xlsx|zip|rar)$", Pattern.CASE_INSENSITIVE);
    // 音频文件后缀正则
    private static final Pattern AUDIO_PATTERN = Pattern.compile("\\.(mp3|wav|ogg|m4a)$", Pattern.CASE_INSENSITIVE);
    // 字幕文件后缀正则
    private static final Pattern SUBTITLE_PATTERN = Pattern.compile("\\.(srt|ass|vtt)$", Pattern.CASE_INSENSITIVE);
    // Ebook文件后缀正则
    private static final Pattern EBOOK_PATTERN = Pattern.compile("\\.(epub|mobi|azw3)$", Pattern.CASE_INSENSITIVE);
    
    @Override
    public List<String> extractImages(String html, String baseUrl) {
        List<String> images = new ArrayList<>();
        Document doc = Jsoup.parse(html, baseUrl);
        
        // 提取<img>标签的图片
        Elements imgElements = doc.select("img[src]");
        imgElements.forEach(img -> {
            String src = img.attr("abs:src");
            if (IMAGE_PATTERN.matcher(src).find()) {
                images.add(src);
            }
        });
        
        // 提取背景图片
        Elements bgElements = doc.select("[style*=background-image]");
        bgElements.forEach(element -> {
            String style = element.attr("style");
            String url = extractUrlFromStyle(style);
            if (url != null && IMAGE_PATTERN.matcher(url).find()) {
                images.add(url);
            }
        });
        
        return images;
    }
    
    @Override
    public List<String> extractVideos(String html, String baseUrl) {
        List<String> videos = new ArrayList<>();
        Document doc = Jsoup.parse(html, baseUrl);
        
        // 提取视频资源并检查加密
        Elements videoElements = doc.select("video[src], source[src]");
        videoElements.forEach(video -> {
            String src = video.attr("abs:src");
            if (VIDEO_PATTERN.matcher(src).find()) {
                // 检查视频加密
                EncryptionInfo encryptionInfo = encryptionDetector.detectEncryption(src);
                if (encryptionInfo.isHasToken() || encryptionInfo.isHasDRM() || 
                    encryptionInfo.isHasHLSEncryption()) {
                    log.info("发现加密视频: {}, 加密信息: {}", src, JSONUtil.toJsonStr(encryptionInfo));
                    // 处理加密视频
                    handleEncryptedVideo(src, encryptionInfo);
                } else {
                    videos.add(src);
                }
            }
        });
        
        return videos;
    }
    
    @Override
    public void downloadResource(Long taskId, String url, String type) {
        // 检查过滤规则
        if (!isAllowedResource(url)) {
            log.warn("资源被过滤: {}", url);
            return;
        }
        
        try {
            String fileName = getFileNameFromUrl(url);
            String savePath = getSavePath(taskId, type, fileName);
            
            // 创建目录
            FileUtil.mkdir(new File(savePath).getParent());
            
            // 使用重试机制下载文件
            RetryUtil.retry(() -> {
                // 创建下载请求
                HttpRequest request = createHttpRequest(url);
                // 添加进度监听
                DownloadProgressListener progressListener = new DownloadProgressListener(taskId, url);
                request.header("Accept-Encoding", "identity");  // 禁用压缩以获取准确的进度
                HttpResponse response = request.execute();
                if (response.isOk()) {
                    long total = response.contentLength();
                    long current = 0;
                    byte[] buffer = new byte[8192];
                    try (FileOutputStream out = new FileOutputStream(new File(savePath))) {
                        InputStream in = response.bodyStream();
                        int len;
                        while ((len = in.read(buffer)) != -1) {
                            out.write(buffer, 0, len);
                            current += len;
                            progressListener.progress(total, current);
                        }
                    } catch (IOException e) {
                        throw new RuntimeException(e);
                    }
                }
                return null;
            }, resourceConfig.getMaxRetries(), 1000, "下载文件");
            
            // 更新下载记录
            updateDownloadStatus(taskId, url, savePath, 2);
        } catch (Exception e) {
            updateDownloadStatus(taskId, url, null, 3, e.getMessage());
        }
    }
    
    private HttpRequest createHttpRequest(String url) {
        HttpRequest request = HttpUtil.createGet(url);
        
        // 设置超时
        request.timeout(resourceConfig.getDownloadTimeout());
        
        // 配置代理
        if (resourceConfig.getProxy().isEnabled()) {
            ResourceConfig.ProxyConfig proxy = resourceConfig.getProxy();
            request.setHttpProxy(proxy.getHost(), proxy.getPort());
            if (StrUtil.isNotEmpty(proxy.getUsername())) {
                request.basicAuth(proxy.getUsername(), proxy.getPassword());
            }
        }
        
        return request;
    }
    
    private boolean isAllowedResource(String url) {
        ResourceConfig.FilterConfig filter = resourceConfig.getFilter();
        
        // 检查域名白名单
        if (!filter.getAllowedDomains().isEmpty()) {
            boolean domainAllowed = filter.getAllowedDomains().stream()
                    .anyMatch(domain -> url.contains(domain));
            if (!domainAllowed) return false;
        }
        
        // 检查文件扩展名
        if (!filter.getAllowedExtensions().isEmpty()) {
            boolean extensionAllowed = filter.getAllowedExtensions().stream()
                    .anyMatch(ext -> url.toLowerCase().endsWith(ext.toLowerCase()));
            if (!extensionAllowed) return false;
        }
        
        // 检查排除URL
        if (filter.getExcludedUrls().stream().anyMatch(url::contains)) {
            return false;
        }
        
        return true;
    }
    
    @Override
    public void mergeVideoSegments(Long taskId, List<String> segments, String outputPath) {
        try {
            // 创建临时目录
            String tempDir = "temp/" + taskId + "/";
            FileUtil.mkdir(tempDir);
            
            // 下载所有片段
            List<File> files = new ArrayList<>();
            for (int i = 0; i < segments.size(); i++) {
                String segmentUrl = segments.get(i);
                String segmentPath = tempDir + "segment_" + i + ".ts";
                HttpUtil.downloadFile(segmentUrl, segmentPath);
                files.add(new File(segmentPath));
            }
            
            // 合并文件
            FileUtil.mkdir(new File(outputPath).getParent());
            mergeFiles(files, outputPath);
            
            // 清理临时文件
            FileUtil.del(tempDir);
        } catch (Exception e) {
            throw new RuntimeException("合并视频片段失败：" + e.getMessage());
        }
    }
    
    @Override
    public List<String> extractFiles(String html, String baseUrl) {
        List<String> files = new ArrayList<>();
        Document doc = Jsoup.parse(html, baseUrl);
        
        // 提取<a>标签的文件链接
        Elements linkElements = doc.select("a[href]");
        linkElements.forEach(link -> {
            String href = link.attr("abs:href");
            if (FILE_PATTERN.matcher(href).find()) {
                files.add(href);
            }
        });
        
        return files;
    }
    
    @Override
    public List<String> extractComics(String html, String baseUrl) {
        List<String> comics = new ArrayList<>();
        Document doc = Jsoup.parse(html, baseUrl);
        
        // 漫画网站特定的图片选择器
        Elements comicElements = doc.select(".comic-image img, .manga-image img, .chapter-image img");
        comicElements.forEach(img -> {
            String src = img.attr("abs:src");
            if (IMAGE_PATTERN.matcher(src).find()) {
                comics.add(src);
            }
        });
        
        return comics;
    }
    
    private String extractUrlFromStyle(String style) {
        if (style == null) return null;
        // 匹配url("...")或url('...')或url(...)格式
        Pattern pattern = Pattern.compile("url\\(['\"]?(.*?)['\"]?\\)");
        java.util.regex.Matcher matcher = pattern.matcher(style);
        if (matcher.find()) {
            return matcher.group(1);
        }
        return null;
    }
    
    private String getFileNameFromUrl(String url) {
        return url.substring(url.lastIndexOf("/") + 1);
    }
    
    private String getSavePath(Long taskId, String type, String fileName) {
        return "download/" + taskId + "/" + type.toLowerCase() + "/" + fileName;
    }
    
    private void updateDownloadStatus(Long taskId, String url, String localPath, int status) {
        updateDownloadStatus(taskId, url, localPath, status, null);
    }
    
    private void updateDownloadStatus(Long taskId, String url, String localPath, int status, String errorMsg) {
        ResourceDownload download = new ResourceDownload();
        download.setTaskId(taskId);
        download.setResourceUrl(url);
        download.setLocalPath(localPath);
        download.setStatus(status);
        download.setErrorMsg(errorMsg);
        download.setCreateTime(new Date());
        download.setUpdateTime(new Date());
        
        resourceDownloadMapper.insert(download);
    }
    
    private void mergeFiles(List<File> files, String outputPath) throws Exception {
        try (FileOutputStream fos = new FileOutputStream(outputPath)) {
            for (File file : files) {
                byte[] bytes = FileUtil.readBytes(file);
                fos.write(bytes);
            }
            fos.flush();
        }
    }
    
    // 添加M3U8视频处理方法
    private void processM3u8Video(Long taskId, String m3u8Url) {
        try {
            // 下载m3u8文件
            String m3u8Content = HttpUtil.get(m3u8Url);
            List<String> segments = parseM3u8(m3u8Content);
            
            // 获取基础URL
            String baseUrl = m3u8Url.substring(0, m3u8Url.lastIndexOf("/") + 1);
            
            // 补全片段URL
            segments = segments.stream()
                    .map(segment -> segment.startsWith("http") ? segment : baseUrl + segment)
                    .collect(Collectors.toList());
            
            // 合并视频片段
            String outputPath = getSavePath(taskId, "VIDEO", "output.mp4");
            mergeVideoSegments(taskId, segments, outputPath);
            
        } catch (Exception e) {
            throw new RuntimeException("处理M3U8视频失败：" + e.getMessage());
        }
    }
    
    private List<String> parseM3u8(String content) {
        List<String> segments = new ArrayList<>();
        String[] lines = content.split("\n");
        for (String line : lines) {
            if (!line.startsWith("#") && line.trim().length() > 0) {
                segments.add(line.trim());
            }
        }
        return segments;
    }
    
    @Override
    public List<String> extractAudios(String html, String baseUrl) {
        List<String> audios = new ArrayList<>();
        Document doc = Jsoup.parse(html, baseUrl);
        
        // 提取音频文件
        Elements audioElements = doc.select("audio[src], source[src]");
        audioElements.forEach(audio -> {
            String src = audio.attr("abs:src");
            if (AUDIO_PATTERN.matcher(src).find()) {
                audios.add(src);
            }
        });
        
        return audios;
    }
    
    @Override
    public List<String> extractSubtitles(String html, String baseUrl) {
        List<String> subtitles = new ArrayList<>();
        Document doc = Jsoup.parse(html, baseUrl);
        
        // 提取字幕文件
        Elements trackElements = doc.select("track[src]");
        trackElements.forEach(track -> {
            String src = track.attr("abs:src");
            if (SUBTITLE_PATTERN.matcher(src).find()) {
                subtitles.add(src);
            }
        });
        
        return subtitles;
    }
    
    @Override
    public void handleEncryptedVideo(String url, EncryptionInfo encryptionInfo) {
        try {
            if (encryptionInfo.isHasToken()) {
                // 处理带token的视频
                handleTokenVideo(url, encryptionInfo.getTokenValue());
            }
            
            if (encryptionInfo.isHasHLSEncryption()) {
                // 处理HLS加密视频
                handleHLSEncryptedVideo(url, encryptionInfo.getKeyUrl());
            }
            
            if (encryptionInfo.isHasDRM()) {
                // 记录DRM视频信息
                log.warn("发现DRM保护视频，暂不支持下载: {}", url);
            }
            
        } catch (Exception e) {
            log.error("处理加密视频失败: {}", e.getMessage(), e);
        }
    }
    
    @Override
    public void handleTokenVideo(String url, String token) {
        // 1. 验证token有效期
        if (isTokenExpired(token)) {
            // 2. 尝试刷新token
            token = refreshToken(url);
        }
        
        // 3. 使用新token下载
        String finalUrl = appendToken(url, token);
        downloadVideo(finalUrl);
    }
    
    public void handleHLSEncryptedVideo(String m3u8Url, String keyUrl) {
        try {
            // 1. 下载密钥
            byte[] key = downloadKey(keyUrl);
            
            // 2. 下载并解密m3u8文件
            String m3u8Content = decryptM3u8(m3u8Url, key);
            
            // 3. 解析视频片段
            List<String> segments = parseM3u8(m3u8Content);
            
            // 4. 下载并解密片段
            for (String segment : segments) {
                byte[] encryptedData = downloadSegment(segment);
                byte[] decryptedData = decryptSegment(encryptedData, key);
                saveSegment(decryptedData);
            }
            
            // 5. 合并片段
            mergeSegments();
            
        } catch (Exception e) {
            log.error("处理HLS加密视频失败: {}", e.getMessage(), e);
        }
    }
    
    private boolean isTokenExpired(String token) {
        try {
            // 实现token过期检查逻辑
            return false;
        } catch (Exception e) {
            log.error("检查token过期失败", e);
            return true;
        }
    }
    
    private String refreshToken(String url) {
        try {
            // 实现token刷新逻辑
            return "";
        } catch (Exception e) {
            log.error("刷新token失败", e);
            return null;
        }
    }
    
    private String appendToken(String url, String token) {
        try {
            // 实现添加token到url的逻辑
            return url + "?token=" + token;
        } catch (Exception e) {
            log.error("添加token失败", e);
            return url;
        }
    }
    
    private void downloadVideo(String url) {
        try {
            // 实现视频下载逻辑
        } catch (Exception e) {
            log.error("下载视频失败", e);
        }
    }
    
    private byte[] downloadKey(String keyUrl) {
        try {
            // 实现密钥下载逻辑
            return new byte[0];
        } catch (Exception e) {
            log.error("下载密钥失败", e);
            return null;
        }
    }
    
    private String decryptM3u8(String content, byte[] key) {
        try {
            // 实现M3U8解密逻辑
            return "";
        } catch (Exception e) {
            log.error("解密M3U8失败", e);
            return null;
        }
    }
    
    private byte[] downloadSegment(String url) {
        try {
            // 实现分片下载逻辑
            return new byte[0];
        } catch (Exception e) {
            log.error("下载分片失败", e);
            return null;
        }
    }
    
    private byte[] decryptSegment(byte[] data, byte[] key) {
        try {
            // 实现分片解密逻辑
            return new byte[0];
        } catch (Exception e) {
            log.error("解密分片失败", e);
            return null;
        }
    }
    
    private void saveSegment(byte[] data) {
        try {
            // 实现分片保存逻辑
        } catch (Exception e) {
            log.error("保存分片失败", e);
        }
    }
    
    private void mergeSegments() {
        try {
            // 实现分片合并逻辑
        } catch (Exception e) {
            log.error("合并分片失败", e);
        }
    }
    
    @Override
    public void postProcessVideo(File videoFile) {
        // 1. 检查视频完整性
        if (!videoProcessor.checkIntegrity(videoFile)) {
            log.error("视频文件损坏: {}", videoFile.getPath());
            return;
        }
        
        // 2. 压缩视频（目标大小100MB）
        videoProcessor.compress(videoFile, 100);
        
        // 3. 添加水印
        File watermark = new File("watermark.png");
        if (watermark.exists()) {
            videoProcessor.addWatermark(videoFile, watermark, WatermarkPosition.BOTTOM_RIGHT);
        }
        
        // 4. 转换格式（如果需要）
        String extension = FileUtil.extName(videoFile);
        if (!"mp4".equalsIgnoreCase(extension)) {
            videoProcessor.convert(videoFile, "mp4");
        }
    }
} 