package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/frank-vpl/gnet-cli/internal/config"
	"github.com/frank-vpl/gnet-cli/internal/logger"
	"github.com/frank-vpl/gnet-cli/internal/network"
	"github.com/frank-vpl/gnet-cli/internal/proxy"
	"github.com/frank-vpl/gnet-cli/internal/share"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).Padding(0, 1)
	accent      = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	dim         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingTop(1)
	bannerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 2)
	fieldStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
)

type step int

const (
	stepProxy step = iota
	stepHTTPPort
	stepSocksPort
	stepBind
	stepMethod
	stepConfirm
	stepRunning
)

type tuiModel struct {
	cfg     config.ProxyConfig
	step    step
	ifaces  []network.InterfaceInfo
	methods []share.Method

	bindChoices   []string
	methodChoices []string

	proxyChoice  int
	bindChoice   int
	methodChoice int

	input textinput.Model
	err   string

	server       *proxy.Server
	logRing      *LogRing
	logLines     []string
	tick         int
	firewallMsg  string
	reachableIPs []string
}

type tickMsg time.Time

func RunWizard(initial config.ProxyConfig) error {
	ifaces, _ := network.ListInterfaces()
	methods := share.Methods()

	bindChoices := []string{"0.0.0.0 (all interfaces)"}
	for _, i := range ifaces {
		for _, ip := range i.IPs {
			if network.IsIP(ip) && !strings.Contains(ip, ":") {
				bindChoices = append(bindChoices, fmt.Sprintf("%s (%s)", ip, i.Name))
			}
		}
	}

	methodChoices := make([]string, 0, len(methods))
	for _, m := range methods {
		methodChoices = append(methodChoices, fmt.Sprintf("%s - %s", m.ID, m.Name))
	}

	ti := textinput.New()
	ti.Placeholder = "8080"
	ti.CharLimit = 5
	ti.Width = 10

	m := tuiModel{
		cfg:           initial,
		ifaces:        ifaces,
		methods:       methods,
		bindChoices:   bindChoices,
		methodChoices: methodChoices,
		input:         ti,
		logRing:       NewLogRing(50),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func (m tuiModel) Init() tea.Cmd {
	return tea.Batch(tickCmd(), textinput.Blink)
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q", "esc":
			if m.step == stepRunning {
				if m.server != nil {
					m.server.Stop()
				}
				return m, tea.Quit
			}
		}
	case tickMsg:
		if m.step == stepRunning && m.server != nil {
			m.logLines = m.logRing.Lines(12)
			m.tick++
			return m, tickCmd()
		}
		return m, tickCmd()
	}

	switch m.step {
	case stepProxy:
		return m.updateProxyChoice(msg)
	case stepHTTPPort, stepSocksPort:
		return m.updatePort(msg)
	case stepBind:
		return m.updateBind(msg)
	case stepMethod:
		return m.updateMethod(msg)
	case stepConfirm:
		return m.updateConfirm(msg)
	case stepRunning:
		return m, nil
	}
	return m, nil
}

func (m *tuiModel) updateProxyChoice(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.proxyChoice > 0 {
				m.proxyChoice--
			}
		case "down", "j":
			if m.proxyChoice < 2 {
				m.proxyChoice++
			}
		case "enter":
			m.cfg.HTTPEnabled = m.proxyChoice == 0 || m.proxyChoice == 2
			m.cfg.Socks5Enabled = m.proxyChoice == 1 || m.proxyChoice == 2
			if m.cfg.HTTPEnabled {
				m.step = stepHTTPPort
				m.input.SetValue(strconv.Itoa(m.cfg.HTTPPort))
				m.input.Placeholder = "8080"
				return m, textinput.Blink
			}
			m.step = stepSocksPort
			m.input.SetValue(strconv.Itoa(m.cfg.Socks5Port))
			m.input.Placeholder = "1080"
			return m, textinput.Blink
		}
	}
	return m, nil
}

