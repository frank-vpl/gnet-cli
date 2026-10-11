<div align="center">

[English](README.md) | فارسی

# 🌐 GNet CLI

**شبکه خصوصی مجازی (VPN) خود را از طریق HTTP/SOCKS5 با هر دستگاهی در شبکه خود به اشتراک بگذارید.**

سیستم خود را به یک سرور پروکسی سبک تبدیل کنید تا تلفن‌ها، تبلت‌ها، تلویزیون‌های هوشمند و سایر دستگاه‌ها بتوانند ترافیک خود را از طریق VPN شما — روی شبکه محلی (LAN)، وای‌فای، هات‌اسپات، USB یا روتر — هدایت کنند.

![Go](https://img.shields.io/badge/Go-1.27+-00ADD8?logo=go&logoColor=white&style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-7B2CBF?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-3DDC84?style=for-the-badge)
![Version](https://img.shields.io/badge/version-v1.0.2-1F6FEB?style=for-the-badge)
[![Build and Test](https://img.shields.io/github/actions/workflow/status/frank-vpl/gnet-cli/build.yaml?label=Build%20%26%20Test&logo=githubactions&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/actions/workflows/build.yaml)
[![Latest Release](https://img.shields.io/github/v/release/frank-vpl/gnet-cli?display_name=release&label=Latest%20Release&logo=github&style=for-the-badge&color=1F6FEB)](https://github.com/frank-vpl/gnet-cli/releases/latest)

[⬇️ دانلود](#️-دانلود) · [🚀 شروع سریع](#-شروع-سریع) · [📚 دستورها](#-دستورها) · [🔌 روش‌های اشتراک‌گذاری](#-روشهای-اشتراکگذاری)

</div>

---

## 🖼️ تصویر صفحه

<p align="center">
  <img src="image.png" alt="GNet CLI Wizard" width="700" />
</p>

*جادوگر راه‌اندازی تعاملی — نوع پروکسی، پورت‌ها، آدرس اتصال و روش اشتراک‌گذاری را در چند ثانیه انتخاب کنید.*

---

## ✨ ویژگی‌ها

- 🔁 **پروکسی HTTP و SOCKS5** — یکی یا هر دو را همزمان با پورت‌های دلخواه اجرا کنید.
- 🧙 **جادوگر تعاملی** — یک رابط کاربری ترمینال (TUI) ظریف شما را گام‌به‌گام در تنظیمات راهنمایی می‌کند.
- 🛡️ **اتوماسیون فایروال** — به‌طور خودکار قوانین فایروال را در **ویندوز، لینوکس (ufw/firewalld/nftables/iptables) و macOS** ایجاد می‌کند (یا دستورهای آماده کپی می‌دهد).
- 📡 **تشخیص خودکار IP** — تمام آدرس‌های قابل‌دسترس را نمایش می‌دهد تا دستگاه‌ها بتوانند متصل شوند.
- 🔌 **۵ روش اشتراک‌گذاری** — شبکه محلی، وای‌فای، هات‌اسپات (AP)، اتصال USB و Port Forward روتر.
- 💾 **تنظیمات پایدار** — تنظیمات شما در `~/.gnet-cli/config.json` ذخیره می‌شود.
- 🐛 **لاگ دیباگ** — لاگ‌گیری اختیاری در فایل و ردیابی اتصال‌های زنده.
- ⏹️ **خاموشی با وقار** — شنونده‌ها و اتصال‌های فعال را به‌طور تمیز می‌بندد.

---

## 📦 نصب

### ساخت از سورس

```bash
git clone https://github.com/frank-vpl/gnet-cli.git
cd gnet-cli
go build -o gnet .
```

> در ویندوز خروجی `gnet.exe` است. پوشه را به `PATH` اضافه کنید یا مستقیماً اجرایش کنید.

### پیش‌نیازها

- [Go](https://go.dev/) نسخه ۱.۲۷ یا بالاتر

---

## ⬇️ دانلود

آخرین باینری از پیش ساخته‌شده برای پلتفرم خود را دریافت کنید — نیازی به نصب Go نیست.

| پلتفرم | معماری | دانلود |
|----------|-------------|----------|
| 🪟 ویندوز | x64 | [![Download](https://img.shields.io/badge/Windows-x64-0078D6?logo=windows&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-windows-x64.zip) |
| 🪟 ویندوز | x86 | [![Download](https://img.shields.io/badge/Windows-x86-0078D6?logo=windows&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-windows-x86.zip) |
| 🐧 لینوکس | x64 | [![Download](https://img.shields.io/badge/Linux-x64-FCC624?logo=linux&logoColor=black&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-linux-x64.zip) |
| 🐧 لینوکس | x86 | [![Download](https://img.shields.io/badge/Linux-x86-FCC624?logo=linux&logoColor=black&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-linux-x86.zip) |
| 🍎 macOS | x64 (Intel) | [![Download](https://img.shields.io/badge/macOS-x64-000000?logo=apple&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-macos-x64.zip) |
| 🍎 macOS | ARM64 (Apple Silicon) | [![Download](https://img.shields.io/badge/macOS-ARM64-000000?logo=apple&logoColor=white&style=for-the-badge)](https://github.com/frank-vpl/gnet-cli/releases/latest/download/gnet-macos-arm64.zip) |

> 💡 هر فایل ZIP شامل باینری `gnet` (یا `gnet.exe`)، `README.md` و `LICENSE` است.
> آن را از حالت فشرده خارج کنید، پوشه را به `PATH` اضافه کنید و `gnet` را اجرا کنید.

تمام نسخه‌ها: [github.com/frank-vpl/gnet-cli/releases](https://github.com/frank-vpl/gnet-cli/releases)

---

## 🔧 افزودن به PATH (اختیاری)

پس از دانلود یا ساخت باینری `gnet`، پوشه آن را به متغیر `PATH` سیستم خود اضافه کنید تا بتوانید `gnet` را از هر پنجره ترمینالی اجرا کنید.

### 🪟 ویندوز

**گزینه A — PowerShell (کاربر فعلی، پیشنهادی):**

```powershell
# C:\Tools\gnet را با پوشه حاوی gnet.exe جایگزین کنید
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Tools\gnet", "User")
```

ترمینال را ببندید و دوباره باز کنید، سپس تأیید کنید:

```powershell
gnet --help
```

**گزینه B — رابط گرافیکی (System Properties):**

1. کلیدهای `Win + R` را فشار دهید، `sysdm.cpl` را تایپ کنید و Enter بزنید.
2. به **Advanced** → **Environment Variables** بروید.
3. در بخش **User variables**، `Path` را انتخاب کنید → **Edit** → **New**.
4. مسیر پوشه حاوی `gnet.exe` را جایگذاری کنید (مثلاً `C:\Tools\gnet`).
5. روی تمام پنجره‌ها **OK** بزنید و ترمینال را دوباره راه‌اندازی کنید.

**گزینه C — اپلیکیشن تنظیمات ویندوز:**

1. **Start** را باز کنید، عبارت `environment variables` را تایپ کنید و **Edit the system environment variables** (یا **Edit environment variables for your account**) را انتخاب کنید.
2. در پنجره‌ای که باز می‌شود، روی **Environment Variables…** کلیک کنید.
3. در بخش **User variables**، `Path` را انتخاب کنید → **Edit** → **New**.
4. مسیر پوشه حاوی `gnet.exe` را جایگذاری کنید (مثلاً `C:\Tools\gnet`).
5. روی تمام پنجره‌ها **OK** بزنید و ترمینال را دوباره راه‌اندازی کنید.

> نکته: همچنین می‌توانید از مسیر **Settings → System → About → Advanced system settings** به این صفحه برسید.

### 🐧 لینوکس

باینری را به پوشه‌ای که از قبل در `PATH` شماست منتقل کنید:

```bash
sudo mv gnet /usr/local/bin/gnet
sudo chmod +x /usr/local/bin/gnet
```

یا پوشه آن را به‌طور دائم به `PATH` اضافه کنید:

```bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
# یا ~/.zshrc اگر از Zsh استفاده می‌کنید
source ~/.bashrc
```

تأیید:

```bash
gnet --help
```

### 🍎 macOS

```bash
sudo mv gnet /usr/local/bin/gnet
sudo chmod +x /usr/local/bin/gnet
```

اگر آن را از طریق Safari دانلود کرده‌اید و macOS آن را مسدود کرد، یک‌بار اجازه دهید:

```bash
xattr -d com.apple.quarantine /usr/local/bin/gnet
```

تأیید:

```bash
gnet --help
```

---

## 🚀 شروع سریع

جادوگر تعاملی را بدون هیچ آرگومانی اجرا کنید:

```bash
gnet
```

یا سرورهای پروکسی را مستقیماً با تنظیمات ذخیره‌شده/پیش‌فرض راه‌اندازی کنید:

```bash
gnet serve
```

خروجی مشابه زیر را خواهید دید:

```
  GNet CLI proxy running
  HTTP   -> 0.0.0.0:8080
  SOCKS5 -> 0.0.0.0:1080
  LAN IP -> 192.168.1.42
  Press Ctrl+C to stop
```

تنظیمات پروکسی دستگاه خود را به یکی از IPهای نمایش‌داده‌شده هدایت کنید:

| نوع   | مثال             |
|--------|---------------------|
| HTTP   | `192.168.1.42:8080` |
| SOCKS5 | `192.168.1.42:1080` |

---

## 📚 دستورها

| دستور            | توضیح                                         |
|--------------------|-----------------------------------------------------|
| `gnet`             | راه‌اندازی جادوگر تنظیمات تعاملی                 |
| `gnet serve`       | شروع سرورهای پروکسی HTTP/SOCKS5                     |
| `gnet config`      | نمایش تنظیمات ذخیره‌شده                        |
| `gnet interfaces`  | فهرست تمام کارت‌های شبکه و آدرس IP آن‌ها  |
| `gnet methods`     | فهرست تمام روش‌های اشتراک‌گذاری با مراحل راه‌اندازی          |

### پرچم‌های `gnet serve`

| پرچم            | نوع   | مقدار پیش‌فرض | توضیح                    |
|-----------------|--------|---------|--------------------------------|
| `--http-port`   | int    | `8080`  | پورت پروکسی HTTP                |
| `--socks-port`  | int    | `1080`  | پورت پروکسی SOCKS5              |
| `--bind`        | string | `0.0.0.0` | آدرس اتصال                 |
| `--http`        | bool   | `true`  | فعال‌سازی پروکسی HTTP              |
| `--socks5`      | bool   | `true`  | فعال‌سازی پروکسی SOCKS5            |
| `--log-file`    | string | –       | نوشتن لاگ‌ها در فایل           |
| `--debug`       | bool   | `false` | فعال‌سازی لاگ‌گیری سطح دیباگ     |

**مثال‌ها**

```bash
# فقط HTTP روی پورت 3128
gnet serve --http-port 3128 --socks5=false

# فقط SOCKS5، متصل به یک اینترفیس خاص
gnet serve --bind 192.168.1.42 --http=false

# لاگ‌گیری دیباگ در فایل
gnet serve --debug --log-file gnet.log
```

---

## 🔌 روش‌های اشتراک‌گذاری

| روش | نیاز به ادمین | بهترین استفاده |
|--------|----------------|----------|
| 🌐 **شبکه محلی / Ethernet** | ❌ | دستگاه‌های سیمی در همان شبکه |
| 📶 **وای‌فای (همین شبکه)** | ❌ | دستگاه‌های روی همان وای‌فای |
| 📱 **هات‌اسپات موبایل / AP** | ✅ | اتصال مستقیم دستگاه‌ها به وای‌فای سیستم شما |
| 🔌 **اتصال USB / RNDIS** | ✅ | اتصال سیمی USB به تلفن‌ها/تبلت‌ها |
| 🏠 **روتر / Port Forward** | ❌ | دسترسی به دستگاه‌ها در شبکه‌های مختلف |

دستور `gnet methods` را اجرا کنید تا مراحل راه‌اندازی مخصوص هر سیستم‌عامل برای هر روش را ببینید.

---

## 🔥 دستورهای فایروال

ابزار GNet CLI سعی می‌کند فایروال را برای پورت‌های خود **به‌طور خودکار باز** کند (پیش‌فرض `8080` HTTP، `1080` SOCKS5). اگر راه‌اندازی خودکار به دسترسی root/admin نیاز داشته باشد، دستور دقیق برای اجرا را چاپ می‌کند. در زیر مرجع کاملی برای تمام سیستم‌عامل‌ها و فایروال‌های اصلی آورده شده است.

> عبارت `<PORT>` را با پورت تنظیم‌شده خود جایگزین کنید (مثلاً `8080` یا `1080`). برای چند پورت از کاما استفاده کنید، مثلاً `8080,1080`.

### 🪟 ویندوز

```powershell
# افزودن قانون ورودی (به‌عنوان Administrator اجرا کنید)
netsh advfirewall firewall add rule name="GNet CLI" dir=in action=allow protocol=TCP localport=8080,1080 profile=any

# بررسی وجود قانون
netsh advfirewall firewall show rule name="GNet CLI"

# حذف قانون
netsh advfirewall firewall delete rule name="GNet CLI"
```

### 🐧 لینوکس

#### ۱. ufw (اوبونتو، دبیان، مینت، Pop!_OS، کالی)

```bash
sudo ufw allow 8080/tcp
sudo ufw allow 1080/tcp
# یا در یک خط:
sudo ufw allow 8080,1080/tcp

# بررسی وضعیت
sudo ufw status verbose

# حذف قانون
sudo ufw delete allow 8080/tcp
sudo ufw delete allow 1080/tcp
```

#### ۲. firewalld (RHEL، سنتی‌اوس، فدورا، Rocky، Alma، openSUSE)

```bash
# دائمی + بارگذاری مجدد
sudo firewall-cmd --permanent --add-port=8080/tcp --add-port=1080/tcp
sudo firewall-cmd --reload

# بررسی
sudo firewall-cmd --list-ports

# حذف
sudo firewall-cmd --permanent --remove-port=8080/tcp --remove-port=1080/tcp
sudo firewall-cmd --reload
```

#### ۳. nftables (Arch، Alpine، دبیان مدرن، فدورا)

```bash
# ایجاد جدول + زنجیره + قانون اختصاصی
sudo nft add table inet gnet
sudo nft add chain inet gnet input '{ type filter hook input priority 0 ; }'
sudo nft add rule inet gnet input tcp dport {8080,1080} accept

# فهرست قوانین
sudo nft list ruleset

# حذف (حذف کل جدول)
sudo nft delete table inet gnet
```

#### ۴. iptables (لینوکس قدیمی)

```bash
# اجازه دادن به پورت‌ها
sudo iptables -A INPUT -p tcp -m multiport --dports 8080,1080 -j ACCEPT

# ذخیره قانون (پایدار)
sudo netfilter-persistent save          # Debian/Ubuntu
sudo service iptables save              # RHEL/CentOS 6

# فهرست قوانین
sudo iptables -L INPUT -n --line-numbers

# حذف قانون (با شماره خط)
sudo iptables -D INPUT 1
```

#### ۵. OpenWrt (LuCI / UCI)

```bash
# اجازه دادن به پورت‌ها در ناحیه WAN
uci add firewall rule
uci set firewall.@rule[-1].name='GNet CLI'
uci set firewall.@rule[-1].src='wan'
uci set firewall.@rule[-1].dest_port='8080 1080'
uci set firewall.@rule[-1].proto='tcp'
uci set firewall.@rule[-1].target='ACCEPT'
uci commit firewall
/etc/init.d/firewall restart
```

### 🍎 macOS

```bash
# فعال‌سازی فایروال
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate on

# اجازه دادن به باینری GNet
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /path/to/gnet
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp /path/to/gnet

# مشاهده وضعیت فایروال
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --listapps

# حذف قانون اپلیکیشن
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --remove /path/to/gnet
```

**یا از طریق رابط گرافیکی:** System Settings → Network → Firewall → Options → **+** → `gnet` را انتخاب کنید → Allow incoming connections.

### 🧹 BSD / pfSense / FreeBSD (pf)

```pf
# /etc/pf.conf (FreeBSD) یا از طریق رابط گرافیکی pfSense:
pass in on any proto tcp from any to any port { 8080, 1080 }
```

```bash
# بارگذاری مجدد pf
sudo pfctl -f /etc/pf.conf && sudo pfctl -e
```

---

## ⚙️ پیکربندی

تنظیمات در `~/.gnet-cli/config.json` ذخیره می‌شود:

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

**قوانین اعتبارسنجی**

- حداقل یکی از پروکسی‌ها (HTTP یا SOCKS5) باید فعال باشد.
- پورت‌ها باید در بازه `1–65535` باشند.
- پورت‌های HTTP و SOCKS5 نمی‌توانند یکسان باشند.

---

## 🛠️ ساختار پروژه

```
gnet-cli/
├── main.go                    # نقطه ورود
├── go.mod
├── go.sum
└── internal/
    ├── cli/
    │   ├── commands.go        # دستورها و پرچم‌های Cobra
    │   ├── tui.go             # جادوگر تعاملی Bubble Tea
    │   └── logring.go         # بافر حلقه‌ای درون‌حافظه‌ای برای لاگ‌های زنده
    ├── config/
    │   └── config.go          # بارگذاری/ذخیره/اعتبارسنجی تنظیمات
    ├── logger/
    │   └── logger.go          # لاگر ساختاریافته
    ├── network/
    │   ├── network.go         # شمارش اینترفیس‌ها و توابع کمکی IP
    │   └── firewall.go        # مدیریت قوانین فایروال (ویندوز + راهنما)
    ├── proxy/
    │   ├── server.go          # چرخه حیات شنونده و خاموشی با وقار
    │   ├── http.go            # هندلر پروکسی HTTP CONNECT
    │   ├── socks5.go          # هندلر پروکسی SOCKS5
    │   └── relay.go           # رله اتصال دوطرفه
    └── share/
        └── methods.go         # روش‌های اشتراک‌گذاری و مراحل راه‌اندازی مخصوص OS
```

---

## 🧪 توسعه

```bash
# اجرای مستقیم
go run .

# اجرای تست‌ها
go test ./...

# ساخت
go build -o gnet .
```

---

## 🚀 CI/CD

این پروژه از **GitHub Actions** برای تست خودکار، کامپایل متقابل و انتشار استفاده می‌کند.

### گردش کار: [Build and Test](.github/workflows/build.yaml)

| شغل | مبدأ اجرا | توضیح |
|-----|---------|-------------|
| `test` | push / PR | اجرای `go test ./...` روی اوبونتو، ویندوز و macOS |
| `build` | push / PR / tag | کامپایل متقابل ۶ باینری (لینوکس/ویندوز/macOS × x86/x64/arm64) و آپلود ZIPها |
| `release` | تگ `v*` یا دستی | ایجاد یک انتشار GitHub با تمام باینری‌ها + جدول دانلود |

### فرآیند انتشار خودکار

1. یک تگ نسخه را پوش کنید، مثلاً `git tag v1.0.0 && git push origin v1.0.0` (یا گردش کار را به‌طور دستی از تب **Actions** اجرا کنید).
2. شغل `build` باینری `gnet` را برای تمام پلتفرم‌ها کامپایل می‌کند و هر باینری را همراه با `README.md` + `LICENSE` در یک ZIP بسته‌بندی می‌کند.
3. شغل `release` یک انتشار GitHub منتشر می‌کند و هر ZIP را به‌عنوان یک فایل قابل‌دانلود ضمیمه می‌کند.

### وضعیت

- **Build & Test:** [![Build and Test](https://img.shields.io/github/actions/workflow/status/frank-vpl/gnet-cli/build.yaml?logo=githubactions&logoColor=white&style=flat-square)](https://github.com/frank-vpl/gnet-cli/actions/workflows/build.yaml)
- **Latest Release:** [![Latest Release](https://img.shields.io/github/v/release/frank-vpl/gnet-cli?display_name=release&logo=github&style=flat-square)](https://github.com/frank-vpl/gnet-cli/releases/latest)

---

## 📄 مجوز

تحت مجوز MIT منتشر شده است. برای جزئیات به [LICENSE](LICENSE) مراجعه کنید.

---

<div align="center">

ساخته‌شده با 🩷 توسط **حسین پیرا** با استفاده از [Cobra](https://github.com/spf13/cobra)، [Bubble Tea](https://github.com/charmbracelet/bubbletea) و [Lip Gloss](https://github.com/charmbracelet/lipgloss)

</div>
