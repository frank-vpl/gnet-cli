package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/frank-vpl/gnet-cli/internal/logger"
)

func handleHTTP(client net.Conn) {
	defer client.Close()
	clientIP := client.RemoteAddr().String()
	logger.Info("HTTP client connected", "client", clientIP)

	br := bufio.NewReader(client)
	reqLine, err := br.ReadString('\n')
	if err != nil {
		logger.Warn("failed to read HTTP request line", "client", clientIP, "err", err)
		return
	}
	reqLine = strings.TrimSpace(reqLine)
	if reqLine == "" {
		return
	}
	logger.Debug("HTTP request", "client", clientIP, "line", reqLine)

	parts := strings.Split(reqLine, " ")
	if len(parts) < 3 {
		sendHTTPError(client, 400, "Bad Request")
		return
	}
	method := parts[0]
	target := parts[1]

	if strings.EqualFold(method, "CONNECT") {
		handleHTTPSConnect(client, target)
		return
	}

	headers, err := readHeaders(br)
	if err != nil {
		sendHTTPError(client, 400, "Bad Request")
		return
	}

	var hostHeader string
	for _, h := range headers {
		if strings.HasPrefix(strings.ToLower(h), "host:") {
			hostHeader = strings.TrimSpace(h[5:])
		}
	}

	requestURL := target
	if !strings.HasPrefix(strings.ToLower(target), "http://") && !strings.HasPrefix(strings.ToLower(target), "https://") {
		if hostHeader == "" {
			sendHTTPError(client, 400, "Bad Request")
			return
		}
		requestURL = "http://" + hostHeader + target
	}

	forwardHTTPRequest(client, method, requestURL, headers)
}

func readHeaders(br *bufio.Reader) ([]string, error) {
	var headers []string
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			return nil, err
		}
		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed == "" {
			break
		}
		headers = append(headers, trimmed)
	}
	return headers, nil
}

func forwardHTTPRequest(client net.Conn, method, rawURL string, headers []string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		logger.Warn("invalid URL", "url", rawURL, "err", err)
		sendHTTPError(client, 400, "Bad Request")
		return
	}
	hostname := u.Hostname()
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	addr := net.JoinHostPort(hostname, port)

	target, err := net.DialTimeout("tcp", addr, 15*time.Second)
	if err != nil {
		logger.Warn("failed to connect target", "addr", addr, "err", err)
		sendHTTPError(client, 502, "Bad Gateway")
		return
	}
	defer target.Close()

	var req bytes.Buffer
	path := u.RequestURI()
	if path == "" {
		path = "/"
	}
	req.WriteString(fmt.Sprintf("%s %s HTTP/1.1\r\n", method, path))

	hostSent := false
	for _, h := range headers {
		lower := strings.ToLower(h)
		if strings.HasPrefix(lower, "proxy-") || strings.HasPrefix(lower, "connection:") {
			continue
		}
		req.WriteString(h)
		req.WriteString("\r\n")
		if strings.HasPrefix(lower, "host:") {
			hostSent = true
		}
	}
	if !hostSent {
		req.WriteString(fmt.Sprintf("Host: %s\r\n", hostname))
	}
	req.WriteString("Connection: close\r\n\r\n")

	if _, err := target.Write(req.Bytes()); err != nil {
		logger.Warn("failed to write to target", "err", err)
		sendHTTPError(client, 502, "Bad Gateway")
		return
	}

	io.Copy(client, target)
}

func handleHTTPSConnect(client net.Conn, target string) {
	host, portStr, err := net.SplitHostPort(target)
	if err != nil {
		host = target
		portStr = "443"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 443
	}
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	targetConn, err := net.DialTimeout("tcp", addr, 15*time.Second)
	if err != nil {
		logger.Warn("CONNECT failed", "addr", addr, "err", err)
		sendHTTPError(client, 502, "Bad Gateway")
		return
	}
	defer targetConn.Close()

	if _, err := client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		logger.Warn("failed to send CONNECT response", "err", err)
		return
	}
	logger.Info("HTTPS tunnel established", "client", client.RemoteAddr(), "target", addr)
	bidirectionalRelay(client, targetConn)
}

func sendHTTPError(client net.Conn, code int, msg string) {
	resp := fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nConnection: close\r\nContent-Length: %d\r\n\r\n%d %s\n",
		code, msg, len(fmt.Sprintf("%d %s\n", code, msg)), code, msg)
	_, _ = client.Write([]byte(resp))
}
