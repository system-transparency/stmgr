package trustpolicy

import (
	"encoding/json"
	"io"

	"system-transparency.org/stboot/trust"
)

func Check(in string, out io.Writer) error {
	var policy trust.Policy

	err := policy.UnmarshalJSON([]byte(in))
	if err != nil {
		return err
	}

	return output(policy, out)
}

func output(tp trust.Policy, out io.Writer) error {
	pretty, err := json.MarshalIndent(tp, "", "  ")
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
