package com.security.service.download;

import com.security.model.entity.DownloadProgress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.concurrent.ConcurrentHashMap;

@Component
public class ResumableDownloader {

    private static Logger log = LoggerFactory.getLogger(ResumableDownloader.class);
    private final ConcurrentHashMap<String, DownloadProgress> progressMap = new ConcurrentHashMap<>();
    
    public void download(String url, String savePath, DownloadProgress progress) {
        File file = new File(savePath);
        long downloadedSize = file.exists() ? file.length() : 0;
        
        try {
            HttpURLConnection conn = (HttpURLConnection) new URL(url).openConnection();
            
            // 支持断点续传
            if (downloadedSize > 0) {
                conn.setRequestProperty("Range", "bytes=" + downloadedSize + "-");
            }
            
            // 执行下载
            if (conn.getResponseCode() == HttpURLConnection.HTTP_PARTIAL) {
                resumeDownload(conn, file, downloadedSize, progress);
            } else {
                normalDownload(conn, file, progress);
            }
            
        } catch (Exception e) {
            log.error("下载失败", e);
            progress.setStatus(3); // 失败状态
        }
    }
    
    private void resumeDownload(HttpURLConnection conn, File file, long downloadedSize, 
                              DownloadProgress progress) throws IOException {
        long totalSize = downloadedSize + conn.getContentLengthLong();
        progress.setTotalSize(totalSize);
        
        try (InputStream in = conn.getInputStream();
             RandomAccessFile raf = new RandomAccessFile(file, "rw")) {
            
            raf.seek(downloadedSize);
            byte[] buffer = new byte[8192];
            int bytesRead;
            
            while ((bytesRead = in.read(buffer)) != -1) {
                raf.write(buffer, 0, bytesRead);
                downloadedSize += bytesRead;
                updateProgress(progress, downloadedSize);
            }
        }
    }
    
    private void normalDownload(HttpURLConnection conn, File file, DownloadProgress progress) 
            throws IOException {
        progress.setTotalSize(conn.getContentLengthLong());
        
        try (InputStream in = conn.getInputStream();
             FileOutputStream out = new FileOutputStream(file)) {
            
            byte[] buffer = new byte[8192];
            int bytesRead;
            long downloadedSize = 0;
            
            while ((bytesRead = in.read(buffer)) != -1) {
                out.write(buffer, 0, bytesRead);
                downloadedSize += bytesRead;
                updateProgress(progress, downloadedSize);
            }
        }
    }
    
    private void updateProgress(DownloadProgress progress, long downloadedSize) {
        progress.setDownloadedSize(downloadedSize);
        progress.setProgress((int) (downloadedSize * 100 / progress.getTotalSize()));
        progressMap.put(progress.getUrl(), progress);
    }
} 