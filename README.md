# 实现的接口列表

## 1. 学习活动服务接口
- **接口路径**: `src/main/java/com/security/service/LearningActivityService.java`
- **功能**: 管理学习活动的进度和状态。
- **方法**:
  - `void startLearningActivity(String activityId);` - 开始学习活动
  - `void completeLearningActivity(String activityId);` - 完成学习活动
  - `void updateLearningProgress(String activityId, int progress);` - 更新学习进度
  - `String getLearningActivityStatus(String activityId);` - 获取学习活动状态

## 2. 视频播放服务接口
- **接口路径**: `src/main/java/com/security/service/VideoPlaybackService.java`
- **功能**: 管理视频播放的进度和状态。
- **方法**:
  - `void startVideo(String videoId);` - 开始播放视频
  - `void pauseVideo(String videoId);` - 暂停视频播放
  - `void resumeVideo(String videoId);` - 恢复视频播放
  - `void updateVideoProgress(String videoId, int progress);` - 更新视频播放进度
  - `String getVideoPlaybackStatus(String videoId);` - 获取视频播放状态

## 3. 学习平台服务接口
- **接口路径**: `src/main/java/com/security/service/LearningPlatformService.java`
- **功能**: 与不同学习平台进行交互。
- **方法**:
  - `void login(String username, String password);` - 登录学习平台
  - `void logout();` - 登出学习平台
  - `LearningActivityService getLearningActivityService();` - 获取学习活动服务
  - `VideoPlaybackService getVideoPlaybackService();` - 获取视频播放服务

## 4. 进度监控服务接口
- **接口路径**: `src/main/java/com/security/service/ProgressMonitoringService.java`
- **功能**: 监控学习和视频播放的进度。
- **方法**:
  - `void monitorLearningProgress(String activityId);` - 监控学习进度
  - `void monitorVideoProgress(String videoId);` - 监控视频播放进度
  - `void reportProgress(String activityId, int progress);` - 报告学习进度
  - `void reportVideoProgress(String videoId, int progress);` - 报告视频播放进度

## 5. 通用爬虫服务接口
- **接口路径**: `src/main/java/com/security/service/GenericCrawlerService.java`
- **功能**: 爬取学习平台的内容。
- **方法**:
  - `void crawlLearningContent(String url);` - 爬取学习内容
  - `void crawlVideoContent(String url);` - 爬取视频内容
  - `void updateProgress(String url, String contentType, int progress);` - 更新进度

## 6. 认证和授权服务接口
- **接口路径**: `src/main/java/com/security/service/AuthenticationService.java`
- **功能**: 测试 Web 应用的身份验证和授权机制。
- **方法**:
  - `boolean testLogin(String url, String username, String password);` - 测试登录
  - `boolean testPasswordReset(String url, String email);` - 测试密码重置
  - `boolean testSessionManagement(String url);` - 测试会话管理

## 7. 输入验证服务接口
- **接口路径**: `src/main/java/com/security/service/InputValidationService.java`
- **功能**: 测试 Web 应用对用户输入的验证。
- **方法**:
  - `boolean testInputValidation(String url, String payload);` - 测试输入验证
  - `boolean testOutputEncoding(String url);` - 测试输出编码

## 8. SQL 注入检测服务接口
- **接口路径**: `src/main/java/com/security/service/SqlInjectionService.java`
- **功能**: 检测 Web 应用是否存在 SQL 注入漏洞。
- **方法**:
  - `boolean detectSqlInjection(String url, String payload);` - 检测 SQL 注入
  - `boolean testForBlindSqlInjection(String url);` - 测试盲 SQL 注入

## 9. CSRF 检测服务接口
- **接口路径**: `src/main/java/com/security/service/CsrfDetectionService.java`
- **功能**: 检测 Web 应用是否存在跨站请求伪造漏洞。
- **方法**:
  - `boolean detectCsrf(String url);` - 检测 CSRF
  - `boolean validateCsrfTokens(String url);` - 验证 CSRF 令牌

## 10. XSS 检测服务接口
- **接口路径**: `src/main/java/com/security/service/XssDetectionService.java`
- **功能**: 检测 Web 应用是否存在跨站脚本漏洞。
- **方法**:
  - `boolean detectXss(String url, String payload);` - 检测 XSS
  - `boolean testForStoredXss(String url);` - 测试存储型 XSS

## 11. 业务逻辑漏洞检测服务接口
- **接口路径**: `src/main/java/com/security/service/BusinessLogicVulnerabilityService.java`
- **功能**: 检测 Web 应用的业务逻辑漏洞。
- **方法**:
  - `boolean detectBusinessLogicVulnerability(String url);` - 检测业务逻辑漏洞
  - `boolean validateBusinessProcesses(String url);` - 验证业务流程

