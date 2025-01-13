package com.security.service.download;

import cn.hutool.http.HttpRequest;
import com.security.exception.BusinessException;
import com.security.model.entity.SegmentStatus;
import com.security.service.video.VideoParserManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import com.google.common.util.concurrent.RateLimiter;
import com.security.config.ResourceConfig;
import com.security.model.entity.DownloadProgress;
import cn.hutool.http.HttpUtil;
import cn.hutool.core.io.FileUtil;
import cn.hutool.core.io.StreamProgress;

import java.io.*;
import java.util.*;
import java.util.concurrent.*;
import java.net.URI;

import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.WebDriver;
import cn.hutool.http.HttpResponse;
import com.security.model.entity.DownloadStats;

import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;
import com.security.service.video.parser.VideoInfo;
import com.security.service.video.parser.VideoSource;

@Component
public class DownloadQueueManager {
    
    private static final Logger log = LoggerFactory.getLogger(DownloadQueueManager.class);
    
    @Autowired
    private WebDriver webDriver;
    
    private final BlockingQueue<DownloadTask> downloadQueue;
    private final ExecutorService downloadExecutor;
    private final RateLimiter rateLimiter;
    
    @Autowired
    private ResourceConfig resourceConfig;
    
    private final Map<String, DownloadStats> downloadStats = new ConcurrentHashMap<>();
    private final int MIN_THREADS = 2;
    private final int MAX_THREADS = 10;
    private volatile int currentThreads = 3;
    
    @Autowired
    private VideoParserManager videoParserManager;
    
    public DownloadQueueManager(ResourceConfig config) {
        this.downloadQueue = new LinkedBlockingQueue<>(config.getDownload().getQueueSize());
        this.downloadExecutor = Executors.newFixedThreadPool(config.getDownload().getThreadPoolSize());
        this.rateLimiter = RateLimiter.create(config.getDownload().getSpeedLimit());
        
        // 启动下载处理线程
        startDownloadProcessor();
    }
    
    private void startDownloadProcessor() {
        downloadExecutor.submit(() -> {
            while (!Thread.currentThread().isInterrupted()) {
                try {
                    DownloadTask task = downloadQueue.take();
                    processDownloadTask(task);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    break;
                } catch (Exception e) {
                    log.error("处理下载任务失败", e);
                }
            }
        });
    }
    
    private void processDownloadTask(DownloadTask task) {
        try {
            String url = task.getUrl();
            if (url == null) {
                throw new BusinessException("下载URL不能为空");
            }
            
            // 根据URL类型选择不同的下载方式
            if (url.contains("blob:")) {
                processBlobDownload(url, task.getSavePath());  // 修改这里，传入正确的参数
            } else if (url.contains(".m3u8")) {
                processM3u8Download(task);
            } else {
                processNormalDownload(task);
            }
            
        } catch (Exception e) {
            log.error("处理下载任务失败: {}", task.getUrl(), e);
            task.setStatus(3);
            task.setErrorMsg(e.getMessage());
        }
    }
    
