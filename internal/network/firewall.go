package network

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/frank-vpl/gnet-cli/internal/logger"
)

type FirewallResult struct {
	OS         string
	RuleAdded  bool
	NeedsAdmin bool
	Message    string
	Ports      []int
}

const ruleName = "GNet CLI"

func EnsureFirewall(ports []int) FirewallResult {
	res := FirewallResult{OS: runtime.GOOS, Ports: ports}

	switch runtime.GOOS {
	case "windows":
		return ensureWindowsFirewall(ports)
	case "linux":
		return ensureLinuxFirewall(ports)
	case "darwin":
		return ensureDarwinFirewall(ports)
	default:
		res.Message = firewallHint(ports)
		return res
	}
}

// ---------------------------------------------------------------------------
// Windows
// ---------------------------------------------------------------------------

func ensureWindowsFirewall(ports []int) FirewallResult {
	res := FirewallResult{OS: "windows", Ports: ports}
	portList := joinPorts(ports)

	if firewallRuleExistsWindows(ruleName) {
		res.RuleAdded = true
		res.Message = fmt.Sprintf("Windows Firewall rule '%s' already allows ports %s", ruleName, portList)
		return res
	}

	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name="+ruleName,
		"dir=in",
		"action=allow",
		"protocol=TCP",
		"localport="+portList,
		"profile=any",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		out := strings.TrimSpace(string(output))
		lower := strings.ToLower(out)
		if strings.Contains(lower, "requires elevation") ||
			strings.Contains(lower, "administrator") ||
			strings.Contains(lower, "access is denied") ||
			errExitCode(err) == 1 {
			res.NeedsAdmin = true
			res.Message = fmt.Sprintf(
				"Windows Firewall blocked the rule (needs Administrator).\n"+
					"  Option 1 (easiest): close this and run gnet again as Administrator  (right-click > Run as administrator).\n"+
					"  Option 2: run this in an elevated PowerShell:\n"+
					"    netsh advfirewall firewall add rule name=\"GNet CLI\" dir=in action=allow protocol=TCP localport=%s profile=any",
				portList)
			logger.Warn("firewall rule requires admin", "err", err)
			return res
		}
		res.Message = fmt.Sprintf("Could not add firewall rule: %v\n%s", err, out)
		logger.Warn("firewall rule add failed", "err", err, "output", out)
		return res
	}

	res.RuleAdded = true
	res.Message = fmt.Sprintf("Added Windows Firewall rule '%s' allowing TCP ports %s", ruleName, portList)
	logger.Info("firewall rule added", "ports", portList)
	return res
}

func firewallRuleExistsWindows(name string) bool {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name="+name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "localport")
}

// ---------------------------------------------------------------------------
// Linux  (ufw -> firewalld -> nftables -> iptables)
// ---------------------------------------------------------------------------

func ensureLinuxFirewall(ports []int) FirewallResult {
	res := FirewallResult{OS: "linux", Ports: ports}

	backends := []func([]int) (string, bool, bool){
		tryUFW,
		tryFirewalld,
		tryNFT,
		tryIPTables,
	}

	for _, fn := range backends {
		msg, ok, needsAdmin := fn(ports)
		if ok {
			res.RuleAdded = true
			res.Message = msg
			logger.Info("firewall rule added", "ports", joinPorts(ports))
			return res
		}
		if needsAdmin {
			res.NeedsAdmin = true
			res.Message = msg
			return res
		}
	}

	res.Message = firewallHint(ports)
	return res
}

// ufw — Ubuntu / Debian and derivatives
func tryUFW(ports []int) (string, bool, bool) {
	if !commandExists("ufw") {
		return "", false, false
	}
	portList := joinPorts(ports)

	out, err, needsAdmin := runWithSudo("ufw", "allow", portList+"/tcp")
	if err == nil {
		return fmt.Sprintf("Added ufw rule allowing TCP ports %s", portList), true, false
	}
	if needsAdmin {
		return fmt.Sprintf(
			"ufw needs root to allow ports %s.\n"+
				"  Run: sudo ufw allow %s/tcp",
			portList, portList), false, true
	}
	logger.Debug("ufw add failed", "err", err, "out", out)
	return out, false, false
}

