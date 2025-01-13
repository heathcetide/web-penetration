package com.security.model.dto.request;


import javax.validation.constraints.NotBlank;

public class AssetScanRequest {
    @NotBlank(message = "域名不能为空")
    private String domain;
    
    private Integer scanType; // 1-端口扫描 2-子域名扫描 3-目录扫描
    
    private String scanParams; // 扫描参数，JSON格式

    public String getDomain() {
        return domain;
    }

    public void setDomain(String domain) {
        this.domain = domain;
    }

    public Integer getScanType() {
        return scanType;
    }

    public void setScanType(Integer scanType) {
        this.scanType = scanType;
    }

    public String getScanParams() {
        return scanParams;
    }

    public void setScanParams(String scanParams) {
        this.scanParams = scanParams;
    }
}