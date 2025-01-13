package com.security.service.impl;

import cn.hutool.core.util.StrUtil;
import cn.hutool.http.HttpRequest;
import cn.hutool.http.HttpResponse;
import cn.hutool.http.HttpUtil;
import cn.hutool.core.util.ReUtil;
import com.alibaba.fastjson2.JSON;
import com.security.exception.BusinessException;
import com.security.config.ResourceConfig;
import com.security.model.dto.response.CrawlerResultDTO;
import com.security.model.entity.WebsiteInfo;
import com.security.mapper.WebsiteInfoMapper;
import com.security.service.IWebCrawlerService;
import com.security.service.IResourceProcessService;
import com.security.service.download.DownloadQueueManager;
import com.security.service.download.DownloadTask;
import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.select.Elements;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Service;
import org.jsoup.nodes.Element;

import java.net.URI;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.Executor;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import java.util.concurrent.ConcurrentHashMap;
import java.util.regex.Matcher;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.WebDriverWait;
import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.logging.LogEntries;
import org.openqa.selenium.logging.LogEntry;
import org.openqa.selenium.logging.LogType;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.JsonNode;
import cn.hutool.core.util.RandomUtil;

@Service
public class WebCrawlerServiceImpl implements IWebCrawlerService {
    
    private static final Logger log = LoggerFactory.getLogger(WebCrawlerServiceImpl.class);
    
    @Autowired
    private WebsiteInfoMapper websiteInfoMapper;
    
    @Autowired
    @Qualifier("scanTaskExecutor")
    private Executor taskExecutor;
    
    @Autowired
    private IResourceProcessService resourceProcessService;
    
    @Autowired
    private DownloadQueueManager downloadQueueManager;

    @Autowired
    private ResourceConfig resourceConfig;
    
    @Autowired
    private WebDriver webDriver;
    
    // 邮箱正则
    private static final Pattern EMAIL_PATTERN = Pattern.compile("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,6}");
    // 电话正则
    private static final Pattern PHONE_PATTERN = Pattern.compile("\\d{3,4}[-]?\\d{7,8}|\\d{11}");
    // API正则
    private static final Pattern API_PATTERN = Pattern.compile("(\\/api\\/[\\w\\-\\/]+)|(\\/v[1-9]\\/[\\w\\-\\/]+)");
    
    // 用于存储爬虫进度的Map
    private static final ConcurrentHashMap<Long, Map<String, Object>> CRAWLER_PROGRESS = new ConcurrentHashMap<>();
    
    private final ObjectMapper objectMapper = new ObjectMapper();
    
    @Override
    public Long createCrawlerTask(String url, Integer depth) {
        // 验证URL
        final String finalUrl = !url.startsWith("http") ? "http://" + url : url;
        
        // 创建任务记录
        WebsiteInfo websiteInfo = new WebsiteInfo();
        websiteInfo.setUrl(finalUrl);
        websiteInfo.setCreateTime(new Date());
        
        // 插入数据并获取ID
        if (websiteInfoMapper.insert(websiteInfo) > 0) {
            Long taskId = websiteInfo.getId();
            log.info("创建爬虫任务成功: taskId={}, url={}", taskId, finalUrl);
            
            // 初始化进度信息
            Map<String, Object> progress = new HashMap<>();
            progress.put("status", 0);  // 0-待开始 1-进行中 2-已完成 3-失败
            progress.put("progress", 0);
            progress.put("step", "初始化");
            progress.put("error", null);
            CRAWLER_PROGRESS.put(taskId, progress);
            
            // 异步执行爬虫任务
            taskExecutor.execute(() -> asyncCrawl(taskId, finalUrl, depth));
            
            return taskId;
        } else {
            throw new BusinessException("创建爬虫任务失败");
        }
    }
    