// firewalld — RHEL / CentOS / Fedora / openSUSE
func tryFirewalld(ports []int) (string, bool, bool) {
	if !commandExists("firewall-cmd") {
		return "", false, false
	}

	// Skip if firewalld is not running.
	stateOut, err := runCmd("firewall-cmd", "--state")
	if err != nil || !strings.Contains(strings.ToLower(stateOut), "running") {
		return "", false, false
	}

	args := []string{"--permanent"}
	for _, p := range ports {
		args = append(args, fmt.Sprintf("--add-port=%d/tcp", p))
	}

	out, err, needsAdmin := runWithSudo("firewall-cmd", args...)
	if err != nil {
		portList := joinPorts(ports)
		if needsAdmin {
			return fmt.Sprintf(
				"firewalld needs root.\n"+
					"  Run: sudo firewall-cmd --permanent --add-port=%s/tcp && sudo firewall-cmd --reload",
				portList), false, true
		}
		logger.Debug("firewalld add failed", "err", err, "out", out)
		return out, false, false
	}

	// Apply the permanent rule to the running configuration.
	_, _, _ = runWithSudo("firewall-cmd", "--reload")

	portList := joinPorts(ports)
	return fmt.Sprintf("Added firewalld rule allowing TCP ports %s", portList), true, false
}

// nftables — modern Linux (replaces iptables)
func tryNFT(ports []int) (string, bool, bool) {
	if !commandExists("nft") {
		return "", false, false
	}
	portSet := strings.Join(intsToStrings(ports), ",")

	// Idempotent setup: table + chain (errors are ignored if they already exist),
	// then append the accept rule.
	setup := [][]string{
		{"add", "table", "inet", "gnet"},
		{"add", "chain", "inet", "gnet", "input", "{", "type", "filter", "hook", "input", "priority", "0", ";"},
		{"add", "rule", "inet", "gnet", "input", "tcp", "dport", "{" + portSet + "}", "accept"},
	}

	lastOut := ""
	for _, c := range setup {
		out, err, needsAdmin := runWithSudo("nft", c...)
		if err != nil {
			if needsAdmin {
				portList := joinPorts(ports)
				return fmt.Sprintf(
					"nftables needs root.\n"+
						"  Run:\n"+
						"    sudo nft add table inet gnet\n"+
						"    sudo nft add chain inet gnet input '{ type filter hook input priority 0 ; }'\n"+
						"    sudo nft add rule inet gnet input tcp dport {%s} accept",
					portList), false, true
			}
			lastOut = out
			// table/chain may already exist — keep going.
		}
	}

	portList := joinPorts(ports)
	_ = lastOut
	return fmt.Sprintf("Added nftables rule allowing TCP ports %s", portList), true, false
}

// iptables — legacy Linux firewall
func tryIPTables(ports []int) (string, bool, bool) {
	if !commandExists("iptables") {
		return "", false, false
	}
	portList := joinPorts(ports)

	out, err, needsAdmin := runWithSudo("iptables",
		"-A", "INPUT", "-p", "tcp", "-m", "multiport", "--dports", portList, "-j", "ACCEPT")
	if err == nil {
		return fmt.Sprintf("Added iptables rule allowing TCP ports %s", portList), true, false
	}
	if needsAdmin {
		return fmt.Sprintf(
			"iptables needs root.\n"+
				"  Run: sudo iptables -A INPUT -p tcp -m multiport --dports %s -j ACCEPT",
			portList), false, true
	}
	logger.Debug("iptables add failed", "err", err, "out", out)
	return out, false, false
}

// ---------------------------------------------------------------------------
// macOS  (Application Firewall via socketfilterfw)
// ---------------------------------------------------------------------------

func ensureDarwinFirewall(ports []int) FirewallResult {
	res := FirewallResult{OS: "darwin", Ports: ports}

	sfw := "/usr/libexec/ApplicationFirewall/socketfilterfw"
	if _, err := os.Stat(sfw); err != nil {
		res.Message = firewallHint(ports)
		return res
	}

	binPath, err := os.Executable()
	if err != nil || binPath == "" {
		binPath = os.Args[0]
	}

	// Make sure the firewall is on (best-effort).
	_, _, _ = runWithSudo(sfw, "--setglobalstate", "on")

	// Add the binary to the firewall list (idempotent: already added = ok).
	out, err, needsAdmin := runWithSudo(sfw, "--add", binPath)
	alreadyAdded := strings.Contains(strings.ToLower(out), "already")
	if err != nil && !alreadyAdded {
		if needsAdmin {
			res.NeedsAdmin = true
			res.Message = fmt.Sprintf(
				"macOS Application Firewall needs root.\n"+
					"  Run:\n"+
					"    sudo %s --setglobalstate on\n"+
					"    sudo %s --add %q\n"+
					"    sudo %s --unblockapp %q",
				sfw, sfw, binPath, sfw, binPath)
			return res
		}
		logger.Debug("socketfilterfw add failed", "err", err, "out", out)
	}

	// Ensure the app is allowed (not blocked).
	_, _, _ = runWithSudo(sfw, "--unblockapp", binPath)

	res.RuleAdded = true
	res.Message = "Allowed GNet CLI in the macOS Application Firewall"
	logger.Info("macos firewall rule added", "binary", binPath)
	return res
}

