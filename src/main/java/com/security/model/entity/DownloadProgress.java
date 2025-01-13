package com.security.model.entity;

import com.baomidou.mybatisplus.annotation.TableName;


@TableName("download_progress")
public class DownloadProgress {
    private Long id;
    private Long taskId;
    private String url;
    private Long totalSize;
    private Long downloadedSize;
    private Integer progress;
    private Integer status; // 0-等待 1-下载中 2-已完成 3-失败 4-暂停
    private Long lastModified;
    private String etag;
    private String rangeSupport;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Long getTaskId() {
        return taskId;
    }

    public void setTaskId(Long taskId) {
        this.taskId = taskId;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public Long getTotalSize() {
        return totalSize;
    }

    public void setTotalSize(Long totalSize) {
        this.totalSize = totalSize;
    }

    public Long getDownloadedSize() {
        return downloadedSize;
    }

    public void setDownloadedSize(Long downloadedSize) {
        this.downloadedSize = downloadedSize;
    }

    public Integer getProgress() {
        return progress;
    }

    public void setProgress(Integer progress) {
        this.progress = progress;
    }

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public Long getLastModified() {
        return lastModified;
    }

    public void setLastModified(Long lastModified) {
        this.lastModified = lastModified;
    }

    public String getEtag() {
        return etag;
    }

    public void setEtag(String etag) {
        this.etag = etag;
    }

    public String getRangeSupport() {
        return rangeSupport;
    }

    public void setRangeSupport(String rangeSupport) {
        this.rangeSupport = rangeSupport;
    }
}