    @Override
    public void basicInfoCrawl(Long taskId, String url) {
        WebsiteInfo info = new WebsiteInfo();
        info.setTaskId(taskId);
        info.setUrl(url);
        
        int maxRetries = 3;
        int retryCount = 0;
        Exception lastException = null;
        
        while (retryCount < maxRetries) {
            try {
                webDriver.get(url);
                // 1. 增加页面加载超时时间
                webDriver.manage().timeouts().pageLoadTimeout(Duration.ofSeconds(30));
                
                // 2. 添加更多的反检测参数
                ((JavascriptExecutor) webDriver).executeScript(
                    "Object.defineProperty(navigator, 'webdriver', {get: () => undefined});"
                );
                
                // 3. 增加等待时间
                WebDriverWait wait = new WebDriverWait(webDriver, Duration.ofSeconds(20));
                
                // 4. 添加随机延迟
                Thread.sleep(RandomUtil.randomInt(1000, 3000));
                
                // 5. 等待页面加载完成
                wait.until(webDriver -> ((JavascriptExecutor) webDriver)
                    .executeScript("return document.readyState").equals("complete"));
                    
                // 6. 检查是否有反爬提示
                if (webDriver.getPageSource().contains("访问受限") || 
                    webDriver.getPageSource().contains("人机验证")) {
                    throw new BusinessException("检测到反爬虫机制，请稍后重试");
                }
                
                // 等待动态内容加载
                Thread.sleep(2000);
                
                // 1. 首先获取XHR请求的数据（优先获取后端资源）
                log.info("开始捕获网络请求: taskId={}", taskId);
                List<Map<String, Object>> networkData = captureNetworkData();
                log.info("网络请求捕获完成: taskId={}, 请求数量={}", taskId, networkData.size());
                
                // 从XHR响应中提取视频
                List<String> videoUrls = new ArrayList<>();
                List<String> networkVideoUrls = extractVideoFromNetwork(networkData);
                log.info("从网络请求中提取到视频: taskId={}, 数量={}, urls={}", 
                        taskId, networkVideoUrls.size(), networkVideoUrls);
                videoUrls.addAll(networkVideoUrls);
                
                // 2. 然后获取页面内容（静态前端资源）
                String html = webDriver.getPageSource();
                Document doc = Jsoup.parse(html, url);
                log.info("页面解析完成: taskId={}, 标题={}", taskId, doc.title());
                
                // 从HTML中提取视频
                List<String> htmlVideoUrls = extractVideoFromHtml(doc);
                log.info("从HTML中提取到视频: taskId={}, 数量={}, urls={}", 
                        taskId, htmlVideoUrls.size(), htmlVideoUrls);
                videoUrls.addAll(htmlVideoUrls);
                
                // 更新数据库
                info.setId(taskId);
                info.setTitle(doc.title());
                info.setVideos(JSON.toJSONString(videoUrls));
                
                // 如果有视频资源，添加到下载队列
                if (!videoUrls.isEmpty()) {
                    log.info("发现视频资源: taskId={}, 总数={}", taskId, videoUrls.size());
                    for (String videoUrl : videoUrls) {
                        String fileName = getFileName(videoUrl);
                        String savePath = getSavePath(taskId, "video", fileName);
                        DownloadTask task = new DownloadTask();
                        task.setUrl(videoUrl);
                        task.setSavePath(savePath);
                        downloadQueueManager.addTask(task);
                        log.info("添加视频下载任务: taskId={}, url={}, savePath={}", 
                                taskId, videoUrl, savePath);
                    }
                }
                
                websiteInfoMapper.updateById(info);
                log.info("基本信息更新完成: taskId={}", taskId);
                
                return; // 如果成功，直接返回
            } catch (Exception e) {
                lastException = e;
                retryCount++;
                log.warn("爬取失败，第{}次重试: taskId={}, error={}", 
                    retryCount, taskId, e.getMessage());
                    
                if (retryCount < maxRetries) {
                    try {
                        // 指数退避重试
                        Thread.sleep(1000 * (long)Math.pow(2, retryCount));
                    } catch (InterruptedException ie) {
                        Thread.currentThread().interrupt();
                        break;
                    }
                }
            }
        }
        
        // 所有重试都失败后抛出异常
        log.error("基本信息爬取失败: taskId={}, error={}", taskId, lastException.getMessage());
        throw new BusinessException("基本信息爬取失败: " + lastException.getMessage());
    }
    
