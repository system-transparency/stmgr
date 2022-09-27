package main

import (
	"os"

	"github.com/system-transparency/stmgr/eval"
	"github.com/system-transparency/stmgr/log"
	"github.com/system-transparency/stmgr/mkiso"
)

const (
	_ = iota
	commandCallPosition
	subcommandCallPosition
	flagsCallPosition
)

func main() {
	if err := run(os.Args); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

// We move the main program logic away into run()
// so that our code is easier to write tests for.
func run(args []string) error {
	const usage = `Usage: stmgr <COMMAND> <SUBCOMMAND> [flags...]
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

Use 'stmgr <COMMAND> -help' for more info.
`

	// Display helptext if not enough arguments are given
	if len(args) < flagsCallPosition {
		log.Print(usage)

		return nil
	}

	// Check which command is requested or display usage
	switch args[commandCallPosition] {
	case "ospkg":
		return ospkgArg(args)
	case "mkiso":
		return mkisoArg(args)
	case "provision":
		return provisionArg(args)
	case "keygen":
		return keygenArg(args)
	default:
		// Display usage on unknown command
		log.Print(usage)

		return nil
	}
}

// Check for ospkg subcommands.
func ospkgArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "create":
		return eval.OspkgCreate(args[flagsCallPosition:])
	case "sign":
		return eval.OspkgSign(args[flagsCallPosition:])

	default:
		// Display usage on unknown subcommand
		log.Print(`SUBCOMMANDS:
	create:
		Create an OS package from the provided operating
		system files.

	sign:
		Sign the provided OS package with your private key.

Use 'stmgr ospkg <SUBCOMMAND> -help' for more info.
`)

		return nil
	}
}

// Check for provision subcommands.
func provisionArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "hostconfig":
		return eval.ProvisionHostconfig(args[flagsCallPosition:])
	default:
		// Display usage on unknown subcommand
		log.Print(`SUBCOMMANDS:
	hostconfig:
		Allows creating host configurations by spawning a TUI in
		which the user can input values into that are converted
		into a host_configuration.json file.

Use 'stmgr provision <SUBCOMMAND> -help' for more info.
`)

		return nil
	}
}

func mkisoArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "create":
		return mkiso.MkisoCreate(args[flagsCallPosition:])
	default:
		log.Print(`SUBCOMMANDS:
	create:
		create an ISO image with an optional host configuration and EFI binary 
Use 'stmgr mkiso <SUBCOMMAND> -help' for more info.
`)
	}
	return nil
}

// Check for keygen subcommands.
func keygenArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "certificate":
		return eval.KeygenCertificate(args[flagsCallPosition:])
	default:
		// Display usage on unknown subcommand
		log.Print(`SUBCOMMANDS:
	certificate:
		Generate certificates for signing OS packages
		using ED25519 keys.

Use 'stmgr keygen <SUBCOMMAND> -help' for more info.
`)

		return nil
	}
}