func (m *tuiModel) updatePort(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			port, err := strconv.Atoi(strings.TrimSpace(m.input.Value()))
			if err != nil || port < 1 || port > 65535 {
				m.err = "Invalid port. Use 1-65535."
				return m, nil
			}
			m.err = ""
			if m.step == stepHTTPPort {
				m.cfg.HTTPPort = port
				if m.cfg.Socks5Enabled {
					m.step = stepSocksPort
					m.input.SetValue(strconv.Itoa(m.cfg.Socks5Port))
					m.input.Placeholder = "1080"
					return m, textinput.Blink
				}
				m.step = stepBind
				return m, nil
			}
			m.cfg.Socks5Port = port
			m.step = stepBind
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *tuiModel) updateBind(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.bindChoice > 0 {
				m.bindChoice--
			}
		case "down", "j":
			if m.bindChoice < len(m.bindChoices)-1 {
				m.bindChoice++
			}
		case "enter":
			choice := m.bindChoices[m.bindChoice]
			if idx := strings.Index(choice, " "); idx > 0 {
				m.cfg.BindAddr = choice[:idx]
			} else {
				m.cfg.BindAddr = "0.0.0.0"
			}
			m.step = stepMethod
			return m, nil
		}
	}
	return m, nil
}

func (m *tuiModel) updateMethod(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.methodChoice > 0 {
				m.methodChoice--
			}
		case "down", "j":
			if m.methodChoice < len(m.methodChoices)-1 {
				m.methodChoice++
			}
		case "enter":
			m.cfg.SharingMethod = m.methods[m.methodChoice].ID
			m.step = stepConfirm
			return m, nil
		}
	}
	return m, nil
}

func (m *tuiModel) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "y", "Y", "enter":
			if err := m.cfg.Validate(); err != nil {
				m.err = err.Error()
				return m, nil
			}
			_ = config.Save(m.cfg)

			ports := activePorts(&m.cfg)
			fw := network.EnsureFirewall(ports)
			m.firewallMsg = fw.Message

			m.reachableIPs = m.collectReachableIPs()

			ring := m.logRing
			inner := logger.NewTextHandler(logger.LevelDebug)
			logger.SetHandler(&ringHandler{ring: ring, inner: inner})

			srv := proxy.New(m.cfg)
			ctx, cancel := context.WithCancel(context.Background())
			if err := srv.Start(ctx); err != nil {
				m.err = "Failed to start: " + err.Error()
				cancel()
				return m, nil
			}
			m.server = srv
			m.step = stepRunning
			go func() {
				sigCh := make(chan os.Signal, 1)
				signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
				<-sigCh
				srv.Stop()
				cancel()
			}()
			return m, nil
		case "n", "N":
			m.step = stepProxy
			return m, nil
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	var b strings.Builder
	b.WriteString(bannerStyle.Render("GNet CLI  —  Share your VPN everywhere") + "\n\n")

	switch m.step {
	case stepProxy:
		b.WriteString(titleStyle.Render("1/5  Select proxy type") + "\n\n")
		opts := []string{"HTTP only", "SOCKS5 only", "Both HTTP & SOCKS5 (recommended)"}
		for i, o := range opts {
			mark := " "
			if i == m.proxyChoice {
				mark = accent.Render(">")
			}
			b.WriteString(fmt.Sprintf("  %s %s\n", mark, o))
		}
		b.WriteString(helpStyle.Render("↑/↓ to choose, Enter to confirm"))
	case stepHTTPPort:
		b.WriteString(titleStyle.Render("2/5  HTTP proxy port") + "\n\n")
		b.WriteString("  Port: " + m.input.View() + "\n")
		b.WriteString(helpStyle.Render("Enter to confirm (default 8080)"))
	case stepSocksPort:
		b.WriteString(titleStyle.Render("3/5  SOCKS5 proxy port") + "\n\n")
		b.WriteString("  Port: " + m.input.View() + "\n")
		b.WriteString(helpStyle.Render("Enter to confirm (default 1080)"))
	case stepBind:
		b.WriteString(titleStyle.Render("4/5  Bind address") + "\n\n")
		for i, c := range m.bindChoices {
			mark := " "
			if i == m.bindChoice {
				mark = accent.Render(">")
			}
			b.WriteString(fmt.Sprintf("  %s %s\n", mark, c))
		}
		b.WriteString(helpStyle.Render("↑/↓ to choose, Enter to confirm"))
	case stepMethod:
		b.WriteString(titleStyle.Render("5/5  Sharing method") + "\n\n")
		for i, c := range m.methodChoices {
			mark := " "
			if i == m.methodChoice {
				mark = accent.Render(">")
			}
			b.WriteString(fmt.Sprintf("  %s %s\n", mark, c))
		}
		if m.methodChoice < len(m.methods) {
			desc := m.methods[m.methodChoice].Description
			b.WriteString("\n" + dim.Render("  "+wrap(desc, 70)))
			if m.methods[m.methodChoice].RequiresRoot {
				b.WriteString("\n" + dim.Render("  Note: may require administrator/root privileges"))
			}
		}
		b.WriteString(helpStyle.Render("↑/↓ to choose, Enter to confirm"))
	case stepConfirm:
		b.WriteString(titleStyle.Render("Confirm configuration") + "\n\n")
		b.WriteString(summary(&m.cfg))
		b.WriteString("\n  Start proxy servers? " + fieldStyle.Render("[Y/n]") + "\n")
		b.WriteString(helpStyle.Render("Y to start, n to reconfigure"))
	case stepRunning:
		st := m.server.Status()
		b.WriteString(titleStyle.Render("Proxy servers running") + "\n\n")
		if st.HTTPRunning {
			b.WriteString(fmt.Sprintf("  %s HTTP   -> %s:%d\n", accent.Render("●"), st.BindAddr, st.HTTPPort))
		}
		if st.Socks5Running {
			b.WriteString(fmt.Sprintf("  %s SOCKS5 -> %s:%d\n", accent.Render("●"), st.BindAddr, st.Socks5Port))
		}
		b.WriteString(fmt.Sprintf("  Active connections: %s\n", accent.Render(strconv.FormatInt(st.ActiveConns, 10))))

		b.WriteString("\n" + dim.Render("  Reachable on this machine (set proxy on your device to one of these):") + "\n")
		for _, ip := range m.reachableIPs {
			ports := ""
			if st.HTTPRunning {
				ports = fmt.Sprintf(":%d", st.HTTPPort)
			}
			if st.Socks5Running {
				if ports != "" {
					ports += " / "
				}
				ports += fmt.Sprintf(":%d", st.Socks5Port)
			}
			b.WriteString(fmt.Sprintf("    %s %s%s\n", accent.Render("→"), ip, dim.Render(ports)))
		}

		if m.firewallMsg != "" {
			fwStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
			b.WriteString("\n" + fwStyle.Render("  Firewall: "+wrap(m.firewallMsg, 70)) + "\n")
		}

		b.WriteString("\n" + dim.Render("  Live log:") + "\n")
		for _, l := range m.logLines {
			b.WriteString("  " + dim.Render(l) + "\n")
		}
		b.WriteString(helpStyle.Render("q or Ctrl+C to stop"))
	}

	if m.err != "" {
		b.WriteString("\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Error: "+m.err))
	}
	return b.String()
}