    private List<Map<String, Object>> captureNetworkData() {
        List<Map<String, Object>> networkData = new ArrayList<>();
        try {
            // 等待页面加载完成
            Thread.sleep(2000);
            
            // 获取性能日志
            LogEntries logs = webDriver.manage().logs().get(LogType.PERFORMANCE);
            if (logs != null) {
                for (LogEntry entry : logs) {
                    try {
                        JsonNode json = objectMapper.readTree(entry.getMessage());
                        JsonNode message = json.get("message");
                        String method = message.get("method").asText();
                        
                        if (method.equals("Network.responseReceived")) {
                            JsonNode params = message.get("params");
                            JsonNode response = params.get("response");
                            String url = response.get("url").asText();
                            String type = response.get("mimeType").asText();
                            
                            // 处理视频相关的响应
                            if (isVideoResponse(type, url)) {
                                String actualUrl = handleVideoUrl(url, type);
                                if (actualUrl != null) {
                                    Map<String, Object> data = new HashMap<>();
                                    data.put("url", actualUrl);
                                    data.put("type", type);
                                    networkData.add(data);
                                    log.info("发现视频资源: url={}, type={}", actualUrl, type);
                                }
                            }
                        }
                    } catch (Exception e) {
                        log.warn("解析网络日志失败: {}", e.getMessage());
                    }
                }
            }
        } catch (Exception e) {
            log.warn("捕获网络数据失败: {}", e.getMessage());
        }
        return networkData;
    }
    
    private boolean isVideoResponse(String type, String url) {
        return (type != null && (
                type.startsWith("video/") ||
                type.contains("mpegurl") ||
                type.contains("mp4") ||
                type.contains("m3u8"))) ||
               (url != null && (
                url.contains(".mp4") ||
                url.contains(".m3u8") ||
                url.contains(".flv") ||
                url.matches(".*blob:.*")));
    }
    
    private String handleVideoUrl(String url, String type) {
        try {
            if (url.startsWith("blob:")) {
                // 对于 blob URL，尝试获取真实的媒体流地址
                String script = 
                    "var xhr = new XMLHttpRequest();" +
                    "xhr.open('GET', '" + url + "', false);" +
                    "xhr.send(null);" +
                    "return xhr.responseURL;";
                
                Object result = ((JavascriptExecutor) webDriver).executeScript(script);
                if (result != null) {
                    String actualUrl = result.toString();
                    log.info("获取到blob真实地址: {} -> {}", url, actualUrl);
                    return actualUrl;
                }
                
                // 如果无法获取真实URL，尝试直接获取媒体源
                script = 
                    "var video = document.querySelector('video');" +
                    "return video ? video.src : null;";
                
                result = ((JavascriptExecutor) webDriver).executeScript(script);
                if (result != null) {
                    String mediaUrl = result.toString();
                    log.info("获取到视频源地址: {}", mediaUrl);
                    return mediaUrl;
                }
            }
            return url;
        } catch (Exception e) {
            log.warn("处理视频URL失败: {}", e.getMessage());
            return null;
        }
    }
    
    private List<String> extractVideoFromHtml(Document doc) {
        List<String> videoUrls = new ArrayList<>();
        
        // 1. 查找 video 标签
        Elements videoElements = doc.select("video source, video");
        for (Element video : videoElements) {
            String src = video.attr("src");
            if (StrUtil.isNotEmpty(src)) {
                videoUrls.add(src);
            }
        }
        
        // 2. 查找 iframe (可能包含嵌入的视频)
        Elements iframes = doc.select("iframe");
        for (Element iframe : iframes) {
            String src = iframe.attr("src");
            if (isVideoUrl(src)) {
                videoUrls.add(src);
            }
        }
        
        // 3. 查找 object 和 embed 标签
        Elements objects = doc.select("object, embed");
        for (Element obj : objects) {
            String data = obj.attr("data");
            if (isVideoUrl(data)) {
                videoUrls.add(data);
            }
        }
        
        return videoUrls;
    }
    
