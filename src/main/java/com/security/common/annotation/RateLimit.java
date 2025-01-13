package com.security.common.annotation;

import java.lang.annotation.*;
import java.util.concurrent.TimeUnit;

@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface RateLimit {
    /**
     * 限流key
     */
    String key() default "";
    
    /**
     * 时间窗口，默认1秒
     */
    long time() default 1;
    
    /**
     * 时间单位，默认秒
     */
    TimeUnit unit() default TimeUnit.SECONDS;
    
    /**
     * 限制次数
     */
    int count() default 100;
} 