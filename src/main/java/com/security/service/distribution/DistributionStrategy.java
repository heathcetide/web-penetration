package com.security.service.distribution;

public interface DistributionStrategy {
    boolean shouldDistributeTo(StorageProvider provider);
} 