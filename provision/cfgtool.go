package provision

import (
	"strings"
	"time"
)

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
