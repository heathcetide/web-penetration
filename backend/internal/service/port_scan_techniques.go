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

// TCP 标志常量
const (
	TCPFlagFIN = 0x01
	TCPFlagSYN = 0x02
	TCPFlagRST = 0x04
	TCPFlagPSH = 0x08
	TCPFlagACK = 0x10
	TCPFlagURG = 0x20
)

// ACK扫描
func (s *PortScanService) scanPortACK(ip string, port int, timeout int) (state string) {
	handle, err := pcap.OpenLive("eth0", 65535, true, pcap.BlockForever)
	if err != nil {
		return "error"
	}
	defer handle.Close()

	// 构建TCP ACK包
	srcPort := uint16(rand.Intn(65535))
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(port),
		ACK:     true,
		Window:  1024,
		Seq:     rand.Uint32(),
	}

	// 发送数据包并分析响应
	if err := s.sendTCPPacket(handle, ip, tcp); err != nil {
		return "error"
	}

	// 如果收到RST，说明端口未被过滤
	if s.receiveTCPResponse(handle, ip, port, timeout, TCPFlagRST) {
		return "unfiltered"
	}
	return "filtered"
}

// FIN扫描
func (s *PortScanService) scanPortFIN(ip string, port int, timeout int) (state string) {
	handle, err := pcap.OpenLive("eth0", 65535, true, pcap.BlockForever)
	if err != nil {
		return "error"
	}
	defer handle.Close()

	// 构建TCP FIN包
	srcPort := uint16(rand.Intn(65535))
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(port),
		FIN:     true,
		Window:  1024,
		Seq:     rand.Uint32(),
	}

	// 发送数据包并分析响应
	if err := s.sendTCPPacket(handle, ip, tcp); err != nil {
		return "error"
	}

	// 如果收到RST，说明端口关闭
	if s.receiveTCPResponse(handle, ip, port, timeout, TCPFlagRST) {
		return "closed"
	}
	return "open|filtered"
}

// NULL扫描
func (s *PortScanService) scanPortNULL(ip string, port int, timeout int) (state string) {
	handle, err := pcap.OpenLive("eth0", 65535, true, pcap.BlockForever)
	if err != nil {
		return "error"
	}
	defer handle.Close()

	// 构建TCP NULL包（无标志位）
	srcPort := uint16(rand.Intn(65535))
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(port),
		Window:  1024,
		Seq:     rand.Uint32(),
	}

	// 发送数据包并分析响应
	if err := s.sendTCPPacket(handle, ip, tcp); err != nil {
		return "error"
	}

	// 如果收到RST，说明端口关闭
	if s.receiveTCPResponse(handle, ip, port, timeout, TCPFlagRST) {
		return "closed"
	}
	return "open|filtered"
}

// XMAS扫描
func (s *PortScanService) scanPortXMAS(ip string, port int, timeout int) (state string) {
	handle, err := pcap.OpenLive("eth0", 65535, true, pcap.BlockForever)
	if err != nil {
		return "error"
	}
	defer handle.Close()

	// 构建TCP XMAS包（FIN、PSH、URG标志位都位）
	srcPort := uint16(rand.Intn(65535))
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(port),
		FIN:     true,
		PSH:     true,
		URG:     true,
		Window:  1024,
		Seq:     rand.Uint32(),
	}

	// 发送数据包并分析响应
	if err := s.sendTCPPacket(handle, ip, tcp); err != nil {
		return "error"
	}

	// 如果收到RST，说明端口关闭
	if s.receiveTCPResponse(handle, ip, port, timeout, TCPFlagRST) {
		return "closed"
	}
	return "open|filtered"
}

// 发送TCP数据包
func (s *PortScanService) sendTCPPacket(handle *pcap.Handle, ip string, tcp *layers.TCP) error {
	// 构建以太网和IP层
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

	// 设置TCP校验和
	tcp.SetNetworkLayerForChecksum(&ip4)

	// 序列化数据包
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := gopacket.SerializeLayers(buffer, opts, &eth, &ip4, tcp); err != nil {
		return err
	}

	// 发送数据包
	return handle.WritePacketData(buffer.Bytes())
}

// 接收TCP响应
func (s *PortScanService) receiveTCPResponse(handle *pcap.Handle, ip string, port int, timeout int, flag uint8) bool {
	// 设置过滤器
	filter := fmt.Sprintf("tcp and src host %s and src port %d", ip, port)
	if err := handle.SetBPFFilter(filter); err != nil {
		return false
	}

	// 设置超时
	start := time.Now()
	for {
		if time.Since(start) > time.Duration(timeout)*time.Second {
			return false
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
		if (tcp.FIN && flag == TCPFlagFIN) ||
			(tcp.SYN && flag == TCPFlagSYN) ||
			(tcp.RST && flag == TCPFlagRST) ||
			(tcp.PSH && flag == TCPFlagPSH) ||
			(tcp.ACK && flag == TCPFlagACK) ||
			(tcp.URG && flag == TCPFlagURG) {
			return true
		}
	}
}
