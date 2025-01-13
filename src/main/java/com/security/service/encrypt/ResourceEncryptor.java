package com.security.service.encrypt;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import javax.crypto.Cipher;
import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.nio.file.Files;
import java.security.MessageDigest;
import java.util.Base64;

@Service
public class ResourceEncryptor {
    
    private static final Logger log = LoggerFactory.getLogger(ResourceEncryptor.class);
    
    private final String encryptionKey;
    private final String algorithm;
    
    public ResourceEncryptor(
            @Value("${resource.video.encryption.custom.encryption-key}") String encryptionKey,
            @Value("${resource.video.encryption.custom.algorithm:AES}") String algorithm) {
        this.encryptionKey = encryptionKey;
        this.algorithm = algorithm;
    }
    
    /**
     * 加密文件
     */
    public void encryptFile(File sourceFile, File targetFile) {
        try {
            // 生成密钥
            SecretKey key = generateKey(encryptionKey);
            
            // 初始化加密器
            Cipher cipher = Cipher.getInstance(algorithm);
            cipher.init(Cipher.ENCRYPT_MODE, key);
            
            // 读取文件并加密
            try (FileInputStream in = new FileInputStream(sourceFile);
                 FileOutputStream out = new FileOutputStream(targetFile)) {
                
                byte[] buffer = new byte[8192];
                int len;
                while ((len = in.read(buffer)) != -1) {
                    byte[] encryptedData = cipher.update(buffer, 0, len);
                    if (encryptedData != null) {
                        out.write(encryptedData);
                    }
                }
                
                byte[] finalData = cipher.doFinal();
                if (finalData != null) {
                    out.write(finalData);
                }
            }
            
            log.info("文件加密成功: {} -> {}", sourceFile.getPath(), targetFile.getPath());
        } catch (Exception e) {
            log.error("文件加密失败: {}", e.getMessage(), e);
        }
    }
    
    /**
     * 解密文件
     */
    public void decryptFile(File sourceFile, File targetFile) {
        try {
            // 生成密钥
            SecretKey key = generateKey(encryptionKey);
            
            // 初始化解密器
            Cipher cipher = Cipher.getInstance(algorithm);
            cipher.init(Cipher.DECRYPT_MODE, key);
            
            // 读取文件并解密
            try (FileInputStream in = new FileInputStream(sourceFile);
                 FileOutputStream out = new FileOutputStream(targetFile)) {
                
                byte[] buffer = new byte[8192];
                int len;
                while ((len = in.read(buffer)) != -1) {
                    byte[] decryptedData = cipher.update(buffer, 0, len);
                    if (decryptedData != null) {
                        out.write(decryptedData);
                    }
                }
                
                byte[] finalData = cipher.doFinal();
                if (finalData != null) {
                    out.write(finalData);
                }
            }
            
            log.info("文件解密成功: {} -> {}", sourceFile.getPath(), targetFile.getPath());
        } catch (Exception e) {
            log.error("文件解密失败: {}", e.getMessage(), e);
        }
    }
    
    /**
     * 生成密钥
     */
    private SecretKey generateKey(String key) throws Exception {
        MessageDigest digest = MessageDigest.getInstance("SHA-256");
        byte[] hash = digest.digest(key.getBytes());
        return new SecretKeySpec(hash, algorithm);
    }
} 