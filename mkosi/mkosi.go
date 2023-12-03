package mkosi

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	RequiredStaticCmds     = "-t uki --bootable --bootloader uki --bios-bootloader none"
	RequiredFedoraPackages = "kernel,systemd-boot-unsigned"
	MinRequiredMkosi       = "mkosi 17"
	MkosiCmd               = "mkosi"
	OSRelease              = "/etc/os-release"
)

// Distro enum
type Distro string

const (
	// Fedora distro
	Fedora  Distro = "fedora"
	Gentoo  Distro = "gentoo"
	Unknown Distro = "unknown"
)

func isMkosiInstalled() bool {
	_, err := exec.LookPath(MkosiCmd)
	if err != nil {
		return false
	}
	cmd := exec.Command(MkosiCmd, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	if string(out) < MinRequiredMkosi {
		return false
	}
	return true
}

func isDistroSupport() (Distro, bool) {
	file, err := os.ReadFile(OSRelease)
	if err != nil {
		return Unknown, false
	}
	delim := strings.Split(string(file), "\n")
	for _, str := range delim {
		if strings.HasPrefix(str, "ID=") {
			switch str[3:] {
			case "fedora":
				return Fedora, true
			case "gentoo":
				return Gentoo, true
			}
			break
		}
	}
	return Unknown, false
}

func RunMkosi(args []string) error {
	if !isMkosiInstalled() {
		return fmt.Errorf("mkosi is not installed")
	}
	strings.Split(RequiredStaticCmds, " ")
	args = append(args, RequiredStaticCmds)
	distro, supported := isDistroSupport()
	if !supported {
		return fmt.Errorf("distro %s is not supported", distro)
	}
	switch distro {
	case Fedora:
		args = append(args, "-p", RequiredFedoraPackages)
	}
	outputDir, err := os.MkdirTemp("/tmp", "mkosi")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}
	args = append(args, "-O", outputDir)
	cmd := exec.Command(MkosiCmd, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkosi failed: %s with command-line %s", string(out), cmd.String())
	}
	fmt.Printf("mkosi succeeded at: %s\n", outputDir)
	return nil
}
