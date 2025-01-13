package com.security.service.retry;

import java.util.function.Predicate;

/**
 * 重试策略配置
 */
public class RetryStrategy {
    
    /**
     * 最大重试次数
     */
    private final int maxAttempts;
    
    /**
     * 初始延迟时间(毫秒)
     */
    private final long initialDelay;
    
    /**
     * 退避乘数
     */
    private final double backoffMultiplier;
    
    /**
     * 重试条件判断
     */
    private final Predicate<Exception> retryPredicate;
    
    /**
     * 构造函数
     */
    public RetryStrategy(int maxAttempts, long initialDelay, double backoffMultiplier, 
            Predicate<Exception> retryPredicate) {
        this.maxAttempts = maxAttempts;
        this.initialDelay = initialDelay;
        this.backoffMultiplier = backoffMultiplier;
        this.retryPredicate = retryPredicate;
    }
    
    /**
     * 默认构造函数
     */
    public RetryStrategy() {
        this(3, 1000, 2.0, e -> true);
    }
    
    public int getMaxAttempts() {
        return maxAttempts;
    }
    
    public long getInitialDelay() {
        return initialDelay;
    }
    
    public double getBackoffMultiplier() {
        return backoffMultiplier;
    }
    
    public Predicate<Exception> getRetryPredicate() {
        return retryPredicate;
    }
    
    /**
     * 计算第n次重试的延迟时间
     */
    public long getDelayForAttempt(int attempt) {
        return (long) (initialDelay * Math.pow(backoffMultiplier, attempt - 1));
    }
    
    /**
     * 判断是否需要重试
     */
    public boolean shouldRetry(Exception e) {
        return retryPredicate.test(e);
    }
} 