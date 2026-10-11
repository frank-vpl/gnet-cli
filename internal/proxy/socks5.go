package proxy

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/frank-vpl/gnet-cli/internal/logger"
)

const (
	socks5Version = 0x05
	cmdConnect    = 0x01
	atypIPv4      = 0x01
	atypDomain    = 0x03
	atypIPv6      = 0x04
)

func handleSOCKS5(client net.Conn) {
	defer client.Close()
	clientIP := client.RemoteAddr().String()
	logger.Info("SOCKS5 client connected", "client", clientIP)

	if err := socksHandshake(client); err != nil {
		logger.Warn("SOCKS5 handshake failed", "client", clientIP, "err", err)
		return
	}

	target, err := socksReadRequest(client)
	if err != nil {
		logger.Warn("SOCKS5 request failed", "client", clientIP, "err", err)
		writeSocksReply(client, 0x01, "", 0)
		return
	}

	targetConn, err := net.DialTimeout("tcp", target, 15*time.Second)
	if err != nil {
		logger.Warn("SOCKS5 connect target failed", "target", target, "err", err)
		writeSocksReply(client, 0x05, "", 0)
		return
	}
	defer targetConn.Close()

	writeSocksReply(client, 0x00, "0.0.0.0", 0)
	logger.Info("SOCKS5 connection established", "client", clientIP, "target", target)
	bidirectionalRelay(client, targetConn)
}

func socksHandshake(client net.Conn) error {
	header := make([]byte, 2)
	if _, err := io.ReadFull(client, header); err != nil {
		return err
	}
	if header[0] != socks5Version {
		return fmt.Errorf("unsupported SOCKS version: %d", header[0])
	}
	nMethods := int(header[1])
	methods := make([]byte, nMethods)
	if _, err := io.ReadFull(client, methods); err != nil {
		return err
	}
	noAuth := false
	for _, m := range methods {
		if m == 0x00 {
			noAuth = true
			break
		}
	}
	if !noAuth {
		client.Write([]byte{socks5Version, 0xFF})
		return fmt.Errorf("no acceptable auth method")
	}
	_, err := client.Write([]byte{socks5Version, 0x00})
	return err
}

func socksReadRequest(client net.Conn) (string, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(client, header); err != nil {
		return "", err
	}
	if header[0] != socks5Version {
		return "", fmt.Errorf("unsupported SOCKS version: %d", header[0])
	}
	if header[1] != cmdConnect {
		return "", fmt.Errorf("unsupported command: %d", header[1])
	}
	atyp := header[3]

	var host string
	switch atyp {
	case atypIPv4:
		ip := make([]byte, 4)
		if _, err := io.ReadFull(client, ip); err != nil {
			return "", err
		}
		host = net.IP(ip).String()
	case atypDomain:
		lenBuf := make([]byte, 1)
		if _, err := io.ReadFull(client, lenBuf); err != nil {
			return "", err
		}
		domain := make([]byte, lenBuf[0])
		if _, err := io.ReadFull(client, domain); err != nil {
			return "", err
		}
		host = string(domain)
	case atypIPv6:
		ip := make([]byte, 16)
		if _, err := io.ReadFull(client, ip); err != nil {
			return "", err
		}
		host = net.IP(ip).String()
	default:
		return "", fmt.Errorf("unsupported address type: %d", atyp)
	}

	portBuf := make([]byte, 2)
	if _, err := io.ReadFull(client, portBuf); err != nil {
		return "", err
	}
	port := binary.BigEndian.Uint16(portBuf)
	return net.JoinHostPort(host, strconv.Itoa(int(port))), nil
}

func writeSocksReply(client net.Conn, rep byte, bindAddr string, bindPort int) {
	ip := net.ParseIP(bindAddr)
	if ip == nil {
		ip = net.IPv4zero
	}
	ip4 := ip.To4()
	reply := []byte{socks5Version, rep, 0x00, atypIPv4}
	if ip4 != nil {
		reply = append(reply, ip4...)
	} else {
		reply[3] = atypIPv6
		reply = append(reply, ip.To16()...)
	}
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, uint16(bindPort))
	reply = append(reply, portBytes...)
	_, _ = client.Write(reply)
}
