package com.security.service.distribution;

import java.io.File;

public interface StorageProvider {
    String getName();
    void store(File file) throws Exception;
} 