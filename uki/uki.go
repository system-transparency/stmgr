package uki

import (
	"bytes"
	"debug/pe"
	_ "embed" // Needed for go:embed directive
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"system-transparency.org/stboot/stlog"
)

// Adapted from https://github.com/Foxboron/sbctl/blob/master/bundles.go

type UKI struct {
	kernel     string
	initramfs  string
	cmdline    string
	osRelease  string
	sbat       string
	appendSbat bool
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
	cmdlineTmpfile, err := os.CreateTemp("/var/tmp", "cmdline.*")
	if err != nil {
		return fmt.Errorf("failed to make temporary file for cmdline")
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
		osreleaseTmpfile, err := os.CreateTemp("/var/tmp", "os-release.*")
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

func (u *UKI) SetSBAT(sbat string, appendSBAT bool) {
	u.sbat = sbat
	u.appendSbat = appendSBAT
}

func (u *UKI) Cleanup() {
	os.Remove(u.cmdline)
	os.Remove(u.osRelease)
}

//go:embed stub/linuxx64.efi.stub
var embeddedStub string

func getStub(stub string) []byte {
	if stub != "" {
		data, err := os.ReadFile(stub)
		if err == nil {
			return data
		}

		stlog.Info("Failed to read %s as uefi stub: %v", stub, err)
	}

	stlog.Info("Using embedded uefi stub")

	// Implies a copy of the non-mutable string.
	return []byte(embeddedStub)
}

func writeStub(f io.Writer, stub string) error {
	stubFile := getStub(stub)
	if _, err := f.Write(stubFile); err != nil {
		return err
	}

	return nil
}

// Check if binary has an existing SBAT section.
func hasSBAT(stub string) bool {
	out, err := exec.Command("objdump", "-h", stub).Output()
	if errors.Is(err, exec.ErrNotFound) {
		return false
	}

	return bytes.Contains(out, []byte(".sbat"))
}

// Fet the SBAT section from the binary.
func getSBAT(stub string) []byte {
	out, err := exec.Command("objcopy", "--dump-section", ".sbat=/dev/stdout", stub).Output()
	if errors.Is(err, exec.ErrNotFound) {
		return []byte{}
	}

	return out
}

//nolint:varnamelen
func writeOsrelease(f io.Writer) error {
	osrelease := []byte(`NAME="stboot"
PRETTY_NAME="System Transparency Boot Loader"
ID=stboot
BUILD_ID=rolling
`)
	if _, err := f.Write(osrelease); err != nil {
		return err
	}

	return nil
}

//nolint:varnamelen
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

//nolint:cyclop,funlen,gocognit
func generateUKI(uki *UKI, stub, out string) error {
	removeSBAT := false

	// If there is an existing SBAT section, we need to remove it
	if hasSBAT(stub) {
		removeSBAT = true
	}

	// If we want to append the sbat section we need to read the
	// existing section and write both to a file.
	if uki.appendSbat && removeSBAT {
		oldSBAT := getSBAT(stub)

		sbatFile, err := os.CreateTemp("/var/tmp", "stmgr-sbat.*.csv")
		if err != nil {
			return fmt.Errorf("failed to make temporary file for stmgr.csv")
		}

		defer os.Remove(sbatFile.Name())

		suppliedSBAT, err := os.ReadFile(uki.sbat)
		if err != nil {
			return err
		}

		if _, err := sbatFile.Write(suppliedSBAT); err != nil {
			return err
		}

		if _, err := sbatFile.Write(oldSBAT); err != nil {
			return err
		}

		uki.sbat = sbatFile.Name()
	}

	type section struct {
		section  string
		file     string
		optional bool
	}

	sections := []section{
		{".osrel", uki.osRelease, false},
		{".cmdline", uki.cmdline, false},
		{".initrd", uki.initramfs, false},
		{".linux", uki.kernel, false},
		// The stub has a default .sbat section we can use
		{".sbat", uki.sbat, true},
	}

	// Because the sections might overlap we need to figure out the sizes of the
	// different sections
	vma, err := getVMA(stub)
	if err != nil {
		return err
	}

	// -p preserves the dates of the files we are embedding into sections
	args := []string{"-p"}

	//nolint:varnamelen
	for _, s := range sections {
		if s.file == "" {
			if s.optional {
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

		// If the SBAT section is present we need to remove it before adding it
		if removeSBAT {
			args = append(args,
				"--remove-section", ".sbat",
			)
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

		//nolint:errorlint
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("exit code was not 0: %d", exitError.ExitCode())
		}
	}

	return nil
}

//nolint:nlreturn
func roundUpToBlockSize(size uint64) uint64 {
	const blockSize = 4096
	return ((size + blockSize - 1) / blockSize) * blockSize
}
