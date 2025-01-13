package com.security.service.video;

import java.util.Map;
import java.util.List;

public class EncryptionInfo {
    private boolean hasToken;           // 是否有token
    private String tokenType;           // token类型
    private String tokenValue;          // token值
    
    private boolean hasDRM;             // 是否有DRM保护
    private String drmType;             // DRM类型
    
    private boolean hasHLSEncryption;   // 是否有HLS加密
    private String keyUrl;              // 密钥URL
    
    private boolean hasContentEncryption; // 是否有内容加密
    private String encryptionMethod;      // 加密方法
    
    private boolean hasRefererCheck;    // 是否有防盗链
    private boolean hasMultipleQualities; // 是否有多个清晰度
    
    private String error;               // 错误信息
    
    private boolean hasCustomEncryption;  // 自定义加密
    private String customEncryptionKey;   // 自定义加密密钥
    
    private boolean hasTimeExpiry;        // 时间戳过期
    private Long expiryTime;              // 过期时间
    
    private boolean hasDASH;              // 是否DASH流
    private boolean hasCommonEncryption;  // 通用加密
    private boolean hasWidevineDRM;       // Widevine DRM
    private boolean hasPlayReadyDRM;      // PlayReady DRM
    
    private Map<String, String> extraParams;  // 额外参数
    private List<String> encryptionLayers;    // 加密层级
    
    public boolean isEncrypted() {
        return hasToken || hasDRM || hasHLSEncryption || 
               hasContentEncryption || hasCustomEncryption ||
               hasTimeExpiry || hasCommonEncryption;
    }

    public boolean isHasToken() {
        return hasToken;
    }

    public void setHasToken(boolean hasToken) {
        this.hasToken = hasToken;
    }

    public String getTokenType() {
        return tokenType;
    }

    public void setTokenType(String tokenType) {
        this.tokenType = tokenType;
    }

    public String getTokenValue() {
        return tokenValue;
    }

    public void setTokenValue(String tokenValue) {
        this.tokenValue = tokenValue;
    }

    public boolean isHasDRM() {
        return hasDRM;
    }

    public void setHasDRM(boolean hasDRM) {
        this.hasDRM = hasDRM;
    }

    public String getDrmType() {
        return drmType;
    }

    public void setDrmType(String drmType) {
        this.drmType = drmType;
    }

    public boolean isHasHLSEncryption() {
        return hasHLSEncryption;
    }

    public void setHasHLSEncryption(boolean hasHLSEncryption) {
        this.hasHLSEncryption = hasHLSEncryption;
    }

    public String getKeyUrl() {
        return keyUrl;
    }

    public void setKeyUrl(String keyUrl) {
        this.keyUrl = keyUrl;
    }

    public boolean isHasContentEncryption() {
        return hasContentEncryption;
    }

    public void setHasContentEncryption(boolean hasContentEncryption) {
        this.hasContentEncryption = hasContentEncryption;
    }

    public String getEncryptionMethod() {
        return encryptionMethod;
    }

    public void setEncryptionMethod(String encryptionMethod) {
        this.encryptionMethod = encryptionMethod;
    }

    public boolean isHasRefererCheck() {
        return hasRefererCheck;
    }

    public void setHasRefererCheck(boolean hasRefererCheck) {
        this.hasRefererCheck = hasRefererCheck;
    }

    public boolean isHasMultipleQualities() {
        return hasMultipleQualities;
    }

    public void setHasMultipleQualities(boolean hasMultipleQualities) {
        this.hasMultipleQualities = hasMultipleQualities;
    }

    public String getError() {
        return error;
    }

    public void setError(String error) {
        this.error = error;
    }

    public boolean isHasCustomEncryption() {
        return hasCustomEncryption;
    }

    public void setHasCustomEncryption(boolean hasCustomEncryption) {
        this.hasCustomEncryption = hasCustomEncryption;
    }

    public String getCustomEncryptionKey() {
        return customEncryptionKey;
    }

    public void setCustomEncryptionKey(String customEncryptionKey) {
        this.customEncryptionKey = customEncryptionKey;
    }

    public boolean isHasTimeExpiry() {
        return hasTimeExpiry;
    }

    public void setHasTimeExpiry(boolean hasTimeExpiry) {
        this.hasTimeExpiry = hasTimeExpiry;
    }

    public Long getExpiryTime() {
        return expiryTime;
    }

    public void setExpiryTime(Long expiryTime) {
        this.expiryTime = expiryTime;
    }

    public boolean isHasDASH() {
        return hasDASH;
    }

    public void setHasDASH(boolean hasDASH) {
        this.hasDASH = hasDASH;
    }

    public boolean isHasCommonEncryption() {
        return hasCommonEncryption;
    }

    public void setHasCommonEncryption(boolean hasCommonEncryption) {
        this.hasCommonEncryption = hasCommonEncryption;
    }

    public boolean isHasWidevineDRM() {
        return hasWidevineDRM;
    }

    public void setHasWidevineDRM(boolean hasWidevineDRM) {
        this.hasWidevineDRM = hasWidevineDRM;
    }

    public boolean isHasPlayReadyDRM() {
        return hasPlayReadyDRM;
    }

    public void setHasPlayReadyDRM(boolean hasPlayReadyDRM) {
        this.hasPlayReadyDRM = hasPlayReadyDRM;
    }

    public Map<String, String> getExtraParams() {
        return extraParams;
    }

    public void setExtraParams(Map<String, String> extraParams) {
        this.extraParams = extraParams;
    }

    public List<String> getEncryptionLayers() {
        return encryptionLayers;
    }

    public void setEncryptionLayers(List<String> encryptionLayers) {
        this.encryptionLayers = encryptionLayers;
    }
}