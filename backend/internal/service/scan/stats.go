package scan

// 统计相关的类型定义
type ServiceCount struct {
    Service string `json:"service"`
    Count   int    `json:"count"`
}

type PortCount struct {
    Port  int `json:"port"`
    Count int `json:"count"`
}

type ProtocolCount struct {
    Protocol string `json:"protocol"`
    Count    int    `json:"count"`
}

// 扫描统计
type ScanStats struct {
    TotalScans      int                    `json:"total_scans"`
    OpenPorts       int                    `json:"open_ports"`
    ClosedPorts     int                    `json:"closed_ports"`
    FilteredPorts   int                    `json:"filtered_ports"`
    UniqueServices  int                    `json:"unique_services"`
    UniqueProtocols int                    `json:"unique_protocols"`
    ServiceVersions map[string][]string    `json:"service_versions"`
    PortsByService  map[string][]int       `json:"ports_by_service"`
    VulnsByService  map[string][]string    `json:"vulns_by_service"`
    TopPorts        []PortCount            `json:"top_ports"`
    TopServices     []ServiceCount         `json:"top_services"`
    Protocols       []ProtocolCount        `json:"protocols"`
} 