package service

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"math/rand"
	"net"
	"time"
)

// SYN扫描
func (s *PortScanService) scanPortSYN(ip string, port int, timeout int) (state string, banner string) {
	// 创建原始套接字
	handle, err := pcap.OpenLive("eth0", 65535, true, pcap.BlockForever)
	if err != nil {
		return "error", ""
	}
	defer handle.Close()

	// 构建TCP SYN包
	srcPort := uint16(rand.Intn(65535))
	eth := layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}

	ip4 := layers.IPv4{
		SrcIP:    net.ParseIP("0.0.0.0"),
		DstIP:    net.ParseIP(ip),
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}

	tcp := layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(port),
		SYN:     true,
	}

	// 设置TCP校验和
	tcp.SetNetworkLayerForChecksum(&ip4)

	// 序列化数据包
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := gopacket.SerializeLayers(buffer, opts, &eth, &ip4, &tcp); err != nil {
		return "error", ""
	}

	// 发送数据包
	if err := handle.WritePacketData(buffer.Bytes()); err != nil {
		return "error", ""
	}

	// 设置过滤器只接收目标端口的响应
	filter := fmt.Sprintf("tcp and src host %s and src port %d", ip, port)
	if err := handle.SetBPFFilter(filter); err != nil {
		return "error", ""
	}

	// 设置超时
	start := time.Now()
	for {
		if time.Since(start) > time.Duration(timeout)*time.Second {
			return "filtered", ""
		}

		// 接收响应包
		data, _, err := handle.ReadPacketData()
		if err != nil {
			continue
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, _ := tcpLayer.(*layers.TCP)
		if tcp.SYN && tcp.ACK {
			return "open", ""
		} else if tcp.RST {
			return "closed", ""
		}
	}
}

// UDP扫描
func (s *PortScanService) scanPortUDP(ip string, port int, timeout int) (state string, banner string) {
	// 创建UDP连接
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", ip, port), time.Duration(timeout)*time.Second)
	if err != nil {
		return "error", ""
	}
	defer conn.Close()

	// 发送探测数据
	probeData := s.getUDPProbe(port)
	if _, err := conn.Write(probeData); err != nil {
		return "error", ""
	}

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	// 读取响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return "open|filtered", ""
		}
		return "filtered", ""
	}

	return "open", string(buffer[:n])
}

// 获取UDP探测数据
func (s *PortScanService) getUDPProbe(port int) []byte {
	// 常见UDP服务的探测数据
	probes := map[int][]byte{
		53:  []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},       // DNS查询
		161: []byte{0x30, 0x26, 0x02, 0x01, 0x01, 0x04, 0x06, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63}, // SNMP查询
		137: []byte{0x80, 0x94, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},       // NetBIOS查询
		// 添加更多UDP服务的探测数据
	}

	if probe, ok := probes[port]; ok {
		return probe
	}
	return []byte("Hello\n") // 默认探测数据
}
