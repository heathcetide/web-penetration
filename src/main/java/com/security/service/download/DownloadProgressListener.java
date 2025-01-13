package com.security.service.download;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import com.security.service.download.ProgressListener;

public class DownloadProgressListener implements ProgressListener {
    private final Long taskId;
    private final String url;
    private long totalSize;
    private long downloadedSize;
    private static Logger log = LoggerFactory.getLogger(DownloadProgressListener.class);
    public DownloadProgressListener(Long taskId, String url) {
        this.taskId = taskId;
        this.url = url;
    }
    
    @Override
    public void start() {
        log.info("开始下载: taskId={}, url={}", taskId, url);
    }
    
    @Override
    public void progress(long total, long current) {
        this.totalSize = total;
        this.downloadedSize = current;
        int progress = (int) (current * 100 / total);
        log.info("下载进度: taskId={}, url={}, progress={}%", taskId, url, progress);
    }
    
    @Override
    public void finish() {
        log.info("下载完成: taskId={}, url={}", taskId, url);
    }
}