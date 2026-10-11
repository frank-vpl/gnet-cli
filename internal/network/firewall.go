package network

import (
	"errors"
	"fmt"
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

func EnsureFirewall(ports []int) FirewallResult {
	res := FirewallResult{OS: runtime.GOOS, Ports: ports}

	if runtime.GOOS != "windows" {
		res.Message = firewallHint(ports)
		return res
	}

	portList := joinPorts(ports)
	ruleName := "GNet CLI"

	exists := firewallRuleExists(ruleName)
	if exists {
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

func firewallRuleExists(name string) bool {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name="+name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(output)), "localport")
}

func RemoveFirewallRule() error {
	if runtime.GOOS != "windows" {
		return nil
	}
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name=GNet CLI")
	return cmd.Run()
}

func firewallHint(ports []int) string {
	portList := joinPorts(ports)
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("Allow TCP ports %s in your firewall, e.g.:\n  sudo ufw allow %s/tcp", portList, portList)
	case "darwin":
		return fmt.Sprintf("Allow incoming connections on TCP ports %s in System Settings > Network > Firewall.", portList)
	default:
		return fmt.Sprintf("Ensure TCP ports %s are open in your firewall.", portList)
	}
}

func joinPorts(ports []int) string {
	strs := make([]string, len(ports))
	for i, p := range ports {
		strs[i] = strconv.Itoa(p)
	}
	return strings.Join(strs, ",")
}

func errExitCode(err error) int {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return -1
}
