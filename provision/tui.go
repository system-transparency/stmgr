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

var (
	app              = tview.NewApplication() //nolint:gochecknoglobals
	mainForm         = tview.NewForm()        //nolint:gochecknoglobals
	interfacesForm   = tview.NewForm()        //nolint:gochecknoglobals
	customForm       = tview.NewForm()        //nolint:gochecknoglobals
	pages            = tview.NewPages()       //nolint:gochecknoglobals
	addrMode         = tview.NewDropDown()    //nolint:gochecknoglobals
	versionField     = tview.NewInputField()  //nolint:gochecknoglobals
	hostIPField      = tview.NewInputField()  //nolint:gochecknoglobals
	gatewayIPField   = tview.NewInputField()  //nolint:gochecknoglobals
	dnsField         = tview.NewInputField()  //nolint:gochecknoglobals
	interfaceField   = tview.NewInputField()  //nolint:gochecknoglobals
	provURLField     = tview.NewInputField()  //nolint:gochecknoglobals
	idField          = tview.NewInputField()  //nolint:gochecknoglobals
	authField        = tview.NewInputField()  //nolint:gochecknoglobals
	customLabelField = tview.NewInputField()  //nolint:gochecknoglobals

	extension = make(map[string]string) //nolint:gochecknoglobals
)

func runInteractive(efi bool) error { //nolint:funlen,cyclop
	cfg := &HostCfgSimplified{}

	versionField.
		SetLabel("Version").
		SetDoneFunc(func(key tcell.Key) {
			if v, ok := evalVersion(versionField.GetText()); ok {
				cfg.Version = v
			} else {
				versionField.SetText("1")
			}
		})

	addrMode.
		SetLabel("Network Mode").
		AddOption("dhcp", func() {
			cfg.IPAddrMode = "dhcp"
		}).
		AddOption("static", func() {
			cfg.IPAddrMode = "static"
		})

	hostIPField.
		SetLabel("Host IP").
		SetDoneFunc(func(key tcell.Key) {
			_, network, ok := evalCIDR(hostIPField.GetText())
			if ok {
				cfg.HostIP = hostIPField.GetText()
				gatewayIPField.SetText(guessGateway(network))
			}
		})

	gatewayIPField.
		SetLabel("Gateway IP").
		SetDoneFunc(func(key tcell.Key) {
			if evalIP(gatewayIPField.GetText()) {
				cfg.DefaultGateway = gatewayIPField.GetText()
			}
		})

	dnsField.
		SetLabel("DNS IP").
		SetDoneFunc(func(key tcell.Key) {
			if evalIP(dnsField.GetText()) {
				cfg.DNSServer = dnsField.GetText()
			}
		})

	interfaceField.
		SetLabel("Network Interfaces").
		SetDoneFunc(func(key tcell.Key) {
			if evalMAC(interfaceField.GetText()) {
				cfg.NetworkInterface = interfaceField.GetText()
			}
		})

	provURLField.
		SetLabel("Provisioning URLs").
		SetDoneFunc(func(key tcell.Key) {
			if evalURLs(provURLField.GetText()) {
				cfg.ProvisioningURLs = strings.Split(provURLField.GetText(), " ")
			}
		})

	idField.
		SetLabel("ID").
		SetDoneFunc(func(key tcell.Key) {
			if idField.GetText() == "" {
				idField.SetText(getRandomHex())
				cfg.ID = idField.GetText()
			}
			if evalRand(idField.GetText()) {
				cfg.ID = idField.GetText()
			}
		})

	authField.
		SetLabel("Authentication").
		SetDoneFunc(func(key tcell.Key) {
			if authField.GetText() == "" {
				authField.SetText(getRandomHex())
				cfg.Auth = authField.GetText()
			}
			if evalRand(authField.GetText()) {
				cfg.Auth = authField.GetText()
			}
		})

	mainForm.
		AddFormItem(versionField).
		AddFormItem(addrMode).
		AddFormItem(hostIPField).
		AddFormItem(gatewayIPField).
		AddFormItem(dnsField).
		AddFormItem(interfaceField).
		AddFormItem(provURLField).
		AddFormItem(idField).
		AddFormItem(authField).
		AddButton("Save", func() {
			cfg.Timestamp = time.Now().Unix()
			cfg = appendCustomData(cfg)
			if err := MarshalCfg(cfg, efi); err != nil {
				app.QueueEvent(tcell.NewEventError(err))
			}
			app.Stop()
		}).
		AddButton("Interfaces", func() {
			pages.SwitchToPage("interfaces")
		}).
		AddButton("Add Field", func() {
			pages.SwitchToPage("custom")
		}).
		AddButton("Exit", func() { app.Stop() }).
		SetBorder(true).
		SetTitle("CFGTOOL").
		SetTitleAlign(tview.AlignLeft)

	if ifaces, err := net.Interfaces(); err == nil {
		for _, iface := range ifaces {
			mac := iface.HardwareAddr.String()
			if !strings.Contains(iface.Flags.String(), net.FlagLoopback.String()) {
				interfacesForm.AddCheckbox(iface.Name+" ["+mac+"]", false, func(checked bool) {
					toggleMAC(checked, mac)
				})
			}
		}
	}

	interfacesForm.
		AddButton("Back", func() {
			pages.SwitchToPage("main")
		}).
		AddButton("Exit", func() { app.Stop() })

	customLabelField.SetLabel("Custom field key")

	customForm.
		AddFormItem(customLabelField).
		AddButton("Add", func() {
			addCustomField()
			pages.SwitchToPage("main")
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("main")
		}).
		AddButton("Exit", func() { app.Stop() })

	pages.
		AddPage("main", mainForm, true, true).
		AddPage("interfaces", interfacesForm, true, false).
		AddPage("custom", customForm, true, false)

	return app.SetRoot(pages, true).SetFocus(pages).Run()
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

func evalMAC(mac string) bool {
	if _, err := net.ParseMAC(mac); err != nil {
		return false
	}

	return true
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

func toggleMAC(state bool, mac string) {
	text := interfaceField.GetText()
	if state {
		if text == "" {
			interfaceField.SetText(mac)

			return
		}

		interfaceField.SetText(text + " " + mac)
	} else {
		interfaceField.SetText(text + " ")
		interfaceField.SetText(strings.ReplaceAll(text, mac+" ", ""))
		text = interfaceField.GetText()
		interfaceField.SetText(strings.TrimSpace(text))
	}
}

func addCustomField() {
	customField := tview.NewInputField()
	customField.
		SetLabel(customLabelField.GetText()).
		SetDoneFunc(func(key tcell.Key) {
			extension[customField.GetLabel()] = customField.GetText()
		})
	customLabelField.SetText("")
	mainForm.AddFormItem(customField)
}

func appendCustomData(cfg *HostCfgSimplified) *HostCfgSimplified {
	if mainForm.GetFormItemCount() <= minimumElementCount {
		return cfg
	}

	cfg.Custom = extension

	return cfg
}
