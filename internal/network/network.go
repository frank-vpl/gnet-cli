package network

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

type InterfaceInfo struct {
	Name         string   `json:"name"`
	MTU          int      `json:"mtu"`
	Flags        string   `json:"flags"`
	IPs          []string `json:"ips"`
	HardwareAddr string   `json:"hardware_addr,omitempty"`
}

func ListInterfaces() ([]InterfaceInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list interfaces: %w", err)
	}
	var out []InterfaceInfo
	for _, ifc := range ifaces {
		info := InterfaceInfo{
			Name:  ifc.Name,
			MTU:   ifc.MTU,
			Flags: ifc.Flags.String(),
		}
		if ifc.HardwareAddr != nil {
			info.HardwareAddr = ifc.HardwareAddr.String()
		}
		addrs, err := ifc.Addrs()
		if err == nil {
			for _, a := range addrs {
				ipStr := a.String()
				if ip, _, err := net.ParseCIDR(ipStr); err == nil {
					ipStr = ip.String()
				}
				info.IPs = append(info.IPs, ipStr)
			}
		}
		out = append(out, info)
	}
	return out, nil
}

func PrimaryIPv4() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("detect primary IP: %w", err)
	}
	defer conn.Close()
	addr := conn.LocalAddr().(*net.UDPAddr)
	return addr.IP.String(), nil
}

func OS() string {
	return runtime.GOOS
}

func IsIP(s string) bool {
	return net.ParseIP(s) != nil
}

func FormatInterfaces(ifaces []InterfaceInfo) string {
	var b strings.Builder
	for _, i := range ifaces {
		fmt.Fprintf(&b, "%-20s %s\n", i.Name, i.Flags)
		if i.HardwareAddr != "" {
			fmt.Fprintf(&b, "  HW:   %s\n", i.HardwareAddr)
		}
		for _, ip := range i.IPs {
			fmt.Fprintf(&b, "  IP:   %s\n", ip)
		}
	}
	return b.String()
}
