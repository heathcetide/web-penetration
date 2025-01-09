package scan

import "regexp"

// ServiceFingerprints 服务指纹库
var ServiceFingerprints = map[string][]*ServiceFingerprint{
    "http": {
        {
            Service: "Apache",
            Pattern: regexp.MustCompile(`Apache(?:/(\d+[\.\d]+[^\s]*)|[^\s]*)`),
        },
        {
            Service: "Nginx",
            Pattern: regexp.MustCompile(`nginx(?:/(\d+[\.\d]+))?`),
        },
    },
    "ssh": {
        {
            Service: "OpenSSH",
            Pattern: regexp.MustCompile(`OpenSSH[_-](\d+[\.\d]+)`),
        },
    },
    "ftp": {
        {
            Service: "vsftpd",
            Pattern: regexp.MustCompile(`vsftpd\s+([\d\.]+)`),
        },
        {
            Service: "ProFTPD",
            Pattern: regexp.MustCompile(`ProFTPD\s+([\d\.]+)`),
        },
    },
    // 添加更多服务指纹
}

// CommonPorts 常用端口服务映射
var CommonPorts = map[int]string{
    21:   "ftp",
    22:   "ssh",
    23:   "telnet",
    25:   "smtp",
    53:   "dns",
    80:   "http",
    110:  "pop3",
    143:  "imap",
    443:  "https",
    3306: "mysql",
    5432: "postgresql",
    6379: "redis",
    // 添加更多端口映射
} 