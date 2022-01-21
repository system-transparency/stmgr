package provision

import (
	"encoding/json"
	"io/fs"
	"os"

	guid "github.com/google/uuid"
	"github.com/system-transparency/efivar/efivarfs"
)

const defaultFilePerm fs.FileMode = 0o600

type HostCfgSimplified struct {
	Version          int               `json:"version"`
	IPAddrMode       string            `json:"network_mode"`
	HostIP           string            `json:"host_ip"`
	DefaultGateway   string            `json:"gateway"`
	DNSServer        string            `json:"dns"`
	NetworkInterface string            `json:"network_interface"`
	ProvisioningURLs []string          `json:"provisioning_urls"`
	ID               string            `json:"identity"`
	Auth             string            `json:"authentication"`
	Timestamp        int64             `json:"timestamp"`
	Custom           map[string]string `json:"custom,omitempty"`
}

func MarshalCfg(cfg *HostCfgSimplified, efi bool) error {
	j, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if efi {
		g, err := guid.Parse("f401f2c1-b005-4be0-8cee-f2e5945bcbe7")
		if err != nil {
			return err
		}

		attrs := efivarfs.AttributeBootserviceAccess | efivarfs.AttributeRuntimeAccess | efivarfs.AttributeNonVolatile

		return efivarfs.WriteVariable("STHostConfig", &g, attrs, j)
	}

	return os.WriteFile("host_configuration.json", j, defaultFilePerm)
}
