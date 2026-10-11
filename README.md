<div align="center">

# 🌐 GNet CLI

**Share your VPN over HTTP/SOCKS5 to any device on your network.**

Turn your machine into a lightweight proxy server so phones, tablets, smart TVs, and other devices can route traffic through your VPN — over LAN, Wi-Fi, hotspot, USB, or router.

![Go](https://img.shields.io/badge/Go-1.27+-00ADD8?logo=go&logoColor=white&style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-7B2CBF?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-3DDC84?style=for-the-badge)
![Version](https://img.shields.io/badge/version-v1.0-1F6FEB?style=for-the-badge)
[![Build and Test](https://img.shields.io/github/actions/workflow/status/frank-vpl/gnet-cli/build.yaml?label=Build%20%26%20Test&logo=githubactions&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/actions/workflows/build.yaml)
[![Latest Release](https://img.shields.io/github/v/release/frank-vpl/gnet-cli?display_name=release&label=Latest%20Release&logo=github&style=for-the-badge&color=1F6FEB)](https://github.com/frank-vpl/gnet-cli/releases/latest)

[⬇️ Download](#️-download) · [🚀 Quick Start](#-quick-start) · [📚 Commands](#-commands) · [🔌 Sharing Methods](#-sharing-methods)

</div>

---

## ✨ Features

- 🔁 **HTTP & SOCKS5 Proxy** — Run one or both simultaneously with custom ports.
- 🧙 **Interactive Wizard** — A polished TUI walks you through setup step by step.
- 🛡️ **Firewall Automation** — Auto-creates the Windows Firewall rule; gives hints for Linux/macOS.
- 📡 **Auto IP Detection** — Shows every reachable address so devices can connect.
- 🔌 **5 Sharing Methods** — LAN, Wi-Fi, Hotspot (AP), USB Tethering, Router Port Forward.
- 💾 **Persistent Config** — Saves your settings to `~/.gnet-cli/config.json`.
- 🐛 **Debug Logging** — Optional file logging and live connection tracking.
- ⏹️ **Graceful Shutdown** — Cleanly closes listeners and active connections.

---

## 📦 Installation

### Build from source

```bash
git clone https://github.com/frank-vpl/gnet-cli.git
cd gnet-cli
go build -o gnet .
```

> On Windows, the output is `gnet.exe`. Add the directory to your `PATH` or run it directly.

### Requirements

- [Go](https://go.dev/) 1.27 or later

---

## ⬇️ Download

Grab the latest pre-built binary for your platform — no Go installation required.

| Platform | Architecture | Download |
|----------|-------------|----------|
| 🪟 Windows | x64 | [![Download](https://img.shields.io/badge/gnet-windows--x64.zip-0078D6?logo=windows&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-windows-x64.zip) |
| 🪟 Windows | x86 | [![Download](https://img.shields.io/badge/gnet-windows--x86.zip-0078D6?logo=windows&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-windows-x86.zip) |
| 🐧 Linux | x64 | [![Download](https://img.shields.io/badge/gnet-linux--x64.zip-FCC624?logo=linux&logoColor=black&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-linux-x64.zip) |
| 🐧 Linux | x86 | [![Download](https://img.shields.io/badge/gnet-linux--x86.zip-FCC624?logo=linux&logoColor=black&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-linux-x86.zip) |
| 🍎 macOS | x64 (Intel) | [![Download](https://img.shields.io/badge/gnet-macos--x64.zip-000000?logo=apple&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-macos-x64.zip) |
| 🍎 macOS | ARM64 (Apple Silicon) | [![Download](https://img.shields.io/badge/gnet-macos--arm64.zip-000000?logo=apple&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-macos-arm64.zip) |

> 💡 Each ZIP contains the `gnet` binary (or `gnet.exe`), `README.md`, and `LICENSE`.
> Extract it, add the folder to your `PATH`, and run `gnet`.

All releases: [github.com/frank-vpl/gnet-cli/releases](https://github.com/frank-vpl/gnet-cli/releases)

---

## 🚀 Quick Start

Launch the interactive wizard with no arguments:

```bash
gnet
```

Or start the proxy servers directly using saved/default config:

```bash
gnet serve
```

You'll see output like:

```
  GNet CLI proxy running
  HTTP   -> 0.0.0.0:8080
  SOCKS5 -> 0.0.0.0:1080
  LAN IP -> 192.168.1.42
  Press Ctrl+C to stop
```

Point your device's proxy settings at one of the listed IPs:

| Type   | Example             |
|--------|---------------------|
| HTTP   | `192.168.1.42:8080` |
| SOCKS5 | `192.168.1.42:1080` |

---

## 📚 Commands

| Command            | Description                                         |
|--------------------|-----------------------------------------------------|
| `gnet`             | Launch the interactive setup wizard                 |
| `gnet serve`       | Start HTTP/SOCKS5 proxy servers                     |
| `gnet config`      | Show the saved configuration                        |
| `gnet interfaces`  | List all network interfaces and their IP addresses  |
| `gnet methods`     | List every sharing method with setup steps          |

### `gnet serve` flags

| Flag            | Type   | Default | Description                    |
|-----------------|--------|---------|--------------------------------|
| `--http-port`   | int    | `8080`  | HTTP proxy port                |
| `--socks-port`  | int    | `1080`  | SOCKS5 proxy port              |
| `--bind`        | string | `0.0.0.0` | Bind address                 |
| `--http`        | bool   | `true`  | Enable HTTP proxy              |
| `--socks5`      | bool   | `true`  | Enable SOCKS5 proxy            |
| `--log-file`    | string | –       | Write logs to a file           |
| `--debug`       | bool   | `false` | Enable debug-level logging     |

**Examples**

```bash
# HTTP only on port 3128
gnet serve --http-port 3128 --socks5=false

# SOCKS5 only, bound to a specific interface
gnet serve --bind 192.168.1.42 --http=false

# Debug logging to a file
gnet serve --debug --log-file gnet.log
```

---

## 🔌 Sharing Methods

| Method | Requires Admin | Best For |
|--------|----------------|----------|
| 🌐 **LAN / Ethernet** | ❌ | Wired devices on the same network |
| 📶 **Wi-Fi (same network)** | ❌ | Devices on the same Wi-Fi |
| 📱 **Mobile Hotspot / AP** | ✅ | Connect devices directly to your machine's Wi-Fi |
| 🔌 **USB Tethering / RNDIS** | ✅ | Wired USB connection to phones/tablets |
| 🏠 **Router / Port Forward** | ❌ | Reach devices across networks |

Run `gnet methods` to see OS-specific setup steps for each method.

---

## ⚙️ Configuration

Config is stored at `~/.gnet-cli/config.json`:

```json
{
  "http_port": 8080,
  "socks5_port": 1080,
  "bind_addr": "0.0.0.0",
  "http_enabled": true,
  "socks5_enabled": true,
  "sharing_method": "lan"
}
```

**Validation rules**

- At least one proxy (HTTP or SOCKS5) must be enabled.
- Ports must be in the range `1–65535`.
- HTTP and SOCKS5 ports cannot be the same.

---

## 🛠️ Project Structure

```
gnet-cli/
├── main.go                    # Entry point
├── go.mod
├── go.sum
└── internal/
    ├── cli/
    │   ├── commands.go        # Cobra commands & flags
    │   ├── tui.go             # Bubble Tea interactive wizard
    │   └── logring.go         # In-memory ring buffer for live logs
    ├── config/
    │   └── config.go          # Config load/save/validate
    ├── logger/
    │   └── logger.go          # Structured logger
    ├── network/
    │   ├── network.go         # Interface enumeration & IP helpers
    │   └── firewall.go        # Firewall rule management (Windows + hints)
    ├── proxy/
    │   ├── server.go          # Listener lifecycle & graceful shutdown
    │   ├── http.go            # HTTP CONNECT proxy handler
    │   ├── socks5.go          # SOCKS5 proxy handler
    │   └── relay.go           # Bidirectional connection relay
    └── share/
        └── methods.go         # Sharing methods & OS-specific setup steps
```

---

## 🧪 Development

```bash
# Run directly
go run .

# Run tests
go test ./...

# Build
go build -o gnet .
```

---

## � CI/CD

The project uses **GitHub Actions** for automated testing, cross-compilation, and releases.

### Workflow: [Build and Test](.github/workflows/build.yaml)

| Job | Trigger | Description |
|-----|---------|-------------|
| `test` | push / PR | Runs `go test ./...` on Ubuntu, Windows, and macOS |
| `build` | push / PR / tag | Cross-compiles 6 binaries (Linux/Windows/macOS × x86/x64/arm64) and uploads ZIPs |
| `release` | tag `v*` or manual | Creates a GitHub Release with all binaries + a download table |

### Automated release process

1. Push a version tag, e.g. `git tag v1.0.0 && git push origin v1.0.0` (or run the workflow manually from the **Actions** tab).
2. The `build` job compiles `gnet` for all platforms and packages each binary with `README.md` + `LICENSE` into a ZIP.
3. The `release` job publishes a GitHub Release and attaches every ZIP as a downloadable asset.

### Status

- **Build & Test:** [![Build and Test](https://img.shields.io/github/actions/workflow/status/frank-vpl/gnet-cli/build.yaml?logo=githubactions&logoColor=white&style=flat-square)](https://github.com/frank-vpl/gnet-cli/actions/workflows/build.yaml)
- **Latest Release:** [![Latest Release](https://img.shields.io/github/v/release/frank-vpl/gnet-cli?display_name=release&logo=github&style=flat-square)](https://github.com/frank-vpl/gnet-cli/releases/latest)

---

## � License

Released under the MIT License. See [LICENSE](LICENSE) for details.

---

<div align="center">

Built with 🩷 by **Hossein Pira** using [Cobra](https://github.com/spf13/cobra), [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss)

</div>
