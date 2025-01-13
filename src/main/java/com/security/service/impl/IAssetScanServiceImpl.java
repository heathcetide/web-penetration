package com.security.service.impl;

import com.security.service.IAssetScanService;
import org.springframework.stereotype.Service;

@Service
public class IAssetScanServiceImpl implements IAssetScanService {
    @Override
    public Long createAssetScanTask(String domain) {
        return 0L;
    }

    @Override
    public void portScan(Long taskId, String target) {

    }

    @Override
    public void subdomainScan(Long taskId, String domain) {

    }

    @Override
    public void directoryScan(Long taskId, String url) {

    }
}