    private void processM3u8Download(DownloadTask task) {
        String tempDir = task.getSavePath() + "_temp";
        boolean downloadSuccess = false;
        DownloadStats stats = new DownloadStats();
        stats.setStartTime(System.currentTimeMillis());
        downloadStats.put(task.getTaskId().toString(), stats);
        
        try {
            // 1. 获取M3U8内容
            String m3u8Content = retryDownload(() -> {
                HttpRequest request = HttpRequest.get(task.getUrl())
                    .timeout(30000)
                    .header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124")
                    .header("Referer", getBaseUrl(task.getUrl()));
                return request.execute().body();
            });

            // 2. 解析所有片段
            List<String> segments = parseM3u8(m3u8Content, task.getUrl());
            if (segments.isEmpty()) {
                throw new RuntimeException("未找到视频片段");
            }
            
            // 3. 创建临时目录
            File tempDirFile = new File(tempDir);
            if (!tempDirFile.exists()) {
                tempDirFile.mkdirs();
            }
            
            // 4. 使用多线程下载
            CountDownLatch downloadLatch = new CountDownLatch(segments.size());
            ExecutorService downloadExecutor = Executors.newFixedThreadPool(currentThreads);
            ScheduledExecutorService monitor = Executors.newSingleThreadScheduledExecutor();
            
            // 启动监控
            monitor.scheduleAtFixedRate(() -> 
                adjustThreadCount(stats), 5, 5, TimeUnit.SECONDS);
            
            // 创建分片下载计划
            List<SegmentStatus> downloadPlan = createDownloadPlan(segments, tempDir);
            stats.getSegmentStatus().putAll(downloadPlan.stream()
                .collect(Collectors.toMap(SegmentStatus::getUrl, s -> s)));
            
            // 下载分片
            for (SegmentStatus segment : downloadPlan) {
                downloadExecutor.submit(() -> downloadSegmentWithRetry(segment, stats));
            }
            
            // 等待所有下载完成
            downloadLatch.await(30, TimeUnit.MINUTES); // 设置超时时间
            downloadExecutor.shutdown();
            
            // 验证下载结果
            if (segments.size() != segments.size()) {
                throw new RuntimeException(String.format("部分分片下载失败: 预期%d个，实际下载%d个", 
                    segments.size(), segments.size()));
            }
            
            downloadSuccess = true;
            log.info("所有分片下载完成，共{}个文件", segments.size());
            
            // 5. 尝试合并文件
            try {
                mergeSegments(tempDir, task.getSavePath());
                // 合并成功后删除临时文件
                FileUtil.del(tempDir);
                log.info("视频合并完成并清理临时文件: {}", task.getSavePath());
            } catch (Exception e) {
                log.error("合并失败，保留临时文件: {}", tempDir, e);
                throw new RuntimeException("文件合并失败，临时文件保留在: " + tempDir, e);
            }
            
        } catch (Exception e) {
            String errorMsg = downloadSuccess ? 
                "分片下载完成但合并失败，临时文件保留在: " + tempDir : 
                "下载过程失败: " + e.getMessage();
            log.error("M3U8处理失败: {}", errorMsg, e);
            updateDownloadProgress(task, -1L, errorMsg);
            
            // 如果下载都没成功，清理临时目录
            if (!downloadSuccess && tempDir != null) {
                try {
                    FileUtil.del(tempDir);
                } catch (Exception ex) {
                    log.warn("清理临时目录失败: {}", tempDir, ex);
                }
            }
        } finally {
            downloadStats.remove(task.getTaskId());
        }
    }
    
    private void processM3u8Download(String m3u8Url, String savePath) {
        try {
            log.info("开始下载m3u8视频: {}", m3u8Url);
            
            // 1. 下载m3u8文件内容
            String m3u8Content = HttpUtil.get(m3u8Url);
            if (m3u8Content == null || m3u8Content.trim().isEmpty()) {
                throw new BusinessException("无法获取m3u8内容");
            }
            
            // 2. 解析m3u8获取视频片段
            List<String> segments = parseM3u8(m3u8Content, m3u8Url);
            if (segments.isEmpty()) {
                throw new BusinessException("未找到视频片段");
            }
            
            // 3. 创建临时目录存放ts文件
            String tempDir = savePath + "_segments";
            FileUtil.mkdir(tempDir);
            
            try {
                // 4. 下载所有片段
                int totalSegments = segments.size();
                for (int i = 0; i < segments.size(); i++) {
                    String segmentUrl = segments.get(i);
                    // 处理相对路径
                    if (!segmentUrl.startsWith("http")) {
                        segmentUrl = getBaseUrl(m3u8Url) + segmentUrl;
                    }
                    
                    String segmentPath = tempDir + File.separator + String.format("%05d", i) + ".ts";
                    
                    // 下载片段
                    try {
                        HttpRequest request = HttpRequest.get(segmentUrl)
                            .header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124")
                            .header("Referer", webDriver.getCurrentUrl())
                            .timeout(30000);
                            
                        HttpResponse response = request.execute();
                        if (response.isOk()) {
                            FileUtil.writeBytes(response.bodyBytes(), segmentPath);
                            // 更新进度
                            double progress = (i + 1.0) / totalSegments * 100;
                            log.info("下载进度: {}/{} ({:.2f}%)", i + 1, totalSegments, progress);
                        } else {
                            throw new BusinessException("片段下载失败，HTTP状态码: " + response.getStatus());
                        }
                    } catch (Exception e) {
                        log.error("下载片段失败: " + segmentUrl, e);
                        throw new BusinessException("片段下载失败: " + e.getMessage());
                    }
                }
                
                // 5. 合并视频片段
                mergeSegments(tempDir, savePath);
                
                // 6. 清理临时文件
                FileUtil.del(tempDir);
                FileUtil.del(savePath + ".temp.ts");
                
                log.info("M3U8视频下载完成: {}", savePath);
                
            } catch (Exception e) {
                // 清理临时文件
                FileUtil.del(tempDir);
                FileUtil.del(savePath + ".temp.ts");
                throw e;
            }
            
        } catch (Exception e) {
            log.error("M3U8视频下载失败", e);
            throw new BusinessException("M3U8视频下载失败: " + e.getMessage());
        }
    }
    
