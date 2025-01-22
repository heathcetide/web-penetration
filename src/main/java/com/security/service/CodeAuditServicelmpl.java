package com.security.service;

import net.sourceforge.pmd.PMD;
import net.sourceforge.pmd.PMDConfiguration;
import net.sourceforge.pmd.Report;
import net.sourceforge.pmd.RuleSet;
import net.sourceforge.pmd.RuleSetFactory;
import org.springframework.stereotype.Service;
import java.io.File;


@Service
public class CodeAuditServiceImpl implements CodeAuditService {

    @Override
    public boolean auditCode(String sourceCode) {
        try {
            // 创建 PMD 配置对象
            PMDConfiguration configuration = new PMDConfiguration();

            // 创建规则集工厂
            RuleSetFactory ruleSetFactory = new RuleSetFactory(configuration);

            // 加载默认规则集
            RuleSet ruleSet = ruleSetFactory.createRuleSet(new File("path/to/ruleset.xml"));

            // 创建 PMD 对象
            PMD pmd = new PMD(configuration);

            // 分析代码
            Report report = pmd.analyze(new File(sourceCode), ruleSet);

            // 检查报告是否包含违规
            return report.getViolationCount() > 0;
        } catch (Exception e) {
            // 处理异常
            e.printStackTrace();
            return false;
        }
    }
}
