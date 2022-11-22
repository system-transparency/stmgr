package eval

import (
	"flag"

	"git.glasklar.is/system-transparency/core/stmgr/log"
	"git.glasklar.is/system-transparency/core/stmgr/ospkg"
)

// OspkgCreate takes arguments like os.Args as a string array
// and maps them to their corresponding flags using the std flag
// package. It then calls ospkg.Create after they are parsed.
func OspkgCreate(args []string) error {
	// Create a custom flag set and register flags
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createOut := createCmd.String("out", "", "OS package output path."+
		" Two files will be created: the archive ZIP file and the descriptor JSON file."+
		" A directory or a filename can be passed."+
		" In case of a filename the file extensions will be set properly."+
		" Default name is system-transparency-os-package.")
	createLabel := createCmd.String("label", "", "Short description of the boot configuration."+
		" Defaults to 'System Transparency OS package <kernel>'.")
	createURL := createCmd.String("url", "", "URL of the OS package zip file in case of network boot mode.")
	createKernel := createCmd.String("kernel", "", "Operating system kernel.")
	createInitramfs := createCmd.String("initramfs", "", "Operating system initramfs.")
	createCmdLine := createCmd.String("cmdline", "", "Kernel command line.")
	createLogLevel := createCmd.String("loglevel", "", "Set loglevel to any of debug, info, warn, error (default) and panic.")

	// Parse which flags are provided to the function
	if err := createCmd.Parse(args); err != nil {
		return err
	}

	// Adjust loglevel
	setLoglevel(*createLogLevel)

	// Print the successfully parsed flags in debug level
	createCmd.Visit(func(f *flag.Flag) {
		log.Debugf("Registered flag %q", f)
	})

	// Call function with parsed flags
	return ospkg.Create(
		&ospkg.CreateArgs{
			OutPath:   *createOut,
			Label:     *createLabel,
			URL:       *createURL,
			Kernel:    *createKernel,
			Initramfs: *createInitramfs,
			Cmdline:   *createCmdLine,
		},
	)
}

// OspkgSign takes arguments like os.Args as a string array
// and maps them to their corresponding flags using the std flag
// package. It then calls ospkg.Sign after they are parsed.
func OspkgSign(args []string) error {
	// Create a custom flag set and register flags
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signKey := signCmd.String("key", "", "Private key for signing.")
	signCert := signCmd.String("cert", "", "Certificate corresponding to the private key.")
	signOSPKG := signCmd.String("ospkg", "", "OS package archive or descriptor file. Both need to be present.")
	signLogLevel := signCmd.String("loglevel", "", "Set loglevel to any of debug, info, warn, error (default) and panic.")

	// Parse which flags are provided to the function
	if err := signCmd.Parse(args); err != nil {
		return err
	}

	// Adjust loglevel
	setLoglevel(*signLogLevel)

	// Print the successfully parsed flags in debug level
	signCmd.Visit(func(f *flag.Flag) {
		log.Debugf("Registered flag %q", f)
	})

	// Call function with parsed flags
	return ospkg.Sign(*signKey, *signCert, *signOSPKG)
}