// ---------------------------------------------------------------------------
// Removal (cleanup) — mirrors EnsureFirewall across platforms
// ---------------------------------------------------------------------------

func RemoveFirewallRule() error {
	switch runtime.GOOS {
	case "windows":
		if !firewallRuleExistsWindows(ruleName) {
			return nil
		}
		cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name="+ruleName)
		return cmd.Run()
	case "linux":
		return removeLinuxFirewall()
	case "darwin":
		sfw := "/usr/libexec/ApplicationFirewall/socketfilterfw"
		binPath, err := os.Executable()
		if err != nil || binPath == "" {
			binPath = os.Args[0]
		}
		_, _, _ = runWithSudo(sfw, "--remove", binPath)
		return nil
	default:
		return nil
	}
}

func removeLinuxFirewall() error {
	// Default ports used for cleanup (RemoveFirewallRule takes no args).
	defaultPorts := []int{8080, 1080}

	if commandExists("ufw") {
		for _, p := range defaultPorts {
			_, _, _ = runWithSudo("ufw", "delete", "allow", strconv.Itoa(p)+"/tcp")
		}
	}
	if commandExists("nft") {
		_, _, _ = runWithSudo("nft", "delete", "table", "inet", "gnet")
	}
	if commandExists("firewall-cmd") {
		stateOut, _ := runCmd("firewall-cmd", "--state")
		if strings.Contains(strings.ToLower(stateOut), "running") {
			args := []string{"--permanent"}
			for _, p := range defaultPorts {
				args = append(args, fmt.Sprintf("--remove-port=%d/tcp", p))
			}
			_, _, _ = runWithSudo("firewall-cmd", args...)
			_, _, _ = runWithSudo("firewall-cmd", "--reload")
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func firewallHint(ports []int) string {
	portList := joinPorts(ports)
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf(
			"Allow TCP ports %s in your firewall, e.g.:\n"+
				"  sudo ufw allow %s/tcp            # ufw (Ubuntu/Debian)\n"+
				"  sudo firewall-cmd --permanent --add-port=%s/tcp && sudo firewall-cmd --reload   # firewalld\n"+
				"  sudo nft add rule inet filter input tcp dport {%s} accept   # nftables\n"+
				"  sudo iptables -A INPUT -p tcp -m multiport --dports %s -j ACCEPT   # iptables",
			portList, portList, portList, portList, portList)
	case "darwin":
		return fmt.Sprintf(
			"Allow incoming connections on TCP ports %s for GNet CLI in\n"+
				"  System Settings > Network > Firewall, or run:\n"+
				"    sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add %q\n"+
				"    sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp %q",
			portList, os.Args[0], os.Args[0])
	default:
		return fmt.Sprintf("Ensure TCP ports %s are open in your firewall.", portList)
	}
}

// runWithSudo runs a command directly first, then retries with `sudo -n`
// (non-interactive). Returns output, error, and whether the failure was due to
// missing privileges.
func runWithSudo(name string, args ...string) (out string, err error, needsAdmin bool) {
	out, err = runCmd(name, args...)
	if err == nil {
		return out, nil, false
	}

	sudoOut, sudoErr := runCmd("sudo", append([]string{"-n", name}, args...)...)
	if sudoErr == nil {
		return sudoOut, nil, false
	}

	combined := out + "\n" + sudoOut
	low := strings.ToLower(combined)
	if strings.Contains(low, "password is required") ||
		strings.Contains(low, "permission denied") ||
		strings.Contains(low, "not in the sudoers") ||
		strings.Contains(low, "a password is required") {
		return combined, sudoErr, true
	}
	return combined, sudoErr, false
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func joinPorts(ports []int) string {
	strs := make([]string, len(ports))
	for i, p := range ports {
		strs[i] = strconv.Itoa(p)
	}
	return strings.Join(strs, ",")
}

func intsToStrings(ports []int) []string {
	strs := make([]string, len(ports))
	for i, p := range ports {
		strs[i] = strconv.Itoa(p)
	}
	return strs
}

func errExitCode(err error) int {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return -1
}
