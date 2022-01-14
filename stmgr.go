// Copyright 2022 the System Transparency Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/system-transparency/stmgr/keygen"
	"github.com/system-transparency/stmgr/provision"
)

const helpText = `
Usage: stmgr <COMMAND> [subcommands...]
	provision:
		Allows creating host configurations by spawning a TUI in
		which the user can input values into that are converted
		into a host_configuration.json file.

	keygen:
		Generate certificates for signing OS packages
		using ED25519 keys.

Use stmgr <COMMAND> -help for more info.
`

func main() {
	log.SetPrefix("stmgr: ")
	log.SetFlags(log.Ltime | log.Lmsgprefix)
	if err := run(os.Args); err != nil {
		log.Printf("ERROR: Runtime error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	// Display helptext if no arguments are given
	if len(args) < 2 {
		log.Print(helpText)
		return nil
	}

	// Evaluate the cli arguments
	switch args[1] {
	case "provision":
		// Provision tool and subcommands
		provisionCmd := flag.NewFlagSet("provision", flag.ContinueOnError)
		efi := provisionCmd.Bool("efi", false, "Store host_configuration.json in the efivarfs")

		provisionCmd.Parse(args[2:])
		return provision.Run(*efi)
	case "keygen":
		// Keygen tool and subcommands
		keygenCmd := flag.NewFlagSet("keygen", flag.ContinueOnError)
		keygenRootCert := keygenCmd.String("rootCert", "", "Root certificate in PEM format to sign the new certificate. Ignored if -isCA is set")
		keygenRootKey := keygenCmd.String("rootKey", "", "Root key in PEM format to sign the new certificate. Ignored if -isCA is set")
		keygenIsCA := keygenCmd.Bool("isCA", false, "Generate self signed root certificate")
		keygenValidFrom := keygenCmd.String("validFrom", "", "Date formatted as RFC822. Defaults to time of creation.")
		keygenValidUntil := keygenCmd.String("validUntil", "", "Date formatted as RFC822. Defaults to time of creation + 72h")
		keygenCertOut := keygenCmd.String("certOut", "", "Output certificate file. Defaults to cert.pem or rootcert.pem is -isCA is set.")
		keygenKeyOut := keygenCmd.String("keyOut", "", "Output key file. Defaults to key.pem or rootkey.pem if -isCA is set.")

		keygenCmd.Parse(args[2:])
		return keygen.Run(*keygenIsCA, *keygenRootCert, *keygenRootKey, *keygenValidFrom, *keygenValidUntil, *keygenCertOut, *keygenKeyOut)
	default:
		// Display helptext on unknown command
		log.Print(helpText)
		return nil
	}
}
