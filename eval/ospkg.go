package eval

import (
	"flag"

	"github.com/system-transparency/stmgr/log"
	"github.com/system-transparency/stmgr/ospkg"
)

func OspkgCreate(args []string) error {
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

	if err := createCmd.Parse(args); err != nil {
		return err
	}

	setLoglevel(*createLogLevel)

	createCmd.Visit(func(f *flag.Flag) {
		log.Debugf("Registered flag %q", f)
	})

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

func OspkgSign(args []string) error {
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
	signKey := signCmd.String("key", "", "Private key for signing.")
	signCert := signCmd.String("cert", "", "Certificate corresponding to the private key.")
	signOSPKG := signCmd.String("ospkg", "", "OS package archive or descriptor file. Both need to be present.")
	signLogLevel := signCmd.String("loglevel", "", "Set loglevel to any of debug, info, warn, error (default) and panic.")

	if err := signCmd.Parse(args); err != nil {
		return err
	}

	setLoglevel(*signLogLevel)

	signCmd.Visit(func(f *flag.Flag) {
		log.Debugf("Registered flag %q", f)
	})

	return ospkg.Sign(*signKey, *signCert, *signOSPKG)
}

func OspkgShow(args []string) error {
	log.Print("Not implemented yet!")

	return nil
}
