package eval

import (
	"errors"
	"flag"
	"os"

	"system-transparency.org/stmgr/hostconfig"
)

// HostConfigCheck takes arguments like os.Args as a string array
// and maps them to their corresponding flags using the std flag
// package. It then calls trustpolicy.Create after they are parsed.
func HostConfigCheck(args []string) error {
	createCmd := flag.NewFlagSet("check", flag.ExitOnError)

	if err := createCmd.Parse(args); err != nil {
		return err
	}

	var json string

	switch createCmd.NArg() {
	case 0:
		return errors.New("missing argument, provide input json data")
	case 1:
		json = createCmd.Arg(0)
	default:
		return errors.New("only one argument allowed")
	}

	return hostconfig.Check(json, os.Stdout)
}
