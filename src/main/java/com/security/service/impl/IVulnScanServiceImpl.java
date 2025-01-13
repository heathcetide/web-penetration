package com.security.service.impl;

import com.security.service.IVulnScanService;
import org.springframework.stereotype.Service;

public class IVulnScanServiceImpl implements IVulnScanService {
    @Override
    public Long createVulnScanTask(Long targetId) {
        return 0L;
    }

    @Override
    public void sqlInjectionScan(Long taskId, String url) {

    }

    @Override
    public void xssScan(Long taskId, String url) {

    }

    @Override
    public void csrfScan(Long taskId, String url) {

    }
}
