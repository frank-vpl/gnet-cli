package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultHTTPPort   = 8080
	DefaultSocks5Port = 1080
	DefaultBindAddr   = "0.0.0.0"
	DefaultConfigDir  = ".gnet-cli"
	configFileName    = "config.json"
)

type ProxyType string

const (
	ProxyHTTP   ProxyType = "http"
	ProxySOCKS5 ProxyType = "socks5"
)

type ProxyConfig struct {
	HTTPPort      int    `json:"http_port"`
	Socks5Port    int    `json:"socks5_port"`
	BindAddr      string `json:"bind_addr"`
	HTTPEnabled   bool   `json:"http_enabled"`
	Socks5Enabled bool   `json:"socks5_enabled"`
	AuthUsername  string `json:"auth_username,omitempty"`
	AuthPassword  string `json:"auth_password,omitempty"`
	SharingMethod string `json:"sharing_method,omitempty"`
}

func Default() ProxyConfig {
	return ProxyConfig{
		HTTPPort:      DefaultHTTPPort,
		Socks5Port:    DefaultSocks5Port,
		BindAddr:      DefaultBindAddr,
		HTTPEnabled:   true,
		Socks5Enabled: true,
	}
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	dir := filepath.Join(home, DefaultConfigDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create config dir: %w", err)
	}
	return filepath.Join(dir, configFileName), nil
}

func Load() (ProxyConfig, error) {
	path, err := ConfigPath()
	if err != nil {
		return Default(), err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg := Default()
			_ = Save(cfg)
			return cfg, nil
		}
		return Default(), fmt.Errorf("read config: %w", err)
	}
	var cfg ProxyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), fmt.Errorf("parse config: %w", err)
	}
	cfg.applyDefaults()
	return cfg, nil
}

func Save(cfg ProxyConfig) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	cfg.applyDefaults()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func (c *ProxyConfig) applyDefaults() {
	if c.HTTPPort == 0 {
		c.HTTPPort = DefaultHTTPPort
	}
	if c.Socks5Port == 0 {
		c.Socks5Port = DefaultSocks5Port
	}
	if c.BindAddr == "" {
		c.BindAddr = DefaultBindAddr
	}
}

func (c ProxyConfig) Validate() error {
	if !c.HTTPEnabled && !c.Socks5Enabled {
		return errors.New("at least one proxy (HTTP or SOCKS5) must be enabled")
	}
	if c.HTTPPort < 1 || c.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.HTTPPort)
	}
	if c.Socks5Port < 1 || c.Socks5Port > 65535 {
		return fmt.Errorf("invalid SOCKS5 port: %d", c.Socks5Port)
	}
	if c.HTTPEnabled && c.Socks5Enabled && c.HTTPPort == c.Socks5Port {
		return errors.New("HTTP and SOCKS5 ports must be different")
	}
	return nil
}
