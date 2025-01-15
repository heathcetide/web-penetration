package com.security.service.processor;

public class ScanResultProcessor {   //蒋浩天无中生有，用来表示扫描结果
    private boolean isVulnerable;
    private String message;

    public ScanResultProcessor(boolean isVulnerable, String message) {
        this.isVulnerable = isVulnerable;
        this.message = message;
    }

    public boolean isVulnerable() {
        return isVulnerable;
    }

    public String getMessage() {
        return message;
    }

    @Override
    public String toString() {
        return "ScanResult{" +
                "isVulnerable=" + isVulnerable +
                ", message='" + message + '\'' +
                '}';
    }
}