    private void downloadSegmentWithRetry(SegmentStatus segment, DownloadStats stats) {
        File segmentFile = new File(segment.getLocalPath());
        long startPos = segmentFile.exists() ? segmentFile.length() : 0;
        
        try {
            HttpRequest request = HttpRequest.get(segment.getUrl())
                .header("Range", "bytes=" + startPos + "-")
                .timeout(30000);
            
            try (RandomAccessFile file = new RandomAccessFile(segmentFile, "rw")) {
                file.seek(startPos);
                
                HttpResponse response = request.execute();
                if (!response.isOk()) {
                    throw new RuntimeException("下载失败: " + response.getStatus());
                }
                
                byte[] buffer = new byte[8192];
                try (InputStream is = response.bodyStream()) {
                    int len;
                    while ((len = is.read(buffer)) != -1) {
                        file.write(buffer, 0, len);
                        updateDownloadStats(stats, len);
                    }
                }
            }
            
            segment.setStatus("COMPLETED");
            
        } catch (Exception e) {
            segment.setRetryCount(segment.getRetryCount() + 1);
            if (segment.getRetryCount() < 3) {
                downloadSegmentWithRetry(segment, stats);
            } else {
                segment.setStatus("FAILED");
                log.error("分片下载失败: {}", segment.getUrl(), e);
            }
        }
    }
    
    private void updateDownloadStats(DownloadStats stats, int bytes) {
        long now = System.currentTimeMillis();
        synchronized (stats) {
            stats.setBytesDownloaded(stats.getBytesDownloaded() + bytes);
            long timeDiff = now - stats.getLastUpdateTime();
            if (timeDiff >= 1000) {
                double speed = (double) stats.getBytesDownloaded() / 
                    ((now - stats.getStartTime()) / 1000.0);
                stats.setCurrentSpeed(speed);
                stats.setLastUpdateTime(now);
            }
        }
    }
    
    private void adjustThreadCount(DownloadStats stats) {
        double currentSpeed = stats.getCurrentSpeed();
        int activeThreads = stats.getActiveThreads();
        
        // 根据下载速度调整线程数
        if (currentSpeed < 500_000 && activeThreads > MIN_THREADS) { // 500KB/s
            currentThreads = Math.max(MIN_THREADS, currentThreads - 1);
        } else if (currentSpeed > 2_000_000 && activeThreads < MAX_THREADS) { // 2MB/s
            currentThreads = Math.min(MAX_THREADS, currentThreads + 1);
        }
        
        log.info("下载速度: {}/s, 当前线程数: {}", 
            FileUtil.readableFileSize(Math.round(currentSpeed)), 
            currentThreads);
    }
    
    private void processNormalDownload(DownloadTask task) {
        // 创建下载进度记录
        DownloadProgress progress = new DownloadProgress();
        progress.setTaskId(task.getTaskId());
        progress.setUrl(task.getUrl());
        
        // 下载文件
        File saveFile = new File(task.getSavePath());
        FileUtil.mkParentDirs(saveFile);
        
        HttpUtil.downloadFile(task.getUrl(), saveFile, new StreamProgress() {
            @Override
            public void start() {
                log.info("开始下载: {}", task.getUrl());
            }
            
            @Override
            public void progress(long total, long current) {
                rateLimiter.acquire((int) current);
                updateDownloadProgress(task, current);
            }
            
            @Override
            public void finish() {
                log.info("下载完成: {}", task.getUrl());
            }
        });
    }
    
    private void updateDownloadProgress(DownloadTask task, long progress) {
        updateDownloadProgress(task, progress, null);
    }
    
    private void updateDownloadProgress(DownloadTask task, long progress, String error) {
        try {
            // 更新下载进度
            task.setDownloadedSize(progress);
            if (error != null) {
                task.setStatus(3); // 失败状态
                task.setErrorMsg(error);
            }
            log.info("下载进度: taskId={}, progress={}", task.getTaskId(), progress);
        } catch (Exception e) {
            log.error("更新下载进度失败", e);
        }
    }
    
