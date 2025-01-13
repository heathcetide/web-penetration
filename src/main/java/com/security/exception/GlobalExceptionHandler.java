package com.security.exception;

import com.security.model.params.Result;
import org.springframework.validation.BindException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
public class GlobalExceptionHandler {
    
    @ExceptionHandler(BusinessException.class)
    public Result handleBusinessException(BusinessException e) {
//        log.error("业务异常：{}", e.getMessage());
        return Result.error(e.getCode(), e.getMessage());
    }
    
    @ExceptionHandler(BindException.class)
    public Result handleBindException(BindException e) {
//        log.error("参数校验异常：{}", e.getBindingResult().getAllErrors().get(0).getDefaultMessage());
        return Result.error(400, e.getBindingResult().getAllErrors().get(0).getDefaultMessage());
    }
    
    @ExceptionHandler(Exception.class)
    public Result handleException(Exception e) {
//        log.error("系统异常：", e);
        return Result.error(500, "系统异常，请联系管理员");
    }
} 