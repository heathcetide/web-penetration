package com.security.service.distribution;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.io.File;
import java.util.List;
import java.util.concurrent.CompletableFuture;

@Service
public class ResourceDistributor {

    private static Logger log = LoggerFactory.getLogger(ResourceDistributor.class);
    private final List<StorageProvider> storageProviders;
    
    public ResourceDistributor(List<StorageProvider> storageProviders) {
        this.storageProviders = storageProviders;
    }
    
    public void distribute(File file, DistributionStrategy strategy) {
        CompletableFuture<?>[] futures = storageProviders.stream()
            .filter(provider -> strategy.shouldDistributeTo(provider))
            .map(provider -> CompletableFuture.runAsync(() -> {
                try {
                    provider.store(file);
                } catch (Exception e) {
                    log.error("分发失败: {} -> {}", file.getName(), provider.getName(), e);
                }
            }))
            .toArray(CompletableFuture[]::new);
        
        CompletableFuture.allOf(futures).join();
    }
} 