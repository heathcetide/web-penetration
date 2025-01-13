package com.security.config;

import io.github.bonigarcia.wdm.WebDriverManager;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

@Configuration
public class WebDriverConfig {

    @Value("${selenium.config.timeout:60}")
    private int timeout;

    @Bean(destroyMethod = "quit")
    public WebDriver webDriver() {
        // 使用WebDriverManager自动管理ChromeDriver
        WebDriverManager.chromedriver().setup();
        
        ChromeOptions options = new ChromeOptions();
        
        // 启用无头模式
        options.addArguments("--headless=new");
        
        // 禁用音频
        options.addArguments("--mute-audio");
        options.addArguments("--autoplay-policy=no-user-gesture-required");
        
        // 允许跨域请求
        options.addArguments("--disable-web-security");
        options.addArguments("--allow-running-insecure-content");
        options.addArguments("--disable-features=IsolateOrigins,site-per-process");
        
        // 禁用开发者工具
        options.addArguments("--remote-debugging-port=0");
        options.addArguments("--remote-allow-origins=*");
        
        // 添加性能优化选项
        options.addArguments("--disable-gpu");                // 禁用GPU加速
        options.addArguments("--no-sandbox");                // 禁用沙箱模式
        options.addArguments("--disable-dev-shm-usage");     // 禁用/dev/shm使用
        options.addArguments("--disable-extensions");        // 禁用扩展
        options.addArguments("--disable-browser-side-navigation"); // 禁用浏览器侧导航
        options.addArguments("--disable-infobars");         // 禁用信息栏
        options.addArguments("--disable-notifications");    // 禁用通知
        options.addArguments("--disable-popup-blocking");   // 禁用弹出窗口阻止
        options.addArguments("--window-size=1920,1080");   // 设置窗口大小
        options.addArguments("--start-maximized");         // 最大化窗口
        options.addArguments("--ignore-certificate-errors"); // 忽略证书错误
        
        // 设置页面加载策略
        options.setPageLoadStrategy(org.openqa.selenium.PageLoadStrategy.EAGER);
        
        // 设置浏览器首选项
        Map<String, Object> prefs = new HashMap<>();
        prefs.put("profile.default_content_setting_values.notifications", 2);
        prefs.put("profile.default_content_settings.popups", 0);
        prefs.put("profile.default_content_setting_values.images", 2);  // 禁用图片加载
        prefs.put("profile.managed_default_content_settings.media_stream", 2);  // 禁用媒体流
        prefs.put("profile.managed_default_content_settings.media", 2);  // 禁用媒体
        options.setExperimentalOption("prefs", prefs);
        
        // 禁用日志
        options.setExperimentalOption("excludeSwitches", 
            new String[]{"enable-automation", "enable-logging"});
        
        try {
            ChromeDriver driver = new ChromeDriver(options);
            
            // 设置各种超时时间
            driver.manage().timeouts().pageLoadTimeout(Duration.ofSeconds(timeout));
            driver.manage().timeouts().scriptTimeout(Duration.ofSeconds(timeout));
            driver.manage().timeouts().implicitlyWait(Duration.ofSeconds(10));
            
            return driver;
        } catch (Exception e) {
            throw new RuntimeException("创建ChromeDriver失败: " + e.getMessage(), e);
        }
    }
} 