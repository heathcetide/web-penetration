package service

import (
	_ "bytes"
	_ "encoding/hex"
	"net"
	"time"
)

// 服务探测配置
type ProbeConfig struct {
	Name       string        // 探测名称
	Probe      []byte        // 探测数据
	MatchRegex string        // 匹配正则
	Ports      []int         // 适用端口
	Protocol   string        // 协议
	Timeout    time.Duration // 超时时间
	SSLPorts   []int         // SSL端口
}

// 常见服务探测配置
var commonProbes = []ProbeConfig{
	{
		Name:     "HTTP",
		Probe:    []byte("HEAD / HTTP/1.0\r\n\r\n"),
		Protocol: "tcp",
		Ports:    []int{80, 8080, 8000, 8008},
		SSLPorts: []int{443, 8443},
	},
	{
		Name:     "SSH",
		Protocol: "tcp",
		Ports:    []int{22},
	},
	{
		Name:     "FTP",
		Protocol: "tcp",
		Ports:    []int{21},
	},
	{
		Name:     "SMTP",
		Probe:    []byte("EHLO test\r\n"),
		Protocol: "tcp",
		Ports:    []int{25, 465, 587},
	},
	{
		Name:     "POP3",
		Protocol: "tcp",
		Ports:    []int{110, 995},
	},
	{
		Name:     "IMAP",
		Protocol: "tcp",
		Ports:    []int{143, 993},
	},
	{
		Name:     "MySQL",
		Protocol: "tcp",
		Ports:    []int{3306},
	},
	{
		Name:     "MSSQL",
		Protocol: "tcp",
		Ports:    []int{1433},
	},
	{
		Name:     "Redis",
		Probe:    []byte("*1\r\n$4\r\ninfo\r\n"),
		Protocol: "tcp",
		Ports:    []int{6379},
	},
	{
		Name:     "MongoDB",
		Protocol: "tcp",
		Ports:    []int{27017},
	},
}

// 探测服务
func (s *PortScanService) probeService(ip string, port int, protocol string) (service string, banner string) {
	// 查找适用的探测配置
	var probes []ProbeConfig
	for _, probe := range commonProbes {
		if probe.Protocol == protocol {
			for _, p := range probe.Ports {
				if p == port {
					probes = append(probes, probe)
					break
				}
			}
		}
	}

	// 如果没有特定配置，使用通用探测
	if len(probes) == 0 {
		probes = []ProbeConfig{
			{
				Name:     "Generic",
				Probe:    []byte("\r\n"),
				Protocol: protocol,
				Timeout:  time.Second * 5,
			},
		}
	}

	// 执行探测
	for _, probe := range probes {
		service, banner = s.executeProbe(ip, port, probe)
		if service != "" {
			return service, banner
		}
	}

	return "unknown", ""
}

// 执行探测
func (s *PortScanService) executeProbe(ip string, port int, probe ProbeConfig) (service string, banner string) {
	// 建立连接
	conn, err := net.DialTimeout(probe.Protocol, net.JoinHostPort(ip, string(port)), probe.Timeout)
	if err != nil {
		return "", ""
	}
	defer conn.Close()

	// 发送探测数据
	if len(probe.Probe) > 0 {
		conn.SetWriteDeadline(time.Now().Add(probe.Timeout))
		if _, err := conn.Write(probe.Probe); err != nil {
			return "", ""
		}
	}

	// 读取响应
	conn.SetReadDeadline(time.Now().Add(probe.Timeout))
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", ""
	}

	banner = string(buffer[:n])
	return probe.Name, banner
}
