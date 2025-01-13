package com.security.model.entity;


import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class DownloadStats {
    private long bytesDownloaded;
    private long startTime;
    private long lastUpdateTime;
    private double currentSpeed; // bytes/s
    private int activeThreads;
    private Map<String, SegmentStatus> segmentStatus = new ConcurrentHashMap<>();

    public long getBytesDownloaded() {
        return bytesDownloaded;
    }

    public void setBytesDownloaded(long bytesDownloaded) {
        this.bytesDownloaded = bytesDownloaded;
    }

    public long getStartTime() {
        return startTime;
    }

    public void setStartTime(long startTime) {
        this.startTime = startTime;
    }

    public long getLastUpdateTime() {
        return lastUpdateTime;
    }

    public void setLastUpdateTime(long lastUpdateTime) {
        this.lastUpdateTime = lastUpdateTime;
    }

    public double getCurrentSpeed() {
        return currentSpeed;
    }

    public void setCurrentSpeed(double currentSpeed) {
        this.currentSpeed = currentSpeed;
    }

    public int getActiveThreads() {
        return activeThreads;
    }

    public void setActiveThreads(int activeThreads) {
        this.activeThreads = activeThreads;
    }

    public Map<String, SegmentStatus> getSegmentStatus() {
        return segmentStatus;
    }

    public void setSegmentStatus(Map<String, SegmentStatus> segmentStatus) {
        this.segmentStatus = segmentStatus;
    }
}


