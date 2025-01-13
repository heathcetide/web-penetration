package com.security.config;

import java.util.List;
import java.util.ArrayList;

public class FilterConfig {
    private List<String> allowedExtensions = new ArrayList<>();
    private List<String> blockedExtensions = new ArrayList<>();
    private Long maxFileSize;
    private boolean skipDuplicates;
    private boolean validateContent;
    
    // getters and setters
    public List<String> getAllowedExtensions() {
        return allowedExtensions;
    }
    
    public void setAllowedExtensions(List<String> allowedExtensions) {
        this.allowedExtensions = allowedExtensions;
    }
    
    public List<String> getBlockedExtensions() {
        return blockedExtensions;
    }
    
    public void setBlockedExtensions(List<String> blockedExtensions) {
        this.blockedExtensions = blockedExtensions;
    }
    
    public Long getMaxFileSize() {
        return maxFileSize;
    }
    
    public void setMaxFileSize(Long maxFileSize) {
        this.maxFileSize = maxFileSize;
    }
    
    public boolean isSkipDuplicates() {
        return skipDuplicates;
    }
    
    public void setSkipDuplicates(boolean skipDuplicates) {
        this.skipDuplicates = skipDuplicates;
    }
    
    public boolean isValidateContent() {
        return validateContent;
    }
    
    public void setValidateContent(boolean validateContent) {
        this.validateContent = validateContent;
    }
} 