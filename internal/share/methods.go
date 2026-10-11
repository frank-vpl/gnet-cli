package share

import (
	"fmt"
	"runtime"
	"strings"
)

type Method struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RequiresRoot bool   `json:"requires_root"`
	Steps       []string `json:"steps"`
}

const (
	MethodLAN     = "lan"
	MethodWiFi    = "wifi"
	MethodHotspot = "hotspot"
	MethodUSB     = "usb"
	MethodRouter  = "router"
)

func Methods() []Method {
	os := runtime.GOOS
	return []Method{
		lanMethod(),
		wifiMethod(),
		hotspotMethod(os),
		usbMethod(os),
		routerMethod(),
	}
}

func Get(id string) (Method, bool) {
	for _, m := range Methods() {
		if m.ID == id {
			return m, true
		}
	}
	return Method{}, false
}

func lanMethod() Method {
	return Method{
		ID:          MethodLAN,
		Name:        "LAN / Ethernet",
		Description: "Share over the local wired network. Devices on the same LAN use your machine's LAN IP.",
		RequiresRoot: false,
		Steps: []string{
			"Run `gnet serve` and note your LAN IP (shown by `gnet status` or `gnet interfaces`).",
			"On the target device, set HTTP proxy to <LAN-IP>:8080 or SOCKS5 to <LAN-IP>:1080.",
			"Ensure your machine's firewall allows inbound traffic on the proxy ports.",
		},
	}
}

func wifiMethod() Method {
	return Method{
		ID:          MethodWiFi,
		Name:        "Wi-Fi (same network)",
		Description: "Share over an existing Wi-Fi network. Other devices on the same Wi-Fi use your Wi-Fi IP.",
		RequiresRoot: false,
		Steps: []string{
			"Connect your machine to the Wi-Fi network.",
			"Run `gnet serve` and note your Wi-Fi IP.",
			"On target devices on the same Wi-Fi, set the proxy to <WI-FI-IP>:<port>.",
		},
	}
}

func hotspotMethod(os string) Method {
	m := Method{
		ID:          MethodHotspot,
		Name:        "Mobile Hotspot / AP",
		Description: "Turn your machine into a Wi-Fi access point so devices connect directly to it.",
		RequiresRoot: true,
	}
	switch os {
	case "windows":
		m.Steps = []string{
			"Open Settings > Network & Internet > Mobile hotspot.",
			"Enable 'Share my Internet connection with other devices'.",
			"Note the network name and password, connect target devices to it.",
			"Run `gnet serve` and use the hotspot IP shown in `gnet interfaces`.",
		}
	case "linux":
		m.Steps = []string{
			"Use NetworkManager: `nmcli device wifi hotspot ifname wlan0 ssid GNet password <pass>`.",
			"Or use create_ap: `sudo create_ap wlan0 eth0 GNet <pass>`.",
			"Run `gnet serve` and set the proxy IP to the AP interface IP on target devices.",
		}
	case "darwin":
		m.Steps = []string{
			"Open System Settings > General > Sharing > Internet Sharing.",
			"Share your connection from (e.g. Ethernet/VPN) to computers using Wi-Fi.",
			"Configure Wi-Fi Options (name/password), then enable Internet Sharing.",
			"Run `gnet serve` and use the Mac's Wi-Fi IP on connected devices.",
		}
	default:
		m.Steps = []string{"Enable a Wi-Fi access point on your OS, then run `gnet serve`."}
	}
	return m
}

func usbMethod(os string) Method {
	m := Method{
		ID:          MethodUSB,
		Name:        "USB Tethering / RNDIS",
		Description: "Share over a USB cable using USB tethering or RNDIS networking.",
		RequiresRoot: true,
	}
	switch os {
	case "windows":
		m.Steps = []string{
			"Connect the device via USB and enable USB tethering on it.",
			"Windows will install an RNDIS network adapter automatically.",
			"Run `gnet interfaces` to find the RNDIS adapter IP (e.g. 192.168.x.x).",
			"Run `gnet serve` and configure the proxy on the device using that IP.",
		}
	case "linux":
		m.Steps = []string{
			"Connect the device via USB and enable USB tethering.",
			"Run `ip addr` to find the usb0 interface IP.",
			"Run `gnet serve` and set the device proxy to the usb0 IP:port.",
		}
	case "darwin":
		m.Steps = []string{
			"Connect the device via USB and enable USB tethering.",
			"Find the USB network interface IP in System Settings > Network.",
			"Run `gnet serve` and set the device proxy to that IP:port.",
		}
	default:
		m.Steps = []string{"Enable USB tethering, find the USB interface IP, then run `gnet serve`."}
	}
	return m
}

func routerMethod() Method {
	return Method{
		ID:          MethodRouter,
		Name:        "Router / Port Forward",
		Description: "Expose the proxy through your router's port forwarding to reach devices across networks.",
		RequiresRoot: false,
		Steps: []string{
			"Log in to your router admin panel.",
			"Add a port-forwarding rule: external port -> <your-machine-IP>:8080 (HTTP) and :1080 (SOCKS5).",
			"Run `gnet serve` on your machine.",
			"Target devices use your router's public IP and the forwarded ports.",
			"Optional: enable UPnP on the router for automatic forwarding.",
		},
	}
}

func FormatMethods() string {
	var b strings.Builder
	for _, m := range Methods() {
		fmt.Fprintf(&b, "[%s] %s\n", m.ID, m.Name)
		fmt.Fprintf(&b, "  %s\n", m.Description)
		if m.RequiresRoot {
			b.WriteString("  (may require administrator/root privileges)\n")
		}
		for i, s := range m.Steps {
			fmt.Fprintf(&b, "    %d. %s\n", i+1, s)
		}
		b.WriteString("\n")
	}
	return b.String()
}
