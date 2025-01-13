package com.security.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import java.util.List;
import java.util.ArrayList;

@Configuration
@ConfigurationProperties(prefix = "resource")
public class ResourceConfig {
    private String downloadPath = "download/";
    private String tempPath = "temp/";
    private Integer maxRetries = 3;
    private Integer downloadTimeout = 30000;
    private Integer maxConcurrent = 5;
    
    // 代理配置
    private ProxyConfig proxy = new ProxyConfig();
    // 过滤规则
    private FilterConfig filter = new FilterConfig();
    // 视频配置
    private VideoConfig video = new VideoConfig();
    
    private DownloadConfig download = new DownloadConfig();
    private IntegrityConfig integrity = new IntegrityConfig();

    public static class ProxyConfig {
        private boolean enabled = false;
        private String host;
        private Integer port;
        private String username;
        private String password;

        public boolean isEnabled() {
            return enabled;
        }

        public void setEnabled(boolean enabled) {
            this.enabled = enabled;
        }

        public String getHost() {
            return host;
        }

        public void setHost(String host) {
            this.host = host;
        }

        public Integer getPort() {
            return port;
        }

        public void setPort(Integer port) {
            this.port = port;
        }

        public String getUsername() {
            return username;
        }

        public void setUsername(String username) {
            this.username = username;
        }

        public String getPassword() {
            return password;
        }

        public void setPassword(String password) {
            this.password = password;
        }
    }

    public static class FilterConfig {
        private List<String> allowedDomains = new ArrayList<>();
        private List<String> allowedExtensions = new ArrayList<>();
        private List<String> excludedUrls = new ArrayList<>();
        private Long maxFileSize = 100 * 1024 * 1024L; // 默认100MB

        public List<String> getAllowedDomains() {
            return allowedDomains;
        }

        public void setAllowedDomains(List<String> allowedDomains) {
            this.allowedDomains = allowedDomains;
        }

        public List<String> getAllowedExtensions() {
            return allowedExtensions;
        }

        public void setAllowedExtensions(List<String> allowedExtensions) {
            this.allowedExtensions = allowedExtensions;
        }

        public List<String> getExcludedUrls() {
            return excludedUrls;
        }

        public void setExcludedUrls(List<String> excludedUrls) {
            this.excludedUrls = excludedUrls;
        }

        public Long getMaxFileSize() {
            return maxFileSize;
        }

        public void setMaxFileSize(Long maxFileSize) {
            this.maxFileSize = maxFileSize;
        }
    }

    public static class DownloadConfig {
        private Integer queueSize = 1000;
        private Integer threadPoolSize = 5;
        private Long speedLimit = 1024 * 1024L; // 默认1MB/s
        private Boolean enableResume = true;

        public Integer getQueueSize() {
            return queueSize;
        }

        public void setQueueSize(Integer queueSize) {
            this.queueSize = queueSize;
        }

        public Integer getThreadPoolSize() {
            return threadPoolSize;
        }

        public void setThreadPoolSize(Integer threadPoolSize) {
            this.threadPoolSize = threadPoolSize;
        }

        public Long getSpeedLimit() {
            return speedLimit;
        }

        public void setSpeedLimit(Long speedLimit) {
            this.speedLimit = speedLimit;
        }

        public Boolean getEnableResume() {
            return enableResume;
        }

        public void setEnableResume(Boolean enableResume) {
            this.enableResume = enableResume;
        }
    }

    public static class IntegrityConfig {
        private Boolean enableVerification = true;
        private String algorithm = "MD5"; // MD5, SHA-1, SHA-256
        private Boolean strictMode = false;

        public Boolean getEnableVerification() {
            return enableVerification;
        }

        public void setEnableVerification(Boolean enableVerification) {
            this.enableVerification = enableVerification;
        }

        public String getAlgorithm() {
            return algorithm;
        }

        public void setAlgorithm(String algorithm) {
            this.algorithm = algorithm;
        }

        public Boolean getStrictMode() {
            return strictMode;
        }

        public void setStrictMode(Boolean strictMode) {
            this.strictMode = strictMode;
        }
    }

    public String getDownloadPath() {
        return downloadPath;
    }

    public void setDownloadPath(String downloadPath) {
        this.downloadPath = downloadPath;
    }

    public String getTempPath() {
        return tempPath;
    }

    public void setTempPath(String tempPath) {
        this.tempPath = tempPath;
    }

    public Integer getMaxRetries() {
        return maxRetries;
    }

    public void setMaxRetries(Integer maxRetries) {
        this.maxRetries = maxRetries;
    }

    public Integer getDownloadTimeout() {
        return downloadTimeout;
    }

    public void setDownloadTimeout(Integer downloadTimeout) {
        this.downloadTimeout = downloadTimeout;
    }

    public Integer getMaxConcurrent() {
        return maxConcurrent;
    }

    public void setMaxConcurrent(Integer maxConcurrent) {
        this.maxConcurrent = maxConcurrent;
    }

    public ProxyConfig getProxy() {
        return proxy;
    }

    public void setProxy(ProxyConfig proxy) {
        this.proxy = proxy;
    }

    public FilterConfig getFilter() {
        return filter;
    }

    public void setFilter(FilterConfig filter) {
        this.filter = filter;
    }

    public DownloadConfig getDownload() {
        return download;
    }

    public void setDownload(DownloadConfig download) {
        this.download = download;
    }

    public IntegrityConfig getIntegrity() {
        return integrity;
    }

    public void setIntegrity(IntegrityConfig integrity) {
        this.integrity = integrity;
    }

    public VideoConfig getVideo() {
        return video;
    }

    public void setVideo(VideoConfig video) {
        this.video = video;
    }
}