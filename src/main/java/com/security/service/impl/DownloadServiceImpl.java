package com.security.service.impl;

import com.security.service.DownloadService;
import org.apache.commons.io.FileUtils;
import java.io.File;
import java.io.IOException;
import java.net.URL;


public class DownloadServiceImpl implements DownloadService {
    @Override
    public void downloadResource(String resourceUrl, String destinationPath) {
        try {
            // 创建URL对象，指向资源地址
            URL url = new URL(resourceUrl);
            // 创建File对象，指向本地目标路径
            File destinationFile = new File(destinationPath);

            // 使用Apache Commons IO库的FileUtils.copyURLToFile方法进行下载
            FileUtils.copyURLToFile(url, destinationFile);
            System.out.println("资源下载成功，保存路径：" + destinationPath);
        } catch (IOException e) {
            System.err.println("资源下载失败：" + e.getMessage());
        }
    }

    /*
    *public static void main(String[] args) {
        DownloadService downloadService = new DownloadService();
        // 示例：下载一个图片资源
        String imageUrl = "https://example.com/image.jpg";
        String localImagePath = "D:/downloaded/image.jpg";
        downloadService.downloadResource(imageUrl, localImagePath);
    }
     */
}
