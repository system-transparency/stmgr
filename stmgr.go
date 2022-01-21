package main

import (
	"flag"
	"os"

	"github.com/system-transparency/stmgr/keygen"
	"github.com/system-transparency/stmgr/logging"
	"github.com/system-transparency/stmgr/ospkg"
	"github.com/system-transparency/stmgr/provision"
)

const (
	_ = iota
	commandCallPosition
	subcommandCallPosition
	flagsCallPosition
)

const (
	usage = `Usage: stmgr <COMMAND> <SUBCOMMAND> [flags...]
COMMANDS:
	ospkg:
		Set of commands related to OS packages. This includes
		creating, signing and analyzing them.

	provision:
		Set of commands to provision a node for system-transparency
		usage, like creating and writing a host configuration.

	keygen:
		Commands to generate different keys and certificates for
		system-transparency.

	build:
		Not yet implemented!

Use 'stmgr <COMMAND> -help' for more info.
`

	ospkgUsage = `SUBCOMMANDS:
	create:
		Create an OS package from the provided operating
		system files.

	sign:
		Sign the provided OS package with your private key.

Use 'stmgr ospkg <SUBCOMMAND> -help' for more info.
`

	provisionUsage = `SUBCOMMANDS:
	hostconfig:
		Allows creating host configurations by spawning a TUI in
		which the user can input values into that are converted
		into a host_configuration.json file.

Use 'stmgr provision <SUBCOMMAND> -help' for more info.
`

	keygenUsage = `SUBCOMMANDS:
	certificate:
		Generate certificates for signing OS packages
		using ED25519 keys.

Use 'stmgr keygen <SUBCOMMAND> -help' for more info.
`
)

func main() {
	log := logging.NewLogger(logging.ErrorLevel)
	if err := run(os.Args, log); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run(args []string, log *logging.Logger) error {
	// Display helptext if no arguments are given
	if len(args) < flagsCallPosition {
		log.Print(usage)

		return nil
	}

	// Check which command is requested or display usage
	switch args[commandCallPosition] {
	case "ospkg":
		return ospkgArg(args, log)
	case "provision":
		return provisionArg(args, log)
	case "keygen":
		return keygenArg(args, log)
	default:
		// Display usage on unknown command
		log.Print(usage)

		return nil
	}
}

// Check for ospkg subcommands.
func ospkgArg(args []string, log *logging.Logger) error {
	switch args[subcommandCallPosition] {
	case "create":
		// Create tool and flags
		createCmd := flag.NewFlagSet("createOSPKG", flag.ExitOnError)
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

		if err := createCmd.Parse(args[flagsCallPosition:]); err != nil {
			return err
		}

		return ospkg.Create(
			&ospkg.Args{
				OutPath:   *createOut,
				Label:     *createLabel,
				URL:       *createURL,
				Kernel:    *createKernel,
				Initramfs: *createInitramfs,
				Cmdline:   *createCmdLine,
			},
		)

	case "sign":
		// Sign tool and flags
		signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
		signKey := signCmd.String("key", "", "Private key for signing.")
		signCert := signCmd.String("cert", "", "Certificate corresponding to the private key.")
		signOSPKG := signCmd.String("ospkg", "", "OS package archive or descriptor file. Both need to be present.")

		if err := signCmd.Parse(args[flagsCallPosition:]); err != nil {
			return err
		}

		return ospkg.Sign(*signKey, *signCert, *signOSPKG)

	case "show":
		// Show tool and flags
		log.Print("Not implemented yet!")

		return nil

	default:
		// Display usage on unknown subcommand
		log.Print(ospkgUsage)

		return nil
	}
}

// Check for provision subcommands.
func provisionArg(args []string, log *logging.Logger) error {
	switch args[subcommandCallPosition] {
	case "hostconfig":
		// Host configuration tool and flags
		hostconfigCmd := flag.NewFlagSet("provision", flag.ExitOnError)
		hostconfigEfi := hostconfigCmd.Bool("efi", false, "Store host_configuration.json in the efivarfs.")
		hostconfigVersion := hostconfigCmd.Int("version", 1, "Hostconfig version.")
		hostconfigAddrMode := hostconfigCmd.String("addrMode", "", "Hostconfig network_mode.")
		hostconfigHostIP := hostconfigCmd.String("hostIP", "", "Hostconfig host_ip.")
		hostconfigGateway := hostconfigCmd.String("gateway", "", "Hostconfig gateway.")
		hostconfigDNS := hostconfigCmd.String("dns", "", "Hostconfig dns.")
		hostconfigInterface := hostconfigCmd.String("interface", "", "Hostconfig network_interface.")
		hostconfigURLs := hostconfigCmd.String("urls", "", "Hostconfig provisioning_urls.")
		hostconfigID := hostconfigCmd.String("id", "", "Hostconfig identity.")
		hostconfigAuth := hostconfigCmd.String("auth", "", "Hostconfig authentication.")

		if err := hostconfigCmd.Parse(args[flagsCallPosition:]); err != nil {
			return err
		}

		return provision.Cfgtool(
			*hostconfigEfi,
			*hostconfigVersion,
			*hostconfigAddrMode,
			*hostconfigHostIP,
			*hostconfigGateway,
			*hostconfigDNS,
			*hostconfigInterface,
			*hostconfigURLs,
			*hostconfigID,
			*hostconfigAuth,
		)

	default:
		// Display usage on unknown subcommand
		log.Print(provisionUsage)

		return nil
	}
}

// Check for keygen subcommands.
func keygenArg(args []string, log *logging.Logger) error {
	switch args[subcommandCallPosition] {
	case "certificate":
		// Certificate tool and flags
		certificateCmd := flag.NewFlagSet("keygen", flag.ExitOnError)
		certificateRootCert := certificateCmd.String("rootCert", "", "Root cert in PEM format to sign the new certificate."+
			" Ignored if -isCA is set.")
		certificateRootKey := certificateCmd.String("rootKey", "", "Root key in PEM format to sign the new certificate."+
			" Ignored if -isCA is set.")
		certificateIsCA := certificateCmd.Bool("isCA", false, "Generate self signed root certificate.")
		certificateValidFrom := certificateCmd.String("validFrom", "", "Date formatted as RFC822."+
			" Defaults to time of creation.")
		certificateValidUntil := certificateCmd.String("validUntil", "", "Date formatted as RFC822."+
			" Defaults to time of creation + 72h.")
		certificateCertOut := certificateCmd.String("certOut", "", "Output certificate file."+
			" Defaults to cert.pem or rootcert.pem is -isCA is set.")
		certificateKeyOut := certificateCmd.String("keyOut", "", "Output key file."+
			" Defaults to key.pem or rootkey.pem if -isCA is set.")

		if err := certificateCmd.Parse(args[flagsCallPosition:]); err != nil {
			return err
		}

		return keygen.Certificate(
			&keygen.Args{
				IsCa:         *certificateIsCA,
				RootCertPath: *certificateRootCert,
				RootKeyPath:  *certificateRootKey,
				NotBefore:    *certificateValidFrom,
				NotAfter:     *certificateValidUntil,
				CertOut:      *certificateCertOut,
				KeyOut:       *certificateKeyOut,
			},
		)

	default:
		// Display usage on unknown subcommand
		log.Print(keygenUsage)

		return nil
	}
}