    private List<String> parseM3u8(String m3u8Content, String baseUrl) {
        List<String> segments = new ArrayList<>();
        String[] lines = m3u8Content.split("\n");
        String masterPlaylistUrl = null;
        
        // 1. 检查是否是master playlist
        for (String line : lines) {
            line = line.trim();
            if (line.startsWith("#EXT-X-STREAM-INF:")) {
                // 获取下一行的URL
                int index = Arrays.asList(lines).indexOf(line);
                if (index + 1 < lines.length) {
                    masterPlaylistUrl = lines[index + 1].trim();
                    break;
                }
            }
        }
        
        // 2. 如果是master playlist，下载实际的playlist
        if (masterPlaylistUrl != null) {
            try {
                String actualUrl = masterPlaylistUrl.startsWith("http") ? 
                    masterPlaylistUrl : getBaseUrl(baseUrl) + masterPlaylistUrl;
                    
                log.info("发现master playlist，获取实际播放列表: {}", actualUrl);
                String playlistContent = HttpUtil.get(actualUrl);
                return parseM3u8(playlistContent, actualUrl); // 递归解析实际的playlist
            } catch (Exception e) {
                log.error("获取实际播放列表失败", e);
                throw e;
            }
        }
        
        // 3. 解析片段
        for (String line : lines) {
            line = line.trim();
            if (!line.startsWith("#") && !line.isEmpty()) {
                // 这是一个媒体片段URL
                segments.add(line);
            } else if (line.startsWith("#EXTINF:")) {
                // 可以获取片段时长信息
                double duration = Double.parseDouble(line.split(":")[1].replace(",", ""));
                log.debug("片段时长: {}秒", duration);
            }
        }
        
        log.info("解析到 {} 个视频片段", segments.size());
        return segments;
    }
    