## 12. 事件管理服务接口
- **接口路径**: `src/main/java/com/security/service/IncidentManagementService.java`
- **功能**: 处理安全事件和响应。
- **方法**:
  - `void logIncident(String incidentDetails);` - 记录事件
  - `void escalateIncident(String incidentId);` - 升级事件
  - `void resolveIncident(String incidentId);` - 解决事件

## 13. 组件和库漏洞扫描服务接口
- **接口路径**: `src/main/java/com/security/service/ComponentVulnerabilityService.java`
- **功能**: 检查 Web 应用使用的组件和库是否存在已知漏洞。
- **方法**:
  - `List<String> scanForComponentVulnerabilities(String url);` - 扫描组件漏洞
  - `boolean checkForLicenseCompliance(String url);` - 检查许可证合规性

## 14. 安全审计服务接口
- **接口路径**: `src/main/java/com/security/service/SecurityAuditService.java`
- **功能**: 对 Web 应用进行安全审计。
- **方法**:
  - `boolean performSecurityAudit(String url);` - 执行安全审计
  - `String getAuditReport(String url);` - 获取审计报告

## 15. 安全意识培训服务接口
- **接口路径**: `src/main/java/com/security/service/SecurityAwarenessTrainingService.java`
- **功能**: 提供安全意识培训。
- **方法**:
  - `void conductTrainingSession(String sessionDetails);` - 进行培训课程
  - `String getTrainingFeedback(String sessionId);` - 获取培训反馈

## 16. 爬虫配置服务接口
- **接口路径**: `src/main/java/com/security/service/CrawlerConfigurationService.java`
- **功能**: 配置爬虫的参数和选项。
- **方法**:
  - `void setUserAgent(String userAgent);` - 设置用户代理
  - `void setMaxDepth(int maxDepth);` - 设置最大爬取深度
  - `void setTimeout(int timeout);` - 设置请求超时
  - `void setMaxPages(int maxPages);` - 设置最大爬取页面数

## 17. 爬虫调度服务接口
- **接口路径**: `src/main/java/com/security/service/CrawlerSchedulerService.java`
- **功能**: 调度和管理爬虫任务。
- **方法**:
  - `void scheduleCrawl(String url);` - 调度爬取任务
  - `void cancelCrawl(String taskId);` - 取消爬取任务
  - `void pauseCrawl(String taskId);` - 暂停爬取任务
  - `void resumeCrawl(String taskId);` - 恢复爬取任务

## 18. 爬虫结果处理服务接口
- **接口路径**: `src/main/java/com/security/service/CrawlResultHandlerService.java`
- **功能**: 处理爬取结果。
- **方法**:
  - `void saveCrawlResults(String url, String results);` - 保存爬取结果
  - `void exportResultsToJson(String url);` - 导出结果为 JSON 格式
  - `void exportResultsToCsv(String url);` - 导出结果为 CSV 格式

## 19. 爬虫状态监控服务接口
- **接口路径**: `src/main/java/com/security/service/CrawlerStatusService.java`
- **功能**: 监控爬虫的运行状态。
- **方法**:
  - `String getCrawlStatus(String taskId);` - 获取爬取任务状态
  - `int getTotalPagesCrawled(String taskId);` - 获取已爬取的总页面数
  - `int getTotalErrors(String taskId);` - 获取爬取过程中发生的错误总数

## 20. 爬虫日志服务接口
- **接口路径**: `src/main/java/com/security/service/CrawlerLogService.java`
- **功能**: 记录和管理爬虫日志。
- **方法**:
  - `void logInfo(String message);` - 记录信息日志
  - `void logError(String message);` - 记录错误日志
  - `String getLog(String taskId);` - 获取特定任务的日志

## 21. 爬虫策略服务接口
- **接口路径**: `src/main/java/com/security/service/CrawlerStrategyService.java`
- **功能**: 定义爬虫的爬取策略。
- **方法**:
  - `void setCrawlDelay(int delay);` - 设置爬取延迟
  - `void setCrawlPolicy(String policy);` - 设置爬取策略（如深度优先、广度优先）
  - `void setAllowedDomains(String[] domains);` - 设置允许爬取的域名

## 22. 爬虫数据分析服务接口
- **接口路径**: `src/main/java/com/security/service/CrawlDataAnalysisService.java`
- **功能**: 分析爬取的数据。
- **方法**:
  - `void analyzeCrawlData(String url);` - 分析爬取的数据
  - `String generateCrawlReport(String url);` - 生成爬取报告