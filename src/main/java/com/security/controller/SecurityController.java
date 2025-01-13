package com.security.controller;

import com.security.model.entity.ScanResult;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/security")
public class SecurityController {

    @PostMapping("/sql-injection")
    public ScanResult testSqlInjection(@RequestParam String url) {
        // 调用 SQL 注入服务
        return null; // 这里返回 ScanResult 对象
    }

    @PostMapping("/xss")
    public ScanResult testXss(@RequestParam String url) {
        // 调用 XSS 服务
        return null; // 这里返回 ScanResult 对象
    }

    @PostMapping("/csrf")
    public ScanResult testCsrf(@RequestParam String url) {
        // 调用 CSRF 服务
        return null; // 这里返回 ScanResult 对象
    }

    @GetMapping("/assets")
    public ScanResult collectAssets(@RequestParam String url) {
        // 调用资产收集服务
        return null; // 这里返回 ScanResult 对象
    }
} 