// Copyright 2022 the System Transparency Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/system-transparency/stmgr/provision"
)

var helpText = `
Usage:
	provision:
		Allows creating host configurations by spawning a TUI in
		which the user can input values into that are converted
		into a host_configuration.json file.
			-efi: Store the output in the efivarfs
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
		log.Println(helpText)
		return nil
	}

	// Provision tool and subcommands
	provisionCmds := flag.NewFlagSet("provision", flag.ContinueOnError)
	efi := provisionCmds.Bool("efi", false, "Store host_configuration.json in the efivarfs")

	switch args[1] {
	case "provision":
		provisionCmds.Parse(args[2:])
		return provision.Run(*efi)
	default:
		log.Println(helpText)
		return nil
	}
}
