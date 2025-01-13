package com.security.service.impl;

import com.security.service.XssDetectionService;

public class XssDetectionServiceImpl implements XssDetectionService {
    @Override
    public boolean detectXss(String url, String payload) {
        //TODO 需要实现
        return false;
    }

    @Override
    public boolean testForStoredXss(String url) {
        // TODO 需要实现
        return false;
    }
}
