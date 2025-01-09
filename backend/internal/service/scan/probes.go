package scan

// ServiceProbe 服务探测规则
type ServiceProbe struct {
    Name     string   // 服务名称
    Port     int      // 默认端口
    Protocol string   // 协议
    Probes   [][]byte // 探测数据
    Patterns []string // 匹配模式
}

// DefaultProbes 默认探测规则
var DefaultProbes = map[string]*ServiceProbe{
    "http": {
        Name:     "HTTP",
        Port:     80,
        Protocol: "tcp",
        Probes: [][]byte{
            []byte("HEAD / HTTP/1.0\r\n\r\n"),
            []byte("GET / HTTP/1.0\r\n\r\n"),
        },
        Patterns: []string{
            `^HTTP/[\d.]+\s+\d+`,
            `Server:\s+([^\r\n]+)`,
        },
    },
    "https": {
        Name:     "HTTPS",
        Port:     443,
        Protocol: "tcp",
        Probes: [][]byte{
            []byte{0x16, 0x03, 0x01}, // SSL/TLS ClientHello
        },
        Patterns: []string{
            `^HTTP/[\d.]+\s+\d+`,
            `Server:\s+([^\r\n]+)`,
        },
    },
    "ssh": {
        Name:     "SSH",
        Port:     22,
        Protocol: "tcp",
        Probes:   [][]byte{},
        Patterns: []string{
            `^SSH-[\d.]+`,
            `OpenSSH_[\d.]+`,
        },
    },
    // 添加更多服务探测规则
} 