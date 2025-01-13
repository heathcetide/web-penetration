package com.security.service.video.security;

import com.security.service.video.EncryptionInfo;
import org.springframework.stereotype.Component;
import java.util.HashMap;
import java.util.Map;

@Component
public class EncryptionStrengthAnalyzer {
    
    private static final Map<String, Integer> ENCRYPTION_SCORES = new HashMap<>();
    
    static {
        // DRM分数
        ENCRYPTION_SCORES.put("WIDEVINE_L1", 100);
        ENCRYPTION_SCORES.put("WIDEVINE_L2", 80);
        ENCRYPTION_SCORES.put("WIDEVINE_L3", 60);
        ENCRYPTION_SCORES.put("PLAYREADY", 70);
        
        // 加密算法分数
        ENCRYPTION_SCORES.put("AES-256-CBC", 90);
        ENCRYPTION_SCORES.put("AES-128-CBC", 70);
        ENCRYPTION_SCORES.put("AES-128-ECB", 50);
        
        // 其他保护措施分数
        ENCRYPTION_SCORES.put("TOKEN", 30);
        ENCRYPTION_SCORES.put("REFERER", 20);
        ENCRYPTION_SCORES.put("TIME_EXPIRY", 25);
    }
    
    public SecurityAssessment analyzeEncryptionStrength(EncryptionInfo info) {
        SecurityAssessment assessment = new SecurityAssessment();
        int totalScore = 0;
        
        // 评估DRM保护
        if (info.isHasWidevineDRM()) {
            totalScore += ENCRYPTION_SCORES.get("WIDEVINE_L1");
            assessment.addProtection("Widevine DRM");
        }
        
        if (info.isHasPlayReadyDRM()) {
            totalScore += ENCRYPTION_SCORES.get("PLAYREADY");
            assessment.addProtection("PlayReady DRM");
        }
        
        // 评估内容加密
        if (info.isHasContentEncryption()) {
            String method = info.getEncryptionMethod();
            totalScore += ENCRYPTION_SCORES.getOrDefault(method, 50);
            assessment.addProtection("Content Encryption: " + method);
        }
        
        // 评估其他保护措施
        if (info.isHasToken()) {
            totalScore += ENCRYPTION_SCORES.get("TOKEN");
            assessment.addProtection("Token Authentication");
        }
        
        if (info.isHasRefererCheck()) {
            totalScore += ENCRYPTION_SCORES.get("REFERER");
            assessment.addProtection("Referer Check");
        }
        
        // 计算最终得分
        assessment.setScore(Math.min(100, totalScore));
        assessment.setLevel(determineSecurityLevel(assessment.getScore()));
        
        return assessment;
    }
    
    private SecurityLevel determineSecurityLevel(int score) {
        if (score >= 90) return SecurityLevel.VERY_HIGH;
        if (score >= 70) return SecurityLevel.HIGH;
        if (score >= 50) return SecurityLevel.MEDIUM;
        if (score >= 30) return SecurityLevel.LOW;
        return SecurityLevel.VERY_LOW;
    }
} 