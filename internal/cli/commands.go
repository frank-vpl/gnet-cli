package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/frank-vpl/gnet-cli/internal/config"
	"github.com/frank-vpl/gnet-cli/internal/logger"
	"github.com/frank-vpl/gnet-cli/internal/network"
	"github.com/frank-vpl/gnet-cli/internal/proxy"
	"github.com/frank-vpl/gnet-cli/internal/share"
)

var (
	flagHTTP    int
	flagSocks   int
	flagBind    string
	flagHTTPOn  bool
	flagSocksOn bool
	flagLogFile string
	flagDebug   bool
)

func Execute() error {
	root := &cobra.Command{
		Use:   "gnet",
		Short: "GNet CLI — share your VPN over HTTP/SOCKS5 to any device",
		Long:  "GNet CLI turns your machine into a proxy server so other devices on your network (LAN, Wi-Fi, hotspot, USB, router) can route traffic through your VPN.\n\nRun without arguments to launch the interactive setup wizard.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			return RunWizard(cfg)
		},
	}

	serve := &cobra.Command{
		Use:   "serve",
		Short: "Start the proxy server(s) with saved or flag-provided config",
		RunE:  runServe,
	}
	serve.Flags().IntVar(&flagHTTP, "http-port", 0, "HTTP proxy port")
	serve.Flags().IntVar(&flagSocks, "socks-port", 0, "SOCKS5 proxy port")
	serve.Flags().StringVar(&flagBind, "bind", "", "Bind address (default 0.0.0.0)")
	serve.Flags().BoolVar(&flagHTTPOn, "http", true, "Enable HTTP proxy")
	serve.Flags().BoolVar(&flagSocksOn, "socks5", true, "Enable SOCKS5 proxy")
	serve.Flags().StringVar(&flagLogFile, "log-file", "", "Write logs to file")
	serve.Flags().BoolVar(&flagDebug, "debug", false, "Enable debug logging")

	cfgCmd := &cobra.Command{
		Use:   "config",
		Short: "Show or edit the saved configuration",
		RunE:  runConfig,
	}

	ifaces := &cobra.Command{
		Use:   "interfaces",
		Short: "List network interfaces and IP addresses",
		RunE:  runInterfaces,
	}

	methods := &cobra.Command{
		Use:   "methods",
		Short: "List available sharing methods and setup steps",
		RunE:  runMethods,
	}

	root.AddCommand(serve, cfgCmd, ifaces, methods)
	return root.Execute()
}

func runServe(cmd *cobra.Command, args []string) error {
	if flagDebug {
		logger.SetLevel(logger.LevelDebug)
	}
	if flagLogFile != "" {
		if err := logger.EnableFileLog(flagLogFile); err != nil {
			return fmt.Errorf("open log file: %w", err)
		}
		defer logger.Close()
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if flagHTTP != 0 {
		cfg.HTTPPort = flagHTTP
	}
	if flagSocks != 0 {
		cfg.Socks5Port = flagSocks
	}
	if flagBind != "" {
		cfg.BindAddr = flagBind
	}
	cfg.HTTPEnabled = flagHTTPOn
	cfg.Socks5Enabled = flagSocksOn
	if err := cfg.Validate(); err != nil {
		return err
	}

	if err := config.Save(cfg); err != nil {
		logger.Warn("could not save config", "err", err)
	}

	ports := activePorts(&cfg)
	fw := network.EnsureFirewall(ports)
	if fw.Message != "" {
		fmt.Printf("  Firewall: %s\n", fw.Message)
		if fw.NeedsAdmin {
			fmt.Println()
		}
	}

	srv := proxy.New(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		return err
	}

	fmt.Printf("\n  GNet CLI proxy running\n")
	if cfg.HTTPEnabled {
		fmt.Printf("  HTTP   -> %s:%d\n", cfg.BindAddr, cfg.HTTPPort)
	}
	if cfg.Socks5Enabled {
		fmt.Printf("  SOCKS5 -> %s:%d\n", cfg.BindAddr, cfg.Socks5Port)
	}
	if primary, err := network.PrimaryIPv4(); err == nil {
		fmt.Printf("  LAN IP -> %s\n", primary)
	}
	fmt.Println("  Press Ctrl+C to stop")
	fmt.Println()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\n  Stopping...")
	srv.Stop()
	wait := make(chan struct{})
	go func() {
		srv.Wait()
		close(wait)
	}()
	select {
	case <-wait:
	case <-time.After(5 * time.Second):
	}
	fmt.Println("  Stopped.")
	return nil
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	path, _ := config.ConfigPath()
	fmt.Printf("Config path: %s\n\n%s\n", path, string(data))
	return nil
}

func runInterfaces(cmd *cobra.Command, args []string) error {
	ifaces, err := network.ListInterfaces()
	if err != nil {
		return err
	}
	fmt.Printf("OS: %s\n\n", network.OS())
	fmt.Print(network.FormatInterfaces(ifaces))
	if primary, err := network.PrimaryIPv4(); err == nil {
		fmt.Printf("\nPrimary outbound IPv4: %s\n", primary)
	}
	return nil
}

func runMethods(cmd *cobra.Command, args []string) error {
	fmt.Print(share.FormatMethods())
	return nil
}
