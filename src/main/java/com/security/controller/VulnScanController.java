package com.security.controller;

import com.security.model.params.Result;
import com.security.model.dto.request.VulnScanRequest;
import com.security.service.IVulnScanService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;

@RestController
@RequestMapping("/api/vuln")
public class VulnScanController {
    @Autowired
    private IVulnScanService vulnScanService;
    
    @PostMapping("/scan")
    public Result createVulnScan(@RequestBody @Valid VulnScanRequest request) {
        Long taskId = vulnScanService.createVulnScanTask(request.getTargetId());
        return Result.success(taskId);
    }
    
    @GetMapping("/result/{taskId}")
    public Result getVulnResult(@PathVariable Long taskId) {
        // 获取漏洞扫描结果
        return Result.success();
    }
} 