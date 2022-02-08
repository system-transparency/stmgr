package provision

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// The nature of tview does not really allow me to write this
// code in a linter friendly way. Maybe I'll get to it when we
// really need to enforce those rules but for now it should
// be sufficient to just mark problematic statements with nolint

const (
	maxIDAndAuthLength      = 64
	maxIDAndAuthLengthBytes = maxIDAndAuthLength / 2
	minimumElementCount     = 9
)

type tui struct {
	app              *tview.Application
	mainForm         *tview.Form
	customForm       *tview.Form
	pages            *tview.Pages
	addrModeMenu     *tview.DropDown
	interfaceMenu    *tview.DropDown
	versionField     *tview.InputField
	hostIPField      *tview.InputField
	gatewayIPField   *tview.InputField
	dnsField         *tview.InputField
	provURLField     *tview.InputField
	idField          *tview.InputField
	authField        *tview.InputField
	customLabelField *tview.InputField
	extension        map[string]string
}

// This function constructs the and manages the terminal UI application.
// It uses tview for that, as it abstracts tcell functions and makes building
// simple UIs easier.
func (t *tui) runInteractive(efi bool) error {
	cfg := &HostCfgSimplified{}

	t.setVersionField(cfg)
	t.setAddrModeMenu(cfg)
	t.setHostIPField(cfg)
	t.setGatewayIPField(cfg)
	t.setDNSField(cfg)
	t.setInterfaceMenu(cfg)
	t.setProvURLField(cfg)
	t.setIDField(cfg)
	t.setAuthField(cfg)
	t.setMainForm(cfg, efi)
	t.setCustomForm()
	t.setPages()

	return t.app.SetRoot(t.pages, true).SetFocus(t.pages).Run()
}

func (t *tui) setVersionField(cfg *HostCfgSimplified) {
	t.versionField.
		SetLabel("Version").
		SetDoneFunc(func(key tcell.Key) {
			if v, ok := evalVersion(t.versionField.GetText()); ok {
				cfg.Version = v
			} else {
				t.versionField.SetText("1")
			}
		})
}

func (t *tui) setAddrModeMenu(cfg *HostCfgSimplified) {
	t.addrModeMenu.
		SetLabel("Network Mode").
		AddOption("dhcp", func() {
			cfg.IPAddrMode = "dhcp"
		}).
		AddOption("static", func() {
			cfg.IPAddrMode = "static"
		})
}

func (t *tui) setHostIPField(cfg *HostCfgSimplified) {
	t.hostIPField.
		SetLabel("Host IP").
		SetDoneFunc(func(key tcell.Key) {
			_, network, ok := evalCIDR(t.hostIPField.GetText())
			if ok {
				cfg.HostIP = t.hostIPField.GetText()
				t.gatewayIPField.SetText(guessGateway(network))
			}
		})
}

func (t *tui) setGatewayIPField(cfg *HostCfgSimplified) {
	t.gatewayIPField.
		SetLabel("Gateway IP").
		SetDoneFunc(func(key tcell.Key) {
			if evalIP(t.gatewayIPField.GetText()) {
				cfg.DefaultGateway = t.gatewayIPField.GetText()
			}
		})
}

func (t *tui) setDNSField(cfg *HostCfgSimplified) {
	t.dnsField.
		SetLabel("DNS IP").
		SetDoneFunc(func(key tcell.Key) {
			if evalIP(t.dnsField.GetText()) {
				cfg.DNSServer = t.dnsField.GetText()
			}
		})
}

func (t *tui) setInterfaceMenu(cfg *HostCfgSimplified) {
	t.interfaceMenu.SetLabel("Network Interfaces")

	if ifaces, err := net.Interfaces(); err != nil {
		t.interfaceMenu.
			AddOption("NONE", func() {
				cfg.NetworkInterface = ""
			})
	} else {
		for _, iface := range ifaces {
			if !strings.Contains(iface.Flags.String(), net.FlagLoopback.String()) {
				mac := iface.HardwareAddr.String()
				t.interfaceMenu.
					AddOption(mac, func() {
						cfg.NetworkInterface = mac
					})
			}
		}
	}
}

func (t *tui) setProvURLField(cfg *HostCfgSimplified) {
	t.provURLField.
		SetLabel("Provisioning URLs").
		SetDoneFunc(func(key tcell.Key) {
			if evalURLs(t.provURLField.GetText()) {
				cfg.ProvisioningURLs = strings.Split(t.provURLField.GetText(), " ")
			}
		})
}

