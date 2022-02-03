package eval

import "testing"

func TestSetLoglevel(t *testing.T) {
	for _, level := range []string{
		"",
		"debug",
		"info",
		"warn",
		"panic",
	} {
		setLoglevel(level)
	}
}
