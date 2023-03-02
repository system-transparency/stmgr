package uki

import (
	"debug/pe"
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"system-transparency.org/stmgr/log"
)

// Adapted from https://github.com/Foxboron/sbctl/blob/master/bundles.go

type UKI struct {
	kernel    string `json:"kernel"`
	initramfs string `json:"initramfs"`
	cmdline   string `json:"cmdline"`
	osRelease string `json:"os_release"`
}

func (u *UKI) SetKernel(kernel string) error {
	if kernel == "" {
		return fmt.Errorf("no kernel specified")
	}
	u.kernel = kernel
	return nil
}
func (u *UKI) SetInitramfs(initramfs string) error {
	if initramfs == "" {
		return fmt.Errorf("no initramfs specified")
	}
	u.initramfs = initramfs
	return nil
}

func (u *UKI) SetCmdline(cmdline string) error {
	cmdlineTmpfile, err := os.CreateTemp("", "cmdline.*")
	if err != nil {
		return fmt.Errorf("failed to make temporary file for os-release")
	}
	if _, err := cmdlineTmpfile.Write([]byte(cmdline + "\n")); err != nil {
		return fmt.Errorf("failed to write cmdline tmp file")
	}

	u.cmdline = cmdlineTmpfile.Name()
	return nil
}

func (u *UKI) SetOSRelease(osrelease string) error {
	// uki.OSRelease = "/etc/os-release"
	if osrelease == "" {
		osreleaseTmpfile, err := os.CreateTemp("", "os-release.*")
		if err != nil {
			return fmt.Errorf("failed to make temporary file for os-release")
		}
		if err := writeOsrelease(osreleaseTmpfile); err != nil {
			return fmt.Errorf("failed to write os-release: %w", err)
		}
		u.osRelease = osreleaseTmpfile.Name()
		return nil
	}
	u.osRelease = osrelease
	return nil

}
func (u *UKI) Cleanup() {
	os.Remove(u.cmdline)
	os.Remove(u.osRelease)
}

//go:embed stub/*
var stubs embed.FS

func getStub(stub string) []byte {
	if stub == "" {
		stub = "/usr/lib/systemd/boot/efi/linuxx64.elf.stub"
	}
	if _, err := os.Stat(stub); os.IsExist(err) {
		if b, err := os.ReadFile(stub); err != nil {
			return b
		}
		log.Infof("Failed to read %s as stub: %v", stub, err)
		log.Infof("Using fallback stub")
	}
	f, err := stubs.ReadFile("stub/linuxx64.efi.stub")
	if err != nil {
		return []byte{}
	}
	return f
}

func writeStub(f io.Writer, stub string) error {
	stubFile := getStub(stub)
	if _, err := f.Write(stubFile); err != nil {
		return err
	}
	return nil
}

func writeOsrelease(f io.Writer) error {
	osrelease := []byte(`NAME="stboot"
PRETTY_NAME="System Transparency Boot Loader"
ID=stboot
BUILD_ID=rolling
`)
	if _, err := f.Write(osrelease); err != nil {
		return nil
	}
	return nil
}

func getVMA(stub string) (uint64, error) {
	e, err := pe.Open(stub)
	if err != nil {
		return 0, err
	}
	e.Close()
	s := e.Sections[len(e.Sections)-1]

	vma := uint64(s.VirtualAddress) + uint64(s.VirtualSize)
	switch e := e.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		vma += uint64(e.ImageBase)
	case *pe.OptionalHeader64:
		vma += e.ImageBase
	}
	vma = roundUpToBlockSize(vma)
	return vma, nil
}

func generateUKI(uki *UKI, stub, out string) error {
	type section struct {
		section string
		file    string
	}
	sections := []section{
		{".osrel", uki.osRelease},
		{".cmdline", uki.cmdline},
		{".initrd", uki.initramfs},
		{".linux", uki.kernel},
	}

	// Because the sections might overlap we need to figure out the sizes of the
	// different sections
	vma, err := getVMA(stub)
	if err != nil {
		return err
	}

	var args []string
	for _, s := range sections {
		if s.file == "" {
			// optional sections
			switch s.section {
			case ".splash":
				continue
			}
		}
		fi, err := os.Stat(s.file)
		if err != nil || fi.IsDir() {
			return err
		}
		var flags string
		switch s.section {
		case ".linux":
			flags = "code,readonly"
		default:
			flags = "data,readonly"
		}
		args = append(args,
			"--add-section", fmt.Sprintf("%s=%s", s.section, s.file),
			"--set-section-flags", fmt.Sprintf("%s=%s", s.section, flags),
			"--change-section-vma", fmt.Sprintf("%s=%#x", s.section, vma),
		)
		vma += roundUpToBlockSize(uint64(fi.Size()))
	}

	args = append(args, stub, out)
	cmd := exec.Command("objcopy", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return err
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("exit code was not 0: %d", exitError.ExitCode())
		}
	}
	return nil
}

func roundUpToBlockSize(size uint64) uint64 {
	const blockSize = 4096
	return ((size + blockSize - 1) / blockSize) * blockSize
}
