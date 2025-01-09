package service

import (
	_ "fmt"
	"net"
	"os"
	"runtime"
	"time"
)

// 路由跟踪服务
type RouteTraceService struct {
	maxHops    int
	timeoutSec time.Duration
}

// 跟踪结果
type TraceResult struct {
	Hop      int           `json:"hop"`
	IP       string        `json:"ip"`
	RTT      time.Duration `json:"rtt"`
	Hostname string        `json:"hostname"`
}

// 创建路由跟踪服务
func NewRouteTraceService(maxHops int, timeoutSec int) *RouteTraceService {
	return &RouteTraceService{
		maxHops:    maxHops,
		timeoutSec: time.Duration(timeoutSec) * time.Second,
	}
}

// 执行路由跟踪
func (s *RouteTraceService) TraceRoute(target string) ([]*TraceResult, error) {
	// 解析目标地址
	ipAddr, err := net.ResolveIPAddr("ip", target)
	if err != nil {
		return nil, err
	}

	var results []*TraceResult
	for ttl := 1; ttl <= s.maxHops; ttl++ {
		// 创建ICMP连接
		conn, err := net.DialTimeout("ip4:icmp", target, s.timeoutSec)
		if err != nil {
			return nil, err
		}

		// 设置TTL
		if err := s.setTTL(conn.(*net.IPConn), ttl); err != nil {
			conn.Close()
			return nil, err
		}

		// 发送ICMP请求
		start := time.Now()
		if _, err := conn.Write(s.createICMPRequest(ttl)); err != nil {
			conn.Close()
			return nil, err
		}

		// 接收响应
		reply := make([]byte, 1500)
		if err := conn.SetReadDeadline(time.Now().Add(s.timeoutSec)); err != nil {
			conn.Close()
			return nil, err
		}

		_, addr, err := conn.(*net.IPConn).ReadFrom(reply)
		conn.Close()

		if err != nil {
			// 超时或其他错误，继续下一跳
			results = append(results, &TraceResult{
				Hop: ttl,
				IP:  "*",
				RTT: s.timeoutSec,
			})
			continue
		}

		// 解析响应
		rtt := time.Since(start)
		hostname, _ := net.LookupAddr(addr.String())

		results = append(results, &TraceResult{
			Hop:      ttl,
			IP:       addr.String(),
			RTT:      rtt,
			Hostname: s.getHostname(hostname),
		})

		// 如果到达目标地址，结束跟踪
		if addr.String() == ipAddr.String() {
			break
		}
	}

	return results, nil
}

// 设置TTL
func (s *RouteTraceService) setTTL(conn *net.IPConn, ttl int) error {
	f, err := conn.File()
	if err != nil {
		return err
	}
	defer f.Close()

	// 根据操作系统设置TTL
	if runtime.GOOS == "windows" {
		return s.setTTLWindows(f, ttl)
	}
	return s.setTTLUnix(f, ttl)
}

// Windows平台设置TTL
func (s *RouteTraceService) setTTLWindows(f *os.File, ttl int) error {
	// Windows下使用WSA函数
	return nil // TODO: 实现Windows下的TTL设置
}

// Unix平台设置TTL
func (s *RouteTraceService) setTTLUnix(f *os.File, ttl int) error {
	// Unix下使用setsockopt
	return nil // TODO: 实现Unix下的TTL设置
}

// 创建ICMP请求
func (s *RouteTraceService) createICMPRequest(seq int) []byte {
	msg := make([]byte, 8)
	msg[0] = 8                // Echo Request
	msg[1] = 0                // Code
	msg[2] = 0                // Checksum
	msg[3] = 0                // Checksum
	msg[4] = 0                // Identifier
	msg[5] = 0                // Identifier
	msg[6] = byte(seq >> 8)   // Sequence Number
	msg[7] = byte(seq & 0xff) // Sequence Number

	// 计算校验和
	s.calculateChecksum(msg)
	return msg
}

// 计算ICMP校验和
func (s *RouteTraceService) calculateChecksum(msg []byte) {
	var sum uint32
	for i := 0; i < len(msg)-1; i += 2 {
		sum += uint32(msg[i+1])<<8 | uint32(msg[i])
	}
	if len(msg)%2 == 1 {
		sum += uint32(msg[len(msg)-1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)

	// 填充校验和
	msg[2] = byte(^uint16(sum))
	msg[3] = byte(^uint16(sum) >> 8)
}

// 获取主机名
func (s *RouteTraceService) getHostname(names []string) string {
	if len(names) > 0 {
		return names[0]
	}
	return ""
}
