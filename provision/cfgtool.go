package provision

import (
	"strings"
	"time"
)

func Cfgtool(efi bool, version int, addrMode, hostIP, gateway, dns, interfaces, urls, id, auth string) error {
	if isDefined(addrMode, hostIP, gateway, dns, interfaces, urls, id, auth) {
		cfg := &HostCfgSimplified{
			Version:          version,
			IPAddrMode:       addrMode,
			HostIP:           hostIP,
			DefaultGateway:   gateway,
			DNSServer:        dns,
			NetworkInterface: interfaces,
			ProvisioningURLs: strings.Split(urls, " "),
			ID:               id,
			Auth:             auth,
			Timestamp:        time.Now().Unix(),
		}

		return MarshalCfg(cfg, efi)
	}

	return runInteractive(efi)
}

func isDefined(s ...string) bool {
	return strings.Join(s, "") != ""
}