    private List<String> extractVideoFromNetwork(List<Map<String, Object>> networkData) {
        List<String> videoUrls = new ArrayList<>();
        
        // 1. 查找视频信息API响应
        Optional<Map<String, Object>> videoInfo = networkData.stream()
            .filter(data -> {
                String url = (String) data.get("url");
                return url != null && url.contains("getinfo") 
                    && url.contains("vv.video.qq.com");
            })
            .findFirst();
            
        if (videoInfo.isPresent()) {
            try {
                Map<String, Object> info = videoInfo.get();
                JsonNode response = objectMapper.readTree((String) info.get("response"));
                // 提取视频URL(具体字段根据API响应结构调整)
                JsonNode urls = response.path("vl").path("vi")
                    .get(0).path("ul").path("ui");
                for (JsonNode url : urls) {
                    videoUrls.add(url.path("url").asText());
                }
            } catch (Exception e) {
                log.error("解析视频信息失败", e);
            }
        }
        
        return videoUrls;
    }
    
    private void extractVideoUrlsFromJson(JsonNode node, List<String> videoUrls) {
        if (node.isObject()) {
            node.fields().forEachRemaining(entry -> {
                String key = entry.getKey().toLowerCase();
                JsonNode value = entry.getValue();
                
                // 检查常见的视频相关字段
                if ((key.contains("video") || key.contains("media") || key.contains("url") || 
                     key.contains("src") || key.contains("stream")) && 
                    value.isTextual()) {
                    String url = value.asText();
                    if (isVideoUrl(url)) {
                        videoUrls.add(url);
                    }
                }
                
                // 递归处理嵌套对象
                extractVideoUrlsFromJson(value, videoUrls);
            });
        } else if (node.isArray()) {
            node.elements().forEachRemaining(element -> 
                extractVideoUrlsFromJson(element, videoUrls));
        }
    }
    
    private void extractVideoUrlsFromText(String content, List<String> videoUrls) {
        // 匹配视频URL的正则表达式
        List<Pattern> patterns = Arrays.asList(
            Pattern.compile("https?://[^\\s<>\"']+?\\.(mp4|m3u8|flv)[^\\s<>\"']*"),
            Pattern.compile("\"(https?://[^\"]+?\\.(mp4|m3u8|flv)[^\"]*)\""),
            Pattern.compile("'(https?://[^']+?\\.(mp4|m3u8|flv)[^']*)'")
        );
        
        for (Pattern pattern : patterns) {
            Matcher matcher = pattern.matcher(content);
            while (matcher.find()) {
                String url = matcher.group(1) != null ? matcher.group(1) : matcher.group(0);
                if (isVideoUrl(url)) {
                    videoUrls.add(url);
                }
            }
        }
    }
    
    private boolean isVideoUrl(String url) {
        if (url == null) return false;
        url = url.toLowerCase();
        return url.contains("youtube.com") || 
               url.contains("youku.com") ||
               url.contains("vimeo.com") ||
               url.contains("bilibili.com") ||
               url.matches(".*\\.(mp4|m3u8|flv).*") ||
               url.contains("/video/") ||
               url.contains("play") ||
               url.contains("stream");
    }
    
    private String detectFramework(Document doc, HttpResponse response) {
        // 检测常见框架特征
        if (doc.select("meta[name=generator]").size() > 0) {
            return doc.select("meta[name=generator]").attr("content");
        }
        
        // 检测前端框架
        if (doc.select("[ng-app]").size() > 0) return "Angular";
        if (doc.select("[data-reactroot]").size() > 0) return "React";
        if (doc.select("[data-v-]").size() > 0) return "Vue";
        
        // 检测后端框架
        String poweredBy = response.header("X-Powered-By");
        if (poweredBy != null) return poweredBy;
        
        return "Unknown";
    }
    
