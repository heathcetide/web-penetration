package com.security.service.compress;

import org.apache.commons.compress.archivers.zip.ZipArchiveEntry;
import org.apache.commons.compress.archivers.zip.ZipArchiveOutputStream;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.io.*;
import java.util.List;

@Service
public class ResourceCompressor {

    private static Logger log = LoggerFactory.getLogger(ResourceCompressor.class);
    
    public void compressFiles(List<File> files, String outputPath) {
        try (ZipArchiveOutputStream zipOut = new ZipArchiveOutputStream(new File(outputPath))) {
            for (File file : files) {
                addToZip(file, file.getName(), zipOut);
            }
        } catch (IOException e) {
            log.error("压缩文件失败", e);
        }
    }
    
    private void addToZip(File file, String entryName, ZipArchiveOutputStream zipOut) 
            throws IOException {
        ZipArchiveEntry entry = new ZipArchiveEntry(file, entryName);
        zipOut.putArchiveEntry(entry);
        
        if (file.isFile()) {
            try (FileInputStream fis = new FileInputStream(file)) {
                byte[] buffer = new byte[8192];
                int bytesRead;
                while ((bytesRead = fis.read(buffer)) != -1) {
                    zipOut.write(buffer, 0, bytesRead);
                }
            }
        }
        
        zipOut.closeArchiveEntry();
    }
} 