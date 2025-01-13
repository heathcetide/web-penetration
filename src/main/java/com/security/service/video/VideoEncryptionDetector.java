package com.security.service.video;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import java.util.regex.Pattern;
import java.util.regex.Matcher;
import java.util.Map;
import java.util.HashMap;
import cn.hutool.http.HttpUtil;
import cn.hutool.http.HttpResponse;


@Service
public class VideoEncryptionDetector {

    private static Logger log = LoggerFactory.getLogger(VideoEncryptionDetector.class);
    // HLS加密特征
    private static final Pattern HLS_ENCRYPTION = Pattern.compile("#EXT-X-KEY:METHOD=(AES-128|SAMPLE-AES)");
    // DRM特征
    private static final Pattern DRM_PATTERN = Pattern.compile("(widevine|playready|fairplay)");
    // Token参数特征
    private static final Pattern TOKEN_PATTERN = Pattern.compile("[?&](token|auth|sign)=([^&]+)");
    // DASH加密特征
    private static final Pattern DASH_ENCRYPTION = Pattern.compile("(\\.mpd|/dash/)");
    private static final Pattern CUSTOM_ENCRYPTION = Pattern.compile("\\.(m3u8|ts|mp4)\\?key=([^&]+)");
    private static final Pattern TIME_EXPIRY = Pattern.compile("expires?=\\d+");
    
    /**
     * 检测视频URL的加密类型
     */
    public EncryptionInfo detectEncryption(String url) {
        EncryptionInfo info = new EncryptionInfo();
        
        try {
            // 1. 检查URL中的token参数
            checkTokenParameters(url, info);
            
            // 2. 检查响应头
            checkResponseHeaders(url, info);
            
            // 3. 如果是m3u8，检查内容
            if (url.endsWith(".m3u8")) {
                checkM3u8Content(url, info);
            }
            
            // 4. 检查防盗链
            checkRefererRestriction(url, info);
            
            // 5. 检查自定义加密
            checkCustomEncryption(url, info);
            
        } catch (Exception e) {
            log.error("检测视频加密失败: {}", e.getMessage(), e);
            info.setError(e.getMessage());
        }
        
        return info;
    }
    
    private void checkTokenParameters(String url, EncryptionInfo info) {
        Matcher matcher = TOKEN_PATTERN.matcher(url);
        if (matcher.find()) {
            info.setHasToken(true);
            info.setTokenType(matcher.group(1));
            info.setTokenValue(matcher.group(2));
        }
    }
    
    private void checkResponseHeaders(String url, EncryptionInfo info) {
        try {
            Map<String, String> headers = new HashMap<>();
            headers.put("User-Agent", "Mozilla/5.0");
            
            HttpResponse response = HttpUtil.createGet(url)
                .addHeaders(headers)
                .execute();
            
            // 检查DRM相关头
            String drmHeader = response.header("X-DRM-Policy");
            if (drmHeader != null) {
                info.setHasDRM(true);
                info.setDrmType(drmHeader);
            }
            
            // 检查内容加密
            String contentEncryption = response.header("Content-Encryption");
            if (contentEncryption != null) {
                info.setHasContentEncryption(true);
                info.setEncryptionMethod(contentEncryption);
            }
        } catch (Exception e) {
            log.warn("检查响应头失败: {}", e.getMessage());
        }
    }
    
    private void checkM3u8Content(String url, EncryptionInfo info) {
        try {
            String content = HttpUtil.get(url);
            
            // 检查HLS加密
            if (HLS_ENCRYPTION.matcher(content).find()) {
                info.setHasHLSEncryption(true);
                // 提取加密密钥URL
                Pattern keyPattern = Pattern.compile("#EXT-X-KEY:.*URI=\"([^\"]+)\"");
                Matcher matcher = keyPattern.matcher(content);
                if (matcher.find()) {
                    info.setKeyUrl(matcher.group(1));
                }
            }
            
            // 检查分片是否加密
            if (content.contains("#EXT-X-STREAM-INF")) {
                info.setHasMultipleQualities(true);
            }
        } catch (Exception e) {
            log.warn("检查M3U8内容失败: {}", e.getMessage());
        }
    }
    
    private void checkRefererRestriction(String url, EncryptionInfo info) {
        try {
            // 不带Referer请求
            HttpResponse response1 = HttpUtil.createGet(url).execute();
            
            // 带随机Referer请求
            Map<String, String> headers = new HashMap<>();
            headers.put("Referer", "https://example.com");
            HttpResponse response2 = HttpUtil.createGet(url)
                .addHeaders(headers)
                .execute();
            
            // 比较响应是否不同
            if (response1.getStatus() != response2.getStatus()) {
                info.setHasRefererCheck(true);
            }
        } catch (Exception e) {
            log.warn("检查防盗链失败: {}", e.getMessage());
        }
    }
    
    private void checkCustomEncryption(String url, EncryptionInfo info) {
        try {
            // 检查自定义加密参数
            Matcher matcher = CUSTOM_ENCRYPTION.matcher(url);
            if (matcher.find()) {
                info.setHasCustomEncryption(true);
                info.setCustomEncryptionKey(matcher.group(2));
            }
            
            // 检查时间戳参数
            if (TIME_EXPIRY.matcher(url).find()) {
                info.setHasTimeExpiry(true);
                info.setExpiryTime(extractExpiryTime(url));
            }
            
            // 检查DASH流
            if (DASH_ENCRYPTION.matcher(url).find()) {
                info.setHasDASH(true);
                checkDASHEncryption(url, info);
            }
        } catch (Exception e) {
            log.warn("检查自定义加密失败: {}", e.getMessage());
        }
    }
    
    private void checkDASHEncryption(String url, EncryptionInfo info) {
        try {
            String mpd = HttpUtil.get(url);
            // 检查DASH加密方案
            if (mpd.contains("ContentProtection")) {
                if (mpd.contains("urn:mpeg:dash:mp4protection:2011")) {
                    info.setHasCommonEncryption(true);
                }
                if (mpd.contains("urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed")) {
                    info.setHasWidevineDRM(true);
                }
                if (mpd.contains("9a04f079-9840-4286-ab92-e65be0885f95")) {
                    info.setHasPlayReadyDRM(true);
                }
            }
        } catch (Exception e) {
            log.warn("检查DASH加密失败: {}", e.getMessage());
        }
    }
    
    private Long extractExpiryTime(String url) {
        try {
            Pattern pattern = Pattern.compile("expires?=(\\d+)");
            Matcher matcher = pattern.matcher(url);
            if (matcher.find()) {
                return Long.parseLong(matcher.group(1));
            }
        } catch (Exception e) {
            log.warn("提取过期时间失败: {}", e.getMessage());
        }
        return null;
    }
} 