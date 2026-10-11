package proxy

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/frank-vpl/gnet-cli/internal/logger"
)

const relayBufferSize = 32 * 1024

func relay(dst net.Conn, src net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := make([]byte, relayBufferSize)
	for {
		_ = src.SetReadDeadline(time.Now().Add(2 * time.Minute))
		n, err := src.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				logger.Debug("relay write failed", "err", werr)
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				logger.Debug("relay read ended", "err", err)
			}
			return
		}
	}
}

func bidirectionalRelay(a, b net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go relay(a, b, &wg)
	go relay(b, a, &wg)
	wg.Wait()
}