    private List<String> extractDynamicUrls(Document doc) {
        List<String> urls = new ArrayList<>();
        
        // 提取API请求URL
        Elements scripts = doc.select("script");
        for (Element script : scripts) {
            String scriptContent = script.html();
            // 匹配常见的API调用模式
            List<String> apiUrls = ReUtil.findAll("(api|service|rest)/\\w+", scriptContent, 0);
            urls.addAll(apiUrls);
            
            // 匹配Ajax URL
            List<String> ajaxUrls = ReUtil.findAll("\\$\\.(?:get|post)\\(['\"]([^'\"]+)['\"]", 
                                                  scriptContent, 1);
            urls.addAll(ajaxUrls);
        }
        
        return urls;
    }
    
    private void processAjaxResponse(Long taskId, String jsonResponse) {
        try {
            // 解析JSON响应
            Map<String, Object> data = JSON.parseObject(jsonResponse);
            
            // 更新数据库
            WebsiteInfo info = new WebsiteInfo();
            info.setId(taskId);
            
            // 处理动态数据...
            websiteInfoMapper.updateById(info);
        } catch (Exception e) {
            log.warn("处理Ajax响应失败: {}", e.getMessage());
        }
    }
    
    private void asyncDownloadResources(Long taskId, List<String> images, 
                                      List<String> videos, List<String> files) {
        // 创建下载任务
        for (String url : images) {
            createDownloadTask(taskId, url, "IMAGE");
        }
        
        for (String url : videos) {
            createDownloadTask(taskId, url, "VIDEO");
        }
        
        for (String url : files) {
            createDownloadTask(taskId, url, "FILE");
        }
    }
    
    private void createDownloadTask(Long taskId, String url, String type) {
        DownloadTask task = new DownloadTask();
        task.setTaskId(taskId);
        task.setUrl(url);
        task.setType(type);
        task.setSavePath(getSavePath(taskId, type, getFileName(url)));
        downloadQueueManager.addTask(task);
    }
    
    @Override
    public List<String> linksCrawl(String url, Integer depth) {
        Set<String> links = new HashSet<>();
        crawlLinks(url, depth, links, new HashSet<>());
        return new ArrayList<>(links);
    }
    
    private void crawlLinks(String url, int depth, Set<String> links, Set<String> visited) {
        if (depth <= 0 || visited.contains(url)) {
            return;
        }
        
        visited.add(url);
        
        try {
            String html = HttpUtil.get(url);
            Document doc = Jsoup.parse(html);
            Elements elements = doc.select("a[href]");
            
            elements.forEach(element -> {
                String link = element.attr("abs:href");
                if (StrUtil.isNotEmpty(link) && link.startsWith("http")) {
                    System.out.println("捕获链接: "+link);
                    links.add(link);
                    crawlLinks(link, depth - 1, links, visited);
                }
            });
        } catch (Exception e) {
            // 记录错误日志
        }
    }
    
    @Override
    public void sensitiveInfoCrawl(Long taskId, String url) {
        try {
            log.info("开始爬取敏感信息: taskId={}, url={}", taskId, url);
            
            // 设置超时时间
            long startTime = System.currentTimeMillis();
            long timeout = 5 * 60 * 1000; // 5分钟超时
            
            // 获取页面内容
            String html = HttpRequest.get(url)
                .timeout(10000)  // 10秒连接超时
                .execute()
                .body();
                
            log.info("页面内容获取完成: taskId={}, contentLength={}", taskId, html.length());
            
            // 提取邮箱
            List<String> emails = new ArrayList<>();
            Matcher emailMatcher = EMAIL_PATTERN.matcher(html);
            while (emailMatcher.find() && !isTimeout(startTime, timeout)) {
                System.out.println("emailMatcher.group()");
                emails.add(emailMatcher.group());
            }
            log.info("邮箱提取完成: taskId={}, count={}", taskId, emails.size());

            // 提取电话
            List<String> phones = new ArrayList<>();
            Matcher phoneMatcher = PHONE_PATTERN.matcher(html);
            while (phoneMatcher.find() && !isTimeout(startTime, timeout)) {
                phones.add(phoneMatcher.group());
            }
            log.info("电话提取完成: taskId={}, count={}", taskId, phones.size());
            
            // 更新数据库
            WebsiteInfo info = new WebsiteInfo();
            info.setId(taskId);
            info.setEmails(JSON.toJSONString(emails));
            info.setPhones(JSON.toJSONString(phones));
            websiteInfoMapper.updateById(info);

            log.info("敏感信息爬取完成: taskId={}", taskId);
            
        } catch (Exception e) {
            log.error("敏感信息爬取失败: taskId={}, error={}", taskId, e.getMessage(), e);
            throw new BusinessException("敏感信息收集失败：" + e.getMessage());
        }
    }
    
