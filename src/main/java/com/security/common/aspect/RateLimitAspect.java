package com.security.common.aspect;

import com.security.common.annotation.RateLimit;
import com.security.exception.BusinessException;
import com.security.util.RedisUtil;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Aspect
@Component
public class RateLimitAspect {
    
    @Autowired
    private RedisUtil redisUtil;
    
    @Around("@annotation(rateLimit)")
    public Object around(ProceedingJoinPoint point, RateLimit rateLimit) throws Throwable {
        String key = rateLimit.key();
        if (key.isEmpty()) {
            MethodSignature signature = (MethodSignature) point.getSignature();
            key = signature.getMethod().getName();
        }
        
        long count = redisUtil.increment(key, 1);
        if (count == 1) {
            redisUtil.expire(key, rateLimit.time(), rateLimit.unit());
        }
        
        if (count > rateLimit.count()) {
            throw new BusinessException(429, "请求过于频繁，请稍后重试");
        }
        
        return point.proceed();
    }
} 