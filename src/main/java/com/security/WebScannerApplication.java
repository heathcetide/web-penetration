package com.security;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@MapperScan("com.security.mapper")
public class WebScannerApplication {
    public static void main(String[] args) {
        SpringApplication.run(WebScannerApplication.class, args);
    }
} 