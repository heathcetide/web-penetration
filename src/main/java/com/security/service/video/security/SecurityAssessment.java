package com.security.service.video.security;

import java.util.ArrayList;
import java.util.List;

public class SecurityAssessment {
    private int score;                          // 安全得分
    private SecurityLevel level;                // 安全等级
    private List<String> protections;           // 保护措施列表
    private List<String> vulnerabilities;       // 漏洞列表
    private List<String> recommendations;       // 建议列表
    
    public SecurityAssessment() {
        this.protections = new ArrayList<>();
        this.vulnerabilities = new ArrayList<>();
        this.recommendations = new ArrayList<>();
    }
    
    public void addProtection(String protection) {
        this.protections.add(protection);
    }
    
    public void addVulnerability(String vulnerability) {
        this.vulnerabilities.add(vulnerability);
    }
    
    public void addRecommendation(String recommendation) {
        this.recommendations.add(recommendation);
    }

    public int getScore() {
        return score;
    }

    public void setScore(int score) {
        this.score = score;
    }

    public SecurityLevel getLevel() {
        return level;
    }

    public void setLevel(SecurityLevel level) {
        this.level = level;
    }

    public List<String> getProtections() {
        return protections;
    }

    public void setProtections(List<String> protections) {
        this.protections = protections;
    }

    public List<String> getVulnerabilities() {
        return vulnerabilities;
    }

    public void setVulnerabilities(List<String> vulnerabilities) {
        this.vulnerabilities = vulnerabilities;
    }

    public List<String> getRecommendations() {
        return recommendations;
    }

    public void setRecommendations(List<String> recommendations) {
        this.recommendations = recommendations;
    }
}