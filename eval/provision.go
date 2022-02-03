package eval

import (
	"flag"
	"strings"

	"github.com/system-transparency/stmgr/log"
	"github.com/system-transparency/stmgr/provision"
)

func ProvisionHostconfig(args []string) error {
	hostconfigCmd := flag.NewFlagSet("provision", flag.ExitOnError)
	hostconfigEfi := hostconfigCmd.Bool("efi", false, "Store host_configuration.json in the efivarfs.")
	hostconfigVersion := hostconfigCmd.Int("version", 1, "Hostconfig version.")
	hostconfigAddrMode := hostconfigCmd.String("addrMode", "", "Hostconfig network_mode.")
	hostconfigHostIP := hostconfigCmd.String("hostIP", "", "Hostconfig host_ip.")
	hostconfigGateway := hostconfigCmd.String("gateway", "", "Hostconfig gateway.")
	hostconfigDNS := hostconfigCmd.String("dns", "", "Hostconfig dns.")
	hostconfigInterface := hostconfigCmd.String("interface", "", "Hostconfig network_interface.")
	hostconfigURLs := hostconfigCmd.String("urls", "", "Hostconfig provisioning_urls.")
	hostconfigID := hostconfigCmd.String("id", "", "Hostconfig identity.")
	hostconfigAuth := hostconfigCmd.String("auth", "", "Hostconfig authentication.")
	hostconfigLogLevel := hostconfigCmd.String("loglevel", "", "Set loglevel to any of debug, info, warn, error (default) and panic.")

	if err := hostconfigCmd.Parse(args); err != nil {
		return err
	}

	setLoglevel(*hostconfigLogLevel)

	hostconfigCmd.Visit(func(f *flag.Flag) {
		log.Debugf("Registered flag %q", f)
	})

	return provision.Cfgtool(
		*hostconfigEfi,
		&provision.HostCfgSimplified{
			Version:          *hostconfigVersion,
			IPAddrMode:       *hostconfigAddrMode,
			HostIP:           *hostconfigHostIP,
			DefaultGateway:   *hostconfigGateway,
			DNSServer:        *hostconfigDNS,
			NetworkInterface: *hostconfigInterface,
			ProvisioningURLs: strings.Split(*hostconfigURLs, " "),
			ID:               *hostconfigID,
			Auth:             *hostconfigAuth,
		},
	)
}