func (t *tui) setIDField(cfg *HostCfgSimplified) {
	t.idField.
		SetLabel("ID").
		SetDoneFunc(func(key tcell.Key) {
			if t.idField.GetText() == "" {
				t.idField.SetText(getRandomHex())
				cfg.ID = t.idField.GetText()
			}
			if evalRand(t.idField.GetText()) {
				cfg.ID = t.idField.GetText()
			}
		})
}

func (t *tui) setAuthField(cfg *HostCfgSimplified) {
	t.authField.
		SetLabel("Authentication").
		SetDoneFunc(func(key tcell.Key) {
			if t.authField.GetText() == "" {
				t.authField.SetText(getRandomHex())
				cfg.Auth = t.authField.GetText()
			}
			if evalRand(t.authField.GetText()) {
				cfg.Auth = t.authField.GetText()
			}
		})
}

func (t *tui) setMainForm(cfg *HostCfgSimplified, efi bool) {
	t.mainForm.
		AddFormItem(t.versionField).
		AddFormItem(t.addrModeMenu).
		AddFormItem(t.hostIPField).
		AddFormItem(t.gatewayIPField).
		AddFormItem(t.dnsField).
		AddFormItem(t.interfaceMenu).
		AddFormItem(t.provURLField).
		AddFormItem(t.idField).
		AddFormItem(t.authField).
		AddButton("Save", func() {
			cfg.Timestamp = time.Now().Unix()
			cfg = t.appendCustomData(cfg)
			if err := MarshalCfg(cfg, efi); err != nil {
				t.app.QueueEvent(tcell.NewEventError(err))
			}
			t.app.Stop()
		}).
		AddButton("Add Field", func() {
			t.pages.SwitchToPage("custom")
		}).
		AddButton("Exit", func() { t.app.Stop() }).
		SetBorder(true).
		SetTitle("CFGTOOL").
		SetTitleAlign(tview.AlignLeft)
}

func (t *tui) setCustomForm() {
	t.customLabelField.SetLabel("Custom field key")

	t.customForm.
		AddFormItem(t.customLabelField).
		AddButton("Add", func() {
			t.addCustomField()
			t.pages.SwitchToPage("main")
		}).
		AddButton("Back", func() {
			t.pages.SwitchToPage("main")
		}).
		AddButton("Exit", func() { t.app.Stop() })
}

func (t *tui) addCustomField() {
	customField := tview.NewInputField()
	customField.
		SetLabel(t.customLabelField.GetText()).
		SetDoneFunc(func(key tcell.Key) {
			t.extension[customField.GetLabel()] = customField.GetText()
		})
	t.customLabelField.SetText("")
	t.mainForm.AddFormItem(customField)
}

func (t *tui) appendCustomData(cfg *HostCfgSimplified) *HostCfgSimplified {
	if t.mainForm.GetFormItemCount() <= minimumElementCount {
		return cfg
	}

	cfg.Custom = t.extension

	return cfg
}

func (t *tui) setPages() {
	t.pages.
		AddPage("main", t.mainForm, true, true).
		AddPage("custom", t.customForm, true, false)
}

func evalVersion(version string) (int, bool) {
	v, err := strconv.Atoi(version)
	if err != nil {
		return 1, false
	}

	return v, true
}

func evalIP(ip string) bool {
	if result := net.ParseIP(ip); result != nil {
		return true
	}

	return false
}

func evalCIDR(cidr string) (string, string, bool) {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", false
	}

	return ip.String(), network.String(), true
}

func evalURLs(urls string) bool {
	for _, url := range strings.Split(urls, " ") {
		if !strings.HasPrefix(url, "http://") || !strings.HasPrefix(url, "https://") {
			return false
		}
	}

	return true
}

func evalRand(s string) bool {
	return len(s) <= maxIDAndAuthLength
}

func getRandomHex() string {
	b := make([]byte, maxIDAndAuthLengthBytes)
	if _, err := rand.Reader.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func guessGateway(s string) string {
	segments := strings.Split(strings.TrimRight(s, "/"), ".")
	i, _ := strconv.Atoi(segments[3])
	i++
	segments[3] = strconv.Itoa(i)

	return strings.Join(segments, ".")
}
