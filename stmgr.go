package main

import (
	"log"
	"os"

	"system-transparency.org/stboot/stlog"
	"system-transparency.org/stmgr/eval"
	"system-transparency.org/stmgr/uki"
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

	ospkg:
		Set of commands related to OS packages. This includes
		creating, signing and analyzing them.

	keygen:
		Commands to generate different keys and certificates for
		system-transparency.

	uki:
		Create a Unified Kernel Image (UKI) for booting stboot and provisioning tools.
		Output formats:
			* ISO

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
	case "ospkg":
		return ospkgArg(args)
	case "uki":
		return ukiArg(args)
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

// Check for ospkg subcommands.
func ospkgArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "create":
		return eval.OspkgCreate(args[flagsCallPosition:])
	case "sign":
		return eval.OspkgSign(args[flagsCallPosition:])
	case "sigsum":
		return eval.OspkgSigsum(args[flagsCallPosition:])
	case "verify":
		return eval.OspkgVerify(args[flagsCallPosition:])
	default:
		// Display usage on unknown subcommand
		log.Print(`SUBCOMMANDS:
	create:
		Create an OS package from the provided operating
		system files.

	sign:
		Sign the provided OS package with your private key.

	sigsum:
		Add a Sigsum proof of logging and corresponding cert
		to the OS package descriptor.

	verify:
		Verify the provided OS package using a Trust policy directory,
		or just a file containing root certificate(s).

Use 'stmgr ospkg <SUBCOMMAND> -help' for more info.
`)

		return nil
	}
}

func ukiArg(args []string) error {
	switch args[subcommandCallPosition] {
	case "create":
		return uki.Create(args[flagsCallPosition:])
	case "to-iso":
		return uki.ToISO(args[flagsCallPosition:])
	default:
		log.Print(`SUBCOMMANDS:
	create:
		Create a unified kernel image with an optional host configuration.

	to-iso:
		Format an already created UKI as a bootale ISO image.

Use 'stmgr uki <SUBCOMMAND> -help' for more info.
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
