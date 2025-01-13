package com.security.service.download;

import com.google.common.util.concurrent.RateLimiter;

import java.io.FilterInputStream;
import java.io.IOException;
import java.io.InputStream;

public class RateLimitInputStream extends FilterInputStream {
    private final RateLimiter rateLimiter;
    
    public RateLimitInputStream(InputStream in, RateLimiter rateLimiter) {
        super(in);
        this.rateLimiter = rateLimiter;
    }
    
    @Override
    public int read(byte[] b, int off, int len) throws IOException {
        int bytesRead = super.read(b, off, len);
        if (bytesRead > 0) {
            rateLimiter.acquire(bytesRead);
        }
        return bytesRead;
    }
} 