    private void processBlobDownload(String blobUrl, String savePath) throws Exception {
        try {
            log.info("开始获取视频资源: {}", blobUrl);
            String tempDir = savePath + "_segments";
            FileUtil.mkdir(tempDir);
            Set<String> downloadedUrls = new HashSet<>();
            
            // 1. 初始化资源监控对象和视频控制
            ((JavascriptExecutor) webDriver).executeScript(
                "// 初始化资源监控\n" +
                "window._videoResources = {\n" +
                "    m3u8: new Set(),\n" +
                "    ts: new Set(),\n" +
                "    mp4: new Set(),\n" +
                "    lastUpdate: Date.now()\n" +
                "};\n" +
                "\n" +
                "// 初始化视频控制\n" +
                "window._videoControl = {\n" +
                "    maxPlaybackRate: 16.0,\n" +
                "    setVideoSpeed: function(video) {\n" +
                "        if (!video) return;\n" +
                "        try {\n" +
                "            video.muted = true;\n" +
                "            video.autoplay = true;\n" +
                "            video.playbackRate = this.maxPlaybackRate;\n" +
                "            video.play().catch(e => console.log('播放失败，忽略错误:', e));\n" +
                "            \n" +
                "            // 监听速度变化\n" +
                "            video.addEventListener('ratechange', () => {\n" +
                "                if (video.playbackRate < this.maxPlaybackRate) {\n" +
                "                    video.playbackRate = this.maxPlaybackRate;\n" +
                "                }\n" +
                "            });\n" +
                "            \n" +
                "            // 监听暂停\n" +
                "            video.addEventListener('pause', () => {\n" +
                "                video.play().catch(e => console.log('重新播放失败，忽略错误:', e));\n" +
                "            });\n" +
                "            \n" +
                "            console.log('视频加速成功:', video.src);\n" +
                "        } catch (e) {\n" +
                "            console.error('设置视频速度失败:', e);\n" +
                "        }\n" +
                "    },\n" +
                "    \n" +
                "    accelerateAllVideos: function() {\n" +
                "        const videos = document.getElementsByTagName('video');\n" +
                "        for (let video of videos) {\n" +
                "            this.setVideoSpeed(video);\n" +
                "        }\n" +
                "    }\n" +
                "};\n" +
                "\n" +
                "// 定期检查和加速视频\n" +
                "setInterval(() => window._videoControl.accelerateAllVideos(), 500);\n" +
                "\n" +
                "// 监听新添加的视频元素\n" +
                "const observer = new MutationObserver(mutations => {\n" +
                "    mutations.forEach(mutation => {\n" +
                "        mutation.addedNodes.forEach(node => {\n" +
                "            if (node.nodeName === 'VIDEO') {\n" +
                "                window._videoControl.setVideoSpeed(node);\n" +
                "            }\n" +
                "        });\n" +
                "    });\n" +
                "});\n" +
                "\n" +
                "observer.observe(document.body, {\n" +
                "    childList: true,\n" +
                "    subtree: true\n" +
                "});\n" +
                "\n" +
                "// 注入网络请求拦截器\n" +
                "(function() {\n" +
                "    // 拦截XHR请求\n" +
                "    const XHR = XMLHttpRequest.prototype;\n" +
                "    const open = XHR.open;\n" +
                "    const send = XHR.send;\n" +
                "    \n" +
                "    XHR.open = function() {\n" +
                "        this._url = arguments[1];\n" +
                "        return open.apply(this, arguments);\n" +
                "    };\n" +
                "    \n" +
                "    XHR.send = function() {\n" +
                "        this.addEventListener('load', function() {\n" +
                "            try {\n" +
                "                const url = this._url;\n" +
                "                const contentType = this.getResponseHeader('content-type');\n" +
                "                \n" +
                "                if (url.includes('.m3u8') || contentType?.includes('application/vnd.apple.mpegurl')) {\n" +
                "                    window._videoResources.m3u8.add(url);\n" +
                "                    window._videoResources.lastUpdate = Date.now();\n" +
                "                } else if (url.includes('.ts')) {\n" +
                "                    window._videoResources.ts.add(url);\n" +
                "                    window._videoResources.lastUpdate = Date.now();\n" +
                "                } else if (url.includes('.mp4') || contentType?.includes('video/mp4')) {\n" +
                "                    window._videoResources.mp4.add(url);\n" +
                "                    window._videoResources.lastUpdate = Date.now();\n" +
                "                }\n" +
                "            } catch (e) {\n" +
                "                console.error('XHR拦截器错误:', e);\n" +
                "            }\n" +
                "        });\n" +
                "        return send.apply(this, arguments);\n" +
                "    };\n" +
                "})();\n" +
                "\n" +
                "// 拦截fetch请求\n" +
                "(function() {\n" +
                "    const originalFetch = window.fetch;\n" +
                "    window.fetch = async function(input, init) {\n" +
                "        const url = typeof input === 'string' ? input : input.url;\n" +
                "        const response = await originalFetch(input, init);\n" +
                "        \n" +
                "        try {\n" +
                "            const contentType = response.headers.get('content-type');\n" +
                "            \n" +
                "            if (url.includes('.m3u8') || contentType?.includes('application/vnd.apple.mpegurl')) {\n" +
                "                window._videoResources.m3u8.add(url);\n" +
                "                window._videoResources.lastUpdate = Date.now();\n" +
                "            } else if (url.includes('.ts')) {\n" +
                "                window._videoResources.ts.add(url);\n" +
                "                window._videoResources.lastUpdate = Date.now();\n" +
                "            } else if (url.includes('.mp4') || contentType?.includes('video/mp4')) {\n" +
                "                window._videoResources.mp4.add(url);\n" +
                "                window._videoResources.lastUpdate = Date.now();\n" +
                "            }\n" +
                "        } catch (e) {\n" +
                "            console.error('Fetch拦截器错误:', e);\n" +
                "        }\n" +
                "        \n" +
                "        return response;\n" +
                "    };\n" +
                "})();"
            );
            
            // 2. 等待一段时间确保脚本注入完成
            log.info("等待视频加载和播放...");
            Thread.sleep(2000);
            
            // 3. 开始监控资源
            log.info("开始监控视频资源...");
            int maxWaitTime = 60;
            int noNewResourceCount = 0;
            long lastUpdateTime = System.currentTimeMillis();
            
            for (int i = 0; i < maxWaitTime && noNewResourceCount < 15; i++) {
                boolean foundNewResource = false;
                
                // 获取当前资源状态
                Long jsLastUpdate = (Long) ((JavascriptExecutor) webDriver).executeScript(
                    "return window._videoResources.lastUpdate;"
                );
                
                if (jsLastUpdate > lastUpdateTime) {
                    lastUpdateTime = jsLastUpdate;
                    foundNewResource = true;
                    
                    // 获取并处理m3u8资源
                    @SuppressWarnings("unchecked")
                    List<String> m3u8Urls = (List<String>) ((JavascriptExecutor) webDriver).executeScript(
                        "return Array.from(window._videoResources.m3u8);"
                    );
                    
                    for (String m3u8Url : m3u8Urls) {
                        if (!downloadedUrls.contains(m3u8Url)) {
                            log.info("处理新的m3u8: {}", m3u8Url);
                            processM3u8Download(m3u8Url, savePath);
                            downloadedUrls.add(m3u8Url);
                        }
                    }
                    
                    // 获取并处理ts资源
                    @SuppressWarnings("unchecked")
                    List<String> tsUrls = (List<String>) ((JavascriptExecutor) webDriver).executeScript(
                        "return Array.from(window._videoResources.ts);"
                    );
                    
                    for (String tsUrl : tsUrls) {
                        if (!downloadedUrls.contains(tsUrl)) {
                            log.info("处理新的ts文件: {}", tsUrl);
                            downloadSingleTsFile(tsUrl, tempDir, downloadedUrls.size(), webDriver.getCurrentUrl());
                            downloadedUrls.add(tsUrl);
                        }
                    }
                    
                    // 获取并处理mp4资源
                    @SuppressWarnings("unchecked")
                    List<String> mp4Urls = (List<String>) ((JavascriptExecutor) webDriver).executeScript(
                        "return Array.from(window._videoResources.mp4);"
                    );
                    
                    for (String mp4Url : mp4Urls) {
                        if (!downloadedUrls.contains(mp4Url)) {
                            log.info("处理新的mp4: {}", mp4Url);
                            DownloadTask mp4Task = new DownloadTask();
                            mp4Task.setUrl(mp4Url);
                            mp4Task.setSavePath(savePath);
                            processNormalDownload(mp4Task);
                            downloadedUrls.add(mp4Url);
                        }
                    }
                }
                
                if (foundNewResource) {
                    noNewResourceCount = 0;
                } else {
                    noNewResourceCount++;
                    if (i % 5 == 0) {
                        log.info("等待新的视频资源... {}秒", i);
                        // 尝试触发更多加载
                        ((JavascriptExecutor) webDriver).executeScript(
                            "window._videoControl.accelerateAllVideos();\n" +
                            "const videos = document.getElementsByTagName('video');\n" +
                            "for (let video of videos) {\n" +
                            "    if (video.duration) {\n" +
                            "        video.currentTime = Math.max(0, video.duration - 5);\n" +
                            "    }\n" +
                            "}"
                        );
                    }
                }
                
                Thread.sleep(1000);
            }
            
            if (downloadedUrls.isEmpty()) {
                throw new BusinessException("未找到任何视频资源");
            }
            
            // 4. 合并视频片段
            if (new File(tempDir).exists() && Objects.requireNonNull(new File(tempDir).list()).length > 0) {
                log.info("开始合并视频片段");
                mergeSegments(tempDir, savePath);
                FileUtil.del(tempDir);
            }
            
            // 5. 清理注入的脚本
            ((JavascriptExecutor) webDriver).executeScript(
                "window._videoResources = null;\n" +
                "window._videoControl = null;\n" +
                "const videos = document.getElementsByTagName('video');\n" +
                "for (let video of videos) {\n" +
                "    video.playbackRate = 1.0;\n" +
                "}"
            );
            
            log.info("视频下载完成: {}", savePath);
            
        } catch (Exception e) {
            log.error("视频下载失败", e);
            throw new BusinessException("视频下载失败: " + e.getMessage());
        }
    }
    
