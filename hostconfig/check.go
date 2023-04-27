package hostconfig

import (
	"encoding/json"
	"io"

	"system-transparency.org/stboot/host"
)

func Check(in string, out io.Writer) error {
	var cfg host.Config

	err := cfg.UnmarshalJSON([]byte(in))
	if err != nil {
		return err
	}

	return output(cfg, out)
}

func output(cfg host.Config, out io.Writer) error {
	pretty, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	pretty = append(pretty, '\n')

	_, err = out.Write(pretty)
	if err != nil {
		return err
	}

	return nil
}
