package com.security.controller;

import com.security.model.params.Result;
import com.security.model.dto.request.CrawlerRequest;
import com.security.service.IWebCrawlerService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;

@RestController
@RequestMapping("/api/crawler")
public class WebCrawlerController {
    
    @Autowired
    private IWebCrawlerService webCrawlerService;
    
    @PostMapping("/start")
//    @RateLimit(key = "crawler", time = 1, count = 10)
    public Result startCrawler(@RequestBody @Valid CrawlerRequest request) {
        Long taskId = webCrawlerService.createCrawlerTask(request.getUrl(), request.getDepth());
        return Result.success(taskId);
    }
    
    @GetMapping("/result/{taskId}")
    public Result getCrawlerResult(@PathVariable Long taskId) {
        return Result.success(webCrawlerService.getCrawlerResult(taskId));
    }
} 