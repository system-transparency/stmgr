package eval

import (
	"flag"

	"system-transparency.org/stboot/stlog"
	"system-transparency.org/stmgr/mkosi"
)

// MkosiBuild takes arguments like os.Args as a string array
// and maps them to their corresponding flags using the std flag
// package. It then calls ospkg.Create after they are parsed.
func MkosiBuild(args []string) error {
	// Create a custom flag set and register flags
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildRootPassword := buildCmd.String("root-password", "", "Set root password.")
	buildSSH := buildCmd.Bool("ssh", false, "Enable SSH.")
	buildHostname := buildCmd.String("hostname", "", "Set hostname.")
	buildCmdline := buildCmd.String("kernel-command-line", "", "Set kernel command line.")
	buildLocale := buildCmd.String("locale", "", "Set locale.")
	buildKeymap := buildCmd.String("keymap", "", "Set keymap.")
	buildTimeZone := buildCmd.String("timezone", "", "Set timezone.")
	buildExtraPackages := buildCmd.String("p", "", "Extra packages to install.")
	buildLogLevel := buildCmd.String("loglevel", "", "Set loglevel to any of debug, info, warn, error (default) and panic.")

	// Parse which flags are provided to the function
	if err := buildCmd.Parse(args); err != nil {
		return err
	}

	// Adjust loglevel
	setLoglevel(*buildLogLevel)

	// Print the successfully parsed flags in debug level
	buildCmd.Visit(func(f *flag.Flag) {
		stlog.Debug("Registered flag %q", f)
	})

	// Call function with parsed flags
	return mkosi.Build(
		&mkosi.BuildArgs{
			RootPassword:  *buildRootPassword,
			SSH:           *buildSSH,
			Hostname:      *buildHostname,
			Cmdline:       *buildCmdline,
			Locale:        *buildLocale,
			Keymap:        *buildKeymap,
			TimeZone:      *buildTimeZone,
			ExtraPackages: *buildExtraPackages,
		},
	)
}

// MkosiSign takes arguments like os.Args as a string array
// and maps them to their corresponding flags using the std flag
// package. It then calls ospkg.Sign after they are parsed.
func MkosiSign(args []string) error {
	// Create a custom flag set and register flags
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signKey := signCmd.String("key", "", "Private key for signing.")
	signCert := signCmd.String("cert", "", "Certificate corresponding to the private key.")
	signUKI := signCmd.String("uki", "", "OS package archive or descriptor file. Both need to be present.")
	signLogLevel := signCmd.String("loglevel", "", "Set loglevel to any of debug, info, warn, error (default) and panic.")

	// Parse which flags are provided to the function
	if err := signCmd.Parse(args); err != nil {
		return err
	}

	// Adjust loglevel
	setLoglevel(*signLogLevel)

	// Print the successfully parsed flags in debug level
	signCmd.Visit(func(f *flag.Flag) {
		stlog.Debug("Registered flag %q", f)
	})

	// Call function with parsed flags
	return mkosi.Sign(*signKey, *signCert, *signUKI)
}

func MkosiVerify(args []string) error {
	// Create a custom flag set and register flags
	verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)
	verifyCert := verifyCmd.String("cert", "", "Certificate corresponding to the private key.")
	verifyUKI := verifyCmd.String("uki", "", "OS package archive or descriptor file. Both need to be present.")
	verifyLogLevel := verifyCmd.String("loglevel", "", "Set loglevel to any of debug, info, warn, error (default) and panic.")

	// Parse which flags are provided to the function
	if err := verifyCmd.Parse(args); err != nil {
		return err
	}

	// Adjust loglevel
	setLoglevel(*verifyLogLevel)

	// Print the successfully parsed flags in debug level
	verifyCmd.Visit(func(f *flag.Flag) {
		stlog.Debug("Registered flag %q", f)
	})

	// Call function with parsed flags
	return mkosi.Verify(*verifyCert, *verifyUKI)
}
