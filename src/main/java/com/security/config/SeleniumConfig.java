package com.security.config;

import io.github.bonigarcia.wdm.WebDriverManager;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import javax.annotation.PostConstruct;
import java.util.Arrays;
import java.util.logging.Level;
import org.openqa.selenium.logging.LoggingPreferences;
import org.openqa.selenium.logging.LogType;

@Configuration
public class SeleniumConfig {
    
    private static Logger log = LoggerFactory.getLogger(SeleniumConfig.class);
    
    @PostConstruct
    void setup() {
        try {
            // 自动下载和配置 ChromeDriver
            WebDriverManager.chromedriver().setup();
            log.info("ChromeDriver setup completed");
        } catch (Exception e) {
            log.error("Failed to setup ChromeDriver: {}", e.getMessage(), e);
            throw new IllegalStateException("Failed to setup ChromeDriver", e);
        }
    }
    
    @Bean(destroyMethod = "quit")
    public WebDriver chromeDriver() {
        try {
            ChromeOptions options = new ChromeOptions();
            options.addArguments("--headless");
            options.addArguments("--disable-gpu");
            options.addArguments("--no-sandbox");
            options.addArguments("--disable-dev-shm-usage");
            options.addArguments("--remote-allow-origins=*");
            
            // 启用性能日志
            LoggingPreferences logPrefs = new LoggingPreferences();
            logPrefs.enable(LogType.PERFORMANCE, Level.ALL);
            options.setCapability("goog:loggingPrefs", logPrefs);
            
            // 添加更多配置以绕过一些网站的检测
            options.addArguments("--disable-blink-features=AutomationControlled");
            options.setExperimentalOption("excludeSwitches", Arrays.asList("enable-automation"));
            
            WebDriver driver = new ChromeDriver(options);
            log.info("ChromeDriver initialized successfully");
            return driver;
            
        } catch (Exception e) {
            log.error("Failed to create ChromeDriver: {}", e.getMessage(), e);
            throw new IllegalStateException("Failed to create ChromeDriver", e);
        }
    }
} 