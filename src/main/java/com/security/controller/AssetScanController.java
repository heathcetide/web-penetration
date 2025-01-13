package com.security.controller;

import com.security.model.params.Result;
import com.security.model.dto.request.AssetScanRequest;
import com.security.service.IAssetScanService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;

@RestController
@RequestMapping("/api/asset")
public class AssetScanController {
    @Autowired
    private IAssetScanService assetScanService;
    
    @PostMapping("/scan")
    public Result createScanTask(@RequestBody @Valid AssetScanRequest request) {
        Long taskId = assetScanService.createAssetScanTask(request.getDomain());
        return Result.success(taskId);
    }
    
    @GetMapping("/result/{taskId}")
    public Result getScanResult(@PathVariable Long taskId) {
        // 获取扫描结果
        return Result.success();
    }
} 