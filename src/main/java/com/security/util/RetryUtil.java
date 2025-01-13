package com.security.util;


import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.function.Supplier;

public class RetryUtil {

    private static Logger log = LoggerFactory.getLogger(RetryUtil.class);

    public static <T> T retry(Supplier<T> supplier, int maxRetries, long delayMillis, String operation) {
        int retryCount = 0;
        while (true) {
            try {
                return supplier.get();
            } catch (Exception e) {
                retryCount++;
                if (retryCount > maxRetries) {
                    throw new RuntimeException(operation + "失败，已重试" + maxRetries + "次", e);
                }
                log.warn("{}失败，第{}次重试, error: {}", operation, retryCount, e.getMessage());
                try {
                    Thread.sleep(delayMillis * retryCount);
                } catch (InterruptedException ie) {
                    Thread.currentThread().interrupt();
                    throw new RuntimeException(ie);
                }
            }
        }
    }
} 