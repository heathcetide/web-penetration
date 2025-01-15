package com.security.service;

/**
 * 下载服务接口，用于下载爬取到的资源。 [蒋浩天]
 */
public interface DownloadService {
        void downloadResource(String resourceUrl, String destinationPath);
}