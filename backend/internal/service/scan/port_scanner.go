package scan

import (
	"fmt"
	"net"
	"time"
)

// PortScanner 端口扫描器
type PortScanner struct {
	config  *ScanConfig
	dialer  *net.Dialer
	timeout time.Duration
}

// NewPortScanner 创建端口扫描器
func NewPortScanner(config *ScanConfig) Scanner {
	if config == nil {
		config = DefaultConfig()
	}
	
	return &PortScanner{
		config: config,
		dialer: &net.Dialer{
			Timeout: config.Timeout,
		},
		timeout: config.Timeout,
	}
}

// Scan 执行端口扫描
func (s *PortScanner) Scan(target string, port int, protocol string) (*ScanResult, error) {
	result := &ScanResult{
		Target:    target,
		Port:      port,
		Protocol:  protocol,
		Status:    StatusClosed,
		Timestamp: time.Now(),
	}

	// 构建地址
	addr := fmt.Sprintf("%s:%d", target, port)

	// TCP扫描
	if protocol == "tcp" {
		conn, err := s.dialer.Dial("tcp", addr)
		if err != nil {
			result.Error = err
			return result, nil
		}
		defer conn.Close()

		result.Status = StatusOpen

		// 获取banner
		if s.config.BannerGrab {
			banner, err := s.grabBanner(conn)
			if err == nil {
				result.Banner = banner
			}
		}
	}

	// UDP扫描
	if protocol == "udp" {
		conn, err := net.ListenPacket("udp", "")
		if err != nil {
			result.Error = err
			return result, nil
		}
		defer conn.Close()

		// 发送UDP探测包
		if err := s.sendUDPProbe(conn, addr); err != nil {
			result.Error = err
			return result, nil
		}

		// 等待响应
		if err := s.receiveUDPResponse(conn); err == nil {
			result.Status = StatusOpen
		}
	}

	return result, nil
}

// Stop 停止扫描
func (s *PortScanner) Stop() {
	// 实现停止逻辑
}

// grabBanner 获取服务banner
func (s *PortScanner) grabBanner(conn net.Conn) (string, error) {
	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(s.timeout))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

// sendUDPProbe 发送UDP探测包
func (s *PortScanner) sendUDPProbe(conn net.PacketConn, addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	// 发送空包
	_, err = conn.WriteTo([]byte{}, udpAddr)
	return err
}

// receiveUDPResponse 接收UDP响应
func (s *PortScanner) receiveUDPResponse(conn net.PacketConn) error {
	conn.SetReadDeadline(time.Now().Add(s.timeout))
	
	buf := make([]byte, 4096)
	_, _, err := conn.ReadFrom(buf)
	return err
} 