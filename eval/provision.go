package eval

import (
	"flag"
	"strings"

	"system-transparency.org/stmgr/log"
	"system-transparency.org/stmgr/provision"
)

// ProvisionHostconfig takes arguments like os.Args as a string array
// and maps them to their corresponding flags using the std flag
// package. It then calls provision.Cfgtool after they are parsed.
func ProvisionHostconfig(args []string) error {
	// Create a custom flag set and register flags
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

	hostconfigBondingMode := hostconfigCmd.String("bondingmode", "", "Set default bonding mode (optional)")
	hostconfigBondName := hostconfigCmd.String("bondname", "", "Set bonding interface name (optional)")
	hostconfigNetwokInterfaces := hostconfigCmd.String("network-interfaces", "", "Space separated list of network interfaces (optional, requires with bonding)")

	// Parse which flags are provided to the function
	if err := hostconfigCmd.Parse(args); err != nil {
		return err
	}

	// Adjust loglevel
	setLoglevel(*hostconfigLogLevel)

	// Print the successfully parsed flags in debug level
	hostconfigCmd.Visit(func(f *flag.Flag) {
		log.Debugf("Registered flag %q", f)
	})

	// Call function with parsed flags
	return provision.Cfgtool(
		*hostconfigEfi,
		&provision.HostCfgSimplified{
			Version:           *hostconfigVersion,
			IPAddrMode:        hostconfigAddrMode,
			HostIP:            hostconfigHostIP,
			DefaultGateway:    hostconfigGateway,
			DNSServer:         hostconfigDNS,
			NetworkInterface:  hostconfigInterface,
			ProvisioningURLs:  strings.Split(*hostconfigURLs, " "),
			ID:                hostconfigID,
			Auth:              hostconfigAuth,
			NetworkInterfaces: strings.Split(*hostconfigNetwokInterfaces, " "),
			BondingMode:       hostconfigBondingMode,
			BondName:          hostconfigBondName,
		},
	)
}
