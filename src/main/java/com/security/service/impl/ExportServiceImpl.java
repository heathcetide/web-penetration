package com.security.service.impl;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.security.model.entity.CrawlResult;
import com.security.service.ExportService;
import com.opencsv.CSVWriter;
import java.io.StringWriter;
import java.util.List;
import java.util.Map;
public class ExportServiceImpl implements ExportService {
    @Override
    public String exportToJson(CrawlResult crawlResult) {
        try {
            ObjectMapper objectMapper = new ObjectMapper();
            // 创建一个Map来存储爬虫结果的各个部分
            Map<String, List<String>> resultMap = Map.of(
                    "baseUrl", List.of(crawlResult.getBaseUrl()),
                    "links", crawlResult.getLinks(),
                    "forms", crawlResult.getForms(),
                    "images", crawlResult.getImages(),
                    "videos", crawlResult.getVideos(),
                    "resources", crawlResult.getResources()
            );
            return objectMapper.writeValueAsString(resultMap);
        } catch (Exception e) {
            System.err.println("导出JSON失败：" + e.getMessage());
            return null;
        }
    }

    @Override
    public String exportToCsv(CrawlResult crawlResult) {
        try (StringWriter writer = new StringWriter(); CSVWriter csvWriter = new CSVWriter(writer)) {
            // 写入表头
            String[] header = {"Type", "URL"};
            csvWriter.writeNext(header);

            // 写入BaseUrl
            csvWriter.writeNext(new String[]{"BaseUrl", crawlResult.getBaseUrl()});

            // 写入Links
            for (String link : crawlResult.getLinks()) {
                csvWriter.writeNext(new String[]{"Link", link});
            }

            // 写入Forms
            for (String form : crawlResult.getForms()) {
                csvWriter.writeNext(new String[]{"Form", form});
            }

            // 写入Images
            for (String image : crawlResult.getImages()) {
                csvWriter.writeNext(new String[]{"Image", image});
            }

            // 写入Videos
            for (String video : crawlResult.getVideos()) {
                csvWriter.writeNext(new String[]{"Video", video});
            }

            // 写入Resources
            for (String resource : crawlResult.getResources()) {
                csvWriter.writeNext(new String[]{"Resource", resource});
            }

            return writer.toString();
        } catch (Exception e) {
            System.err.println("导出CSV失败：" + e.getMessage());
            return null;
        }
    }
}
