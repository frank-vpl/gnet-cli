package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/frank-vpl/gnet-cli/internal/config"
	"github.com/frank-vpl/gnet-cli/internal/logger"
)

type Status struct {
	HTTPRunning   bool   `json:"http_running"`
	Socks5Running bool   `json:"socks5_running"`
	HTTPPort      int    `json:"http_port"`
	Socks5Port    int    `json:"socks5_port"`
	BindAddr      string `json:"bind_addr"`
	ActiveConns   int64  `json:"active_conns"`
}

type Server struct {
	cfg config.ProxyConfig

	httpLn     net.Listener
	socksLn    net.Listener
	actualBind string
	stopCh     chan struct{}
	wg         sync.WaitGroup
	conns      sync.Map
	connID     atomic.Int64
	running    atomic.Bool
}

func New(cfg config.ProxyConfig) *Server {
	return &Server{cfg: cfg, stopCh: make(chan struct{})}
}

func (s *Server) Start(ctx context.Context) error {
	if !s.cfg.HTTPEnabled && !s.cfg.Socks5Enabled {
		return errors.New("no proxy enabled")
	}
	if s.running.Load() {
		return errors.New("server already running")
	}
	s.stopCh = make(chan struct{})

	if s.cfg.HTTPEnabled {
		ln, bind, err := s.listenWithFallback(s.cfg.HTTPPort)
		if err != nil {
			return fmt.Errorf("listen HTTP on port %d: %w", s.cfg.HTTPPort, err)
		}
		s.actualBind = bind
		s.httpLn = ln
		s.wg.Add(1)
		go s.acceptLoop(ln, handleHTTP, "HTTP")
		logger.Info("HTTP proxy listening", "addr", fmt.Sprintf("%s:%d", bind, s.cfg.HTTPPort))
	}

	if s.cfg.Socks5Enabled {
		ln, bind, err := s.listenWithFallback(s.cfg.Socks5Port)
		if err != nil {
			s.closeListeners()
			return fmt.Errorf("listen SOCKS5 on port %d: %w", s.cfg.Socks5Port, err)
		}
		if s.actualBind == "" {
			s.actualBind = bind
		}
		s.socksLn = ln
		s.wg.Add(1)
		go s.acceptLoop(ln, handleSOCKS5, "SOCKS5")
		logger.Info("SOCKS5 proxy listening", "addr", fmt.Sprintf("%s:%d", bind, s.cfg.Socks5Port))
	}

	s.running.Store(true)
	go func() {
		select {
		case <-ctx.Done():
		case <-s.stopCh:
		}
		s.shutdown()
	}()
	return nil
}

func (s *Server) acceptLoop(ln net.Listener, handler func(net.Conn), name string) {
	defer s.wg.Done()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if s.running.Load() {
				logger.Warn("accept failed", "server", name, "err", err)
			}
			return
		}
		id := s.connID.Add(1)
		s.conns.Store(id, conn)
		s.wg.Add(1)
		go func(c net.Conn, cid int64) {
			defer s.wg.Done()
			defer s.conns.Delete(cid)
			handler(c)
		}(conn, id)
	}
}

func (s *Server) Stop() {
	close(s.stopCh)
}

func (s *Server) shutdown() {
	if !s.running.CompareAndSwap(true, false) {
		return
	}
	logger.Info("shutting down proxy servers")
	s.closeListeners()
	s.conns.Range(func(_, v any) bool {
		if c, ok := v.(net.Conn); ok {
			_ = c.Close()
		}
		return true
	})
}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) closeListeners() {
	if s.httpLn != nil {
		_ = s.httpLn.Close()
	}
	if s.socksLn != nil {
		_ = s.socksLn.Close()
	}
}

func (s *Server) Status() Status {
	bind := s.cfg.BindAddr
	if s.actualBind != "" {
		bind = s.actualBind
	}
	st := Status{
		HTTPPort:   s.cfg.HTTPPort,
		Socks5Port: s.cfg.Socks5Port,
		BindAddr:   bind,
	}
	if s.httpLn != nil {
		st.HTTPRunning = true
	}
	if s.socksLn != nil {
		st.Socks5Running = true
	}
	count := int64(0)
	s.conns.Range(func(_, _ any) bool {
		count++
		return true
	})
	st.ActiveConns = count
	return st
}

func (s *Server) IsRunning() bool {
	return s.running.Load()
}

func (s *Server) listenWithFallback(port int) (net.Listener, string, error) {
	requested := s.cfg.BindAddr
	addr := fmt.Sprintf("%s:%d", requested, port)
	ln, err := net.Listen("tcp", addr)
	if err == nil {
		return ln, requested, nil
	}

	if requested == "0.0.0.0" || requested == "" || requested == "::" {
		return nil, "", err
	}

	logger.Warn("binding to specific address failed, falling back to 0.0.0.0",
		"requested", requested, "port", port, "err", err)
	fallback := "0.0.0.0"
	ln, err = net.Listen("tcp", fmt.Sprintf("%s:%d", fallback, port))
	if err != nil {
		return nil, "", err
	}
	return ln, fallback, nil
}
