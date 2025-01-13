package com.security.common.constant;

public class RedisKeyConstants {
    private static final String PREFIX = "web_scanner:";
    
    public static final String TASK_RESULT_KEY = PREFIX + "task:result:";
    public static final String TASK_LOCK_KEY = PREFIX + "task:lock:";
    public static final String SCAN_RATE_LIMIT = PREFIX + "rate_limit:";
    
    public static String getTaskResultKey(Long taskId) {
        return TASK_RESULT_KEY + taskId;
    }
    
    public static String getTaskLockKey(Long taskId) {
        return TASK_LOCK_KEY + taskId;
    }
    
    public static String getScanRateLimitKey(String target) {
        return SCAN_RATE_LIMIT + target;
    }
} 