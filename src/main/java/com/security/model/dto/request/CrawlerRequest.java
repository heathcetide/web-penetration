package com.security.model.dto.request;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.Max;

public class CrawlerRequest {
    @NotBlank(message = "URL不能为空")
    private String url;
    
    @Max(value = 5, message = "最大爬取深度不能超过5")
    private Integer depth = 1;
    
    private Boolean includeExternal = false;  // 是否包含外部链接
    private Boolean collectSensitive = true;  // 是否收集敏感信息
    private Boolean analyzeJs = true;         // 是否分析JS文件

    public @NotBlank(message = "URL不能为空") String getUrl() {
        return url;
    }

    public void setUrl(@NotBlank(message = "URL不能为空") String url) {
        this.url = url;
    }

    public @Max(value = 5, message = "最大爬取深度不能超过5") Integer getDepth() {
        return depth;
    }

    public void setDepth(@Max(value = 5, message = "最大爬取深度不能超过5") Integer depth) {
        this.depth = depth;
    }

    public Boolean getIncludeExternal() {
        return includeExternal;
    }

    public void setIncludeExternal(Boolean includeExternal) {
        this.includeExternal = includeExternal;
    }

    public Boolean getCollectSensitive() {
        return collectSensitive;
    }

    public void setCollectSensitive(Boolean collectSensitive) {
        this.collectSensitive = collectSensitive;
    }

    public Boolean getAnalyzeJs() {
        return analyzeJs;
    }

    public void setAnalyzeJs(Boolean analyzeJs) {
        this.analyzeJs = analyzeJs;
    }
}