    // 检查是否超时
    private boolean isTimeout(long startTime, long timeout) {
        return System.currentTimeMillis() - startTime > timeout;
    }
    
    @Override
    public void jsAnalysis(Long taskId, String url) {
        try {
            Document doc = Jsoup.connect(url).get();
            Elements scripts = doc.select("script[src]");
            
            scripts.forEach(script -> {
                String jsUrl = script.attr("abs:src");
                if (StrUtil.isNotEmpty(jsUrl)) {
                    try {
                        String jsContent = HttpUtil.get(jsUrl);
                        // 分析JS内容，提取API和敏感信息
                        List<String> apis = ReUtil.findAll(API_PATTERN, jsContent, 0);
                        // 更新数据库
                        WebsiteInfo info = new WebsiteInfo();
                        info.setId(taskId);
                        info.setLinks(JSON.toJSONString(apis));
                        websiteInfoMapper.updateById(info);
                    } catch (Exception e) {
                        log.warn("爬取js信息错误， {}",e.getMessage());
                    }
                }
            });
        } catch (Exception e) {
            throw new BusinessException("JS分析失败：" + e.getMessage());
        }
    }
    
    @Override
    public CrawlerResultDTO getCrawlerResult(Long taskId) {
        log.info("获取爬虫结果: taskId={}", taskId);
        
        // 从数据库查询结果
        WebsiteInfo info = websiteInfoMapper.selectById(taskId);
        if (info == null) {
            log.warn("未找到爬虫任务: taskId={}", taskId);
            throw new BusinessException("未找到爬虫任务");
        }
        
        // 转换为DTO
        CrawlerResultDTO dto = new CrawlerResultDTO();
        dto.setTaskId(info.getId());
        dto.setUrl(info.getUrl());
        dto.setTitle(info.getTitle());
        dto.setDescription(info.getDescription());
        dto.setKeywords(info.getKeywords());
        dto.setServer(info.getServer());
        dto.setFramework(info.getFramework());
        dto.setCreateTime(info.getCreateTime());
        
        // 转换JSON字符串为List
        if (StrUtil.isNotEmpty(info.getLinks())) {
            List<String> links = JSON.parseArray(info.getLinks(), String.class);
            // 区分内部链接和外部链接
            String domain = getDomain(info.getUrl());
            dto.setInternalLinks(links.stream()
                    .filter(link -> link.contains(domain))
                    .collect(Collectors.toList()));
            dto.setExternalLinks(links.stream()
                    .filter(link -> !link.contains(domain))
                    .collect(Collectors.toList()));
        }
        
        if (StrUtil.isNotEmpty(info.getEmails())) {
            dto.setEmails(JSON.parseArray(info.getEmails(), String.class));
        }
        
        if (StrUtil.isNotEmpty(info.getPhones())) {
            dto.setPhones(JSON.parseArray(info.getPhones(), String.class));
        }
        
        // 添加进度信息
        Map<String, Object> progress = CRAWLER_PROGRESS.get(taskId);
        if (progress != null) {
            dto.setStatus((Integer) progress.get("status"));
            dto.setProgress((Integer) progress.get("progress"));
            dto.setCurrentStep((String) progress.get("step"));
            dto.setErrorMsg((String) progress.get("error"));
        }
        
        log.info("爬虫结果: {}", JSON.toJSONString(dto));
        return dto;
    }
    
    private String getDomain(String url) {
        try {
            URI uri = new URI(url);
            String domain = uri.getHost();
            return domain.startsWith("www.") ? domain.substring(4) : domain;
        } catch (Exception e) {
            return url;
        }
    }
    