func summary(c *config.ProxyConfig) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  HTTP:        %s (port %d)\n", boolStr(c.HTTPEnabled), c.HTTPPort))
	b.WriteString(fmt.Sprintf("  SOCKS5:      %s (port %d)\n", boolStr(c.Socks5Enabled), c.Socks5Port))
	b.WriteString(fmt.Sprintf("  Bind:        %s\n", c.BindAddr))
	if method, ok := share.Get(c.SharingMethod); ok {
		b.WriteString(fmt.Sprintf("  Method:      %s\n", method.Name))
		b.WriteString("\n" + dim.Render("  Steps:") + "\n")
		for i, s := range method.Steps {
			b.WriteString("    " + dim.Render(fmt.Sprintf("%d. %s", i+1, wrap(s, 70))) + "\n")
		}
	}
	return b.String()
}

func boolStr(b bool) string {
	if b {
		return accent.Render("enabled")
	}
	return dim.Render("disabled")
}

func wrap(s string, width int) string {
	if len(s) <= width {
		return s
	}
	var b strings.Builder
	words := strings.Fields(s)
	line := ""
	for _, w := range words {
		if len(line)+len(w)+1 > width {
			b.WriteString(line + "\n  ")
			line = w
		} else {
			if line != "" {
				line += " "
			}
			line += w
		}
	}
	b.WriteString(line)
	return b.String()
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func activePorts(c *config.ProxyConfig) []int {
	var ports []int
	if c.HTTPEnabled {
		ports = append(ports, c.HTTPPort)
	}
	if c.Socks5Enabled {
		ports = append(ports, c.Socks5Port)
	}
	return ports
}

func (m *tuiModel) collectReachableIPs() []string {
	var ips []string
	seen := map[string]bool{}
	for _, ifc := range m.ifaces {
		if !strings.Contains(strings.ToLower(ifc.Flags), "up") {
			continue
		}
		if strings.Contains(strings.ToLower(ifc.Flags), "loopback") {
			continue
		}
		for _, ip := range ifc.IPs {
			if strings.Contains(ip, ":") {
				continue
			}
			if !seen[ip] {
				seen[ip] = true
				ips = append(ips, ip)
			}
		}
	}
	return ips
}
