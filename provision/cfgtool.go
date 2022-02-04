package provision

import (
	"strings"
	"time"
)

// Cfgtool is used to provision a stboot node by creating a
// host configuration JSON file. If any value of the host
// configuration was set to anything but an empty string, then
// the tool will create a file with the info provided and,
// if efi is set to true, write into the efivars or to disk.
// If none of the values in the host configuration are set, it
// will launch in interactive mode and will allow a user to set
// the values using a terminal UI.
func Cfgtool(efi bool, cfg *HostCfgSimplified) error {
	if isDefined(
		cfg.IPAddrMode,
		cfg.HostIP,
		cfg.DefaultGateway,
		cfg.DNSServer,
		cfg.NetworkInterface,
		cfg.ID,
		cfg.Auth,
	) {
		cfg.Timestamp = time.Now().Unix()

		return MarshalCfg(cfg, efi)
	}

	return runInteractive(efi)
}

func isDefined(s ...string) bool {
	return len(strings.Join(s, "")) != 0
}
