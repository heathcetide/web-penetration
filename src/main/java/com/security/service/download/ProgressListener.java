package com.security.service.download;

public interface ProgressListener {
    void start();
    void progress(long total, long current);
    void finish();
} 