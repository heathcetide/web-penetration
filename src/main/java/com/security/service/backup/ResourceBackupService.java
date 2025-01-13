package com.security.service.backup;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;

@Service
public class ResourceBackupService {
    
    private static final Logger log = LoggerFactory.getLogger(ResourceBackupService.class);
    
    private final String backupPath;
    
    public ResourceBackupService(@Value("${resource.backup.path:backup/}") String backupPath) {
        this.backupPath = backupPath;
        // 确保备份目录存在
        new File(backupPath).mkdirs();
    }
    
    /**
     * 备份资源文件
     */
    public void backup(File sourceFile) {
        try {
            // 创建备份文件路径
            String relativePath = sourceFile.getParent().replace(File.separator, "_");
            File backupFile = new File(backupPath, relativePath + "_" + sourceFile.getName());
            
            // 确保父目录存在
            backupFile.getParentFile().mkdirs();
            
            // 复制文件
            Files.copy(sourceFile.toPath(), backupFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
            
            log.info("资源备份成功: {} -> {}", sourceFile.getPath(), backupFile.getPath());
        } catch (Exception e) {
            log.error("资源备份失败: {}", e.getMessage(), e);
        }
    }
    
    /**
     * 从备份恢复文件
     */
    public boolean restore(String backupFileName, String targetPath) {
        try {
            File backupFile = new File(backupPath, backupFileName);
            if (!backupFile.exists()) {
                log.warn("备份文件不存在: {}", backupFile.getPath());
                return false;
            }
            
            File targetFile = new File(targetPath);
            targetFile.getParentFile().mkdirs();
            
            Files.copy(backupFile.toPath(), targetFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
            
            log.info("资源恢复成功: {} -> {}", backupFile.getPath(), targetFile.getPath());
            return true;
        } catch (Exception e) {
            log.error("资源恢复失败: {}", e.getMessage(), e);
            return false;
        }
    }
} 