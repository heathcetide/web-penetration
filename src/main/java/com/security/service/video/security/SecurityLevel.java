package com.security.service.video.security;

public enum SecurityLevel {
    VERY_HIGH("非常高", "使用了多层高强度加密保护"),
    HIGH("高", "使用了标准的加密保护"),
    MEDIUM("中等", "使用了基本的保护措施"),
    LOW("低", "保护措施较弱"),
    VERY_LOW("非常低", "几乎没有保护措施");
    
    private final String name;
    private final String description;
    
    SecurityLevel(String name, String description) {
        this.name = name;
        this.description = description;
    }
    
    public String getName() {
        return name;
    }
    
    public String getDescription() {
        return description;
    }
} 