    private void asyncCrawl(Long taskId, String url, Integer depth) {
        try {
            long startTime = System.currentTimeMillis();
            long timeout = 5 * 60 * 1000; // 15分钟总超时
            
            updateProgress(taskId, 1, 0, "开始爬取");
            
            // 基本信息爬取 (20%)
            if (isTimeout(startTime, timeout)) {
                throw new BusinessException("爬虫任务超时");
            }
            updateProgress(taskId, 1, 10, "爬取基本信息");
            basicInfoCrawl(taskId, url);
            updateProgress(taskId, 1, 20, "基本信息爬取完成");
            
            // 链接爬取 (40%)
            if (isTimeout(startTime, timeout)) {
                throw new BusinessException("爬虫任务超时");
            }
            updateProgress(taskId, 1, 30, "爬取页面链接");
            List<String> links = linksCrawl(url, depth);
            WebsiteInfo info = new WebsiteInfo();
            info.setId(taskId);
            info.setLinks(JSON.toJSONString(links));
            websiteInfoMapper.updateById(info);
            updateProgress(taskId, 1, 40, "链接爬取完成");
            
            // 敏感信息爬取 (70%)
            if (isTimeout(startTime, timeout)) {
                throw new BusinessException("爬虫任务超时");
            }
            updateProgress(taskId, 1, 50, "爬取敏感信息");
            sensitiveInfoCrawl(taskId, url);
            updateProgress(taskId, 1, 70, "敏感信息爬取完成");
            
            // JS分析 (100%)
            if (isTimeout(startTime, timeout)) {
                throw new BusinessException("爬虫任务超时");
            }
            updateProgress(taskId, 1, 80, "分析JS文件");
            jsAnalysis(taskId, url);
            
            // 完成
            updateProgress(taskId, 2, 100, "爬取完成");
            
        } catch (Exception e) {
            String errorMsg = e instanceof BusinessException ? e.getMessage() : "爬虫任务异常：" + e.getMessage();
            log.error("异步爬虫任务异常: taskId={}, error={}", taskId, errorMsg, e);
            updateProgress(taskId, 3, 0, "爬取失败", errorMsg);
        } finally {
            // 可选：任务完成后一段时间清除进度信息
            taskExecutor.execute(() -> {
                try {
                    Thread.sleep(30 * 60 * 1000); // 30分钟后清除
                    CRAWLER_PROGRESS.remove(taskId);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            });
        }
    }
    
    private void updateProgress(Long taskId, Integer status, Integer progress, String step) {
        updateProgress(taskId, status, progress, step, null);
    }
    
    private void updateProgress(Long taskId, Integer status, Integer progress, String step, String error) {
        Map<String, Object> progressInfo = new HashMap<>();
        progressInfo.put("status", status);
        progressInfo.put("progress", progress);
        progressInfo.put("step", step);
        progressInfo.put("error", error);
        CRAWLER_PROGRESS.put(taskId, progressInfo);
        
        log.info("爬虫进度更新: taskId={}, status={}, progress={}, step={}", 
                taskId, status, progress, step);
    }
    
    private String getFileName(String url) {
        try {
            String fileName = url.substring(url.lastIndexOf('/') + 1);
            // 移除查询参数
            if (fileName.contains("?")) {
                fileName = fileName.substring(0, fileName.indexOf("?"));
            }
            return fileName;
        } catch (Exception e) {
            return "download_" + System.currentTimeMillis();
        }
    }
    
    private String getSavePath(Long taskId, String type, String fileName) {
        return resourceConfig.getDownloadPath() + "/" + taskId + "/" + type.toLowerCase() + "/" + fileName;
    }
    
    private String getBaseUrl(String url) {
        try {
            URI uri = new URI(url);
            String baseUrl = uri.getScheme() + "://" + uri.getHost();
            if (uri.getPort() != -1) {
                baseUrl += ":" + uri.getPort();
            }
            return baseUrl;
        } catch (Exception e) {
            return url;
        }
    }
} 