    private void downloadSingleTsFile(String tsUrl, String tempDir, int index, String referer) throws Exception {
        String segmentPath = tempDir + File.separator + String.format("%05d", index) + ".ts";
        int maxRetries = 3;
        int retryDelay = 1000;
        Exception lastException = null;
        
        for (int retry = 0; retry < maxRetries; retry++) {
            try {
                if (retry > 0) {
                    log.info("重试下载ts文件: {}, 第{}次尝试", tsUrl, retry + 1);
                    Thread.sleep(retryDelay * retry);
                }
                
                HttpRequest request = HttpRequest.get(tsUrl)
                    .header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124")
                    .header("Referer", referer)
                    .timeout(30000);
                    
                HttpResponse response = request.execute();
                if (!response.isOk()) {
                    throw new BusinessException("HTTP状态码: " + response.getStatus());
                }
                
                byte[] data = response.bodyBytes();
                if (data == null || data.length == 0) {
                    throw new BusinessException("下载的数据为空");
                }
                
                // 验证文件格式
                if (!isValidTsFile(data)) {
                    throw new BusinessException("无效的ts文件格式");
                }
                
                // 创建父目录（如果不存在）
                File file = new File(segmentPath);
                if (!file.getParentFile().exists()) {
                    file.getParentFile().mkdirs();
                }
                
                // 写入文件
                FileUtil.writeBytes(data, file);
                
                // 验证写入的文件
                if (!file.exists() || file.length() == 0) {
                    throw new BusinessException("文件写入失败");
                }
                
                log.info("ts文件下载成功: {}, 大小: {} bytes", tsUrl, data.length);
                return;
                
            } catch (Exception e) {
                lastException = e;
                if (retry == maxRetries - 1) {
                    throw new BusinessException("ts文件下载失败: " + e.getMessage());
                }
            }
        }
        
        if (lastException != null) {
            throw lastException;
        }
    }
    
    private String getBaseUrl(String url) {
        try {
            URI uri = new URI(url);
            String baseUrl = uri.getScheme() + "://" + uri.getHost();
            if (uri.getPort() != -1) {
                baseUrl += ":" + uri.getPort();
            }
            return baseUrl + uri.getPath().substring(0, uri.getPath().lastIndexOf('/') + 1);
        } catch (Exception e) {
            return url.substring(0, url.lastIndexOf('/') + 1);
        }
    }
    
