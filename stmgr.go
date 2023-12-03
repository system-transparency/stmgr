package main

import (
	"log"
	"os"

	"system-transparency.org/stboot/stlog"
	"system-transparency.org/stmgr/eval"
)

const (
	_ = iota
	commandCallPosition
	subcommandCallPosition
	flagsCallPosition
)

func main() {
	log.SetFlags(0)

	if err := run(os.Args); err != nil {
		stlog.Error("%s", err)
		os.Exit(1)
	}
}

// We move the main program logic away into run()
// so that our code is easier to write tests for.
func run(args []string) error {
	const usage = `Usage: stmgr <COMMAND> <SUBCOMMAND> [flags...]
COMMANDS:
	trustpolicy:
		Manage trust policy files for stboot.

	hostconfig:
		Manage host configuration files for stboot.

	mkosi:
		Set of commands related to mkosi UKI. This includes
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
	case "trustpolicy":
		return trustpolicyArg(args)
	case "hostconfig":
		return hostconfigArg(args)
	case "mkosi":
		return mkosiArg(args)
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

func trustpolicyArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "check":
		return eval.TrustPolicyCheck(args[flagsCallPosition:])
	default:
		log.Print(`SUBCOMMANDS:
	check:
		Create valid trust policy by checking the provided JSON.
		
Use 'stmgr trustpolicy <SUBCOMMAND> -help' for more info.
`)
	}

	return nil
}

func hostconfigArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "check":
		return eval.HostConfigCheck(args[flagsCallPosition:])
	default:
		log.Print(`SUBCOMMANDS:
	check:
		Create valid host configuration by checking the provided JSON.
		
Use 'stmgr hostconfig <SUBCOMMAND> -help' for more info.
`)
	}

	return nil
}

// Check for mkosi subcommands.
func mkosiArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "build":
		return eval.MkosiBuild(args[flagsCallPosition:])
	case "sign":
		return eval.MkosiSign(args[flagsCallPosition:])
	case "verify":
		return eval.MkosiVerify(args[flagsCallPosition:])

	default:
		// Display usage on unknown subcommand
		log.Print(`SUBCOMMANDS:
	build:
		Build an mkosi UKI from the provided operating
		system files.

	sign:
		Sign the provided mkosi UKI with your private key.
	
	verify:
		Verify the provided mkosi UKI with your public key.

Use 'stmgr mkosi <SUBCOMMAND> -help' for more info.
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
