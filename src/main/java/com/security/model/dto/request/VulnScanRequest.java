package com.security.model.dto.request;

import javax.validation.constraints.NotNull;

public class VulnScanRequest {
    @NotNull(message = "目标ID不能为空")
    private Long targetId;
    
    private Integer vulnType; // 1-SQL注入 2-XSS 3-CSRF
    
    private String scanParams; // 扫描参数，JSON格式

    public @NotNull(message = "目标ID不能为空") Long getTargetId() {
        return targetId;
    }

    public void setTargetId(@NotNull(message = "目标ID不能为空") Long targetId) {
        this.targetId = targetId;
    }

    public Integer getVulnType() {
        return vulnType;
    }

    public void setVulnType(Integer vulnType) {
        this.vulnType = vulnType;
    }

    public String getScanParams() {
        return scanParams;
    }

    public void setScanParams(String scanParams) {
        this.scanParams = scanParams;
    }
}