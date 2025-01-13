package com.security.service;

/**
 * 静态代码分析服务接口，用于对 Web 应用的源代码进行静态分析。[刘铭昊]
 */
public interface StaticCodeAnalysisService {

    /**
     * 对源代码执行静态代码分析。
     *
     * @param sourceCode 要分析的源代码
     * @return 如果代码符合安全要求返回 true，否则返回 false
     */
    boolean performStaticCodeAnalysis(String sourceCode);

    /**
     * 获取源代码的分析报告。
     *
     * @param sourceCode 要分析的源代码
     * @return 分析报告的详细内容
     */
    String getAnalysisReport(String sourceCode);

    /**
     * 检查源代码中是否存在常见的安全漏洞。
     *
     * @param sourceCode 要分析的源代码
     * @return 如果存在常见漏洞返回 true，否则返回 false
     */
    boolean checkForCommonVulnerabilities(String sourceCode);

    /**
     * 验证源代码是否符合指定的代码标准。
     *
     * @param sourceCode 要分析的源代码
     * @return 如果符合标准返回 true，否则返回 false
     */
    boolean validateCodeStandards(String sourceCode);

    /**
     * 分析源代码的复杂性。
     *
     * @param sourceCode 要分析的源代码
     * @return 如果代码复杂性较高返回 true，否则返回 false
     */
    boolean analyzeCodeComplexity(String sourceCode);

    /**
     * 识别源代码中的死代码（未被使用的代码）。
     *
     * @param sourceCode 要分析的源代码
     * @return 如果存在死代码返回 true，否则返回 false
     */
    boolean identifyDeadCode(String sourceCode);
}
