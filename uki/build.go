package uki

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/diskfs/go-diskfs/filesystem"
)

//nolint:varnamelen
func writeDiskFs(fs filesystem.FileSystem, file, diskPath string) error {
	if path := filepath.Dir(diskPath); path != "/" {
		if err := fs.Mkdir(path); err != nil {
			return err
		}
	}

	//nolint:nosnakecase
	rw, err := fs.OpenFile(diskPath, os.O_CREATE|os.O_RDWR)
	if err != nil {
		return fmt.Errorf("failed to make %s on the disk image", diskPath)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read %s", filepath.Base(file))
	}

	_, err = rw.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write %s", filepath.Base(file))
	}

	return nil
}

//nolint:varnamelen
func createTempFilename() (string, error) {
	// Go only allows us to create templated filenames if we make one then delete
	// it.
	f, err := os.CreateTemp("/var/tmp", "stmgr.*.vfat")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary vfat file: %w", err)
	}

	f.Close()
	os.RemoveAll(f.Name())

	return f.Name(), nil
}

//nolint:funlen,cyclop
func Create(args []string) error {
	ukiCmd := flag.NewFlagSet("uki", flag.ExitOnError)
	out := ukiCmd.String("out", "stmgr", "output path with format as suffix (default: stmgr)")
	initramfs := ukiCmd.String("initramfs", "/tmp/initramfs.linux_amd64.cpio", "initramfs of a system, commonly from u-root")
	cmdline := ukiCmd.String("cmdline", "", "additional cmdline options for the kernel")
	osrelease := ukiCmd.String("osrelease", "", "os-release file for the uki")
	kernel := ukiCmd.String("kernel", "", "kernel or EFI binary to boot")
	force := ukiCmd.Bool("force", false, "remove existing files (default: false)")
	format := ukiCmd.String("format", "iso", "output format iso or uki (default: iso)")
	stub := ukiCmd.String("stub", "", "UKI stub location (defaults to an embedded stub)")
	sbat := ukiCmd.String("sbat", "", "SBAT metadata")
	appendSbat := ukiCmd.Bool("append-sbat", false, "Append SBAT metadata to the existing section (default: false)")

	if err := ukiCmd.Parse(args); err != nil {
		return err
	}

	//nolint:godox
	// TODO: Use slice.Contains when we want generics
	if *format != "iso" && *format != "uki" {
		return fmt.Errorf("format needs to be one of 'iso' or 'uki'")
	}

	outputFile := *out
	if !strings.HasSuffix(outputFile, *format) {
		outputFile = fmt.Sprintf("%s.%s", outputFile, *format)
	}

	if *force {
		os.Remove(outputFile)
	}

	uki := &UKI{}

	if err := uki.SetCmdline(*cmdline); err != nil {
		return fmt.Errorf("failed setting cmdline: %w", err)
	}

	if err := uki.SetKernel(*kernel); err != nil {
		return fmt.Errorf("failed setting kernel: %w", err)
	}

	if err := uki.SetOSRelease(*osrelease); err != nil {
		return fmt.Errorf("failed setting os-release: %w", err)
	}

	if err := uki.SetInitramfs(*initramfs); err != nil {
		return fmt.Errorf("failed setting initramfs")
	}

	// SBAT section is optional
	uki.SetSBAT(*sbat, *appendSbat)

	// Write the stub file to a temporary file
	stubTmpfile, err := os.CreateTemp("", "stub.*.efi")
	if err != nil {
		return fmt.Errorf("failed to make temporary file for stub")
	}

	defer os.Remove(stubTmpfile.Name())

	if err := writeStub(stubTmpfile, *stub); err != nil {
		return fmt.Errorf("failed to write stub to temporary file")
	}

	var ukiFilename string
	if *format == "uki" {
		ukiFilename = outputFile
	} else {
		// File we write for the UKI
		stmgrUkiTmpfile, err := os.CreateTemp("", "stmgr-uki.*.efi")
		if err != nil {
			return fmt.Errorf("failed to make temporary file for the UKI")
		}
		ukiFilename = stmgrUkiTmpfile.Name()
	}

	if err := generateUKI(uki, stubTmpfile.Name(), ukiFilename); err != nil {
		return fmt.Errorf("failed to write UKI: %w", err)
	}

	//nolint:godox
	// TODO: More output formats
	if *format == "iso" {
		// We care about the name, not the file. Create the file, delete it and use it's name
		tmpfilename, err := createTempFilename()
		if err != nil {
			return fmt.Errorf("failed to make temporary filename: %w", err)
		}

		if err := mkvfat(tmpfilename, ukiFilename); err != nil {
			return fmt.Errorf("failed to make vfat partition: %w", err)
		}

		if err := mkiso(outputFile, tmpfilename); err != nil {
			return fmt.Errorf("failed to make iso: %w", err)
		}
	}

	return nil
}