    private void mergeSegments(String segmentDir, String outputPath) {
        try {
            File dir = new File(segmentDir);
            File[] files = dir.listFiles((d, name) -> name.endsWith(".ts"));
            if (files == null || files.length == 0) {
                log.error("合并失败：目录 {} 中没有找到ts文件", segmentDir);
                throw new RuntimeException("没有找到ts文件");
            }
            
            log.info("开始合并 {} 个ts文件", files.length);
            
            // 按文件名数字顺序排序
            Arrays.sort(files, (f1, f2) -> {
                String n1 = f1.getName().replace(".ts", "");
                String n2 = f2.getName().replace(".ts", "");
                return n1.compareTo(n2);
            });
            
            // 验证所有ts文件的大小
            for (File file : files) {
                if (file.length() == 0) {
                    log.error("发现空文件: {}", file.getName());
                    throw new RuntimeException("文件 " + file.getName() + " 大小为0");
                }
            }
            
            // 创建临时文件用于合并ts
            String tempTsPath = outputPath + ".temp.ts";
            
            // 先合并所有ts文件
            try (FileOutputStream fos = new FileOutputStream(tempTsPath)) {
                byte[] buffer = new byte[1024 * 1024]; // 1MB buffer
                int totalFiles = files.length;
                int processedFiles = 0;
                long totalBytes = 0;
                
                for (File file : files) {
                    try (FileInputStream fis = new FileInputStream(file)) {
                        int len;
                        long fileBytes = 0;
                        while ((len = fis.read(buffer)) != -1) {
                            fos.write(buffer, 0, len);
                            fileBytes += len;
                        }
                        totalBytes += fileBytes;
                        processedFiles++;
                        log.info("合并进度: {}/{} ({}), 当前文件大小: {} bytes", 
                            processedFiles, totalFiles, 
                            String.format("%.1f%%", (processedFiles * 100.0 / totalFiles)),
                            fileBytes);
                    }
                }
                fos.flush();
                log.info("ts文件合并完成，总大小: {} bytes", totalBytes);
            }
            
            // 检查合并后的文件
            File tempFile = new File(tempTsPath);
            if (!tempFile.exists() || tempFile.length() == 0) {
                throw new RuntimeException("合并后的临时文件无效");
            }
            
            // 使用ffmpeg将ts转换为mp4，并添加文字水印
            String outputMp4Path = outputPath.endsWith(".mp4") ? outputPath : outputPath + ".mp4";
            
            // 构建ffmpeg命令，添加文字水印
            List<String> command = new ArrayList<>();
            command.add("ffmpeg");
            command.add("-i");
            command.add(tempTsPath);
            // 添加文字水印
            command.add("-vf");
            // drawtext参数说明：
            // fontfile：字体文件路径（使用Windows系统自带的微软雅黑）
            // text：水印文字
            // fontsize：字体大小
            // fontcolor：字体颜色（白色，半透明）
            // x：水印位置x坐标（w-tw-10表示距离右边10像素）
            // y：水印位置y坐标（h-th-10表示距离底部10像素）
            // box：是否添加文字背景框
            // boxcolor：背景框颜色（黑色，半透明）C:\Windows\Fonts\STKAITI.TTF
            command.add("drawtext=fontfile='C\\:/Windows/Fonts/STKAITI.TTF':text='恩师阙楠林爬的':fontsize=36:" +
                       "fontcolor=white@0.5:x=w-tw-10:y=h-th-10:box=1:boxcolor=black@0.5");
            command.add("-c:a");
            command.add("copy");  // 复制音频流
            command.add("-bsf:a");
            command.add("aac_adtstoasc");
            command.add("-y");
            command.add(outputMp4Path);
            
            log.info("开始转换为MP4并添加水印: {}", outputMp4Path);
            ProcessBuilder processBuilder = new ProcessBuilder(command);
            Process process = processBuilder.start();
            
            // 读取ffmpeg的输出
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getErrorStream()))) {
                String line;
                while ((line = reader.readLine()) != null) {
                    log.debug("ffmpeg: {}", line);
                }
            }
            
            int exitCode = process.waitFor();
            if (exitCode != 0) {
                throw new RuntimeException("ffmpeg转换失败，退出码: " + exitCode);
            }
            
            // 检查输出文件
            File outputFile = new File(outputMp4Path);
            if (!outputFile.exists() || outputFile.length() == 0) {
                throw new RuntimeException("转换后的MP4文件无效");
            }
            
            log.info("视频转换和水印添加完成: {} ({} bytes)", outputMp4Path, outputFile.length());
            
            // 删除临时ts文件
            if (!new File(tempTsPath).delete()) {
                log.warn("临时文件删除失败: {}", tempTsPath);
            }
            
        } catch (Exception e) {
            log.error("合并/转换文件失败", e);
            throw new RuntimeException("合并/转换文件失败: " + e.getMessage(), e);
        }
    }
    
    public void addTask(DownloadTask task) {
        try {
            downloadQueue.put(task);
            log.info("添加下载任务: {}", task.getUrl());
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            log.error("添加下载任务失败: {}", e.getMessage());
        }
    }
    
    private <T> T retryDownload(Callable<T> task) throws Exception {
        int maxRetries = 3;
        int retryDelay = 5000; // 5秒
        
        Exception lastException = null;
        for (int i = 0; i < maxRetries; i++) {
            try {
                return task.call();
            } catch (Exception e) {
                lastException = e;
                log.warn("下载重试 {}/{}: {}", i + 1, maxRetries, e.getMessage());
                if (i < maxRetries - 1) {
                    Thread.sleep(retryDelay);
                }
            }
        }
        throw lastException;
    }
    
    private List<SegmentStatus> createDownloadPlan(List<String> segments, String tempDir) {
        List<SegmentStatus> plan = new ArrayList<>();
        for (int i = 0; i < segments.size(); i++) {
            SegmentStatus status = new SegmentStatus();
            status.setUrl(segments.get(i));
            status.setStatus("PENDING");
            status.setLocalPath(tempDir + File.separator + String.format("%03d", i) + ".ts");
            plan.add(status);
        }
        return plan;
    }
    
    private void processVideoDownload(DownloadTask task) {
        try {
            Map<String, Object> context = new HashMap<>();
            context.put("webDriver", webDriver);
            context.put("headers", task.getHeaders());
            
            VideoInfo videoInfo = videoParserManager.parseVideo(task.getUrl(), context);
            
            for (VideoSource source : videoInfo.getSources()) {
                DownloadTask sourceTask = new DownloadTask();
                sourceTask.setUrl(source.getUrl());
                sourceTask.setSavePath(task.getSavePath());
                sourceTask.setHeaders(source.getHeaders());
                
                switch (source.getType()) {
                    case "DIRECT":
                        processNormalDownload(sourceTask);
                        break;
                    case "M3U8":
                        processM3u8Download(sourceTask);
                        break;
                    case "BLOB":
                        processBlobDownload(sourceTask.getUrl(), sourceTask.getSavePath());  // 修改这里，传入正确的参数
                        break;
                    default:
                        throw new BusinessException("不支持的视频类型: " + source.getType());
                }
            }
        } catch (Exception e) {
            log.error("视频下载失败: {}", task.getUrl(), e);
            throw new BusinessException("视频下载失败: " + e.getMessage());
        }
    }
    
    private void downloadTsFiles(List<String> tsUrls, String tempDir) {
        int totalSegments = tsUrls.size();
        int retryCount = 3;
        int retryDelay = 1000; // 1秒
        
        for (int i = 0; i < tsUrls.size(); i++) {
            String tsUrl = tsUrls.get(i);
            String segmentPath = tempDir + File.separator + String.format("%05d", i) + ".ts";
            boolean downloadSuccess = false;
            
            for (int retry = 0; retry < retryCount && !downloadSuccess; retry++) {
                try {
                    if (retry > 0) {
                        log.info("重试下载 {}/{}, 第{}次尝试", i + 1, totalSegments, retry + 1);
                        Thread.sleep(retryDelay * retry);
                    }
                    
                    HttpRequest request = HttpRequest.get(tsUrl)
                        .header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124")
                        .header("Referer", webDriver.getCurrentUrl())
                        .timeout(30000);
                        
                    HttpResponse response = request.execute();
                    if (response.isOk()) {
                        byte[] data = response.bodyBytes();
                        if (data == null || data.length == 0) {
                            throw new BusinessException("下载的数据为空");
                        }
                        
                        // 验证文件格式（检查ts文件头）
                        if (!isValidTsFile(data)) {
                            throw new BusinessException("无效的ts文件格式");
                        }
                        
                        FileUtil.writeBytes(data, segmentPath);
                        
                        // 验证写入的文件
                        File segment = new File(segmentPath);
                        if (!segment.exists() || segment.length() == 0) {
                            throw new BusinessException("文件写入失败");
                        }
                        
                        downloadSuccess = true;
                        double progress = (i + 1.0) / totalSegments * 100;
                        log.info("下载进度: {}/{} ({}%), 文件大小: {} bytes",
                            i + 1,totalSegments, String.format("%.2f",progress),  data.length);
                            
                    } else {
                        throw new BusinessException("HTTP状态码: " + response.getStatus());
                    }
                } catch (Exception e) {
                    log.error("下载分片失败 {}/{}: {}", i + 1, totalSegments, e.getMessage());
                    if (retry == retryCount - 1) {
                        throw new BusinessException("分片下载失败: " + e.getMessage());
                    }
                }
            }
        }
    }
    
    private boolean isValidTsFile(byte[] data) {
        // TS文件的同步字节是0x47 (71)
        if (data.length < 188) { // TS包的标准大小是188字节
            return false;
        }
        return data[0] == 0x47;
    }
} 