package service

import (
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"gorm.io/gorm"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// 主机发现服务
type HostDiscoveryService struct {
	db *gorm.DB
}

// 主机发现结果
type HostDiscoveryResult struct {
	IP        string
	IsAlive   bool
	TTL       int
	OS        string
	OpenPorts []int
	Hostnames []string
	MAC       string
	Vendor    string
	LastSeen  time.Time
}

// 主机发现方法
func (s *HostDiscoveryService) DiscoverHosts(cidr string, methods []string) ([]*HostDiscoveryResult, error) {
	var results []*HostDiscoveryResult
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// 解析CIDR
	ipNet, err := s.parseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	// 获取所有IP
	ips := s.getAllIPs(ipNet)

	// 创建工作池
	jobs := make(chan string, len(ips))
	for _, ip := range ips {
		jobs <- ip
	}
	close(jobs)

	// 启动工作协程
	workers := 100
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range jobs {
				result := &HostDiscoveryResult{IP: ip}

				// 执行多种主机发现方法
				for _, method := range methods {
					switch method {
					case "icmp":
						if alive := s.icmpPing(ip); alive {
							result.IsAlive = true
						}
					case "tcp":
						if alive := s.tcpPing(ip); alive {
							result.IsAlive = true
						}
					case "arp":
						if mac := s.arpScan(ip); mac != "" {
							result.IsAlive = true
							result.MAC = mac
							result.Vendor = s.lookupVendor(mac)
						}
					case "udp":
						if alive := s.udpPing(ip); alive {
							result.IsAlive = true
						}
					}

					if result.IsAlive {
						break
					}
				}

				// 如果主机存活，获取更多信息
				if result.IsAlive {
					result.TTL = s.getTTL(ip)
					result.OS = s.guessOS(result.TTL)
					result.Hostnames = s.reverseDNS(ip)
					result.LastSeen = time.Now()

					mutex.Lock()
					results = append(results, result)
					mutex.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return results, nil
}

// ICMP Ping
func (s *HostDiscoveryService) icmpPing(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	return cmd.Run() == nil
}

// TCP Ping (SYN to port 80)
func (s *HostDiscoveryService) tcpPing(ip string) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", ip), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// ARP扫描
func (s *HostDiscoveryService) arpScan(ip string) string {
	// 获取本地接口
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		// 跳过loopback和down的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 发送ARP请求
		handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
		if err != nil {
			continue
		}
		defer handle.Close()

		// 构建ARP请求包
		eth := layers.Ethernet{
			SrcMAC:       iface.HardwareAddr,
			DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			EthernetType: layers.EthernetTypeARP,
		}

		arp := layers.ARP{
			AddrType:          layers.LinkTypeEthernet,
			Protocol:          layers.EthernetTypeIPv4,
			HwAddressSize:     6,
			ProtAddressSize:   4,
			Operation:         layers.ARPRequest,
			SourceHwAddress:   []byte(iface.HardwareAddr),
			SourceProtAddress: []byte(net.ParseIP("0.0.0.0").To4()),
			DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
			DstProtAddress:    []byte(net.ParseIP(ip).To4()),
		}

		// 序列化和发送
		buffer := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}

		gopacket.SerializeLayers(buffer, opts, &eth, &arp)
		if err := handle.WritePacketData(buffer.Bytes()); err != nil {
			continue
		}

		// 等待ARP响应
		start := time.Now()
		for time.Since(start) < time.Second {
			data, _, err := handle.ReadPacketData()
			if err != nil {
				continue
			}

			packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}

			arp := arpLayer.(*layers.ARP)
			if arp.Operation != layers.ARPReply || !bytes.Equal(arp.SourceProtAddress, net.ParseIP(ip).To4()) {
				continue
			}

			return net.HardwareAddr(arp.SourceHwAddress).String()
		}
	}

	return ""
}

// UDP Ping
func (s *HostDiscoveryService) udpPing(ip string) bool {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:53", ip), time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	// 发送DNS查询包
	query := []byte{
		0x00, 0x00, // Transaction ID
		0x01, 0x00, // Flags
		0x00, 0x01, // Questions
		0x00, 0x00, // Answer RRs
		0x00, 0x00, // Authority RRs
		0x00, 0x00, // Additional RRs
	}

	if _, err := conn.Write(query); err != nil {
		return false
	}

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(time.Second))

	// 读取响应
	buffer := make([]byte, 512)
	_, err = conn.Read(buffer)
	return err == nil
}

// 获取TTL值
func (s *HostDiscoveryService) getTTL(ip string) int {
	cmd := exec.Command("ping", "-c", "1", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0
	}

	// 解析TTL值
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "ttl=") {
			ttlStr := strings.Split(strings.Split(line, "ttl=")[1], " ")[0]
			ttl := 0
			fmt.Sscanf(ttlStr, "%d", &ttl)
			return ttl
		}
	}
	return 0
}

// 根据TTL猜测操作系统
func (s *HostDiscoveryService) guessOS(ttl int) string {
	switch {
	case ttl <= 64:
		return "Linux/Unix"
	case ttl <= 128:
		return "Windows"
	case ttl <= 255:
		return "Cisco/Network"
	default:
		return "Unknown"
	}
}

// 反向DNS查询
func (s *HostDiscoveryService) reverseDNS(ip string) []string {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil
	}
	return names
}

// MAC厂商查询
func (s *HostDiscoveryService) lookupVendor(mac string) string {
	// 可以实现本地MAC厂商数据库查询
	// 或者调用在线API
	return "Unknown"
}

// 解析CIDR
func (s *HostDiscoveryService) parseCIDR(cidr string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	return ipNet, err
}

// 获取CIDR范围内所有IP
func (s *HostDiscoveryService) getAllIPs(ipNet *net.IPNet) []string {
	var ips []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); s.nextIP(ip) {
		ips = append(ips, ip.String())
	}
	return ips
}

// 获取下一个IP
func (s *HostDiscoveryService) nextIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
